package encoders

import (
	"errors"
	"fmt"
	"strconv"

	"code.cloudfoundry.org/lager"

	"github.com/3d0c/gmf"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// FFMPEGEncode function is responsible for encoding the file
func FFMPEGEncode(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("ffmpeg-encode")
	log.Info("started", lager.Data{"job": jobID})
	defer log.Info("finished")

	gmf.LogSetLevel(gmf.AV_LOG_FATAL)
	job, _ := dbInstance.RetrieveJob(jobID)

	// create input context
	inputCtx, err := gmf.NewInputCtx(job.LocalSource)
	if err != nil {
		log.Error("input-failed", err)
		return err
	}
	defer inputCtx.CloseInputAndRelease()

	// create output context
	outputCtx, err := gmf.NewOutputCtx(job.LocalDestination)
	if err != nil {
		log.Error("output-failed", err)
		return err
	}
	defer outputCtx.CloseOutputAndRelease()

	job.Status = types.JobEncoding
	job.Progress = "0%"
	dbInstance.UpdateJob(job.ID, job)

	//get audio and video stream and the streaMap
	streamMap, srcVideoStream, srcAudioStream, err := getAudioVideoStreamSource(inputCtx, outputCtx, job)
	if err != nil {
		return err
	}
	//calculate total number of frames
	totalFrames := float64(srcVideoStream.NbFrames() + srcAudioStream.NbFrames())
	//process all frames and update the job progress
	err = processAllFramesAndUpdateJobProgress(inputCtx, outputCtx, streamMap, job, dbInstance, totalFrames)
	if err != nil {
		return err
	}

	err = processNewFrames(inputCtx, outputCtx, streamMap)
	if err != nil {
		return err
	}

	if job.Progress != "100%" {
		job.Progress = "100%"
		dbInstance.UpdateJob(job.ID, job)
	}

	return nil
}

func processNewFrames(inputCtx *gmf.FmtCtx, outputCtx *gmf.FmtCtx, streamMap map[int]int) error {
	for i := 0; i < outputCtx.StreamsCnt(); i++ {
		inputStream, err := getStream(inputCtx, 0)
		if err != nil {
			return err
		}
		outputStream, err := getStream(outputCtx, streamMap[inputStream.Index()])
		if err != nil {
			return err
		}

		frame := gmf.NewFrame()

		for {
			if p, ready, _ := frame.FlushNewPacket(outputStream.CodecCtx()); ready {
				configurePacket(p, outputStream, frame)
				if err := outputCtx.WritePacket(p); err != nil {
					return err
				}
				gmf.Release(p)
			} else {
				break
			}
			outputStream.Pts++
		}

		gmf.Release(frame)
	}

	return nil
}

func processAllFramesAndUpdateJobProgress(inputCtx *gmf.FmtCtx, outputCtx *gmf.FmtCtx, streamMap map[int]int, job types.Job, dbInstance db.Storage, totalFrames float64) error {
	var lastDelta int64
	framesCount := float64(0)
	for packet := range inputCtx.GetNewPackets() {
		inputStream, err := getStream(inputCtx, packet.StreamIndex())
		if err != nil {
			return err
		}
		outputStream, err := getStream(outputCtx, streamMap[inputStream.Index()])
		if err != nil {
			return err
		}

		for frame := range packet.Frames(inputStream.CodecCtx()) {
			err := processFrame(inputStream, outputStream, packet, frame, outputCtx, &lastDelta)
			if err != nil {
				return err
			}

			outputStream.Pts++
			framesCount++
			percentage := fmt.Sprintf("%.2f", framesCount/totalFrames*100) + "%"
			if percentage != job.Progress {
				job.Progress = percentage
				dbInstance.UpdateJob(job.ID, job)
			}
		}

		gmf.Release(packet)
	}
	return nil
}

func getStream(context *gmf.FmtCtx, streamIndex int) (*gmf.Stream, error) {
	return context.GetStream(streamIndex)
}

func getAudioVideoStreamSource(inputCtx *gmf.FmtCtx, outputCtx *gmf.FmtCtx, job types.Job) (map[int]int, *gmf.Stream, *gmf.Stream, error) {
	streamMap := make(map[int]int, 0)

	// add video stream to streamMap
	srcVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		return nil, nil, nil, errors.New("unable to get the best video stream inside the input context")
	}
	videoCodec := getVideoCodec(job)
	inputIndex, outputIndex, err := addStream(job, videoCodec, outputCtx, srcVideoStream)
	if err != nil {
		return nil, nil, nil, err
	}
	streamMap[inputIndex] = outputIndex

	// add audio stream to streamMap
	srcAudioStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		return nil, nil, nil, errors.New("unable to get the best audio stream inside the input context")
	}
	audioCodec := getAudioCodec(job)
	inputIndex, outputIndex, err = addStream(job, audioCodec, outputCtx, srcAudioStream)
	if err != nil {
		return nil, nil, nil, err
	}
	streamMap[inputIndex] = outputIndex
	if err := outputCtx.WriteHeader(); err != nil {
		return nil, nil, nil, err
	}

	return streamMap, srcVideoStream, srcAudioStream, nil
}

func configureAudioFrame(packet *gmf.Packet, inputStream *gmf.Stream, outputStream *gmf.Stream, frame *gmf.Frame, lastDelta *int64) {
	fsTb := gmf.AVR{Num: 1, Den: inputStream.CodecCtx().SampleRate()}
	outTb := gmf.AVR{Num: 1, Den: inputStream.CodecCtx().SampleRate()}

	frame.SetPts(packet.Pts())

	pts := gmf.RescaleDelta(inputStream.TimeBase(), frame.Pts(), fsTb.AVRational(), frame.NbSamples(), lastDelta, outTb.AVRational())

	frame.SetNbSamples(outputStream.CodecCtx().FrameSize())
	frame.SetFormat(outputStream.CodecCtx().SampleFmt())
	frame.SetChannelLayout(outputStream.CodecCtx().ChannelLayout())
	frame.SetPts(pts)
}

func configurePacket(packet *gmf.Packet, outputStream *gmf.Stream, frame *gmf.Frame) *gmf.Packet {
	if packet.Pts() != gmf.AV_NOPTS_VALUE {
		packet.SetPts(gmf.RescaleQ(packet.Pts(), outputStream.CodecCtx().TimeBase(), outputStream.TimeBase()))
	}

	if packet.Dts() != gmf.AV_NOPTS_VALUE {
		packet.SetDts(gmf.RescaleQ(packet.Dts(), outputStream.CodecCtx().TimeBase(), outputStream.TimeBase()))
	}

	packet.SetStreamIndex(outputStream.Index())

	return packet
}

func processFrame(inputStream *gmf.Stream, outputStream *gmf.Stream, packet *gmf.Packet, frame *gmf.Frame, outputCtx *gmf.FmtCtx, lastDelta *int64) error {
	if outputStream.IsAudio() {
		configureAudioFrame(packet, inputStream, outputStream, frame, lastDelta)
	} else {
		frame.SetPts(outputStream.Pts)
	}

	if newPacket, ready, _ := frame.EncodeNewPacket(outputStream.CodecCtx()); ready {
		configurePacket(newPacket, outputStream, frame)
		if err := outputCtx.WritePacket(newPacket); err != nil {
			return err
		}
		gmf.Release(newPacket)
	}

	return nil
}

func addStream(job types.Job, codecName string, oc *gmf.FmtCtx, inputStream *gmf.Stream) (int, int, error) {
	var codecContext *gmf.CodecCtx
	var outputStream *gmf.Stream

	codec, err := gmf.FindEncoder(codecName)
	if err != nil {
		return 0, 0, err
	}

	if outputStream = oc.NewStream(codec); outputStream == nil {
		return 0, 0, errors.New("unable to create stream in output context")
	}
	defer gmf.Release(outputStream)

	if codecContext = gmf.NewCodecCtx(codec); codecContext == nil {
		return 0, 0, errors.New("unable to create codec context")
	}
	defer gmf.Release(codecContext)

	// https://ffmpeg.org/pipermail/ffmpeg-devel/2008-January/046900.html
	if oc.IsGlobalHeader() {
		codecContext.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}

	if codec.IsExperimental() {
		codecContext.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	if codecContext.Type() == gmf.AVMEDIA_TYPE_AUDIO {
		err := setAudioCtxParams(codecContext, inputStream, job)
		if err != nil {
			return 0, 0, err
		}
	}

	if codecContext.Type() == gmf.AVMEDIA_TYPE_VIDEO {
		err := setVideoCtxParams(codecContext, inputStream, job)
		if err != nil {
			return 0, 0, err
		}
	}

	if err := codecContext.Open(nil); err != nil {
		return 0, 0, err
	}

	outputStream.SetCodecCtx(codecContext)

	return inputStream.Index(), outputStream.Index(), nil
}

func getProfile(job types.Job) int {
	profiles := map[string]int{
		"baseline": gmf.FF_PROFILE_H264_BASELINE,
		"main":     gmf.FF_PROFILE_H264_MAIN,
		"high":     gmf.FF_PROFILE_H264_HIGH,
	}

	if job.Preset.Video.Profile != "" {
		return profiles[job.Preset.Video.Profile]
	}
	return gmf.FF_PROFILE_H264_MAIN
}

func getVideoCodec(job types.Job) string {
	codecs := map[string]string{
		"h264":   "libx264",
		"vp8":    "libvpx",
		"vp9":    "libvpx-vp9",
		"theora": "libtheora",
		"aac":    "aac",
	}

	if codec, ok := codecs[job.Preset.Video.Codec]; ok {
		return codec
	}
	return "libx264"
}

func getAudioCodec(job types.Job) string {
	codecs := map[string]string{
		"aac":    "aac",
		"vorbis": "vorbis",
	}
	if codec, ok := codecs[job.Preset.Audio.Codec]; ok {
		return codec
	}
	return "aac"
}

func getResolution(job types.Job, inputWidth int, inputHeight int) (int, int) {
	var width, height int
	if job.Preset.Video.Width == "" && job.Preset.Video.Height == "" {
		return inputWidth, inputHeight
	} else if job.Preset.Video.Width == "" {
		height, _ = strconv.Atoi(job.Preset.Video.Height)
		width = (inputWidth * height) / inputHeight
	} else if job.Preset.Video.Height == "" {
		width, _ = strconv.Atoi(job.Preset.Video.Width)
		height = (inputHeight * width) / inputWidth
	} else {
		width, _ = strconv.Atoi(job.Preset.Video.Width)
		height, _ = strconv.Atoi(job.Preset.Video.Height)
	}
	return width, height
}

func setAudioCtxParams(codecContext *gmf.CodecCtx, ist *gmf.Stream, job types.Job) error {
	bitrate, err := strconv.Atoi(job.Preset.Audio.Bitrate)
	if err != nil {
		return err
	}

	codecContext.SetBitRate(bitrate)
	codecContext.SetSampleFmt(ist.CodecCtx().SampleFmt())
	codecContext.SetSampleRate(ist.CodecCtx().SampleRate())
	codecContext.SetChannels(ist.CodecCtx().Channels())
	codecContext.SelectChannelLayout()
	codecContext.SelectSampleRate()
	return nil
}

func setVideoCtxParams(codecContext *gmf.CodecCtx, ist *gmf.Stream, job types.Job) error {
	codecContext.SetTimeBase(gmf.AVR{Num: 1, Den: 25}) // what is this

	if job.Preset.Video.Codec == "h264" {
		profile := getProfile(job)
		codecContext.SetProfile(profile)
	}

	gop, err := strconv.Atoi(job.Preset.Video.GopSize)
	if err != nil {
		return err
	}

	width, height := getResolution(job, ist.CodecCtx().Width(), ist.CodecCtx().Height())

	bitrate, err := strconv.Atoi(job.Preset.Video.Bitrate)
	if err != nil {
		return err
	}

	codecContext.SetDimension(width, height)
	codecContext.SetGopSize(gop)
	codecContext.SetBitRate(bitrate)
	codecContext.SetPixFmt(ist.CodecCtx().PixFmt())

	return nil
}
