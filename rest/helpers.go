package rest

import (
	"fmt"
	"net/http"
)

// HTTPError is a helper to return errors on handlers
func HTTPError(w http.ResponseWriter, httpErr int, msg string, err error) {
	w.WriteHeader(httpErr)
	fmt.Fprintf(w, `{"error": "%s: %s"}`, msg, err.Error())
}
