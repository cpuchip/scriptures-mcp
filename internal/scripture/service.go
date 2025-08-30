package scripture

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Scripture represents a scripture verse
type Scripture struct {
	Book       string `json:"book"`
	Collection string `json:"collection"` // e.g., "Book of Mormon", "New Testament"
	Chapter    int    `json:"chapter"`
	Verse      int    `json:"verse"`
	Text       string `json:"text"`
	Reference  string `json:"reference"`
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
	collections map[string][]string   // Map of collection name to list of book names
}

// NewService creates a new scripture service
func NewService() *Service {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
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
	// Prefer compressed archive
	if zipBytes, err := embeddedData.ReadFile("data/scriptures.zip"); err == nil {
		if err := s.loadFromZipBytes(zipBytes, "embedded zip"); err != nil {
			fmt.Printf("Warning: failed to load embedded zip: %v (falling back to discrete files)\n", err)
		} else {
			return
		}
	}
	// Fallback: discrete JSON files (development fallback if embed pattern changed)
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
	// If a compressed archive exists, prefer it
	zipPath := filepath.Join(dir, "scriptures.zip")
	if data, err := os.ReadFile(zipPath); err == nil {
		if err := s.loadFromZipBytes(data, zipPath); err == nil {
			return
		} else {
			fmt.Printf("Warning: could not load %s: %v (falling back to discrete files)\n", zipPath, err)
		}
	}
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

// loadFromZipBytes loads scriptures from an in-memory zip archive.
func (s *Service) loadFromZipBytes(data []byte, label string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := f.Name
		if !strings.HasSuffix(name, ".json") { // skip non-json
			continue
		}
		rc, err := f.Open()
		if err != nil {
			fmt.Printf("Warning: could not open %s in %s: %v\n", name, label, err)
			continue
		}
		fileBytes, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			fmt.Printf("Warning: could not read %s in %s: %v\n", name, label, err)
			continue
		}
		s.parseAndStore(fileBytes, name)
	}
	return nil
}

// parseAndStore parses raw JSON scripture data and stores verses in memory.
func (s *Service) parseAndStore(data []byte, label string) {
	var scriptureData ScriptureData
	if err := json.Unmarshal(data, &scriptureData); err != nil {
		fmt.Printf("Warning: Could not parse %s: %v\n", label, err)
		return
	}
	
	// Determine collection name from filename
	collection := getCollectionName(label)
	
	// Track books in this collection
	var booksInCollection []string
	
	for _, book := range scriptureData.Books {
		booksInCollection = append(booksInCollection, book.Book)
		for _, chapter := range book.Chapters {
			for _, verse := range chapter.Verses {
				s.scriptures[book.Book] = append(s.scriptures[book.Book], Scripture{
					Book:       book.Book,
					Collection: collection,
					Chapter:    chapter.Chapter,
					Verse:      verse.Verse,
					Text:       verse.Text,
					Reference:  verse.Reference,
				})
			}
		}
	}
	
	// Store collection mapping
	s.collections[collection] = booksInCollection
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

// getCollectionName converts filename to readable collection name
func getCollectionName(filename string) string {
	switch {
	case strings.Contains(filename, "book-of-mormon"):
		return "Book of Mormon"
	case strings.Contains(filename, "doctrine-and-covenants"):
		return "Doctrine and Covenants"
	case strings.Contains(filename, "pearl-of-great-price"):
		return "Pearl of Great Price"
	case strings.Contains(filename, "old-testament"):
		return "Old Testament"
	case strings.Contains(filename, "new-testament"):
		return "New Testament"
	default:
		return "Unknown"
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

	// Get optional book filter
	book := ""
	if bookVal, exists := arguments["book"]; exists {
		if bookStr, ok := bookVal.(string); ok {
			book = bookStr
		}
	}

	// Get optional collection filter
	collection := ""
	if collectionVal, exists := arguments["collection"]; exists {
		if collectionStr, ok := collectionVal.(string); ok {
			collection = collectionStr
		}
	}

	// Perform the search with filters
	results := s.performSearchWithFilters(query, limit, book, collection)

	if len(results) == 0 {
		filterInfo := ""
		if book != "" {
			filterInfo = fmt.Sprintf(" in book '%s'", book)
		} else if collection != "" {
			filterInfo = fmt.Sprintf(" in collection '%s'", collection)
		}
		return mcp.NewToolResultText(fmt.Sprintf("No scriptures found matching '%s'%s. Try different keywords or check spelling.", query, filterInfo)), nil
	}

	response := fmt.Sprintf("Scripture Search Results for '%s'", query)
	if book != "" {
		response += fmt.Sprintf(" in book '%s'", book)
	} else if collection != "" {
		response += fmt.Sprintf(" in collection '%s'", collection)
	}
	response += ":\n\n"

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
	return s.performSearchWithFilters(query, limit, "", "")
}

// performSearchWithFilters performs a keyword search with optional book and collection filters
func (s *Service) performSearchWithFilters(query string, limit int, book string, collection string) []Scripture {
	var results []Scripture
	queryLower := strings.ToLower(query)
	collectionLower := strings.ToLower(collection)

	// Define search order to ensure consistent results
	var searchOrder []string
	if book != "" {
		// Search only in specified book
		if _, exists := s.scriptures[book]; exists {
			searchOrder = []string{book}
		}
	} else if collection != "" {
		// Search only in books from specified collection
		for collectionName, books := range s.collections {
			if strings.ToLower(collectionName) == collectionLower {
				searchOrder = books
				break
			}
		}
	} else {
		// Search all books in consistent order
		for bookName := range s.scriptures {
			searchOrder = append(searchOrder, bookName)
		}
		sort.Strings(searchOrder) // Ensure consistent order
	}

	// Search through scriptures in determined order
	for _, bookName := range searchOrder {
		if bookScriptures, exists := s.scriptures[bookName]; exists {
			for _, scripture := range bookScriptures {
				// Apply filters
				if book != "" && !strings.EqualFold(scripture.Book, book) {
					continue
				}
				if collection != "" && !strings.EqualFold(scripture.Collection, collection) {
					continue
				}

				// Check if text matches query
				if strings.Contains(strings.ToLower(scripture.Text), queryLower) ||
					strings.Contains(strings.ToLower(scripture.Book), queryLower) {
					results = append(results, scripture)
					if len(results) >= limit {
						return results
					}
				}
			}
		}
	}

	// Sort results for consistency (by Collection, Book, Chapter, Verse)
	sort.Slice(results, func(i, j int) bool {
		if results[i].Collection != results[j].Collection {
			return results[i].Collection < results[j].Collection
		}
		if results[i].Book != results[j].Book {
			return results[i].Book < results[j].Book
		}
		if results[i].Chapter != results[j].Chapter {
			return results[i].Chapter < results[j].Chapter
		}
		return results[i].Verse < results[j].Verse
	})

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

//getChapter retrieves an entire chapter from loaded data
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

// ListBooks lists all available books, optionally filtered by collection
func (s *Service) ListBooks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	collection := ""
	if collectionVal, exists := arguments["collection"]; exists {
		if collectionStr, ok := collectionVal.(string); ok {
			collection = collectionStr
		}
	}

	if collection != "" {
		// List books in specific collection
		collectionLower := strings.ToLower(collection)
		for collectionName, books := range s.collections {
			if strings.ToLower(collectionName) == collectionLower {
				response := fmt.Sprintf("Books in %s:\n\n", collectionName)
				for i, book := range books {
					response += fmt.Sprintf("%d. %s\n", i+1, book)
				}
				return mcp.NewToolResultText(response), nil
			}
		}
		return mcp.NewToolResultText(fmt.Sprintf("Collection '%s' not found.", collection)), nil
	}

	// List all collections and their books
	response := "Scripture Collections and Books:\n\n"
	collectionNames := make([]string, 0, len(s.collections))
	for name := range s.collections {
		collectionNames = append(collectionNames, name)
	}
	sort.Strings(collectionNames)

	for _, collectionName := range collectionNames {
		books := s.collections[collectionName]
		response += fmt.Sprintf("## %s (%d books)\n", collectionName, len(books))
		for _, book := range books {
			response += fmt.Sprintf("- %s\n", book)
		}
		response += "\n"
	}

	return mcp.NewToolResultText(response), nil
}

// ListCollections lists all available scripture collections
func (s *Service) ListCollections(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	response := "Available Scripture Collections:\n\n"
	collectionNames := make([]string, 0, len(s.collections))
	for name := range s.collections {
		collectionNames = append(collectionNames, name)
	}
	sort.Strings(collectionNames)

	for i, name := range collectionNames {
		bookCount := len(s.collections[name])
		response += fmt.Sprintf("%d. %s (%d books)\n", i+1, name, bookCount)
	}

	return mcp.NewToolResultText(response), nil
}

// GetTermCounts counts occurrences of terms with optional filtering
func (s *Service) GetTermCounts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	terms, ok := arguments["terms"].([]interface{})
	if !ok || len(terms) == 0 {
		return mcp.NewToolResultError("terms array cannot be empty"), nil
	}

	// Convert interface{} slice to string slice
	var termStrings []string
	for _, term := range terms {
		if termStr, ok := term.(string); ok {
			termStrings = append(termStrings, termStr)
		}
	}

	if len(termStrings) == 0 {
		return mcp.NewToolResultError("no valid terms provided"), nil
	}

	// Get optional filters
	book := ""
	if bookVal, exists := arguments["book"]; exists {
		if bookStr, ok := bookVal.(string); ok {
			book = bookStr
		}
	}

	collection := ""
	if collectionVal, exists := arguments["collection"]; exists {
		if collectionStr, ok := collectionVal.(string); ok {
			collection = collectionStr
		}
	}

	ignoreCommon := true // default to ignore common words
	if ignoreVal, exists := arguments["ignore_common_words"]; exists {
		if ignoreBool, ok := ignoreVal.(bool); ok {
			ignoreCommon = ignoreBool
		}
	}

	// Count terms
	termCounts := s.countTerms(termStrings, book, collection, ignoreCommon)

	// Format response
	response := "Term Counts"
	if book != "" {
		response += fmt.Sprintf(" in book '%s'", book)
	} else if collection != "" {
		response += fmt.Sprintf(" in collection '%s'", collection)
	}
	response += ":\n\n"

	for _, term := range termStrings {
		count := termCounts[strings.ToLower(term)]
		response += fmt.Sprintf("'%s': %d occurrences\n", term, count)
	}

	return mcp.NewToolResultText(response), nil
}

// countTerms counts occurrences of terms with filtering options
func (s *Service) countTerms(terms []string, book string, collection string, ignoreCommon bool) map[string]int {
	counts := make(map[string]int)
	
	// Common words to ignore if ignoreCommon is true
	commonWords := map[string]bool{
		"a": true, "an": true, "and": true, "are": true, "as": true, "at": true, "be": true, "by": true,
		"for": true, "from": true, "has": true, "he": true, "in": true, "is": true, "it": true,
		"its": true, "of": true, "on": true, "that": true, "the": true, "to": true, "was": true,
		"will": true, "with": true, "his": true, "her": true, "him": true, "she": true, "they": true,
		"their": true, "them": true, "this": true, "these": true, "those": true, "have": true,
	}

	// Initialize counts
	for _, term := range terms {
		counts[strings.ToLower(term)] = 0
	}

	// Determine which books to search
	var searchBooks []string
	if book != "" {
		searchBooks = []string{book}
	} else if collection != "" {
		collectionLower := strings.ToLower(collection)
		for collectionName, books := range s.collections {
			if strings.ToLower(collectionName) == collectionLower {
				searchBooks = books
				break
			}
		}
	} else {
		for bookName := range s.scriptures {
			searchBooks = append(searchBooks, bookName)
		}
	}

	// Count occurrences
	for _, bookName := range searchBooks {
		if bookScriptures, exists := s.scriptures[bookName]; exists {
			for _, scripture := range bookScriptures {
				text := strings.ToLower(scripture.Text)
				words := strings.FieldsFunc(text, func(r rune) bool {
					return !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '\'')
				})

				for _, word := range words {
					word = strings.ToLower(strings.Trim(word, "'"))
					if ignoreCommon && commonWords[word] {
						continue
					}
					for _, term := range terms {
						if word == strings.ToLower(term) {
							counts[strings.ToLower(term)]++
						}
					}
				}
			}
		}
	}

	return counts
}
