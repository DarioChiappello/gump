# GUMP - Go Unified Management Package
## Overview
GUMP is an advanced configuration management package for Go that simplifies handling complex configurations. It combines multiple sources (JSON files, environment variables), offers real-time change monitoring, high-performance caching, and a fluent builder system.

## Key Features
### Core Capabilities
Multi-source Loading: JSON, environment variables, and more

Smart Merging: Hierarchical configuration merging

Robust Validation: Verify required keys and data types

Typed Access: Strings, integers, booleans with default values

Dot Notation: Simple access to nested values (app.database.host)

### New Features
🔍 ConfigWatcher:

Continuous file monitoring with fsnotify

Periodic verification with configurable intervals

Callbacks for real-time change reactions

Multi-file support

Cache integration (auto-invalidation)

🔄 EnvLoader:

Load environment variables with configurable prefix

Automatic conversion of ENV_VARIABLES to env.variables

Direct integration with Config structure

🧱 ConfigBuilder:

Builder pattern for fluent configuration

Combine multiple sources (JSON, ENV, etc.)

Cumulative error handling

MustBuild version for safe initialization

⚡ ConfigWithCache:

Thread-safe cache with sync.RWMutex

Selective or full cache invalidation

Strongly typed cached values

Compatible with all Get* methods

Security & Robustness Improvements
Nil validation in merge operations

Comprehensive type handling and conversions

Improved nested map management

Exhaustive type checking

Error recovery mechanisms

Optimizations
Consistent API receivers

Elimination of redundant conversions

Optimized recursive map handling

Concurrent access safety

### Installation

```
bash
go get github.com/DarioChiappello/gump
```

### Basic Usage

```
go
package main

import (
	"fmt"
	"time"
	
	"github.com/DarioChiappello/gump/config"
)

func main() {
	// Fluent multi-source configuration
	cfg, err := config.NewConfigBuilder().
		AddJSONFile("base_config.json").
		AddJSONFile("overrides.json").
		AddEnv("APP_").
		Build()
	
	if err != nil {
		panic(err)
	}

	// Cached configuration
	cachedCfg := config.NewConfigWithCache(cfg)

	// Change watcher
	watcher, _ := config.NewConfigWatcher(cachedCfg, 5*time.Second, "config.json")
	watcher.OnReload(func(c *config.Config) {
		fmt.Println("Configuration updated!")
		cachedCfg.InvalidateCache() // Clear cache on change
	})
	go watcher.Start()
	defer watcher.Stop()

	// Value access
	dbHost := cachedCfg.GetString("database.host", "localhost")
	dbPort := cachedCfg.GetInt("database.port", 5432)
	
	fmt.Printf("Connecting to %s:%d\n", dbHost, dbPort)
}
```

### Advanced Examples
#### ConfigBuilder

```
go
builder := config.NewConfigBuilder().
	AddJSONFile("base.json").
	AddJSONFile("env/production.json").
	AddEnv("APP_")

cfg, err := builder.Build()
if err != nil {
	// Handle cumulative errors
}

// Safe version for initializations
cfg = builder.MustBuild()
EnvLoader
go
cfg := config.NewConfig()
loader := config.NewEnvLoader("APP_", cfg)

// Load environment variables
err := loader.Load()
if err != nil {
	panic(err)
}

// APP_DB_HOST → db.host
host := cfg.GetString("db.host", "localhost")
ConfigWatcher
go
cfg := config.NewConfig()
cfg.LoadFromJSON("config.json")

// Create watcher with 5-second checks
watcher, _ := config.NewConfigWatcher(cfg, 5*time.Second, "config.json")

// Register change callback
watcher.OnReload(func(c *config.Config) {
	fmt.Println("Configuration updated in real-time!")
})

// Start monitoring
go watcher.Start()
defer watcher.Stop()

// Keep application running
select {}
ConfigWithCache
go
cfg := config.NewConfig()
cfg.LoadFromJSON("config.json")

// Create cached version
cachedCfg := config.NewConfigWithCache(cfg)

// First access - loads and caches
val := cachedCfg.GetString("complex.key", "default")

// Subsequent accesses - uses cache
val = cachedCfg.GetString("complex.key", "default")

// Invalidate specific key
cachedCfg.InvalidateKey("complex.key")

// Invalidate entire cache
cachedCfg.InvalidateAllCache()
```

### Key Benefits
Maintainability: Clear component responsibilities

Extensibility: Easy to add new source types

Clarity: Descriptive errors with context

Consistency: Uniform operation behavior

Security: Robust handling of edge cases

Performance: Cache access for complex configurations

Reactivity: Real-time configuration updates

### Running Tests
#### Run all tests
```
go test -v ./...
```


# Contributions
Contributions are welcome! Please open an issue to discuss significant changes before submitting PRs.

License
MIT License