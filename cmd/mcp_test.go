package cmd

import (
	"bytes"
	"encoding/json"
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

func TestMCPCmd_Help(t *testing.T) {
	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetArgs([]string{"mcp", "--help"})

	err := RootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Model Context Protocol")
	assert.Contains(t, out.String(), "--transport")
}
