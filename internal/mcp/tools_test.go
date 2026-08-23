package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/config"
)

func setupTestWorkspace(t *testing.T) string {
	dir := t.TempDir()
	// Create mock .batteryrc
	cfg := config.BatteryConfig{
		Version:   "1.0.0",
		Structure: config.StructureMultiRepo,
		Barrels: []config.BarrelConfig{
			{Name: "auth", Path: "./barrels/auth"},
			{Name: "web", Path: "./barrels/web"},
		},
	}
	_, err := config.SaveConfig(cfg, dir, false)
	require.NoError(t, err)

	// Create barrel folders with .cooper/definition/tech-stack.md
	for _, b := range cfg.Barrels {
		bDir := filepath.Join(dir, b.Path)
		techDir := filepath.Join(bDir, ".cooper", "definition")
		err := os.MkdirAll(techDir, 0755)
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(techDir, "tech-stack.md"), []byte("# Tech Stack\nGo 1.23+ and PostgreSQL\n"), 0644)
		require.NoError(t, err)
	}

	return dir
}

func TestTools_BatteryStatusAndListBarrels(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir, "1.4.1")
	RegisterDefaultTools(srv)

	// 1. Check tools/list has all registered tools
	listResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(1),
		Method:  "tools/list",
	})
	require.Nil(t, listResp.Error)
	toolsResult := listResp.Result.(ListToolsResult)
	assert.GreaterOrEqual(t, len(toolsResult.Tools), 5)

	// 2. Call battery_status
	statusResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(2),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_status"}`),
	})
	require.Nil(t, statusResp.Error)
	statusResult := statusResp.Result.(CallToolResult)
	assert.False(t, statusResult.IsError)
	assert.Contains(t, statusResult.Content[0].Text, `"structure": "multi-repo"`)
	assert.Contains(t, statusResult.Content[0].Text, `"cli_version": "v1.4.1"`)
	assert.Contains(t, statusResult.Content[0].Text, `"config_version": "1.0.0"`)
	assert.Contains(t, statusResult.Content[0].Text, `"barrels":`)

	// 3. Call battery_list_barrels
	listBarrelsResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(3),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_list_barrels"}`),
	})
	require.Nil(t, listBarrelsResp.Error)
	listBarrelsResult := listBarrelsResp.Result.(CallToolResult)
	assert.False(t, listBarrelsResult.IsError)
	assert.Contains(t, listBarrelsResult.Content[0].Text, "auth")
	assert.Contains(t, listBarrelsResult.Content[0].Text, "Go 1.23+")
}

func TestTools_BatteryInitDispatchAndTrackStatus(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultTools(srv)

	trackID := "track_test_feature_20260817"

	// 1. Call battery_init_track
	initResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(10),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_init_track","arguments":{"track_id":"` + trackID + `","barrels":["auth","web"],"name":"Test Feature"}}`),
	})
	require.Nil(t, initResp.Error)
	initResult := initResp.Result.(CallToolResult)
	assert.False(t, initResult.IsError)
	assert.Contains(t, initResult.Content[0].Text, "Initialized track")

	// 2. Call battery_dispatch_track
	dispatchResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(11),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_dispatch_track","arguments":{"track_id":"` + trackID + `"}}`),
	})
	require.Nil(t, dispatchResp.Error)
	dispatchResult := dispatchResp.Result.(CallToolResult)
	assert.False(t, dispatchResult.IsError)
	assert.Contains(t, dispatchResult.Content[0].Text, "Dispatched")

	// Verify plan.md was omitted in dispatched barrel
	barrelPlanPath := filepath.Join(dir, "barrels", "auth", ".cooper", "active", trackID, "plan.md")
	_, err := os.Stat(barrelPlanPath)
	assert.True(t, os.IsNotExist(err))

	// 3. Call battery_track_status
	statusResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(12),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_track_status","arguments":{"track_id":"` + trackID + `"}}`),
	})
	require.Nil(t, statusResp.Error)
	statusResult := statusResp.Result.(CallToolResult)
	assert.False(t, statusResult.IsError)
	assert.Contains(t, statusResult.Content[0].Text, trackID)
	assert.Contains(t, statusResult.Content[0].Text, "auth")
}

func TestTools_ValidationErrors(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultTools(srv)

	// Missing track_id in init
	res := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(20),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_init_track","arguments":{}}`),
	})
	require.Nil(t, res.Error)
	callRes := res.Result.(CallToolResult)
	assert.True(t, callRes.IsError)
	assert.Contains(t, callRes.Content[0].Text, "track_id is required")

	// Missing track_id in dispatch
	res = srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(21),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_dispatch_track","arguments":{}}`),
	})
	require.Nil(t, res.Error)
	callRes = res.Result.(CallToolResult)
	assert.True(t, callRes.IsError)

	// Missing track_id in status
	res = srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(22),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_track_status","arguments":{}}`),
	})
	require.Nil(t, res.Error)
	callRes = res.Result.(CallToolResult)
	assert.True(t, callRes.IsError)

	// Missing barrel in init_barrel_tech_stack
	res = srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(23),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_init_barrel_tech_stack","arguments":{}}`),
	})
	require.Nil(t, res.Error)
	callRes = res.Result.(CallToolResult)
	assert.True(t, callRes.IsError)
}

func TestTools_InitBarrelTechStack(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultTools(srv)
	RegisterDefaultResources(srv)

	// Create a new un-scaffolded barrel
	newBarrelDir := filepath.Join(dir, "barrels", "billing")
	err := os.MkdirAll(newBarrelDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(newBarrelDir, "package.json"), []byte(`{"name":"billing"}`), 0644)
	require.NoError(t, err)

	// Call battery_init_barrel_tech_stack
	res := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(30),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_init_barrel_tech_stack","arguments":{"barrel":"./barrels/billing","framework":"Express","test_runner":"jest"}}`),
	})
	require.Nil(t, res.Error)
	callRes := res.Result.(CallToolResult)
	assert.False(t, callRes.IsError)
	assert.Contains(t, callRes.Content[0].Text, "tech-stack.md")

	// Verify file was created
	techPath := filepath.Join(newBarrelDir, ".cooper", "definition", "tech-stack.md")
	assert.FileExists(t, techPath)
	data, err := os.ReadFile(techPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "Express")
	assert.Contains(t, string(data), "jest")
}
