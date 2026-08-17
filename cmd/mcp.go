package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"
	"github.com/twoboots/battery/internal/mcp"
)

var (
	mcpTransport string
)

var mcpCmd = &cobra.Command{
	Use:     "mcp",
	Aliases: []string{"serve"},
	Short:   "Start the Model Context Protocol (MCP) server over stdio",
	Long: `Start the Battery Model Context Protocol (MCP) server over stdio.

Exposes Battery orchestration tools, living context resources, and prompt templates to AI coding assistants (e.g. Antigravity, Claude Code, Cursor, Windsurf).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMCPServer(cmd.InOrStdin(), cmd.OutOrStdout(), getWorkingDir())
	},
}

func runMCPServer(in io.Reader, out io.Writer, cwd string) error {
	srv := mcp.NewServer(cwd)
	mcp.RegisterDefaultTools(srv)
	mcp.RegisterDefaultResources(srv)
	mcp.RegisterDefaultPrompts(srv)

	ctx := context.Background()
	return srv.Serve(ctx, in, out)
}

func init() {
	mcpCmd.Flags().StringVarP(&mcpTransport, "transport", "t", "stdio", "Transport protocol to use (stdio)")
	RootCmd.AddCommand(mcpCmd)
}
