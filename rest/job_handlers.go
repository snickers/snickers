package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/lib"
	"github.com/snickers/snickers/types"
)

// CreateJob creates a job
func CreateJob(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var jobInput types.JobInput
	if err := json.NewDecoder(r.Body).Decode(&jobInput); err != nil {
		HTTPError(w, http.StatusBadRequest, "unpacking job", err)
		return
	}

	preset, err := dbInstance.RetrievePreset(jobInput.PresetName)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	//TODO should we move this to lib?
	var job types.Job
	job.ID = uniuri.New()
	job.Source = jobInput.Source
	job.Destination = jobInput.Destination
	job.Preset = preset
	job.Status = types.JobCreated
	_, err = dbInstance.StoreJob(job)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "storing job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}
	fmt.Fprintf(w, "%s", result)
}

// ListJobs lists all jobs
func ListJobs(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	jobs, _ := dbInstance.GetJobs()
	result, err := json.Marshal(jobs)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting jobs", err)
		return
	}

	fmt.Fprintf(w, "%s", string(result))
}

// GetJobDetails returns the details of a given job
func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}
	fmt.Fprintf(w, "%s", result)
}

// StartJob triggers an encoding process
func StartJob(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	go lib.StartJob(job)
}
