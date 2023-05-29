package workspace

import (
	"github.com/kontainment/engine/api/server/internal/router"
	"github.com/kontainment/engine/api/types"
	"github.com/kontainment/engine/pkg/plugin"
)

type WorkspaceRouter struct {
	PluginClient plugin.Client
	//TODO: add a mutex to prevent concurrent access to this cache
	WorkspaceCache map[string]*types.Workspace
}

var _ router.Router = &WorkspaceRouter{}

func (wr *WorkspaceRouter) Routes() []router.Route {
	return []router.Route{
		router.NewPostRoute("/kontainment/workspace", wr.postWorkspace),
		router.NewDeleteRoute("/kontainment/workspace", wr.deleteWorkspace),
		router.NewGetRoute("/kontainment/workspaces", wr.listWorkspaces),
		// TODO: Add a new route for updating a workspace
	}
}
