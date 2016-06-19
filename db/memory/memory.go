package memory

import "github.com/flavioribeiro/snickers/types"

// Database struct that persists configurations
type Database struct {
	Presets map[string]types.Preset
}

var instance *Database

// GetDatabase returns database singleton
func GetDatabase() (*Database, error) {
	if instance != nil {
		return instance, nil
	}
	instance = &Database{}
	instance.Presets = map[string]types.Preset{}
	return instance, nil
}

//ClearDatabase clears the database
func (r *Database) ClearDatabase() {
	instance.Presets = map[string]types.Preset{}
}

//StorePreset stores preset information
func (r *Database) StorePreset(preset types.Preset) map[string]types.Preset {
	r.Presets[preset.Name] = preset
	return r.Presets
}

//RetrievePreset retrieves one preset from the database
func (r *Database) RetrievePreset(presetName string) types.Preset {
	return r.Presets[presetName]
}

//UpdatPreset updates a preset
func (r *Database) UpdatePreset(presetName string, newPreset types.Preset) types.Preset {
	r.Presets[presetName] = newPreset
	return newPreset
}

//GetPresets retrieves all presets of the database
func (r *Database) GetPresets() []types.Preset {
	res := make([]types.Preset, 0, len(r.Presets))
	for _, value := range r.Presets {
		res = append(res, value)
	}
	return res
}
