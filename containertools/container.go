package containertools

// Container represents a container object
type Container interface {
	Id() string
	Name() string
	Image() Image
	Ports() []Port
	Volumes() []Volume
}

type ContainerOption func(*container)

func WithName(name string) ContainerOption {
	return func(c *container) {
		c.name = name
	}
}

func WithImage(image Image) ContainerOption {
	return func(c *container) {
		c.image = image
	}
}

func WithPorts(ports []Port) ContainerOption {
	return func(c *container) {
		c.ports = ports
	}
}

func WithVolumes(volumes []Volume) ContainerOption {
	return func(c *container) {
		c.volumes = volumes
	}
}

func WithContainerId(id string) ContainerOption {
	return func(c *container) {
		c.id = id
	}
}

func NewContainer(opts ...ContainerOption) Container {
	container := &container{}
	for _, opt := range opts {
		opt(container)
	}

	return container
}

var _ Container = &container{}

// container is a generic implementation of the Container interface
type container struct {
	id      string
	name    string
	image   Image
	ports   []Port
	volumes []Volume
}

func (c *container) Id() string { return c.id }

func (c *container) Name() string { return c.name }

func (c *container) Image() Image { return c.image }

func (c *container) Ports() []Port { return c.ports }

func (c *container) Volumes() []Volume { return c.volumes }
