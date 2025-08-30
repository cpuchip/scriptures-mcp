package main

import (
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/cpuchip/scriptures-mcp/internal/scripture"
)

func main() {
	// Create a new MCP server
	mcpServer := server.NewMCPServer(
		"LDS Scriptures MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	
	// Initialize scripture service
	scriptureService := scripture.NewService()
	
	// Create and register search_scriptures tool
	searchTool := mcp.NewTool("search_scriptures",
		mcp.WithDescription("Search for scriptures by keyword or phrase across all standard works"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The keyword or phrase to search for in scripture text"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return (default: 10)"),
		),
		mcp.WithString("book",
			mcp.Description("Optional: filter results to a specific book (e.g., '1 Nephi', 'John')"),
		),
		mcp.WithString("collection",
			mcp.Description("Optional: filter results to a specific collection (e.g., 'Book of Mormon', 'New Testament')"),
		),
	)
	mcpServer.AddTool(searchTool, scriptureService.SearchScriptures)
	
	// Create and register get_scripture tool
	getScriptureTool := mcp.NewTool("get_scripture",
		mcp.WithDescription("Retrieve specific scripture verses by reference"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Scripture reference like '1 Nephi 3:7' or 'John 3:16-17'"),
		),
	)
	mcpServer.AddTool(getScriptureTool, scriptureService.GetScripture)
	
	// Create and register get_chapter tool
	getChapterTool := mcp.NewTool("get_chapter",
		mcp.WithDescription("Retrieve complete chapters from scriptures"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Chapter reference like '1 Nephi 3' or 'Matthew 5'"),
		),
	)
	mcpServer.AddTool(getChapterTool, scriptureService.GetChapter)
	
	// Create and register list_collections tool
	listCollectionsTool := mcp.NewTool("list_collections",
		mcp.WithDescription("List all available scripture collections (Book of Mormon, New Testament, etc.)"),
	)
	mcpServer.AddTool(listCollectionsTool, scriptureService.ListCollections)
	
	// Create and register list_books tool
	listBooksTool := mcp.NewTool("list_books",
		mcp.WithDescription("List all available books, optionally filtered by collection"),
		mcp.WithString("collection",
			mcp.Description("Optional: filter to books within a specific collection"),
		),
	)
	mcpServer.AddTool(listBooksTool, scriptureService.ListBooks)
	
	// Create and register term_counts tool
	termCountsTool := mcp.NewTool("term_counts",
		mcp.WithDescription("Count occurrences of specific terms with optional filtering"),
		mcp.WithArray("terms",
			mcp.Required(),
			mcp.Items(mcp.WithString("term", mcp.Description("Term to count"))),
			mcp.Description("Array of terms to count occurrences of"),
		),
		mcp.WithString("book",
			mcp.Description("Optional: filter to a specific book"),
		),
		mcp.WithString("collection",
			mcp.Description("Optional: filter to a specific collection"),
		),
		mcp.WithBoolean("ignore_common_words",
			mcp.Description("Whether to ignore common words like 'the', 'and', etc. (default: true)"),
		),
	)
	mcpServer.AddTool(termCountsTool, scriptureService.GetTermCounts)
	
	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}