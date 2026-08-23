package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoBoots/battery/internal/track"
)

func TestResources_DefaultResourcesAndDynamicResolution(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultResources(srv)

	// 1. List resources
	listResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(1),
		Method:  "resources/list",
	})
	require.Nil(t, listResp.Error)
	resList := listResp.Result.(ListResourcesResult)
	assert.NotEmpty(t, resList.Resources)

	// 2. Read battery://topology
	readTopResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(2),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://topology"}`),
	})
	require.Nil(t, readTopResp.Error)
	topContent := readTopResp.Result.(ReadResourceResult)
	assert.Len(t, topContent.Contents, 1)
	assert.Equal(t, "application/json", topContent.Contents[0].MIMEType)
	assert.Contains(t, topContent.Contents[0].Text, "multi-repo")

	// 3. Read barrel tech-stack
	readTechResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(3),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://barrels/auth/tech-stack"}`),
	})
	require.Nil(t, readTechResp.Error)
	techContent := readTechResp.Result.(ReadResourceResult)
	assert.Len(t, techContent.Contents, 1)
	assert.Equal(t, "text/markdown", techContent.Contents[0].MIMEType)
	assert.Contains(t, techContent.Contents[0].Text, "Go 1.23+")

	// 4. Read track resource after creating a track
	trackID := "track_res_test_20260817"
	_, err := track.InitTrack(dir, trackID, []string{"auth"}, track.InitTrackOptions{Name: "Resource Test"})
	require.NoError(t, err)

	readTrackResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(4),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://tracks/` + trackID + `"}`),
	})
	require.Nil(t, readTrackResp.Error)
	trackContent := readTrackResp.Result.(ReadResourceResult)
	assert.Len(t, trackContent.Contents, 1)
	assert.Equal(t, "application/json", trackContent.Contents[0].MIMEType)
	assert.Contains(t, trackContent.Contents[0].Text, trackID)
}

func TestPrompts_DefaultPrompts(t *testing.T) {
	dir := setupTestWorkspace(t)
	srv := NewServer(dir)
	RegisterDefaultPrompts(srv)

	// 1. List prompts
	listResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(10),
		Method:  "prompts/list",
	})
	pList := listResp.Result.(ListPromptsResult)
	assert.NotEmpty(t, pList.Prompts)

	// 2. Get prompt
	getResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(11),
		Method:  "prompts/get",
		Params:  json.RawMessage(`{"name":"plan_multi_barrel_track","arguments":{"track_id":"track_demo","goal":"Add OAuth authentication"}}`),
	})
	require.Nil(t, getResp.Error)
	pGet := getResp.Result.(GetPromptResult)
	assert.NotEmpty(t, pGet.Messages)
	assert.Contains(t, pGet.Messages[0].Content.Text, "track_demo")
	assert.Contains(t, pGet.Messages[0].Content.Text, "Add OAuth authentication")
}
