package lib

import (
	"os"
	"path"
	"strconv"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/gonfig"
	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(jobID string) error {
	cfg, _ := gonfig.FromJsonFile("../config.json")
	swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)

	job.LocalSource = swapDir + path.Base(job.Source)

	outputDir := swapDir + "dest/" + string(job.ID)
	os.MkdirAll(outputDir, 0777)
	job.LocalDestination = outputDir + "/" + path.Base(job.Source)

	job.Status = types.JobDownloading
	dbInstance.UpdateJob(job.ID, job)

	respch, _ := grab.GetAsync(swapDir, job.Source)

	resp := <-respch
	for !resp.IsComplete() {
		job, _ = dbInstance.RetrieveJob(jobID)
		percentage := strconv.FormatInt(int64(resp.BytesTransferred()*100/resp.Size), 10)
		if job.Details != percentage {
			job.Details = percentage
			dbInstance.UpdateJob(job.ID, job)
		}
	}

	if resp.Error != nil {
		job.Status = types.JobError
		job.Details = string(resp.Error.Error())
		dbInstance.UpdateJob(job.ID, job)
		return resp.Error
	}

	return nil
}
