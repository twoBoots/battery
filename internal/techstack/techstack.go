package techstack

import (
	"os"
	"path/filepath"
	"strings"
)

// CooperTechStackInfo holds parsed metadata regarding a barrel's tech-stack specification.
type CooperTechStackInfo struct {
	Exists   bool   `json:"exists"`
	FilePath string `json:"filePath,omitempty"`
	Summary  string `json:"summary"`
	Content  string `json:"content,omitempty"`
}

// IsSubBattery checks if the specified directory path is itself a battery orchestrator.
func IsSubBattery(dirPath string) bool {
	candidateFiles := []string{".batteryrc", ".batteryrc.json", "battery.config.json"}
	for _, filename := range candidateFiles {
		fi, err := os.Stat(filepath.Join(dirPath, filename))
		if err == nil && !fi.IsDir() {
			return true
		}
	}
	return false
}

// ResolveBarrelTechStack resolves and reads the Cooper tech-stack.md inside a barrel directory.
func ResolveBarrelTechStack(barrelDirPath string) CooperTechStackInfo {
	candidatePaths := []string{
		filepath.Join(barrelDirPath, ".cooper", "definition", "tech-stack.md"),
		filepath.Join(barrelDirPath, "conductor", "tech-stack.md"),
		filepath.Join(barrelDirPath, "tech-stack.md"),
	}

	for _, filePath := range candidatePaths {
		fi, err := os.Stat(filePath)
		if err == nil && !fi.IsDir() {
			data, err := os.ReadFile(filePath)
			if err == nil {
				content := string(data)
				summary := SummarizeTechStackMarkdown(content)
				return CooperTechStackInfo{
					Exists:   true,
					FilePath: filePath,
					Summary:  summary,
					Content:  content,
				}
			}
		}
	}

	return CooperTechStackInfo{
		Exists:  false,
		Summary: "No Cooper tech-stack.md defined",
	}
}

// SummarizeTechStackMarkdown parses markdown bullets or headers to form a concise summary.
func SummarizeTechStackMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	highlights := make([]string, 0, 3)

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			cleaned := strings.TrimSpace(strings.TrimLeft(line, "-*"))
			if len(cleaned) > 0 {
				highlights = append(highlights, cleaned)
			}
		}
		if len(highlights) >= 3 {
			break
		}
	}

	if len(highlights) > 0 {
		return strings.Join(highlights, " | ")
	}

	// Fallback: first non-empty header or line that does not start with '#'
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			return line
		}
	}

	return "Cooper tech-stack defined"
}
