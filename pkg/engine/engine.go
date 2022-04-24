package engine

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/nicklasfrahm/k3se/pkg/sshx"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

const (
	InstallerURL = "https://get.k3s.io"
)

// Engine is a type that encapsulates parts of the installation logic.
type Engine struct {
	Logger *zerolog.Logger

	sync.Mutex
	installer []byte

	Spec *Config
}

// New creates a new Engine.
func New(options ...Option) (*Engine, error) {
	opts, err := GetDefaultOptions().Apply(options...)
	if err != nil {
		return nil, err
	}

	return &Engine{
		Logger: opts.Logger,
	}, nil
}

// Installer returns the downloaded the k3s installer.
func (e *Engine) Installer() ([]byte, error) {
	// Lock engine to prevent concurrent access to installer cache.
	e.Lock()

	if len(e.installer) == 0 {
		resp, err := http.Get(InstallerURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		e.installer, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	e.Unlock()

	return e.installer, nil
}

// SetSpec configures the desired state of the cluster. Note
// that the config will only be applied if the verification
// succeeds.
func (e *Engine) SetSpec(config *Config) error {
	if err := config.Verify(); err != nil {
		return err
	}

	e.Spec = config

	return nil
}

// Configure uploads the installer and the configuration
// prior to a node prior to running the installation.
func (e *Engine) Configure(node *Node) error {
	// Upload the installer.
	installer, err := e.Installer()
	if err != nil {
		return err
	}

	if err := node.Upload("/tmp/k3se/install.sh", bytes.NewReader(installer)); err != nil {
		return err
	}

	// Create the node configuration.
	config := e.Spec.Cluster.Merge(&node.Config)

	// TODO: Configure the "advertise address" based on the first SAN and modify the
	// kubeconfig accordingly.

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := node.Upload("/tmp/k3se/config.yaml", bytes.NewReader(configBytes)); err != nil {
		return err
	}

	// TODO: Upload configuration and move it to appropriate location using "sudo".

	return nil
}

// Cleanup removes all temporary files from the node.
func (e *Engine) Cleanup(node *Node) error {
	return node.Do(sshx.Cmd{
		Cmd: "rm -rf /tmp/k3se",
	})
}
