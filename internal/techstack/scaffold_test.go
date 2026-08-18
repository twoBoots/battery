package techstack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInferTechStack_Go(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module myapp\n\ngo 1.23.0\n"), 0644)
	require.NoError(t, err)

	inferred := InferTechStack(dir)
	assert.Contains(t, inferred.Language, "Go")
	assert.Contains(t, inferred.TestRunner, "go test")
	assert.Equal(t, ">80%", inferred.CoverageThreshold)
}

func TestInferTechStack_NodeVariants(t *testing.T) {
	// 1. pnpm
	dirPnpm := t.TempDir()
	_ = os.WriteFile(filepath.Join(dirPnpm, "package.json"), []byte(`{}`), 0644)
	_ = os.WriteFile(filepath.Join(dirPnpm, "pnpm-lock.yaml"), []byte(``), 0644)
	infPnpm := InferTechStack(dirPnpm)
	assert.Contains(t, infPnpm.TestRunner, "pnpm test")

	// 2. bun
	dirBun := t.TempDir()
	_ = os.WriteFile(filepath.Join(dirBun, "package.json"), []byte(`{}`), 0644)
	_ = os.WriteFile(filepath.Join(dirBun, "bun.lockb"), []byte(``), 0644)
	infBun := InferTechStack(dirBun)
	assert.Contains(t, infBun.TestRunner, "bun test")

	// 3. yarn
	dirYarn := t.TempDir()
	_ = os.WriteFile(filepath.Join(dirYarn, "package.json"), []byte(`{}`), 0644)
	_ = os.WriteFile(filepath.Join(dirYarn, "yarn.lock"), []byte(``), 0644)
	infYarn := InferTechStack(dirYarn)
	assert.Contains(t, infYarn.TestRunner, "yarn test")

	// 4. npm default
	dirNpm := t.TempDir()
	_ = os.WriteFile(filepath.Join(dirNpm, "package.json"), []byte(`{}`), 0644)
	infNpm := InferTechStack(dirNpm)
	assert.Contains(t, infNpm.TestRunner, "npm test")
}

func TestInferTechStack_Rust(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte("[package]\nname = \"core\"\n"), 0644)
	require.NoError(t, err)

	inferred := InferTechStack(dir)
	assert.Equal(t, "Rust", inferred.Language)
	assert.Contains(t, inferred.TestRunner, "cargo test")
}

func TestInferTechStack_Python(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "pyproject.toml"), []byte("[tool.pytest]\n"), 0644)
	require.NoError(t, err)

	inferred := InferTechStack(dir)
	assert.Contains(t, inferred.Language, "Python")
	assert.Contains(t, inferred.TestRunner, "pytest")
}

func TestInferTechStack_Java(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "pom.xml"), []byte("<project></project>"), 0644)
	require.NoError(t, err)

	inferred := InferTechStack(dir)
	assert.Contains(t, inferred.Language, "Java")
	assert.Contains(t, inferred.TestRunner, "mvn test")
}

func TestInferTechStack_Deno(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "deno.json"), []byte(`{"tasks":{"test":"deno test"}}`), 0644)
	require.NoError(t, err)

	inferred := InferTechStack(dir)
	assert.Contains(t, inferred.Language, "Deno")
	assert.Contains(t, inferred.TestRunner, "deno test")
}

func TestInferTechStack_GenericFallback(t *testing.T) {
	dir := t.TempDir()
	inferred := InferTechStack(dir)
	assert.Contains(t, inferred.Language, "Generic")
	assert.Contains(t, inferred.TestRunner, "make test")
}

func TestScaffoldBarrelTechStack_SuccessAndErrors(t *testing.T) {
	// Nonexistent dir
	_, err := ScaffoldBarrelTechStack("/nonexistent/path/123", ScaffoldOptions{})
	require.Error(t, err)

	dir := t.TempDir()
	err = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module myapp\n\ngo 1.23.0\n"), 0644)
	require.NoError(t, err)

	res, err := ScaffoldBarrelTechStack(dir, ScaffoldOptions{})
	require.NoError(t, err)
	assert.True(t, res.Created)
	assert.False(t, res.Overwritten)
	assert.FileExists(t, res.TechStackPath)

	content, err := os.ReadFile(res.TechStackPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Technology Stack & Platform Contracts")
	assert.Contains(t, string(content), "Go")
	assert.Contains(t, string(content), "go test")

	// Attempting to scaffold again without force should fail
	_, err = ScaffoldBarrelTechStack(dir, ScaffoldOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	// Overwrite with force and overrides
	res2, err := ScaffoldBarrelTechStack(dir, ScaffoldOptions{
		Force:             true,
		Language:          "Go 1.24",
		Framework:         "Gin",
		TestRunner:        "go test -v -race ./...",
		Linter:            "golangci-lint",
		CoverageThreshold: "90%",
	})
	require.NoError(t, err)
	assert.True(t, res2.Overwritten)

	content2, err := os.ReadFile(res2.TechStackPath)
	require.NoError(t, err)
	assert.Contains(t, string(content2), "Gin")
	assert.Contains(t, string(content2), "-race")
	assert.Contains(t, string(content2), "90%")
}
