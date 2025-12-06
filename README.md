# go-dotpath

> Dot-notation access and mutation for nested Go maps.

## Installation

```sh
go get -u github.com/afeiship/go-dotpath
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/afeiship/go-dotpath"
)

func main() {
    // Create configuration
    config := map[string]any{
        "app": map[string]any{
            "name": "demo",
        },
    }

    dp := dotpath.New(config)

    // Get values
    name := dp.GetString("app.name")  // "demo"
    fmt.Println(name)

    // Set values (auto-creates nested maps)
    dp.Set("server.port", 8080)
    port := dp.GetInt("server.port")  // 8080
    fmt.Println(port)

    // Check if path exists
    if dp.Has("server.host") {
        fmt.Println("Host exists")
    }

    // Update with deep merge
    updates := map[string]any{
        "app": map[string]any{
            "version": "1.0.0",
        },
    }
    dp.Update(updates)
}
```

## JSON/YAML Support

Built-in support for both JSON and YAML with unified DotIO interface:

```go
// Create DotIO with format
io := dotpath.NewDotIO(dotpath.JSON)  // or dotpath.YAML

// Load data from various sources
io.LoadFromFile("config.json")
io.LoadFromString(`{"key": "value"}`)
io.LoadFromBytes([]byte("key: value"))

// Use dot-path operations
io.Set("app.version", "2.0.0")
port := io.GetInt("server.port")

// Save in any format (dynamic format switching!)
io.SetFormat(dotpath.YAML)
io.SaveToFile("config.yaml")
```

### DotIO Methods
- `NewDotIO(format Format) *DotIO`
- `NewDotIOWithData(data map[string]any, format Format) *DotIO`
- `LoadFromFile(path string) error`
- `LoadFromString(s string) error`
- `LoadFromBytes(data []byte) error`
- `SaveToFile(path string) error`
- `ToString() (string, error)`
- `ToBytes() ([]byte, error)`
- `SetFormat(format Format)`
- `GetFormat() Format`
- `DotPath() *DotPath`

### Quick Examples

```go
// JSON usage
io := dotpath.NewDotIO(dotpath.JSON)
io.LoadFromFile("config.json")
io.Set("debug", true)
io.SaveToFile("config.json")

// YAML usage
io := dotpath.NewDotIO(dotpath.YAML)
io.LoadFromFile("config.yaml")
io.Set("timeout", 30)
io.SaveToFile("config.yaml")

// Format switching
io := dotpath.NewDotIO(dotpath.JSON)
io.LoadFromFile("data.json")
io.SetFormat(dotpath.YAML)  // Switch output format
io.SaveToFile("data.yaml")
```

## Core API

### Constructor
- `New(m map[string]any) *DotPath`

### Getters
- `Get(path string) (any, bool)`
- `GetString(path string) string`
- `GetInt(path string) int`
- `GetBool(path string) bool`
- `GetFloat64(path string) float64`
- `Has(path string) bool`

### Setters
- `Set(path string, value any)`
- `Update(other map[string]any)`

### Access
- `Data() map[string]any`