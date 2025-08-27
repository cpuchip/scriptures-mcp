package main

import (
	"strings"
	"testing"
)

func TestValidateJSONRPC(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid initialize message",
			input:       `{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}`,
			expectError: false,
		},
		{
			name:        "Valid search_scriptures call",
			input:       `{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "faith", "limit": 10}}}`,
			expectError: false,
		},
		{
			name:        "Original problematic query with double quotes",
			input:       `{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": ""see Jesus walking", "limit": 10}}}`,
			expectError: true,
			errorMsg:    "extra quotes in query string",
		},
		{
			name:        "Missing method field",
			input:       `{"jsonrpc": "2.0", "id": 1, "params": {}}`,
			expectError: true,
			errorMsg:    "Missing 'method' field",
		},
		{
			name:        "Wrong JSON-RPC version",
			input:       `{"jsonrpc": "1.0", "id": 1, "method": "initialize", "params": {}}`,
			expectError: true,
			errorMsg:    "Invalid JSON-RPC version",
		},
		{
			name:        "Missing query argument for search_scriptures",
			input:       `{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"limit": 10}}}`,
			expectError: true,
			errorMsg:    "Missing required 'query' argument",
		},
		{
			name:        "Query with surrounding quotes",
			input:       `{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "\"faith\"", "limit": 10}}}`,
			expectError: true,
			errorMsg:    "contains extra quotes",
		},
		{
			name:        "Empty line",
			input:       ``,
			expectError: false,
		},
		{
			name:        "Completely malformed JSON",
			input:       `{this is not json}`,
			expectError: true,
			errorMsg:    "Invalid JSON syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONRPC(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, but got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %q, but got: %v", tt.input, err)
				}
			}
		})
	}
}

func TestValidateJSONRPCMultipleMessages(t *testing.T) {
	// Test the exact scenario from the problem statement
	validMessage := `{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}`
	invalidMessage := `{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": ""see Jesus walking", "limit": 10}}}`
	correctedMessage := `{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "see Jesus walking", "limit": 10}}}`

	// Valid message should pass
	if err := validateJSONRPC(validMessage); err != nil {
		t.Errorf("Valid initialize message should pass validation: %v", err)
	}

	// Invalid message should fail with specific error
	if err := validateJSONRPC(invalidMessage); err == nil {
		t.Errorf("Invalid message with double quotes should fail validation")
	} else if !strings.Contains(err.Error(), "extra quotes") {
		t.Errorf("Error should mention extra quotes issue: %v", err)
	}

	// Corrected message should pass
	if err := validateJSONRPC(correctedMessage); err != nil {
		t.Errorf("Corrected message should pass validation: %v", err)
	}
}