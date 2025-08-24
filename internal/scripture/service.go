package scripture

import (
	"encoding/json"
	"fmt"
	"strings"
	"regexp"
)

// Scripture represents a scripture verse
type Scripture struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

// ScriptureReference represents a parsed scripture reference
type ScriptureReference struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse,omitempty"`
	EndVerse int   `json:"endVerse,omitempty"`
}

// Service handles scripture operations
type Service struct {
	// In a real implementation, this would connect to a scripture database or API
	// For now, we'll use some sample data
}

// NewService creates a new scripture service
func NewService() *Service {
	return &Service{}
}

// SearchScriptures searches for scriptures by keyword or phrase
func (s *Service) SearchScriptures(params json.RawMessage) (interface{}, error) {
	var args struct {
		Query string `json:"query"`
		Limit int    `json:"limit,omitempty"`
	}
	
	if err := json.Unmarshal(params, &args); err != nil {
		return nil, fmt.Errorf("invalid search parameters: %v", err)
	}
	
	if args.Query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}
	
	if args.Limit == 0 {
		args.Limit = 10
	}
	
	// For demonstration, return sample search results
	// In a real implementation, this would query a scripture database
	results := s.performSearch(args.Query, args.Limit)
	
	response := fmt.Sprintf("Scripture Search Results for '%s':\n\n", args.Query)
	for i, result := range results {
		response += fmt.Sprintf("%d. %s %d:%d - %s\n\n", i+1, result.Book, result.Chapter, result.Verse, result.Text)
	}
	
	if len(results) == 0 {
		response = fmt.Sprintf("No scriptures found matching '%s'. Try different keywords or check spelling.", args.Query)
	}
	
	return response, nil
}

// GetScripture retrieves a specific scripture reference
func (s *Service) GetScripture(params json.RawMessage) (interface{}, error) {
	var args struct {
		Reference string `json:"query"`
	}
	
	if err := json.Unmarshal(params, &args); err != nil {
		return nil, fmt.Errorf("invalid reference parameters: %v", err)
	}
	
	if args.Reference == "" {
		return nil, fmt.Errorf("scripture reference cannot be empty")
	}
	
	// Parse the reference
	ref, err := s.parseReference(args.Reference)
	if err != nil {
		return nil, fmt.Errorf("invalid scripture reference: %v", err)
	}
	
	// Get the scripture(s)
	scriptures := s.getScripturesByReference(ref)
	
	if len(scriptures) == 0 {
		return fmt.Sprintf("Scripture reference '%s' not found.", args.Reference), nil
	}
	
	response := fmt.Sprintf("Scripture Reference: %s\n\n", args.Reference)
	for _, scripture := range scriptures {
		response += fmt.Sprintf("%s %d:%d - %s\n\n", scripture.Book, scripture.Chapter, scripture.Verse, scripture.Text)
	}
	
	return response, nil
}

// GetChapter retrieves a full chapter from scriptures
func (s *Service) GetChapter(params json.RawMessage) (interface{}, error) {
	var args struct {
		Reference string `json:"query"`
	}
	
	if err := json.Unmarshal(params, &args); err != nil {
		return nil, fmt.Errorf("invalid chapter parameters: %v", err)
	}
	
	if args.Reference == "" {
		return nil, fmt.Errorf("chapter reference cannot be empty")
	}
	
	// Parse the reference (should be book chapter format)
	ref, err := s.parseChapterReference(args.Reference)
	if err != nil {
		return nil, fmt.Errorf("invalid chapter reference: %v", err)
	}
	
	// Get the entire chapter
	scriptures := s.getChapter(ref.Book, ref.Chapter)
	
	if len(scriptures) == 0 {
		return fmt.Sprintf("Chapter '%s' not found.", args.Reference), nil
	}
	
	response := fmt.Sprintf("%s Chapter %d\n\n", ref.Book, ref.Chapter)
	for _, scripture := range scriptures {
		response += fmt.Sprintf("%d. %s\n\n", scripture.Verse, scripture.Text)
	}
	
	return response, nil
}

// performSearch performs a keyword search (sample implementation)
func (s *Service) performSearch(query string, limit int) []Scripture {
	// Sample scripture data - in a real implementation, this would query a database
	sampleScriptures := []Scripture{
		{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "And it came to pass that I, Nephi, said unto my father: I will go and do the things which the Lord hath commanded, for I know that the Lord giveth no commandments unto the children of men, save he shall prepare a way for them that they may accomplish the thing which he commandeth them."},
		{Book: "2 Nephi", Chapter: 2, Verse: 25, Text: "Adam fell that men might be; and men are, that they might have joy."},
		{Book: "Alma", Chapter: 32, Verse: 21, Text: "And now as I said concerning faith—faith is not to have a perfect knowledge of things; therefore if ye have faith ye hope for things which are not seen, which are true."},
		{Book: "Moroni", Chapter: 10, Verse: 4, Text: "And when ye shall receive these things, I would exhort you that ye would ask God, the Eternal Father, in the name of Christ, if these things are not true; and if ye shall ask with a sincere heart, with real intent, having faith in Christ, he will manifest the truth of it unto you, by the power of the Holy Ghost."},
		{Book: "John", Chapter: 3, Verse: 16, Text: "For God so loved the world, that he gave his only begotten Son, that whosoever believeth in him should not perish, but have everlasting life."},
		{Book: "Matthew", Chapter: 5, Verse: 16, Text: "Let your light so shine before men, that they may see your good works, and glorify your Father which is in heaven."},
		{Book: "D&C", Chapter: 76, Verse: 22, Text: "And now, after the many testimonies which have been given of him, this is the testimony, last of all, which we give of him: That he lives!"},
		{Book: "Moses", Chapter: 1, Verse: 39, Text: "For behold, this is my work and my glory—to bring to pass the immortality and eternal life of man."},
	}
	
	var results []Scripture
	queryLower := strings.ToLower(query)
	
	for _, scripture := range sampleScriptures {
		if strings.Contains(strings.ToLower(scripture.Text), queryLower) ||
		   strings.Contains(strings.ToLower(scripture.Book), queryLower) {
			results = append(results, scripture)
			if len(results) >= limit {
				break
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
	chapter := parseInt(matches[2])
	verse := parseInt(matches[3])
	endVerse := verse
	
	if matches[4] != "" {
		endVerse = parseInt(matches[4])
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
	chapter := parseInt(matches[2])
	
	return &ScriptureReference{
		Book:    book,
		Chapter: chapter,
	}, nil
}

// getScripturesByReference retrieves scriptures by reference
func (s *Service) getScripturesByReference(ref *ScriptureReference) []Scripture {
	// Sample implementation - would query database in real version
	sampleData := map[string]map[int]map[int]string{
		"1 Nephi": {
			3: {
				7: "And it came to pass that I, Nephi, said unto my father: I will go and do the things which the Lord hath commanded, for I know that the Lord giveth no commandments unto the children of men, save he shall prepare a way for them that they may accomplish the thing which he commandeth them.",
			},
		},
		"John": {
			3: {
				16: "For God so loved the world, that he gave his only begotten Son, that whosoever believeth in him should not perish, but have everlasting life.",
			},
		},
		"Moroni": {
			10: {
				4: "And when ye shall receive these things, I would exhort you that ye would ask God, the Eternal Father, in the name of Christ, if these things are not true; and if ye shall ask with a sincere heart, with real intent, having faith in Christ, he will manifest the truth of it unto you, by the power of the Holy Ghost.",
			},
		},
	}
	
	var results []Scripture
	
	if bookData, exists := sampleData[ref.Book]; exists {
		if chapterData, exists := bookData[ref.Chapter]; exists {
			for verse := ref.Verse; verse <= ref.EndVerse; verse++ {
				if text, exists := chapterData[verse]; exists {
					results = append(results, Scripture{
						Book:    ref.Book,
						Chapter: ref.Chapter,
						Verse:   verse,
						Text:    text,
					})
				}
			}
		}
	}
	
	return results
}

// getChapter retrieves an entire chapter
func (s *Service) getChapter(book string, chapter int) []Scripture {
	// Sample chapter data - would query database in real version
	if book == "1 Nephi" && chapter == 3 {
		return []Scripture{
			{Book: "1 Nephi", Chapter: 3, Verse: 1, Text: "And it came to pass that I, Nephi, returned from speaking with the Lord, to the tent of my father."},
			{Book: "1 Nephi", Chapter: 3, Verse: 2, Text: "And it came to pass that he spake unto me, saying: Behold I have dreamed a dream, wherein the Lord hath commanded me that thou and thy brethren shall return to Jerusalem."},
			{Book: "1 Nephi", Chapter: 3, Verse: 3, Text: "For behold, Laban hath the record of the Jews and also a genealogy of my forefathers, and they are engraven upon plates of brass."},
			{Book: "1 Nephi", Chapter: 3, Verse: 4, Text: "Wherefore, the Lord hath commanded me that thou and thy brothers should go unto the house of Laban, and seek the records, and bring them down hither into the wilderness."},
			{Book: "1 Nephi", Chapter: 3, Verse: 5, Text: "And now, behold thy brothers murmur, saying it is a hard thing which I have required of them; but behold I have not required it of them, but it is a commandment of the Lord."},
			{Book: "1 Nephi", Chapter: 3, Verse: 6, Text: "Therefore go, my son, and thou shalt be favored of the Lord, because thou hast not murmured."},
			{Book: "1 Nephi", Chapter: 3, Verse: 7, Text: "And it came to pass that I, Nephi, said unto my father: I will go and do the things which the Lord hath commanded, for I know that the Lord giveth no commandments unto the children of men, save he shall prepare a way for them that they may accomplish the thing which he commandeth them."},
		}
	}
	
	return []Scripture{}
}

// parseInt converts string to int, returns 0 if conversion fails
func parseInt(s string) int {
	var result int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int(r-'0')
		}
	}
	return result
}