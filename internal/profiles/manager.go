package profiles

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "github.com/OlaHulleberg/codzure/internal/config"
)

type Manager struct {
    dir      string
    profiles string
}

func NewManager() (*Manager, error) {
    home, err := os.UserHomeDir()
    if err != nil { return nil, err }
    dir := filepath.Join(home, ".codzure")
    prof := filepath.Join(dir, "profiles")
    if err := os.MkdirAll(prof, 0o755); err != nil { return nil, err }
    return &Manager{dir: dir, profiles: prof}, nil
}

func (m *Manager) currentProfilePath() string {
    return filepath.Join(m.dir, "current-profile.txt")
}

func (m *Manager) legacyEnvPath() string {
    return filepath.Join(m.dir, "current.env")
}

func (m *Manager) profileFile(name string) string { return filepath.Join(m.profiles, name+".json") }

func (m *Manager) GetCurrent() (string, error) {
    b, err := os.ReadFile(m.currentProfilePath())
    if err != nil { return "", err }
    return strings.TrimSpace(string(b)), nil
}

func (m *Manager) SetCurrent(name string) error {
    return os.WriteFile(m.currentProfilePath(), []byte(name+"\n"), 0o644)
}

// GetCurrentConfig loads current profile config, migrating legacy current.env if present
func (m *Manager) GetCurrentConfig(version string) (*config.Config, error) {
    // Migrate legacy file if needed
    if _, err := os.Stat(m.legacyEnvPath()); err == nil {
        // no profiles? create default from legacy
        if _, err := os.Stat(m.profileFile("default")); os.IsNotExist(err) {
            cfg, err := readEnvFile(m.legacyEnvPath())
            if err != nil { return nil, err }
            if err := writeJSONFile(m.profileFile("default"), cfg); err != nil { return nil, err }
            // backup legacy
            _ = os.Rename(m.legacyEnvPath(), m.legacyEnvPath()+".bak")
            _ = m.SetCurrent("default")
        }
    }

    // If no current profile, return an error so caller can trigger interactive GUI
    name, err := m.GetCurrent()
    if err != nil || name == "" {
        return nil, fmt.Errorf("no current profile configured; run 'codzure manage config'")
    }
    return readJSONFile(m.profileFile(name))
}

func (m *Manager) SaveCurrentConfig(cfg *config.Config) error {
    name, err := m.GetCurrent()
    if err != nil || name == "" { name = "default"; _ = m.SetCurrent(name) }
    return writeJSONFile(m.profileFile(name), cfg)
}

// Load loads a specific profile by name
func (m *Manager) Load(profileName string) (*config.Config, error) {
    return readJSONFile(m.profileFile(profileName))
}

func readEnvFile(path string) (*config.Config, error) {
    b, err := os.ReadFile(path)
    if err != nil { return nil, err }
    cfg := &config.Config{}
    lines := strings.Split(string(b), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") { continue }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 { continue }
        k := parts[0]; v := parts[1]
        switch strings.ToUpper(k) {
        case "CODZURE_SUBSCRIPTION": cfg.Subscription = v
        case "CODZURE_GROUP": cfg.Group = v
        case "CODZURE_RESOURCE": cfg.Resource = v
        case "CODZURE_LOCATION": cfg.Location = v
        case "CODZURE_DEPLOYMENT": cfg.Deployment = v
        }
    }
    return cfg, nil
}

func writeJSONFile(path string, cfg *config.Config) error {
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return err }
    b, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil { return err }
    return os.WriteFile(path, b, 0o644)
}

func readJSONFile(path string) (*config.Config, error) {
    b, err := os.ReadFile(path)
    if err != nil { return nil, err }
    var cfg config.Config
    if err := json.Unmarshal(b, &cfg); err != nil { return nil, err }
    return &cfg, nil
}

// prompt helpers removed; interactive GUI will handle first-run

func (m *Manager) Validate(cfg *config.Config) error {
    if strings.TrimSpace(cfg.Subscription) == "" || strings.TrimSpace(cfg.Group) == "" || strings.TrimSpace(cfg.Resource) == "" {
        return errors.New("subscription/group/resource must be set; run 'codzure manage config'")
    }
    return nil
}
