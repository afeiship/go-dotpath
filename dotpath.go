package dotpath

// DotPath wraps map[string]any to provide dot-notation access and mutation
type DotPath struct {
	data map[string]any
}

// New creates a new DotPath instance wrapping the provided map
func New(m map[string]any) *DotPath {
	if m == nil {
		m = make(map[string]any)
	}
	return &DotPath{data: m}
}

// Get retrieves a value by dot path and returns the value along with a boolean indicating success
func (dp *DotPath) Get(path string) (any, bool) {
	if path == "" {
		return dp.data, true
	}

	parts := splitPath(path)
	current := dp.data

	for i, key := range parts {
		if i == len(parts)-1 {
			// Last part - return the value
			if val, exists := current[key]; exists {
				return val, true
			}
			return nil, false
		}

		// Navigate to nested map
		if val, exists := current[key]; exists {
			if nextMap, ok := val.(map[string]any); ok {
				current = nextMap
			} else {
				// Path exists but is not a map
				return nil, false
			}
		} else {
			// Path doesn't exist
			return nil, false
		}
	}

	return nil, false
}

// GetString retrieves a string value by dot path, returns empty string if not found or type mismatch
func (dp *DotPath) GetString(path string) string {
	if val, exists := dp.Get(path); exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt retrieves an int value by dot path, returns 0 if not found or type mismatch
func (dp *DotPath) GetInt(path string) int {
	if val, exists := dp.Get(path); exists {
		// Handle both int and float64 that are whole numbers
		switch v := val.(type) {
		case int:
			return v
		case float64:
			if v == float64(int(v)) {
				return int(v)
			}
		}
	}
	return 0
}

// GetBool retrieves a bool value by dot path, returns false if not found or type mismatch
func (dp *DotPath) GetBool(path string) bool {
	if val, exists := dp.Get(path); exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// GetFloat64 retrieves a float64 value by dot path, returns 0.0 if not found or type mismatch
func (dp *DotPath) GetFloat64(path string) float64 {
	if val, exists := dp.Get(path); exists {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return 0.0
}

// Has checks if a path exists in the data structure
func (dp *DotPath) Has(path string) bool {
	_, exists := dp.Get(path)
	return exists
}

// Set sets a value by dot path, automatically creating nested maps as needed
func (dp *DotPath) Set(path string, value any) {
	if path == "" {
		return
	}

	parts := splitPath(path)
	current := dp.data

	for i, key := range parts {
		if i == len(parts)-1 {
			// Last part - set the value
			current[key] = value
			return
		}

		// Navigate or create nested map
		if val, exists := current[key]; exists {
			if nextMap, ok := val.(map[string]any); ok {
				current = nextMap
			} else {
				// Path exists but is not a map, replace with new map
				newMap := make(map[string]any)
				current[key] = newMap
				current = newMap
			}
		} else {
			// Create new nested map
			newMap := make(map[string]any)
			current[key] = newMap
			current = newMap
		}
	}
}

// Update performs a deep recursive merge of another map into the current data
func (dp *DotPath) Update(other map[string]any) {
	if other == nil {
		return
	}
	dp.data = deepMerge(dp.data, other)
}

// Data returns the underlying map for serialization, inspection, etc.
func (dp *DotPath) Data() map[string]any {
	return dp.data
}

// splitPath splits a dot path into individual parts
func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	// Simple split on dots - this is sufficient for the requirements
	parts := make([]string, 0)
	start := 0

	for i, char := range path {
		if char == '.' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}

	// Add the last part
	if start < len(path) {
		parts = append(parts, path[start:])
	}

	return parts
}

// deepMerge recursively merges two maps, with values from src taking precedence
func deepMerge(dest, src map[string]any) map[string]any {
	if dest == nil {
		dest = make(map[string]any)
	}

	for key, srcVal := range src {
		if destVal, exists := dest[key]; exists {
			// If both values are maps, merge them recursively
			if destMap, ok := destVal.(map[string]any); ok {
				if srcMap, ok := srcVal.(map[string]any); ok {
					dest[key] = deepMerge(destMap, srcMap)
					continue
				}
			}
		}
		// Either key doesn't exist or values are not both maps - use src value
		dest[key] = srcVal
	}

	return dest
}