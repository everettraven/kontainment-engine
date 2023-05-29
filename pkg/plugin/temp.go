package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/kontainment/devcontainers-go/pkg/devcontainers"
	"github.com/kontainment/engine/pkg/containertools"
)

// TODO: Remove this in place of creating an actual plugin
var _ Runner = &DockerRunner{}

type DockerRunner struct {
	Runtime containertools.ContainerRuntime
}

func (dr *DockerRunner) RunCreate(pluginPath string, devcontainer devcontainers.DevContainer) error {
	imageSplit := strings.Split(devcontainer.Image, ":")
	// pull image
	img := containertools.NewImage(
		containertools.WithRepository(imageSplit[0]),
		containertools.WithTag(imageSplit[1]),
	)
	rc, err := dr.Runtime.ImagePull(context.Background(), img, dockertypes.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pulling image: %w", err)
	}
	defer rc.Close()
	// Read all the contents of the image pull read closer
	// so that the image finishes pulling properly. It *seems* the image
	// doesn't finish pulling unless this is done. Should probably
	// investigate this further to determine if the is *actually* the case
	_, _ = io.ReadAll(rc)

	// TODO: create the workspace files locally and mount as a volume

	// create container
	container := containertools.NewContainer(
		containertools.WithImage(img),
		containertools.WithName(fmt.Sprintf("kontainment-workspace-%s", devcontainer.Name)),
	)
	container, err = dr.Runtime.ContainerCreate(context.Background(), container)
	if err != nil {
		return fmt.Errorf("creating container %w", err)
	}

	// start container
	err = dr.Runtime.ContainerStart(context.Background(), container.Id(), dockertypes.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	return nil
}

func (dr *DockerRunner) RunDelete(pluginPath string, devcontainer devcontainers.DevContainer) error {
	// get list of containers filtered by name
	containers, err := dr.Runtime.ContainerList(context.Background(),
		dockertypes.ContainerListOptions{},
	)
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}

	containerList := []containertools.Container{}
	for _, container := range containers {
		if container.Name() == fmt.Sprintf("/kontainment-workspace-%s", devcontainer.Name) {
			containerList = append(containerList, container)
		}
	}

	// there should really only be a single container in the list
	if len(containerList) > 1 {
		return errors.New("more than 1 container in the list")
	}

	if len(containerList) > 0 {
		timeout := time.Duration(1) * time.Second
		err = dr.Runtime.ContainerStop(context.Background(), containerList[0].Id(), &timeout)
		if err != nil {
			return fmt.Errorf("stopping container: %w", err)
		}

		err = dr.Runtime.ContainerDelete(context.Background(), containerList[0].Id(), dockertypes.ContainerRemoveOptions{})
		if err != nil {
			return fmt.Errorf("deleting container: %w", err)
		}
	}

	return nil
}
