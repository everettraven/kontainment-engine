package containertools

// Port represents a port to expose on a container
type Port interface {
	HostPort() string
	ContainerPort() string
}

type PortOption func(*port)

func WithHostPort(hostPort string) PortOption {
	return func(p *port) {
		p.host = hostPort
	}
}

func WithContainerPort(containerPort string) PortOption {
	return func(p *port) {
		p.container = containerPort
	}
}

func NewPort(opts ...PortOption) *port {
	p := &port{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

var _ Port = &port{}

// port is a generic implementation of the Port interface
type port struct {
	host      string
	container string
}

func (p *port) HostPort() string { return p.host }

func (p *port) ContainerPort() string { return p.container }
