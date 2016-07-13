package lib

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(jobID string) error {
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)

	basePath := swapDir + string(job.ID)

	sourceDir := basePath + "/src/"
	job.LocalSource = sourceDir + path.Base(job.Source)

	outputDir := basePath + "/dst/"
	os.MkdirAll(sourceDir, 0777)
	os.MkdirAll(outputDir, 0777)
	outputFilename := strings.Split(path.Base(job.Source), ".")[0] + "_" + job.Preset.Name + "." + job.Preset.Container
	job.LocalDestination = outputDir + outputFilename
	job.Destination = job.Destination + outputFilename
	job.Status = types.JobDownloading
	job.Details = "0%"
	dbInstance.UpdateJob(job.ID, job)

	respch, _ := grab.GetAsync(basePath+"/src/", job.Source)

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
		return resp.Error
	}

	return nil
}
