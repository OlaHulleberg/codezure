package cmd

import (
    "fmt"
    "strings"
    "github.com/OlaHulleberg/codzure/internal/interactive"
    "github.com/OlaHulleberg/codzure/internal/profiles"
    "github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Interactive configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Run interactive wizard
        return runInteractiveConfig()
    },
}

var configListCmd = &cobra.Command{
    Use:   "list",
    Short: "List current configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        pm, err := profiles.NewManager()
        if err != nil { return err }
        cfg, err := pm.GetCurrentConfig(Version)
        if err != nil { return err }
        fmt.Println("Current Configuration:")
        fmt.Printf("  subscription: %s\n", cfg.Subscription)
        fmt.Printf("  group:        %s\n", cfg.Group)
        fmt.Printf("  resource:     %s\n", cfg.Resource)
        fmt.Printf("  location:     %s\n", cfg.Location)
        fmt.Printf("  deployment:   %s\n", cfg.Deployment)
        if cfg.Thinking != "" {
            fmt.Printf("  thinking:     %s\n", cfg.Thinking)
        }
        return nil
    },
}

var configSetCmd = &cobra.Command{
    Use:   "set <key> <value>",
    Short: "Set a configuration value",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        key := strings.ToLower(args[0])
        val := args[1]
        pm, err := profiles.NewManager()
        if err != nil { return err }
        cfg, err := pm.GetCurrentConfig(Version)
        if err != nil { return err }
        switch key {
        case "subscription": cfg.Subscription = val
        case "group": cfg.Group = val
        case "resource": cfg.Resource = val
        case "location": cfg.Location = val
        case "deployment": cfg.Deployment = val
        case "thinking": cfg.Thinking = val
        default:
            return fmt.Errorf("unknown key: %s", key)
        }
        return pm.SaveCurrentConfig(cfg)
    },
}

func init() {
    configCmd.AddCommand(configListCmd)
    configCmd.AddCommand(configSetCmd)
}

func runInteractiveConfig() error {
    pm, err := profiles.NewManager(); if err != nil { return err }
    return interactive.RunInteractiveConfig(Version, pm)
}
