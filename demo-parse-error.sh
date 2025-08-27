#!/bin/bash

# Demonstration script for JSON parse error issue and solution

echo "=== JSON-RPC Parse Error Demonstration ==="
echo

# Build the server if not already built
if [ ! -f "./scriptures-mcp" ]; then
    echo "Building scriptures-mcp server..."
    go build -o scriptures-mcp .
    echo
fi

echo "1. Demonstrating the PARSE ERROR with malformed JSON:"
echo "   Query with double quotes: \"\"see Jesus walking\""
echo

MALFORMED_JSON='{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": ""see Jesus walking", "limit": 10}}}'

echo "$MALFORMED_JSON" | ./scriptures-mcp
echo

echo "2. Using the JSON validator to diagnose the problem:"
echo

echo "$MALFORMED_JSON" | go run tools/validate-json.go
echo

echo "3. Demonstrating the SOLUTION with correctly formatted JSON:"
echo "   Query with proper quotes: \"see Jesus walking\""
echo

CORRECTED_JSON='{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "see Jesus walking", "limit": 10}}}'

echo "$CORRECTED_JSON" | ./scriptures-mcp
echo

echo "=== Summary ==="
echo "❌ Problem: Extra quotes in JSON string: \"\"text\""  
echo "✅ Solution: Proper JSON formatting: \"text\""
echo "🔧 Tool: Use 'go run tools/validate-json.go' to check your JSON"
echo
echo "The -32700 Parse Error is returned by the JSON-RPC parser when"
echo "it encounters malformed JSON. Always validate your JSON syntax"
echo "before sending messages to the MCP server."