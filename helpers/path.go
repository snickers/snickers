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
func GetLocalSourcePath(configPath string, jobID string) string {
	sourceDir := getBaseDir(configPath, jobID) + "/src/"
	os.MkdirAll(sourceDir, 0777)

	return sourceDir
}

// GetLocalDestination builds the path and filename
// of the local destination file
func GetLocalDestination(configPath string, dbInstance db.Storage, jobID string) string {
	destinationDir := getBaseDir(configPath, jobID) + "/dst/"
	os.MkdirAll(destinationDir, 0777)
	return destinationDir + GetOutputFilename(dbInstance, jobID)
}

// GetOutputFilename build the destination path with
// the output filename
func GetOutputFilename(dbInstance db.Storage, jobID string) string {
	job, _ := dbInstance.RetrieveJob(jobID)
	return strings.Split(path.Base(job.Source), ".")[0] + "_" + job.Preset.Name + "." + job.Preset.Container
}

func getBaseDir(configPath string, jobID string) string {
	cfg, err := gonfig.FromJsonFile(configPath)
	if err != nil {
		return "Error loading configuration File"
	}
	swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")

	return swapDir + jobID
}
