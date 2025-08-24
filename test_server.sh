#!/bin/bash

# Test script for the MCP Scripture Server

echo "Testing MCP Scripture Server..."

# Create test input file
cat > /tmp/test_requests.jsonl << 'EOF'
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "faith"}}}
{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "get_scripture", "arguments": {"query": "1 Nephi 3:7"}}}
{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "get_chapter", "arguments": {"query": "1 Nephi 3"}}}
EOF

echo "Running test requests..."
./scriptures-mcp < /tmp/test_requests.jsonl > /tmp/test_output.json 2>&1 &
SERVER_PID=$!

# Give the server time to process
sleep 2

# Kill the server
kill $SERVER_PID 2>/dev/null

echo "Test output:"
cat /tmp/test_output.json

echo "Test completed."