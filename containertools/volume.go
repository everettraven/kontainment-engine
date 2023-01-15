package containertools

// Volume represents a volume to attach for a container
type Volume interface {
	HostPath() string
	MountPath() string
}

type VolumeOption func(*volume)

func WithHostPath(path string) VolumeOption {
	return func(v *volume) {
		v.host = path
	}
}

func WithContainerPath(path string) VolumeOption {
	return func(v *volume) {
		v.mount = path
	}
}

func NewVolume(opts ...VolumeOption) Volume {
	vol := &volume{}
	for _, opt := range opts {
		opt(vol)
	}

	return vol
}

var _ Volume = &volume{}

// volume is a generic implementation of the Volume interface
type volume struct {
	host  string
	mount string
}

func (v *volume) HostPath() string  { return v.host }
func (v *volume) MountPath() string { return v.mount }
