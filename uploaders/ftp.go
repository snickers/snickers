package uploaders

import (
	"net/url"
	"os"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/secsy/goftp"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// FTPUpload uploades the file using FTP. Job Destination should be
// in format: ftp://login:password@host/path
func FTPUpload(logger lager.Logger, dbInstance db.Storage, jobID string) error {
	log := logger.Session("ftp-upload")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		log.Error("retrieving-job", err)
		return err
	}

	job.Status = types.JobUploading
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		log.Error("updating-job", err)
		return err
	}

	u, err := url.Parse(job.Destination)
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

	localFile, err := os.Open(job.LocalDestination)
	if err != nil {
		log.Error("opening-local-destination-failed", err)
		return err
	}

	err = client.Store("."+u.Path, localFile)
	if err != nil {
		log.Error("storing-file-failed", err)
		return err
	}

	return nil
}
