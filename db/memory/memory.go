package memory

import (
	"errors"
	"sync"

	"github.com/snickers/snickers"
)

// Database struct that persists configurations
type Database struct {
	mtx sync.RWMutex

	presets map[string]snickers.Preset
	jobs    map[string]snickers.Job
}

var instance *Database

// GetDatabase returns database singleton
func GetDatabase() (*Database, error) {
	if instance != nil {
		return instance, nil
	}
	instance = &Database{}
	instance.presets = map[string]snickers.Preset{}
	instance.jobs = map[string]snickers.Job{}
	return instance, nil
}

// ClearDatabase clears the database
func (r *Database) ClearDatabase() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	instance.presets = map[string]snickers.Preset{}
	instance.jobs = map[string]snickers.Job{}
	return nil
}

// StorePreset stores preset information
func (r *Database) StorePreset(preset snickers.Preset) (snickers.Preset, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.presets[preset.Name] = preset
	return preset, nil
}

// RetrievePreset retrieves one preset from the database
func (r *Database) RetrievePreset(presetName string) (snickers.Preset, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.presets[presetName]; ok {
		return val, nil
	}
	return snickers.Preset{}, errors.New("preset not found")
}

// UpdatePreset updates a preset
func (r *Database) UpdatePreset(presetName string, newPreset snickers.Preset) (snickers.Preset, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.presets[presetName] = newPreset
	return newPreset, nil
}

// GetPresets retrieves all presets of the database
func (r *Database) GetPresets() ([]snickers.Preset, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	res := make([]snickers.Preset, 0, len(r.presets))
	for _, value := range r.presets {
		res = append(res, value)
	}
	return res, nil
}

// DeletePreset deletes a preset from the database
func (r *Database) DeletePreset(presetName string) (snickers.Preset, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if val, ok := r.presets[presetName]; ok {
		delete(r.presets, presetName)
		return val, nil
	}
	return snickers.Preset{}, errors.New("preset not found")
}

// StoreJob stores job information
func (r *Database) StoreJob(job snickers.Job) (snickers.Job, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.jobs[job.ID] = job
	return job, nil
}

// RetrieveJob retrieves one job from the database
func (r *Database) RetrieveJob(jobID string) (snickers.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.jobs[jobID]; ok {
		return val, nil
	}
	return snickers.Job{}, errors.New("job not found")
}

// UpdateJob updates a job
func (r *Database) UpdateJob(jobID string, newJob snickers.Job) (snickers.Job, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.jobs[jobID] = newJob
	return newJob, nil
}

//GetJobs retrieves all jobs of the database
func (r *Database) GetJobs() ([]snickers.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	res := make([]snickers.Job, 0, len(r.jobs))
	for _, value := range r.jobs {
		res = append(res, value)
	}
	return res, nil
}
