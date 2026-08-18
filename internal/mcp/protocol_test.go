package mcp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONRPC_Serialization(t *testing.T) {
	reqJSON := `{"jsonrpc":"2.0","id":1,"method":"ping"}`
	var req Request
	err := json.Unmarshal([]byte(reqJSON), &req)
	require.NoError(t, err)
	assert.Equal(t, "2.0", req.JSONRPC)
	assert.Equal(t, "ping", req.Method)
	assert.NotNil(t, req.ID)

	resp := NewResponse(req.ID, map[string]string{"status": "pong"})
	respBytes, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.Contains(t, string(respBytes), `"jsonrpc":"2.0"`)
	assert.Contains(t, string(respBytes), `"status":"pong"`)

	errResp := NewErrorResponse(req.ID, MethodNotFoundCode, "Method not found", nil)
	errBytes, err := json.Marshal(errResp)
	require.NoError(t, err)
	assert.Contains(t, string(errBytes), `Method not found`)
	assert.Contains(t, string(errBytes), `-32601`)
}

func TestMCP_InitializeSerialization(t *testing.T) {
	initResult := InitializeResult{
		ProtocolVersion: LatestProtocolVersion,
		ServerInfo: Implementation{
			Name:    "battery-mcp",
			Version: "v1.2.0",
		},
		Capabilities: ServerCapabilities{
			Tools:     &ToolsCapability{},
			Resources: &ResourcesCapability{},
			Prompts:   &PromptsCapability{},
		},
	}

	data, err := json.Marshal(initResult)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"protocolVersion":"2024-11-05"`)
	assert.Contains(t, string(data), `"name":"battery-mcp"`)
	assert.Contains(t, string(data), `"tools":`)
	assert.Contains(t, string(data), `"resources":`)
	assert.Contains(t, string(data), `"prompts":`)
}

func TestMCP_ToolStructures(t *testing.T) {
	tool := Tool{
		Name:        "battery_status",
		Description: "Check workspace health",
		InputSchema: ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"verbose": map[string]string{"type": "boolean"},
			},
		},
	}

	data, err := json.Marshal(tool)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"battery_status"`)

	result := NewTextResult("all systems go", false)
	assert.False(t, result.IsError)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "text", result.Content[0].Type)
	assert.Equal(t, "all systems go", result.Content[0].Text)

	errResult := NewErrorResult("failed to load")
	assert.True(t, errResult.IsError)
	assert.Equal(t, "failed to load", errResult.Content[0].Text)
}
