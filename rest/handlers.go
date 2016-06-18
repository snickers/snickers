package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flavioribeiro/snickers/db/memory"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// TODO a great page for API root location
	fmt.Fprint(w, "Snickers")
}

func CreatePreset(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create preset")
}

func UpdatePreset(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "update preset")
}

func ListPresets(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := memory.GetDatabase()
	if err != nil {
		fmt.Fprint(w, "error while creating database")
	}

	var result []string
	presets := dbInstance.GetPresets()

	for _, preset := range presets {
		presetJson, err := json.Marshal(preset)
		if err != nil {
			fmt.Fprint(w, "error while marshaling preset")
		}
		result = append(result, string(presetJson))
	}

	fmt.Fprintf(w, "%s", result)
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
