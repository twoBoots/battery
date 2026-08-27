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

// RegisterDefaultResources registers resources on the MCP server.
func RegisterDefaultResources(s *Server) {
	// 1. battery://topology
	s.RegisterResource(Resource{
		URI:         "battery://topology",
		Name:        "Battery Topology & Configuration",
		Description: "Merged active configuration representing canonical .batteryrc and local overrides",
		MIMEType:    "application/json",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		effCfg, err := config.GetEffectiveConfig(s.Cwd())
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to get effective config: %w", err)
		}
		for i := range effCfg.Barrels {
			b := &effCfg.Barrels[i]
			bPath := b.Path
			if !filepath.IsAbs(bPath) {
				bPath = filepath.Join(s.Cwd(), bPath)
			}
			b.AbsolutePath = bPath
			if fi, statErr := os.Stat(bPath); statErr == nil && fi.IsDir() {
				b.Exists = true
				ctxInfo := techstack.ResolveBarrelContext(s.Cwd(), bPath, *b)
				b.CooperTechStack = ctxInfo.Summary
				b.HasProfile = ctxInfo.HasProfile
				b.ProfilePath = ctxInfo.ProfilePath
			}
		}
		data, err := json.MarshalIndent(effCfg, "", "  ")
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to marshal topology: %w", err)
		}
		return ReadResourceResult{
			Contents: []ResourceContent{
				{
					URI:      uri,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})

	// 2. battery://barrels/{name}/tech-stack
	s.RegisterResource(Resource{
		URI:         "battery://barrels/{name}/tech-stack",
		Name:        "Barrel Tech Stack",
		Description: "Resolved Cooper tech stack and development guidelines for a specific barrel",
		MIMEType:    "text/markdown",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		// URI: battery://barrels/<name>/tech-stack
		prefix := "battery://barrels/"
		suffix := "/tech-stack"
		if !strings.HasPrefix(uri, prefix) || !strings.HasSuffix(uri, suffix) {
			return ReadResourceResult{}, fmt.Errorf("invalid barrel tech-stack uri: %s", uri)
		}
		name := uri[len(prefix) : len(uri)-len(suffix)]

		effCfg, err := config.GetEffectiveConfig(s.Cwd())
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to load configuration: %w", err)
		}

		var targetBarrel *config.EffectiveBarrel
		for _, b := range effCfg.Barrels {
			if b.Name == name {
				targetBarrel = &b
				break
			}
		}
		if targetBarrel == nil {
			return ReadResourceResult{}, fmt.Errorf("barrel %q not found in configuration", name)
		}

		bPath := targetBarrel.Path
		if !filepath.IsAbs(bPath) {
			bPath = filepath.Join(s.Cwd(), bPath)
		}

		ts := techstack.ResolveBarrelTechStack(bPath)
		techContent := ts.Content
		if techContent == "" {
			profile := techstack.ResolveBarrelProfile(s.Cwd(), *targetBarrel)
			if profile.Content != "" {
				techContent = profile.Content
			} else if targetBarrel.Tech != "" || targetBarrel.Role != "" {
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("# Tech Stack & Context for %s\n\n", name))
				if targetBarrel.Tech != "" {
					sb.WriteString(fmt.Sprintf("- Tech: %s\n", targetBarrel.Tech))
				}
				if targetBarrel.Role != "" {
					sb.WriteString(fmt.Sprintf("- Role: %s\n", targetBarrel.Role))
				}
				techContent = sb.String()
			} else {
				techContent = fmt.Sprintf("# Tech Stack for %s\n%s\n", name, ts.Summary)
			}
		}

		return ReadResourceResult{
			Contents: []ResourceContent{
				{
					URI:      uri,
					MIMEType: "text/markdown",
					Text:     techContent,
				},
			},
		}, nil
	})

	// 3. battery://barrels/{name}/docs
	s.RegisterResource(Resource{
		URI:         "battery://barrels/{name}/docs",
		Name:        "Barrel Documentation Profile",
		Description: "Orchestrator-level architectural profile, commands, and agent guidelines for a specific barrel",
		MIMEType:    "text/markdown",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		prefix := "battery://barrels/"
		suffix := "/docs"
		if !strings.HasPrefix(uri, prefix) || !strings.HasSuffix(uri, suffix) {
			return ReadResourceResult{}, fmt.Errorf("invalid barrel docs uri: %s", uri)
		}
		name := uri[len(prefix) : len(uri)-len(suffix)]

		effCfg, err := config.GetEffectiveConfig(s.Cwd())
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to load configuration: %w", err)
		}

		var targetBarrel *config.EffectiveBarrel
		for _, b := range effCfg.Barrels {
			if b.Name == name {
				targetBarrel = &b
				break
			}
		}
		if targetBarrel == nil {
			return ReadResourceResult{}, fmt.Errorf("barrel %q not found in configuration", name)
		}

		profile := techstack.ResolveBarrelProfile(s.Cwd(), *targetBarrel)
		docContent := profile.Content
		if docContent == "" {
			docContent = fmt.Sprintf("# Barrel Profile: %s\n\nNo barrel documentation profile found at docs/barrels/%s.md\n", name, name)
		}

		return ReadResourceResult{
			Contents: []ResourceContent{
				{
					URI:      uri,
					MIMEType: "text/markdown",
					Text:     docContent,
				},
			},
		}, nil
	})

	// 3. battery://tracks/{track_id}
	s.RegisterResource(Resource{
		URI:         "battery://tracks/{track_id}",
		Name:        "Track Status & Progress",
		Description: "Aggregated status, task counts, and participating barrels for an active or archived track",
		MIMEType:    "application/json",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		prefix := "battery://tracks/"
		if !strings.HasPrefix(uri, prefix) {
			return ReadResourceResult{}, fmt.Errorf("invalid track uri: %s", uri)
		}
		trackID := strings.TrimPrefix(uri, prefix)

		st, err := track.GetMultiBarrelTrackStatus(s.Cwd(), trackID)
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to get track status: %w", err)
		}

		data, err := json.MarshalIndent(st, "", "  ")
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to marshal track status: %w", err)
		}

		return ReadResourceResult{
			Contents: []ResourceContent{
				{
					URI:      uri,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})

	// 4. battery://framework-status
	s.RegisterResource(Resource{
		URI:         "battery://framework-status",
		Name:        "Battery Framework Alignment & Status",
		Description: "Workspace standards and agent skills alignment status against canonical Cooper/Battery templates",
		MIMEType:    "application/json",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		rep, err := framework.InspectFrameworkStatus(s.Cwd(), "", s.Version())
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to inspect framework status: %w", err)
		}
		data, err := json.MarshalIndent(rep, "", "  ")
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to marshal framework status: %w", err)
		}
		return ReadResourceResult{
			Contents: []ResourceContent{
				{
					URI:      uri,
					MIMEType: "application/json",
					Text:     string(data),
				},
			},
		}, nil
	})

	// 5. battery://templates/{name}
	s.RegisterResource(Resource{
		URI:         "battery://templates/{name}",
		Name:        "Upstream Framework Template",
		Description: "Canonical Cooper/Battery template content by name (e.g. skills/cooper-rfc, docs/BATTERY.md)",
		MIMEType:    "text/markdown",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		prefix := "battery://templates/"
		if !strings.HasPrefix(uri, prefix) {
			return ReadResourceResult{}, fmt.Errorf("invalid template uri: %s", uri)
		}
		tmplName := strings.TrimPrefix(uri, prefix)
		content, err := framework.GetTemplate(tmplName)
		if err != nil {
			return ReadResourceResult{}, fmt.Errorf("failed to get template %q: %w", tmplName, err)
		}
		return ReadResourceResult{
			Contents: []ResourceContent{
				{
					URI:      uri,
					MIMEType: "text/markdown",
					Text:     content,
				},
			},
		}, nil
	})
}
