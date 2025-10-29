package cmd

import (
	"fmt"
	"github.com/OlaHulleberg/codezure/internal/azure"
	"github.com/OlaHulleberg/codezure/internal/profiles"
	"github.com/spf13/cobra"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Model operations",
}

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments (models) in current resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		cfg, err := pm.GetCurrentConfig(Version)
		if err != nil {
			return err
		}
		deps, err := azure.ListDeployments(cfg.Subscription, cfg.Resource, cfg.Group)
		if err != nil {
			return err
		}
		fmt.Println("Available deployments:")
		for _, d := range deps {
			fmt.Printf("  %s (model=%s)\n", d.Name, d.ModelName)
		}
		return nil
	},
}

func init() {
	manageCmd.AddCommand(modelsCmd)
	modelsCmd.AddCommand(modelsListCmd)
}
