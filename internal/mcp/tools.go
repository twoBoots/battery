package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/twoBoots/battery/internal/config"
	"github.com/twoBoots/battery/internal/framework"
	"github.com/twoBoots/battery/internal/techstack"
	"github.com/twoBoots/battery/internal/track"
)

// RegisterDefaultTools registers all battery orchestration tools onto the MCP server.
func RegisterDefaultTools(s *Server) {
	// 1. battery_status
	s.RegisterTool(Tool{
		Name:        "battery_status",
		Description: "Returns workspace topology, merged canonical and local barrel configuration, barrel connectivity, active tracks, and framework alignment status.",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"verbose": map[string]interface{}{
					"type":        "boolean",
					"description": "Include additional details in the report",
				},
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		effCfg, err := config.GetEffectiveConfig(s.Cwd())
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to get effective config: %v", err)), nil
		}

		tracks, _ := track.ListTracks(s.Cwd())
		fwReport, _ := framework.InspectFrameworkStatus(s.Cwd(), "", s.Version())

		type StatusReport struct {
			Structure       config.ProjectStructure          `json:"structure"`
			Version         string                           `json:"version"`
			CLIVersion      string                           `json:"cli_version"`
			ConfigVersion   string                           `json:"config_version"`
			BarrelsCount    int                              `json:"barrels_count"`
			Barrels         []config.EffectiveBarrel         `json:"barrels"`
			ActiveTracks    []track.TrackMetadata            `json:"active_tracks"`
			FrameworkStatus *framework.FrameworkStatusReport `json:"framework_status,omitempty"`
		}

		report := StatusReport{
			Structure:       effCfg.Structure,
			Version:         effCfg.Version,
			CLIVersion:      s.Version(),
			ConfigVersion:   effCfg.Version,
			BarrelsCount:    len(effCfg.Barrels),
			Barrels:         effCfg.Barrels,
			ActiveTracks:    tracks,
			FrameworkStatus: fwReport,
		}

		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to marshal status: %v", err)), nil
		}
		return NewTextResult(string(data), false), nil
	})

	// 2. battery_list_barrels
	s.RegisterTool(Tool{
		Name:        "battery_list_barrels",
		Description: "Lists all registered barrels in the battery, their paths, connectivity status, and resolved Cooper tech stacks.",
		InputSchema: ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		effCfg, err := config.GetEffectiveConfig(s.Cwd())
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to get effective config: %v", err)), nil
		}

		type BarrelDetail struct {
			Name            string            `json:"name"`
			Path            string            `json:"path"`
			Type            config.BarrelType `json:"type,omitempty"`
			Source          string            `json:"source"`
			Exists          bool              `json:"exists"`
			CooperTechStack string            `json:"cooper_tech_stack,omitempty"`
		}

		var list []BarrelDetail
		for _, b := range effCfg.Barrels {
			bPath := b.Path
			if !filepath.IsAbs(bPath) {
				bPath = filepath.Join(s.Cwd(), bPath)
			}
			_, statErr := os.Stat(bPath)
			exists := statErr == nil

			tsSummary := ""
			if exists {
				ts := techstack.ResolveBarrelTechStack(bPath)
				tsSummary = ts.Summary
			}

			list = append(list, BarrelDetail{
				Name:            b.Name,
				Path:            b.Path,
				Type:            b.Type,
				Source:          b.Source,
				Exists:          exists,
				CooperTechStack: tsSummary,
			})
		}

		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to marshal barrels list: %v", err)), nil
		}
		return NewTextResult(string(data), false), nil
	})

	// 3. battery_init_track
	s.RegisterTool(Tool{
		Name:        "battery_init_track",
		Description: "Initializes a new multi-barrel track in .cooper/active/<track_id>/ with proposal, design, spec deltas, and plan.",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"track_id": map[string]interface{}{
					"type":        "string",
					"description": "Unique identifier for the track (e.g. track_auth_20260817)",
				},
				"barrels": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": "Target barrel names to participate in this track",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Human-readable name or title of the track",
				},
				"force": map[string]interface{}{
					"type":        "boolean",
					"description": "Force overwrite if track directory already exists",
				},
			},
			Required: []string{"track_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		trackID, _ := args["track_id"].(string)
		if strings.TrimSpace(trackID) == "" {
			return NewErrorResult("track_id is required"), nil
		}

		var barrels []string
		if rawBarrels, ok := args["barrels"].([]interface{}); ok {
			for _, b := range rawBarrels {
				if str, ok := b.(string); ok && str != "" {
					barrels = append(barrels, str)
				}
			}
		}

		name, _ := args["name"].(string)
		force, _ := args["force"].(bool)

		meta, err := track.InitTrack(s.Cwd(), trackID, barrels, track.InitTrackOptions{
			Name:  name,
			Force: force,
		})
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to initialize track: %v", err)), nil
		}

		return NewTextResult(fmt.Sprintf("Initialized track %q in .cooper/active/%s (participating barrels: %s)", meta.TrackID, meta.TrackID, strings.Join(meta.Barrels, ", ")), false), nil
	})

	// 4. battery_dispatch_track
	s.RegisterTool(Tool{
		Name:        "battery_dispatch_track",
		Description: "Dispatches track spec deltas and contracts into target barrel worktrees, omitting plan.md to preserve local planning autonomy.",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"track_id": map[string]interface{}{
					"type":        "string",
					"description": "Track identifier to dispatch",
				},
				"force": map[string]interface{}{
					"type":        "boolean",
					"description": "Force overwrite existing spec deltas in barrels",
				},
			},
			Required: []string{"track_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		trackID, _ := args["track_id"].(string)
		if strings.TrimSpace(trackID) == "" {
			return NewErrorResult("track_id is required"), nil
		}
		force, _ := args["force"].(bool)

		results, err := track.DispatchTrack(s.Cwd(), trackID, track.DispatchTrackOptions{
			Force: force,
		})
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to dispatch track: %v", err)), nil
		}

		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to marshal dispatch results: %v", err)), nil
		}
		return NewTextResult(fmt.Sprintf("Dispatched track %s to %d barrel(s):\n%s", trackID, len(results), string(data)), false), nil
	})

	// 5. battery_track_status
	s.RegisterTool(Tool{
		Name:        "battery_track_status",
		Description: "Inspects and aggregates track execution progress and task completion percentages across all participating barrels.",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"track_id": map[string]interface{}{
					"type":        "string",
					"description": "Track ID to query",
				},
			},
			Required: []string{"track_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		trackID, _ := args["track_id"].(string)
		if strings.TrimSpace(trackID) == "" {
			return NewErrorResult("track_id is required"), nil
		}

		st, err := track.GetMultiBarrelTrackStatus(s.Cwd(), trackID)
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to get track status: %v", err)), nil
		}

		data, err := json.MarshalIndent(st, "", "  ")
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to marshal track status: %v", err)), nil
		}
		return NewTextResult(string(data), false), nil
	})

	// 6. battery_init_barrel_tech_stack
	s.RegisterTool(Tool{
		Name:        "battery_init_barrel_tech_stack",
		Description: "Scaffolds or updates .cooper/definition/tech-stack.md and code styleguides for a specific barrel package or repository.",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"barrel": map[string]interface{}{
					"type":        "string",
					"description": "Barrel identifier (registered name or filesystem path to barrel)",
				},
				"language": map[string]interface{}{
					"type":        "string",
					"description": "Primary programming language (e.g. 'Go 1.23', 'TypeScript')",
				},
				"framework": map[string]interface{}{
					"type":        "string",
					"description": "Framework name (e.g. 'Next.js', 'Gin')",
				},
				"test_runner": map[string]interface{}{
					"type":        "string",
					"description": "Test runner command (e.g. 'go test ./...', 'npm test')",
				},
				"linter": map[string]interface{}{
					"type":        "string",
					"description": "Linter command (e.g. 'golangci-lint run', 'eslint .')",
				},
				"coverage_threshold": map[string]interface{}{
					"type":        "string",
					"description": "Target coverage threshold (e.g. '>80%')",
				},
				"force": map[string]interface{}{
					"type":        "boolean",
					"description": "Overwrite existing tech-stack.md if present",
				},
			},
			Required: []string{"barrel"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		barrelArg, _ := args["barrel"].(string)
		if strings.TrimSpace(barrelArg) == "" {
			return NewErrorResult("barrel argument is required"), nil
		}

		targetPath := barrelArg
		effCfg, err := config.GetEffectiveConfig(s.Cwd())
		if err == nil {
			for _, b := range effCfg.Barrels {
				if b.Name == barrelArg {
					targetPath = b.Path
					break
				}
			}
		}

		if !filepath.IsAbs(targetPath) {
			targetPath = filepath.Join(s.Cwd(), targetPath)
		}

		lang, _ := args["language"].(string)
		frameworkName, _ := args["framework"].(string)
		testRunner, _ := args["test_runner"].(string)
		linter, _ := args["linter"].(string)
		cov, _ := args["coverage_threshold"].(string)
		force, _ := args["force"].(bool)

		res, err := techstack.ScaffoldBarrelTechStack(targetPath, techstack.ScaffoldOptions{
			Language:          lang,
			Framework:         frameworkName,
			TestRunner:        testRunner,
			Linter:            linter,
			CoverageThreshold: cov,
			Force:             force,
		})
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to scaffold barrel tech stack: %v", err)), nil
		}

		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to marshal scaffold result: %v", err)), nil
		}
		return NewTextResult(string(data), false), nil
	})

	// 7. battery_framework_status
	s.RegisterTool(Tool{
		Name:        "battery_framework_status",
		Description: "Inspects workspace or barrel standards against canonical Cooper/Battery framework templates to detect version alignment and local customizations.",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"barrel": map[string]interface{}{
					"type":        "string",
					"description": "Optional registered barrel name or directory path. If omitted, inspects workspace root.",
				},
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		barrelArg, _ := args["barrel"].(string)
		targetRelPath := barrelArg
		if barrelArg != "" {
			effCfg, err := config.GetEffectiveConfig(s.Cwd())
			if err == nil {
				for _, b := range effCfg.Barrels {
					if b.Name == barrelArg {
						targetRelPath = b.Path
						break
					}
				}
			}
		}

		rep, err := framework.InspectFrameworkStatus(s.Cwd(), targetRelPath, s.Version())
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to inspect framework status: %v", err)), nil
		}

		data, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to marshal framework status: %v", err)), nil
		}
		return NewTextResult(string(data), false), nil
	})

	// 8. battery_get_template
	s.RegisterTool(Tool{
		Name:        "battery_get_template",
		Description: "Retrieves the upstream canonical markdown content of a Cooper skill or framework document by template name (e.g. 'skills/cooper-rfc', 'docs/BATTERY.md').",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Template identifier from canonical catalog (e.g. 'skills/cooper-review', 'docs/COOPER.md')",
				},
			},
			Required: []string{"name"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		name, _ := args["name"].(string)
		name = strings.TrimSpace(name)
		if name == "" {
			return NewErrorResult("name is required"), nil
		}

		content, err := framework.GetTemplate(name)
		if err != nil {
			return NewErrorResult(fmt.Sprintf("failed to get template %q: %v", name, err)), nil
		}

		return NewTextResult(content, false), nil
	})
}
