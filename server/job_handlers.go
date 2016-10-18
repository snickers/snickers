package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/snickers/snickers/pipeline"
	"github.com/snickers/snickers/types"
)

// CreateJob creates a job
func (sn *SnickersServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("create-job")
	log.Debug("started")
	defer log.Debug("finished")

	var jobInput types.JobInput
	if err := json.NewDecoder(r.Body).Decode(&jobInput); err != nil {
		log.Error("failed-unpacking-job", err)
		HTTPError(w, http.StatusBadRequest, "unpacking job", err)
		return
	}

	preset, err := sn.db.RetrievePreset(jobInput.PresetName)
	if err != nil {
		log.Error("failed-retrieving-preset", err)
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	var job types.Job

	job.ID = uniuri.New()
	job.Source = jobInput.Source
	job.Destination = jobInput.Destination
	job.Preset = preset
	job.Status = types.JobCreated
	_, err = sn.db.StoreJob(job)
	if err != nil {
		log.Error("failed-storing-job", err)
		HTTPError(w, http.StatusBadRequest, "storing job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		log.Error("failed-packaging-job-data", err)
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", result)
	log.Info("created", lager.Data{"id": job.ID})
}

// ListJobs lists all jobs
func (sn *SnickersServer) ListJobs(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("list-jobs")
	log.Debug("started")
	defer log.Debug("finished")

	jobs, err := sn.db.GetJobs()
	if err != nil {
		log.Error("failed-getting-jobs", err)
		HTTPError(w, http.StatusBadRequest, "getting jobs", err)
		return
	}

	result, err := json.Marshal(jobs)
	if err != nil {
		log.Error("failed-packaging-jobs", err)
		HTTPError(w, http.StatusBadRequest, "packing jobs data", err)
		return
	}

	fmt.Fprintf(w, "%s", string(result))
}

// GetJobDetails returns the details of a given job
func (sn *SnickersServer) GetJobDetails(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("get-job-details")
	log.Debug("started")
	defer log.Debug("finished")

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := sn.db.RetrieveJob(jobID)
	if err != nil {
		log.Error("failed-retrieving-job", err)
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		log.Error("failed-packaging-job-data", err)
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}

	fmt.Fprintf(w, "%s", result)
	log.Info("got-job-details", lager.Data{"id": job.ID})
}

// StartJob triggers an encoding process
func (sn *SnickersServer) StartJob(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("start-job")
	log.Debug("started")
	defer log.Debug("finished")

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := sn.db.RetrieveJob(jobID)
	if err != nil {
		log.Error("failed-retrieving-job", err)
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	log.Debug("starting-job", lager.Data{"id": job.ID})
	w.WriteHeader(http.StatusOK)
	go pipeline.StartJob(log, sn.config, sn.db, job)
}
