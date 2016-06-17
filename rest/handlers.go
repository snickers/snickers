package rest

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Snickers")
}

func CreatePreset(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create preset")
}

func UpdatePreset(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "update preset")
}

func ListPresets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "list presets")
}

func GetPresetDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get preset details")
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create job")
}

func ListJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "list jobs")
}

func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get job details")
}
