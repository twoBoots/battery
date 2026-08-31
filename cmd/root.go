package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	Version = "1.4.2"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "battery",
	Short: "Multi-Repository SDD Orchestrator",
	Long: fmt.Sprintf(`🔋 battery - Multi-Repository SDD Orchestrator (v%s)

Coordinates cross-cutting feature epics, shared specs, and worktrees
across a collection of barrels for human developers and AI agents alike.`, Version),
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// ResetFlags recursively resets all flags on a command and its subcommands to their default values.
func ResetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, child := range cmd.Commands() {
		ResetFlags(child)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	defer ResetFlags(RootCmd)
	return RootCmd.Execute()
}

func init() {
	RootCmd.Version = Version
	RootCmd.SetVersionTemplate("battery v{{.Version}}\n")

	// Add 'ls' and 'list' aliases directly at the root for barrel listing
	RootCmd.AddCommand(listAliasCmd)
	RootCmd.AddCommand(lsAliasCmd)
}

var listAliasCmd = &cobra.Command{
	Use:    "list",
	Hidden: true,
	Short:  "Alias for 'barrel list'",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBarrelList(cmd.OutOrStdout(), ".")
	},
}

var lsAliasCmd = &cobra.Command{
	Use:    "ls",
	Hidden: true,
	Short:  "Alias for 'barrel list'",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBarrelList(cmd.OutOrStdout(), ".")
	},
}

func getWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
