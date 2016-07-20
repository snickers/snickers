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

// JSONHandler adds json headers
func JSONHandler(actual http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		actual.ServeHTTP(w, r)
	})
}
