package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ClientTarget represents an AI assistant / editor MCP configuration destination.
type ClientTarget struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	ConfigPath  string `json:"configPath"`
	Detected    bool   `json:"detected"`
}

// InstallResult represents the outcome of configuring an MCP client.
type InstallResult struct {
	ClientID    string `json:"clientId"`
	DisplayName string `json:"displayName"`
	ConfigPath  string `json:"configPath"`
	Created     bool   `json:"created"`
	Updated     bool   `json:"updated"`
	Error       error  `json:"-"`
}

// GetClaudeDesktopConfigPath returns the OS-specific path for Claude Desktop config.
func GetClaudeDesktopConfigPath(homeDir string) string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Claude", "claude_desktop_config.json")
	default: // linux, bsd, etc.
		return filepath.Join(homeDir, ".config", "Claude", "claude_desktop_config.json")
	}
}

// GetSupportedClients returns a list of supported MCP clients with their target configuration paths and detection status.
func GetSupportedClients(cwd string, homeDir string) []ClientTarget {
	if homeDir == "" {
		if h, err := os.UserHomeDir(); err == nil {
			homeDir = h
		}
	}

	targets := []struct {
		id          string
		displayName string
		configPath  string
		detectPaths []string
	}{
		{
			id:          "cursor",
			displayName: "Cursor IDE (.cursor/mcp.json)",
			configPath:  filepath.Join(cwd, ".cursor", "mcp.json"),
			detectPaths: []string{
				filepath.Join(cwd, ".cursor"),
				filepath.Join(homeDir, ".cursor"),
			},
		},
		{
			id:          "antigravity",
			displayName: "Google Antigravity / Gemini (.gemini/settings.json)",
			configPath:  filepath.Join(cwd, ".gemini", "settings.json"),
			detectPaths: []string{
				filepath.Join(cwd, ".gemini"),
				filepath.Join(homeDir, ".gemini"),
			},
		},
		{
			id:          "claude-desktop",
			displayName: "Anthropic Claude Desktop",
			configPath:  GetClaudeDesktopConfigPath(homeDir),
			detectPaths: []string{
				filepath.Dir(GetClaudeDesktopConfigPath(homeDir)),
			},
		},
		{
			id:          "claude-code",
			displayName: "Anthropic Claude Code (~/.claude.json)",
			configPath:  filepath.Join(homeDir, ".claude.json"),
			detectPaths: []string{
				filepath.Join(homeDir, ".claude.json"),
				filepath.Join(homeDir, ".claude"),
			},
		},
		{
			id:          "windsurf",
			displayName: "Windsurf IDE (~/.codeium/windsurf/mcp_config.json)",
			configPath:  filepath.Join(homeDir, ".codeium", "windsurf", "mcp_config.json"),
			detectPaths: []string{
				filepath.Join(homeDir, ".codeium"),
			},
		},
		{
			id:          "vscode",
			displayName: "VS Code (.vscode/mcp.json)",
			configPath:  filepath.Join(cwd, ".vscode", "mcp.json"),
			detectPaths: []string{
				filepath.Join(cwd, ".vscode"),
			},
		},
	}

	results := make([]ClientTarget, len(targets))
	for i, t := range targets {
		detected := false
		for _, p := range t.detectPaths {
			if _, err := os.Stat(p); err == nil {
				detected = true
				break
			}
		}
		results[i] = ClientTarget{
			ID:          t.id,
			DisplayName: t.displayName,
			ConfigPath:  t.configPath,
			Detected:    detected,
		}
	}

	return results
}

// MergeMCPServerConfig takes an existing JSON byte slice and adds or updates the specified server in mcpServers.
func MergeMCPServerConfig(existingJSON []byte, serverName string, command string, args []string) ([]byte, error) {
	var root map[string]interface{}

	if len(existingJSON) > 0 && strings.TrimSpace(string(existingJSON)) != "" {
		if err := json.Unmarshal(existingJSON, &root); err != nil {
			return nil, fmt.Errorf("failed to parse existing JSON config: %w", err)
		}
	} else {
		root = make(map[string]interface{})
	}

	var mcpServers map[string]interface{}
	if rawServers, ok := root["mcpServers"]; ok {
		if m, ok := rawServers.(map[string]interface{}); ok {
			mcpServers = m
		} else {
			mcpServers = make(map[string]interface{})
		}
	} else {
		mcpServers = make(map[string]interface{})
	}

	batteryEntry := map[string]interface{}{
		"command": command,
		"args":    args,
	}

	mcpServers[serverName] = batteryEntry
	root["mcpServers"] = mcpServers

	return json.MarshalIndent(root, "", "  ")
}

// InstallClients configures the requested client IDs with the battery MCP server.
func InstallClients(cwd string, homeDir string, clientIDs []string) ([]InstallResult, error) {
	supported := GetSupportedClients(cwd, homeDir)
	supportedMap := make(map[string]ClientTarget)
	for _, c := range supported {
		supportedMap[c.ID] = c
	}

	var selectedTargets []ClientTarget
	for _, id := range clientIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		if trimmed == "all" {
			selectedTargets = supported
			break
		}
		target, ok := supportedMap[trimmed]
		if !ok {
			return nil, fmt.Errorf("unknown client ID: '%s'", trimmed)
		}
		selectedTargets = append(selectedTargets, target)
	}

	var results []InstallResult
	for _, target := range selectedTargets {
		res := InstallResult{
			ClientID:    target.ID,
			DisplayName: target.DisplayName,
			ConfigPath:  target.ConfigPath,
		}

		dir := filepath.Dir(target.ConfigPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			res.Error = fmt.Errorf("failed to create directory %s: %w", dir, err)
			results = append(results, res)
			continue
		}

		var existingContent []byte
		var fileExisted bool
		if data, err := os.ReadFile(target.ConfigPath); err == nil {
			existingContent = data
			fileExisted = true
		}

		merged, err := MergeMCPServerConfig(existingContent, "battery", "battery", []string{"mcp"})
		if err != nil {
			res.Error = err
			results = append(results, res)
			continue
		}

		// Ensure trailing newline
		mergedWithNewline := append(merged, '\n')
		if err := os.WriteFile(target.ConfigPath, mergedWithNewline, 0644); err != nil {
			res.Error = fmt.Errorf("failed to write config to %s: %w", target.ConfigPath, err)
			results = append(results, res)
			continue
		}

		if fileExisted {
			res.Updated = true
		} else {
			res.Created = true
		}

		results = append(results, res)
	}

	return results, nil
}
