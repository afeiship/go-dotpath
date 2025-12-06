# DotIO - JSON/YAML Adapter Documentation

The `DotIO` struct provides a unified interface for handling JSON and YAML data formats using the go-dotpath library. It combines the power of dot-path navigation with convenient file I/O operations.

## Overview

`DotIO` acts as an adapter that wraps the core `DotPath` functionality and adds format-specific serialization/deserialization capabilities. It supports both JSON and YAML formats, allowing you to work with configuration files, data interchange formats, and more.

## Features

- **Dual format support**: JSON and YAML
- **File I/O operations**: Load from and save to files
- **String operations**: Work with in-memory strings
- **Byte operations**: Work with raw byte data
- **Format switching**: Change between JSON and YAML at runtime
- **Full dot-path access**: All `DotPath` methods available

## Basic Usage

### Creating a DotIO Instance

```go
// Create with empty data
jsonIO := dotpath.NewDotIO(dotpath.JSON)
yamlIO := dotpath.NewDotIO(dotpath.YAML)

// Create with existing data
data := map[string]any{
    "server": map[string]any{
        "port": 8080,
        "host": "localhost",
    },
}
io := dotpath.NewDotIOWithData(data, dotpath.JSON)
```

### Loading Data

```go
// From file
err := jsonIO.LoadFromFile("config.json")

// From string
jsonStr := `{"server": {"port": 8080}}`
err := jsonIO.LoadFromString(jsonStr)

// From bytes
err := jsonIO.LoadFromBytes([]byte(jsonStr))
```

### Saving Data

```go
// To file with intelligent comment preservation (YAML)
err := yamlIO.Save()           // Save back to original file
err := yamlIO.SaveAs("new.yaml")  // Save to new path

// To file (standard format)
err := jsonIO.SaveAs("output.json")

// To string
str, err := jsonIO.ToString()

// To bytes
data, err := jsonIO.ToBytes()
```

### Working with Dot Paths

```go
// Get values
port := jsonIO.GetInt("server.port")
host := jsonIO.GetString("server.host")
debug := jsonIO.GetBool("debug.enabled")

// Set values
jsonIO.Set("server.timeout", 30)
jsonIO.Set("features.newFeature", true)

// Check path existence
if jsonIO.Has("database") {
    // database path exists
}

// Update with deep merge
newData := map[string]any{
    "server": map[string]any{
        "ssl": true,
    },
}
jsonIO.Update(newData)
```

### Format Switching

```go
// Create as JSON
io := dotpath.NewDotIO(dotpath.JSON)
io.LoadFromFile("config.json")

// Switch to YAML format
io.SetFormat(dotpath.YAML)
io.SaveToFile("config.yaml")
```

## API Reference

### Constructors

- `NewDotIO(format Format) *DotIO` - Creates new DotIO with empty data
- `NewDotIOWithData(data map[string]any, format Format) *DotIO` - Creates with existing data

### Format Constants

- `JSON` - JSON format
- `YAML` - YAML format

### Loading Methods

- `LoadFromFile(path string) error` - Load from file
- `LoadFromString(s string) error` - Load from string
- `LoadFromBytes(data []byte) error` - Load from byte slice

### Saving Methods

- `Save() error` - Save back to original file (YAML: preserves comments & order)
- `SaveAs(path string) error` - Save to specified path (YAML: preserves comments & order)
- `ToString() (string, error)` - Convert to string
- `ToBytes() ([]byte, error)` - Convert to byte slice

### Format Methods

- `GetFormat() Format` - Get current format
- `SetFormat(format Format)` - Change format

### DotPath Access

- `DotPath() *DotPath` - Get underlying DotPath instance

### Convenience Methods (delegated to DotPath)

- `Get(path string) (any, bool)` - Get any value
- `GetString(path string) string` - Get string value
- `GetInt(path string) int` - Get int value
- `GetBool(path string) bool` - Get bool value
- `GetFloat64(path string) float64` - Get float64 value
- `Has(path string) bool` - Check if path exists
- `Set(path string, value any)` - Set value
- `Update(other map[string]any)` - Deep merge update
- `Data() map[string]any` - Get underlying data

### Static Helper Functions

- `LoadJSONAdapter(path string) (*DotIO, error)` - Quick JSON file loader
- `LoadYAMLAdapter(path string) (*DotIO, error)` - Quick YAML file loader

## Error Handling

All methods return descriptive errors:

```go
if err := jsonIO.LoadFromFile("config.json"); err != nil {
    // Handle specific error types
    if os.IsNotExist(err) {
        log.Fatal("Config file not found")
    }
    log.Fatalf("Failed to load config: %v", err)
}
```

## Complete Example

```go
package main

import (
    "log"
    "github.com/feizheng/go-dotpath"
)

func main() {
    // Load JSON configuration
    config, err := dotpath.LoadJSONAdapter("config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Read configuration values
    port := config.GetInt("server.port")
    host := config.GetString("server.host")

    log.Printf("Server: %s:%d", host, port)

    // Update configuration
    config.Set("server.timeout", 30)
    config.Set("debug.enabled", true)

    // Save as YAML
    config.SetFormat(dotpath.YAML)
    if err := config.SaveAs("config.yaml"); err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }

    log.Println("Configuration saved as YAML")
}
```

## YAML Comment & Order Preservation

The DotIO adapter provides intelligent YAML file handling that preserves comments and field order when saving files.

### Features

- **Comment Preservation**: All YAML comments are preserved during save operations
- **Order Preservation**: Field order is maintained from the original file
- **Smart Source Selection**: Automatically chooses the best source for comment preservation

### Behavior

| Method | Target Exists | Comment Source | Use Case |
|--------|---------------|---------------|----------|
| `Save()` | N/A | Original file | Update the file you loaded from |
| `SaveAs(path)` | No | Original file | Create new copy with original comments |
| `SaveAs(path)` | Yes | Target file | Update existing file preserving its comments |

### Example

```go
// Original config.yaml has comments:
// # Database configuration
// database:
//   host: localhost
//   port: 5432
// # App settings
// app:
//   name: myapp

adapter, err := dotpath.LoadYAMLAdapter("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Modify data
adapter.Set("database.username", "admin")
adapter.Set("app.debug", true)

// Save back - preserves original comments and order
err = adapter.Save()
// Result: Comments and order preserved, new fields added appropriately

// Save as new file - preserves comments from original
err = adapter.SaveAs("backup.yaml")
// Result: New file has all comments from original config.yaml
```

### Limitations

- **JSON**: No comment preservation (JSON doesn't support comments)
- **YAML**: Only preserves line comments starting with `#`
- **Structure**: Major structural changes may affect comment placement

## Integration with Core Library

`DotIO` is built on top of the core `DotPath` library and provides a layer of convenience for file operations and format handling. All dot-path operations work identically to the core library:

```go
// These are equivalent:
port1 := config.GetInt("server.port")           // Via DotIO
port2 := config.DotPath().GetInt("server.port")  // Via underlying DotPath
```

## Best Practices

1. **Use appropriate format**: JSON for APIs and web configs, YAML for human-readable configs
2. **Handle errors properly**: All file operations can fail
3. **Batch operations**: Load once, then perform multiple dot-path operations
4. **Format switching**: Use `SetFormat()` to convert between formats without data loss
5. **Static helpers**: Use `LoadJSONAdapter()`/`LoadYAMLAdapter()` for quick file loading

## Dependencies

- Go standard library (`encoding/json`, `os`)
- `gopkg.in/yaml.v3` for YAML support
- Core go-dotpath library