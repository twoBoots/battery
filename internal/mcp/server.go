package mcp

import (
	"strings"

	benderMCP "github.com/twoBoots/bender/pkg/mcp"
)

type (
	Server              = benderMCP.Server
	Tool                = benderMCP.Tool
	ToolInputSchema     = benderMCP.ToolInputSchema
	ToolHandler         = benderMCP.ToolHandler
	CallToolResult      = benderMCP.CallToolResult
	ContentItem         = benderMCP.ContentItem
	Resource            = benderMCP.Resource
	ResourceContent     = benderMCP.ResourceContent
	ReadResourceResult  = benderMCP.ReadResourceResult
	ResourceHandler     = benderMCP.ResourceHandler
	Prompt              = benderMCP.Prompt
	PromptArgument      = benderMCP.PromptArgument
	PromptMessage       = benderMCP.PromptMessage
	GetPromptResult     = benderMCP.GetPromptResult
	PromptHandler       = benderMCP.PromptHandler
	Request             = benderMCP.Request
	Response            = benderMCP.Response
	RPCError            = benderMCP.RPCError
	Implementation      = benderMCP.Implementation
	InitializeParams    = benderMCP.InitializeParams
	InitializeResult    = benderMCP.InitializeResult
	ServerCapabilities  = benderMCP.ServerCapabilities
	ClientCapabilities  = benderMCP.ClientCapabilities
	ToolsCapability     = benderMCP.ToolsCapability
	ResourcesCapability = benderMCP.ResourcesCapability
	PromptsCapability   = benderMCP.PromptsCapability
	ListToolsResult     = benderMCP.ListToolsResult
	CallToolParams      = benderMCP.CallToolParams
	ListResourcesResult = benderMCP.ListResourcesResult
	ReadResourceParams  = benderMCP.ReadResourceParams
	ListPromptsResult   = benderMCP.ListPromptsResult
	GetPromptParams     = benderMCP.GetPromptParams
	ClientTarget        = benderMCP.ClientTarget
	InstallResult       = benderMCP.InstallResult
	InstallerOptions    = benderMCP.InstallerOptions
)

const (
	JSONRPCVersion        = benderMCP.JSONRPCVersion
	LatestProtocolVersion = benderMCP.LatestProtocolVersion
	ParseErrorCode        = benderMCP.ParseErrorCode
	InvalidRequestCode    = benderMCP.InvalidRequestCode
	MethodNotFoundCode    = benderMCP.MethodNotFoundCode
	InvalidParamsCode     = benderMCP.InvalidParamsCode
	InternalErrorCode     = benderMCP.InternalErrorCode
)

var (
	NewTextResult              = benderMCP.NewTextResult
	NewErrorResult             = benderMCP.NewErrorResult
	NewResponse                = benderMCP.NewResponse
	NewErrorResponse           = benderMCP.NewErrorResponse
	InstallClients             = benderMCP.InstallClients
	GetSupportedClients        = benderMCP.GetSupportedClients
	MergeMCPServerConfig       = benderMCP.MergeMCPServerConfig
	GetClaudeDesktopConfigPath = benderMCP.GetClaudeDesktopConfigPath
)

// NewServer creates a new Battery MCP Server instance backed by Bender's MCP server engine.
func NewServer(cwd string, versions ...string) *Server {
	ver := "dev"
	if len(versions) > 0 && strings.TrimSpace(versions[0]) != "" {
		ver = strings.TrimSpace(versions[0])
	}
	return benderMCP.NewServer("battery-mcp", ver, cwd)
}
