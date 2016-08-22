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
	Handler http.HandlerFunc
	Path    string
	Method  string
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
	path := fmt.Sprintf("/%s", strings.Trim(args.Path, "/"))
	router.r.Methods(args.Method).Path(fmt.Sprintf("%s", path)).HandlerFunc(JSONHandler(args.Handler))
}
