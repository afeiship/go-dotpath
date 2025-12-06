You are building adapter packages for the `go-dotpath` core library (which operates on `map[string]interface{}` via dot-path access).

Create two separate adapter packages:

---

### 1. Package: `jsonadapter`
- Path: `jsonadapter/`
- Purpose: Load/save dotpath-compatible data from/to JSON

**Requirements:**
- Depends on: `github.com/yourname/go-dotpath` (use placeholder import; actual path may vary)
- Uses only Go standard library (`encoding/json`)
- Provides:
  - `LoadFromBytes(data []byte) (*dotpath.DotPath, error)`
  - `LoadFromString(s string) (*dotpath.DotPath, error)`
  - `LoadFromFile(path string) (*dotpath.DotPath, error)`
  - `SaveToFile(dp *dotpath.DotPath, path string) error`
  - `ToBytes(dp *dotpath.DotPath) ([]byte, error)`
  - `ToString(dp *dotpath.DotPath) (string, error)`

> All functions must handle JSON ↔ `map[string]interface{}` conversion and wrap/unwrap with `dotpath.New()`.

---

### 2. Package: `yamladapter`
- Path: `yamladapter/`
- Purpose: Load/save dotpath-compatible data from/to YAML

**Requirements:**
- Depends on: 
  - `github.com/yourname/go-dotpath`
  - `gopkg.in/yaml.v3`
- Provides:
  - `LoadFromBytes(data []byte) (*dotpath.DotPath, error)`
  - `LoadFromString(s string) (*dotpath.DotPath, error)`
  - `LoadFromFile(path string) (*dotpath.DotPath, error)`
  - `SaveToFile(dp *dotpath.DotPath, path string) error`
  - `ToBytes(dp *dotpath.DotPath) ([]byte, error)`
  - `ToString(dp *dotpath.DotPath) (string, error)`

> Use `yaml.Unmarshal` and `yaml.Marshal` internally. Do not modify the core `dotpath` logic.

---

### Shared Rules for Both Adapters:
- Each adapter must be **self-contained** in its own directory.
- Do **not** duplicate core logic — only handle serialization/deserialization.
- Return descriptive errors (e.g., file not found, invalid JSON/YAML).
- Include a simple usage example in `examples/json_example.go` and `examples/yaml_example.go`.
- Place unit tests in `__tests__/jsonadapter/` and `__tests__/yamladapter/` respectively.
- Use consistent function signatures across both adapters for ergonomic parity.

---

### Example Usage (YAML):
```go
dp, err := yamladapter.LoadFromFile("config.yaml")
if err != nil { /* handle */ }
port := dp.GetInt("server.port")
dp.Set("server.timeout", 30)
yamladapter.SaveToFile(dp, "config.yaml")