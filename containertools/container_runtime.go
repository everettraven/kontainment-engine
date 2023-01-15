package containertools

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
)

//TODO: Move each interface definition and implementation
// to it's own file

// ContainerRuntime is an interface that is meant to convey all the functionality
// that the Kontainment Engine will need from any given container runtime. This
// will help Kontainment be configurable and support more than one container
// runtime in the future
type ContainerRuntime interface {
	// ImagePull pulls an image from an OCI registry
	ImagePull(ctx context.Context, img Image, opts types.ImagePullOptions) (io.ReadCloser, error)
	// ImageList lists all images currently on the host
	ImageList(ctx context.Context, opts types.ImageListOptions) ([]Image, error)
	// ImageBuild builds an image
	ImageBuild(ctx context.Context, reader io.Reader, opts types.ImageBuildOptions) (io.ReadCloser, error)
	// ContainerCreate creates a container
	ContainerCreate(ctx context.Context, container Container) (Container, error)
	// ContainerStart starts a container
	ContainerStart(ctx context.Context, id string, opts types.ContainerStartOptions) error
	// ContainerStop stops a container
	ContainerStop(ctx context.Context, id string, timeout *time.Duration) error
	// ContainerDelete deletes a container
	ContainerDelete(ctx context.Context, id string, opts types.ContainerRemoveOptions) error
	// ContainerExecCreate creates an exec process in a container
	ContainerExecCreate(ctx context.Context, id string, cfg types.ExecConfig) (string, error)
	// ContainerExecAttach starts and attaches to an exec process
	// running in a container
	ContainerExecAttach(ctx context.Context, id string, chk types.ExecStartCheck) (HijackedResponse, error)
	// ContainerList lists containers currently running on the host
	ContainerList(ctx context.Context, opts types.ContainerListOptions) ([]Container, error)
	// CopyFromContainer copies files from within a container
	// onto the host machine so they can be used in a volume
	CopyFromContainer(ctx context.Context, id, src string) (io.ReadCloser, error)
}
