package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/kontainment/engine/containertools"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var _ containertools.ContainerRuntime = &DockerRuntime{}

// DockerRuntime is an implementation of the
// ContainerRuntime interface that uses Docker
// for container operations
type DockerRuntime struct {
	client *client.Client
}

type DockerRuntimeOption func(*DockerRuntime)

func WithDockerClient(client *client.Client) DockerRuntimeOption {
	return func(dr *DockerRuntime) {
		dr.client = client
	}
}

func NewDockerRuntime(opts ...DockerRuntimeOption) *DockerRuntime {
	dr := &DockerRuntime{}

	for _, opt := range opts {
		opt(dr)
	}

	return dr
}

func (dr *DockerRuntime) ImagePull(ctx context.Context, img containertools.Image, opts types.ImagePullOptions) (io.ReadCloser, error) {
	return dr.client.ImagePull(ctx, imageRef(img), opts)
}

func (dr *DockerRuntime) ImageList(ctx context.Context, opts types.ImageListOptions) ([]containertools.Image, error) {
	imgList, err := dr.client.ImageList(ctx, opts)
	if err != nil {
		return nil, err
	}

	images := []containertools.Image{}

	for _, img := range imgList {
		if len(img.RepoTags) > 0 {
			repoTag := strings.Split(img.RepoTags[0], ":")
			images = append(images,
				containertools.NewImage(
					containertools.WithImageId(img.ID),
					containertools.WithRepository(repoTag[0]),
					containertools.WithTag(repoTag[1]),
				),
			)
		}
	}

	return images, nil
}

func (dr *DockerRuntime) ImageBuild(ctx context.Context, reader io.Reader, opts types.ImageBuildOptions) (io.ReadCloser, error) {
	resp, err := dr.client.ImageBuild(ctx, reader, opts)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (dr *DockerRuntime) ContainerCreate(ctx context.Context, cont containertools.Container) (containertools.Container, error) {
	created, err := dr.client.ContainerCreate(ctx,
		&container.Config{
			Image: imageRef(cont.Image()),
			Cmd:   []string{"tail", "-f", "/dev/null"},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		&v1.Platform{},
		cont.Name(),
	)
	if err != nil {
		return nil, err
	}

	return containertools.NewContainer(
		containertools.WithContainerId(created.ID),
		containertools.WithName(cont.Name()),
		containertools.WithImage(cont.Image()),
	), nil
}

func (dr *DockerRuntime) ContainerStart(ctx context.Context, id string, opts types.ContainerStartOptions) error {
	return dr.client.ContainerStart(ctx, id, opts)
}

func (dr *DockerRuntime) ContainerStop(ctx context.Context, id string, timeout *time.Duration) error {
	return dr.client.ContainerStop(ctx, id, timeout)
}

func (dr *DockerRuntime) ContainerDelete(ctx context.Context, id string, opts types.ContainerRemoveOptions) error {
	return dr.client.ContainerRemove(ctx, id, opts)
}

func (dr *DockerRuntime) ContainerExecCreate(ctx context.Context, containerId string, cfg types.ExecConfig) (string, error) {
	resp, err := dr.client.ContainerExecCreate(ctx, containerId, cfg)
	if err != nil {
		return "", nil
	}

	return resp.ID, nil
}

func (dr *DockerRuntime) ContainerExecAttach(ctx context.Context, execId string, chk types.ExecStartCheck) (containertools.HijackedResponse, error) {
	resp, err := dr.client.ContainerExecAttach(ctx, execId, chk)
	if err != nil {
		return nil, err
	}

	return containertools.NewHijackedResponse(
		containertools.WithConn(resp.Conn),
		containertools.WithReader(resp.Reader),
	), nil
}

func (dr *DockerRuntime) ContainerList(ctx context.Context, opts types.ContainerListOptions) ([]containertools.Container, error) {
	resp, err := dr.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	containerList := []containertools.Container{}

	for _, container := range resp {
		containerList = append(containerList,
			containertools.NewContainer(
				containertools.WithContainerId(container.ID),
				containertools.WithImage(containertools.NewImage(
					containertools.WithImageId(container.ImageID),
					containertools.WithRepository(strings.Split(container.Image, ":")[0]),
					containertools.WithTag(strings.Split(container.Image, ":")[1]),
				)),
				containertools.WithName(container.Names[0]),
				//TODO: Add ports and volumes and other information
			),
		)
	}

	return containerList, nil
}

func (dr *DockerRuntime) CopyFromContainer(ctx context.Context, id, src string) (io.ReadCloser, error) {
	rc, _, err := dr.client.CopyFromContainer(ctx, id, src)
	return rc, err
}

func imageRef(img containertools.Image) string {
	return fmt.Sprintf("%s:%s", img.Repository(), img.Tag())
}
