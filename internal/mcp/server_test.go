package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_InitializeAndPing(t *testing.T) {
	srv := NewServer(t.TempDir())

	// 1. Initialize
	initReq := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(1),
		Method:  "initialize",
		Params:  json.RawMessage(`{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0"}}`),
	}

	resp := srv.HandleRequest(context.Background(), initReq)
	require.Nil(t, resp.Error)
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, "1", string(*resp.ID))

	resultMap, ok := resp.Result.(InitializeResult)
	require.True(t, ok)
	assert.Equal(t, LatestProtocolVersion, resultMap.ProtocolVersion)
	assert.Equal(t, "battery-mcp", resultMap.ServerInfo.Name)
	assert.Equal(t, "v1.4.0", resultMap.ServerInfo.Version)
	assert.NotNil(t, resultMap.Capabilities.Tools)

	// 2. Initialized notification
	notif := Request{
		JSONRPC: JSONRPCVersion,
		Method:  "notifications/initialized",
	}
	notifResp := srv.HandleRequest(context.Background(), notif)
	assert.Nil(t, notifResp.ID) // Notification returns empty/no-id response

	// 3. Ping
	pingReq := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(2),
		Method:  "ping",
	}
	pingResp := srv.HandleRequest(context.Background(), pingReq)
	assert.Nil(t, pingResp.Error)
	assert.Equal(t, "2", string(*pingResp.ID))
}

func TestServer_ServeStdio(t *testing.T) {
	srv := NewServer(t.TempDir())

	inputLines := []string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"ping"}`,
	}
	inBuf := bytes.NewBufferString(strings.Join(inputLines, "\n") + "\n")
	var outBuf bytes.Buffer

	err := srv.Serve(context.Background(), inBuf, &outBuf)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	require.Len(t, lines, 2) // initialize response and ping response (notification produces no output)

	var r1 Response
	err = json.Unmarshal([]byte(lines[0]), &r1)
	require.NoError(t, err)
	assert.Equal(t, "1", string(*r1.ID))

	var r2 Response
	err = json.Unmarshal([]byte(lines[1]), &r2)
	require.NoError(t, err)
	assert.Equal(t, "2", string(*r2.ID))
}

func TestServer_ToolRegistrationAndCall(t *testing.T) {
	srv := NewServer(t.TempDir())

	srv.RegisterTool(Tool{
		Name:        "echo_tool",
		Description: "Echoes input back",
	}, func(ctx context.Context, args map[string]interface{}) (CallToolResult, error) {
		msg, _ := args["msg"].(string)
		if msg == "err" {
			return CallToolResult{}, errors.New("echo error")
		}
		return NewTextResult("echo: "+msg, false), nil
	})

	// 1. List tools
	listReq := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(3),
		Method:  "tools/list",
	}
	listResp := srv.HandleRequest(context.Background(), listReq)
	require.Nil(t, listResp.Error)
	toolsResult, ok := listResp.Result.(ListToolsResult)
	require.True(t, ok)
	require.Len(t, toolsResult.Tools, 1)
	assert.Equal(t, "echo_tool", toolsResult.Tools[0].Name)

	// 2. Call tool successfully
	callReq := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(4),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"echo_tool","arguments":{"msg":"hello world"}}`),
	}
	callResp := srv.HandleRequest(context.Background(), callReq)
	require.Nil(t, callResp.Error)
	callResult, ok := callResp.Result.(CallToolResult)
	require.True(t, ok)
	assert.False(t, callResult.IsError)
	assert.Equal(t, "echo: hello world", callResult.Content[0].Text)

	// 3. Call tool with error
	callErrReq := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(5),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"echo_tool","arguments":{"msg":"err"}}`),
	}
	callErrResp := srv.HandleRequest(context.Background(), callErrReq)
	require.Nil(t, callErrResp.Error) // In MCP, tool errors return CallToolResult with isError: true
	callErrResult, ok := callErrResp.Result.(CallToolResult)
	require.True(t, ok)
	assert.True(t, callErrResult.IsError)
	assert.Equal(t, "echo error", callErrResult.Content[0].Text)

	// 4. Call unknown tool
	callUnknown := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(6),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"unknown"}`),
	}
	callUnknownResp := srv.HandleRequest(context.Background(), callUnknown)
	require.NotNil(t, callUnknownResp.Error)
	assert.Equal(t, MethodNotFoundCode, callUnknownResp.Error.Code)

	// 5. Malformed params
	callBadParams := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(7),
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":`),
	}
	badResp := srv.HandleRequest(context.Background(), callBadParams)
	require.NotNil(t, badResp.Error)
	assert.Equal(t, InvalidParamsCode, badResp.Error.Code)
}

func TestServer_ResourcesAndPromptsRegistration(t *testing.T) {
	srv := NewServer(t.TempDir())

	srv.RegisterResource(Resource{
		URI:      "battery://test",
		Name:     "test",
		MIMEType: "text/plain",
	}, func(ctx context.Context, uri string) (ReadResourceResult, error) {
		if uri == "battery://test-err" {
			return ReadResourceResult{}, errors.New("read failed")
		}
		return ReadResourceResult{
			Contents: []ResourceContent{
				{URI: uri, MIMEType: "text/plain", Text: "test content"},
			},
		}, nil
	})

	srv.RegisterPrompt(Prompt{
		Name: "test_prompt",
	}, func(ctx context.Context, args map[string]string) (GetPromptResult, error) {
		if args["fail"] == "yes" {
			return GetPromptResult{}, errors.New("prompt failed")
		}
		return GetPromptResult{
			Description: "test prompt",
			Messages: []PromptMessage{
				{Role: "user", Content: ContentItem{Type: "text", Text: "do it"}},
			},
		}, nil
	})

	// List & read resources
	resListResp := srv.HandleRequest(context.Background(), Request{JSONRPC: JSONRPCVersion, ID: rawID(10), Method: "resources/list"})
	require.Nil(t, resListResp.Error)

	resReadResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(11),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://test"}`),
	})
	require.Nil(t, resReadResp.Error)

	// Resource not found
	resMissingResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(12),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":"battery://nonexistent"}`),
	})
	require.NotNil(t, resMissingResp.Error)

	// Resource bad params
	resBadParams := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(13),
		Method:  "resources/read",
		Params:  json.RawMessage(`{"uri":`),
	})
	require.NotNil(t, resBadParams.Error)

	// List & get prompts
	pListResp := srv.HandleRequest(context.Background(), Request{JSONRPC: JSONRPCVersion, ID: rawID(20), Method: "prompts/list"})
	require.Nil(t, pListResp.Error)

	pGetResp := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(21),
		Method:  "prompts/get",
		Params:  json.RawMessage(`{"name":"test_prompt"}`),
	})
	require.Nil(t, pGetResp.Error)

	// Prompt not found
	pMissing := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(22),
		Method:  "prompts/get",
		Params:  json.RawMessage(`{"name":"unknown"}`),
	})
	require.NotNil(t, pMissing.Error)

	// Prompt bad params
	pBad := srv.HandleRequest(context.Background(), Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(23),
		Method:  "prompts/get",
		Params:  json.RawMessage(`{"name":`),
	})
	require.NotNil(t, pBad.Error)
}

func TestServer_ContextCancellation(t *testing.T) {
	srv := NewServer(t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	inBuf := bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"ping"}` + "\n")
	var outBuf bytes.Buffer

	err := srv.Serve(ctx, inBuf, &outBuf)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestServer_ErrorResilience(t *testing.T) {
	srv := NewServer(t.TempDir())

	// 1. Unknown method
	unknownReq := Request{
		JSONRPC: JSONRPCVersion,
		ID:      rawID(99),
		Method:  "non_existent_method",
	}
	resp := srv.HandleRequest(context.Background(), unknownReq)
	require.NotNil(t, resp.Error)
	assert.Equal(t, MethodNotFoundCode, resp.Error.Code)

	// 2. Malformed JSON line in stream
	inBuf := bytes.NewBufferString("not valid json\n{\"jsonrpc\":\"2.0\",\"id\":10,\"method\":\"ping\"}\n")
	var outBuf bytes.Buffer

	err := srv.Serve(context.Background(), inBuf, &outBuf)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	require.Len(t, lines, 2)

	var parseErrResp Response
	err = json.Unmarshal([]byte(lines[0]), &parseErrResp)
	require.NoError(t, err)
	require.NotNil(t, parseErrResp.Error)
	assert.Equal(t, ParseErrorCode, parseErrResp.Error.Code)

	var validResp Response
	err = json.Unmarshal([]byte(lines[1]), &validResp)
	require.NoError(t, err)
	assert.Nil(t, validResp.Error)
	assert.Equal(t, "10", string(*validResp.ID))
}

func rawID(id int) *json.RawMessage {
	data, _ := json.Marshal(id)
	msg := json.RawMessage(data)
	return &msg
}
