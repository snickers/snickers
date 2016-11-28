package encoders

import (
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

	hlsConfig := buildHLSConfig(job)
	return segmenter.Segment(hlsConfig)
}

func buildHLSConfig(job types.Job) segmenter.HLSConfig {
	return segmenter.HLSConfig{
		SourceFile:      job.LocalSource,
		FileBase:        job.LocalDestination,
		SegmentDuration: 10,
	}
}
