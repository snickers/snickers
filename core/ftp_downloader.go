package core

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

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

	u, err = url.Parse(job.Source)
	if err != nil {
		return err
	}

	ftp := new(ftp.FTP)
	ftp.Connect(user.Host, 21)

	ftp.Login(u.User.Username(), u.User.Password())
	if ftp.Code == 530 {
		return Error("login failure")
	}

	ftp.Pwd()
	ftp.Quit()

	return nil
}
