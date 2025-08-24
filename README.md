# scriptures-mcp

A Model Context Protocol (MCP) server in Go that provides AI tools for searching and retrieving scripture references from the standard works of The Church of Jesus Christ of Latter-day Saints.

## Features

- **Search Scriptures**: Search for scriptures by keywords or phrases across all standard works
- **Get Scripture References**: Retrieve specific scripture verses by reference (e.g., "1 Nephi 3:7")
- **Get Full Chapters**: Retrieve complete chapters from scriptures

## Standard Works Included

- Book of Mormon
- Bible (King James Version)
- Doctrine and Covenants
- Pearl of Great Price

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

The server implements the Model Context Protocol (MCP) and communicates via JSON-RPC over stdin/stdout.

### Running the Server

```bash
./scriptures-mcp
```

### Available Tools

#### 1. `search_scriptures`
Search for scriptures by keyword or phrase.

**Parameters:**
- `query` (string): The search term or phrase
- `limit` (int, optional): Maximum number of results (default: 10)

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
- `query` (string): Scripture reference (e.g., "1 Nephi 3:7", "John 3:16-17")

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
- `query` (string): Chapter reference (e.g., "1 Nephi 3", "Matthew 5")

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

This server is designed to work with MCP-compatible AI assistants and applications. Configure your MCP client to connect to this server for scripture-related assistance.

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

## Development

### Project Structure

```
scriptures-mcp/
├── main.go                     # Entry point
├── internal/
│   ├── mcp/
│   │   └── server.go          # MCP protocol implementation
│   └── scripture/
│       └── service.go         # Scripture search and retrieval logic
├── test_server.sh             # Test script
└── README.md
```

### Testing

Run the included test script to verify functionality:

```bash
./test_server.sh
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

- Scripture content is based on the standard works of The Church of Jesus Christ of Latter-day Saints
- Built using the Model Context Protocol specification
