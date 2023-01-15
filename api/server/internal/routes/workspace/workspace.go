package workspace

import (
	"github.com/kontainment/engine/api/server/internal/router"
	"github.com/kontainment/engine/containertools"
)

type WorkspaceRouter struct {
	ContainerRuntime containertools.ContainerRuntime
}

var _ router.Router = &WorkspaceRouter{}

func (wr *WorkspaceRouter) Routes() []router.Route {
	return []router.Route{
		router.NewPostRoute("/kontainment/workspace", wr.postWorkspace),
		router.NewDeleteRoute("/kontainment/workspace", wr.deleteWorkspace),
		router.NewGetRoute("/kontainment/workspaces", wr.listWorkspaces),
	}
}
