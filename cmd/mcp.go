package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/twoboots/battery/internal/mcp"
)

var (
	mcpTransport             string
	mcpInstallClients        string
	mcpInstallAll            bool
	mcpInstallNonInteractive bool
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

var mcpInstallCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"setup", "configure"},
	Short:   "Configure Battery MCP server in AI coding assistants (Cursor, Antigravity, Claude, Windsurf, VS Code)",
	Long: `Automatically detect and configure Battery's MCP server in supported AI assistant configuration files.

Supported clients:
  * cursor          - Cursor IDE (.cursor/mcp.json)
  * antigravity     - Google Antigravity / agy (~/.gemini/config/mcp_config.json)
  * claude-desktop  - Anthropic Claude Desktop
  * claude-code     - Anthropic Claude Code (~/.claude.json)
  * windsurf        - Windsurf IDE (~/.codeium/windsurf/mcp_config.json)
  * vscode          - VS Code (.vscode/mcp.json)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var clients []string
		clientFlag, _ := cmd.Flags().GetString("client")
		if clientFlag != "" {
			for _, c := range strings.Split(clientFlag, ",") {
				if trimmed := strings.TrimSpace(c); trimmed != "" {
					clients = append(clients, trimmed)
				}
			}
		}
		allFlag, _ := cmd.Flags().GetBool("all")
		nonIntFlag, _ := cmd.Flags().GetBool("non-interactive")
		return runMCPInstall(cmd.OutOrStdout(), getWorkingDir(), "", clients, allFlag, nonIntFlag)
	},
}

func runMCPServer(in io.Reader, out io.Writer, cwd string) error {
	srv := mcp.NewServer(cwd, Version)
	mcp.RegisterDefaultTools(srv)
	mcp.RegisterDefaultResources(srv)
	mcp.RegisterDefaultPrompts(srv)

	ctx := context.Background()
	return srv.Serve(ctx, in, out)
}

func runMCPInstall(out io.Writer, cwd string, homeDir string, clientIDs []string, all bool, nonInteractive bool) error {
	isNonInteractive := nonInteractive ||
		os.Getenv("CI") == "true" ||
		!isTerminal()

	supported := mcp.GetSupportedClients(cwd, homeDir)

	var selectedIDs []string

	if all {
		for _, s := range supported {
			selectedIDs = append(selectedIDs, s.ID)
		}
	} else if len(clientIDs) > 0 {
		selectedIDs = clientIDs
	} else if isNonInteractive {
		// Non-interactive without specific flags: install into detected clients
		for _, s := range supported {
			if s.Detected {
				selectedIDs = append(selectedIDs, s.ID)
			}
		}
		if len(selectedIDs) == 0 {
			fmt.Fprintln(out, "ℹ️  No AI assistant configurations automatically detected. Specify clients via --client or --all.")
			return nil
		}
	} else {
		// Interactive prompt
		options := make([]huh.Option[string], len(supported))
		var defaultSelected []string
		for i, s := range supported {
			label := s.DisplayName
			if s.Detected {
				label += " [Detected]"
				defaultSelected = append(defaultSelected, s.ID)
			}
			options[i] = huh.NewOption(label, s.ID)
		}

		var chosen []string
		if len(defaultSelected) > 0 {
			chosen = defaultSelected
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select AI assistants to configure for Battery MCP (Space to toggle, Enter to confirm):").
					Options(options...).
					Value(&chosen),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		selectedIDs = chosen
	}

	if len(selectedIDs) == 0 {
		fmt.Fprintln(out, "ℹ️  No AI clients selected. MCP configuration unchanged.")
		return nil
	}

	fmt.Fprintln(out, "🔌 Configuring Battery MCP Server...")

	results, err := mcp.InstallClients(cwd, homeDir, selectedIDs)
	if err != nil {
		return err
	}

	for _, res := range results {
		if res.Error != nil {
			fmt.Fprintf(out, "  [✗] Failed %s (%s): %v\n", res.DisplayName, res.ConfigPath, res.Error)
		} else if res.Created {
			fmt.Fprintf(out, "  [✓] Configured %s -> Created %s\n", res.DisplayName, res.ConfigPath)
		} else if res.Updated {
			fmt.Fprintf(out, "  [✓] Configured %s -> Updated %s\n", res.DisplayName, res.ConfigPath)
		}
	}

	fmt.Fprintln(out, "\n✨ MCP configuration completed!")
	return nil
}

func init() {
	mcpCmd.Flags().StringVarP(&mcpTransport, "transport", "t", "stdio", "Transport protocol to use (stdio)")

	mcpInstallCmd.Flags().StringVarP(&mcpInstallClients, "client", "c", "", "Comma-separated list of clients to configure (cursor, antigravity, claude-desktop, claude-code, windsurf, vscode)")
	mcpInstallCmd.Flags().BoolVarP(&mcpInstallAll, "all", "a", false, "Configure all supported AI clients")
	mcpInstallCmd.Flags().BoolVarP(&mcpInstallNonInteractive, "non-interactive", "y", false, "Run non-interactively configuring detected clients")

	mcpCmd.AddCommand(mcpInstallCmd)
	RootCmd.AddCommand(mcpCmd)
}
