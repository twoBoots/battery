package techstack_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/techstack"
)

func TestIsSubBattery(t *testing.T) {
	tempDir := t.TempDir()

	assert.False(t, techstack.IsSubBattery(tempDir))

	// Create .batteryrc
	err := os.WriteFile(filepath.Join(tempDir, ".batteryrc"), []byte("{}"), 0o644)
	require.NoError(t, err)
	assert.True(t, techstack.IsSubBattery(tempDir))

	// Remove .batteryrc and create .batteryrc.json
	_ = os.Remove(filepath.Join(tempDir, ".batteryrc"))
	err = os.WriteFile(filepath.Join(tempDir, ".batteryrc.json"), []byte("{}"), 0o644)
	require.NoError(t, err)
	assert.True(t, techstack.IsSubBattery(tempDir))

	// Remove .batteryrc.json and create battery.config.json
	_ = os.Remove(filepath.Join(tempDir, ".batteryrc.json"))
	err = os.WriteFile(filepath.Join(tempDir, "battery.config.json"), []byte("{}"), 0o644)
	require.NoError(t, err)
	assert.True(t, techstack.IsSubBattery(tempDir))
}

func TestResolveBarrelTechStack_CooperDefinition(t *testing.T) {
	tempDir := t.TempDir()
	cooperDir := filepath.Join(tempDir, ".cooper", "definition")
	err := os.MkdirAll(cooperDir, 0o755)
	require.NoError(t, err)

	content := `# Tech Stack
- Go 1.23
- Cobra CLI
- Huh TUI
- SQLite
`
	err = os.WriteFile(filepath.Join(cooperDir, "tech-stack.md"), []byte(content), 0o644)
	require.NoError(t, err)

	info := techstack.ResolveBarrelTechStack(tempDir)
	assert.True(t, info.Exists)
	assert.Equal(t, filepath.Join(cooperDir, "tech-stack.md"), info.FilePath)
	assert.Equal(t, "Go 1.23 | Cobra CLI | Huh TUI", info.Summary)
	assert.Equal(t, content, info.Content)
}

func TestResolveBarrelTechStack_ConductorFallback(t *testing.T) {
	tempDir := t.TempDir()
	conductorDir := filepath.Join(tempDir, "conductor")
	err := os.MkdirAll(conductorDir, 0o755)
	require.NoError(t, err)

	content := `# Conductor Stack
* TypeScript
* React
`
	err = os.WriteFile(filepath.Join(conductorDir, "tech-stack.md"), []byte(content), 0o644)
	require.NoError(t, err)

	info := techstack.ResolveBarrelTechStack(tempDir)
	assert.True(t, info.Exists)
	assert.Equal(t, filepath.Join(conductorDir, "tech-stack.md"), info.FilePath)
	assert.Equal(t, "TypeScript | React", info.Summary)
}

func TestResolveBarrelTechStack_RootFallback(t *testing.T) {
	tempDir := t.TempDir()
	content := "Python 3.12 and FastAPI Backend"
	err := os.WriteFile(filepath.Join(tempDir, "tech-stack.md"), []byte(content), 0o644)
	require.NoError(t, err)

	info := techstack.ResolveBarrelTechStack(tempDir)
	assert.True(t, info.Exists)
	assert.Equal(t, "Python 3.12 and FastAPI Backend", info.Summary)
}

func TestResolveBarrelTechStack_Missing(t *testing.T) {
	tempDir := t.TempDir()

	info := techstack.ResolveBarrelTechStack(tempDir)
	assert.False(t, info.Exists)
	assert.Equal(t, "No Cooper tech-stack.md defined", info.Summary)
}

func TestSummarizeTechStackMarkdown_FallbackHeader(t *testing.T) {
	content := `# Only Header
`
	summary := techstack.SummarizeTechStackMarkdown(content)
	assert.Equal(t, "Cooper tech-stack defined", summary)

	content2 := `
Just plain text line without bullets
`
	summary2 := techstack.SummarizeTechStackMarkdown(content2)
	assert.Equal(t, "Just plain text line without bullets", summary2)
}
