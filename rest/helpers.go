package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/flavioribeiro/gonfig"
)

// HTTPError is a helper to return errors on handlers
func HTTPError(w http.ResponseWriter, httpErr int, msg string, err error) {
	w.WriteHeader(httpErr)
	fmt.Fprintf(w, `{"error": "%s: %s"}`, msg, err.Error())
}

// JSONHandler adds json headers
func JSONHandler(actual http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		actual.ServeHTTP(w, r)
	})
}

// GetLogOutput returns the output we want to use
// for http requests log
func GetLogOutput() io.Writer {
	var logOutput io.Writer
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	logfile, _ := cfg.GetString("LOGFILE", "")
	if logfile == "" {
		logOutput = ioutil.Discard
	} else {
		logOutput = os.Stdout
	}

	return logOutput
}
