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
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "faith", "limit": 2}}}' | ./scriptures-mcp 2>/dev/null | jq -r '.result.content[0].text'

echo
echo "2. Getting a specific scripture reference (1 Nephi 3:7):"
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_scripture", "arguments": {"query": "1 Nephi 3:7"}}}' | ./scriptures-mcp 2>/dev/null | jq -r '.result.content[0].text'

echo
echo "3. Getting a full chapter (1 Nephi 3, first 3 verses shown):"
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "get_chapter", "arguments": {"query": "1 Nephi 3"}}}' | ./scriptures-mcp 2>/dev/null | jq -r '.result.content[0].text' | head -10

echo
echo "=== End of Examples ==="