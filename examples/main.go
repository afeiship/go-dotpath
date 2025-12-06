package main

import (
	"fmt"

	"github.com/afeiship/go-dotpath"
)

func main() {
	fmt.Println("Go-DotPath Example Usage")
	fmt.Println("========================")

	// Create initial configuration
	config := map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "demo-app",
			"version": "1.0.0",
		},
		"debug": true,
	}

	// Create new DotPath instance
	dp := dotpath.New(config)
	fmt.Printf("Initial config: %+v\n\n", dp.Data())

	// Demonstrate getters
	fmt.Println("=== Getter Methods ===")

	// Get string
	appName := dp.GetString("app.name")
	fmt.Printf("App name: %s\n", appName)

	// Get bool
	debug := dp.GetBool("debug")
	fmt.Printf("Debug mode: %t\n", debug)

	// Get non-existent path (returns zero value)
	missing := dp.GetString("app.missing")
	fmt.Printf("Missing value (should be empty): '%s'\n", missing)

	// Check if path exists
	hasMissing := dp.Has("app.missing")
	fmt.Printf("Has 'app.missing': %t\n\n", hasMissing)

	// Demonstrate setters
	fmt.Println("=== Setter Methods ===")

	// Set existing path
	dp.Set("app.version", "2.0.0")
	fmt.Printf("Updated app version: %s\n", dp.GetString("app.version"))

	// Set new nested path (auto-creates maps)
	dp.Set("server.host", "localhost")
	dp.Set("server.port", 8080)
	if server, exists := dp.Get("server"); exists {
		fmt.Printf("Server config: %+v\n", server)
	}

	// Set deeply nested path
	dp.Set("database.connection.pool.max", 100)
	fmt.Printf("Database max pool size: %d\n", dp.GetInt("database.connection.pool.max"))
	if db, exists := dp.Get("database"); exists {
		fmt.Printf("Database config: %+v\n", db)
	}

	// Demonstrate update/merge
	fmt.Println("\n=== Update/Merge ===")

	updates := map[string]interface{}{
		"app": map[string]interface{}{
			"version": "2.1.0", // Overwrites existing
			"author":  "Fei Zheng", // Adds new
		},
		"newFeature": map[string]interface{}{
			"enabled": true,
		},
	}

	dp.Update(updates)
	fmt.Printf("After update: %+v\n", dp.Data())

	// Demonstrate all getter types
	fmt.Println("\n=== All Getter Types ===")

	dp.Set("examples.string", "hello world")
	dp.Set("examples.integer", 42)
	dp.Set("examples.float", 3.14159)
	dp.Set("examples.boolean", true)

	fmt.Printf("String: %s\n", dp.GetString("examples.string"))
	fmt.Printf("Integer: %d\n", dp.GetInt("examples.integer"))
	fmt.Printf("Float64: %f\n", dp.GetFloat64("examples.float"))
	fmt.Printf("Boolean: %t\n", dp.GetBool("examples.boolean"))

	// Demonstrate edge cases
	fmt.Println("\n=== Edge Cases ===")

	// Type mismatch (getting int as string)
	wrongType := dp.GetString("examples.integer")
	fmt.Printf("Getting int as string: '%s' (empty due to type mismatch)\n", wrongType)

	// Deep path that doesn't exist
	if deepMissing, exists := dp.Get("very.deep.nested.path"); exists {
		fmt.Printf("Very deep missing path: %v, exists: %t\n", deepMissing, exists)
	} else {
		fmt.Println("Very deep missing path: (nil), exists: false")
	}

	// Working with the underlying map directly
	fmt.Println("\n=== Direct Map Access ===")

	data := dp.Data()
	if serverConfig, exists := data["server"]; exists {
		if serverMap, ok := serverConfig.(map[string]interface{}); ok {
			fmt.Printf("Server host from raw map: %v\n", serverMap["host"])
		}
	}

	// Example of configuration loading simulation
	fmt.Println("\n=== Simulated Config Loading ===")

	// Simulate loading configuration from different sources
	baseConfig := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "my-app",
		},
		"database": map[string]interface{}{
			"driver": "sqlite",
		},
	}

	prodConfig := map[string]interface{}{
		"database": map[string]interface{}{
			"driver": "postgres",
			"host":   "prod-db.example.com",
			"port":   5432,
		},
		"debug": false,
	}

	// Start with base config
	appConfig := dotpath.New(baseConfig)

	// Update with production config
	appConfig.Update(prodConfig)

	fmt.Printf("Final configuration:\n")
	fmt.Printf("App name: %s\n", appConfig.GetString("app.name"))
	fmt.Printf("Database driver: %s\n", appConfig.GetString("database.driver"))
	fmt.Printf("Database host: %s\n", appConfig.GetString("database.host"))
	fmt.Printf("Database port: %d\n", appConfig.GetInt("database.port"))
	fmt.Printf("Debug mode: %t\n", appConfig.GetBool("debug"))
}

