package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flavioribeiro/snickers/db"
	"github.com/flavioribeiro/snickers/types"
	"github.com/gorilla/mux"
)

// CreatePreset creates a preset
func CreatePreset(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var preset types.Preset
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
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

// UpdatePreset updates a preset
func UpdatePreset(w http.ResponseWriter, r *http.Request) {
	dbInstance, err := db.GetDatabase()
	if err != nil {
		HTTPError(w, http.StatusBadRequest, "getting database", err)
		return
	}

	var preset types.Preset
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
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

// ListPresets list all presets available
func ListPresets(w http.ResponseWriter, r *http.Request) {
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

// GetPresetDetails returns the details of a given preset
func GetPresetDetails(w http.ResponseWriter, r *http.Request) {
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
