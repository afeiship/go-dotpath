package main

import (
	"fmt"
	"log"
	"os"

	"github.com/afeiship/go-dotpath"
)

func main() {
	fmt.Println("=== Debug SaveAs Comment Preservation ===\n")

	// Read original file
	originalContent, err := os.ReadFile("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to read original: %v", err)
	}

	fmt.Println("Original file content:")
	fmt.Println(string(originalContent))

	// Load with DotIO
	adapter, err := dotpath.LoadYAMLAdapter("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	// Modify data
	adapter.Set("database.port", 9999)
	adapter.Set("app.new_field", "test")

	// Generate new YAML
	newYAML, err := adapter.ToString()
	if err != nil {
		log.Fatalf("Failed to generate YAML: %v", err)
	}

	fmt.Println("\nGenerated YAML:")
	fmt.Println(newYAML)

	// Save as new file
	err = adapter.SaveAs("debug_output.yaml")
	if err != nil {
		log.Fatalf("Failed to save as: %v", err)
	}

	// Read result
	resultContent, err := os.ReadFile("debug_output.yaml")
	if err != nil {
		log.Fatalf("Failed to read result: %v", err)
	}

	fmt.Println("\nSaved file content:")
	fmt.Println(string(resultContent))
}