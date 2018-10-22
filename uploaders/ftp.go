package uploaders

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
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
		client.Close()
		return err
	}

	fileInfo, err := os.Stat(job.LocalDestination)
	if err != nil {
		log.Error("get-destination-info", err)
		client.Close()
		return err
	}

	remotePath := "." + path.Dir(u.Path)
	log.Info("check-remote-path", lager.Data{"path": remotePath})

	_, err = client.Stat(remotePath)
	if err != nil {
		err = nil
		_, errMk := client.Mkdir(remotePath)
		if errMk != nil {
			log.Error("no-create-path", errMk)
			client.Close()
			return errMk
		}
	}

	if fileInfo.IsDir() {
		base := path.Base(job.LocalDestination)
		client.Mkdir(remotePath + "/" + base)
		files, err := ioutil.ReadDir(job.LocalDestination)
		if err != nil {
			log.Error("listing-files", err)
			client.Close()
			return err
		}
		for _, file := range files {
			localFile, err := os.Open(job.LocalDestination + "/" + file.Name())
			defer localFile.Close()
			if err != nil {
				log.Error("opening-local-destination-failed", err)
				client.Close()
				return err
			}
			client.Store("."+u.Path+"/"+file.Name(), localFile)

		}

	} else {
		localFile, err := os.Open(job.LocalDestination)
		defer localFile.Close()

		fi, err := localFile.Stat()
		if err != nil {
			return err
		}

		var tracker = NewProgressTracker64(fi.Size(), &job, dbInstance)
		reader := io.TeeReader(localFile, tracker)
		tracker.Start()

		if err != nil {
			log.Error("opening-local-destination-failed", err)
			client.Close()
			return err
		}

		err = client.Store("."+u.Path, reader)
		if err != nil {
			log.Error("storing-file-failed", err)
			tracker.Finish()
			client.Close()
			return err
		}
		tracker.Finish()
	}
	client.Close()
	return err
}
