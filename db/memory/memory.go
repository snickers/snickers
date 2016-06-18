package memory

import "github.com/flavioribeiro/snickers/db"

// Database struct that persists configurations
type Database struct {
	Presets map[string]db.Preset
}

var instance *Database

// GetDatabase returns database singleton
func GetDatabase() (*Database, error) {
	if instance != nil {
		return instance, nil
	}
	instance = &Database{}
	instance.Presets = map[string]db.Preset{}
	return instance, nil
}

func ClearDatabase() {
	instance.Presets = map[string]db.Preset{}
}

//StorePreset stores preset information
func (r *Database) StorePreset(preset db.Preset) map[string]db.Preset {
	r.Presets[preset.Name] = preset
	return r.Presets
}

//RetrievePreset retrieves one preset from the database
func (r *Database) RetrievePreset(presetName string) db.Preset {
	return r.Presets[presetName]
}

//GetPresets retrieves all presets of the database
func (r *Database) GetPresets() []db.Preset {
	res := make([]db.Preset, 0, len(r.Presets))
	for _, value := range r.Presets {
		res = append(res, value)
	}
	return res
}
