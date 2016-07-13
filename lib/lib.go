package lib

import (
	"strings"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(jobID string) error

// StartJob starts the job
func StartJob(job types.Job) {
	dbInstance, _ := db.GetDatabase()

	downloadFunc := GetDownloadFunc(job.Source)
	if err := downloadFunc(job.ID); err != nil {
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	if err := FFMPEGEncode(job.ID); err != nil {
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	if err := S3Upload(job.ID); err != nil {
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
