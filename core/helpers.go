package core

import (
	"io"
	"os"

	"github.com/flavioribeiro/gonfig"
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
