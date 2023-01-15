package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kontainment/engine/api/types"
)

type KontainmentClient struct {
	client http.Client
}

func NewClient() *KontainmentClient {
	return &KontainmentClient{
		client: http.Client{},
	}
}

func (c *KontainmentClient) CreateWorkspace(ctx context.Context, wksp *types.Workspace) error {
	// marshal to bytes to set up as a reader
	body, err := json.Marshal(wksp)
	if err != nil {
		return fmt.Errorf("marshalling workspace: %w", err)
	}

	resp, err := c.client.Post("http://127.0.0.1:8080/kontainment/workspace", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating workspace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading error from response: %w", err)
		}

		apiError := &types.ApiError{}
		err = json.Unmarshal(errBytes, apiError)
		if err != nil {
			return fmt.Errorf("unmarshalling error to ApiError type: %w", err)
		}
		return fmt.Errorf("engine failed to create workspace: %s", apiError.Msg)
	}

	return nil
}

func (c *KontainmentClient) DeleteWorkspace(ctx context.Context, wksp *types.Workspace) error {
	// marshal to bytes to set up as a reader
	body, err := json.Marshal(wksp)
	if err != nil {
		return fmt.Errorf("marshalling workspace: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8080/kontainment/workspace", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building the DELETE request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("deleting workspace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading error from response: %w", err)
		}

		apiError := &types.ApiError{}
		err = json.Unmarshal(errBytes, apiError)
		if err != nil {
			return fmt.Errorf("unmarshalling error to ApiError type: %w", err)
		}
		return fmt.Errorf("engine failed to delete workspace: %s", apiError.Msg)
	}

	return nil
}

func (c *KontainmentClient) ListWorkspaces(ctx context.Context) (*types.WorkspaceList, error) {
	resp, err := c.client.Get("http://127.0.0.1:8080/kontainment/workspaces")
	if err != nil {
		return nil, fmt.Errorf("getting workspaces: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading error from response: %w", err)
		}

		apiError := &types.ApiError{}
		err = json.Unmarshal(errBytes, apiError)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling error to ApiError type: %w", err)
		}
		return nil, fmt.Errorf("engine failed to list workspaces: %s", apiError.Msg)
	}

	// unmarshal json to workspace list
	wkspList := &types.WorkspaceList{}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	err = json.Unmarshal(respBytes, wkspList)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response: %w", err)
	}

	return wkspList, nil
}
