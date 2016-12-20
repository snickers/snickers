package server

import (
	"fmt"
	"net/http"
)

// HTTPError is a helper to return errors on handlers
func HTTPError(w http.ResponseWriter, httpErr int, msg string, err error) {
	msg := fmt.Sprintf(`{"error": "%s: %s"}`, msg, err.Error())
	http.Error(w, msg, httpErr)
	return
}

// JSONHandler adds json headers
func JSONHandler(actual http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		actual.ServeHTTP(w, r)
	})
}
