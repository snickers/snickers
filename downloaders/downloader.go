package downloaders

import (
	"path"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/helpers"
	"github.com/snickers/snickers/types"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error

// GetDownloadFunc returns the download function
// based on the job source.
func GetDownloadFunc(jobSource string) DownloadFunc {
	if strings.Contains(jobSource, "amazonaws") {
		return S3Download
	} else if strings.HasPrefix(jobSource, "ftp://") {
		return FTPDownload
	}

	return HTTPDownload
}

// SetupJob is responsible for set the initial state for a given
// job before starting. It sets local source and destination
// paths and the final destination as well.
func SetupJob(jobID string, dbInstance db.Storage, config gonfig.Gonfig) (types.Job, error) {
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return types.Job{}, err
	}

	localSource, err := helpers.GetLocalSourcePath(config, job.ID)
	if err != nil {
		return types.Job{}, err
	}
	job.LocalSource = localSource + path.Base(job.Source)

	job.LocalDestination, err = helpers.GetLocalDestination(config, dbInstance, jobID)
	if err != nil {
		return types.Job{}, err
	}

	job.Destination, err = helpers.GetOutputFilename(dbInstance, jobID)
	if err != nil {
		return types.Job{}, err
	}

	job.Status = types.JobDownloading
	job.Details = "0%"
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		return types.Job{}, err
	}

	return job, nil
}
