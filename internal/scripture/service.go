package scripture

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Scripture represents a scripture verse
type Scripture struct {
	Book      string `json:"book"`
	Chapter   int    `json:"chapter"`
	Verse     int    `json:"verse"`
	Text      string `json:"text"`
	Reference string `json:"reference"`
}

// ScriptureReference represents a parsed scripture reference
type ScriptureReference struct {
	Book     string `json:"book"`
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse,omitempty"`
	EndVerse int    `json:"endVerse,omitempty"`
}

// Service handles scripture operations
type Service struct {
	scriptures map[string][]Scripture // Map of book name to scriptures
}

// NewService creates a new scripture service
func NewService() *Service {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	service.loadScriptures()
	return service
}

// loadScriptures loads scripture data from JSON files
func (s *Service) loadScriptures() {
	// Priority order:
	// 1. SCRIPTURES_DATA_DIR override (external directory)
	// 2. Embedded data (data/*.json in this package)
	// 3. Executable-relative ./data (backward compatibility)

	if override := os.Getenv("SCRIPTURES_DATA_DIR"); override != "" {
		s.loadFromDir(override)
		if len(s.scriptures) > 0 {
			return
		}
		fmt.Printf("Warning: no scripture data loaded from override dir '%s'; falling back to embedded/exe data\n", override)
	}

	// Attempt embedded data
	s.loadFromEmbedded()
	if len(s.scriptures) > 0 {
		return
	}

	// Fallback: executable-relative data directory (legacy layout)
	if exePath, err := os.Executable(); err == nil && exePath != "" {
		baseDir := filepath.Dir(exePath)
		s.loadFromDir(filepath.Join(baseDir, "data"))
	}
}

// loadFromEmbedded loads scripture JSON from the embedded filesystem.
func (s *Service) loadFromEmbedded() {
	if embeddedData == (fs.FS)(nil) { // Shouldn't happen, but guard anyway
		return
	}
	files := scriptureJSONFilenames()
	for _, f := range files {
		data, err := embeddedData.ReadFile("data/" + f)
		if err != nil {
			fmt.Printf("Warning: embedded read failed %s: %v\n", f, err)
			continue
		}
		s.parseAndStore(data, f)
	}
}

// loadFromDir loads scripture JSON files from a real directory on disk.
func (s *Service) loadFromDir(dir string) {
	files := scriptureJSONFilenames()
	for _, f := range files {
		path := filepath.Join(dir, f)
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Warning: Could not read %s: %v\n", path, err)
			continue
		}
		s.parseAndStore(data, f)
	}
}

// parseAndStore parses raw JSON scripture data and stores verses in memory.
func (s *Service) parseAndStore(data []byte, label string) {
	var scriptureData ScriptureData
	if err := json.Unmarshal(data, &scriptureData); err != nil {
		fmt.Printf("Warning: Could not parse %s: %v\n", label, err)
		return
	}
	for _, book := range scriptureData.Books {
		for _, chapter := range book.Chapters {
			for _, verse := range chapter.Verses {
				s.scriptures[book.Book] = append(s.scriptures[book.Book], Scripture{
					Book:      book.Book,
					Chapter:   chapter.Chapter,
					Verse:     verse.Verse,
					Text:      verse.Text,
					Reference: verse.Reference,
				})
			}
		}
	}
}

// scriptureJSONFilenames returns the list of scripture JSON files expected.
func scriptureJSONFilenames() []string {
	return []string{
		"book-of-mormon.json",
		"doctrine-and-covenants.json",
		"pearl-of-great-price.json",
		"old-testament.json",
		"new-testament.json",
	}
}

// ScriptureData represents the structure of the scripture JSON files
type ScriptureData struct {
	Books []struct {
		Book     string `json:"book"`
		Chapters []struct {
			Chapter int `json:"chapter"`
			Verses  []struct {
				Verse     int    `json:"verse"`
				Text      string `json:"text"`
				Reference string `json:"reference"`
			} `json:"verses"`
		} `json:"chapters"`
	} `json:"books"`
}

// loadScriptureFile loads scriptures from a single JSON file
func (s *Service) loadScriptureFile(filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s: %v\n", filepath, err)
		return
	}

	var scriptureData ScriptureData
	if err := json.Unmarshal(data, &scriptureData); err != nil {
		fmt.Printf("Warning: Could not parse %s: %v\n", filepath, err)
		return
	}

	for _, book := range scriptureData.Books {
		for _, chapter := range book.Chapters {
			for _, verse := range chapter.Verses {
				scripture := Scripture{
					Book:      book.Book,
					Chapter:   chapter.Chapter,
					Verse:     verse.Verse,
					Text:      verse.Text,
					Reference: verse.Reference,
				}
				s.scriptures[book.Book] = append(s.scriptures[book.Book], scripture)
			}
		}
	}
}

// SearchScriptures searches for scriptures by keyword or phrase
func (s *Service) SearchScriptures(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	query, ok := arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("search query cannot be empty"), nil
	}

	limit := 10 // default
	if limitVal, exists := arguments["limit"]; exists {
		if limitFloat, ok := limitVal.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Perform the search
	results := s.performSearch(query, limit)

	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No scriptures found matching '%s'. Try different keywords or check spelling.", query)), nil
	}

	response := fmt.Sprintf("Scripture Search Results for '%s':\n\n", query)
	for i, result := range results {
		response += fmt.Sprintf("%d. %s %d:%d - %s\n\n", i+1, result.Book, result.Chapter, result.Verse, result.Text)
	}

	return mcp.NewToolResultText(response), nil
}

// GetScripture retrieves a specific scripture reference
func (s *Service) GetScripture(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	query, ok := arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("scripture reference cannot be empty"), nil
	}

	// Parse the reference
	ref, err := s.parseReference(query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid scripture reference: %v", err)), nil
	}

	// Get the scripture(s)
	scriptures := s.getScripturesByReference(ref)

	if len(scriptures) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Scripture reference '%s' not found.", query)), nil
	}

	response := fmt.Sprintf("Scripture Reference: %s\n\n", query)
	for _, scripture := range scriptures {
		response += fmt.Sprintf("%s %d:%d - %s\n\n", scripture.Book, scripture.Chapter, scripture.Verse, scripture.Text)
	}

	return mcp.NewToolResultText(response), nil
}

// GetChapter retrieves a full chapter from scriptures
func (s *Service) GetChapter(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	query, ok := arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("chapter reference cannot be empty"), nil
	}

	// Parse the reference (should be book chapter format)
	ref, err := s.parseChapterReference(query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("invalid chapter reference: %v", err)), nil
	}

	// Get the entire chapter
	scriptures := s.getChapter(ref.Book, ref.Chapter)

	if len(scriptures) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("Chapter '%s' not found.", query)), nil
	}

	response := fmt.Sprintf("%s Chapter %d\n\n", ref.Book, ref.Chapter)
	for _, scripture := range scriptures {
		response += fmt.Sprintf("%d. %s\n\n", scripture.Verse, scripture.Text)
	}

	return mcp.NewToolResultText(response), nil
}

// performSearch performs a keyword search through loaded scripture data
func (s *Service) performSearch(query string, limit int) []Scripture {
	var results []Scripture
	queryLower := strings.ToLower(query)

	// Search through all loaded scriptures
	for _, bookScriptures := range s.scriptures {
		for _, scripture := range bookScriptures {
			if strings.Contains(strings.ToLower(scripture.Text), queryLower) ||
				strings.Contains(strings.ToLower(scripture.Book), queryLower) {
				results = append(results, scripture)
				if len(results) >= limit {
					return results
				}
			}
		}
	}

	return results
}

// parseReference parses a scripture reference like "1 Nephi 3:7" or "John 3:16-17"
func (s *Service) parseReference(reference string) (*ScriptureReference, error) {
	// Simple regex to parse references like "1 Nephi 3:7" or "John 3:16-17"
	re := regexp.MustCompile(`^(.+?)\s+(\d+):(\d+)(?:-(\d+))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(reference))

	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid reference format. Use format like '1 Nephi 3:7' or 'John 3:16-17'")
	}

	book := strings.TrimSpace(matches[1])
	chapter, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid chapter number: %s", matches[2])
	}
	verse, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid verse number: %s", matches[3])
	}
	endVerse := verse

	if matches[4] != "" {
		endVerse, err = strconv.Atoi(matches[4])
		if err != nil {
			return nil, fmt.Errorf("invalid end verse number: %s", matches[4])
		}
	}

	return &ScriptureReference{
		Book:     book,
		Chapter:  chapter,
		Verse:    verse,
		EndVerse: endVerse,
	}, nil
}

// parseChapterReference parses a chapter reference like "1 Nephi 3"
func (s *Service) parseChapterReference(reference string) (*ScriptureReference, error) {
	// Simple regex to parse chapter references like "1 Nephi 3"
	re := regexp.MustCompile(`^(.+?)\s+(\d+)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(reference))

	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid chapter reference format. Use format like '1 Nephi 3'")
	}

	book := strings.TrimSpace(matches[1])
	chapter, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid chapter number: %s", matches[2])
	}

	return &ScriptureReference{
		Book:    book,
		Chapter: chapter,
	}, nil
}

// getScripturesByReference retrieves scriptures by reference from loaded data
func (s *Service) getScripturesByReference(ref *ScriptureReference) []Scripture {
	var results []Scripture

	// Find scriptures matching the reference
	if bookScriptures, exists := s.scriptures[ref.Book]; exists {
		for _, scripture := range bookScriptures {
			if scripture.Chapter == ref.Chapter &&
				scripture.Verse >= ref.Verse &&
				scripture.Verse <= ref.EndVerse {
				results = append(results, scripture)
			}
		}
	}

	return results
}

// getChapter retrieves an entire chapter from loaded data
func (s *Service) getChapter(book string, chapter int) []Scripture {
	var results []Scripture

	// Find all scriptures in the specified book and chapter
	if bookScriptures, exists := s.scriptures[book]; exists {
		for _, scripture := range bookScriptures {
			if scripture.Chapter == chapter {
				results = append(results, scripture)
			}
		}
	}

	return results
}
