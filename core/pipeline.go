package core

import (
	"strings"

	"github.com/pivotal-golang/lager"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(jobID string) error

// StartJob starts the job
func StartJob(job snickers.Job) {
	//TODO: replace this to use the one initialized on the server
	log := lager.NewLogger("snickers")
	log.Session("start-job", lager.Data{
		"id": job.ID,
	})
	dbInstance, _ := db.GetDatabase()

	log.Info("starting", lager.Data{
		"status":      job.Status,
		"source":      job.Source,
		"destination": job.Destination,
	})

	log.Info("downloading")
	downloadFunc := GetDownloadFunc(job.Source)
	if err := downloadFunc(job.ID); err != nil {
		log.Error("download failed", err)
		job.Status = snickers.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("encoding")
	if err := FFMPEGEncode(job.ID); err != nil {
		log.Error("encode failed", err)
		job.Status = snickers.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("uploading")
	if err := S3Upload(job.ID); err != nil {
		log.Error("upload failed", err)
		job.Status = snickers.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("erasing temporary files")
	if err := CleanSwap(job.ID); err != nil {
		log.Error("erasing temporary files failed", err)
	}

	job.Status = snickers.JobFinished
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
