#!/bin/bash

# Test script for the MCP Scripture Server

echo "Testing MCP Scripture Server..."

# Build the server if not already built
if [ ! -f "./scriptures-mcp" ]; then
    echo "Building server..."
    go build -o scriptures-mcp .
fi

echo "1. Testing initialization and tools list..."
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
) | ./scriptures-mcp 2>/dev/null

echo
echo "2. Testing search_scriptures tool..."
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "faith", "limit": 1}}}'
) | ./scriptures-mcp 2>/dev/null

echo
echo "3. Testing get_scripture tool..."
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_scripture", "arguments": {"query": "1 Nephi 3:7"}}}'
) | ./scriptures-mcp 2>/dev/null

echo
echo "4. Testing get_chapter tool..."
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_chapter", "arguments": {"query": "Moroni 10"}}}'
) | ./scriptures-mcp 2>/dev/null | head -20

echo
echo "Test completed successfully."