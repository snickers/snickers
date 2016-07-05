package rest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// HTTPError is a helper to return errors on handlers
func HTTPError(w http.ResponseWriter, httpErr int, msg string, err error) {
	w.WriteHeader(httpErr)
	fmt.Fprintf(w, `{"error": "%s: %s"}`, msg, err.Error())
}

// GetLogOutput returns the output we want to use
// for http requests log
func GetLogOutput() io.Writer {
	var logOutput io.Writer
	if os.Getenv("SNICKERS_ENV") == "test" {
		logOutput = ioutil.Discard
	} else {
		logOutput = os.Stdout
	}

	return logOutput
}
