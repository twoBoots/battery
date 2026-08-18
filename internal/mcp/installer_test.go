package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMergeMCPServerConfig_EmptyOrNew(t *testing.T) {
	merged, err := MergeMCPServerConfig(nil, "battery", "battery", []string{"mcp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(merged, &data); err != nil {
		t.Fatalf("failed to unmarshal merged json: %v", err)
	}

	servers, ok := data["mcpServers"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected mcpServers map, got %T", data["mcpServers"])
	}

	battery, ok := servers["battery"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected battery entry in mcpServers")
	}

	if battery["command"] != "battery" {
		t.Errorf("expected command 'battery', got %v", battery["command"])
	}
}

func TestMergeMCPServerConfig_PreservesExistingServers(t *testing.T) {
	existingJSON := []byte(`{
		"theme": "dark",
		"mcpServers": {
			"filesystem": {
				"command": "npx",
				"args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
			}
		}
	}`)

	merged, err := MergeMCPServerConfig(existingJSON, "battery", "battery", []string{"mcp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(merged, &data); err != nil {
		t.Fatalf("failed to unmarshal merged json: %v", err)
	}

	if data["theme"] != "dark" {
		t.Errorf("expected theme to be preserved as 'dark', got %v", data["theme"])
	}

	servers := data["mcpServers"].(map[string]interface{})
	if _, ok := servers["filesystem"]; !ok {
		t.Errorf("expected existing 'filesystem' server to be preserved")
	}
	if _, ok := servers["battery"]; !ok {
		t.Errorf("expected 'battery' server to be added")
	}
}

func TestMergeMCPServerConfig_InvalidJSON(t *testing.T) {
	_, err := MergeMCPServerConfig([]byte("invalid-json"), "battery", "battery", []string{"mcp"})
	if err == nil {
		t.Errorf("expected error for invalid json, got nil")
	}
}

func TestMergeMCPServerConfig_NonMapMCPServers(t *testing.T) {
	existing := []byte(`{"mcpServers": "not-a-map"}`)
	merged, err := MergeMCPServerConfig(existing, "battery", "battery", []string{"mcp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(merged, &data); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	servers, ok := data["mcpServers"].(map[string]interface{})
	if !ok || servers["battery"] == nil {
		t.Errorf("expected mcpServers to be replaced with valid map containing battery")
	}
}

func TestGetSupportedClients_Detection(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	workspaceDir := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create .cursor, .gemini, .codeium in dirs
	if err := os.MkdirAll(filepath.Join(workspaceDir, ".cursor"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(workspaceDir, ".gemini"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(homeDir, ".codeium"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(homeDir, ".claude.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	clients := GetSupportedClients(workspaceDir, homeDir)
	for _, c := range clients {
		switch c.ID {
		case "cursor", "antigravity", "windsurf", "claude-code":
			if !c.Detected {
				t.Errorf("expected client %s to be detected", c.ID)
			}
		}
	}
}

func TestGetSupportedClients_AntigravityConfig(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	workspaceDir := filepath.Join(tmpDir, "workspace")

	clients := GetSupportedClients(workspaceDir, homeDir)
	var antigravity *ClientTarget
	for _, c := range clients {
		if c.ID == "antigravity" {
			target := c
			antigravity = &target
			break
		}
	}
	if antigravity == nil {
		t.Fatalf("antigravity client not found")
	}

	expectedPath := filepath.Join(homeDir, ".gemini", "config", "mcp_config.json")
	if antigravity.ConfigPath != expectedPath {
		t.Errorf("expected configPath %s, got %s", expectedPath, antigravity.ConfigPath)
	}

	expectedDisplayName := "Google Antigravity / agy (~/.gemini/config/mcp_config.json)"
	if antigravity.DisplayName != expectedDisplayName {
		t.Errorf("expected displayName %q, got %q", expectedDisplayName, antigravity.DisplayName)
	}
}

func TestGetSupportedClients_AntigravityDetection(t *testing.T) {
	// Test detection via ~/.gemini/config
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	workspaceDir := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(filepath.Join(homeDir, ".gemini", "config"), 0755); err != nil {
		t.Fatal(err)
	}
	clients := GetSupportedClients(workspaceDir, homeDir)
	for _, c := range clients {
		if c.ID == "antigravity" && !c.Detected {
			t.Errorf("expected antigravity detected via ~/.gemini/config")
		}
	}

	// Test detection via workspace/.agents
	tmpDir2 := t.TempDir()
	homeDir2 := filepath.Join(tmpDir2, "home")
	workspaceDir2 := filepath.Join(tmpDir2, "workspace")
	if err := os.MkdirAll(filepath.Join(workspaceDir2, ".agents"), 0755); err != nil {
		t.Fatal(err)
	}
	clients2 := GetSupportedClients(workspaceDir2, homeDir2)
	for _, c := range clients2 {
		if c.ID == "antigravity" && !c.Detected {
			t.Errorf("expected antigravity detected via workspace/.agents")
		}
	}
}

func TestGetSupportedClients_DefaultHome(t *testing.T) {
	tmpDir := t.TempDir()
	clients := GetSupportedClients(tmpDir, "")
	if len(clients) == 0 {
		t.Errorf("expected non-empty clients when homeDir is empty")
	}
}

func TestInstallClients_FileCreationAndUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	workspaceDir := filepath.Join(tmpDir, "workspace")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	results, err := InstallClients(workspaceDir, homeDir, []string{"all"})
	if err != nil {
		t.Fatalf("InstallClients error: %v", err)
	}

	if len(results) != 6 {
		t.Fatalf("expected 6 results, got %d", len(results))
	}

	for _, r := range results {
		if r.Error != nil {
			t.Errorf("client %s had error: %v", r.ClientID, r.Error)
		}
		if !r.Created {
			t.Errorf("client %s was expected to be created", r.ClientID)
		}

		// Verify file exists on disk and contains battery server
		content, err := os.ReadFile(r.ConfigPath)
		if err != nil {
			t.Errorf("failed to read written file %s: %v", r.ConfigPath, err)
			continue
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal(content, &parsed); err != nil {
			t.Errorf("failed to unmarshal file %s: %v", r.ConfigPath, err)
			continue
		}
		servers, ok := parsed["mcpServers"].(map[string]interface{})
		if !ok || servers["battery"] == nil {
			t.Errorf("file %s missing mcpServers.battery", r.ConfigPath)
		}
	}

	// Update test
	updateResults, err := InstallClients(workspaceDir, homeDir, []string{"cursor", "antigravity"})
	if err != nil {
		t.Fatalf("update InstallClients error: %v", err)
	}
	for _, r := range updateResults {
		if !r.Updated {
			t.Errorf("expected client %s to be updated on second pass", r.ClientID)
		}
	}
}

func TestInstallClients_InvalidClient(t *testing.T) {
	tmpDir := t.TempDir()
	_, err := InstallClients(tmpDir, tmpDir, []string{"unknown-client"})
	if err == nil {
		t.Errorf("expected error for unknown-client, got nil")
	}
}

func TestInstallClients_EmptyList(t *testing.T) {
	tmpDir := t.TempDir()
	res, err := InstallClients(tmpDir, tmpDir, []string{"", "  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 results for empty list, got %d", len(res))
	}
}

func TestGetClaudeDesktopConfigPath(t *testing.T) {
	path := GetClaudeDesktopConfigPath("/mock/home")
	if path == "" {
		t.Errorf("expected non-empty claude desktop config path")
	}

	t.Setenv("APPDATA", "/custom/appdata")
	path2 := GetClaudeDesktopConfigPath("/mock/home")
	if path2 == "" {
		t.Errorf("expected non-empty claude desktop config path with APPDATA")
	}
}
