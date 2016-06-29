package lib

import (
	"fmt"

	"github.com/flavioribeiro/snickers/db"
)

func ChangeJobStatus(jobID string, newStatus string) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		fmt.Println("ERROR", err.Error())
		//		panic(err)
	}
	job, err := dbInstance.RetrieveJob(jobID)
	job.Status = newStatus
	dbInstance.UpdateJob(job.ID, job)
}

func ChangeJobDetails(jobID string, newDetails string) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		fmt.Println("ERROR", err.Error())
		//panic(err)
	}
	job, err := dbInstance.RetrieveJob(jobID)
	job.Details = newDetails
	dbInstance.UpdateJob(job.ID, job)
}
