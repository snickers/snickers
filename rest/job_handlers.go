package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flavioribeiro/snickers/db"
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	jobs, _ := dbInstance.GetJobs()
	result, err := json.Marshal(jobs)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting presets", err)
		return
	}

	fmt.Fprintf(w, string(result))
}

// GetJobDetails returns the details of a given job
func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get job details")
}
