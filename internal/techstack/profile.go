package techstack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/twoBoots/battery/internal/config"
)

// BarrelProfileInfo holds parsed metadata regarding an orchestrator-level barrel profile.
type BarrelProfileInfo struct {
	Exists   bool   `json:"exists"`
	FilePath string `json:"filePath,omitempty"`
	Summary  string `json:"summary,omitempty"`
	Content  string `json:"content,omitempty"`
}

// ResolveBarrelProfile searches for a profile markdown document for the given barrel.
func ResolveBarrelProfile(batteryCwd string, barrel config.EffectiveBarrel) BarrelProfileInfo {
	candidatePaths := make([]string, 0, 3)

	if strings.TrimSpace(barrel.Docs) != "" {
		docPath := strings.TrimSpace(barrel.Docs)
		if !filepath.IsAbs(docPath) {
			docPath = filepath.Join(batteryCwd, docPath)
		}
		candidatePaths = append(candidatePaths, docPath)
	}

	candidatePaths = append(candidatePaths,
		filepath.Join(batteryCwd, "docs", "barrels", barrel.Name+".md"),
		filepath.Join(batteryCwd, ".cooper", "barrels", barrel.Name+".md"),
	)

	for _, filePath := range candidatePaths {
		fi, err := os.Stat(filePath)
		if err == nil && !fi.IsDir() {
			data, err := os.ReadFile(filePath)
			if err == nil {
				content := string(data)
				summary := SummarizeTechStackMarkdown(content)
				return BarrelProfileInfo{
					Exists:   true,
					FilePath: filePath,
					Summary:  summary,
					Content:  content,
				}
			}
		}
	}

	return BarrelProfileInfo{
		Exists: false,
	}
}

// ScaffoldBarrelProfile generates docs/barrels/<barrelName>.md with a starter profile template.
func ScaffoldBarrelProfile(batteryCwd, barrelName string, force bool) (string, error) {
	name := strings.TrimSpace(barrelName)
	if name == "" {
		return "", fmt.Errorf("barrel name cannot be empty")
	}

	targetDir := filepath.Join(batteryCwd, "docs", "barrels")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	targetFile := filepath.Join(targetDir, name+".md")
	if fileExists(targetFile) && !force {
		return "", fmt.Errorf("barrel profile already exists at %s (use --force to overwrite)", targetFile)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Barrel Profile: %s\n\n", name))
	sb.WriteString("## Role & Responsibilities\n")
	sb.WriteString("<!-- Concise description of domain responsibilities and purpose -->\n\n")
	sb.WriteString("## Tech Stack & Runtime\n")
	sb.WriteString("- Language / Runtime: \n")
	sb.WriteString("- Build System: \n")
	sb.WriteString("- Key Dependencies: \n\n")
	sb.WriteString("## Development & Build Commands\n")
	sb.WriteString("- Build: \n")
	sb.WriteString("- Test: \n")
	sb.WriteString("- Lint: \n\n")
	sb.WriteString("## Interface Contracts & Integration\n")
	sb.WriteString("- Interfaces: \n")
	sb.WriteString("- Upstream / Downstream Barrels: \n\n")
	sb.WriteString("## AI Agent Guidelines\n")
	sb.WriteString("- Notes for AI coding assistants regarding directory layout, constraints, and operational patterns.\n")

	if err := os.WriteFile(targetFile, []byte(sb.String()), 0o644); err != nil {
		return "", fmt.Errorf("failed to write barrel profile %s: %w", targetFile, err)
	}

	return targetFile, nil
}
