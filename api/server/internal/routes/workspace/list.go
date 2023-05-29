package workspace

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kontainment/engine/api/types"
)

func (wr *WorkspaceRouter) listWorkspaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	wkspList := &types.WorkspaceList{
		Workspaces: []types.Workspace{},
	}

	for _, wksp := range wr.WorkspaceCache {
		wkspList.Workspaces = append(wkspList.Workspaces, *wksp)
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
