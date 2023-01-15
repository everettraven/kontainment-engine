package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/kontainment/engine/api/types"
	"github.com/kontainment/engine/containertools"
)

func (wr *WorkspaceRouter) postWorkspace(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// Parse Workspace definition from request
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("reading the request body: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	workspace := &types.Workspace{}
	err = json.Unmarshal(data, workspace)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("unmarshalling the JSON from the request body: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	// pull image
	img := containertools.NewImage(
		containertools.WithRepository(workspace.Image.Repo),
		containertools.WithTag(workspace.Image.Tag),
	)
	rc, err := wr.ContainerRuntime.ImagePull(context.Background(), img, dockertypes.ImagePullOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("pulling the image: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}
	defer rc.Close()
	// Read all the contents of the image pull read closer
	// so that the image finishes pulling properly. It *seems* the image
	// doesn't finish pulling unless this is done. Should probably
	// investigate this further to determine if the is *actually* the case
	_, _ = io.ReadAll(rc)

	// create container
	container := containertools.NewContainer(
		containertools.WithImage(img),
		containertools.WithName(fmt.Sprintf("kontainment-workspace-%s", workspace.Name)),
	)
	container, err = wr.ContainerRuntime.ContainerCreate(context.Background(), container)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("creating container: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	// start container
	err = wr.ContainerRuntime.ContainerStart(context.Background(), container.Id(), dockertypes.ContainerStartOptions{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("starting the container: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
