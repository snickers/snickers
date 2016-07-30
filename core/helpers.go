package core

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
)

// GetLogOutput returns the output we want to use
// for logging.
func GetLogOutput() io.Writer {
	var logOutput io.Writer
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	logfile, _ := cfg.GetString("LOGFILE", "")
	if logfile == "" {
		logOutput = os.Stderr
	} else {
		f, err := os.Create(logfile)
		if err != nil {
			panic(err)
		}

		logOutput = f
	}

	return logOutput
}

// GetLocalSourcePath builds the path and filename for
// the local source file
func GetLocalSourcePath(jobID string) string {
	sourceDir := getBaseDir(jobID) + "/src/"
	os.MkdirAll(sourceDir, 0777)

	return sourceDir
}

// GetLocalDestination builds the path and filename
// of the local destination file
func GetLocalDestination(jobID string) string {
	destinationDir := getBaseDir(jobID) + "/dst/"
	os.MkdirAll(destinationDir, 0777)
	return destinationDir + GetOutputFilename(jobID)

}

// GetOutputFilename build the destination path with
// the output filename
func GetOutputFilename(jobID string) string {
	dbInstance, _ := db.GetDatabase()
	job, _ := dbInstance.RetrieveJob(jobID)
	return strings.Split(path.Base(job.Source), ".")[0] + "_" + job.Preset.Name + "." + job.Preset.Container
}

func getBaseDir(jobID string) string {
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")

	return swapDir + jobID
}
