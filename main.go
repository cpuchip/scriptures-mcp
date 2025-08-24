package main

import (
	"log"

	"github.com/cpuchip/scriptures-mcp/internal/mcp"
	"github.com/cpuchip/scriptures-mcp/internal/scripture"
)

func main() {
	server := mcp.NewServer()
	
	// Register scripture tools
	scriptureService := scripture.NewService()
	server.RegisterTool("search_scriptures", "Search for scriptures by keyword or phrase", scriptureService.SearchScriptures)
	server.RegisterTool("get_scripture", "Get a specific scripture reference", scriptureService.GetScripture)
	server.RegisterTool("get_chapter", "Get a full chapter from scriptures", scriptureService.GetChapter)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}