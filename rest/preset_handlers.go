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
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var preset types.Preset
	respData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(respData, &preset)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "unpacking preset", err)
		return
	}

	_, err = dbInstance.StorePreset(preset)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "storing preset", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdatePreset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var preset types.Preset
	respData, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(respData, &preset)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "unpacking preset", err)
		return
	}

	_, err = dbInstance.UpdatePreset(preset.Name, preset)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "updating preset", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ListPresets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	presets, _ := dbInstance.GetPresets()
	result, err := json.Marshal(presets)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting presets", err)
		return
	}

	fmt.Fprintf(w, string(result))
}

func GetPresetDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	vars := mux.Vars(r)
	presetName := vars["presetName"]
	preset, err := dbInstance.RetrievePreset(presetName)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	result, err := json.Marshal(preset)
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "packing preset data", err)
		return
	}
	fmt.Fprintf(w, "%s", result)
}
