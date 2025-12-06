package main

import (
	"fmt"
	"log"

	"github.com/afeiship/go-dotpath"
)

func main() {
	// Test YAML formatting with the provided config
	fmt.Println("=== Testing YAML Output Formatting ===\n")

	// Load the YAML file
	adapter, err := dotpath.LoadYAMLAdapter("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	// Get the current YAML output
	output, err := adapter.ToString()
	if err != nil {
		log.Fatalf("Failed to convert to string: %v", err)
	}

	fmt.Println("Current YAML output:")
	fmt.Printf("---\n%s---\n\n", output)

	// Test modifying data and checking output
	adapter.Set("database.username", "admin")
	adapter.Set("database.password", "secret")
	adapter.Set("app.debug", true)
	adapter.Set("features", []string{"auth", "logging", "metrics"})

	modifiedOutput, err := adapter.ToString()
	if err != nil {
		log.Fatalf("Failed to get modified output: %v", err)
	}

	fmt.Println("Modified YAML output:")
	fmt.Printf("---\n%s---\n\n", modifiedOutput)

	// Test access to nested values
	fmt.Println("Testing data access:")
	fmt.Printf("Database host: %s\n", adapter.GetString("database.host"))
	fmt.Printf("Database port: %d\n", adapter.GetInt("database.port"))
	fmt.Printf("App name: %s\n", adapter.GetString("app.name"))
	fmt.Printf("App version: %s\n", adapter.GetString("app.version"))
	fmt.Printf("App debug: %t\n", adapter.GetBool("app.debug"))
}