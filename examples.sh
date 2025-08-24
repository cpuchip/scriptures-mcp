#!/bin/bash

# Example usage of the scriptures-mcp server

echo "=== MCP Scripture Server Example Usage ==="
echo

# Build the server if not already built
if [ ! -f "./scriptures-mcp" ]; then
    echo "Building server..."
    go build -o scriptures-mcp .
fi

echo "1. Searching for scriptures about 'faith':"
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "example", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "faith", "limit": 2}}}'
) | ./scriptures-mcp 2>/dev/null | grep -E '"result".*"content"' | tail -n1 | jq -r '.result.content[0].text'

echo
echo "2. Getting a specific scripture reference (1 Nephi 3:7):"
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "example", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_scripture", "arguments": {"query": "1 Nephi 3:7"}}}'
) | ./scriptures-mcp 2>/dev/null | grep -E '"result".*"content"' | tail -n1 | jq -r '.result.content[0].text'

echo
echo "3. Getting a full chapter (Moroni 10:4-5):"
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "example", "version": "1.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_scripture", "arguments": {"query": "Moroni 10:4-5"}}}'
) | ./scriptures-mcp 2>/dev/null | grep -E '"result".*"content"' | tail -n1 | jq -r '.result.content[0].text'

echo
echo "=== End of Examples ==="