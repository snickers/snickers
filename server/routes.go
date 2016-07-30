package server

import (
	"net/http"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/gorilla/mux"

	"github.com/snickers/snickers/core"
)

// Route maps methods to endpoints
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes stores all routes
type Routes []Route

// NewRouter creates a new router for HTTP requests
func NewRouter() *mux.Router {
	log.SetHandler(text.New(core.GetLogOutput()))

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := httplog.New(JSONHandler(route.HandlerFunc))
		router.Methods(route.Method).Path(route.Pattern).Handler(handler)
	}

	return router
}

var routes = Routes{
	// Preset routes
	Route{"POST", "/presets", CreatePreset},
	Route{"PUT", "/presets", UpdatePreset},
	Route{"GET", "/presets", ListPresets},
	Route{"GET", "/presets/{presetName}", GetPresetDetails},
	Route{"DELETE", "/presets/{presetName}", DeletePreset},

	// Job routes
	Route{"POST", "/jobs", CreateJob},
	Route{"GET", "/jobs", ListJobs},
	Route{"GET", "/jobs/{jobID}", GetJobDetails},
	Route{"POST", "/jobs/{jobID}/start", StartJob},
}
