package lib

import (
	"strconv"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

// StartJob starts the job
func StartJob(jobID string) {
	download(jobID)
}

func download(jobID string) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		panic(err)
	}

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		panic(err)
	}

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
