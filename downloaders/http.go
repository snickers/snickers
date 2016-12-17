package downloaders

import (
	"strconv"

	"code.cloudfoundry.org/lager"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error {
	log := logger.Session("http-download")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	respch, err := grab.GetAsync(job.LocalSource, job.Source)
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
		if job.Progress != percentage {
			job.Progress = percentage + "%"
			dbInstance.UpdateJob(job.ID, job)
		}
	}

	if resp.Error != nil {
		return resp.Error
	}

	return nil
}
