package downloaders

import (
	"path"
	"strconv"

	"code.cloudfoundry.org/lager"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/helpers"
	"github.com/snickers/snickers/types"
)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error {
	log := logger.Session("http-download")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		log.Error("retrieving-job", err)
		return err
	}

	localSource, err := helpers.GetLocalSourcePath(config, job.ID)
	if err != nil {
		return err
	}
	job.LocalSource = localSource + path.Base(job.Source)

	job.LocalDestination, err = helpers.GetLocalDestination(config, dbInstance, jobID)
	if err != nil {
		return err
	}

	job.Destination, err = helpers.GetOutputFilename(dbInstance, jobID)
	if err != nil {
		return err
	}

	job.Status = types.JobDownloading
	job.Details = "0%"

	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		log.Error("updating-job", err)
		return err
	}

	respch, err := grab.GetAsync(localSource, job.Source)
	if err != nil {
		return nil
	}

	resp := <-respch
	for !resp.IsComplete() {
		job, err = dbInstance.RetrieveJob(jobID)
		if err != nil {
			return err
		}

		percentage := strconv.FormatInt(int64(resp.BytesTransferred()*100/resp.Size), 10)
		if job.Details != percentage {
			job.Details = percentage + "%"
			dbInstance.UpdateJob(job.ID, job)
		}
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}
