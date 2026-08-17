package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type ToolHandler func(ctx context.Context, args map[string]interface{}) (CallToolResult, error)
type ResourceHandler func(ctx context.Context, uri string) (ReadResourceResult, error)
type PromptHandler func(ctx context.Context, args map[string]string) (GetPromptResult, error)

type toolEntry struct {
	tool    Tool
	handler ToolHandler
}

type resourceEntry struct {
	resource Resource
	handler  ResourceHandler
}

type promptEntry struct {
	prompt  Prompt
	handler PromptHandler
}

// Server implements an MCP stdio server.
type Server struct {
	cwd         string
	version     string
	mu          sync.RWMutex
	tools       map[string]toolEntry
	resources   map[string]resourceEntry
	prompts     map[string]promptEntry
	initialized bool
}

// NewServer creates a new MCP Server instance.
func NewServer(cwd string) *Server {
	return &Server{
		cwd:       cwd,
		version:   "v1.2.0",
		tools:     make(map[string]toolEntry),
		resources: make(map[string]resourceEntry),
		prompts:   make(map[string]promptEntry),
	}
}

// RegisterTool adds a tool definition and its handler.
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[tool.Name] = toolEntry{tool: tool, handler: handler}
}

// RegisterResource adds a resource definition and its handler.
func (s *Server) RegisterResource(resource Resource, handler ResourceHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resources[resource.URI] = resourceEntry{resource: resource, handler: handler}
}

// RegisterPrompt adds a prompt definition and its handler.
func (s *Server) RegisterPrompt(prompt Prompt, handler PromptHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts[prompt.Name] = promptEntry{prompt: prompt, handler: handler}
}

// HandleRequest processes an individual JSON-RPC request.
func (s *Server) HandleRequest(ctx context.Context, req Request) Response {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch req.Method {
	case "initialize":
		s.initialized = true
		return NewResponse(req.ID, InitializeResult{
			ProtocolVersion: LatestProtocolVersion,
			ServerInfo: Implementation{
				Name:    "battery-mcp",
				Version: s.version,
			},
			Capabilities: ServerCapabilities{
				Tools:     &ToolsCapability{ListChanged: false},
				Resources: &ResourcesCapability{Subscribe: false, ListChanged: false},
				Prompts:   &PromptsCapability{ListChanged: false},
			},
		})

	case "notifications/initialized":
		// Notifications have no response
		return Response{JSONRPC: JSONRPCVersion}

	case "ping":
		return NewResponse(req.ID, map[string]interface{}{})

	case "tools/list":
		toolsList := make([]Tool, 0, len(s.tools))
		for _, entry := range s.tools {
			toolsList = append(toolsList, entry.tool)
		}
		return NewResponse(req.ID, ListToolsResult{Tools: toolsList})

	case "tools/call":
		var params CallToolParams
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &params); err != nil {
				return NewErrorResponse(req.ID, InvalidParamsCode, fmt.Sprintf("invalid tool params: %v", err), nil)
			}
		}
		entry, exists := s.tools[params.Name]
		if !exists {
			return NewErrorResponse(req.ID, MethodNotFoundCode, fmt.Sprintf("unknown tool: %s", params.Name), nil)
		}
		res, err := entry.handler(ctx, params.Arguments)
		if err != nil {
			return NewResponse(req.ID, NewErrorResult(err.Error()))
		}
		return NewResponse(req.ID, res)

	case "resources/list":
		resList := make([]Resource, 0, len(s.resources))
		for _, entry := range s.resources {
			resList = append(resList, entry.resource)
		}
		return NewResponse(req.ID, ListResourcesResult{Resources: resList})

	case "resources/read":
		var params ReadResourceParams
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &params); err != nil {
				return NewErrorResponse(req.ID, InvalidParamsCode, fmt.Sprintf("invalid resource params: %v", err), nil)
			}
		}
		entry, exists := s.resources[params.URI]
		if !exists {
			// Also support pattern matching prefix if needed
			return NewErrorResponse(req.ID, MethodNotFoundCode, fmt.Sprintf("resource not found: %s", params.URI), nil)
		}
		res, err := entry.handler(ctx, params.URI)
		if err != nil {
			return NewErrorResponse(req.ID, InternalErrorCode, err.Error(), nil)
		}
		return NewResponse(req.ID, res)

	case "prompts/list":
		pList := make([]Prompt, 0, len(s.prompts))
		for _, entry := range s.prompts {
			pList = append(pList, entry.prompt)
		}
		return NewResponse(req.ID, ListPromptsResult{Prompts: pList})

	case "prompts/get":
		var params GetPromptParams
		if len(req.Params) > 0 {
			if err := json.Unmarshal(req.Params, &params); err != nil {
				return NewErrorResponse(req.ID, InvalidParamsCode, fmt.Sprintf("invalid prompt params: %v", err), nil)
			}
		}
		entry, exists := s.prompts[params.Name]
		if !exists {
			return NewErrorResponse(req.ID, MethodNotFoundCode, fmt.Sprintf("prompt not found: %s", params.Name), nil)
		}
		res, err := entry.handler(ctx, params.Arguments)
		if err != nil {
			return NewErrorResponse(req.ID, InternalErrorCode, err.Error(), nil)
		}
		return NewResponse(req.ID, res)

	default:
		return NewErrorResponse(req.ID, MethodNotFoundCode, fmt.Sprintf("method not found: %s", req.Method), nil)
	}
}

// Serve reads JSON-RPC requests from in and writes JSON-RPC responses to out.
func (s *Server) Serve(ctx context.Context, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	// Buffer can handle large JSON inputs up to 10MB
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			errResp := NewErrorResponse(nil, ParseErrorCode, fmt.Sprintf("parse error: %v", err), nil)
			data, _ := json.Marshal(errResp)
			out.Write(append(data, '\n'))
			continue
		}

		resp := s.HandleRequest(ctx, req)
		// If request was a notification (no ID and not an error), do not respond
		if req.ID == nil && resp.Error == nil {
			continue
		}

		data, err := json.Marshal(resp)
		if err != nil {
			continue
		}
		out.Write(append(data, '\n'))
	}

	return scanner.Err()
}
