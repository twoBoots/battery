package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/track"
)

var trackCmd = &cobra.Command{
	Use:     "track",
	Aliases: []string{"tracks"},
	Short:   "Manage multi-barrel tracks and spec dispatch",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTrackList(cmd.OutOrStdout(), getWorkingDir())
	},
}

var (
	trackInitName         string
	trackInitType         string
	trackInitBarrels      string
	trackInitCapabilities string
)

var trackInitCmd = &cobra.Command{
	Use:   "init <track_id>",
	Short: "Initialize a new multi-barrel track in battery",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		trackID := strings.TrimSpace(args[0])
		cwd := getWorkingDir()

		var barrels []string
		if strings.TrimSpace(trackInitBarrels) != "" {
			for _, b := range strings.Split(trackInitBarrels, ",") {
				b = strings.TrimSpace(b)
				if b != "" {
					barrels = append(barrels, b)
				}
			}
		}

		var caps []string
		if strings.TrimSpace(trackInitCapabilities) != "" {
			for _, c := range strings.Split(trackInitCapabilities, ",") {
				c = strings.TrimSpace(c)
				if c != "" {
					caps = append(caps, c)
				}
			}
		}

		return runTrackInit(cmd.OutOrStdout(), cwd, trackID, trackInitName, trackInitType, barrels, caps)
	},
}

var dispatchForce bool

var trackDispatchCmd = &cobra.Command{
	Use:   "dispatch <track_id>",
	Short: "Dispatch track specs and metadata to target barrels",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		trackID := strings.TrimSpace(args[0])
		cwd := getWorkingDir()
		return runTrackDispatch(cmd.OutOrStdout(), cwd, trackID, dispatchForce)
	},
}

var trackStatusCmd = &cobra.Command{
	Use:   "status [<track_id>]",
	Short: "Show multi-barrel execution status for a track",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd := getWorkingDir()
		if len(args) == 0 {
			return runTrackList(cmd.OutOrStdout(), cwd)
		}
		return runTrackStatus(cmd.OutOrStdout(), cwd, strings.TrimSpace(args[0]))
	},
}

var trackListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all active and archived tracks in battery",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTrackList(cmd.OutOrStdout(), getWorkingDir())
	},
}

func init() {
	trackInitCmd.Flags().StringVarP(&trackInitName, "name", "n", "", "Human readable title for the track")
	trackInitCmd.Flags().StringVarP(&trackInitType, "type", "t", "feature", "Track type ('feature', 'fix', 'refactor', 'chore')")
	trackInitCmd.Flags().StringVarP(&trackInitBarrels, "barrels", "b", "", "Comma-separated list of target barrels (defaults to all)")
	trackInitCmd.Flags().StringVarP(&trackInitCapabilities, "caps", "c", "", "Comma-separated list of capabilities")

	trackDispatchCmd.Flags().BoolVarP(&dispatchForce, "force", "f", false, "Overwrite existing barrel spec deltas")

	trackCmd.AddCommand(trackInitCmd)
	trackCmd.AddCommand(trackDispatchCmd)
	trackCmd.AddCommand(trackStatusCmd)
	trackCmd.AddCommand(trackListCmd)

	RootCmd.AddCommand(trackCmd)
}

func runTrackInit(out io.Writer, cwd, trackID, name, trackType string, barrels, caps []string) error {
	// If no barrels specified, default to all registered barrels in .batteryrc
	if len(barrels) == 0 {
		effCfg, err := config.GetEffectiveConfig(cwd)
		if err == nil {
			for _, b := range effCfg.Barrels {
				barrels = append(barrels, b.Name)
			}
		}
	}

	opts := track.InitTrackOptions{
		Name:         name,
		Type:         track.TrackType(trackType),
		Capabilities: caps,
	}

	meta, err := track.InitTrack(cwd, trackID, barrels, opts)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "✓ Initialized track '%s' (%s)\n", meta.TrackID, meta.Name)
	fmt.Fprintf(out, "  Directory : .cooper/active/%s\n", meta.TrackID)
	fmt.Fprintf(out, "  Barrels   : %s\n", strings.Join(meta.Barrels, ", "))
	fmt.Fprintf(out, "\nRun 'battery track dispatch %s' to seed specs into target barrels.\n", meta.TrackID)
	return nil
}

func runTrackDispatch(out io.Writer, cwd, trackID string, force bool) error {
	results, err := track.DispatchTrack(cwd, trackID, track.DispatchTrackOptions{Force: force})
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "\n🚀 Dispatched track '%s' across %d barrel(s):\n\n", trackID, len(results))
	for _, r := range results {
		if r.Created {
			fmt.Fprintf(out, "  🟢 %s: Spec seeded at %s (plan.md omitted for local planning)\n", r.BarrelName, r.TargetDir)
		} else {
			fmt.Fprintf(out, "  🔴 %s: Failed - %s\n", r.BarrelName, r.Error)
		}
	}
	fmt.Fprintln(out)
	return nil
}

func runTrackStatus(out io.Writer, cwd, trackID string) error {
	status, err := track.GetMultiBarrelTrackStatus(cwd, trackID)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "\n📊 Multi-Barrel Track Status: %s (%s)\n", status.TrackID, status.Name)
	fmt.Fprintln(out, "==========================================================")
	fmt.Fprintf(out, "Overall Status : %s\n", status.Status)
	fmt.Fprintf(out, "Barrels Count  : %d\n\n", len(status.Barrels))

	for _, b := range status.Barrels {
		icon := "🟢"
		if b.Location == track.LocationMissing {
			icon = "🔴"
		} else if b.Status == track.TrackStatusPlanning {
			icon = "🟡"
		}

		fmt.Fprintf(out, "  %s %s [%s]\n", icon, b.BarrelName, b.Location)
		fmt.Fprintf(out, "     Status     : %s\n", b.Status)
		if b.ActivePlanTasks > 0 {
			fmt.Fprintf(out, "     Progress   : %d%% (%d/%d tasks completed)\n", b.PercentComplete(), b.CompletedTasks, b.ActivePlanTasks)
		} else if b.Location == track.LocationArchive {
			fmt.Fprintln(out, "     Progress   : 100% (Completed & Archived)")
		} else {
			fmt.Fprintln(out, "     Progress   : Planning (plan.md not yet started)")
		}
		if len(b.SpecDeltas) > 0 {
			fmt.Fprintf(out, "     Spec Deltas: %s\n", strings.Join(b.SpecDeltas, ", "))
		}
		fmt.Fprintln(out)
	}

	return nil
}

func runTrackList(out io.Writer, cwd string) error {
	tracks, err := track.ListTracks(cwd)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "\n📋 Battery Tracks (%d)\n\n", len(tracks))
	if len(tracks) == 0 {
		fmt.Fprintln(out, "  No active or archived tracks found. Run 'battery track init <track_id>' to start a track.")
		fmt.Fprintln(out)
		return nil
	}

	for _, t := range tracks {
		statusTag := fmt.Sprintf("[%s]", t.Status)
		barrels := strings.Join(t.Barrels, ", ")
		if barrels == "" {
			barrels = "none"
		}
		fmt.Fprintf(out, "  • %-28s %-12s Barrels: %s\n", t.TrackID, statusTag, barrels)
	}
	fmt.Fprintln(out)
	return nil
}
