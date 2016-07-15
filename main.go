package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/rest"
)

func main() {
	currentDir, _ := os.Getwd()
	cfg, _ := gonfig.FromJsonFile(currentDir + "/config.json")
	port, _ := cfg.GetString("PORT", "8080")
	fmt.Println("Starting Snickers on port", port)
	log.Fatal(http.ListenAndServe(":"+port, rest.NewRouter()))
}
