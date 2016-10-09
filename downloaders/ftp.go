package downloaders

import (
	"net/url"
	"os"
	"path"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/secsy/goftp"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/helpers"
	"github.com/snickers/snickers/types"
)

// FTPDownload downloads the file from FTP. Job Source should be
// in format: ftp://login:password@host/path
func FTPDownload(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("ftp-download")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		log.Error("retrieving-job", err)
		return err
	}

	job.LocalSource = helpers.GetLocalSourcePath(job.ID) + path.Base(job.Source)
	job.LocalDestination = helpers.GetLocalDestination(dbInstance, jobID)
	job.Destination = helpers.GetOutputFilename(dbInstance, jobID)
	job.Status = types.JobDownloading
	job.Details = "0%"
	dbInstance.UpdateJob(job.ID, job)

	u, err := url.Parse(job.Source)
	if err != nil {
		return err
	}

	pw, isSet := u.User.Password()
	if !isSet {
		pw = ""
	}

	config := goftp.Config{
		User:               u.User.Username(),
		Password:           pw,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(config, u.Host+":21")
	if err != nil {
		log.Error("dial-config-failed", err)
		return err
	}

	outputFile, err := os.Create(job.LocalSource)
	if err != nil {
		log.Error("creating-local-source-failed", err)
		return err
	}

	err = client.Retrieve(u.Path, outputFile)
	if err != nil {
		log.Error("retrieving-output-failed", err)
		return err
	}

	return nil
}
