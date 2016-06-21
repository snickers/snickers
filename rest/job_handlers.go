package rest

import (
	"fmt"
	"net/http"
)

// CreateJob creates a job
func CreateJob(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create job")
}

// StartJob triggers an encoding process
func StartJob(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "start job")
}

// ListJobs lists all jobs
func ListJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "list jobs")
}

// GetJobDetails returns the details of a given job
func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get job details")
}
