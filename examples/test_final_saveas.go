package main

import (
	"fmt"
	"log"

	"github.com/afeiship/go-dotpath"
)

func main() {
	fmt.Println("=== Final SaveAs Test with Comment Preservation ===\n")

	// Load file with comments
	adapter, err := dotpath.LoadYAMLAdapter("test_config_with_comments.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	fmt.Println("Original file:")
	fmt.Println("----------------------------------------")
	content, _ := adapter.ToString()
	fmt.Println(content)

	// Modify data
	adapter.Set("database.username", "admin")
	adapter.Set("app.debug", true)
	adapter.Set("app.new_feature", "enabled")

	// SaveAs to new file (should preserve comments from original)
	err = adapter.SaveAs("final_output.yaml")
	if err != nil {
		log.Fatalf("Failed to save as: %v", err)
	}

	fmt.Println("✅ Saved as final_output.yaml (should preserve comments)")
	fmt.Println("Please check final_output.yaml to verify comments are preserved!")
}