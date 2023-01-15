package server

import (
	"net"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/kontainment/engine/api/server/config"
	"github.com/kontainment/engine/api/server/internal/router"
	"github.com/kontainment/engine/api/server/internal/routes/workspace"
	"github.com/kontainment/engine/containertools/runtimes/docker"
)

type KontainmentServer struct {
	config  config.KontainmentConfig
	routers []router.Router
	addr    string
}

// TODO: Add ability to configure the port
func NewKontainmentServer() (*KontainmentServer, error) {
	dockerCli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &KontainmentServer{
		config: config.KontainmentConfig{},
		routers: []router.Router{
			&workspace.WorkspaceRouter{
				ContainerRuntime: docker.NewDockerRuntime(docker.WithDockerClient(dockerCli)),
			},
		},
		addr: "127.0.0.1:8080",
	}, nil
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
