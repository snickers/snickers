package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"github.com/snickers/snickers/types"
)

// CreatePreset creates a preset
func (sn *SnickersServer) CreatePreset(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("create-preset")
	log.Debug("started")
	defer log.Debug("finished")

	var preset types.Preset
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
		log.Error("failed-unpacking-preset", err)
		HTTPError(w, http.StatusBadRequest, "unpacking preset", err)
		return
	}

	_, err := sn.db.StorePreset(preset)
	if err != nil {
		log.Error("failed-storing-preset", err)
		HTTPError(w, http.StatusBadRequest, "storing preset", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	result, err := json.Marshal(preset)
	log.Info("preset-created", lager.Data{"preset-name": preset.Name})
	fmt.Fprintf(w, "%s", result)
}

// UpdatePreset updates a preset
func (sn *SnickersServer) UpdatePreset(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("update-preset")
	log.Debug("started")
	defer log.Debug("finished")

	var preset types.Preset
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
		log.Error("failed-unpacking-preset", err)
		HTTPError(w, http.StatusBadRequest, "unpacking preset", err)
		return
	}

	_, err := sn.db.RetrievePreset(preset.Name)
	if err != nil {
		log.Error("failed-retrieving-preset", err)
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	_, err = sn.db.UpdatePreset(preset.Name, preset)
	if err != nil {
		log.Error("failed-updating-preset", err)
		HTTPError(w, http.StatusBadRequest, "updating preset", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	result, err := json.Marshal(preset)
	fmt.Fprintf(w, "%s", result)
	log.Info("preset-updated", lager.Data{"preset-name": preset.Name})
}

// ListPresets list all presets available
func (sn *SnickersServer) ListPresets(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("list-presets")
	log.Debug("started")
	defer log.Debug("finished")

	presets, err := sn.db.GetPresets()
	if err != nil {
		log.Error("failed-getting-preset", err)
		HTTPError(w, http.StatusBadRequest, "getting presets", err)
		return
	}
	result, err := json.Marshal(presets)
	if err != nil {
		log.Error("failed-parsing-preset", err)
		HTTPError(w, http.StatusBadRequest, "parsing presets", err)
		return
	}

	fmt.Fprintf(w, string(result))
}

// GetPresetDetails returns the details of a given preset
func (sn *SnickersServer) GetPresetDetails(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("get-preset-details")
	log.Debug("started")
	defer log.Debug("finished")

	vars := mux.Vars(r)
	presetName := vars["presetName"]
	preset, err := sn.db.RetrievePreset(presetName)
	if err != nil {
		log.Error("failed-retrieving-preset", err)
		HTTPError(w, http.StatusBadRequest, "retrieving preset", err)
		return
	}

	result, err := json.Marshal(preset)
	if err != nil {
		log.Error("failed-packing-preset-data", err)
		HTTPError(w, http.StatusBadRequest, "packing preset data", err)
		return
	}

	fmt.Fprintf(w, "%s", result)
}

// DeletePreset creates a preset
func (sn *SnickersServer) DeletePreset(w http.ResponseWriter, r *http.Request) {
	log := sn.logger.Session("delete-preset")
	log.Debug("started")
	defer log.Debug("finished")

	vars := mux.Vars(r)
	presetName := vars["presetName"]
	_, err := sn.db.DeletePreset(presetName)
	if err != nil {
		log.Error("failed-deleting-preset", err)
		HTTPError(w, http.StatusBadRequest, "deleting preset", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
