# LDS Scriptures MCP Server

A Model Context Protocol (MCP) server in Go that provides AI tools for searching and retrieving scripture references from the standard works of The Church of Jesus Christ of Latter-day Saints.

## Features

### Scripture Tools
The server provides three main AI tools:

1. **`search_scriptures`**: Search for scriptures by keywords or phrases across all standard works
2. **`get_scripture`**: Retrieve specific scripture verses by reference
3. **`get_chapter`**: Retrieve complete chapters from scriptures

### Standard Works Coverage
- Book of Mormon
- Bible (King James Version) 
- Doctrine and Covenants
- Pearl of Great Price

## Implementation Details

### Built With
- **[mcp-go](https://github.com/mark3labs/mcp-go)**: Go SDK for building MCP servers
- **[scriptures-json](https://github.com/bcbooks/scriptures-json)**: Complete JSON scripture database

### Scripture Reference Parsing
Supports multiple reference formats:
- Single verses: `"1 Nephi 3:7"`, `"John 3:16"`
- Verse ranges: `"John 3:16-17"`, `"Matthew 5:3-12"` 
- Full chapters: `"1 Nephi 3"`, `"Matthew 5"`

## Installation

### Prerequisites
- Go 1.19 or later

### Build from Source
```bash
git clone https://github.com/cpuchip/scriptures-mcp.git
cd scriptures-mcp
go build -o scriptures-mcp .
```

## Usage

### Running the Server
```bash
./scriptures-mcp
```

The server implements the Model Context Protocol (MCP) and communicates via JSON-RPC over stdin/stdout.

### Available Tools

#### 1. `search_scriptures`
Search for scriptures by keyword or phrase.

**Parameters:**
- `query` (string, required): The search term or phrase
- `limit` (number, optional): Maximum number of results (default: 10)

**Example:**
```json
{
  "name": "search_scriptures",
  "arguments": {
    "query": "faith",
    "limit": 5
  }
}
```

#### 2. `get_scripture`
Retrieve a specific scripture reference.

**Parameters:**
- `query` (string, required): Scripture reference (e.g., "1 Nephi 3:7", "John 3:16-17")

**Example:**
```json
{
  "name": "get_scripture", 
  "arguments": {
    "query": "1 Nephi 3:7"
  }
}
```

#### 3. `get_chapter`
Retrieve a full chapter from scriptures.

**Parameters:**
- `query` (string, required): Chapter reference (e.g., "1 Nephi 3", "Matthew 5")

**Example:**
```json
{
  "name": "get_chapter",
  "arguments": {
    "query": "1 Nephi 3"
  }
}
```

## MCP Integration

This server is designed to work with MCP-compatible AI assistants and applications.

### Example MCP Configuration
```json
{
  "mcpServers": {
    "scriptures": {
      "command": "./scriptures-mcp",
      "args": []
    }
  }
}
```

### Testing
```bash
# Test search functionality
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "search_scriptures", "arguments": {"query": "faith", "limit": 3}}}' | ./scriptures-mcp

# Test specific scripture lookup  
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "1.0.0", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "get_scripture", "arguments": {"query": "1 Nephi 3:7"}}}' | ./scriptures-mcp
```

## Data Sources

The scripture data comes from the [scriptures-json](https://github.com/bcbooks/scriptures-json) repository, which provides:
- Complete text of all standard works
- JSON format with books, chapters, and verses
- Proper verse references and citations
- Based on official 2013 LDS edition

## Benefits for AI Assistants

This MCP server enables AI assistants to:
- **Search scriptures contextually** based on user questions about gospel topics
- **Provide accurate scripture references** with proper citations  
- **Quote complete passages** when discussing doctrinal concepts
- **Cross-reference related verses** across different standard works

## Development

### Project Structure
```
scriptures-mcp/
├── main.go                     # Entry point
├── data/                       # Scripture JSON files
│   ├── book-of-mormon.json
│   ├── doctrine-and-covenants.json  
│   ├── new-testament.json
│   ├── old-testament.json
│   └── pearl-of-great-price.json
├── internal/
│   └── scripture/
│       └── service.go          # Scripture search and retrieval logic
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── test_server.sh             # Test script
├── examples.sh                # Usage examples
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Scripture content from [bcbooks/scriptures-json](https://github.com/bcbooks/scriptures-json) 
- Built using [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)
- Based on the Model Context Protocol specification
