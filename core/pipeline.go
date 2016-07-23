package core

import (
	"strings"

	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(jobID string) error

// StartJob starts the job
func StartJob(job types.Job) {
	dbInstance, _ := db.GetDatabase()

	log.SetHandler(text.New(GetLogOutput()))
	ctx := log.WithFields(log.Fields{
		"id":          job.ID,
		"status":      job.Status,
		"source":      job.Source,
		"destination": job.Destination,
	})

	ctx.Info("downloading")
	downloadFunc := GetDownloadFunc(job.Source)
	if err := downloadFunc(job.ID); err != nil {
		ctx.WithError(err).Error("download failed")
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	ctx.Info("encoding")
	if err := FFMPEGEncode(job.ID); err != nil {
		ctx.WithError(err).Error("encode failed")
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	ctx.Info("uploading")
	if err := S3Upload(job.ID); err != nil {
		ctx.WithError(err).Error("upload failed")
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
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
