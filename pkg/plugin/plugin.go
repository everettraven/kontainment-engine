package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/kontainment/devcontainers-go/pkg/devcontainers"
)

const (
	defaultPluginDirectory = ".kontainment/plugins"
)

// Client is a client for interacting with kontainment plugins
type Client interface {
	// CreateWorkspace will create a new workspace using a plugin
	CreateWorkspace(plugin string, devcontainer devcontainers.DevContainer) error

	// DeleteWorkspace will delete an existing workspace using a plugin
	DeleteWorkspace(plugin string, devcontainer devcontainers.DevContainer) error
}

type ClientOption func(*client)

func WithRunner(r Runner) ClientOption {
	return func(c *client) {
		c.runner = r
	}
}

func NewClient(opts ...ClientOption) Client {
	c := &client{
		runner: NewRunner(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type client struct {
	runner Runner
}

func (c *client) CreateWorkspace(plugin string, devcontainer devcontainers.DevContainer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home dir: %w", err)
	}
	ppath := path.Join(home, defaultPluginDirectory, plugin)
	return c.runner.RunCreate(ppath, devcontainer)
}

func (c *client) DeleteWorkspace(plugin string, devcontainer devcontainers.DevContainer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home dir: %w", err)
	}
	ppath := path.Join(home, defaultPluginDirectory, plugin)
	return c.runner.RunCreate(ppath, devcontainer)
}

type Runner interface {
	RunCreate(pluginPath string, devcontainer devcontainers.DevContainer) error
	RunDelete(pluginPath string, devcontainer devcontainers.DevContainer) error
}

func NewRunner() Runner {
	return &runner{}
}

type runner struct{}

func (r *runner) RunCreate(pluginPath string, devcontainer devcontainers.DevContainer) error {
	dcBytes, err := json.Marshal(devcontainer)
	if err != nil {
		return fmt.Errorf("marshalling devcontainer: %w", err)
	}

	cmd := exec.Command(pluginPath, "create", string(dcBytes))
	return cmd.Run()
}

// TODO: We probably don't have to marshal and pass the *whole* devcontainer here
func (r *runner) RunDelete(pluginPath string, devcontainer devcontainers.DevContainer) error {
	dcBytes, err := json.Marshal(devcontainer)
	if err != nil {
		return fmt.Errorf("marshalling devcontainer: %w", err)
	}

	cmd := exec.Command(pluginPath, "delete", string(dcBytes))
	return cmd.Run()
}
