package types

// These constants are used on the status field of Job type
const (
	JobCreated     = JobStatus("created")
	JobDownloading = JobStatus("downloading")
	JobEncoding    = JobStatus("encoding")
	JobUploading   = JobStatus("uploading")
	JobFinished    = JobStatus("finished")
	JobError       = JobStatus("error")
)

// JobStatus represents the status of a job
type JobStatus string

// Job is the set of parameters of a given job
type Job struct {
	ID               string    `json:"id"`
	Source           string    `json:"source"`
	Destination      string    `json:"destination"`
	Preset           Preset    `json:"preset"`
	Status           JobStatus `json:"status"`
	Details          string    `json:"details"`
	Progress         string    `json:"progress"`
	LocalSource      string    `json:"-"`
	LocalDestination string    `json:"-"`
}

// JobInput stores the information passed from the
// user when creating a job.
type JobInput struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	PresetName  string `json:"preset"`
}
