package interactive

import (
    "fmt"
    "github.com/OlaHulleberg/codzure/internal/azure"
    "github.com/OlaHulleberg/codzure/internal/config"
    "github.com/OlaHulleberg/codzure/internal/profiles"
)

// RunInteractiveConfig runs an interactive configuration wizard using Bubbletea selector
func RunInteractiveConfig(currentVersion string, mgr *profiles.Manager) error {
    cfg, err := mgr.GetCurrentConfig(currentVersion)
    if err != nil {
        // No current profile; start with defaults and proceed with interactive GUI
        cfg = &config.Config{}
    }
    currentProfile, err := mgr.GetCurrent()
    if err != nil {
        // No current profile yet; use 'default' for messages. SaveCurrentConfig will set it.
        currentProfile = "default"
    }

    // Subscriptions
    subs, err := azure.ListSubscriptions(); if err != nil { return fmt.Errorf("failed to list subscriptions: %w", err) }
    subOpts := make([]SelectOption, len(subs))
    for i, s := range subs { subOpts[i] = SelectOption{ID: s.ID, Display: fmt.Sprintf("%s (%s)", s.Name, s.ID)} }
    subID, err := InteractiveSelect("Select Subscription", "Type to filter subscriptions...", subOpts, cfg.Subscription)
    if err != nil { return fmt.Errorf("subscription selection failed: %w", err) }

    // Resources
    resList, err := azure.ListOpenAIResources(subID); if err != nil { return fmt.Errorf("failed to list resources: %w", err) }
    resOpts := make([]SelectOption, len(resList))
    for i, r := range resList { resOpts[i] = SelectOption{ID: r.Name, Display: fmt.Sprintf("%s — rg=%s, region=%s", r.Name, r.ResourceGroup, r.Location)} }
    resName, err := InteractiveSelect("Select Azure OpenAI Resource", "Type to filter resources...", resOpts, cfg.Resource)
    if err != nil { return fmt.Errorf("resource selection failed: %w", err) }
    var res azure.OpenAIResource
    for _, r := range resList { if r.Name == resName { res = r; break } }
    endpoint, err := azure.GetEndpoint(subID, res.Name, res.ResourceGroup)
    if err != nil { return fmt.Errorf("failed to get endpoint: %w", err) }

    // Deployments (Models)
    deps, err := azure.ListDeployments(subID, res.Name, res.ResourceGroup); if err != nil { return fmt.Errorf("failed to list deployments: %w", err) }
    depOpts := make([]SelectOption, len(deps))
    for i, d := range deps { depOpts[i] = SelectOption{ID: d.Name, Display: fmt.Sprintf("%s — model=%s", d.Name, d.ModelName)} }
    depName, err := InteractiveSelect("Select Model Deployment", "Type to filter models...", depOpts, cfg.Deployment)
    if err != nil { return fmt.Errorf("deployment selection failed: %w", err) }

    // Thinking level
    tl := azure.ThinkingLevels()
    tlOpts := make([]SelectOption, len(tl))
    for i, s := range tl {
        desc := s
        switch s { case "low": desc = "low — fastest, cheapest"; case "medium": desc = "medium — balanced"; case "high": desc = "high — deepest reasoning" }
        tlOpts[i] = SelectOption{ID: s, Display: desc}
    }
    thinking, err := InteractiveSelect("Select Thinking Level", "Type to filter levels...", tlOpts, cfg.Thinking)
    if err != nil { thinking = cfg.Thinking } // allow cancel to keep existing

    // Update cfg
    cfg.Subscription = subID
    cfg.Group = res.ResourceGroup
    cfg.Resource = res.Name
    cfg.Location = res.Location
    cfg.Endpoint = endpoint
    cfg.Deployment = depName
    if thinking != "" { cfg.Thinking = thinking }

    if err := mgr.SaveCurrentConfig(cfg); err != nil { return fmt.Errorf("failed to save config: %w", err) }
    fmt.Printf("\n✓ Configuration saved successfully to profile '%s'!\n", currentProfile)
    fmt.Printf("\nConfiguration:\n")
    fmt.Printf("  Subscription: %s\n", cfg.Subscription)
    fmt.Printf("  ResourceGrp:  %s\n", cfg.Group)
    fmt.Printf("  Resource:     %s\n", cfg.Resource)
    fmt.Printf("  Region:       %s\n", cfg.Location)
    fmt.Printf("  Endpoint:     %s\n", cfg.Endpoint)
    fmt.Printf("  Deployment:   %s\n", cfg.Deployment)
    if cfg.Thinking != "" { fmt.Printf("  Thinking:     %s\n", cfg.Thinking) }
    return nil
}
