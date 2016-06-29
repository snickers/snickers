package lib

import (
	"strconv"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

// StartJob starts the job
func StartJob(job types.Job) {
	download(job.ID)
}

// TODO we should have different "download"
// drivers for different protocols (s3,ftp,http)
func download(jobID string) {
	dbInstance, err := db.GetDatabase()
	job, err := dbInstance.RetrieveJob(jobID)

	changeJobStatus(job.ID, types.JobDownloading)

	respch, err := grab.GetAsync(".", job.Source)
	if err != nil {
		changeJobStatus(job.ID, types.JobError)
		changeJobDetails(job.ID, err.Error())
		return
	}

	resp := <-respch
	for !resp.IsComplete() {
		job, _ = dbInstance.RetrieveJob(jobID)
		percentage := strconv.FormatInt(int64(resp.BytesTransferred()*100/resp.Size), 10)
		if job.Details != percentage {
			changeJobDetails(job.ID, percentage)
		}
	}

	encode(job)
}

func encode(job types.Job) {
	changeJobStatus(job.ID, types.JobEncoding)
	changeJobDetails(job.ID, "0%")
}
