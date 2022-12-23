package workspace

import "net/http"

func (wr *WorkspaceRouter) postWorkspace(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("POST to /workspace"))

	// Parse Workspace definition from request
	// Create container
}
