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
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

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

	outputFile, err := os.Create("/tmp/aeeez.jpg")
	if err != nil {
		return err
	}

	err = client.Retrieve("/video/vows_640.jpg", outputFile)
	if err != nil {
		return err
	}

	return nil
}
