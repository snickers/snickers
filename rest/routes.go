package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	logOutput := GetLogOutput()
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler := handlers.LoggingHandler(logOutput, route.HandlerFunc)
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

	// Job routes
	Route{"POST", "/jobs", CreateJob},
	Route{"POST", "/jobs/{jobID}/start", StartJob},
	Route{"GET", "/jobs", ListJobs},
	Route{"GET", "/jobs/{jobID}", GetJobDetails},
}
