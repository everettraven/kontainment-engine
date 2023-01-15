package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/kontainment/engine/api/types"
)

func (wr *WorkspaceRouter) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// get list of containers filtered by name
	containers, err := wr.ContainerRuntime.ContainerList(context.Background(),
		dockertypes.ContainerListOptions{All: true},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("listing containers: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	wkspList := &types.WorkspaceList{
		Workspaces: []types.Workspace{},
	}

	for _, container := range containers {
		if strings.Contains(container.Name(), "kontainment-workspace") {
			wksp := types.NewWorkspace(
				types.WithName(container.Name()),
				types.WithImage(container.Image()),
			)
			wkspList.Workspaces = append(wkspList.Workspaces, *wksp)
		}
	}

	jsonData, err := json.Marshal(wkspList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("marshalling the response: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
