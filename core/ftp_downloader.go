package core

import (
	"errors"
	"net/url"

	"github.com/smallfish/ftp"
	"github.com/snickers/snickers/db"
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

	ftp := new(ftp.FTP)
	ftp.Connect(u.Host, 21)

	password, isSet := u.User.Password()
	if !isSet {
		password = ""
	}

	ftp.Login(u.User.Username(), password)
	if ftp.Code == 530 {
		return errors.New("login failure")
	}

	ftp.Pwd()
	ftp.Quit()

	return nil
}
