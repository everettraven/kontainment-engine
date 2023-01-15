package types

import "github.com/kontainment/engine/containertools"

type image struct {
	Id   string
	Repo string
	Tag  string
}

// Workspace represents a Kontainment workspace.
// It contains all the information needed for a
// workspace to be started.
type Workspace struct {
	// Name is the name of the workspace. This name
	// should be unique for each workspace created
	Name string

	// Image is the image to be used to build the workspace
	Image image
}

// WorkspaceOptions represents a function type to
// configure values of a Workspace
type WorkspaceOptions func(*Workspace)

// WithName sets the name of the Workspace
// to the name provided as a parameter
func WithName(name string) WorkspaceOptions {
	return func(w *Workspace) {
		w.Name = name
	}
}

// WithImage sets the image to use for the Workspace
// to the image that is provided as a parameter
func WithImage(img containertools.Image) WorkspaceOptions {
	return func(w *Workspace) {
		w.Image = image{
			Id:   img.Id(),
			Repo: img.Repository(),
			Tag:  img.Tag(),
		}
	}
}

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
