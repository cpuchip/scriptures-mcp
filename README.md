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

### Scripture Data Attribution

This project uses scripture data from the excellent [bcbooks/scriptures-json](https://github.com/bcbooks/scriptures-json) repository. We are grateful to the maintainers of this project for providing:

- **Complete text** of all standard works of The Church of Jesus Christ of Latter-day Saints
- **Structured JSON format** with books, chapters, and verses for easy parsing
- **Accurate references** and proper verse citations
- **2013 LDS edition** text ensuring consistency with official sources

#### Data Files Included:
- `book-of-mormon.json` - Complete Book of Mormon (2.7MB)
- `old-testament.json` - Old Testament from KJV Bible (7.8MB)  
- `new-testament.json` - New Testament from KJV Bible (2.5MB)
- `doctrine-and-covenants.json` - Complete Doctrine and Covenants (891KB)
- `pearl-of-great-price.json` - Complete Pearl of Great Price (269KB)

**Total scripture database: ~13.6MB** providing access to 66+ books and 31,000+ verses.

#### Keeping Data Up to Date

The scripture data in this repository is currently a snapshot from [scriptures-json](https://github.com/bcbooks/scriptures-json). 

**Automated Data Sync:**
We provide a script to automatically sync with the latest scripture data:

```bash
# Run the data sync script
./sync-data.sh
```

This script will:
1. Clone the latest version of the scriptures-json repository
2. Copy the required JSON files to the `/data` directory  
3. Show you what files were updated
4. Clean up temporary files

**Manual Data Update:**
For manual updates:
1. Visit the [scriptures-json repository](https://github.com/bcbooks/scriptures-json)
2. Download the latest JSON files from their repository
3. Replace the files in the `/data` directory of this project

**Future Enhancement:** We plan to implement automated data synchronization via GitHub Actions to keep the scripture data current with updates from the source repository.

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
├── main.go                         # Entry point
├── main_test.go                   # Main package tests
├── sync-data.sh                   # Scripture data sync script
├── data/                          # Scripture JSON files
│   ├── book-of-mormon.json
│   ├── doctrine-and-covenants.json  
│   ├── new-testament.json
│   ├── old-testament.json
│   └── pearl-of-great-price.json
├── internal/
│   └── scripture/
│       ├── service.go             # Scripture search and retrieval logic
│       └── service_test.go        # Comprehensive unit tests
├── .github/
│   └── workflows/
│       └── ci.yml                 # GitHub Actions CI/CD pipeline
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── test_server.sh                 # Integration test script
├── examples.sh                    # Usage examples
└── README.md
```

### Running Tests

This project includes comprehensive unit tests using Go's built-in testing framework:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

#### Test Coverage
The tests cover:
- Scripture data loading from JSON files
- Search functionality across all scriptures
- Scripture reference parsing (both verse and chapter references)  
- Scripture retrieval by specific references
- Chapter retrieval functionality
- Error handling for invalid inputs
- MCP tool integration

### CI/CD Pipeline

The project includes a comprehensive GitHub Actions workflow that:
- **Runs tests** on every push and pull request
- **Builds binaries** for multiple platforms:
  - Linux (amd64, arm64, 386, arm)
  - macOS (amd64, arm64) 
  - Windows (amd64, arm64, 386)
- **Creates releases** automatically on main branch pushes
- **Uploads artifacts** for easy download

### Building from Source

```bash
# Clone the repository
git clone https://github.com/cpuchip/scriptures-mcp.git
cd scriptures-mcp

# Download dependencies
go mod download

# Run tests
go test ./...

# Build for current platform
go build -o scriptures-mcp .

# Build for specific platform (cross-compilation)
GOOS=windows GOARCH=amd64 go build -o scriptures-mcp.exe .
GOOS=linux GOARCH=arm64 go build -o scriptures-mcp-linux-arm64 .
GOOS=darwin GOARCH=arm64 go build -o scriptures-mcp-darwin-arm64 .
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
