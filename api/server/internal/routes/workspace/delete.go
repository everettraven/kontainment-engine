package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/kontainment/engine/api/types"
	"github.com/kontainment/engine/containertools"
)

func (wr *WorkspaceRouter) deleteWorkspace(w http.ResponseWriter, r *http.Request) {
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

	// get list of containers filtered by name
	containers, err := wr.ContainerRuntime.ContainerList(context.Background(),
		dockertypes.ContainerListOptions{},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("listing containers: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	containerList := []containertools.Container{}
	for _, container := range containers {
		if container.Name() == fmt.Sprintf("/kontainment-workspace-%s", workspace.Name) {
			containerList = append(containerList, container)
		}
	}

	// there should really only be a single container in the list
	if len(containerList) > 1 {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("More than one container with the name %q - Can't pick one", workspace.Name))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	if len(containerList) > 0 {
		timeout := time.Duration(1) * time.Second
		err = wr.ContainerRuntime.ContainerStop(context.Background(), containerList[0].Id(), &timeout)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			apiError := types.NewApiError(fmt.Sprintf("stopping container: %s", err.Error()))
			// TODO: Should consider how to handle if this errors
			errBytes, _ := json.Marshal(apiError)
			w.Write(errBytes)
			return
		}

		err = wr.ContainerRuntime.ContainerDelete(context.Background(), containerList[0].Id(), dockertypes.ContainerRemoveOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			apiError := types.NewApiError(fmt.Sprintf("deleting container: %s", err.Error()))
			// TODO: Should consider how to handle if this errors
			errBytes, _ := json.Marshal(apiError)
			w.Write(errBytes)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
