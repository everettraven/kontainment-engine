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
	ImagePull(context.Context, string, types.ImagePullOptions) (Image, error)
	// ImageList lists all images currently on the host
	ImageList(context.Context, types.ImageListOptions) ([]Image, error)
	// ImageBuild builds an image
	ImageBuild(context.Context, io.Reader, types.ImageBuildOptions) (Image, error)
	// ContainerCreate creates a container
	ContainerCreate(context.Context, Container) (Container, error)
	// ContainerStart starts a container
	ContainerStart(context.Context, string, types.ContainerStartOptions) (Container, error)
	// ContainerStop stops a container
	ContainerStop(context.Context, string, *time.Duration) (Container, error)
	// ContainerDelete deletes a container
	ContainerDelete(context.Context, string, types.ContainerRemoveOptions) (Container, error)
	// ContainerExecCreate creates an exec process in a container
	ContainerExecCreate(context.Context, string, types.ExecConfig) (ContainerExec, error)
	// ContainerExecAttach starts and attaches to an exec process
	// running in a container
	ContainerExecAttach(context.Context, string, types.ExecStartCheck) (HijackedResponse, error)
	// ContainerList lists containers currently running on the host
	ContainerList(context.Context, types.ContainerListOptions) ([]Container, error)
	// CopyFromContainer copies files from within a container
	// onto the host machine so they can be used in a volume
	CopyFromContainer(context.Context, string, string) (io.ReadCloser, error)
}
