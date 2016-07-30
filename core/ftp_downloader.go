package core

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/secsy/goftp"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// FTPDownload downloads the file from FTP. Job Source should be
// in format: ftp://login:password@host/path
func FTPDownload(jobID string) error {
	dbInstance, err := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	job.LocalSource = GetLocalSourcePath(job.ID) + path.Base(job.Source)
	job.LocalDestination = GetLocalDestination(jobID)
	job.Destination = GetOutputFilename(jobID)
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
		return err
	}

	outputFile, err := os.Create(job.LocalSource)
	if err != nil {
		return err
	}

	err = client.Retrieve(u.Path, outputFile)
	if err != nil {
		return err
	}

	return nil
}
