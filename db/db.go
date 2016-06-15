package db

// Preset define the set of parameters of a given preset
type Preset struct {
	Name         string      `json:"name,omitempty"`
	Description  string      `json:"description,omitempty"`
	Container    string      `json:"container,omitempty"`
	Profile      string      `json:"profile,omitempty"`
	ProfileLevel string      `json:"profileLevel,omitempty"`
	RateControl  string      `json:"rateControl,omitempty"`
	Video        VideoPreset `json:"video"`
	Audio        AudioPreset `json:"audio"`
}

// VideoPreset define the set of parameters for video on a given preset
type VideoPreset struct {
	Width         string `json:"width,omitempty"`
	Height        string `json:"height,omitempty"`
	Codec         string `json:"codec,omitempty"`
	Bitrate       string `json:"bitrate,omitempty"`
	GopSize       string `json:"gopSize,omitempty"`
	GopMode       string `json:"gopMode,omitempty"`
	InterlaceMode string `json:"interlaceMode,omitempty"`
}

// AudioPreset define the set of parameters for audio on a given preset
type AudioPreset struct {
	Codec   string `json:"codec,omitempty"`
	Bitrate string `json:"bitrate,omitempty"`
}
