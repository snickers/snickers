package core

import (
	"net/url"
	"os"
	"time"

	"github.com/secsy/goftp"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// FTPUpload uploades the file using FTP. Job Destination should be
// in format: ftp://login:password@host/path
func FTPUpload(jobID string) error {
	dbInstance, err := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)

	job.Status = types.JobUploading
	dbInstance.UpdateJob(job.ID, job)

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
		return err
	}

	localFile, err := os.Open(job.LocalDestination)
	if err != nil {
		return err
	}

	err = client.Store(u.Path, localFile)
	if err != nil {
		return err
	}

	return nil
}
