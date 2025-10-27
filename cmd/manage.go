package cmd

import (
    "github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
    Use:   "manage",
    Short: "Manage configuration and resources",
}

func init() {
    manageCmd.AddCommand(configCmd)
    manageCmd.AddCommand(modelsCmd)
    manageCmd.AddCommand(profilesCmd)
    manageCmd.AddCommand(versionCmd)
    manageCmd.AddCommand(updateCmd)
}
