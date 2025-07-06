# 🧠 GUMP - Go Unified Management Package

> **Advanced, reactive, and extensible configuration management for Go.**

---

## 📘 Overview

**GUMP** simplifies complex configuration management in Go by unifying multiple sources like JSON files and environment variables. It supports real-time change detection, efficient caching, and a fluent builder API for declarative configuration setups.

---

## ✨ Key Features

### 🔧 Core Capabilities

- **Multi-source Loading**: JSON files, environment variables, and more  
- **Smart Merging**: Hierarchical config merging with override support  
- **Robust Validation**: Ensure required keys and types are correct  
- **Typed Access**: Strong typing with sensible defaults  
- **Dot Notation**: Easy access to nested fields (`app.database.host`)

---

### 🚀 New Modules

#### 🔍 ConfigWatcher  
> Real-time configuration updates with callbacks  
- Monitors files via `fsnotify`  
- Periodic verification with customizable intervals  
- Callback support for change events  
- Multi-file support  
- Integrated cache invalidation

#### 🔄 EnvLoader  
> Simplified environment variable integration  
- Prefix-based ENV loading (e.g., `APP_`)  
- Auto-converts `ENV_VARIABLE` → `env.variable`  
- Seamlessly merges into your config structure

#### 🧱 ConfigBuilder  
> Fluent configuration builder with chaining  
- Compose from multiple sources  
- Cumulative error aggregation  
- `MustBuild()` for safe one-liners

#### ⚡ ConfigWithCache  
> Thread-safe cache layer  
- Fast access with `sync.RWMutex`  
- Selective or full invalidation  
- Compatible with all `Get*` methods  
- Strongly typed and performant

---

### 🔐 Security & Robustness

- Nil-safe merge operations  
- Extensive type conversions and checks  
- Nested map and slice handling  
- Error recovery and panic protection

---

### ⚙️ Optimizations

- Uniform API receivers  
- Zero-copy conversions  
- Recursive map efficiency  
- Thread-safe concurrent access

---

## 📦 Installation

```bash
go get github.com/DarioChiappello/gump
```

---

## 🚀 Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/DarioChiappello/gump/config"
)

func main() {
	cfg, err := config.NewConfigBuilder().
		AddJSONFile("base_config.json").
		AddJSONFile("overrides.json").
		AddEnv("APP_").
		Build()

	if err != nil {
		panic(err)
	}

	cachedCfg := config.NewConfigWithCache(cfg)

	watcher, _ := config.NewConfigWatcher(cachedCfg, 5*time.Second, "config.json")
	watcher.OnReload(func(c *config.Config) {
		fmt.Println("Configuration updated!")
		cachedCfg.InvalidateCache()
	})
	go watcher.Start()
	defer watcher.Stop()

	dbHost := cachedCfg.GetString("database.host", "localhost")
	dbPort := cachedCfg.GetInt("database.port", 5432)

	fmt.Printf("Connecting to %s:%d\n", dbHost, dbPort)
}
```

---

## 🧪 Advanced Usage

### 🔨 ConfigBuilder

```go
builder := config.NewConfigBuilder().
	AddJSONFile("base.json").
	AddJSONFile("env/production.json").
	AddEnv("APP_")

cfg, err := builder.Build()
if err != nil {
	// Handle cumulative errors
}

cfg = builder.MustBuild()
```

---

### 🌱 EnvLoader

```go
cfg := config.NewConfig()
loader := config.NewEnvLoader("APP_", cfg)

err := loader.Load()
if err != nil {
	panic(err)
}

host := cfg.GetString("db.host", "localhost")
```

---

### 🔁 ConfigWatcher

```go
cfg := config.NewConfig()
cfg.LoadFromJSON("config.json")

watcher, _ := config.NewConfigWatcher(cfg, 5*time.Second, "config.json")

watcher.OnReload(func(c *config.Config) {
	fmt.Println("Configuration updated in real-time!")
})

go watcher.Start()
defer watcher.Stop()

select {}
```

---

### 🧠 ConfigWithCache

```go
cfg := config.NewConfig()
cfg.LoadFromJSON("config.json")

cachedCfg := config.NewConfigWithCache(cfg)

val := cachedCfg.GetString("complex.key", "default")

cachedCfg.InvalidateKey("complex.key")
cachedCfg.InvalidateAllCache()
```

---

## ✅ Benefits

- 🧩 **Modular**: Clean separation of logic  
- 🔌 **Extensible**: Easy to add new sources or hooks  
- 🧼 **Readable**: Context-aware, descriptive errors  
- 🛡️ **Reliable**: Handles edge cases and bad inputs gracefully  
- ⚡ **Fast**: In-memory cache for high-performance access  
- 🔄 **Reactive**: Instant config reload on file changes  

---

## 🧪 Running Tests

```bash
go test -v ./...
```

---

## 🤝 Contributions

We welcome contributions!  
Please [open an issue](https://github.com/DarioChiappello/gump/issues) to discuss large changes before submitting PRs.

---

## 📄 License

**MIT License**  
See the [LICENSE](./LICENSE) file for details.
