package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/afeiship/go-dotpath"
)

func main() {
	fmt.Println("=== Testing YAML with Comments and Order Preservation ===\n")

	// Read original file
	originalContent, err := os.ReadFile("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to read original file: %v", err)
	}

	fmt.Println("Original content:")
	fmt.Println(string(originalContent))

	// Load original YAML with DotIO
	adapter := dotpath.NewDotIO(dotpath.YAML)
	err = adapter.LoadFromFile("test_config.yaml")
	if err != nil {
		log.Fatalf("Failed to load YAML: %v", err)
	}

	// Modify some data
	adapter.Set("database.username", "admin")
	adapter.Set("database.password", "secret")
	adapter.Set("app.debug", true)
	adapter.Set("app.environment", "production")

	// Generate new YAML output
	yamlOutput, err := adapter.ToString()
	if err != nil {
		log.Fatalf("Failed to generate YAML: %v", err)
	}

	fmt.Println("\nStandard YAML output:")
	fmt.Println(yamlOutput)

	// Create output with preserved comments and order
	finalOutput := mergeYAMLWithPreservedOrder(string(originalContent), yamlOutput)

	fmt.Println("\nMerged output with preserved comments and order:")
	fmt.Println(finalOutput)

	// Save to new file
	err = os.WriteFile("test_config_with_comments.yaml", []byte(finalOutput), 0644)
	if err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}

	fmt.Println("\n✅ Saved to test_config_with_comments.yaml")
}

type YAMLSection struct {
	comment string
	key     string
	indent  string
	content []string
}

func mergeYAMLWithPreservedOrder(original, newYAML string) string {
	originalLines := strings.Split(original, "\n")

	var sections []YAMLSection
	var currentSection *YAMLSection

	// Parse original YAML sections
	for i, line := range originalLines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			// This is a comment, start new section
			if currentSection != nil {
				sections = append(sections, *currentSection)
			}
			currentSection = &YAMLSection{
				comment: line,
				indent:  extractIndent(line),
			}
		} else if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			// This is content
			if currentSection == nil {
				currentSection = &YAMLSection{}
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

func parseYAMLToMap(yamlStr string) map[string]*YAMLSection {
	lines := strings.Split(yamlStr, "\n")
	result := make(map[string]*YAMLSection)

	var currentSection *YAMLSection

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
			currentSection = &YAMLSection{
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

func extractIndent(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " "))]
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

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