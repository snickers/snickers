package main

import (
	"github.com/snickers/snickers/rest"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", rest.NewRouter()))
}
