package containertools

// ContainerExec represents the response that is received
// when creating an exec process in a container. It is used
// when attaching to the exec process
type ContainerExec interface {
	Id() string
}

var _ ContainerExec = &containerExec{}

// containerExec is a generic implementation of the ContainerExec interface
type containerExec struct {
	id string
}
