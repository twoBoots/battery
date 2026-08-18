package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/twoboots/battery/internal/config"
	"github.com/twoboots/battery/internal/discovery"
)

type InitOptions struct {
	Structure      string
	Local          bool
	NonInteractive bool
	Force          bool
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize battery configuration (.batteryrc)",
	RunE: func(cmd *cobra.Command, args []string) error {
		structure, _ := cmd.Flags().GetString("structure")
		local, _ := cmd.Flags().GetBool("local")
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
		yes, _ := cmd.Flags().GetBool("yes")
		force, _ := cmd.Flags().GetBool("force")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		opts := InitOptions{
			Structure:      structure,
			Local:          local,
			NonInteractive: nonInteractive || yes,
			Force:          force || overwrite,
		}
		return runInit(cmd.OutOrStdout(), getWorkingDir(), opts)
	},
}

func init() {
	initCmd.Flags().StringP("structure", "s", "", "Specify structure: 'multi-repo', 'monorepo', or 'custom'")
	initCmd.Flags().Bool("local", false, "Save to .batteryrc.local instead of canonical .batteryrc")
	initCmd.Flags().Bool("non-interactive", false, "Run non-interactively with auto-discovered topology")
	initCmd.Flags().BoolP("yes", "y", false, "Alias for --non-interactive")
	initCmd.Flags().BoolP("force", "f", false, "Force overwrite existing configuration")
	initCmd.Flags().Bool("overwrite", false, "Alias for --force")

	RootCmd.AddCommand(initCmd)
}

func isTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func hasExistingConfig(cwd string, isLocal bool) bool {
	filename := config.ConfigFilename
	if isLocal {
		filename = config.LocalConfigFilename
	}
	filePath := filepath.Join(cwd, filename)
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func promptExistingConfigAction(targetFilename string) (string, error) {
	var selected string
	options := []huh.Option[string]{
		huh.NewOption(fmt.Sprintf("Continue setup and preserve current config in %s (Recommended)", targetFilename), "preserve"),
		huh.NewOption(fmt.Sprintf("Overwrite %s and start clean", targetFilename), "overwrite"),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Existing Battery configuration detected in %s:", targetFilename)).
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return "preserve", err
	}

	return selected, nil
}

func runInit(out io.Writer, cwd string, opts InitOptions) error {
	isNonInteractive := opts.NonInteractive ||
		os.Getenv("CI") == "true" ||
		!isTerminal()

	fmt.Fprintln(out, "🔋 Initializing Battery Configuration...")

	targetFilename := config.ConfigFilename
	if opts.Local {
		targetFilename = config.LocalConfigFilename
	}

	if hasExistingConfig(cwd, opts.Local) {
		if isNonInteractive {
			if opts.Force {
				fmt.Fprintf(out, "  [!] Overwriting existing configuration in %s...\n", targetFilename)
			} else {
				fmt.Fprintf(out, "  [✓] Existing configuration detected in %s; preserving current configuration.\n", targetFilename)
				var existingStructure config.ProjectStructure
				var numBarrels int

				if opts.Local {
					if localCfg, err := config.LoadLocalConfig(cwd); err == nil && localCfg != nil {
						existingStructure = localCfg.Structure
						numBarrels = len(localCfg.Barrels)
					}
				} else {
					if canonCfg, err := config.LoadConfig(cwd); err == nil && canonCfg != nil {
						existingStructure = canonCfg.Structure
						numBarrels = len(canonCfg.Barrels)
					}
				}

				if existingStructure == "" {
					existingStructure = config.StructureMultiRepo
				}

				fmt.Fprintf(out, "   Structure: %s\n", existingStructure)
				fmt.Fprintf(out, "   Barrels  : %d registered\n\n", numBarrels)
				return nil
			}
		} else {
			action, err := promptExistingConfigAction(targetFilename)
			if err != nil {
				return err
			}
			if action == "preserve" {
				fmt.Fprintf(out, "  [✓] Preserving current configuration in %s.\n", targetFilename)
				var existingStructure config.ProjectStructure
				var numBarrels int

				if opts.Local {
					if localCfg, err := config.LoadLocalConfig(cwd); err == nil && localCfg != nil {
						existingStructure = localCfg.Structure
						numBarrels = len(localCfg.Barrels)
					}
				} else {
					if canonCfg, err := config.LoadConfig(cwd); err == nil && canonCfg != nil {
						existingStructure = canonCfg.Structure
						numBarrels = len(canonCfg.Barrels)
					}
				}

				if existingStructure == "" {
					existingStructure = config.StructureMultiRepo
				}

				fmt.Fprintf(out, "   Structure: %s\n", existingStructure)
				fmt.Fprintf(out, "   Barrels  : %d registered\n\n", numBarrels)

				var configureMCP bool
				confirmMCPForm := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().
							Title("Would you like to configure Battery MCP server for your AI assistant? (Cursor, Claude, Antigravity, Windsurf, VS Code)").
							Value(&configureMCP),
					),
				)
				if err := confirmMCPForm.Run(); err == nil && configureMCP {
					if err := runMCPInstall(out, cwd, "", nil, false, false); err != nil {
						fmt.Fprintf(out, "  [!] MCP configuration notice: %v\n", err)
					}
				}

				return nil
			}

			fmt.Fprintf(out, "  [!] Overwriting %s and starting clean...\n", targetFilename)
		}
	}

	discovered := discovery.DiscoverCandidateBarrels(cwd)

	var selectedStructure config.ProjectStructure
	if opts.Structure != "" {
		selectedStructure = config.ProjectStructure(opts.Structure)
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

	if _, err := config.SaveConfig(newConfig, cwd, opts.Local); err != nil {
		return err
	}

	fmt.Fprintf(out, "\n🎉 Battery initialized successfully in %s!\n", targetFilename)
	fmt.Fprintf(out, "   Structure: %s\n", selectedStructure)
	fmt.Fprintf(out, "   Barrels  : %d registered\n\n", len(finalBarrels))

	if !isNonInteractive {
		var configureMCP bool
		confirmMCPForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Would you like to configure Battery MCP server for your AI assistant? (Cursor, Claude, Antigravity, Windsurf, VS Code)").
					Value(&configureMCP),
			),
		)
		if err := confirmMCPForm.Run(); err == nil && configureMCP {
			if err := runMCPInstall(out, cwd, "", nil, false, false); err != nil {
				fmt.Fprintf(out, "  [!] MCP configuration notice: %v\n", err)
			}
		}
	}

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
