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
	if os.Getenv("SNICKERS_ENV") == "test" {
		logOutput = ioutil.Discard
	} else {
		logOutput = os.Stdout
	}

	return logOutput
}
