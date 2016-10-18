package pipeline

import (
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/downloaders"
	"github.com/snickers/snickers/encoders"
	"github.com/snickers/snickers/types"
	"github.com/snickers/snickers/uploaders"
)

// StartJob starts the job
func StartJob(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, job types.Job) {
	log := logger.Session("start-job", lager.Data{
		"id":          job.ID,
		"status":      job.Status,
		"source":      job.Source,
		"destination": job.Destination,
	})
	defer log.Info("finished")

	log.Info("downloading")
	downloadFunc := downloaders.GetDownloadFunc(job.Source)
	if err := downloadFunc(log, config, dbInstance, job.ID); err != nil {
		log.Error("download failed", err)
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("encoding")
	encodeFunc := encoders.GetEncodeFunc(job)
	if err := encodeFunc(logger, dbInstance, job.ID); err != nil {
		log.Error("encode failed", err)
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("uploading")
	uploadFunc := uploaders.GetUploadFunc(job.Destination)
	if err := uploadFunc(logger, dbInstance, job.ID); err != nil {
		log.Error("upload failed", err)
		job.Status = types.JobError
		job.Details = err.Error()
		dbInstance.UpdateJob(job.ID, job)
		return
	}

	log.Info("erasing temporary files")
	if err := CleanSwap(dbInstance, job.ID); err != nil {
		log.Error("erasing temporary files failed", err)
	}

	job.Status = types.JobFinished
	dbInstance.UpdateJob(job.ID, job)
}

// CleanSwap removes LocalSource and LocalDestination
// files/directories.
func CleanSwap(dbInstance db.Storage, jobID string) error {
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	err = os.RemoveAll(job.LocalSource)
	if err != nil {
		return err
	}

	err = os.RemoveAll(job.LocalDestination)
	return err
}
