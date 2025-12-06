# JSON Configuration Examples with go-dotpath

This document provides practical examples of using go-dotpath to manage JSON configuration files in Go applications.

## 1. Basic Configuration File

### config.json
```json
{
  "app": {
    "name": "my-service",
    "version": "1.0.0",
    "debug": true,
    "author": "development-team"
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "cors": {
      "enabled": true,
      "origins": ["http://localhost:3000", "http://localhost:8080"],
      "credentials": true
    }
  },
  "database": {
    "driver": "postgres",
    "host": "localhost",
    "port": 5432,
    "database": "myapp",
    "sslmode": "disable",
    "pool": {
      "max": 100,
      "min": 10,
      "timeout": 30
    }
  },
  "logging": {
    "level": "info",
    "format": "json",
    "outputs": ["stdout", "file"],
    "file": {
      "path": "/var/log/app.log",
      "maxSize": "100MB",
      "rotate": true
    }
  }
}
```

### Go Implementation
```go
package main

import (
    "fmt"
    "log"
    "github.com/afeiship/go-dotpath"
)

func main() {
    // Load configuration
    io := dotpath.NewDotIO(dotpath.JSON)
    if err := io.LoadFromFile("config.json"); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Access configuration values
    appName := io.GetString("app.name")
    debug := io.GetBool("app.debug")
    dbHost := io.GetString("database.host")
    dbPort := io.GetInt("database.port")
    maxPool := io.GetInt("database.pool.max")

    fmt.Printf("Application: %s (Debug: %t)\n", appName, debug)
    fmt.Printf("Database: %s:%d (Pool max: %d)\n", dbHost, dbPort, maxPool)

    // Access nested configuration
    corsOrigins := io.Data()["server"].(map[string]any)["cors"].(map[string]any)["origins"].([]any)
    fmt.Printf("CORS Origins: %v\n", corsOrigins)

    // Update configuration
    io.Set("app.version", "2.0.0")
    io.Set("database.pool.max", 200)
    io.Set("monitoring.enabled", true)

    // Add new nested configuration
    io.Set("api.rateLimit", map[string]any{
        "requests": 1000,
        "window": "1m",
        "burst": 5000,
    })

    // Save updated configuration
    if err := io.SaveToFile("config.json"); err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }

    fmt.Println("Configuration updated successfully!")
}
```

## 2. Environment-Specific Configurations

### config.base.json
```json
{
  "app": {
    "name": "my-service"
  },
  "server": {
    "port": 8080
  },
  "database": {
    "driver": "sqlite"
  },
  "logging": {
    "level": "info"
  }
}
```

### config.dev.json
```json
{
  "app": {
    "environment": "development"
  },
  "server": {
    "debug": true,
    "timeout": 30
  },
  "database": {
    "host": "localhost"
  },
  "logging": {
    "level": "debug"
  }
}
```

### config.prod.json
```json
{
  "app": {
    "environment": "production"
  },
  "server": {
    "debug": false,
    "timeout": 10
  },
  "database": {
    "driver": "postgres",
    "host": "prod-db.example.com",
    "sslmode": "require"
  },
  "logging": {
    "level": "warn"
  }
}
```

### Go Implementation
```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/afeiship/go-dotpath"
)

func loadConfiguration(environment string) (*dotpath.DotIO, error) {
    // Load base configuration
    baseIO := dotpath.NewDotIO(dotpath.JSON)
    if err := baseIO.LoadFromFile("config.base.json"); err != nil {
        return nil, fmt.Errorf("failed to load base config: %w", err)
    }

    // Load environment-specific override
    envFile := fmt.Sprintf("config.%s.json", environment)
    envIO := dotpath.NewDotIO(dotpath.JSON)

    if _, err := os.Stat(envFile); err == nil {
        if err := envIO.LoadFromFile(envFile); err != nil {
            return nil, fmt.Errorf("failed to load %s: %w", envFile, err)
        }
        // Merge configurations
        baseIO.Update(envIO.Data())
    }

    return baseIO, nil
}

func main() {
    env := os.Getenv("GO_ENV")
    if env == "" {
        env = "development" // Default to development
    }

    fmt.Printf("Loading configuration for environment: %s\n", env)

    config, err := loadConfiguration(env)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    fmt.Printf("Environment: %s\n", config.GetString("app.environment"))
    fmt.Printf("Debug Mode: %t\n", config.GetBool("server.debug"))

    // Save final merged configuration
    if err := config.SaveToFile("config.json"); err != nil {
        log.Fatalf("Failed to save merged config: %v", err)
    }

    fmt.Println("Configuration loaded and saved successfully!")
}
```

## 3. Configuration Validation

### Go Implementation
```go
package main

import (
    "fmt"
    "log"
    "github.com/afeiship/go-dotpath"
)

type ConfigValidator struct {
    required []string
    checks  map[string]func(*dotpath.DotIO) error
}

func (cv *ConfigValidator) AddRequired(path string) *ConfigValidator {
    cv.required = append(cv.required, path)
    return cv
}

func (cv *ConfigValidator) AddCheck(path string, check func(*dotpath.DotIO) error) *ConfigValidator {
    if cv.checks == nil {
        cv.checks = make(map[string]func(*dotpath.DotIO) error)
    }
    cv.checks[path] = check
    return cv
}

func (cv *ConfigValidator) Validate(io *dotpath.DotIO) error {
    // Check required fields
    for _, path := range cv.required {
        if !io.Has(path) {
            return fmt.Errorf("missing required configuration: %s", path)
        }
    }

    // Run custom checks
    for path, check := range cv.checks {
        if err := check(io); err != nil {
            return fmt.Errorf("validation failed for %s: %w", path, err)
        }
    }

    return nil
}

func main() {
    io := dotpath.NewDotIO(dotpath.JSON)
    if err := io.LoadFromFile("config.json"); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Setup validator
    validator := &ConfigValidator{}
    validator.
        AddRequired("app.name").
        AddRequired("server.port").
        AddRequired("database.host").
        AddCheck("server.port", func(io *dotpath.DotIO) error {
            port := io.GetInt("server.port")
            if port <= 0 || port > 65535 {
                return fmt.Errorf("invalid port number: %d", port)
            }
            return nil
        }).
        AddCheck("database.driver", func(io *dotpath.DotIO) error {
            driver := io.GetString("database.driver")
            validDrivers := []string{"postgres", "mysql", "sqlite"}
            for _, valid := range validDrivers {
                if driver == valid {
                    return nil
                }
            }
            return fmt.Errorf("invalid database driver: %s", driver)
        })

    // Validate configuration
    if err := validator.Validate(io); err != nil {
        log.Fatalf("Configuration validation failed: %v", err)
    }

    fmt.Println("Configuration is valid!")
    fmt.Printf("App: %s\n", io.GetString("app.name"))
    fmt.Printf("Server Port: %d\n", io.GetInt("server.port"))
}
```

## 4. Hot Configuration Reloading

### Go Implementation
```go
package main

import (
    "log"
    "os/signal"
    "syscall"
    "time"
    "github.com/afeiship/go-dotpath"
)

type ConfigManager struct {
    config *dotpath.DotIO
    file   string
    mtimes time.Time
}

func NewConfigManager(configFile string) *ConfigManager {
    cm := &ConfigManager{
        file: configFile,
    }
    cm.Load()
    cm.WatchForChanges()
    return cm
}

func (cm *ConfigManager) Load() error {
    io := dotpath.NewDotIO(dotpath.JSON)
    if err := io.LoadFromFile(cm.file); err != nil {
        return err
    }

    // Get file modification time
    if info, err := os.Stat(cm.file); err == nil {
        cm.mtimes = info.ModTime()
    }

    cm.config = io
    return nil
}

func (cm *ConfigManager) Get() *dotpath.DotIO {
    return cm.config
}

func (cm *ConfigManager) WatchForChanges() {
    go func() {
        for {
            time.Sleep(5 * time.Second)

            if info, err := os.Stat(cm.file); err == nil {
                if info.ModTime().After(cm.mtimes) {
                    cm.mtimes = info.ModTime()
                    if err := cm.Load(); err != nil {
                        log.Printf("Error reloading config: %v", err)
                    } else {
                        log.Println("Configuration reloaded successfully")
                    }
                }
            }
        }
    }()
}

func (cm *ConfigManager) Save() error {
    return cm.config.SaveToFile(cm.file)
}

func main() {
    cm := NewConfigManager("config.json")

    // Set up signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        for sig := range sigChan {
            log.Printf("Received signal: %v. Shutting down...\n", sig)
            break
        }
    }()

    // Main application loop
    for {
        config := cm.Get()

        // Use configuration values
        appName := config.GetString("app.name")
        debug := config.GetBool("app.debug")
        port := config.GetInt("server.port")

        fmt.Printf("App: %s, Port: %d, Debug: %t\n", appName, port, debug)

        // Simulate work
        time.Sleep(10 * time.Second)
    }
}
```

## 5. API Configuration

### api-config.json
```json
{
  "endpoints": {
    "users": {
      "path": "/api/v1/users",
      "method": "GET",
      "auth": true
    },
    "posts": {
      "path": "/api/v1/posts",
      "method": ["GET", "POST"],
      "auth": true
    },
    "health": {
      "path": "/health",
      "method": "GET",
      "auth": false
    }
  },
  "rateLimiting": {
    "global": {
      "requests": 1000,
      "window": "1m",
      "burst": 5000
    },
    "endpoints": {
      "/api/v1/users": {
        "requests": 100,
        "window": "1m"
      }
    }
  },
  "cors": {
    "allowedOrigins": ["http://localhost:3000", "https://example.com"],
    "allowedMethods": ["GET", "POST", "PUT", "DELETE"],
    "allowedHeaders": ["Content-Type", "Authorization"]
  }
}
```

### Go Implementation
```go
package api

import (
    "fmt"
    "net/http"
    "github.com/afeiship/go-dotpath"
)

type APIConfig struct {
    io *dotpath.DotIO
}

func LoadAPIConfig(configFile string) (*APIConfig, error) {
    io := dotpath.NewDotIO(dotpath.JSON)
    if err := io.LoadFromFile(configFile); err != nil {
        return nil, err
    }
    return &APIConfig{io: io}, nil
}

func (ac *APIConfig) GetEndpointConfig(endpoint string) (map[string]any, error) {
    if !ac.io.Has("endpoints." + endpoint) {
        return nil, fmt.Errorf("endpoint not found: %s", endpoint)
    }
    return ac.io.Get("endpoints." + endpoint).(map[string]any), nil
}

func (ac *APIConfig) IsAuthenticated(endpoint string) bool {
    config, err := ac.GetEndpointConfig(endpoint)
    if err != nil {
        return false
    }
    return config["auth"].(bool)
}

func (ac *APIConfig) GetRateLimit(endpoint string) (int, string) {
    // Check endpoint-specific rate limit first
    if ac.io.Has("rateLimiting.endpoints." + endpoint) {
        epConfig := ac.io.Get("rateLimit.endpoints." + endpoint).(map[string]any)
        return epConfig["requests"].(int), epConfig["window"].(string)
    }

    // Fall back to global rate limit
    global := ac.io.Get("rateLimiting.global").(map[string]any)
    return global["requests"].(int), global["window"].(string)
}

func (ac *APIConfig) SetupRoutes(mux *http.ServeMux) {
    endpoints := ac.io.Data()["endpoints"].(map[string]any)
    for name, config := range endpoints {
        epConfig := config.(map[string]any)
        path := epConfig["path"].(string)

        handler := func(w http.ResponseWriter, r *http.Request) {
            // Apply rate limiting
            requests, window := ac.GetRateLimit(name)
            // Rate limiting implementation would go here

            // Check authentication
            if ac.IsAuthenticated(name) {
                // Check Authorization header
                authHeader := r.Header.Get("Authorization")
                if authHeader == "" {
                    http.Error(w, "Unauthorized", http.StatusUnauthorized)
                    return
                }
            }

            // Handle request
            fmt.Fprintf(w, "Endpoint: %s, Method: %s", name, r.Method)
        }

        // Handle multiple HTTP methods if specified
        methods := epConfig["method"]
        if methodSlice, ok := methods.([]any); ok {
            for _, method := range methodSlice {
                mux.HandleFunc(path, handler)
                if method != "GET" {
                    mux.HandleFunc(path, handler)
                }
            }
        } else {
            // Single method (GET by default)
            mux.HandleFunc(path, handler)
        }
    }
}
```

## 6. Testing with Mock Configuration

### config.test.json
```json
{
  "test_mode": true,
  "database": {
    "mock": true,
    "data": {
      "users": [
        {"id": 1, "name": "Test User"},
        {"id": 2, "name": "Another User"}
      ]
    }
  }
}
```

### Go Implementation
```go
package test

import (
    "github.com/afeiship/go-dotpath"
)

func SetupTestConfig() *dotpath.DotIO {
    io := dotpath.NewDotIO(dotpath.JSON)

    // Load test configuration if exists, create default otherwise
    if _, err := os.Stat("config.test.json"); os.IsNotExist(err) {
        // Create default test configuration
        testConfig := map[string]any{
            "test_mode": true,
            "server": map[string]any{
                "port": 3001,
                "host": "localhost",
            },
            "database": map[string]any{
                "mock": true,
                "data": map[string]any{
                    "users": []any{
                        map[string]any{"id": 1, "name": "Test User"},
                        map[string]any{"id": 2, "name": "Another User"},
                    },
                },
            },
        }

        io = dotpath.NewDotIOWithData(testConfig, dotpath.JSON)
    } else {
        io.LoadFromFile("config.test.json")
    }

    return io.DotPath()
}
```

## Best Practices

### 1. Use Environment Variables
```go
func GetConfigPath() string {
    if path := os.Getenv("CONFIG_PATH"); path != "" {
        return path
    }
    return "config.json"
}
```

### 2. Provide Default Values
```go
func GetDefaultConfig() map[string]any {
    return map[string]any{
        "app": map[string]any{
            "name": "default-app",
            "version": "1.0.0",
            "debug": false,
        },
        "server": map[string]any{
            "port": 8080,
            "host": "0.0.0.0",
        },
    }
}
```

### 3. Use Configuration Schema Validation
```go
type ConfigSchema struct {
    App    AppConfig    `json:"app"`
    Server ServerConfig `json:"server"`
    DB     DatabaseConfig `json:"database"`
}

type AppConfig struct {
    Name    string `json:"name" validate:"required"`
    Version string `json:"version" validate:"required"`
    Debug   bool   `json:"debug"`
}

func (c *ConfigSchema) Validate() error {
    return nil
}
```

These examples demonstrate practical ways to use go-dotpath for JSON configuration management in real applications.