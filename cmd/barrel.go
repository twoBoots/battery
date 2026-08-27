package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/techstack"
)

var barrelCmd = &cobra.Command{
	Use:     "barrel",
	Aliases: []string{"barrels"},
	Short:   "Manage registered barrels",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBarrelList(cmd.OutOrStdout(), getWorkingDir())
	},
}

var barrelListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all registered barrels and resolved Cooper tech stacks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBarrelList(cmd.OutOrStdout(), getWorkingDir())
	},
}

var (
	addName  string
	addType  string
	addRole  string
	addTech  string
	addDocs  string
	addJira  string
	addLocal bool
)

var barrelAddCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Add a barrel to configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pathArg := strings.TrimSpace(args[0])
		name := strings.TrimSpace(addName)
		if name == "" {
			name = config.InferBarrelName(pathArg)
		}

		barrelType := config.BarrelTypeBarrel
		if strings.ToLower(addType) == "battery" {
			barrelType = config.BarrelTypeBattery
		}

		b := config.BarrelConfig{
			Name: name,
			Path: pathArg,
			Type: barrelType,
			Role: strings.TrimSpace(addRole),
			Tech: strings.TrimSpace(addTech),
			Docs: strings.TrimSpace(addDocs),
			Jira: strings.TrimSpace(addJira),
		}

		cwd := getWorkingDir()
		if _, err := config.AddBarrel(b, cwd, addLocal); err != nil {
			return err
		}

		targetFile := config.ConfigFilename
		if addLocal {
			targetFile = config.LocalConfigFilename
		}
		fmt.Fprintf(cmd.OutOrStdout(), "✓ Added barrel '%s' (%s) to %s\n", name, pathArg, targetFile)
		return nil
	},
}

var removeLocal bool

var barrelRemoveCmd = &cobra.Command{
	Use:     "remove <name_or_path>",
	Aliases: []string{"rm"},
	Short:   "Remove a barrel from configuration",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := strings.TrimSpace(args[0])
		cwd := getWorkingDir()
		if _, err := config.RemoveBarrel(identifier, cwd, removeLocal); err != nil {
			return err
		}

		targetFile := config.ConfigFilename
		if removeLocal {
			targetFile = config.LocalConfigFilename
		}
		fmt.Fprintf(cmd.OutOrStdout(), "✓ Removed barrel '%s' from %s\n", identifier, targetFile)
		return nil
	},
}

var (
	barrelInitLang      string
	barrelInitFramework string
	barrelInitTest      string
	barrelInitLinter    string
	barrelInitCov       string
	barrelInitForce     bool
)

var barrelInitCmd = &cobra.Command{
	Use:   "init <path|name>",
	Short: "Scaffold .cooper/definition/tech-stack.md for a barrel package or repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := strings.TrimSpace(args[0])
		return runBarrelInit(cmd.OutOrStdout(), getWorkingDir(), target, techstack.ScaffoldOptions{
			Language:          barrelInitLang,
			Framework:         barrelInitFramework,
			TestRunner:        barrelInitTest,
			Linter:            barrelInitLinter,
			CoverageThreshold: barrelInitCov,
			Force:             barrelInitForce,
		})
	},
}

func init() {
	barrelAddCmd.Flags().StringVarP(&addName, "name", "n", "", "Custom name for the barrel")
	barrelAddCmd.Flags().StringVarP(&addType, "type", "t", "barrel", "Specify barrel type ('barrel' or 'battery')")
	barrelAddCmd.Flags().StringVar(&addRole, "role", "", "Domain role and responsibility description")
	barrelAddCmd.Flags().StringVar(&addTech, "tech", "", "Primary tech stack and runtime summary")
	barrelAddCmd.Flags().StringVar(&addDocs, "docs", "", "Path to orchestrator barrel documentation profile")
	barrelAddCmd.Flags().StringVar(&addJira, "jira", "", "Jira project or issue tracker mapping")
	barrelAddCmd.Flags().BoolVar(&addLocal, "local", false, "Add to local developer overrides (.batteryrc.local)")

	barrelRemoveCmd.Flags().BoolVar(&removeLocal, "local", false, "Remove from local developer overrides (.batteryrc.local)")

	barrelInitCmd.Flags().StringVar(&barrelInitLang, "language", "", "Override primary language (e.g. 'Go 1.23', 'TypeScript')")
	barrelInitCmd.Flags().StringVar(&barrelInitFramework, "framework", "", "Override framework (e.g. 'Next.js', 'Gin')")
	barrelInitCmd.Flags().StringVar(&barrelInitTest, "test-runner", "", "Override test runner command (e.g. 'go test ./...')")
	barrelInitCmd.Flags().StringVar(&barrelInitLinter, "linter", "", "Override linter command (e.g. 'golangci-lint run')")
	barrelInitCmd.Flags().StringVar(&barrelInitCov, "coverage-threshold", "", "Override coverage threshold (e.g. '>80%')")
	barrelInitCmd.Flags().BoolVarP(&barrelInitForce, "force", "f", false, "Overwrite existing tech-stack.md")

	barrelCmd.AddCommand(barrelListCmd)
	barrelCmd.AddCommand(barrelAddCmd)
	barrelCmd.AddCommand(barrelRemoveCmd)
	barrelCmd.AddCommand(barrelInitCmd)

	RootCmd.AddCommand(barrelCmd)
}

func runBarrelInit(out io.Writer, cwd, target string, opts techstack.ScaffoldOptions) error {
	targetPath := target
	effCfg, err := config.GetEffectiveConfig(cwd)
	if err == nil {
		for _, b := range effCfg.Barrels {
			if b.Name == target {
				targetPath = b.Path
				break
			}
		}
	}

	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(cwd, targetPath)
	}

	res, err := techstack.ScaffoldBarrelTechStack(targetPath, opts)
	if err != nil {
		return err
	}

	status := "Created"
	if res.Overwritten {
		status = "Updated"
	}

	fmt.Fprintf(out, "✓ %s Cooper tech stack at %s\n", status, res.TechStackPath)
	if res.StyleguidePath != "" {
		fmt.Fprintf(out, "✓ Styleguide: %s\n", res.StyleguidePath)
	}
	return nil
}

func runBarrelList(out io.Writer, cwd string) error {
	effective, err := config.GetEffectiveConfig(cwd)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "\n🛢️  Registered Barrels (%d) [Structure: %s]\n\n", len(effective.Barrels), effective.Structure)

	if len(effective.Barrels) == 0 {
		fmt.Fprintln(out, "  No barrels registered yet. Run 'battery barrel add <path>' or 'battery init'.")
		fmt.Fprintln(out)
		return nil
	}

	for _, barrel := range effective.Barrels {
		absPath := filepath.Join(cwd, barrel.Path)
		exists := false
		fi, err := os.Stat(absPath)
		if err == nil && fi.IsDir() {
			exists = true
		}

		existsTag := "✓"
		if !exists {
			existsTag = "✗ (missing)"
		}

		sourceTag := "[canonical]"
		if barrel.Source == "local" {
			sourceTag = "[local override]"
		}

		typeTag := ""
		if barrel.Type == config.BarrelTypeBattery {
			typeTag = " [sub-battery]"
		}

		subBat := ""
		techInfo := techstack.CooperTechStackInfo{Summary: "N/A"}
		if exists {
			techInfo = techstack.ResolveBarrelTechStack(absPath)
			if techstack.IsSubBattery(absPath) {
				subBat = " [contains .batteryrc]"
			}
		}

		fmt.Fprintf(out, "  • %s %s%s%s\n", barrel.Name, sourceTag, typeTag, subBat)
		fmt.Fprintf(out, "    Path   : %s (%s)\n", barrel.Path, existsTag)
		fmt.Fprintf(out, "    Cooper : %s\n\n", techInfo.Summary)
	}

	return nil
}
