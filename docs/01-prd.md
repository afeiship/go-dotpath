
Please implement a Go package named `go-dotpath` that provides dot-notation access and mutation for nested configuration data represented as `map[string]interface{}`.

**Requirements:**

1. **Core Type**: Define a type (e.g., `DotPath`) that wraps `map[string]interface{}`.
2. **Constructor**: Provide a `New(m map[string]interface{}) *DotPath` function.
3. **Getter Methods** (safe, zero-value on missing/type mismatch):
   - `Get(path string) (interface{}, bool)`
   - `GetString(path string) string`
   - `GetInt(path string) int`
   - `GetBool(path string) bool`
   - `GetFloat64(path string) float64`
   - `Has(path string) bool`  // alias for key existence
4. **Setter & Mutators**:
   - `Set(path string, value interface{})` — must auto-create missing nested maps (e.g., `"a.b.c"` creates `a → b → c`)
   - `Update(other map[string]interface{})` — perform deep recursive merge
5. **Data Access**:
   - `Data() map[string]interface{}` — returns the underlying map (for serialization, inspection, etc.)
6. **Path Syntax**: Use dot notation (`"server.port"`, `"x.y.z"`). Do not support arrays or advanced paths — only string-keyed nested maps.
7. **Dependencies**: The package must **only use the Go standard library**. No external imports (no YAML, JSON, etc.).
8. **Error Handling**: Never panic. Return zero values or `false` for missing/invalid paths.
9. **Code Quality**: Idiomatic Go, clear comments, and efficient navigation.

**Example Usage:**
```go
m := map[string]interface{}{
    "app": map[string]interface{}{
        "name": "demo",
    },
}
dp := dotpath.New(m)
name := dp.GetString("app.name")   // "demo"
dp.Set("server.port", 8080)        // creates nested map if needed
```


## Deliverables:
- A single file dotpath.go with the full implementation
- A separate test file in the __tests__ directory: __tests__/dotpath_test.go
- An example file in examples/main.go showing basic usage
- No file I/O, encoding, or format-specific logic (keep it generic)