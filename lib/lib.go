package lib

import (
	"strings"

	"github.com/flavioribeiro/snickers/types"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(jobID string) error

// StartJob starts the job
func StartJob(job types.Job) {
	downloadFunc := GetDownloadFunc(job.Source)
	if err := downloadFunc(job.ID); err != nil {
		return
	}

	if err := FFMPEGEncode(job.ID); err != nil {
		return
	}

	if err := S3Upload(job.ID); err != nil {
		return
	}
}

func GetDownloadFunc(jobSource string) DownloadFunc {
	if strings.Contains(jobSource, "amazonaws") {
		return S3Download
	}

	return HTTPDownload
}
