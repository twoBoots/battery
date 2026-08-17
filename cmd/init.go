package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/twoboots/battery/internal/config"
	"github.com/twoboots/battery/internal/discovery"
)

var (
	initStructure      string
	initLocal          bool
	initNonInteractive bool
	initYes            bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize battery configuration (.batteryrc)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(cmd.OutOrStdout(), getWorkingDir())
	},
}

func init() {
	initCmd.Flags().StringVarP(&initStructure, "structure", "s", "", "Specify structure: 'multi-repo', 'monorepo', or 'custom'")
	initCmd.Flags().BoolVar(&initLocal, "local", false, "Save to .batteryrc.local instead of canonical .batteryrc")
	initCmd.Flags().BoolVar(&initNonInteractive, "non-interactive", false, "Run non-interactively with auto-discovered topology")
	initCmd.Flags().BoolVarP(&initYes, "yes", "y", false, "Alias for --non-interactive")

	RootCmd.AddCommand(initCmd)
}

func isTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func runInit(out io.Writer, cwd string) error {
	isNonInteractive := initNonInteractive ||
		initYes ||
		os.Getenv("CI") == "true" ||
		!isTerminal()

	fmt.Fprintln(out, "🔋 Initializing Battery Configuration...")

	discovered := discovery.DiscoverCandidateBarrels(cwd)

	var selectedStructure config.ProjectStructure
	if initStructure != "" {
		selectedStructure = config.ProjectStructure(initStructure)
	} else if isNonInteractive {
		selectedStructure = discovered.Structure
	} else {
		structure, err := promptStructure(discovered.Structure)
		if err != nil {
			return err
		}
		selectedStructure = structure
	}

	var finalBarrels []config.BarrelConfig

	if isNonInteractive {
		fmt.Fprintf(out, "  [✓] Topology selected: %s\n", selectedStructure)
		finalBarrels = discovered.Barrels
		if len(finalBarrels) > 0 {
			names := make([]string, len(finalBarrels))
			for i, b := range finalBarrels {
				names[i] = b.Name
			}
			fmt.Fprintf(out, "  [✓] Discovered %d barrels: %s\n", len(finalBarrels), strings.Join(names, ", "))
		} else {
			fmt.Fprintln(out, "  [✓] No candidate barrels automatically discovered.")
		}
	} else {
		barrels, err := promptBarrels(out, selectedStructure, discovered.Barrels)
		if err != nil {
			return err
		}
		finalBarrels = barrels
	}

	newConfig := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: selectedStructure,
		Barrels:   finalBarrels,
	}

	targetFilename := config.ConfigFilename
	if initLocal {
		targetFilename = config.LocalConfigFilename
	}

	if _, err := config.SaveConfig(newConfig, cwd, initLocal); err != nil {
		return err
	}

	fmt.Fprintf(out, "\n🎉 Battery initialized successfully in %s!\n", targetFilename)
	fmt.Fprintf(out, "   Structure: %s\n", selectedStructure)
	fmt.Fprintf(out, "   Barrels  : %d registered\n\n", len(finalBarrels))

	return nil
}

func promptStructure(detected config.ProjectStructure) (config.ProjectStructure, error) {
	var selected string

	multiRepoLabel := "Multi-Repo (Barrels in sibling repositories)"
	if detected == config.StructureMultiRepo {
		multiRepoLabel += " [Detected]"
	}

	monorepoLabel := "Monorepo   (Barrels in subdirectories / packages)"
	if detected == config.StructureMonorepo {
		monorepoLabel += " [Detected]"
	}

	customLabel := "Custom     (Heterogeneous, hybrid, or nested layout)"
	if detected == config.StructureCustom {
		customLabel += " [Detected]"
	}

	options := []huh.Option[string]{
		huh.NewOption(multiRepoLabel, string(config.StructureMultiRepo)),
		huh.NewOption(monorepoLabel, string(config.StructureMonorepo)),
		huh.NewOption(customLabel, string(config.StructureCustom)),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the project structure for this battery:").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return detected, err
	}

	return config.ProjectStructure(selected), nil
}

func promptBarrels(out io.Writer, _ config.ProjectStructure, discovered []config.BarrelConfig) ([]config.BarrelConfig, error) {
	barrels := make([]config.BarrelConfig, 0)

	if len(discovered) > 0 {
		var selectedPaths []string
		options := make([]huh.Option[string], len(discovered))
		for i, b := range discovered {
			options[i] = huh.NewOption(fmt.Sprintf("%s (%s)", b.Name, b.Path), b.Path)
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select candidate barrel folders to include (Space to toggle, Enter to confirm):").
					Options(options...).
					Value(&selectedPaths),
			),
		)

		if err := form.Run(); err != nil {
			return nil, err
		}

		for _, path := range selectedPaths {
			for _, b := range discovered {
				if b.Path == path {
					barrels = append(barrels, b)
					break
				}
			}
		}

		fmt.Fprintf(out, "  [✓] Selected %d barrel(s).\n", len(barrels))
	} else {
		fmt.Fprintln(out, "  [i] No candidate barrels automatically discovered.")
	}

	var addCustom bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to add any additional custom barrel paths?").
				Value(&addCustom),
		),
	)

	if err := confirmForm.Run(); err != nil {
		return barrels, nil
	}

	for addCustom {
		var customPath string
		pathForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Enter barrel path (Leave blank to finish):").
					Value(&customPath),
			),
		)
		if err := pathForm.Run(); err != nil {
			break
		}

		trimmedPath := strings.TrimSpace(customPath)
		if trimmedPath == "" {
			break
		}

		inferredName := config.InferBarrelName(trimmedPath)
		var customName string
		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(fmt.Sprintf("Name for '%s':", trimmedPath)).
					Value(&customName),
			),
		)
		if err := nameForm.Run(); err != nil {
			break
		}

		finalName := strings.TrimSpace(customName)
		if finalName == "" {
			finalName = inferredName
		}

		barrels = append(barrels, config.BarrelConfig{
			Name: finalName,
			Path: trimmedPath,
		})
		fmt.Fprintf(out, "  [✓] Added '%s' (%s)\n", finalName, trimmedPath)
	}

	return barrels, nil
}
