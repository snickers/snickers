package core

import (
	"path"
	"strconv"

	"github.com/cavaliercoder/grab"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(jobID string) error {
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	job.LocalSource = GetLocalSourcePath(job.ID) + path.Base(job.Source)
	job.LocalDestination = GetLocalDestination(jobID)
	job.Destination = GetOutputFilename(jobID)
	job.Status = types.JobDownloading
	job.Details = "0%"
	dbInstance.UpdateJob(job.ID, job)

	respch, _ := grab.GetAsync(GetLocalSourcePath(job.ID), job.Source)

	resp := <-respch
	for !resp.IsComplete() {
		job, _ = dbInstance.RetrieveJob(jobID)
		percentage := strconv.FormatInt(int64(resp.BytesTransferred()*100/resp.Size), 10)
		if job.Details != percentage {
			job.Details = percentage + "%"
			dbInstance.UpdateJob(job.ID, job)
		}
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}
