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
	var cc *gmf.CodecCtx
	var ost *gmf.Stream

	codec, err := gmf.FindEncoder(codecName)
	if err != nil {
		return 0, 0, err
	}

	if ost = oc.NewStream(codec); ost == nil {
		return 0, 0, errors.New("unable to create stream in output context")
	}
	defer gmf.Release(ost)

	if cc = gmf.NewCodecCtx(codec); cc == nil {
		return 0, 0, errors.New("unable to create codec context")
	}
	defer gmf.Release(cc)

	// https://ffmpeg.org/pipermail/ffmpeg-devel/2008-January/046900.html
	if oc.IsGlobalHeader() {
		cc.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}

	if codec.IsExperimental() {
		cc.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	if cc.Type() == gmf.AVMEDIA_TYPE_AUDIO {
		bitrate, err := strconv.Atoi(job.Preset.Audio.Bitrate)
		if err != nil {
			return 0, 0, err
		}

		gop, err := strconv.Atoi(job.Preset.Video.GopSize)
		if err != nil {
			return 0, 0, err
		}

		cc.SetBitRate(bitrate)
		cc.SetGopSize(gop)
		cc.SetSampleFmt(ist.CodecCtx().SampleFmt())
		cc.SetSampleRate(ist.CodecCtx().SampleRate())
		cc.SetChannels(ist.CodecCtx().Channels())
		cc.SelectChannelLayout()
		cc.SelectSampleRate()
	}

	if cc.Type() == gmf.AVMEDIA_TYPE_VIDEO {
		cc.SetTimeBase(gmf.AVR{Num: 1, Den: 25}) // what is this

		if job.Preset.Video.Codec == "h264" {
			profile := GetProfile(job)
			cc.SetProfile(profile)
		}

		width, height := GetResolution(job, ist.CodecCtx().Width(), ist.CodecCtx().Height())

		bitrate, err := strconv.Atoi(job.Preset.Video.Bitrate)
		if err != nil {
			return 0, 0, err
		}

		cc.SetDimension(width, height)
		cc.SetBitRate(bitrate)
		cc.SetPixFmt(ist.CodecCtx().PixFmt())
	}

	if err := cc.Open(nil); err != nil {
		return 0, 0, err
	}

	ost.SetCodecCtx(cc)

	return ist.Index(), ost.Index(), nil
}

// FFMPEGEncode function is responsible for encoding the file
func FFMPEGEncode(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("ffmpeg-encode")
	log.Info("started", lager.Data{"job": jobID})
	defer log.Info("finished")

	gmf.LogSetLevel(gmf.AV_LOG_FATAL)
	job, _ := dbInstance.RetrieveJob(jobID)
	stMap := make(map[int]int, 0)
	var lastDelta int64

	inputCtx, err := gmf.NewInputCtx(job.LocalSource)
	if err != nil {
		log.Error("input-failed", err)
		return err
	}
	defer inputCtx.CloseInputAndRelease()

	outputCtx, err := gmf.NewOutputCtx(job.LocalDestination)
	if err != nil {
		log.Error("output-failed", err)
		return err
	}
	defer outputCtx.CloseOutputAndRelease()

	job.Status = types.JobEncoding
	job.Details = "0%"
	dbInstance.UpdateJob(job.ID, job)

	srcVideoStream, _ := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)

	videoCodec := GetCodec(job)

	log.Info("add-stream-start", lager.Data{"code": videoCodec})
	i, o, err := addStream(job, videoCodec, outputCtx, srcVideoStream)
	if err != nil {
		log.Error("add-stream-failed", err)
		return err
	}
	log.Info("add-stream-finished")
	stMap[i] = o

	srcAudioStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		return err
	}

	audioCodec := "aac" // default codec
	if job.Preset.Audio.Codec != "aac" {
		audioCodec = job.Preset.Audio.Codec
	}

	i, o, err = addStream(job, audioCodec, outputCtx, srcAudioStream)
	if err != nil {
		return err
	}
	stMap[i] = o

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
		ost, err := outputCtx.GetStream(stMap[ist.Index()])
		if err != nil {
			return err
		}

		for frame := range packet.Frames(ist.CodecCtx()) {
			if ost.IsAudio() {
				fsTb := gmf.AVR{Num: 1, Den: ist.CodecCtx().SampleRate()}
				outTb := gmf.AVR{Num: 1, Den: ist.CodecCtx().SampleRate()}

				frame.SetPts(packet.Pts())

				pts := gmf.RescaleDelta(ist.TimeBase(), frame.Pts(), fsTb.AVRational(), frame.NbSamples(), &lastDelta, outTb.AVRational())

				frame.
					SetNbSamples(ost.CodecCtx().FrameSize()).
					SetFormat(ost.CodecCtx().SampleFmt()).
					SetChannelLayout(ost.CodecCtx().ChannelLayout()).
					SetPts(pts)
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
		ost, err := outputCtx.GetStream(stMap[ist.Index()])
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

// GetProfile returns the GMF profile number based on job profile string
func GetProfile(job types.Job) int {
	profiles := map[string]int{
		"baseline": gmf.FF_PROFILE_H264_BASELINE,
		"main":     gmf.FF_PROFILE_H264_MAIN,
		"high":     gmf.FF_PROFILE_H264_HIGH,
	}

	if job.Preset.Video.Profile {
		return profiles[job.Preset.Video.Profile]
	}
	return gmf.FF_PROFILE_H264_MAIN
}

// GetCodec returns the right codec
func GetCodec(job types.Job) string {
	codecs := map[string]int{
		"h264":   "libx264",
		"vp8":    "libvpx",
		"vp9":    "libvpx-vp9",
		"theora": "libtheora",
	}

	if job.Preset.Video.Codec {
		return codecs[job.Preset.Video.Codec]
	}
	return "libx264"
}

// GetResolution calculate the output resolution based on the preset and input source
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
