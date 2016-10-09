package helpers

import (
	"os"
	"path"
	"strings"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
)

// GetLocalSourcePath builds the path and filename for
// the local source file
func GetLocalSourcePath(jobID string) string {
	sourceDir := getBaseDir(jobID) + "/src/"
	os.MkdirAll(sourceDir, 0777)

	return sourceDir
}

// GetLocalDestination builds the path and filename
// of the local destination file
func GetLocalDestination(dbInstance db.Storage, jobID string) string {
	destinationDir := getBaseDir(jobID) + "/dst/"
	os.MkdirAll(destinationDir, 0777)
	return destinationDir + GetOutputFilename(dbInstance, jobID)
}

// GetOutputFilename build the destination path with
// the output filename
func GetOutputFilename(dbInstance db.Storage, jobID string) string {
	job, _ := dbInstance.RetrieveJob(jobID)
	return strings.Split(path.Base(job.Source), ".")[0] + "_" + job.Preset.Name + "." + job.Preset.Container
}

func getBaseDir(jobID string) string {
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")

	return swapDir + jobID
}
