package scripture

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestService_SearchStability tests that search results are consistent
func TestService_SearchStability(t *testing.T) {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
	}

	// Add test data with Collection field
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Collection: "Book of Mormon", Chapter: 3, Verse: 7, Text: "I will go and do the things which the Lord hath commanded", Reference: "1 Nephi 3:7"},
		{Book: "1 Nephi", Collection: "Book of Mormon", Chapter: 17, Verse: 50, Text: "If God had commanded me to do all things I could do them", Reference: "1 Nephi 17:50"},
	}
	service.scriptures["John"] = []Scripture{
		{Book: "John", Collection: "New Testament", Chapter: 3, Verse: 16, Text: "For God so loved the world", Reference: "John 3:16"},
	}
	service.collections["Book of Mormon"] = []string{"1 Nephi"}
	service.collections["New Testament"] = []string{"John"}

	// Test that multiple calls to the same search return the same results in the same order
	query := "God"
	limit := 10

	// Perform search multiple times
	results1 := service.performSearch(query, limit)
	results2 := service.performSearch(query, limit)
	results3 := service.performSearch(query, limit)

	// Check that all results are identical
	if len(results1) != len(results2) || len(results2) != len(results3) {
		t.Errorf("Search results have different lengths: %d, %d, %d", len(results1), len(results2), len(results3))
	}

	for i := range results1 {
		if results1[i] != results2[i] {
			t.Errorf("Search result %d differs between calls 1 and 2", i)
		}
		if results2[i] != results3[i] {
			t.Errorf("Search result %d differs between calls 2 and 3", i)
		}
	}

	// Verify order is consistent (should be sorted by Collection, Book, Chapter, Verse)
	for i := 1; i < len(results1); i++ {
		prev := results1[i-1]
		curr := results1[i]
		
		if prev.Collection > curr.Collection {
			t.Errorf("Results not sorted by collection: %s > %s", prev.Collection, curr.Collection)
		} else if prev.Collection == curr.Collection {
			if prev.Book > curr.Book {
				t.Errorf("Results not sorted by book within collection: %s > %s", prev.Book, curr.Book)
			} else if prev.Book == curr.Book {
				if prev.Chapter > curr.Chapter {
					t.Errorf("Results not sorted by chapter within book: %d > %d", prev.Chapter, curr.Chapter)
				} else if prev.Chapter == curr.Chapter && prev.Verse > curr.Verse {
					t.Errorf("Results not sorted by verse within chapter: %d > %d", prev.Verse, curr.Verse)
				}
			}
		}
	}
}

// TestService_SearchWithFilters tests search filtering functionality
func TestService_SearchWithFilters(t *testing.T) {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
	}

	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Collection: "Book of Mormon", Chapter: 3, Verse: 7, Text: "I will go and do the things which the Lord hath commanded", Reference: "1 Nephi 3:7"},
	}
	service.scriptures["2 Nephi"] = []Scripture{
		{Book: "2 Nephi", Collection: "Book of Mormon", Chapter: 9, Verse: 28, Text: "O that cunning plan of the evil one! O the vainness, and the frailties, and the foolishness of men!", Reference: "2 Nephi 9:28"},
	}
	service.scriptures["John"] = []Scripture{
		{Book: "John", Collection: "New Testament", Chapter: 3, Verse: 16, Text: "For God so loved the world", Reference: "John 3:16"},
	}
	service.collections["Book of Mormon"] = []string{"1 Nephi", "2 Nephi"}
	service.collections["New Testament"] = []string{"John"}

	tests := []struct {
		name             string
		query            string
		book             string
		collection       string
		expectedCount    int
		expectedContains string
	}{
		{
			name:             "Search in specific book",
			query:            "the",
			book:             "1 Nephi",
			expectedCount:    1,
			expectedContains: "1 Nephi",
		},
		{
			name:             "Search in specific collection",
			query:            "the",
			collection:       "Book of Mormon",
			expectedCount:    2,
			expectedContains: "Nephi",
		},
		{
			name:             "Search in New Testament collection",
			query:            "God",
			collection:       "New Testament",
			expectedCount:    1,
			expectedContains: "John",
		},
		{
			name:          "Search with no matches in book",
			query:         "nonexistent",
			book:          "1 Nephi",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := service.performSearchWithFilters(tt.query, 10, tt.book, tt.collection)
			
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
			
			if tt.expectedCount > 0 && tt.expectedContains != "" {
				found := false
				for _, result := range results {
					if strings.Contains(result.Book, tt.expectedContains) || strings.Contains(result.Text, tt.expectedContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected results to contain '%s'", tt.expectedContains)
				}
			}
		})
	}
}

// TestService_ListCollections tests the list collections functionality
func TestService_ListCollections(t *testing.T) {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
	}

	service.collections["Book of Mormon"] = []string{"1 Nephi", "2 Nephi"}
	service.collections["New Testament"] = []string{"Matthew", "John"}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}
	result, err := service.ListCollections(context.Background(), request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.IsError {
		t.Error("Expected success but got error result")
	}

	// The test just verifies that we don't get errors - actual content verification would require more complex content parsing
}

// TestService_ListBooks tests the list books functionality
func TestService_ListBooks(t *testing.T) {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
	}

	service.collections["Book of Mormon"] = []string{"1 Nephi", "2 Nephi"}
	service.collections["New Testament"] = []string{"Matthew", "John"}

	tests := []struct {
		name       string
		collection string
	}{
		{
			name:       "List all books",
			collection: "",
		},
		{
			name:       "List books in Book of Mormon",
			collection: "Book of Mormon",
		},
		{
			name:       "List books in New Testament",
			collection: "New Testament",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{}
			if tt.collection != "" {
				args["collection"] = tt.collection
			}

			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: args,
				},
			}
			result, err := service.ListBooks(context.Background(), request)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result.IsError {
				t.Error("Expected success but got error result")
			}
		})
	}
}

// TestService_GetTermCounts tests the term counting functionality
func TestService_GetTermCounts(t *testing.T) {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
	}

	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Collection: "Book of Mormon", Chapter: 3, Verse: 7, Text: "I will go and do the things which the Lord hath commanded", Reference: "1 Nephi 3:7"},
		{Book: "1 Nephi", Collection: "Book of Mormon", Chapter: 3, Verse: 8, Text: "And it came to pass that when my father had heard these words he was exceedingly glad, for he knew that I had been blessed of the Lord.", Reference: "1 Nephi 3:8"},
	}
	service.collections["Book of Mormon"] = []string{"1 Nephi"}

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"terms": []interface{}{"Lord", "the", "faith"},
				"book":  "1 Nephi",
			},
		},
	}
	result, err := service.GetTermCounts(context.Background(), request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.IsError {
		t.Error("Expected success but got error result")
	}
}

// TestService_SearchScripturesWithFilters tests the enhanced search scriptures with filters
func TestService_SearchScripturesWithFilters(t *testing.T) {
	service := &Service{
		scriptures:  make(map[string][]Scripture),
		collections: make(map[string][]string),
	}

	// Add test data
	service.scriptures["1 Nephi"] = []Scripture{
		{Book: "1 Nephi", Collection: "Book of Mormon", Chapter: 3, Verse: 7, Text: "I will go and do the things which the Lord hath commanded", Reference: "1 Nephi 3:7"},
	}
	service.scriptures["John"] = []Scripture{
		{Book: "John", Collection: "New Testament", Chapter: 3, Verse: 16, Text: "For God so loved the world", Reference: "John 3:16"},
	}
	service.collections["Book of Mormon"] = []string{"1 Nephi"}
	service.collections["New Testament"] = []string{"John"}

	tests := []struct {
		name        string
		arguments   map[string]interface{}
		expectError bool
	}{
		{
			name: "Search with book filter",
			arguments: map[string]interface{}{
				"query": "God",
				"book":  "John",
			},
			expectError: false,
		},
		{
			name: "Search with collection filter",
			arguments: map[string]interface{}{
				"query":      "Lord",
				"collection": "Book of Mormon",
			},
			expectError: false,
		},
		{
			name: "Search with invalid book",
			arguments: map[string]interface{}{
				"query": "God",
				"book":  "NonExistent",
			},
			expectError: false, // Should just return no results
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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