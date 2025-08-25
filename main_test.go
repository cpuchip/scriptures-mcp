package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestMainFunction(t *testing.T) {
	// This is mainly to ensure main() can be called without panicking
	// In a real test, we would mock the server.ServeStdio call
	// For now, we just test that the main function exists and is structured correctly
	
	// We can't easily test main() directly since it calls server.ServeStdio which blocks
	// Instead, we'll test that we can import and the necessary components exist
	
	// This test ensures the main package compiles correctly
	t.Log("Main package compiles successfully")
}