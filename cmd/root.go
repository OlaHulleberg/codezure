package cmd

import (
	"fmt"
	"github.com/OlaHulleberg/codezure/internal/config"
	"github.com/OlaHulleberg/codezure/internal/interactive"
	"github.com/OlaHulleberg/codezure/internal/launcher"
	"github.com/OlaHulleberg/codezure/internal/profiles"
	"github.com/OlaHulleberg/codezure/internal/updater"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	codezureProfileFlag string
	Version             = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "codezure",
	Short: "Launch Codex CLI with Azure OpenAI configuration",
	Long:  `codezure configures Azure OpenAI env and launches Codex CLI, or prints env for manual use.`,
	Args:  cobra.ArbitraryArgs, // Accept any args for passthrough to Codex
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&codezureProfileFlag, "codezure-profile", "", "Use a specific codezure profile for this run")

	// Allow unknown flags to pass through to Codex CLI
	rootCmd.FParseErrWhitelist.UnknownFlags = true
	// Register subcommands
	rootCmd.AddCommand(manageCmd)
}

func runRoot(cmd *cobra.Command, args []string) error {
	// Collect passthrough args for Codex CLI
	passthroughArgs := collectPassthroughArgs()

	// Check for updates in background
	go updater.CheckForUpdates(Version)

	// Load configuration from profile
	pm, err := profiles.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create profile manager: %w", err)
	}

	var cfg *config.Config
	profile := codezureProfileFlag
	if profile != "" {
		// Load specific profile
		cfg, err = pm.Load(profile)
		if err != nil {
			return fmt.Errorf("failed to load profile '%s': %w", profile, err)
		}
		fmt.Printf("Using profile: %s\n\n", profile)
	} else {
		// First-run: if no current profile, trigger interactive GUI
		if _, e := pm.GetCurrent(); e != nil {
			// Launch interactive config to save current profile
			if err := interactive.RunInteractiveConfig(Version, pm); err != nil {
				return err
			}
		}
		// Load current profile
		cfg, err = pm.GetCurrentConfig(Version)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	// Validate configuration
	if err := pm.Validate(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	return launcher.Launch(passthroughArgs)
}

// collectPassthroughArgs separates codezure flags from Codex CLI args
func collectPassthroughArgs() []string {
	if len(os.Args) <= 1 {
		return nil
	}

	var passthroughArgs []string
	// codezure flags and whether they require a value as the next arg
	codezureFlags := map[string]bool{
		"--codezure-profile": true,
	}

	skip := false
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if skip {
			skip = false
			continue
		}

		// Check if this is a codezure flag
		if strings.HasPrefix(arg, "--codezure-") {
			// Check if it's a flag with value (--flag=value or --flag value)
			if strings.Contains(arg, "=") {
				// --flag=value format, skip entirely
				continue
			}
			if requires, ok := codezureFlags[arg]; ok && requires {
				// --flag value format, skip this and next arg
				skip = true
				continue
			}
			// codezure boolean flag; skip this single arg
			continue
		}

		// This is a passthrough arg
		passthroughArgs = append(passthroughArgs, arg)
	}

	return passthroughArgs
}
