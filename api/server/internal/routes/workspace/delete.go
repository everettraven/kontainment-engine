package workspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kontainment/engine/api/types"
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

	if err = wr.PluginClient.DeleteWorkspace(workspace.Plugin, workspace.DevContainer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := types.NewApiError(fmt.Sprintf("deleting workspace: %s", err.Error()))
		// TODO: Should consider how to handle if this errors
		errBytes, _ := json.Marshal(apiError)
		w.Write(errBytes)
		return
	}

	delete(wr.WorkspaceCache, workspace.DevContainer.Name)

	w.WriteHeader(http.StatusNoContent)
}
