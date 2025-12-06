# YAML Configuration Examples with go-dotpath

This document provides practical examples of using go-dotpath to manage YAML configuration files in Go applications.

## 1. Basic Application Configuration

### config.yaml
```yaml
app:
  name: "my-service"
  version: "1.0.0"
  environment: "production"
  debug: false

server:
  host: "0.0.0.0"
  port: 8080
  timeout: 30s
  workers: 4
  cors:
    enabled: true
    origins:
      - "http://localhost:3000"
      - "http://localhost:8080"
    credentials: true

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  database: "myapp"
  sslmode: "require"
  pool:
    max: 100
    min: 10
    timeout: 30s

logging:
  level: "info"
  format: "json"
  outputs:
    - "stdout"
    - "file"
  file:
    path: "/var/log/app.log"
    max_size: "100MB"
    rotate: true

monitoring:
  enabled: true
  metrics:
    enabled: true
    path: "/metrics"
  health_check:
    enabled: true
    path: "/health"
    timeout: 10s
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
    io := dotpath.NewDotIO(dotpath.YAML)
    if err := io.LoadFromFile("config.yaml"); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Access configuration values
    appName := io.GetString("app.name")
    env := io.GetString("app.environment")
    debug := io.GetBool("app.debug")
    port := io.GetInt("server.port")
    timeout := io.GetString("server.timeout")
    poolMax := io.GetInt("database.pool.max")

    fmt.Printf("Application: %s (Env: %s, Debug: %t)\n", appName, env, debug)
    fmt.Printf("Server: Port %d, Timeout %s\n", port, timeout)
    fmt.Printf("Database Pool Max: %d\n", poolMax)

    // Access nested configuration
    origins := io.Data()["server"].(map[string]any)["cors"].(map[string]any)["origins"].([]any)
    fmt.Printf("CORS Origins: %v\n", origins)

    // Update configuration
    io.Set("app.version", "2.0.0")
    io.Set("monitoring.tracing.enabled", true)
    io.Set("database.pool.max", 200)

    // Save updated configuration with comment preservation
    if err := io.Save(); err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }

    fmt.Println("Configuration updated successfully!")
}
```

## 2. Multi-Environment Configuration

### config.base.yaml
```yaml
app:
  name: "my-service"
  version: "1.0.0"

server:
  port: 8080
  host: "0.0.0.0"

database:
  driver: "sqlite"

logging:
  level: "info"
  format: "text"

monitoring:
  enabled: false
```

### config.dev.yaml
```yaml
app:
  environment: "development"
  debug: true

server:
  debug: true
  timeout: 30s
  read_timeout: 60s
  write_timeout: 60s

database:
  host: "localhost"
  port: 5432

logging:
  level: "debug"
  format: "console"

monitoring:
  enabled: true
  debug: true
```

### config.prod.yaml
```yaml
app:
  environment: "production"
  debug: false

server:
  debug: false
  timeout: 10s
  read_timeout: 15s
  write_timeout: 15s
  tls:
    enabled: true
    cert_file: "/etc/ssl/certs/server.crt"
    key_file: "/etc/ssl/private/server.key"

database:
  driver: "postgres"
  host: "prod-db.example.com"
  sslmode: "require"
  backup:
    enabled: true
    schedule: "0 2 * * *"
    retention_days: 30

logging:
  level: "warn"
  format: "json"
  outputs:
    - "file"
    - "syslog"

monitoring:
  enabled: true
  alerts:
    enabled: true
    webhook_url: "https://hooks.slack.com/..."
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
    baseIO := dotpath.NewDotIO(dotpath.YAML)
    if err := baseIO.LoadFromFile("config.base.yaml"); err != nil {
        return nil, fmt.Errorf("failed to load base config: %w", err)
    }

    // Load environment-specific override
    envFile := fmt.Sprintf("config.%s.yaml", environment)
    envIO := dotpath.NewDotIO(dotpath.YAML)

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
        env = "development"
    }

    fmt.Printf("Loading configuration for environment: %s\n", env)

    config, err := loadConfiguration(env)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    fmt.Printf("Environment: %s\n", config.GetString("app.environment"))
    fmt.Printf("Debug Mode: %t\n", config.GetBool("app.debug"))
    fmt.Printf("Database Driver: %s\n", config.GetString("database.driver"))

    // Save final merged configuration
    if err := config.SaveToFile("config.yaml"); err != nil {
        log.Fatalf("Failed to save merged config: %v", err)
    }

    fmt.Println("Configuration loaded and saved successfully!")
}
```

## 3. Service Discovery Configuration

### services.yaml
```yaml
services:
  auth-service:
    host: "auth-service.example.com"
    port: 8081
    protocol: "https"
    health_check: "/health"
    timeout: 5s
    retries: 3
    circuit_breaker:
      enabled: true
      failure_threshold: 5
      recovery_timeout: 30s

  user-service:
    host: "user-service.example.com"
    port: 8082
    protocol: "https"
    health_check: "/health"
    timeout: 10s
    retries: 3
    circuit_breaker:
      enabled: true
      failure_threshold: 10
      recovery_timeout: 60s

  database:
    primary:
      host: "db-primary.example.com"
      port: 5432
      connection_pool:
        max: 100
        min: 10
        timeout: 30s
    replica:
      host: "db-replica.example.com"
      port: 5432
      connection_pool:
        max: 50
        min: 5
        timeout: 30s

  cache:
    primary:
      host: "cache-primary.example.com"
      port: 6379
      db: 0
      timeout: 2s
    replica:
      host: "cache-replica.example.com"
      port: 6379
      db: 0
      timeout: 2s

  message_queue:
    type: "rabbitmq"
    host: "mq.example.com"
    port: 5672
    vhost: "/"
    exchanges:
      - name: "events"
        type: "topic"
        durable: true
      - name: "commands"
        type: "direct"
        durable: true
```

### Go Implementation
```go
package main

import (
    "fmt"
    "time"
    "github.com/afeiship/go-dotpath"
)

type ServiceConfig struct {
    Host     string
    Port     int
    Protocol string
    Timeout  time.Duration
    Retries  int
}

func main() {
    io := dotpath.NewDotIO(dotpath.YAML)
    if err := io.LoadFromFile("services.yaml"); err != nil {
        log.Fatalf("Failed to load services config: %v", err)
    }

    // Iterate through all services
    services := io.Data()["services"].(map[string]any)

    for serviceName, serviceData := range services {
        if serviceName == "database" || serviceName == "cache" || serviceName == "message_queue" {
            continue // Skip complex nested services for this example
        }

        config := serviceData.(map[string]any)
        serviceConfig := ServiceConfig{
            Host:     config["host"].(string),
            Port:     int(config["port"].(int)),
            Protocol: config["protocol"].(string),
            Retries:  int(config["retries"].(int)),
        }

        // Parse timeout
        if timeoutStr, ok := config["timeout"].(string); ok {
            if timeout, err := time.ParseDuration(timeoutStr); err == nil {
                serviceConfig.Timeout = timeout
            }
        }

        fmt.Printf("Service %s: %s://%s:%d (Timeout: %v, Retries: %d)\n",
            serviceName, serviceConfig.Protocol, serviceConfig.Host,
            serviceConfig.Port, serviceConfig.Timeout, serviceConfig.Retries)

        // Check circuit breaker configuration
        cbPath := fmt.Sprintf("services.%s.circuit_breaker", serviceName)
        if io.Has(cbPath) {
            cbEnabled := io.GetBool(cbPath + ".enabled")
            cbThreshold := io.GetInt(cbPath + ".failure_threshold")
            cbRecovery := io.GetString(cbPath + ".recovery_timeout")
            fmt.Printf("  Circuit Breaker: enabled=%t, threshold=%d, recovery=%s\n",
                cbEnabled, cbThreshold, cbRecovery)
        }
    }

    // Handle database configuration
    dbPrimaryPath := "services.database.primary"
    dbHost := io.GetString(dbPrimaryPath + ".host")
    dbPort := io.GetInt(dbPrimaryPath + ".port")
    dbPoolMax := io.GetInt(dbPrimaryPath + ".connection_pool.max")

    fmt.Printf("Database Primary: %s:%d (Pool max: %d)\n", dbHost, dbPort, dbPoolMax)
}
```

## 4. Feature Flag Configuration

### features.yaml
```yaml
features:
  new_ui:
    enabled: true
    rollout_percentage: 50
    user_groups:
      - "beta_testers"
      - "internal_team"
    conditions:
      min_version: "2.0.0"
      max_version: ""

  analytics:
    enabled: true
    providers:
      - name: "google_analytics"
        enabled: true
        config:
          tracking_id: "GA-XXXXXXX"
      - name: "mixpanel"
        enabled: false

  maintenance_mode:
    enabled: false
    message: "System under maintenance. Please try again later."
    override_users:
      - "admin"
      - "support"

  rate_limiting:
    enabled: true
    default_limits:
      requests: 1000
      window: "1m"
      burst: 5000
    endpoint_limits:
      "/api/v1/users":
        requests: 100
        window: "1m"
      "/api/v1/auth":
        requests: 10
        window: "5m"

  a_b_testing:
    enabled: true
    experiments:
      - name: "new_checkout_flow"
        enabled: true
        traffic_percentage: 25
        variations:
          - name: "control"
            percentage: 50
          - name: "variant_a"
            percentage: 50
      - name: "pricing_display"
        enabled: true
        traffic_percentage: 50
        variations:
          - name: "monthly"
            percentage: 33
          - name: "yearly"
            percentage: 33
          - name: "lifetime"
            percentage: 34
```

### Go Implementation
```go
package main

import (
    "fmt"
    "math/rand"
    "github.com/afeiship/go-dotpath"
)

type FeatureManager struct {
    io *dotpath.DotIO
}

func NewFeatureManager(configFile string) (*FeatureManager, error) {
    io := dotpath.NewDotIO(dotpath.YAML)
    if err := io.LoadFromFile(configFile); err != nil {
        return nil, err
    }
    return &FeatureManager{io: io}, nil
}

func (fm *FeatureManager) IsFeatureEnabled(featureName string, userID string) bool {
    featurePath := fmt.Sprintf("features.%s", featureName)

    if !fm.io.Has(featurePath) {
        return false
    }

    if !fm.io.GetBool(featurePath + ".enabled") {
        return false
    }

    // Check user overrides
    overrideUsers := fm.io.GetStringSlice(featurePath + ".override_users")
    for _, user := range overrideUsers {
        if user == userID {
            return true
        }
    }

    // Check rollout percentage
    if fm.io.Has(featurePath + ".rollout_percentage") {
        percentage := fm.GetInt(featurePath + ".rollout_percentage")
        // Simple hash-based rollout (in production, use consistent hashing)
        hash := fm.hashUserID(userID)
        return (hash % 100) < percentage
    }

    return true
}

func (fm *FeatureManager) hashUserID(userID string) int {
    hash := 0
    for _, char := range userID {
        hash = hash*31 + int(char)
    }
    return hash
}

func (fm *FeatureManager) GetABTestVariant(experimentName string, userID string) string {
    expPath := fmt.Sprintf("features.a_b_testing.experiments")

    experiments := fm.io.Data()["features"].(map[string]any)["a_b_testing"].(map[string]any)["experiments"].([]any)

    for _, exp := range experiments {
        expData := exp.(map[string]any)
        if expData["name"].(string) == experimentName && expData["enabled"].(bool) {
            trafficPercentage := int(expData["traffic_percentage"].(int))
            if (fm.hashUserID(userID) % 100) >= trafficPercentage {
                return "" // User not in experiment
            }

            variations := expData["variations"].([]any)
            hash := fm.hashUserID(userID)
            percentage := hash % 100

            for _, varData := range variations {
                varInfo := varData.(map[string]any)
                varPercentage := int(varInfo["percentage"].(int))
                if percentage < varPercentage {
                    return varInfo["name"].(string)
                }
                percentage -= varPercentage
            }
        }
    }

    return ""
}

func (fm *FeatureManager) UpdateFeature(featureName string, updates map[string]any) error {
    featurePath := fmt.Sprintf("features.%s", featureName)
    for key, value := range updates {
        fm.io.Set(featurePath+"."+key, value)
    }
    return fm.io.SaveToFile("features.yaml")
}

func main() {
    fm, err := NewFeatureManager("features.yaml")
    if err != nil {
        log.Fatalf("Failed to load feature manager: %v", err)
    }

    // Check feature flags for a user
    userID := "user123"
    fmt.Printf("Feature status for user: %s\n", userID)
    fmt.Printf("New UI: %t\n", fm.IsFeatureEnabled("new_ui", userID))
    fmt.Printf("Analytics: %t\n", fm.IsFeatureEnabled("analytics", userID))
    fmt.Printf("Rate Limiting: %t\n", fm.IsFeatureEnabled("rate_limiting", userID))

    // Get A/B test variant
    variant := fm.GetABTestVariant("new_checkout_flow", userID)
    if variant != "" {
        fmt.Printf("A/B Test Variant for '%s': %s\n", "new_checkout_flow", variant)
    }

    // Update feature flag
    err = fm.UpdateFeature("maintenance_mode", map[string]any{
        "enabled": true,
        "message": "Scheduled maintenance starting at 2 AM",
    })
    if err != nil {
        fmt.Printf("Failed to update feature: %v\n", err)
    } else {
        fmt.Println("Maintenance mode enabled successfully!")
    }
}
```

## 5. Kubernetes Configuration

### k8s-config.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: production
data:
  app.yaml: |
    app:
      name: "my-service"
      version: "1.0.0"
      environment: "production"

    server:
      port: 8080
      host: "0.0.0.0"

    database:
      driver: "postgres"
      host: "postgres-service"
      port: 5432

---
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: production
type: Opaque
stringData:
  secrets.yaml: |
    database:
      username: "app_user"
      password: "secure_password"
      ssl_cert: |
        -----BEGIN CERTIFICATE-----
        MIIBIjANBgkqh...
        -----END CERTIFICATE-----

    auth:
      jwt_secret: "your_jwt_secret_here"
      encryption_key: "your_encryption_key"
```

### Go Implementation
```go
package main

import (
    "fmt"
    "os"
    "github.com/afeiship/go-dotpath"
)

type KubernetesConfig struct {
    config *dotpath.DotIO
    secrets *dotpath.DotIO
}

func LoadKubernetesConfig() (*KubernetesConfig, error) {
    // Load main config from ConfigMap
    configPath := os.Getenv("CONFIG_PATH")
    if configPath == "" {
        configPath = "/etc/config/app.yaml"
    }

    config := dotpath.NewDotIO(dotpath.YAML)
    if err := config.LoadFromFile(configPath); err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    // Load secrets from Secret
    secretsPath := os.Getenv("SECRETS_PATH")
    if secretsPath == "" {
        secretsPath = "/etc/secrets/secrets.yaml"
    }

    secrets := dotpath.NewDotIO(dotpath.YAML)
    if err := secrets.LoadFromFile(secretsPath); err != nil {
        return nil, fmt.Errorf("failed to load secrets: %w", err)
    }

    return &KubernetesConfig{
        config: config,
        secrets: secrets,
    }, nil
}

func (kc *KubernetesConfig) GetDatabaseConfig() (map[string]any, error) {
    dbConfig := make(map[string]any)

    // Copy main database config
    dbConfigData := kc.config.Data()["database"].(map[string]any)
    for k, v := range dbConfigData {
        dbConfig[k] = v
    }

    // Add secrets
    dbSecrets := kc.secrets.Data()["database"].(map[string]any)
    for k, v := range dbSecrets {
        dbConfig[k] = v
    }

    return dbConfig, nil
}

func (kc *KubernetesConfig) UpdateConfig(path string, value any) error {
    kc.config.Set(path, value)

    // Write back to file if running outside Kubernetes
    if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
        return kc.config.SaveToFile("/etc/config/app.yaml")
    }

    return nil
}

func main() {
    kc, err := LoadKubernetesConfig()
    if err != nil {
        log.Fatalf("Failed to load Kubernetes config: %v", err)
    }

    fmt.Printf("App Name: %s\n", kc.config.GetString("app.name"))
    fmt.Printf("Environment: %s\n", kc.config.GetString("app.environment"))

    // Get complete database configuration
    dbConfig, err := kc.GetDatabaseConfig()
    if err != nil {
        log.Fatalf("Failed to get database config: %v", err)
    }

    fmt.Printf("Database: %s@%s:%d/%s\n",
        dbConfig["username"],
        dbConfig["host"],
        dbConfig["port"],
        dbConfig["database"])

    // Get auth secrets
    jwtSecret := kc.secrets.GetString("auth.jwt_secret")
    if jwtSecret != "" {
        fmt.Println("JWT Secret: ***REDACTED***")
    }

    // Update configuration
    err = kc.UpdateConfig("server.workers", 8)
    if err != nil {
        fmt.Printf("Failed to update config: %v\n", err)
    } else {
        fmt.Println("Configuration updated successfully!")
    }
}
```

## 6. Configuration Validation

### Go Implementation
```go
package main

import (
    "fmt"
    "log"
    "regexp"
    "github.com/afeiship/go-dotpath"
)

type YAMLConfigValidator struct {
    io       *dotpath.DotIO
    required []string
    checks   map[string]func(*dotpath.DotIO) error
}

func NewYAMLConfigValidator(io *dotpath.DotIO) *YAMLConfigValidator {
    return &YAMLConfigValidator{
        io:     io,
        checks: make(map[string]func(*dotpath.DotIO) error),
    }
}

func (vcv *YAMLConfigValidator) AddRequired(path string) *YAMLConfigValidator {
    vcv.required = append(vcv.required, path)
    return vcv
}

func (vcv *YAMLConfigValidator) AddCheck(path string, check func(*dotpath.DotIO) error) *YAMLConfigValidator {
    vcv.checks[path] = check
    return vcv
}

func (vcv *YAMLConfigValidator) Validate() error {
    // Check required fields
    for _, path := range vcv.required {
        if !vcv.io.Has(path) {
            return fmt.Errorf("missing required configuration: %s", path)
        }
    }

    // Run custom checks
    for path, check := range vcv.checks {
        if err := check(vcv.io); err != nil {
            return fmt.Errorf("validation failed for %s: %w", path, err)
        }
    }

    return nil
}

func main() {
    io := dotpath.NewDotIO(dotpath.YAML)
    if err := io.LoadFromFile("config.yaml"); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Setup validator
    validator := NewYAMLConfigValidator(io).
        AddRequired("app.name").
        AddRequired("app.environment").
        AddRequired("server.port").
        AddRequired("database.driver").
        AddCheck("app.name", func(io *dotpath.DotIO) error {
            name := io.GetString("app.name")
            if len(name) < 3 {
                return fmt.Errorf("app name must be at least 3 characters")
            }
            if matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]*$`, name); !matched {
                return fmt.Errorf("app name must contain only lowercase letters, numbers, and hyphens")
            }
            return nil
        }).
        AddCheck("server.port", func(io *dotpath.DotIO) error {
            port := io.GetInt("server.port")
            if port <= 0 || port > 65535 {
                return fmt.Errorf("invalid port number: %d", port)
            }
            return nil
        }).
        AddCheck("server.timeout", func(io *dotpath.DotIO) error {
            timeout := io.GetString("server.timeout")
            if timeout == "" {
                return fmt.Errorf("server timeout cannot be empty")
            }
            // Validate duration format
            if matched, _ := regexp.MatchString(`^\d+[smh]$", timeout); !matched {
                return fmt.Errorf("invalid timeout format: %s (expected format: 30s, 5m, 1h)", timeout)
            }
            return nil
        }).
        AddCheck("database.driver", func(io *dotpath.DotIO) error {
            driver := io.GetString("database.driver")
            validDrivers := []string{"postgres", "mysql", "sqlite", "mongodb"}
            for _, valid := range validDrivers {
                if driver == valid {
                    return nil
                }
            }
            return fmt.Errorf("invalid database driver: %s", driver)
        }).
        AddCheck("logging.level", func(io *dotpath.DotIO) error {
            level := io.GetString("logging.level")
            validLevels := []string{"debug", "info", "warn", "error", "fatal"}
            for _, valid := range validLevels {
                if level == valid {
                    return nil
                }
            }
            return fmt.Errorf("invalid logging level: %s", level)
        }).
        AddCheck("monitoring.enabled", func(io *dotpath.DotIO) error {
            enabled := io.GetBool("monitoring.enabled")
            if enabled {
                // If monitoring is enabled, health_check should also be enabled
                if !io.GetBool("monitoring.health_check.enabled") {
                    return fmt.Errorf("monitoring is enabled but health_check is disabled")
                }
            }
            return nil
        })

    // Validate configuration
    if err := validator.Validate(); err != nil {
        log.Fatalf("Configuration validation failed: %v", err)
    }

    fmt.Println("YAML configuration is valid!")
    fmt.Printf("App: %s (Env: %s)\n", io.GetString("app.name"), io.GetString("app.environment"))
    fmt.Printf("Server: Port %d, Timeout %s\n", io.GetInt("server.port"), io.GetString("server.timeout"))
    fmt.Printf("Database: %s\n", io.GetString("database.driver"))
    fmt.Printf("Logging Level: %s\n", io.GetString("logging.level"))
    fmt.Printf("Monitoring: %t\n", io.GetBool("monitoring.enabled"))
}
```

## Best Practices for YAML Configuration

### 1. Use Environment-Specific Files
```yaml
# config.yaml (base)
app:
  name: "my-service"

# config.dev.yaml
app:
  environment: "development"

# config.prod.yaml
app:
  environment: "production"
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
            "timeout": "30s",
        },
    }
}
```

### 3. Use Configuration Environment Variables
```go
func GetConfigPath() string {
    if path := os.Getenv("CONFIG_PATH"); path != "" {
        return path
    }
    if env := os.Getenv("GO_ENV"); env != "" {
        return fmt.Sprintf("config.%s.yaml", env)
    }
    return "config.yaml"
}
```

### 4. Validate Configuration Structure
```go
func ValidateYAMLStructure(io *dotpath.DotIO) error {
    // Check for valid YAML structure
    data := io.Data()
    if data == nil {
        return fmt.Errorf("empty configuration")
    }

    // Validate top-level keys
    validKeys := []string{"app", "server", "database", "logging", "monitoring"}
    for key := range data {
        valid := false
        for _, validKey := range validKeys {
            if key == validKey {
                valid = true
                break
            }
        }
        if !valid {
            return fmt.Errorf("unknown top-level key: %s", key)
        }
    }

    return nil
}
```

## Comment & Order Preservation

The go-dotpath library provides intelligent YAML comment and order preservation when saving files.

### Example with Comments

```yaml
# Database configuration
database:
  host: localhost
  port: 5432
  timeout: 30s  # Connection timeout in seconds

# Application settings
app:
  name: "my-app"
  version: "1.0.0"
  debug: false  # Set to true for development
```

```go
package main

import (
    "log"
    "github.com/afeiship/go-dotpath"
)

func main() {
    // Load YAML with comments
    adapter, err := dotpath.LoadYAMLAdapter("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // Modify configuration
    adapter.Set("database.username", "admin")
    adapter.Set("app.debug", true)
    adapter.Set("monitoring.enabled", true)

    // Save back to original file - preserves comments and order
    err = adapter.Save()
    if err != nil {
        log.Fatal(err)
    }

    // Save to new file - preserves comments from original
    err = adapter.SaveAs("config-backup.yaml")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Resulting Files

**Original `config.yaml` after `Save()`:**
```yaml
# Database configuration
database:
  host: localhost
  port: 5432
  timeout: 30s  # Connection timeout in seconds
  username: admin

# Application settings
app:
  name: "my-app"
  version: "1.0.0"
  debug: true  # Set to true for development
```

**New `config-backup.yaml` after `SaveAs()`:**
```yaml
# Database configuration
database:
  host: localhost
  port: 5432
  timeout: 30s  # Connection timeout in seconds
  username: admin

# Application settings
app:
  name: "my-app"
  version: "1.0.0"
  debug: true  # Set to true for development

monitoring:
  enabled: true
```

### Best Practices

1. **Use `Save()`** for updating the same file you loaded from
2. **Use `SaveAs()`** for creating backups or copies with preserved comments
3. **Comments** are preserved exactly as they appear in the source file
4. **Field order** is maintained from the original structure

These examples demonstrate practical ways to use go-dotpath for YAML configuration management in real applications.