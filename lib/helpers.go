package lib

import (
	"fmt"
	"io"
	"os"

	"github.com/flavioribeiro/gonfig"
)

// GetLogOutput returns the output we want to use
// for http requests log
func GetLogOutput() io.Writer {
	var logOutput io.Writer
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	logfile, _ := cfg.GetString("LOGFILE", "")
	if logfile == "" {
		logOutput = os.Stderr
	} else {
		fmt.Println("Logging requests on", logfile)
		f, err := os.Create(logfile)
		if err != nil {
			panic(err)
		}

		logOutput = f
	}

	return logOutput
}
