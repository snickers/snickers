package main

import (
	"github.com/flavioribeiro/snickers/rest"
	"log"
	"net/http"
)

func main() {
	router := rest.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
