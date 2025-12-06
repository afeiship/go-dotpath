# Go-DotPath Adapters

This document describes the JSON and YAML adapters for the go-dotpath package.

## Overview

The adapters provide seamless integration between `map[string]any` data structures and JSON/YAML formats, allowing you to easily load, modify, and save configuration files.

## JSON Adapter

### Installation

The JSON adapter is built into the main package and requires no additional dependencies.

### Functions

```go
// Load functions
dp, err := dotpath.LoadJSONFromBytes([]byte(`{"key": "value"}`))
dp, err := dotpath.LoadJSONFromString(`{"key": "value"}`)
dp, err := dotpath.LoadJSONFromFile("config.json")

// Save functions
err := dp.SaveToJSONFile("output.json")
bytes, err := dp.ToJSONBytes()
str, err := dp.ToJSONString()
```

### Example Usage

```go
// Load configuration
dp, err := dotpath.LoadJSONFromFile("config.json")
if err != nil {
    log.Fatal(err)
}

// Access data
port := dp.GetInt("server.port")
debug := dp.GetBool("debug")

// Modify data
dp.Set("server.port", 9090)
dp.Set("new.feature", true)

// Save back to file
err = dp.SaveToJSONFile("config.json")
```

## YAML Adapter

### Installation

First, install the YAML dependency:

```bash
go get gopkg.in/yaml.v3
```

### Functions

```go
// Load functions
dp, err := dotpath.LoadYAMLFromBytes([]byte("key: value"))
dp, err := dotpath.LoadYAMLFromString("key: value")
dp, err := dotpath.LoadYAMLFromFile("config.yaml")

// Save functions
err := dp.SaveToYAMLFile("output.yaml")
bytes, err := dp.ToYAMLBytes()
str, err := dp.ToYAMLString()
```

### Example Usage

```go
// Load configuration
dp, err := dotpath.LoadYAMLFromFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

// Access data
port := dp.GetInt("server.port")
debug := dp.GetBool("debug")

// Modify data
dp.Set("server.port", 9090)
dp.Set("new.feature", true)

// Save back to file
err = dp.SaveToYAMLFile("config.yaml")
```

## Configuration Management Pattern

Both adapters support a common pattern for configuration management:

### Environment-Specific Configurations

```go
// Base configuration
baseConfig := `
app:
  name: "my-app"
  version: "1.0.0"
`

// Development overrides
devConfig := `
database:
  host: "localhost"
debug: true
`

// Production overrides
prodConfig := `
database:
  host: "prod.example.com"
debug: false
`

// Load base and merge with environment-specific config
dp, _ := dotpath.LoadYAMLFromString(baseConfig)
envDP, _ := dotpath.LoadYAMLFromString(prodConfig)
dp.Update(envDP.Data())
```

### File Operations

```go
// Load from file
dp, err := dotpath.LoadJSONFromFile("config.json")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Make modifications
dp.Set("app.version", "2.0.0")
dp.Set("feature.enabled", true)

// Create backup
dp.SaveToJSONFile("config.backup.json")

// Save updated config
dp.SaveToJSONFile("config.json")
```

## Error Handling

Both adapters provide comprehensive error handling:

```go
dp, err := dotpath.LoadJSONFromFile("config.json")
if err != nil {
    switch {
    case os.IsNotExist(err):
        log.Println("Config file not found, using defaults")
        dp = dotpath.New(defaultConfig)
    case strings.Contains(err.Error(), "invalid character"):
        log.Printf("Invalid JSON in config file: %v", err)
        return
    default:
        log.Fatalf("Unexpected error loading config: %v", err)
    }
}
```

## Performance Considerations

- Both adapters use streaming operations for large files
- JSON adapter uses standard library `encoding/json`
- YAML adapter uses `gopkg.in/yaml.v3` for better performance
- Consider format based on your use case:
  - JSON: Better performance, machine-readable
  - YAML: More human-readable, supports comments

## Examples

See the `examples/` directory for complete working examples:

- `json_example.go` - Comprehensive JSON adapter demo
- `yaml_example.go` - YAML adapter demo (requires build tag)
- `adapters_example.go` - Side-by-side comparison

## Testing

Run tests for adapters:

```bash
# Test JSON adapter
go test -run "TestDotpath.*JSON"

# Test YAML adapter
go test -run "TestDotpath.*YAML"

# Test all functionality
go test -v
```

## Best Practices

1. **Use appropriate format**: JSON for APIs/machine consumption, YAML for human-edited configs
2. **Handle errors gracefully**: Always check error returns from adapter functions
3. **Create backups**: Save configuration backups before making changes
4. **Use environment-specific configs**: Merge base config with environment overrides
5. **Validate configurations**: Use dotpath access to validate required fields exist