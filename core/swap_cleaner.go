package core

import (
	"os"

	"github.com/snickers/snickers/db"
)

// CleanSwap removes LocalSource and LocalDestination
// files/directories.
func CleanSwap(dbInstance db.Storage, jobID string) error {
	job, _ := dbInstance.RetrieveJob(jobID)

	err := os.RemoveAll(job.LocalSource)
	if err != nil {
		return err
	}

	err = os.RemoveAll(job.LocalDestination)
	if err != nil {
		return err
	}

	return nil
}
