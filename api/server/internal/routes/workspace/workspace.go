package workspace

import "github.com/kontainment/engine/api/server/internal/router"

type WorkspaceRouter struct{}

var _ router.Router = &WorkspaceRouter{}

func (wr *WorkspaceRouter) Routes() []router.Route {
	return []router.Route{
		router.NewPostRoute("/kontainment/workspace", wr.postWorkspace),
	}
}
