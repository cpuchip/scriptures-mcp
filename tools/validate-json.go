package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// validateJSONRPC validates a JSON-RPC message and provides detailed error information
func validateJSONRPC(line string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil // Skip empty lines
	}

	// First check if it's valid JSON at all
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		// Check for common double-quote issue in query strings
		if strings.Contains(line, `""`) && strings.Contains(err.Error(), "invalid character") {
			return fmt.Errorf("Invalid JSON syntax - likely caused by extra quotes in query string.\nFound '\"\"' which should be '\"'\nError: %v\nLine: %s", err, line)
		}
		return fmt.Errorf("Invalid JSON syntax: %v\nLine: %s", err, line)
	}

	// Check if it has the basic JSON-RPC structure
	var baseMsg struct {
		JSONRPC string      `json:"jsonrpc"`
		ID      interface{} `json:"id"`
		Method  string      `json:"method"`
	}

	if err := json.Unmarshal(raw, &baseMsg); err != nil {
		return fmt.Errorf("Failed to parse JSON-RPC structure: %v", err)
	}

	if baseMsg.JSONRPC != "2.0" {
		return fmt.Errorf("Invalid JSON-RPC version: expected '2.0', got '%s'", baseMsg.JSONRPC)
	}

	if baseMsg.Method == "" {
		return fmt.Errorf("Missing 'method' field in JSON-RPC message")
	}

	// Additional validation for tools/call method
	if baseMsg.Method == "tools/call" {
		var toolCall struct {
			Params struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			} `json:"params"`
		}
		
		if err := json.Unmarshal(raw, &toolCall); err != nil {
			return fmt.Errorf("Failed to parse tools/call params: %v", err)
		}

		if toolCall.Params.Name == "" {
			return fmt.Errorf("Missing tool 'name' in tools/call params")
		}

		if toolCall.Params.Arguments == nil {
			return fmt.Errorf("Missing 'arguments' in tools/call params")
		}

		// Check for common query argument issues
		if toolCall.Params.Name == "search_scriptures" || toolCall.Params.Name == "get_scripture" || toolCall.Params.Name == "get_chapter" {
			if query, exists := toolCall.Params.Arguments["query"]; exists {
				if queryStr, ok := query.(string); ok {
					if strings.HasPrefix(queryStr, "\"") || strings.HasSuffix(queryStr, "\"") {
						return fmt.Errorf("Query argument contains extra quotes: %q\nRemove the extra quotes around the query string", queryStr)
					}
				}
			} else {
				return fmt.Errorf("Missing required 'query' argument for %s", toolCall.Params.Name)
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("JSON-RPC Validator for scriptures-mcp")
		fmt.Println("Usage:")
		fmt.Println("  echo 'JSON-RPC message' | go run validate-json.go")
		fmt.Println("  go run validate-json.go < input.json")
		fmt.Println("")
		fmt.Println("This tool validates JSON-RPC messages for common formatting errors")
		fmt.Println("that can cause 'Parse error' (-32700) responses.")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	lineNum := 0
	hasErrors := false

	fmt.Println("Validating JSON-RPC messages...")
	fmt.Println("=====================================")

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		if err := validateJSONRPC(line); err != nil {
			fmt.Printf("❌ Line %d: %v\n", lineNum, err)
			hasErrors = true
		} else if strings.TrimSpace(line) != "" {
			fmt.Printf("✅ Line %d: Valid JSON-RPC message\n", lineNum)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=====================================")
	if hasErrors {
		fmt.Println("❌ Validation completed with errors")
		fmt.Println("\nCommon JSON formatting issues:")
		fmt.Println("• Extra quotes: \"\"text\" should be \"text\"")
		fmt.Println("• Missing commas between fields")
		fmt.Println("• Unescaped quotes inside strings")
		fmt.Println("• Missing closing braces or brackets")
		os.Exit(1)
	} else {
		fmt.Println("✅ All JSON-RPC messages are valid")
	}
}