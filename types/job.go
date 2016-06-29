package types

// These constants are used on the status field of Job type
const (
	JobCreated     = "created"
	JobDownloading = "downloading"
	JobEncoding    = "encoding"
	JobUploading   = "uploading"
	JobFinished    = "finished"
	JobError       = "error"
)

// Job is the set of parameters of a given job
type Job struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Preset      Preset `json:"preset"`
	Status      string `json:"status"`
	Details     string `json:"progress"`
}

// JobInput stores the information passed from the
// user when creating a job.
type JobInput struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	PresetName  string `json:"preset"`
}
