package encoders

import (
	"errors"
	"strconv"

	"code.cloudfoundry.org/lager"

	"github.com/3d0c/gmf"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

func addStream(job types.Job, codecName string, oc *gmf.FmtCtx, ist *gmf.Stream) (int, int, error) {
	var codecContext *gmf.CodecCtx
	var ost *gmf.Stream

	codec, err := gmf.FindEncoder(codecName)
	if err != nil {
		return 0, 0, err
	}

	if ost = oc.NewStream(codec); ost == nil {
		return 0, 0, errors.New("unable to create stream in output context")
	}
	defer gmf.Release(ost)

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
		err := setAudioCtxParams(codecContext, ist, job)
		if err != nil {
			return 0, 0, err
		}
	}

	if codecContext.Type() == gmf.AVMEDIA_TYPE_VIDEO {
		err := setVideoCtxParams(codecContext, ist, job)
		if err != nil {
			return 0, 0, err
		}
	}

	if err := codecContext.Open(nil); err != nil {
		return 0, 0, err
	}

	ost.SetCodecCtx(codecContext)

	return ist.Index(), ost.Index(), nil
}

// FFMPEGEncode function is responsible for encoding the file
func FFMPEGEncode(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("ffmpeg-encode")
	log.Info("started", lager.Data{"job": jobID})
	defer log.Info("finished")

	gmf.LogSetLevel(gmf.AV_LOG_FATAL)
	job, _ := dbInstance.RetrieveJob(jobID)
	streamMap := make(map[int]int, 0)
	var lastDelta int64

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
	job.Details = "0%"
	dbInstance.UpdateJob(job.ID, job)

	// add video stream to streamMap
	srcVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		return err
	}
	videoCodec := getVideoCodec(job)

	i, o, err := addStream(job, videoCodec, outputCtx, srcVideoStream)
	if err != nil {
		return err
	}
	streamMap[i] = o

	// add audio stream to streamMap
	srcAudioStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		return err
	}
	audioCodec := getAudioCodec(job)

	i, o, err = addStream(job, audioCodec, outputCtx, srcAudioStream)
	if err != nil {
		return err
	}
	streamMap[i] = o

	if err := outputCtx.WriteHeader(); err != nil {
		return err
	}

	totalFrames := float64(srcVideoStream.NbFrames() + srcAudioStream.NbFrames())
	framesCount := float64(0)

	for packet := range inputCtx.GetNewPackets() {
		ist, err := inputCtx.GetStream(packet.StreamIndex())
		if err != nil {
			return err
		}
		ost, err := outputCtx.GetStream(streamMap[ist.Index()])
		if err != nil {
			return err
		}

		for frame := range packet.Frames(ist.CodecCtx()) {
			if ost.IsAudio() {
				fsTb := gmf.AVR{Num: 1, Den: ist.CodecCtx().SampleRate()}
				outTb := gmf.AVR{Num: 1, Den: ist.CodecCtx().SampleRate()}

				frame.SetPts(packet.Pts())

				pts := gmf.RescaleDelta(ist.TimeBase(), frame.Pts(), fsTb.AVRational(), frame.NbSamples(), &lastDelta, outTb.AVRational())

				frame.SetNbSamples(ost.CodecCtx().FrameSize())
				frame.SetFormat(ost.CodecCtx().SampleFmt())
				frame.SetChannelLayout(ost.CodecCtx().ChannelLayout())
				frame.SetPts(pts)
			} else {
				frame.SetPts(ost.Pts)
			}

			if p, ready, _ := frame.EncodeNewPacket(ost.CodecCtx()); ready {
				if p.Pts() != gmf.AV_NOPTS_VALUE {
					p.SetPts(gmf.RescaleQ(p.Pts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				if p.Dts() != gmf.AV_NOPTS_VALUE {
					p.SetDts(gmf.RescaleQ(p.Dts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				p.SetStreamIndex(ost.Index())

				if err := outputCtx.WritePacket(p); err != nil {
					return err
				}
				gmf.Release(p)
			}

			ost.Pts++
			framesCount++
			percentage := string(strconv.FormatInt(int64(framesCount/totalFrames*100), 10) + "%")
			if percentage != job.Details {
				job.Details = percentage
				dbInstance.UpdateJob(job.ID, job)
			}
		}
		gmf.Release(packet)
	}

	for i := 0; i < outputCtx.StreamsCnt(); i++ {
		ist, err := inputCtx.GetStream(0)
		if err != nil {
			return err
		}
		ost, err := outputCtx.GetStream(streamMap[ist.Index()])
		if err != nil {
			return err
		}

		frame := gmf.NewFrame()

		for {
			if p, ready, _ := frame.FlushNewPacket(ost.CodecCtx()); ready {
				if p.Pts() != gmf.AV_NOPTS_VALUE {
					p.SetPts(gmf.RescaleQ(p.Pts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				if p.Dts() != gmf.AV_NOPTS_VALUE {
					p.SetDts(gmf.RescaleQ(p.Dts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				p.SetStreamIndex(ost.Index())

				if err := outputCtx.WritePacket(p); err != nil {
					return err
				}
				gmf.Release(p)
			} else {
				gmf.Release(p)
				break
			}

			ost.Pts++
		}

		gmf.Release(frame)
	}
	if job.Details != "100%" {
		job.Details = "100%"
		dbInstance.UpdateJob(job.ID, job)
	}

	return nil
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

func GetResolution(job types.Job, inputWidth int, inputHeight int) (int, int) {
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

	width, height := GetResolution(job, ist.CodecCtx().Width(), ist.CodecCtx().Height())

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
