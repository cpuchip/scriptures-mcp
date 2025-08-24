package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Handler     func(params json.RawMessage) (interface{}, error)
}

// Server represents the MCP server
type Server struct {
	tools map[string]*Tool
}

// NewServer creates a new MCP server
func NewServer() *Server {
	return &Server{
		tools: make(map[string]*Tool),
	}
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(name, description string, handler func(params json.RawMessage) (interface{}, error)) {
	s.tools[name] = &Tool{
		Name:        name,
		Description: description,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query or scripture reference",
				},
			},
		},
		Handler: handler,
	}
}

// Start starts the MCP server
func (s *Server) Start() error {
	log.Println("Starting MCP Scripture Server...")
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		var request JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			log.Printf("Error parsing request: %v", err)
			continue
		}
		
		response := s.handleRequest(&request)
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
	}
	
	return scanner.Err()
}

// handleRequest processes a JSON-RPC request
func (s *Server) handleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{
						"listChanged": false,
					},
				},
				"serverInfo": map[string]interface{}{
					"name":    "scriptures-mcp",
					"version": "1.0.0",
				},
			},
		}
	
	case "tools/list":
		var toolsList []map[string]interface{}
		for _, tool := range s.tools {
			toolsList = append(toolsList, map[string]interface{}{
				"name":        tool.Name,
				"description": tool.Description,
				"inputSchema": tool.InputSchema,
			})
		}
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": toolsList,
			},
		}
	
	case "tools/call":
		return s.handleToolCall(req)
	
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

// handleToolCall processes a tool call request
func (s *Server) handleToolCall(req *JSONRPCRequest) *JSONRPCResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}
	
	tool, exists := s.tools[params.Name]
	if !exists {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Tool not found",
			},
		}
	}
	
	result, err := tool.Handler(params.Arguments)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}
	
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}