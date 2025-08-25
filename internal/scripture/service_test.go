package scripture

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// Test data for scripture testing
var testScriptureData = ScriptureData{
	Books: []struct {
		Book     string `json:"book"`
		Chapters []struct {
			Chapter int `json:"chapter"`
			Verses  []struct {
				Verse     int    `json:"verse"`
				Text      string `json:"text"`
				Reference string `json:"reference"`
			} `json:"verses"`
		} `json:"chapters"`
	}{
		{
			Book: "1 Nephi",
			Chapters: []struct {
				Chapter int `json:"chapter"`
				Verses  []struct {
					Verse     int    `json:"verse"`
					Text      string `json:"text"`
					Reference string `json:"reference"`
				} `json:"verses"`
			}{
				{
					Chapter: 3,
					Verses: []struct {
						Verse     int    `json:"verse"`
						Text      string `json:"text"`
						Reference string `json:"reference"`
					}{
						{Verse: 7, Text: "And it came to pass that I, Nephi, said unto my father: I will go and do the things which the Lord hath commanded, for I know that the Lord giveth no commandments unto the children of men, save he shall prepare a way for them that they may accomplish the thing which he commandeth them.", Reference: "1 Nephi 3:7"},
						{Verse: 8, Text: "And it came to pass that when my father had heard these words he was exceedingly glad, for he knew that I had been blessed of the Lord.", Reference: "1 Nephi 3:8"},
					},
				},
				{
					Chapter: 17,
					Verses: []struct {
						Verse     int    `json:"verse"`
						Text      string `json:"text"`
						Reference string `json:"reference"`
					}{
						{Verse: 50, Text: "And I said unto them: If God had commanded me to do all things I could do them. If he should command me that I should say unto this water, be thou earth, it should be earth; and if I should say it, it would be done.", Reference: "1 Nephi 17:50"},
					},
				},
			},
		},
		{
			Book: "John",
			Chapters: []struct {
				Chapter int `json:"chapter"`
				Verses  []struct {
					Verse     int    `json:"verse"`
					Text      string `json:"text"`
					Reference string `json:"reference"`
				} `json:"verses"`
			}{
				{
					Chapter: 3,
					Verses: []struct {
						Verse     int    `json:"verse"`
						Text      string `json:"text"`
						Reference string `json:"reference"`
					}{
						{Verse: 16, Text: "For God so loved the world, that he gave his only begotten Son, that whosoever believeth in him should not perish, but have everlasting life.", Reference: "John 3:16"},
						{Verse: 17, Text: "For God sent not his Son into the world to condemn the world; but that the world through him might be saved.", Reference: "John 3:17"},
					},
				},
			},
		},
	},
}

// createTestDataFile creates a temporary JSON file with test data
func createTestDataFile(t *testing.T, filename string, data ScriptureData) string {
	tmpDir := t.TempDir()
	filepath := filepath.Join(tmpDir, filename)
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	
	err = os.WriteFile(filepath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	return filepath
}

func TestService_NewService(t *testing.T) {
	// Create a service (this will try to load from data directory)
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Verify service is initialized correctly
	if service.scriptures == nil {
		t.Error("Expected scriptures map to be initialized")
	}
}

func TestService_loadScriptureFile(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Create test data file
	testFile := createTestDataFile(t, "test-scripture.json", testScriptureData)
	
	// Load the test file
	service.loadScriptureFile(testFile)
	
	// Verify data was loaded correctly
	if len(service.scriptures) != 2 {
		t.Errorf("Expected 2 books, got %d", len(service.scriptures))
	}
	
	// Check 1 Nephi data
	nephiScriptures := service.scriptures["1 Nephi"]
	if len(nephiScriptures) != 3 {
		t.Errorf("Expected 3 verses for 1 Nephi, got %d", len(nephiScriptures))
	}
	
	// Check John data
	johnScriptures := service.scriptures["John"]
	if len(johnScriptures) != 2 {
		t.Errorf("Expected 2 verses for John, got %d", len(johnScriptures))
	}
	
	// Verify specific verse
	found := false
	for _, scripture := range nephiScriptures {
		if scripture.Chapter == 3 && scripture.Verse == 7 {
			found = true
			expectedText := "And it came to pass that I, Nephi, said unto my father: I will go and do the things which the Lord hath commanded, for I know that the Lord giveth no commandments unto the children of men, save he shall prepare a way for them that they may accomplish the thing which he commandeth them."
			if scripture.Text != expectedText {
				t.Errorf("Expected correct verse text for 1 Nephi 3:7")
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find 1 Nephi 3:7")
	}
}

func TestService_parseReference(t *testing.T) {
	service := &Service{}
	
	tests := []struct {
		name        string
		reference   string
		expected    *ScriptureReference
		expectError bool
	}{
		{
			name:      "Single verse",
			reference: "1 Nephi 3:7",
			expected: &ScriptureReference{
				Book:     "1 Nephi",
				Chapter:  3,
				Verse:    7,
				EndVerse: 7,
			},
			expectError: false,
		},
		{
			name:      "Verse range",
			reference: "John 3:16-17",
			expected: &ScriptureReference{
				Book:     "John",
				Chapter:  3,
				Verse:    16,
				EndVerse: 17,
			},
			expectError: false,
		},
		{
			name:        "Invalid format",
			reference:   "Invalid reference",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Missing verse",
			reference:   "1 Nephi 3",
			expected:    nil,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseReference(tt.reference)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.Book != tt.expected.Book {
				t.Errorf("Expected book '%s', got '%s'", tt.expected.Book, result.Book)
			}
			if result.Chapter != tt.expected.Chapter {
				t.Errorf("Expected chapter %d, got %d", tt.expected.Chapter, result.Chapter)
			}
			if result.Verse != tt.expected.Verse {
				t.Errorf("Expected verse %d, got %d", tt.expected.Verse, result.Verse)
			}
			if result.EndVerse != tt.expected.EndVerse {
				t.Errorf("Expected end verse %d, got %d", tt.expected.EndVerse, result.EndVerse)
			}
		})
	}
}

func TestService_parseChapterReference(t *testing.T) {
	service := &Service{}
	
	tests := []struct {
		name        string
		reference   string
		expected    *ScriptureReference
		expectError bool
	}{
		{
			name:      "Valid chapter reference",
			reference: "1 Nephi 3",
			expected: &ScriptureReference{
				Book:    "1 Nephi",
				Chapter: 3,
			},
			expectError: false,
		},
		{
			name:      "Another valid reference",
			reference: "John 15",
			expected: &ScriptureReference{
				Book:    "John",
				Chapter: 15,
			},
			expectError: false,
		},
		{
			name:        "Invalid format",
			reference:   "Invalid reference",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "With verse number",
			reference:   "1 Nephi 3:7",
			expected:    nil,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.parseChapterReference(tt.reference)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.Book != tt.expected.Book {
				t.Errorf("Expected book '%s', got '%s'", tt.expected.Book, result.Book)
			}
			if result.Chapter != tt.expected.Chapter {
				t.Errorf("Expected chapter %d, got %d", tt.expected.Chapter, result.Chapter)
			}
		})
	}
}

func TestService_performSearch(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "I will go and do the things which the Lord hath commanded", Reference: "1 Nephi 3:7"},
		{Book: "1 Nephi", Chapter: 17, Verse: 50, Text: "If God had commanded me to do all things I could do them", Reference: "1 Nephi 17:50"},
	}
	service.scriptures["John"] = []Scripture{
		{Book: "John", Chapter: 3, Verse: 16, Text: "For God so loved the world", Reference: "John 3:16"},
	}
	
	tests := []struct {
		name           string
		query          string
		limit          int
		expectedCount  int
		shouldContain  string
	}{
		{
			name:          "Search for 'God'",
			query:         "God",
			limit:         10,
			expectedCount: 2,
			shouldContain: "God",
		},
		{
			name:          "Search for 'Lord'",
			query:         "Lord",
			limit:         10,
			expectedCount: 1,
			shouldContain: "Lord",
		},
		{
			name:          "Search with limit",
			query:         "God",
			limit:         1,
			expectedCount: 1,
			shouldContain: "God",
		},
		{
			name:          "No matches",
			query:         "nonexistent",
			limit:         10,
			expectedCount: 0,
			shouldContain: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := service.performSearch(tt.query, tt.limit)
			
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
			
			if tt.expectedCount > 0 {
				found := false
				for _, result := range results {
					if strings.Contains(result.Text, tt.shouldContain) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected results to contain '%s'", tt.shouldContain)
				}
			}
		})
	}
}

func TestService_getScripturesByReference(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "I will go and do", Reference: "1 Nephi 3:7"},
		{Book: "1 Nephi", Chapter: 3, Verse: 8, Text: "And it came to pass", Reference: "1 Nephi 3:8"},
		{Book: "1 Nephi", Chapter: 17, Verse: 50, Text: "If God had commanded", Reference: "1 Nephi 17:50"},
	}
	
	tests := []struct {
		name           string
		reference      *ScriptureReference
		expectedCount  int
	}{
		{
			name: "Single verse",
			reference: &ScriptureReference{
				Book:     "1 Nephi",
				Chapter:  3,
				Verse:    7,
				EndVerse: 7,
			},
			expectedCount: 1,
		},
		{
			name: "Verse range",
			reference: &ScriptureReference{
				Book:     "1 Nephi",
				Chapter:  3,
				Verse:    7,
				EndVerse: 8,
			},
			expectedCount: 2,
		},
		{
			name: "Non-existent book",
			reference: &ScriptureReference{
				Book:     "Non-existent",
				Chapter:  1,
				Verse:    1,
				EndVerse: 1,
			},
			expectedCount: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := service.getScripturesByReference(tt.reference)
			
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestService_getChapter(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "I will go and do", Reference: "1 Nephi 3:7"},
		{Book: "1 Nephi", Chapter: 3, Verse: 8, Text: "And it came to pass", Reference: "1 Nephi 3:8"},
		{Book: "1 Nephi", Chapter: 17, Verse: 50, Text: "If God had commanded", Reference: "1 Nephi 17:50"},
	}
	
	tests := []struct {
		name          string
		book          string
		chapter       int
		expectedCount int
	}{
		{
			name:          "Chapter with 2 verses",
			book:          "1 Nephi",
			chapter:       3,
			expectedCount: 2,
		},
		{
			name:          "Chapter with 1 verse",
			book:          "1 Nephi",
			chapter:       17,
			expectedCount: 1,
		},
		{
			name:          "Non-existent chapter",
			book:          "1 Nephi",
			chapter:       99,
			expectedCount: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := service.getChapter(tt.book, tt.chapter)
			
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestService_SearchScriptures(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "I will go and do the things which the Lord hath commanded", Reference: "1 Nephi 3:7"},
	}
	
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "Valid search",
			arguments: map[string]interface{}{
				"query": "Lord",
				"limit": 5.0,
			},
			expectError: false,
		},
		{
			name: "Empty query",
			arguments: map[string]interface{}{
				"query": "",
			},
			expectError: true,
		},
		{
			name:        "Missing query",
			arguments:   map[string]interface{}{},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a proper CallToolRequest struct instead of mocking
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.arguments,
				},
			}
			result, err := service.SearchScriptures(context.Background(), request)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.expectError {
				if !result.IsError {
					t.Error("Expected error result but got success")
				}
			} else {
				if result.IsError {
					t.Error("Expected success but got error result")
				}
			}
		})
	}
}

func TestService_GetScripture(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "I will go and do", Reference: "1 Nephi 3:7"},
	}
	
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "Valid reference",
			arguments: map[string]interface{}{
				"query": "1 Nephi 3:7",
			},
			expectError: false,
		},
		{
			name: "Invalid reference format",
			arguments: map[string]interface{}{
				"query": "invalid reference",
			},
			expectError: true,
		},
		{
			name:        "Missing query",
			arguments:   map[string]interface{}{},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.arguments,
				},
			}
			result, err := service.GetScripture(context.Background(), request)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.expectError {
				if !result.IsError {
					t.Error("Expected error result but got success")
				}
			} else {
				if result.IsError {
					t.Error("Expected success but got error result")
				}
			}
		})
	}
}

func TestService_GetChapter(t *testing.T) {
	service := &Service{
		scriptures: make(map[string][]Scripture),
	}
	
	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "I will go and do", Reference: "1 Nephi 3:7"},
		{Book: "1 Nephi", Chapter: 3, Verse: 8, Text: "And it came to pass", Reference: "1 Nephi 3:8"},
	}
	
	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "Valid chapter reference",
			arguments: map[string]interface{}{
				"query": "1 Nephi 3",
			},
			expectError: false,
		},
		{
			name: "Invalid chapter format",
			arguments: map[string]interface{}{
				"query": "invalid reference",
			},
			expectError: true,
		},
		{
			name:        "Missing query",
			arguments:   map[string]interface{}{},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.arguments,
				},
			}
			result, err := service.GetChapter(context.Background(), request)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.expectError {
				if !result.IsError {
					t.Error("Expected error result but got success")
				}
			} else {
				if result.IsError {
					t.Error("Expected success but got error result")
				}
			}
		})
	}
}

