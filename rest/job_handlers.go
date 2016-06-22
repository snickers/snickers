package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
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
	job.Source = jobInput.Source
	job.Destination = jobInput.Destination
	job.Preset = preset
	_, err = dbInstance.StoreJob(job)
	if err != nil {
		fmt.Println("-> ", err)
		HTTPError(w, http.StatusBadRequest, "storing job", err)
		return
	}
	w.WriteHeader(http.StatusOK)
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
