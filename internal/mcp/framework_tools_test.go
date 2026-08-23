package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/framework"
)

func TestTools_FrameworkStatusAndGetTemplate(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultTools(srv)
	RegisterDefaultResources(srv)
	RegisterDefaultPrompts(srv)

	// 1. Call battery_framework_status on workspace root
	res := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(100),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_framework_status"}`),
	})
	require.Nil(t, res.Error)
	callRes := res.Result.(CallToolResult)
	assert.False(t, callRes.IsError)
	assert.Contains(t, callRes.Content[0].Text, "cliVersion")
	assert.Contains(t, callRes.Content[0].Text, "skills/cooper-rfc")

	// 2. Call battery_framework_status on barrel "auth"
	barrelRes := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(101),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_framework_status","arguments":{"barrel":"auth"}}`),
	})
	require.Nil(t, barrelRes.Error)
	barrelCallRes := barrelRes.Result.(CallToolResult)
	assert.False(t, barrelCallRes.IsError)
	assert.Contains(t, barrelCallRes.Content[0].Text, "barrels/auth")

	// 3. Call battery_get_template with valid template
	tmplRes := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(102),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_get_template","arguments":{"name":"skills/cooper-rfc"}}`),
	})
	require.Nil(t, tmplRes.Error)
	tmplCallRes := tmplRes.Result.(CallToolResult)
	assert.False(t, tmplCallRes.IsError)
	assert.Contains(t, tmplCallRes.Content[0].Text, "cooper-rfc")

	// 4. Call battery_get_template with invalid template
	invalidRes := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(103),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_get_template","arguments":{"name":"non-existent"}}`),
	})
	require.Nil(t, invalidRes.Error)
	invalidCallRes := invalidRes.Result.(CallToolResult)
	assert.True(t, invalidCallRes.IsError)

	// 5. Call battery_get_template without name
	missingArgRes := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(104),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_get_template","arguments":{}}`),
	})
	require.Nil(t, missingArgRes.Error)
	missingArgCallRes := missingArgRes.Result.(CallToolResult)
	assert.True(t, missingArgCallRes.IsError)
	assert.Contains(t, missingArgCallRes.Content[0].Text, "name is required")
}

func TestResources_FrameworkStatusAndTemplates(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultResources(srv)

	// 1. Read battery://framework-status
	resStatus := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(110),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://framework-status"}`),
	})
	require.Nil(t, resStatus.Error)
	statusContent := resStatus.Result.(ReadResourceResult)
	assert.Len(t, statusContent.Contents, 1)
	assert.Equal(t, "application/json", statusContent.Contents[0].MIMEType)
	assert.Contains(t, statusContent.Contents[0].Text, "cliVersion")

	// 2. Read battery://templates/skills/cooper-rfc
	resTmpl := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(111),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://templates/skills/cooper-rfc"}`),
	})
	require.Nil(t, resTmpl.Error)
	tmplContent := resTmpl.Result.(ReadResourceResult)
	assert.Len(t, tmplContent.Contents, 1)
	assert.Equal(t, "text/markdown", tmplContent.Contents[0].MIMEType)
	assert.Contains(t, tmplContent.Contents[0].Text, "cooper-rfc")

	// 3. Read invalid template URI
	resInvalid := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(112),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://templates/unknown-template"}`),
	})
	require.NotNil(t, resInvalid.Error)
}

func TestPrompts_GuideFrameworkUpgradeTrack(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultPrompts(srv)

	// 1. Check prompt in list
	listResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(120),
		Method:  "prompts/list",
	})
	require.Nil(t, listResp.Error)
	pList := listResp.Result.(ListPromptsResult)
	foundPrompt := false
	for _, p := range pList.Prompts {
		if p.Name == "guide_framework_upgrade_track" {
			foundPrompt = true
			break
		}
	}
	assert.True(t, foundPrompt, "expected prompt 'guide_framework_upgrade_track' in prompts/list")

	// 2. Get prompt with arguments
	getResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(121),
		Method:  "prompts/get",
		Params:  json.RawMessage(`{"name":"guide_framework_upgrade_track","arguments":{"track_id":"track_upgrade_demo","barrel":"auth"}}`),
	})
	require.Nil(t, getResp.Error)
	pGet := getResp.Result.(GetPromptResult)
	assert.NotEmpty(t, pGet.Messages)
	assert.Contains(t, pGet.Messages[0].Content.Text, "track_upgrade_demo")
	assert.Contains(t, pGet.Messages[0].Content.Text, "auth")
	assert.Contains(t, pGet.Messages[0].Content.Text, "battery_framework_status")
	assert.Contains(t, pGet.Messages[0].Content.Text, "battery_get_template")
}

func TestTools_EnrichedBatteryStatus(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultTools(srv)

	// Populate a file in dir so framework inspection runs
	rfcContent, _ := framework.GetTemplate("skills/cooper-rfc")
	rfcPath := filepath.Join(dir, ".agents", "skills", "cooper-rfc", "SKILL.md")
	_ = os.MkdirAll(filepath.Dir(rfcPath), 0755)
	_ = os.WriteFile(rfcPath, []byte(rfcContent), 0644)

	res := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(130),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"battery_status"}`),
	})
	require.Nil(t, res.Error)
	callRes := res.Result.(CallToolResult)
	assert.False(t, callRes.IsError)
	assert.Contains(t, callRes.Content[0].Text, `"framework_status":`)
}
