package pipeline

import (
	"net/url"
	"os"
	"path"

	"code.cloudfoundry.org/lager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/downloaders"
	"github.com/snickers/snickers/encoders"
	"github.com/snickers/snickers/helpers"
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

	log.Info("setup")
	job, err := SetupJob(job.ID, dbInstance, config)
	if err != nil {
		log.Error("setup-job failed", err)
		return
	}

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

// SetupJob is responsible for set the initial state for a given
// job before starting. It sets local source and destination
// paths and the final destination as well.
func SetupJob(jobID string, dbInstance db.Storage, config gonfig.Gonfig) (types.Job, error) {
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return types.Job{}, err
	}

	localSource, err := helpers.GetLocalSourcePath(config, job.ID)
	if err != nil {
		return types.Job{}, err
	}
	job.LocalSource = localSource + path.Base(job.Source)

	job.LocalDestination, err = helpers.GetLocalDestination(config, dbInstance, jobID)
	if err != nil {
		return types.Job{}, err
	}

	u, err := url.Parse(job.Destination)
	if err != nil {
		return types.Job{}, err
	}
	outputFilename, err := helpers.GetOutputFilename(dbInstance, jobID)
	if err != nil {
		return types.Job{}, err
	}
	u.Path = path.Join(u.Path, outputFilename)
	job.Destination = u.String()

	job.Status = types.JobDownloading
	job.Details = "0%"
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		return types.Job{}, err
	}

	return job, nil
}
