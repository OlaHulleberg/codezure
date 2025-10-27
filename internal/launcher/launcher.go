package launcher

import (
    "fmt"
    "os"
    "os/exec"
    "github.com/OlaHulleberg/codzure/internal/azure"
    "github.com/OlaHulleberg/codzure/internal/profiles"
)

func Launch(passthrough []string) error {
    key, endpoint, err := azure.FetchKeyAndEndpoint()
    if err != nil { return err }
    pm, _ := profiles.NewManager()
    cfg, _ := pm.GetCurrentConfig("dev")

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
