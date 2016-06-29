package lib

import (
	"fmt"
	"strconv"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

// StartJob starts the job
func StartJob(jobID string) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		// we should update the status
		// of the job saying it failed here
	}
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		// we should update the status
		// of the job saying it failed here
	}
	download(job)
}

func download(job types.Job) {
	changeJobStatus(job.ID, types.JobDownloading)
	respch, err := grab.GetAsync(".", job.Source)
	if err != nil {
		// we should update the status
		// of the job saying it failed here
	}
	fmt.Printf("Initializing download...\n")
	resp := <-respch
	for !resp.IsComplete() {
		changeJobProgress(job.ID, strconv.FormatInt(int64(resp.BytesTransferred()*100/resp.Size), 10))
	}
	encode(job)
}

func encode(job types.Job) {
	changeJobStatus(job.ID, types.JobEncoding)
}

func changeJobStatus(jobID string, newStatus string) {
	fmt.Println("Updating Job Status", jobID, newStatus)
	dbInstance, err := db.GetDatabase()
	if err != nil {
		// we should update the status
		// of the job saying it failed here
	}
	job, err := dbInstance.RetrieveJob(jobID)
	job.Status = newStatus
	dbInstance.UpdateJob(job.ID, job)
}

func changeJobProgress(jobID string, newProgress string) {
	fmt.Println("Updating Job Progress", jobID, newProgress)
	dbInstance, err := db.GetDatabase()
	if err != nil {
		// we should update the status
		// of the job saying it failed here
	}
	job, err := dbInstance.RetrieveJob(jobID)
	job.Progress = newProgress
	dbInstance.UpdateJob(job.ID, job)
}
