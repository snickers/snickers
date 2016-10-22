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
func GetLocalSourcePath(config gonfig.Gonfig, jobID string) (string, error) {
	baseDir, err := getBaseDir(config, jobID)
	sourceDir := baseDir + "/src/"
	if err != nil {
		return "", err
	}

	os.MkdirAll(sourceDir, 0700)

	return sourceDir, nil
}

// GetLocalDestination builds the path and filename
// of the local destination file
func GetLocalDestination(config gonfig.Gonfig, dbInstance db.Storage, jobID string) (string, error) {
	baseDir, err := getBaseDir(config, jobID)
	if err != nil {
		return "", err
	}

	destinationDir := baseDir + "/dst/"
	if err != nil {
		return "", err
	}

	os.MkdirAll(destinationDir, 0700)
	outputFilename, err := GetOutputFilename(dbInstance, jobID)
	if err != nil {
		return "", err
	}

	return destinationDir + outputFilename, nil
}

// GetOutputFilename build the destination path with
// the output filename
func GetOutputFilename(dbInstance db.Storage, jobID string) (string, error) {
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return "", err
	}

	return strings.Split(path.Base(job.Source), ".")[0] + "_" + job.Preset.Name + "." + job.Preset.Container, nil
}

func getBaseDir(config gonfig.Gonfig, jobID string) (string, error) {
	swapDir, err := config.GetString("SWAP_DIRECTORY", "")
	if err != nil {
		return "", err
	}

	return swapDir + jobID, nil
}
