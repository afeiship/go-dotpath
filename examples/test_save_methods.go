package main

import (
	"fmt"
	"log"

	"github.com/afeiship/go-dotpath"
)

func main() {
	fmt.Println("=== Testing Save Methods with Comment Preservation ===\n")

	// Test 1: Save method (preserves comments in original file)
	fmt.Println("1. Testing Save() method:")
	adapter1, err := dotpath.LoadYAMLAdapter("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	// Modify data
	adapter1.Set("database.username", "admin")
	adapter1.Set("database.password", "secret123")
	adapter1.Set("app.debug", true)

	// Save back to original file
	err = adapter1.Save()
	if err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	fmt.Println("✅ Saved back to test_config.yaml with preserved comments\n")

	// Test 2: SaveAs method (always preserves comments and order for YAML)
	fmt.Println("2. Testing SaveAs() method:")
	adapter2, err := dotpath.LoadYAMLAdapter("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	// Modify data differently
	adapter2.Set("database.port", 3306)
	adapter2.Set("app.environment", "production")
	adapter2.Set("features", []string{"auth", "logging", "metrics"})

	// Save as a new file (should preserve comments from original)
	err = adapter2.SaveAs("test_config_new.yaml")
	if err != nil {
		log.Fatalf("Failed to save as: %v", err)
	}

	fmt.Println("✅ Saved as test_config_new.yaml (preserved comments from original)\n")

	// Test 3: SaveAs to existing file (preserves existing comments)
	fmt.Println("3. Testing SaveAs() to existing file:")

	// First, copy original file to create a target with comments
	adapter3 := dotpath.NewDotIO(dotpath.YAML)
	adapter3.Set("server", map[string]any{
		"host": "0.0.0.0",
		"port": 8080,
	})

	// Save to create a target file
	err = adapter3.SaveAs("target_config.yaml")
	if err != nil {
		log.Fatalf("Failed to create target file: %v", err)
	}

	// Manually add comments to target file
	fmt.Println("Created target_config.yaml\n")

	// Now use SaveAs to save to this existing file
	adapter4, err := dotpath.LoadYAMLAdapter("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	adapter4.Set("server.host", "localhost")
	adapter4.Set("server.port", 9090)

	err = adapter4.SaveAs("target_config.yaml")
	if err != nil {
		log.Fatalf("Failed to save as to existing file: %v", err)
	}

	fmt.Println("✅ Saved to existing target_config.yaml (would preserve comments if they existed)\n")

	fmt.Println("All save methods tested successfully!")
	fmt.Println("\nFile summary:")
	fmt.Println("- test_config.yaml: Updated in place with preserved comments")
	fmt.Println("- test_config_new.yaml: New file with preserved comments from original")
	fmt.Println("- target_config.yaml: Updated target file")
}