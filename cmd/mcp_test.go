package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/internal/config"
	"github.com/twoboots/battery/internal/mcp"
)

func TestMCPCmd_Execution(t *testing.T) {
	tempDir := t.TempDir()

	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "core", Path: "./barrels/core"},
		},
	}
	_, err := config.SaveConfig(cfg, tempDir, false)
	require.NoError(t, err)

	input := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}` + "\n" +
		`{"jsonrpc":"2.0","id":3,"method":"ping"}` + "\n"

	inBuf := bytes.NewBufferString(input)
	var outBuf bytes.Buffer

	err = runMCPServer(inBuf, &outBuf, tempDir)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	require.Len(t, lines, 3)

	var initResp mcp.Response
	err = json.Unmarshal([]byte(lines[0]), &initResp)
	require.NoError(t, err)
	assert.Nil(t, initResp.Error)

	var listToolsResp mcp.Response
	err = json.Unmarshal([]byte(lines[1]), &listToolsResp)
	require.NoError(t, err)
	assert.Nil(t, listToolsResp.Error)

	var pingResp mcp.Response
	err = json.Unmarshal([]byte(lines[2]), &pingResp)
	require.NoError(t, err)
	assert.Nil(t, pingResp.Error)
}

func TestMCPServe_Alias(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	input := `{"jsonrpc":"2.0","id":1,"method":"ping"}` + "\n"
	inBuf := bytes.NewBufferString(input)
	var outBuf bytes.Buffer

	err := runMCPServer(inBuf, &outBuf, tempDir)
	require.NoError(t, err)
	assert.Contains(t, outBuf.String(), `"result":{}`)
}

func TestMCPCmd_Help(t *testing.T) {
	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"mcp", "--help"})

	err := RootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Model Context Protocol")
	assert.Contains(t, out.String(), "--transport")
	_ = RootCmd.Flags().Set("help", "false")
	_ = mcpCmd.Flags().Set("help", "false")
}

func TestMCPInstallCmd_Help(t *testing.T) {
	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"mcp", "install", "--help"})

	err := RootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Supported clients")
	assert.Contains(t, out.String(), "--client")
	assert.Contains(t, out.String(), "--all")
	_ = RootCmd.Flags().Set("help", "false")
	_ = mcpCmd.Flags().Set("help", "false")
	_ = mcpInstallCmd.Flags().Set("help", "false")
}

func TestMCPInstall_WithClientsFlag(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	workspaceDir := filepath.Join(tempDir, "workspace")
	require.NoError(t, os.MkdirAll(homeDir, 0755))
	require.NoError(t, os.MkdirAll(workspaceDir, 0755))

	var out bytes.Buffer
	err := runMCPInstall(&out, workspaceDir, homeDir, []string{"cursor", "antigravity"}, false, true)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "Configuring Battery MCP Server")
	assert.Contains(t, out.String(), "Cursor IDE")
	assert.Contains(t, out.String(), "Google Antigravity")

	// Verify Cursor file exists
	cursorFile := filepath.Join(workspaceDir, ".cursor", "mcp.json")
	assert.FileExists(t, cursorFile)

	// Verify Antigravity file exists
	antigravityFile := filepath.Join(homeDir, ".gemini", "config", "mcp_config.json")
	assert.FileExists(t, antigravityFile)
}

func TestMCPInstall_AllFlag(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	workspaceDir := filepath.Join(tempDir, "workspace")
	require.NoError(t, os.MkdirAll(homeDir, 0755))
	require.NoError(t, os.MkdirAll(workspaceDir, 0755))

	var out bytes.Buffer
	err := runMCPInstall(&out, workspaceDir, homeDir, nil, true, true)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "MCP configuration completed")
	assert.FileExists(t, filepath.Join(workspaceDir, ".cursor", "mcp.json"))
	assert.FileExists(t, filepath.Join(workspaceDir, ".vscode", "mcp.json"))
	assert.FileExists(t, filepath.Join(homeDir, ".claude.json"))
}

func TestMCPInstall_NonInteractiveNoDetections(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	workspaceDir := filepath.Join(tempDir, "workspace")
	require.NoError(t, os.MkdirAll(homeDir, 0755))
	require.NoError(t, os.MkdirAll(workspaceDir, 0755))

	var out bytes.Buffer
	err := runMCPInstall(&out, workspaceDir, homeDir, nil, false, true)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "No AI assistant configurations automatically detected")
}

func TestMCPInstall_NonInteractiveWithDetections(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	workspaceDir := filepath.Join(tempDir, "workspace")
	require.NoError(t, os.MkdirAll(homeDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(workspaceDir, ".cursor"), 0755))

	var out bytes.Buffer
	err := runMCPInstall(&out, workspaceDir, homeDir, nil, false, true)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Configuring Battery MCP Server")
	assert.Contains(t, out.String(), "Cursor IDE")
	assert.FileExists(t, filepath.Join(workspaceDir, ".cursor", "mcp.json"))
}

func TestMCPInstall_CLIExecution(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	_ = RootCmd.Flags().Set("help", "false")
	_ = mcpCmd.Flags().Set("help", "false")
	_ = mcpInstallCmd.Flags().Set("help", "false")

	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"mcp", "install", "--client", "cursor,vscode", "--non-interactive"})

	err := RootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Configuring Battery MCP Server")
	assert.FileExists(t, filepath.Join(tempDir, ".cursor", "mcp.json"))
	assert.FileExists(t, filepath.Join(tempDir, ".vscode", "mcp.json"))
}

func TestMCPInstall_CLIExecution_All(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	_ = RootCmd.Flags().Set("help", "false")
	_ = mcpCmd.Flags().Set("help", "false")
	_ = mcpInstallCmd.Flags().Set("help", "false")

	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"mcp", "install", "--all", "--non-interactive"})

	err := RootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Configuring Battery MCP Server")
	assert.FileExists(t, filepath.Join(tempDir, ".cursor", "mcp.json"))
}

func TestMCPInstall_InvalidClientError(t *testing.T) {
	tempDir := t.TempDir()
	var out bytes.Buffer
	err := runMCPInstall(&out, tempDir, tempDir, []string{"nonexistent-client"}, false, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown client ID")
}
