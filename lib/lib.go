package lib

import (
	"github.com/flavioribeiro/snickers/types"
)

// StartJob starts the job
func StartJob(job types.Job) {
	var err error
	err = HTTPDownload(job.ID)
	if err != nil {
		return
	}

	err = FFMPEGEncode(job.ID)
	if err != nil {
		return
	}
}
