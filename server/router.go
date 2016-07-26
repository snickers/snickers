package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Router struct {
	r        *mux.Router
	notFound http.Handler
	sub      map[string]*mux.Router
}

type RouterArguments struct {
	Handler    http.HandlerFunc
	Path       string
	PathPrefix string
	Method     string
}

func NewRouter() *Router {
	return &Router{
		r:   mux.NewRouter(),
		sub: make(map[string]*mux.Router),
	}
}

func (router *Router) Handler() http.Handler {
	return router.r
}

func (router *Router) AddHandler(args RouterArguments) {
	var r *mux.Router

	if sub, ok := router.sub[args.PathPrefix]; ok {
		r = sub
	} else {
		r = router.r
	}

	var prefix, path string
	if args.PathPrefix != "" {
		prefix = fmt.Sprintf("/%s", strings.Trim(args.PathPrefix, "/"))
	}
	path = fmt.Sprintf("/%s", strings.Trim(args.Path, "/"))
	r.Methods(args.Method).Path(fmt.Sprintf("%s%s", prefix, path)).HandlerFunc(args.Handler)
}
