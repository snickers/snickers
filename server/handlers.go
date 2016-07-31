package server

import (
	"fmt"
	"net/http"
)

func (server *SnickersServer) pingHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, `{"ping":"pong"}`)
}
