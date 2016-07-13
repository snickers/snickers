package lib

import (
	"errors"
	"strconv"

	"github.com/3d0c/gmf"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

func addStream(codecName string, oc *gmf.FmtCtx, ist *gmf.Stream) (int, int, error) {
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

	if oc.IsGlobalHeader() {
		cc.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}

	if codec.IsExperimental() {
		cc.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	if cc.Type() == gmf.AVMEDIA_TYPE_AUDIO {
		cc.SetSampleFmt(ist.CodecCtx().SampleFmt())
		cc.SetSampleRate(ist.CodecCtx().SampleRate())
		cc.SetChannels(ist.CodecCtx().Channels())
		cc.SelectChannelLayout()
		cc.SelectSampleRate()

	}

	if cc.Type() == gmf.AVMEDIA_TYPE_VIDEO {
		cc.SetTimeBase(gmf.AVR{Num: 1, Den: 25})
		cc.SetProfile(gmf.FF_PROFILE_MPEG4_SIMPLE)
		cc.SetDimension(ist.CodecCtx().Width(), ist.CodecCtx().Height())
		cc.SetPixFmt(ist.CodecCtx().PixFmt())
	}

	if err := cc.Open(nil); err != nil {
		return 0, 0, err
	}

	ost.SetCodecCtx(cc)

	return ist.Index(), ost.Index(), nil
}

// FFMPEGEncode function is responsible for encoding the file
func FFMPEGEncode(jobID string) error {
	gmf.LogSetLevel(gmf.AV_LOG_FATAL)
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	srcFileName := job.LocalSource
	dstFileName := job.LocalDestination
	stMap := make(map[int]int, 0)
	var lastDelta int64

	inputCtx, err := gmf.NewInputCtx(srcFileName)
	if err != nil {
		return err
	}
	defer inputCtx.CloseInputAndRelease()

	outputCtx, err := gmf.NewOutputCtx(dstFileName)
	if err != nil {
		return err
	}
	defer outputCtx.CloseOutputAndRelease()

	job.Status = types.JobEncoding
	job.Details = "0%"
	dbInstance.UpdateJob(job.ID, job)

	srcVideoStream, _ := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	i, o, err := addStream("mpeg4", outputCtx, srcVideoStream)
	if err != nil {
		return err
	}
	stMap[i] = o

	srcAudioStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	if err != nil {
		return err
	}

	i, o, err = addStream("aac", outputCtx, srcAudioStream)
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
