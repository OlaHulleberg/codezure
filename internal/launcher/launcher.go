package launcher

import (
	"fmt"
	"github.com/OlaHulleberg/codzure/internal/azure"
	"github.com/OlaHulleberg/codzure/internal/profiles"
	"github.com/OlaHulleberg/codzure/internal/secrets"
	"os"
	"os/exec"
)

func Launch(passthrough []string) error {
	pm, _ := profiles.NewManager()
	cfg, _ := pm.GetCurrentConfig("dev")

	// Determine auth mode (default azure-cli)
	auth := cfg.Auth
	if auth == "" {
		auth = "azure-cli"
	}

	var key string
	var endpoint string
	var err error

	switch auth {
	case "api-key":
		// Fetch API key from OS keychain
		profileName, e := pm.GetCurrent()
		if e != nil || profileName == "" {
			profileName = "default"
		}
		key, err = secrets.GetKey(profileName)
		if err != nil {
			return fmt.Errorf("failed to retrieve API key from keychain for profile '%s': %w", profileName, err)
		}
		endpoint = cfg.Endpoint
		if endpoint == "" {
			return fmt.Errorf("endpoint not set in profile; run 'codzure manage config' to configure")
		}
	case "azure-cli":
		key, endpoint, err = azure.FetchKeyAndEndpoint()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown auth mode: %s", auth)
	}

	os.Setenv("AZURE_OPENAI_API_KEY", key)
	os.Setenv("AZURE_OPENAI_ENDPOINT", endpoint)
	os.Setenv("AZURE_OPENAI_DEPLOYMENT", cfg.Deployment)
	os.Setenv("OPENAI_API_KEY", key)
	os.Setenv("OPENAI_BASE_URL", endpoint+"/openai/v1")

	if _, err := exec.LookPath("codex"); err == nil {
		cmd := exec.Command("codex", passthrough...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
	return fmt.Errorf("codex CLI not found on PATH; install Codex and ensure it's on your PATH")
}

// runLaunch is a small adapter used by cmd/root
func runLaunch(args []string) error { return Launch(args) }
