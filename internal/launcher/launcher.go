package launcher

import (
	"fmt"
	"github.com/OlaHulleberg/codezure/internal/azure"
	"github.com/OlaHulleberg/codezure/internal/profiles"
	"github.com/OlaHulleberg/codezure/internal/secrets"
	"os"
	"os/exec"
	"strings"
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
			return fmt.Errorf("endpoint not set in profile; run 'codezure manage config' to configure")
		}
	case "azure-cli":
		key, endpoint, err = azure.FetchKeyAndEndpoint()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown auth mode: %s", auth)
	}

	// Build child process environment; avoid mutating global env
	env := os.Environ()
	env = append(env, "CODEZURE_API_KEY="+key)

	if _, err := exec.LookPath("codex"); err == nil {
		// Build Codex overrides (non-destructive) and append to passthrough.
		prof, e := pm.GetCurrent()
		if e != nil || prof == "" {
			prof = "default"
		}
		// Configure Codex via runtime overrides (no system file writes).
		// Add only keys the user hasn't already specified.
		if !hasOverrideKey(passthrough, "model_provider") {
			// Use codezure provider wired to Azure Responses API
			passthrough = append(passthrough, "--config", "model_provider=\"codezure\"")
			passthrough = append(passthrough, "--config", "model_providers.codezure.name=\"Codezure\"")
			ep := strings.TrimRight(strings.TrimSpace(endpoint), "/")
			base := ep + "/openai/v1"
			passthrough = append(passthrough, "--config", fmt.Sprintf("model_providers.codezure.base_url=%q", base))
			passthrough = append(passthrough, "--config", "model_providers.codezure.env_key=\"CODEZURE_API_KEY\"")
			passthrough = append(passthrough, "--config", "model_providers.codezure.wire_api=\"responses\"")
		}
		if !hasOverrideKey(passthrough, "model") && !hasModelFlag(passthrough) {
			if strings.TrimSpace(cfg.Deployment) != "" {
				passthrough = append(passthrough, "--config", fmt.Sprintf("model=%q", strings.TrimSpace(cfg.Deployment)))
			}
		}
		if strings.TrimSpace(cfg.Thinking) != "" && !hasOverrideKey(passthrough, "model_reasoning_effort") {
			passthrough = append(passthrough, "--config", fmt.Sprintf("model_reasoning_effort=%q", strings.TrimSpace(cfg.Thinking)))
		}

		cmd := exec.Command("codex", passthrough...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		return cmd.Run()
	}
	return fmt.Errorf("codex CLI not found on PATH; install Codex and ensure it's on your PATH")
}

// runLaunch is a small adapter used by cmd/root
func runLaunch(args []string) error { return Launch(args) }

// hasConfigFlag returns true if passthrough contains --config/-c
func hasConfigFlag(args []string) bool {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--config" || a == "-c" {
			return true
		}
		// Handle --config=path form
		if strings.HasPrefix(a, "--config=") {
			return true
		}
	}
	return false
}

// hasModelFlag returns true if passthrough contains --model/-m
func hasModelFlag(args []string) bool {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--model" || a == "-m" { // handle next arg as value
			return true
		}
		if strings.HasPrefix(a, "--model=") || strings.HasPrefix(a, "-m=") {
			return true
		}
	}
	return false
}

// hasOverrideKey reports if any -c/--config override sets the given top-level key or nested path
// e.g., hasOverrideKey(args, "model_provider") or hasOverrideKey(args, "model_reasoning_effort").
func hasOverrideKey(args []string, key string) bool {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--config" || a == "-c" {
			if i+1 < len(args) {
				v := args[i+1]
				if strings.HasPrefix(v, key+"=") || strings.HasPrefix(v, key+".") {
					return true
				}
			}
		} else if strings.HasPrefix(a, "--config=") {
			v := strings.TrimPrefix(a, "--config=")
			if strings.HasPrefix(v, key+"=") || strings.HasPrefix(v, key+".") {
				return true
			}
		}
	}
	return false
}
