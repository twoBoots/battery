package techstack

import (
	"strings"

	"github.com/twoBoots/battery/internal/config"
)

// BarrelContextInfo represents the resolved hybrid architectural and tech stack context.
type BarrelContextInfo struct {
	Summary       string `json:"summary"`
	Source        string `json:"source"` // "cooper", "metadata", "profile", "none"
	HasCooperSpec bool   `json:"hasCooperSpec"`
	HasProfile    bool   `json:"hasProfile"`
	ProfilePath   string `json:"profilePath,omitempty"`
	TechPath      string `json:"techPath,omitempty"`
	Tech          string `json:"tech,omitempty"`
	Role          string `json:"role,omitempty"`
	Docs          string `json:"docs,omitempty"`
	Jira          string `json:"jira,omitempty"`
}

// ResolveBarrelContext computes the hybrid context for a barrel following the fallback hierarchy.
func ResolveBarrelContext(batteryCwd, barrelAbsPath string, barrel config.EffectiveBarrel) BarrelContextInfo {
	ctx := BarrelContextInfo{
		Tech: strings.TrimSpace(barrel.Tech),
		Role: strings.TrimSpace(barrel.Role),
		Docs: strings.TrimSpace(barrel.Docs),
		Jira: strings.TrimSpace(barrel.Jira),
	}

	// 1. Check Profile existence
	profile := ResolveBarrelProfile(batteryCwd, barrel)
	if profile.Exists {
		ctx.HasProfile = true
		ctx.ProfilePath = profile.FilePath
	}

	// 2. Check Cooper living spec
	if barrelAbsPath != "" {
		cooperInfo := ResolveBarrelTechStack(barrelAbsPath)
		if cooperInfo.Exists {
			ctx.HasCooperSpec = true
			ctx.TechPath = cooperInfo.FilePath
			ctx.Source = "cooper"
			ctx.Summary = cooperInfo.Summary
			return ctx
		}
	}

	// 3. Check inline metadata
	if ctx.Tech != "" || ctx.Role != "" {
		ctx.Source = "metadata"
		if ctx.Tech != "" && ctx.Role != "" {
			ctx.Summary = ctx.Tech + " (" + ctx.Role + ")"
		} else if ctx.Tech != "" {
			ctx.Summary = ctx.Tech
		} else {
			ctx.Summary = ctx.Role
		}
		return ctx
	}

	// 4. Check profile markdown summary
	if ctx.HasProfile && profile.Summary != "" {
		ctx.Source = "profile"
		ctx.Summary = profile.Summary
		return ctx
	}

	// 5. Default fallback
	ctx.Source = "none"
	ctx.Summary = "No Cooper tech-stack or profile defined"
	return ctx
}
