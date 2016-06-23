package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
	"github.com/gorilla/mux"
)

// CreateJob creates a job
func CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var jobInput types.JobInput
	respData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(respData, &jobInput)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "unpacking job", err)
		return
	}

	preset, err := dbInstance.RetrievePreset(jobInput.PresetName)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	var job types.Job
	job.ID = uniuri.New()
	job.Source = jobInput.Source
	job.Destination = jobInput.Destination
	job.Preset = preset
	_, err = dbInstance.StoreJob(job)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "storing job", err)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		HTTPError(w, http.StatusBadRequest, "getting jobs", err)
		return
	}

	fmt.Fprintf(w, string(result))
}

// GetJobDetails returns the details of a given job
func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

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
	fmt.Fprint(w, "start job")
}
