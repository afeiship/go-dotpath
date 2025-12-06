package dotpath

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	originalPath string // Track the original file path for Save() method
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
	dio.originalPath = path // Store the original path
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


// saveStandard performs standard file save without comment preservation
func (dio *DotIO) saveStandard(path string) error {
	data, err := dio.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// SaveAs saves data to specified path with preservation options
// For YAML files, always attempts to preserve comments and order
func (dio *DotIO) SaveAs(path string) error {
	// For YAML files, always preserve comments and order
	if dio.format == YAML {
		// Check if target file exists
		if _, err := os.Stat(path); err == nil {
			// File exists, preserve its comments and order
			return dio.saveWithPreservation(path, path)
		} else {
			// File doesn't exist, try to preserve from original file if available
			if dio.originalPath != "" {
				return dio.saveWithPreservation(path, dio.originalPath)
			}
		}
	}

	// For JSON files, use standard save
	return dio.saveStandard(path)
}

// saveWithPreservation saves YAML file preserving comments from specified source file
func (dio *DotIO) saveWithPreservation(targetPath, sourcePath string) error {
	// Read source file for comment preservation
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Generate new YAML content
	newYAML, err := dio.ToString()
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	// Preserve comments and order from source file
	preservedYAML := mergeYAMLWithPreservedOrder(string(sourceContent), newYAML)

	// Write to target file
	if err := os.WriteFile(targetPath, []byte(preservedYAML), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", targetPath, err)
	}

	return nil
}

// Save saves data back to the original file path
// For YAML files, attempts to preserve comments and order
// Returns an error if the adapter was not loaded from a file
func (dio *DotIO) Save() error {
	if dio.originalPath == "" {
		return fmt.Errorf("no original file path to save to - use SaveAs() instead")
	}

	// Simply delegate to SaveAs with the original path
	return dio.SaveAs(dio.originalPath)
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
		// Use a custom YAML marshaller for better formatting
		return dio.marshalYAML(data)
	default:
		return nil, fmt.Errorf("unsupported format: %s", dio.format)
	}
}

// marshalYAML provides better YAML formatting with consistent indentation
func (dio *DotIO) marshalYAML(data any) ([]byte, error) {
	// Use yaml.Marshal to generate YAML, then apply formatting optimizations
	result, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Apply simple formatting improvements
	return formatYAML(string(result)), nil
}

// formatYAML applies formatting improvements to YAML output
func formatYAML(yamlStr string) []byte {
	// Remove trailing whitespace from each line
	lines := strings.Split(yamlStr, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}

	// Join lines and ensure single trailing newline
	result := strings.Join(lines, "\n")
	// Remove multiple trailing newlines
	result = strings.TrimRight(result, "\n")
	if result != "" {
		result += "\n"
	}

	return []byte(result)
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

// YAML preservation types and functions

type yamlSection struct {
	comment string
	key     string
	indent  string
	content []string
}

// mergeYAMLWithPreservedOrder merges new YAML data with original preserving comments and order
func mergeYAMLWithPreservedOrder(original, newYAML string) string {
	originalLines := strings.Split(original, "\n")

	var sections []yamlSection
	var currentSection *yamlSection

	// Parse original YAML sections
	for i, line := range originalLines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			// This is a comment, start new section
			if currentSection != nil {
				sections = append(sections, *currentSection)
			}
			currentSection = &yamlSection{
				comment: line,
				indent:  extractIndent(line),
			}
		} else if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			// This is content
			if currentSection == nil {
				currentSection = &yamlSection{}
			}

			if currentSection.key == "" {
				currentSection.key = extractKey(line)
			}
			currentSection.content = append(currentSection.content, line)

			// Look ahead for nested content
			if i+1 < len(originalLines) {
				nextLine := originalLines[i+1]
				if strings.HasPrefix(nextLine, " ") && !strings.HasPrefix(strings.TrimSpace(nextLine), "#") {
					// Nested content
					for j := i + 1; j < len(originalLines); j++ {
						nestedLine := originalLines[j]
						if strings.HasPrefix(nestedLine, " ") && !strings.HasPrefix(strings.TrimSpace(nestedLine), "#") {
							currentSection.content = append(currentSection.content, nestedLine)
							i = j
						} else {
							break
						}
					}
				}
			}
			sections = append(sections, *currentSection)
			currentSection = nil
		}
	}

	// Parse new YAML data into map for easy lookup
	newData := parseYAMLToMap(newYAML)

	// Reconstruct with preserved order and updated values
	var result []string
	for _, section := range sections {
		if section.comment != "" {
			result = append(result, section.comment)
		}

		if section.key != "" {
			if newValue, exists := newData[section.key]; exists {
				result = append(result, newValue.content...)
			}
		}
	}

	// Add any new sections that weren't in original
	for _, section := range sections {
		if section.key == "" && len(section.content) > 0 {
			// This might be new content without comment
			result = append(result, section.content...)
		}
	}

	return strings.Join(result, "\n")
}

// parseYAMLToMap parses YAML string into a map of sections
func parseYAMLToMap(yamlStr string) map[string]*yamlSection {
	lines := strings.Split(yamlStr, "\n")
	result := make(map[string]*yamlSection)

	var currentSection *yamlSection

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(trimmed, "#") {
			// Top-level key
			if currentSection != nil {
				result[currentSection.key] = currentSection
			}

			key := extractKey(line)
			currentSection = &yamlSection{
				key:     key,
				content: []string{line},
			}
		} else if currentSection != nil {
			// Nested content
			currentSection.content = append(currentSection.content, line)
		}
	}

	if currentSection != nil {
		result[currentSection.key] = currentSection
	}

	return result
}

// extractIndent extracts leading spaces from a line
func extractIndent(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " "))]
}

// extractKey extracts the key from a YAML line
func extractKey(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return ""
	}

	// Handle nested keys
	parts := strings.Split(trimmed, ":")
	if len(parts) > 0 {
		key := strings.TrimSpace(parts[0])
		// Remove list item marker if present
		return strings.TrimPrefix(key, "- ")
	}

	return ""
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

