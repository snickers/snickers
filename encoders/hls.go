package encoders

import (
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/snickers/hls/segmenter"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// HLSEncode function is responsible for encoding adaptive bitrate outputs
func HLSEncode(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("hls-encode")
	log.Info("started", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	job.Details = "0%"
	job.Status = types.JobEncoding
	dbInstance.UpdateJob(job.ID, job)

	err = encodeInH264(logger, dbInstance, jobID)
	if err != nil {
		return err
	}
	hlsConfig := buildHLSConfig(job)
	err = segmenter.Segment(hlsConfig)
	if err != nil {
		return err
	}
	job.Details = "100%"
	return nil
}

func encodeInH264(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	mp4File, err := ioutil.TempFile("", "snckrs_")
	mp4File.Close()
	mp4Filename := mp4File.Name() + ".mp4"
	os.Create(mp4Filename)
	defer os.Remove(mp4Filename)
	if err != nil {
		return err
	}

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	oldLocalDestination := job.LocalDestination
	job.LocalDestination = mp4Filename
	_, err = dbInstance.UpdateJob(jobID, job)
	if err != nil {
		return err
	}
	err = FFMPEGEncode(logger, dbInstance, jobID)
	if err != nil {
		return err
	}
	job.LocalDestination = oldLocalDestination
	job.LocalSource = mp4Filename
	_, err = dbInstance.UpdateJob(jobID, job)
	if err != nil {
		return err
	}
	return nil
}

func buildHLSConfig(job types.Job) segmenter.HLSConfig {
	return segmenter.HLSConfig{
		SourceFile:      job.LocalSource,
		FileBase:        job.LocalDestination,
		SegmentDuration: 10,
	}
}
