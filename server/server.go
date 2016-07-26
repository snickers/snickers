package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/pivotal-golang/lager"
	"github.com/snickers/snickers"
)

type SnickersServer struct {
	net.Listener
	logger lager.Logger

	listenAddr    string
	listenNetwork string
	router        *Router
	server        *http.Server
	backend       snickers.Backend
}

func New(log lager.Logger, listenNetwork, listenAddr string, be snickers.Backend) *SnickersServer {
	s := &SnickersServer{
		logger:        log,
		listenAddr:    listenAddr,
		listenNetwork: listenNetwork,
		router:        NewRouter(),
		backend:       be,
	}

	// Set up routes
	routes := map[Route]RouterArguments{
		Ping: RouterArguments{Path: Routes[Ping].Path, Method: Routes[Ping].Method, Handler: s.pingHandler},
		//TODO: DRY
		CreateJob:        RouterArguments{Path: Routes[CreateJob].Path, Method: Routes[CreateJob].Method, Handler: s.CreateJob},
		ListJobs:         RouterArguments{Path: Routes[ListJobs].Path, Method: Routes[ListJobs].Method, Handler: s.ListJobs},
		GetJobDetails:    RouterArguments{Path: Routes[GetJobDetails].Path, Method: Routes[GetJobDetails].Method, Handler: s.GetJobDetails},
		StartJob:         RouterArguments{Path: Routes[StartJob].Path, Method: Routes[StartJob].Method, Handler: s.StartJob},
		CreatePreset:     RouterArguments{Path: Routes[CreatePreset].Path, Method: Routes[CreatePreset].Method, Handler: s.CreatePreset},
		UpdatePreset:     RouterArguments{Path: Routes[UpdatePreset].Path, Method: Routes[UpdatePreset].Method, Handler: s.UpdatePreset},
		ListPresets:      RouterArguments{Path: Routes[ListPresets].Path, Method: Routes[ListPresets].Method, Handler: s.StartJob},
		GetPresetDetails: RouterArguments{Path: Routes[GetPresetDetails].Path, Method: Routes[GetPresetDetails].Method, Handler: s.GetPresetDetails},
		DeletePreset:     RouterArguments{Path: Routes[DeletePreset].Path, Method: Routes[DeletePreset].Method, Handler: s.DeletePreset},
	}
	for _, route := range routes {
		s.router.AddHandler(RouterArguments{Path: route.Path, Method: route.Method, Handler: route.Handler})
	}

	s.server = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.router.r.ServeHTTP(w, r)
		}),
	}

	return s
}

func (sn *SnickersServer) Start() error {
	var err error

	sn.logger.Info("snickers-start-server", lager.Data{"listenAddr": sn.listenAddr})
	sn.Listener, err = net.Listen(sn.listenNetwork, sn.listenAddr)
	if err != nil {
		fmt.Println(err)
		sn.logger.Error("snickers-failed-starting-server", err)
		return err
	}

	go sn.server.Serve(sn.Listener)

	return nil
}

func (sn *SnickersServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sn.router.Handler()
}

func (sn *SnickersServer) Stop() error {
	if sn.listenNetwork == "unix" {
		if err := os.Remove(sn.listenAddr); err != nil {
			sn.logger.Info("failed-to-stop-server", lager.Data{"listenAddr": sn.listenAddr})
			return err
		}
	}

	return sn.Listener.Close()
}
