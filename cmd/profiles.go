package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/OlaHulleberg/codezure/internal/config"
	"github.com/OlaHulleberg/codezure/internal/profiles"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		entries, err := os.ReadDir(pmProfilesDir(pm))
		if err != nil {
			return err
		}
		current, _ := pm.GetCurrent()
		for _, e := range entries {
			if filepath.Ext(e.Name()) == ".json" {
				name := e.Name()[:len(e.Name())-5]
				mark := ""
				if name == current {
					mark = " *"
				}
				fmt.Printf("%s%s\n", name, mark)
			}
		}
		return nil
	},
}

func pmProfilesDir(pm *profiles.Manager) string { return filepath.Join(pmDir(pm), "profiles") }
func pmDir(pm *profiles.Manager) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".codezure")
}

var configSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current config as profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		cfg, err := pm.GetCurrentConfig(Version)
		if err != nil {
			return err
		}
		name := args[0]
		return os.WriteFile(filepath.Join(pmProfilesDir(pm), name+".json"), mustJSON(cfg), 0o644)
	},
}

var configSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch current profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		name := args[0]
		if _, err := os.Stat(filepath.Join(pmProfilesDir(pm), name+".json")); err != nil {
			return err
		}
		return pm.SetCurrent(name)
	},
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		current, _ := pm.GetCurrent()
		if current == args[0] {
			return fmt.Errorf("cannot delete current profile")
		}
		return os.Remove(filepath.Join(pmProfilesDir(pm), args[0]+".json"))
	},
}

var configRenameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a profile",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		old, new := args[0], args[1]
		return os.Rename(filepath.Join(pmProfilesDir(pm), old+".json"), filepath.Join(pmProfilesDir(pm), new+".json"))
	},
}

var configCopyCmd = &cobra.Command{
	Use:   "copy <src> <dst>",
	Short: "Copy a profile",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pm, err := profiles.NewManager()
		if err != nil {
			return err
		}
		src, dst := args[0], args[1]
		b, err := os.ReadFile(filepath.Join(pmProfilesDir(pm), src+".json"))
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(pmProfilesDir(pm), dst+".json"), b, 0o644)
	},
}

func mustJSON(cfg *config.Config) []byte {
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return b
}

func init() {
	manageCmd.AddCommand(profilesCmd)
	configCmd.AddCommand(configSaveCmd)
	configCmd.AddCommand(configSwitchCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configRenameCmd)
	configCmd.AddCommand(configCopyCmd)
}
