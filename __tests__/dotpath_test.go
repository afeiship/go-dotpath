package dotpath_test

import (
	"testing"

	"github.com/afeiship/go-dotpath"
)

func TestNew(t *testing.T) {
	// Test with nil map
	dp := dotpath.New(nil)
	if dp.Data() == nil {
		t.Error("Expected non-nil map when creating with nil")
	}

	// Test with existing map
	m := map[string]interface{}{"key": "value"}
	dp = dotpath.New(m)
	if dp.GetString("key") != "value" {
		t.Error("Expected existing map to be preserved")
	}
}

func TestGet(t *testing.T) {
	m := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "demo",
			"port": 8080,
		},
		"debug": true,
	}
	dp := dotpath.New(m)

	// Test existing paths
	if val, exists := dp.Get("app.name"); !exists || val != "demo" {
		t.Errorf("Expected 'demo', got %v, exists: %v", val, exists)
	}

	if val, exists := dp.Get("app.port"); !exists || val != 8080 {
		t.Errorf("Expected 8080, got %v, exists: %v", val, exists)
	}

	if val, exists := dp.Get("debug"); !exists || val != true {
		t.Errorf("Expected true, got %v, exists: %v", val, exists)
	}

	// Test non-existing paths
	if _, exists := dp.Get("app.nonexistent"); exists {
		t.Error("Expected non-existent path to return false")
	}

	if _, exists := dp.Get("nonexistent"); exists {
		t.Error("Expected non-existent path to return false")
	}

	// Test empty path returns entire map
	if val, exists := dp.Get(""); !exists || val == nil {
		t.Error("Expected empty path to return entire map")
	}
}

func TestGetString(t *testing.T) {
	m := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "demo",
			"port": 8080, // int type
		},
		"empty": "",
	}
	dp := dotpath.New(m)

	// Test existing string
	if name := dp.GetString("app.name"); name != "demo" {
		t.Errorf("Expected 'demo', got '%s'", name)
	}

	// Test non-string type returns empty
	if port := dp.GetString("app.port"); port != "" {
		t.Errorf("Expected empty string for int type, got '%s'", port)
	}

	// Test empty string
	if empty := dp.GetString("empty"); empty != "" {
		t.Errorf("Expected empty string, got '%s'", empty)
	}

	// Test non-existing path
	if missing := dp.GetString("missing"); missing != "" {
		t.Errorf("Expected empty string for missing path, got '%s'", missing)
	}
}

func TestGetInt(t *testing.T) {
	m := map[string]interface{}{
		"port":         8080,
		"floatPort":    8080.0,
		"invalidFloat": 808.5,
		"string":       "not a number",
	}
	dp := dotpath.New(m)

	// Test existing int
	if port := dp.GetInt("port"); port != 8080 {
		t.Errorf("Expected 8080, got %d", port)
	}

	// Test float64 that is a whole number
	if port := dp.GetInt("floatPort"); port != 8080 {
		t.Errorf("Expected 8080, got %d", port)
	}

	// Test float64 that is not a whole number
	if port := dp.GetInt("invalidFloat"); port != 0 {
		t.Errorf("Expected 0 for non-whole float, got %d", port)
	}

	// Test non-numeric type
	if num := dp.GetInt("string"); num != 0 {
		t.Errorf("Expected 0 for string type, got %d", num)
	}

	// Test non-existing path
	if missing := dp.GetInt("missing"); missing != 0 {
		t.Errorf("Expected 0 for missing path, got %d", missing)
	}
}

func TestGetBool(t *testing.T) {
	m := map[string]interface{}{
		"debug":    true,
		"disabled": false,
		"number":   1,
		"string":   "true",
	}
	dp := dotpath.New(m)

	// Test existing bool
	if debug := dp.GetBool("debug"); !debug {
		t.Errorf("Expected true, got %v", debug)
	}

	if disabled := dp.GetBool("disabled"); disabled {
		t.Errorf("Expected false, got %v", disabled)
	}

	// Test non-bool type
	if num := dp.GetBool("number"); num {
		t.Errorf("Expected false for number type, got %v", num)
	}

	// Test non-existing path
	if missing := dp.GetBool("missing"); missing {
		t.Errorf("Expected false for missing path, got %v", missing)
	}
}

func TestGetFloat64(t *testing.T) {
	m := map[string]interface{}{
		"pi":      3.14159,
		"whole":   42.0,
		"integer": 42,
		"string":  "not a number",
	}
	dp := dotpath.New(m)

	// Test existing float64
	if pi := dp.GetFloat64("pi"); pi != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", pi)
	}

	// Test whole number float64
	if whole := dp.GetFloat64("whole"); whole != 42.0 {
		t.Errorf("Expected 42.0, got %f", whole)
	}

	// Test int type (should return 0)
	if num := dp.GetFloat64("integer"); num != 0 {
		t.Errorf("Expected 0 for int type, got %f", num)
	}

	// Test non-numeric type
	if num := dp.GetFloat64("string"); num != 0 {
		t.Errorf("Expected 0 for string type, got %f", num)
	}

	// Test non-existing path
	if missing := dp.GetFloat64("missing"); missing != 0 {
		t.Errorf("Expected 0 for missing path, got %f", missing)
	}
}

func TestHas(t *testing.T) {
	m := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "demo",
		},
	}
	dp := dotpath.New(m)

	// Test existing paths
	if !dp.Has("app.name") {
		t.Error("Expected Has to return true for existing path")
	}

	if !dp.Has("app") {
		t.Error("Expected Has to return true for existing top-level key")
	}

	// Test non-existing paths
	if dp.Has("app.nonexistent") {
		t.Error("Expected Has to return false for non-existent path")
	}

	if dp.Has("nonexistent") {
		t.Error("Expected Has to return false for non-existent top-level key")
	}
}

func TestSet(t *testing.T) {
	dp := dotpath.New(nil)

	// Test setting simple value
	dp.Set("name", "demo")
	if dp.GetString("name") != "demo" {
		t.Errorf("Expected 'demo', got %v", dp.GetString("name"))
	}

	// Test setting nested value (auto-creation)
	dp.Set("server.port", 8080)
	if !dp.Has("server") {
		t.Error("Expected 'server' key to be created")
	}

	if dp.GetInt("server.port") != 8080 {
		t.Errorf("Expected 8080, got %v", dp.GetInt("server.port"))
	}

	// Test setting deeply nested value
	dp.Set("a.b.c.d", "deep")
	if val := dp.GetString("a.b.c.d"); val != "deep" {
		t.Errorf("Expected 'deep', got '%s'", val)
	}

	// Test overwriting existing value
	dp.Set("name", "updated")
	if dp.GetString("name") != "updated" {
		t.Error("Expected value to be overwritten")
	}

	// Test setting with empty path (should do nothing)
	originalData := dp.Data()
	dp.Set("", "should not be set")
	if len(dp.Data()) != len(originalData) {
		t.Error("Empty path should not modify data")
	}
}

func TestSetWithConflict(t *testing.T) {
	dp := dotpath.New(nil)

	// Set a non-map value
	dp.Set("config", "value")

	// Try to set nested value under the non-map
	dp.Set("config.nested", "should work")

	// The original value should be replaced with a map
	if dp.GetString("config.nested") != "should work" {
		t.Error("Should be able to set nested value under non-map")
	}
}

func TestUpdate(t *testing.T) {
	dp := dotpath.New(map[string]interface{}{
		"app": map[string]interface{}{
			"name": "demo",
			"port": 3000,
		},
		"debug": true,
	})

	other := map[string]interface{}{
		"app": map[string]interface{}{
			"port":    8080, // Should overwrite existing
			"version": "1.0.0", // Should be added
		},
		"newKey": "newValue",
	}

	dp.Update(other)

	// Check merged values
	if dp.GetString("app.name") != "demo" {
		t.Error("Original nested value should be preserved")
	}

	if dp.GetInt("app.port") != 8080 {
		t.Error("Should overwrite existing nested value")
	}

	if dp.GetString("app.version") != "1.0.0" {
		t.Error("Should add new nested value")
	}

	if dp.GetString("newKey") != "newValue" {
		t.Error("Should add new top-level value")
	}

	if !dp.GetBool("debug") {
		t.Error("Original top-level value should be preserved")
	}
}

func TestUpdateWithNil(t *testing.T) {
	dp := dotpath.New(map[string]interface{}{"key": "value"})
	originalSize := len(dp.Data())

	dp.Update(nil)
	if len(dp.Data()) != originalSize {
		t.Error("Update with nil should not modify data")
	}
}

func TestData(t *testing.T) {
	original := map[string]interface{}{
		"key": "value",
		"nested": map[string]interface{}{
			"inner": true,
		},
	}
	dp := dotpath.New(original)

	data := dp.Data()
	if data == nil {
		t.Error("Data should return the original map")
	}

	// Verify it's the same map (by modifying it)
	data["newKey"] = "newValue"
	if dp.GetString("newKey") != "newValue" {
		t.Error("Data should return the actual underlying map")
	}
}

func TestDeepMergeViaUpdate(t *testing.T) {
	// Test deep merge functionality via Update method
	dp := dotpath.New(map[string]interface{}{
		"a": 1,
		"b": 2,
		"nested": map[string]interface{}{
			"x": 10,
			"y": 20,
		},
	})

	updates := map[string]interface{}{
		"b": 3,           // Should overwrite
		"c": 4,           // Should be added
		"nested": map[string]interface{}{
			"y": 30,       // Should overwrite nested
			"z": 40,       // Should be added nested
		},
	}

	dp.Update(updates)
	data := dp.Data()

	// Verify merge results
	if data["a"] != 1 {
		t.Error("Original value should be preserved")
	}
	if data["b"] != 3 {
		t.Error("Source value should overwrite destination")
	}
	if data["c"] != 4 {
		t.Error("New value should be added")
	}

	// Check nested merge
	nested := data["nested"].(map[string]interface{})
	if nested["x"] != 10 {
		t.Error("Original nested value should be preserved")
	}
	if nested["y"] != 30 {
		t.Error("Source nested value should overwrite destination")
	}
	if nested["z"] != 40 {
		t.Error("New nested value should be added")
	}
}

// Example test from the PRD
func TestExampleFromPRD(t *testing.T) {
	m := map[string]interface{}{
		"app": map[string]interface{}{
			"name": "demo",
		},
	}
	dp := dotpath.New(m)

	name := dp.GetString("app.name")
	if name != "demo" {
		t.Errorf("Expected 'demo', got '%s'", name)
	}

	dp.Set("server.port", 8080)

	if port := dp.GetInt("server.port"); port != 8080 {
		t.Errorf("Expected 8080, got %d", port)
	}
}