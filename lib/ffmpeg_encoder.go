package lib

import (
	"errors"
	"fmt"

	. "github.com/3d0c/gmf"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

func FFMPEGEncode(jobID string) error {
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	return encode(job)
}

func addStream(codecName string, oc *FmtCtx, ist *Stream) (int, int, error) {
	var cc *CodecCtx
	var ost *Stream

	codec, err := FindEncoder(codecName)
	if err != nil {
		return 0, 0, err
	}

	if ost = oc.NewStream(codec); ost == nil {
		fmt.Println("unable to create stream in output context")
	}
	defer Release(ost)

	if cc = NewCodecCtx(codec); cc == nil {
		fmt.Println("unable to create codec context")
	}

	defer Release(cc)

	if oc.IsGlobalHeader() {
		cc.SetFlag(CODEC_FLAG_GLOBAL_HEADER)
	}

	if codec.IsExperimental() {
		cc.SetStrictCompliance(FF_COMPLIANCE_EXPERIMENTAL)
	}

	if cc.Type() == AVMEDIA_TYPE_AUDIO {
		cc.SetSampleFmt(ist.CodecCtx().SampleFmt())
		cc.SetSampleRate(ist.CodecCtx().SampleRate())
		cc.SetChannels(ist.CodecCtx().Channels())
		cc.SelectChannelLayout()
		cc.SelectSampleRate()

	}

	if cc.Type() == AVMEDIA_TYPE_VIDEO {
		cc.SetTimeBase(AVR{1, 25})
		cc.SetProfile(FF_PROFILE_MPEG4_SIMPLE)
		cc.SetDimension(ist.CodecCtx().Width(), ist.CodecCtx().Height())
		cc.SetPixFmt(ist.CodecCtx().PixFmt())
	}

	if err := cc.Open(nil); err != nil {
		fmt.Println(err.Error())
	}

	ost.SetCodecCtx(cc)

	return ist.Index(), ost.Index(), nil
}

// Encode function is responsible for encoding the file
func encode(job types.Job) error {
	dbInstance, _ := db.GetDatabase()
	srcFileName := job.LocalSource
	dstFileName := job.LocalDestination
	stMap := make(map[int]int, 0)
	var lastDelta int64

	inputCtx, err := NewInputCtx(srcFileName)
	if err != nil {
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return err
	}
	defer inputCtx.CloseInputAndRelease()

	outputCtx, err := NewOutputCtx(dstFileName)
	if err != nil {
		return err
	}
	defer outputCtx.CloseOutputAndRelease()

	srcVideoStream, err := inputCtx.GetBestStream(AVMEDIA_TYPE_VIDEO)
	if err != nil {
		return errors.New("No video stream found in " + srcFileName)
	}

	i, o, _ := addStream("mpeg4", outputCtx, srcVideoStream)
	stMap[i] = o

	srcAudioStream, err := inputCtx.GetBestStream(AVMEDIA_TYPE_AUDIO)
	if err != nil {
		return errors.New("No audio stream found in " + srcFileName)
	}

	i, o, _ = addStream("aac", outputCtx, srcAudioStream)
	stMap[i] = o

	if err := outputCtx.WriteHeader(); err != nil {
		return err
	}

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
				fsTb := AVR{1, ist.CodecCtx().SampleRate()}
				outTb := AVR{1, ist.CodecCtx().SampleRate()}

				frame.SetPts(packet.Pts())

				pts := RescaleDelta(ist.TimeBase(), frame.Pts(), fsTb.AVRational(), frame.NbSamples(), &lastDelta, outTb.AVRational())

				frame.
					SetNbSamples(ost.CodecCtx().FrameSize()).
					SetFormat(ost.CodecCtx().SampleFmt()).
					SetChannelLayout(ost.CodecCtx().ChannelLayout()).
					SetPts(pts)
			} else {
				frame.SetPts(ost.Pts)
			}

			if p, ready, _ := frame.EncodeNewPacket(ost.CodecCtx()); ready {
				if p.Pts() != AV_NOPTS_VALUE {
					p.SetPts(RescaleQ(p.Pts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				if p.Dts() != AV_NOPTS_VALUE {
					p.SetDts(RescaleQ(p.Dts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				p.SetStreamIndex(ost.Index())

				if err := outputCtx.WritePacket(p); err != nil {
					return err
				}
				Release(p)
			}

			ost.Pts++
		}
		Release(packet)
	}

	// Flush encoders
	// @todo refactor it (should be a better way)
	for i := 0; i < outputCtx.StreamsCnt(); i++ {
		ist, err := inputCtx.GetStream(0)
		if err != nil {
			return err
		}
		ost, err := outputCtx.GetStream(stMap[ist.Index()])
		if err != nil {
			return err
		}

		frame := NewFrame()

		for {
			if p, ready, _ := frame.FlushNewPacket(ost.CodecCtx()); ready {
				if p.Pts() != AV_NOPTS_VALUE {
					p.SetPts(RescaleQ(p.Pts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				if p.Dts() != AV_NOPTS_VALUE {
					p.SetDts(RescaleQ(p.Dts(), ost.CodecCtx().TimeBase(), ost.TimeBase()))
				}

				p.SetStreamIndex(ost.Index())

				if err := outputCtx.WritePacket(p); err != nil {
					return err
				}
				Release(p)
			} else {
				Release(p)
				break
			}

			ost.Pts++
		}

		Release(frame)
	}

	return nil
}
