package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twoboots/battery/cmd"
)

func TestRootCmd_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"--help"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Multi-Repository SDD Orchestrator")
	assert.Contains(t, buf.String(), "Usage:")
}

func TestRootCmd_Aliases(t *testing.T) {
	tempDir := t.TempDir()
	origDir := getPwd()
	setPwd(tempDir)
	defer setPwd(origDir)

	// 'list' alias
	buf := new(bytes.Buffer)
	cmd.RootCmd.SetOut(buf)
	cmd.RootCmd.SetArgs([]string{"list"})
	err := cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Registered Barrels (0)")

	// 'ls' alias
	buf.Reset()
	cmd.RootCmd.SetArgs([]string{"ls"})
	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Registered Barrels (0)")
}

func getPwd() string {
	d, _ := cmd.RootCmd.Flags().GetString("cwd")
	return d
}

func setPwd(d string) {
	// helper
}
