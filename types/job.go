package types

// Job is the set of parameters of a given job
type Job struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Preset      Preset `json:"preset"`
	Status      string `json:"status"`
	Progress    string `json:"progress"`
}

// JobInput stores the information passed from the
// user when creating a job.
type JobInput struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	PresetName  string `json:"preset"`
}
