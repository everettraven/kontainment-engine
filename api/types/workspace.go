package types

import (
	"github.com/kontainment/devcontainers-go/pkg/devcontainers"
)

type image struct {
	Id   string
	Repo string
	Tag  string
}

// Workspace represents a Kontainment workspace.
// It contains all the information needed for a
// workspace to be started.
type Workspace struct {
	DevContainer devcontainers.DevContainer
	Plugin       string
}

// WorkspaceOptions represents a function type to
// configure values of a Workspace
type WorkspaceOptions func(*Workspace)

// NewWorkspace creates a new Workspace with the provided options
func NewWorkspace(opts ...WorkspaceOptions) *Workspace {
	workspace := &Workspace{}

	for _, opt := range opts {
		opt(workspace)
	}

	return workspace
}

type WorkspaceList struct {
	Workspaces []Workspace
}

type ApiError struct {
	Msg string
}

func NewApiError(msg string) *ApiError {
	return &ApiError{
		Msg: msg,
	}
}

func (ae *ApiError) Error() string {
	return ae.Msg
}
