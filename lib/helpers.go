package lib

import (
	"fmt"

	"github.com/flavioribeiro/snickers/db"
)

func changeJobStatus(jobID string, newStatus string) {
	fmt.Println("Updating Job Status", jobID, newStatus)
	dbInstance, err := db.GetDatabase()
	if err != nil {
		panic(err)
	}
	job, err := dbInstance.RetrieveJob(jobID)
	job.Status = newStatus
	dbInstance.UpdateJob(job.ID, job)
}

func changeJobDetails(jobID string, newDetails string) {
	fmt.Println("Updating Job Details", jobID, newDetails)
	dbInstance, err := db.GetDatabase()
	if err != nil {
		panic(err)
	}
	job, err := dbInstance.RetrieveJob(jobID)
	job.Details = newDetails
	dbInstance.UpdateJob(job.ID, job)
}
