package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kontainment/engine/api/server/config"
	"github.com/kontainment/engine/api/server/internal/router"
	"github.com/kontainment/engine/api/server/internal/routes/workspace"
)

type KontainmentServer struct {
	config  config.KontainmentConfig
	routers []router.Router
	addr    string
}

func NewKontainmentServer() *KontainmentServer {
	return &KontainmentServer{
		config: config.KontainmentConfig{},
		routers: []router.Router{
			&workspace.WorkspaceRouter{},
		},
		addr: "127.0.0.1:8080",
	}
}

func (ks *KontainmentServer) Serve() error {
	gMux := mux.NewRouter()
	for _, router := range ks.routers {
		for _, route := range router.Routes() {
			gMux.Path(route.Pattern()).Methods(route.Method()).HandlerFunc(route.Handler())
		}
	}

	srv := &http.Server{
		Handler: gMux,
	}

	listener, err := net.Listen("tcp", ks.addr)
	if err != nil {
		return err
	}

	return srv.Serve(listener)
}
