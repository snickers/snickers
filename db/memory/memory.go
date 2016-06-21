package memory

import "github.com/flavioribeiro/snickers/types"

// Database struct that persists configurations
type Database struct {
	presets map[string]types.Preset
	jobs    map[string]types.Job
}

var instance *Database

// GetDatabase returns database singleton
func GetDatabase() (*Database, error) {
	if instance != nil {
		return instance, nil
	}
	instance = &Database{}
	instance.presets = map[string]types.Preset{}
	return instance, nil
}

// ClearDatabase clears the database
func (r *Database) ClearDatabase() error {
	instance.presets = map[string]types.Preset{}
	instance.jobs = map[string]types.Job{}
	return nil
}

// StorePreset stores preset information
func (r *Database) StorePreset(preset types.Preset) (map[string]types.Preset, error) {
	r.presets[preset.Name] = preset
	return r.presets, nil
}

// RetrievePreset retrieves one preset from the database
func (r *Database) RetrievePreset(presetName string) (types.Preset, error) {
	return r.presets[presetName], nil
}

// UpdatePreset updates a preset
func (r *Database) UpdatePreset(presetName string, newPreset types.Preset) (types.Preset, error) {
	r.presets[presetName] = newPreset
	return newPreset, nil
}

// GetPresets retrieves all presets of the database
func (r *Database) GetPresets() ([]types.Preset, error) {
	res := make([]types.Preset, 0, len(r.presets))
	for _, value := range r.presets {
		res = append(res, value)
	}
	return res, nil
}

// StoreJob stores job information
func (r *Database) StoreJob(job types.Job) (map[string]types.Job, error) {
	r.jobs[job.ID] = job
	return r.jobs, nil
}

// RetrieveJob retrieves one job from the database
func (r *Database) RetrieveJob(jobID string) (types.Job, error) {
	return r.jobs[jobID], nil
}

// UpdateJob updates a job
func (r *Database) UpdateJob(jobID string, newJob types.Job) (types.Job, error) {
	r.jobs[jobID] = newJob
	return newJob, nil
}

//GetJobs retrieves all jobs of the database
func (r *Database) GetJobs() ([]types.Job, error) {
	res := make([]types.Job, 0, len(r.jobs))
	for _, value := range r.jobs {
		res = append(res, value)
	}
	return res, nil
}
