package db

import (
	"errors"
	"github.com/snickers/snickers/types"
	"sync"
)

// Database struct that persists configurations
type memoryDatabase struct {
	mtx sync.RWMutex

	presets map[string]types.Preset
	jobs    map[string]types.Job
}

var databaseInit sync.Once
var memoryInstance *memoryDatabase

// GetDatabase returns database singleton
func getMemoryDatabase() (Storage, error) {
	databaseInit.Do(func() {
		memoryInstance = &memoryDatabase{}
		memoryInstance.presets = map[string]types.Preset{}
		memoryInstance.jobs = map[string]types.Job{}
	})

	return memoryInstance, nil
}

// ClearDatabase clears the database
func (r *memoryDatabase) ClearDatabase() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	memoryInstance.presets = map[string]types.Preset{}
	memoryInstance.jobs = map[string]types.Job{}
	return nil
}

// StorePreset stores preset information
func (r *memoryDatabase) StorePreset(preset types.Preset) (types.Preset, error) {

	//prevent replacing existing preset
	if _, err := r.RetrievePreset(preset.Name); err == nil {
		return types.Preset{}, errors.New("Error 409: Preset already exists, please update instead.")
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.presets[preset.Name] = preset
	return preset, nil
}

// RetrievePreset retrieves one preset from the database
func (r *memoryDatabase) RetrievePreset(presetName string) (types.Preset, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.presets[presetName]; ok {
		return val, nil
	}
	return types.Preset{}, errors.New("preset not found")
}

// UpdatePreset updates a preset
func (r *memoryDatabase) UpdatePreset(presetName string, newPreset types.Preset) (types.Preset, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.presets[presetName] = newPreset
	return newPreset, nil
}

// GetPresets retrieves all presets of the database
func (r *memoryDatabase) GetPresets() ([]types.Preset, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	res := make([]types.Preset, 0, len(r.presets))
	for _, value := range r.presets {
		res = append(res, value)
	}
	return res, nil
}

// DeletePreset deletes a preset from the database
func (r *memoryDatabase) DeletePreset(presetName string) (types.Preset, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if val, ok := r.presets[presetName]; ok {
		delete(r.presets, presetName)
		return val, nil
	}
	return types.Preset{}, errors.New("preset not found")
}

// StoreJob stores job information
func (r *memoryDatabase) StoreJob(job types.Job) (types.Job, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.jobs[job.ID] = job
	return job, nil
}

// RetrieveJob retrieves one job from the database
func (r *memoryDatabase) RetrieveJob(jobID string) (types.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.jobs[jobID]; ok {
		return val, nil
	}
	return types.Job{}, errors.New("job not found")
}

// UpdateJob updates a job
func (r *memoryDatabase) UpdateJob(jobID string, newJob types.Job) (types.Job, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.jobs[jobID] = newJob
	return newJob, nil
}

//GetJobs retrieves all jobs of the database
func (r *memoryDatabase) GetJobs() ([]types.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	res := make([]types.Job, 0, len(r.jobs))
	for _, value := range r.jobs {
		res = append(res, value)
	}
	return res, nil
}
