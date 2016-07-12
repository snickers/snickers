package lib

import (
	"github.com/flavioribeiro/snickers/types"
)

// StartJob starts the job
func StartJob(job types.Job) {
	if err := HTTPDownload(job.ID); err != nil {
		return
	}

	if err := FFMPEGEncode(job.ID); err != nil {
		return
	}
}
