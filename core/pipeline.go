package core

import (
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(logger lager.Logger, dbInstance db.Storage, jobID string) error

// StartJob starts the job
func StartJob(logger lager.Logger, dbInstance db.Storage, job types.Job) {
	log := logger.Session("start-job", lager.Data{
		"id":          job.ID,
		"status":      job.Status,
		"source":      job.Source,
		"destination": job.Destination,
	})
	defer log.Info("finished")

	log.Info("downloading")
	downloadFunc := GetDownloadFunc(job.Source)
	if err := downloadFunc(log, dbInstance, job.ID); err != nil {
		log.Error("download failed", err)
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("encoding")
	if err := FFMPEGEncode(logger, dbInstance, job.ID); err != nil {
		log.Error("encode failed", err)
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("uploading")
	if err := S3Upload(logger, dbInstance, job.ID); err != nil {
		log.Error("upload failed", err)
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("erasing temporary files")
	if err := CleanSwap(dbInstance, job.ID); err != nil {
		log.Error("erasing temporary files failed", err)
	}

	job.Status = types.JobFinished
	dbInstance.UpdateJob(job.ID, job)
}

// GetDownloadFunc returns the download function
// based on the job source.
func GetDownloadFunc(jobSource string) DownloadFunc {
	if strings.Contains(jobSource, "amazonaws") {
		return S3Download
	}

	return HTTPDownload
}
