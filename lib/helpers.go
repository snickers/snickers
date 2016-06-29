package lib

import (
	"github.com/flavioribeiro/snickers/db"
)

func ChangeJobStatus(jobID string, newStatus string) {
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	job.Status = newStatus
	dbInstance.UpdateJob(job.ID, job)
}

func ChangeJobDetails(jobID string, newDetails string) {
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	job.Details = newDetails
	dbInstance.UpdateJob(job.ID, job)
}
