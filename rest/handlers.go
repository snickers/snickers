package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
	"github.com/gorilla/mux"
)

func CreatePreset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		fmt.Fprint(w, "error while creating preset")
	}

	var preset types.Preset
	respData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(respData, &preset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = dbInstance.StorePreset(preset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdatePreset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		fmt.Fprint(w, "error while creating database")
		return
	}

	var preset types.Preset
	respData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(respData, &preset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = dbInstance.UpdatePreset(preset.Name, preset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ListPresets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		fmt.Fprint(w, "error while creating database")
		return
	}

	presets, _ := dbInstance.GetPresets()
	result, err := json.Marshal(presets)
	if err != nil {
		fmt.Fprint(w, "error getting presets")
		return
	}

	fmt.Fprintf(w, string(result))
}

func GetPresetDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		fmt.Fprint(w, "error while creating database")
		return
	}

	vars := mux.Vars(r)
	presetName := vars["presetName"]
	preset, err := dbInstance.RetrievePreset(presetName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(preset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s", result)
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create job")
}

func StartJob(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "start job")
}

func ListJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "list jobs")
}

func GetJobDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get job details")
}
