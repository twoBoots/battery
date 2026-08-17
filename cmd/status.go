package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/twoboots/battery/internal/config"
	"github.com/twoboots/battery/internal/techstack"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show workspace structure and barrel connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatus(cmd.OutOrStdout(), getWorkingDir())
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

func runStatus(out io.Writer, cwd string) error {
	effective, err := config.GetEffectiveConfig(cwd)
	if err != nil {
		return err
	}

	localConfig, err := config.LoadLocalConfig(cwd)
	if err != nil {
		return err
	}

	localOverrideStatus := "None"
	if localConfig != nil {
		localOverrideStatus = "Active (.batteryrc.local)"
	}

	fmt.Fprintln(out, "\n🔋 Battery Workspace Status")
	fmt.Fprintln(out, "============================")
	fmt.Fprintf(out, "Version       : %s\n", effective.Version)
	fmt.Fprintf(out, "Structure     : %s\n", effective.Structure)
	fmt.Fprintf(out, "Local Overrides: %s\n", localOverrideStatus)
	fmt.Fprintf(out, "Total Barrels : %d\n\n", len(effective.Barrels))

	if len(effective.Barrels) == 0 {
		fmt.Fprintf(out, "No barrels currently configured in this battery.\n\n")
		return nil
	}

	validCount := 0
	for _, barrel := range effective.Barrels {
		absPath := filepath.Join(cwd, barrel.Path)
		exists := false
		fi, err := os.Stat(absPath)
		if err == nil && fi.IsDir() {
			exists = true
		}

		if exists {
			validCount++
		}

		statusIcon := "🔴"
		if exists {
			statusIcon = "🟢"
		}

		subBattery := ""
		techSummary := "Directory not found"

		if exists {
			if techstack.IsSubBattery(absPath) {
				subBattery = " [Sub-Battery Orchestrator]"
			}
			techInfo := techstack.ResolveBarrelTechStack(absPath)
			techSummary = techInfo.Summary
		}

		fmt.Fprintf(out, "  %s %s (%s)%s\n", statusIcon, barrel.Name, barrel.Path, subBattery)
		fmt.Fprintf(out, "     Tech Stack : %s\n", techSummary)
	}

	fmt.Fprintf(out, "\nSummary: %d/%d barrels connected.\n\n", validCount, len(effective.Barrels))
	return nil
}
