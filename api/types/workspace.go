package types

// Workspace represents a Kontainment workspace.
// It contains all the information needed for a
// workspace to be started.
type Workspace struct {
	// Name is the name of the workspace. This name
	// should be unique for each workspace created
	Name string
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

// NewWorkspace creates a new Workspace with the provided options
func NewWorkspace(opts ...WorkspaceOptions) *Workspace {
	workspace := &Workspace{}

	for _, opt := range opts {
		opt(workspace)
	}

	return workspace
}
