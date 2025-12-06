package dotpath

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Format represents the data format for the adapter
type Format string

const (
	JSON Format = "json"
	YAML Format = "yaml"
)

// DotIO provides a unified interface for JSON/YAML operations
type DotIO struct {
	dp    *DotPath
	format Format
}

// NewDotIO creates a new DotIO with the specified format
func NewDotIO(format Format) *DotIO {
	return &DotIO{
		dp:    New(nil),
		format: format,
	}
}

// NewDotIOWithData creates a new DotIO with existing data
func NewDotIOWithData(data map[string]any, format Format) *DotIO {
	return &DotIO{
		dp:    New(data),
		format: format,
	}
}

// LoadFromFile loads data from file
func (dio *DotIO) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var result map[string]any

	switch dio.format {
	case JSON:
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	case YAML:
		if err := yaml.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", dio.format)
	}

	dio.dp = New(result)
	return nil
}

// LoadFromString loads data from string
func (dio *DotIO) LoadFromString(s string) error {
	return dio.LoadFromBytes([]byte(s))
}

// LoadFromBytes loads data from bytes
func (dio *DotIO) LoadFromBytes(data []byte) error {
	var result map[string]any

	switch dio.format {
	case JSON:
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	case YAML:
		if err := yaml.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", dio.format)
	}

	dio.dp = New(result)
	return nil
}

// SaveToFile saves data to file
func (dio *DotIO) SaveToFile(path string) error {
	data, err := dio.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// ToString converts data to string
func (dio *DotIO) ToString() (string, error) {
	data, err := dio.ToBytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToBytes converts data to bytes
func (dio *DotIO) ToBytes() ([]byte, error) {
	data := dio.dp.Data()

	switch dio.format {
	case JSON:
		return json.MarshalIndent(data, "", "  ")
	case YAML:
		return yaml.Marshal(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", dio.format)
	}
}

// GetFormat returns the current format
func (dio *DotIO) GetFormat() Format {
	return dio.format
}

// SetFormat changes the format
func (dio *DotIO) SetFormat(format Format) {
	dio.format = format
}

// Get the underlying DotPath instance for all dot-path operations
func (dio *DotIO) DotPath() *DotPath {
	return dio.dp
}

// Convenience methods that delegate to DotPath

// Get retrieves a value by path
func (dio *DotIO) Get(path string) (any, bool) {
	return dio.dp.Get(path)
}

// GetString retrieves a string value by path
func (dio *DotIO) GetString(path string) string {
	return dio.dp.GetString(path)
}

// GetInt retrieves an int value by path
func (dio *DotIO) GetInt(path string) int {
	return dio.dp.GetInt(path)
}

// GetBool retrieves a bool value by path
func (dio *DotIO) GetBool(path string) bool {
	return dio.dp.GetBool(path)
}

// GetFloat64 retrieves a float64 value by path
func (dio *DotIO) GetFloat64(path string) float64 {
	return dio.dp.GetFloat64(path)
}

// Has checks if a path exists
func (dio *DotIO) Has(path string) bool {
	return dio.dp.Has(path)
}

// Set sets a value by path
func (dio *DotIO) Set(path string, value any) {
	dio.dp.Set(path, value)
}

// Update performs a deep merge with another map
func (dio *DotIO) Update(other map[string]any) {
	dio.dp.Update(other)
}

// Data returns the underlying map
func (dio *DotIO) Data() map[string]any {
	return dio.dp.Data()
}

// Static helper functions for creating adapters from existing files

// LoadJSONAdapter creates a JSON DotIO from file
func LoadJSONAdapter(path string) (*DotIO, error) {
	adapter := NewDotIO(JSON)
	if err := adapter.LoadFromFile(path); err != nil {
		return nil, err
	}
	return adapter, nil
}

// LoadYAMLAdapter creates a YAML DotIO from file
func LoadYAMLAdapter(path string) (*DotIO, error) {
	adapter := NewDotIO(YAML)
	if err := adapter.LoadFromFile(path); err != nil {
		return nil, err
	}
	return adapter, nil
}

