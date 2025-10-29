package azure

import (
	"fmt"
	"github.com/OlaHulleberg/codzure/internal/profiles"
	"os"
	"os/exec"
	"strings"
)

func requireAz() error {
	if _, err := exec.LookPath("az"); err != nil {
		return fmt.Errorf("Azure CLI (az) not found. Install from https://aka.ms/azcli and run 'az login'.")
	}
	return nil
}

func FetchKeyAndEndpoint() (string, string, error) {
	if err := requireAz(); err != nil {
		return "", "", err
	}
	pm, err := profiles.NewManager()
	if err != nil {
		return "", "", err
	}
	cfg, err := pm.GetCurrentConfig("dev")
	if err != nil {
		return "", "", err
	}
	if err := pm.Validate(cfg); err != nil {
		return "", "", err
	}
	if err := runCmd("az", "account", "set", "--subscription", cfg.Subscription); err != nil {
		return "", "", err
	}
	keyBytes, err := runCmdOutput("az", "cognitiveservices", "account", "keys", "list",
		"--name", cfg.Resource, "--resource-group", cfg.Group, "--query", "key1", "-o", "tsv")
	if err != nil {
		return "", "", err
	}
	endBytes, err := runCmdOutput("az", "cognitiveservices", "account", "show",
		"--name", cfg.Resource, "--resource-group", cfg.Group, "--query", "properties.endpoint", "-o", "tsv")
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(string(keyBytes)), strings.TrimSpace(string(endBytes)), nil
}

// GetEndpoint returns the endpoint URL for a given resource
func GetEndpoint(subscription, resource, group string) (string, error) {
	if err := requireAz(); err != nil {
		return "", err
	}
	if err := runCmd("az", "account", "set", "--subscription", subscription); err != nil {
		return "", err
	}
	out, err := runCmdOutput("az", "cognitiveservices", "account", "show",
		"--name", resource, "--resource-group", group, "--query", "properties.endpoint", "-o", "tsv")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}
