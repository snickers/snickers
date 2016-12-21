package server

import "net/http"

type Route int

const (
	CreateJob Route = iota
	DeleteJob
	ListJobs
	GetJobDetails
	StartJob
	CreatePreset
	UpdatePreset
	ListPresets
	GetPresetDetails
	DeletePreset
)

var Routes = map[Route]RouterArguments{
	//Job routes
	CreateJob:     RouterArguments{Path: "/jobs", Method: http.MethodPost},
	ListJobs:      RouterArguments{Path: "/jobs", Method: http.MethodGet},
	GetJobDetails: RouterArguments{Path: "/jobs/{jobID}", Method: http.MethodGet},
	DeleteJob:     RouterArguments{Path: "/jobs/{jobID}", Method: http.MethodDelete},
	StartJob:      RouterArguments{Path: "/jobs/{jobID}/start", Method: http.MethodPost},

	//Preset routes
	CreatePreset:     RouterArguments{Path: "/presets", Method: http.MethodPost},
	UpdatePreset:     RouterArguments{Path: "/presets", Method: http.MethodPut},
	ListPresets:      RouterArguments{Path: "/presets", Method: http.MethodGet},
	GetPresetDetails: RouterArguments{Path: "/presets/{presetName}", Method: http.MethodGet},
	DeletePreset:     RouterArguments{Path: "/presets/{presetName}", Method: http.MethodDelete},
}
