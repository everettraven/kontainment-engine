package containertools

// Image represents a container image
type Image interface {
	Repository() string
	Tag() string
	Id() string
}

type ImageOption func(*image)

func WithRepository(repo string) ImageOption {
	return func(i *image) {
		i.repo = repo
	}
}

func WithTag(tag string) ImageOption {
	return func(i *image) {
		i.tag = tag
	}
}

func WithImageId(id string) ImageOption {
	return func(i *image) {
		i.id = id
	}
}

func NewImage(opts ...ImageOption) *image {
	img := &image{}
	for _, opt := range opts {
		opt(img)
	}

	return img
}

var _ Image = &image{}

// image is a generic implementation of the Image interface
type image struct {
	repo string
	tag  string
	id   string
}

func (i *image) Repository() string { return i.repo }
func (i *image) Tag() string        { return i.tag }
func (i *image) Id() string         { return i.id }
