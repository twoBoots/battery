package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/twoboots/battery/internal/updater"
)

var (
	updaterClientOverride *updater.Client

	updateCheckFlag         bool
	updateForceFlag         bool
	updateTargetVersionFlag string
	updateRepoFlag          string
	updateExecPathFlag      string
)

// SetUpdaterClient overrides the default updater client (used primarily in tests).
func SetUpdaterClient(c *updater.Client) {
	updaterClientOverride = c
}

// UpdateCmd represents the 'battery update' command.
var UpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"self-update"},
	Short:   "Update Battery CLI to the latest (or specified) version",
	Long: `Checks GitHub Releases for new versions of Battery, downloads the
appropriate platform binary, and updates the local installation in place.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := updater.Options{
			Repo:           updateRepoFlag,
			TargetVersion:  updateTargetVersionFlag,
			CurrentVersion: Version,
			ExecutablePath: updateExecPathFlag,
			Force:          updateForceFlag,
			CheckOnly:      updateCheckFlag,
		}
		return RunUpdate(cmd.OutOrStdout(), cmd.ErrOrStderr(), opts, updaterClientOverride)
	},
}

func init() {
	UpdateCmd.Flags().BoolVarP(&updateCheckFlag, "check", "c", false, "Check for available updates without applying")
	UpdateCmd.Flags().BoolVarP(&updateForceFlag, "force", "f", false, "Force re-download and reinstall even if up to date")
	UpdateCmd.Flags().StringVarP(&updateTargetVersionFlag, "target-version", "t", "", "Target a specific release version tag (e.g. v1.3.0)")
	UpdateCmd.Flags().StringVar(&updateRepoFlag, "repo", updater.DefaultRepo, "GitHub repository to check for releases")
	UpdateCmd.Flags().StringVar(&updateExecPathFlag, "exec-path", "", "Target executable path to overwrite (advanced)")
	_ = UpdateCmd.Flags().MarkHidden("exec-path")
	_ = UpdateCmd.Flags().MarkHidden("repo")

	RootCmd.AddCommand(UpdateCmd)
}

// RunUpdate performs the update operation and writes user-friendly logs to the provided output writers.
func RunUpdate(out, errOut io.Writer, opts updater.Options, client *updater.Client) error {
	if client == nil {
		client = updater.NewClient()
	}
	if opts.CurrentVersion == "" {
		opts.CurrentVersion = Version
	}

	fmt.Fprintf(out, "🔋 Checking for updates (current version: v%s)...\n", opts.CurrentVersion)

	res, err := updater.SelfUpdateWithClient(client, opts)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if opts.CheckOnly {
		if res.UpdateAvailable {
			fmt.Fprintf(out, "✨ %s\nRun 'battery update' to install the new version.\n", res.Message)
		} else {
			fmt.Fprintf(out, "✅ %s\n", res.Message)
		}
		return nil
	}

	if res.Updated {
		fmt.Fprintf(out, "✅ %s\n", res.Message)
	} else {
		fmt.Fprintf(out, "ℹ️ %s\n", res.Message)
	}

	return nil
}
