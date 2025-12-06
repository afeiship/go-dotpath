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

## API

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