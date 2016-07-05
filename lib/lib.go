package lib

import (
	"os"
	"strconv"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

type nextStep func(types.Job)

// StartJob starts the job
func StartJob(job types.Job) {
	Download(job.ID, encode)
}

// Download function downloads sources using
// http protocol.
//
// TODO we should have different "download"
// drivers for different protocols (s3,ftp,http)
func Download(jobID string, next nextStep) {
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)

	ChangeJobStatus(job.ID, types.JobDownloading)

	respch, _ := grab.GetAsync(os.Getenv("SNICKERS_SWAPDIR"), job.Source)

	resp := <-respch
	for !resp.IsComplete() {
		job, _ = dbInstance.RetrieveJob(jobID)
		percentage := strconv.FormatInt(int64(resp.BytesTransferred()*100/resp.Size), 10)
		if job.Details != percentage {
			ChangeJobDetails(job.ID, percentage)
		}
	}

	if resp.Error != nil {
		ChangeJobStatus(job.ID, types.JobError)
		ChangeJobDetails(job.ID, string(resp.Error.Error()))
		return
	}

	next(job)
}

func encode(job types.Job) {
	ChangeJobStatus(job.ID, types.JobEncoding)
	ChangeJobDetails(job.ID, "0%")
}
