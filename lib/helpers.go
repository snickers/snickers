package lib

import (
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

func ChangeJobStatus(jobID string, newStatus types.JobStatus) {
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
