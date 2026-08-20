package framework

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/*
var embeddedTemplates embed.FS

const (
	StatusUpToDate          = "up_to_date"
	StatusCustomizedLocally = "customized_locally"
	StatusOutdated          = "outdated"
	StatusMissing           = "missing"
)

// TemplateInfo describes a canonical framework asset.
type TemplateInfo struct {
	Name         string `json:"name"`
	TargetPath   string `json:"targetPath"`
	Category     string `json:"category"`
	Description  string `json:"description"`
	EmbeddedPath string `json:"-"`
}

// FileStatus represents the inspection result for a single framework file in the workspace.
type FileStatus struct {
	TemplateName          string `json:"templateName"`
	TargetPath            string `json:"targetPath"`
	Category              string `json:"category"`
	Status                string `json:"status"` // up_to_date | customized_locally | outdated | missing
	HasLocalModifications bool   `json:"hasLocalModifications"`
	Details               string `json:"details,omitempty"`
}

// FrameworkStatusReport aggregates the inspection results across the target workspace or barrel.
type FrameworkStatusReport struct {
	CLIVersion       string       `json:"cliVersion"`
	StandardsVersion string       `json:"standardsVersion"`
	Target           string       `json:"target"`
	TargetPath       string       `json:"targetPath"`
	UpToDate         bool         `json:"upToDate"`
	UpdateAvailable  bool         `json:"updateAvailable"`
	Summary          string       `json:"summary"`
	Files            []FileStatus `json:"files"`
}

var canonicalCatalog = []TemplateInfo{
	{
		Name:         "skills/cooper-setup",
		TargetPath:   ".agents/skills/cooper-setup/SKILL.md",
		Category:     "skill",
		Description:  "Cooper project initialization & infrastructure setup skill",
		EmbeddedPath: "templates/skills/cooper-setup/SKILL.md",
	},
	{
		Name:         "skills/cooper-new-track",
		TargetPath:   ".agents/skills/cooper-new-track/SKILL.md",
		Category:     "skill",
		Description:  "Cooper track inception, living spec deltas & TDD plan generation skill",
		EmbeddedPath: "templates/skills/cooper-new-track/SKILL.md",
	},
	{
		Name:         "skills/cooper-implement",
		TargetPath:   ".agents/skills/cooper-implement/SKILL.md",
		Category:     "skill",
		Description:  "Cooper TDD execution, Git notes metadata & phase checkpoint skill",
		EmbeddedPath: "templates/skills/cooper-implement/SKILL.md",
	},
	{
		Name:         "skills/cooper-review",
		TargetPath:   ".agents/skills/cooper-review/SKILL.md",
		Category:     "skill",
		Description:  "Cooper code & living specification review skill",
		EmbeddedPath: "templates/skills/cooper-review/SKILL.md",
	},
	{
		Name:         "skills/cooper-rfc",
		TargetPath:   ".agents/skills/cooper-rfc/SKILL.md",
		Category:     "skill",
		Description:  "Cooper two-tier RFC architecture & collaborative review skill",
		EmbeddedPath: "templates/skills/cooper-rfc/SKILL.md",
	},
	{
		Name:         "skills/cooper-status",
		TargetPath:   ".agents/skills/cooper-status/SKILL.md",
		Category:     "skill",
		Description:  "Cooper workspace & worktree status overview skill",
		EmbeddedPath: "templates/skills/cooper-status/SKILL.md",
	},
	{
		Name:         "docs/COOPER.md",
		TargetPath:   ".cooper/COOPER.md",
		Category:     "framework_doc",
		Description:  "Canonical Cooper SDD Framework reference manual",
		EmbeddedPath: "templates/docs/COOPER.md",
	},
	{
		Name:         "docs/BATTERY.md",
		TargetPath:   ".cooper/BATTERY.md",
		Category:     "framework_doc",
		Description:  "Battery Multi-Barrel Orchestrator architecture reference manual",
		EmbeddedPath: "templates/docs/BATTERY.md",
	},
	{
		Name:         "definition/workflow.md",
		TargetPath:   ".cooper/definition/workflow.md",
		Category:     "definition",
		Description:  "Cooper standard SDD workflow & phase definitions",
		EmbeddedPath: "templates/definition/workflow.md",
	},
}

// ListTemplates returns all registered canonical framework templates.
func ListTemplates() []TemplateInfo {
	copied := make([]TemplateInfo, len(canonicalCatalog))
	copy(copied, canonicalCatalog)
	return copied
}

// GetTemplate retrieves the raw markdown content of a canonical template by name.
func GetTemplate(name string) (string, error) {
	for _, tmpl := range canonicalCatalog {
		if tmpl.Name == name {
			data, err := embeddedTemplates.ReadFile(tmpl.EmbeddedPath)
			if err != nil {
				return "", fmt.Errorf("failed to read embedded template %q: %w", name, err)
			}
			return string(data), nil
		}
	}
	return "", fmt.Errorf("template %q not found in canonical catalog", name)
}

// normalizeContent standardizes line endings and trailing whitespace for accurate comparison.
func normalizeContent(b []byte) []byte {
	str := strings.ReplaceAll(string(b), "\r\n", "\n")
	return bytes.TrimSpace([]byte(str))
}

// InspectFrameworkStatus inspects target workspace or barrel directory against canonical framework standards.
func InspectFrameworkStatus(cwd string, barrelRelPath string, cliVersion string) (*FrameworkStatusReport, error) {
	targetDir := cwd
	targetName := "workspace_root"

	if barrelRelPath != "" {
		if filepath.IsAbs(barrelRelPath) {
			targetDir = barrelRelPath
		} else {
			targetDir = filepath.Join(cwd, barrelRelPath)
		}
		targetName = barrelRelPath
	}

	report := &FrameworkStatusReport{
		CLIVersion:       cliVersion,
		StandardsVersion: cliVersion,
		Target:           targetName,
		TargetPath:       targetDir,
		UpToDate:         true,
		UpdateAvailable:  false,
		Files:            make([]FileStatus, 0, len(canonicalCatalog)),
	}

	allMatch := true
	hasCustom := false
	hasMissing := false

	for _, tmpl := range canonicalCatalog {
		upstreamRaw, err := embeddedTemplates.ReadFile(tmpl.EmbeddedPath)
		if err != nil {
			continue
		}
		normUpstream := normalizeContent(upstreamRaw)
		upstreamHash := sha256.Sum256(normUpstream)

		targetFilePath := filepath.Join(targetDir, tmpl.TargetPath)
		localData, err := os.ReadFile(targetFilePath)

		status := StatusUpToDate
		hasLocalMod := false
		details := ""

		if err != nil {
			if os.IsNotExist(err) {
				status = StatusMissing
				details = "File is not present in workspace"
				hasMissing = true
				allMatch = false
			} else {
				status = StatusMissing
				details = fmt.Sprintf("Error reading file: %v", err)
				allMatch = false
			}
		} else {
			normLocal := normalizeContent(localData)
			localHash := sha256.Sum256(normLocal)

			if !bytes.Equal(upstreamHash[:], localHash[:]) {
				status = StatusCustomizedLocally
				hasLocalMod = true
				hasCustom = true
				details = "Locally modified or customized by team"
				allMatch = false
			}
		}

		report.Files = append(report.Files, FileStatus{
			TemplateName:          tmpl.Name,
			TargetPath:            tmpl.TargetPath,
			Category:              tmpl.Category,
			Status:                status,
			HasLocalModifications: hasLocalMod,
			Details:               details,
		})
	}

	report.UpToDate = allMatch
	report.UpdateAvailable = !allMatch

	if allMatch {
		report.Summary = fmt.Sprintf("All project standards and skills match Cooper/Battery v%s.", cliVersion)
	} else {
		var parts []string
		if hasCustom {
			parts = append(parts, "local customizations detected")
		}
		if hasMissing {
			parts = append(parts, "missing template files")
		}
		report.Summary = fmt.Sprintf("Workspace has divergence from Cooper/Battery v%s (%s). An upgrade track is recommended.", cliVersion, strings.Join(parts, ", "))
	}

	return report, nil
}
