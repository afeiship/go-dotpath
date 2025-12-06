package main

import (
	"fmt"
	"log"

	"github.com/afeiship/go-dotpath"
)

func main() {
	fmt.Println("Go-DotPath Demo")
	fmt.Println("================")

	// Create DotIO with JSON format
	io := dotpath.NewDotIO(dotpath.JSON)

	// Load JSON data
	jsonData := `{
		"app": {
			"name": "my-app",
			"version": "1.0.0"
		},
		"server": {
			"port": 8080,
			"host": "localhost"
		}
	}`

	if err := io.LoadFromString(jsonData); err != nil {
		log.Fatal(err)
	}

	// Use dot-path operations
	fmt.Printf("App: %s v%s\n", io.GetString("app.name"), io.GetString("app.version"))
	fmt.Printf("Server: %s:%d\n", io.GetString("server.host"), io.GetInt("server.port"))

	// Modify data
	io.Set("app.version", "2.0.0")
	io.Set("debug", true)

	// Switch to YAML and save
	io.SetFormat(dotpath.YAML)
	yamlStr, err := io.ToString()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nAs YAML:")
	fmt.Println(yamlStr)

	// Switch back to JSON
	io.SetFormat(dotpath.JSON)
	jsonStr, _ := io.ToString()
	fmt.Println("As JSON:")
	fmt.Println(jsonStr)
}