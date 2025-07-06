package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/DarioChiappello/gump/config"
)

func main() {
	// Get actual directory
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting working directory: %v", err))
	}

	// Config routes files
	baseConfigPath := filepath.Join(wd, "../..", "testdata", "base_config.json")
	emergencyConfigPath := filepath.Join(wd, "../..", "testdata", "emergency.json")
	overrideConfigPath := filepath.Join(wd, "../..", "testdata", "override.json")

	// --- Example 1: ConfigBuilder advanced use ---
	fmt.Println("\n=== Ejemplo 1: ConfigBuilder avanzado ===")
	builder := config.NewConfigBuilder().
		WithJSON(baseConfigPath).
		WithJSON(emergencyConfigPath).
		WithJSON(overrideConfigPath)

	// env vars config
	os.Setenv("APP_DB_PORT", "3306")
	os.Setenv("APP_LOGGING_LEVEL", "trace")
	builder.WithEnv("APP_")

	builtCfg, err := builder.Build()
	if err != nil {
		panic(fmt.Sprintf("Error building config: %v", err))
	}

	dbHost, _ := builtCfg.GetString("db.host")
	dbPort, _ := builtCfg.GetInt("db.port")
	dbSSL, _ := builtCfg.GetBool("db.ssl")
	logLevel, _ := builtCfg.GetString("logging.level")
	fmt.Println("Configuración construida:")
	fmt.Println("db.host =", dbHost)
	fmt.Println("db.port =", dbPort) // Overwrite by env var
	fmt.Println("db.ssl  =", dbSSL)
	fmt.Println("logging.level =", logLevel) // Overwrite by env var

	// --- Example 2: Config with cache ---
	fmt.Println("\n=== Example 2: Config with cache ===")
	cachedCfg := config.NewConfigWithCache(builtCfg)

	// First access - without cache
	start := time.Now()
	cachedHost, _ := cachedCfg.GetString("db.host")
	cachedTime := time.Since(start)
	fmt.Printf("First access: %s (took %v)\n", cachedHost, cachedTime)

	// Second access - from cache
	start = time.Now()
	cachedHost, _ = cachedCfg.GetString("db.host")
	cachedTime = time.Since(start)
	fmt.Printf("Acceso desde caché: %s (tomó %v)\n", cachedHost, cachedTime)

	// Invalidate cache and access again
	cachedCfg.InvalidateCache()
	start = time.Now()
	cachedHost, _ = cachedCfg.GetString("db.host")
	cachedTime = time.Since(start)
	fmt.Printf("After invalidate cache: %s (took %v)\n", cachedHost, cachedTime)

	// --- Example 3: Config observer (Watcher) ---
	fmt.Println("\n=== Example 3: Config Watcher ===")
	watcherCfg := config.NewConfig()
	if err := watcherCfg.LoadFromJSON(baseConfigPath); err != nil {
		panic(fmt.Sprintf("Error loading base config: %v", err))
	}

	// Create watcher
	watcher, err := config.NewConfigWatcher(watcherCfg, 1*time.Second, baseConfigPath, overrideConfigPath)
	if err != nil {
		panic(fmt.Sprintf("Error creating watcher: %v", err))
	}

	// Register callback
	watcher.OnReload(func(c *config.Config) {
		fmt.Println("\n¡Configuration reloaded!")
		host, _ := c.GetString("db.host")
		port, _ := c.GetInt("db.port")
		ssl, _ := c.GetBool("db.ssl")
		fmt.Printf("New values: host=%s, port=%d, ssl=%t\n", host, port, ssl)

		// Invalidate cache in cache config
		cachedCfg.InvalidateCache()
	})

	// Init watcher on second plane
	go watcher.Start()
	defer watcher.Stop()

	fmt.Println("Observer initialized. Modify config files to watch changes...")

	// Simulate changes after a time
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("\nModifying base_config.json...")
		modifyFile(baseConfigPath, `{"app": {"name": "GUMP"}, "db": {"host": "watched_host", "port": 9999, "ssl": true}}`)

		time.Sleep(3 * time.Second)
		fmt.Println("\nModifying override.json...")
		modifyFile(overrideConfigPath, `{"db": {"host": "watched_host_override"}, "logging": {"level": "trace"}}`)
	}()

	// Mantener el programa en ejecución
	time.Sleep(10 * time.Second)
	fmt.Println("\nStopping observer...")

	// --- Example 4: Detailed Validation ---
	fmt.Println("\n=== Ejemplo 4: Validación avanzada ===")
	err = builtCfg.Validate([]string{
		"db.host",
		"db.port",
		"app.name",
		"logging.level",
		"nonexistent.key",
		"app.name.invalid",
	})

	if err != nil {
		fmt.Println("Validation Errors:")
		switch e := err.(type) {
		case *config.KeyError:
			fmt.Printf("- Required key missing: %s\n", e.Key)
		case *config.PathError:
			fmt.Printf("- Invalid Route: %s (segment: %s)\n", e.Key, e.Segment)
		default:
			fmt.Printf("- Error: %v\n", err)
		}
	}

	// --- Example 5: Complex Types and convertions ---
	fmt.Println("\n=== Example 5: Complex Types and convertions ===")
	complexCfg := config.NewConfig()
	complexCfg.SetData(map[string]interface{}{
		"int_as_float":   15.0,
		"int_as_string":  "25",
		"bool_as_string": "true",
		"nested": map[string]interface{}{
			"value": "100",
		},
	})

	// Implicit Convertions
	intFromFloat, _ := complexCfg.GetInt("int_as_float")
	intFromString, _ := complexCfg.GetInt("int_as_string")
	boolFromString, _ := complexCfg.GetBool("bool_as_string")
	nestedInt, _ := complexCfg.GetInt("nested.value")

	fmt.Println("Convertions:")
	fmt.Printf("int_as_float (15.0) -> int: %d\n", intFromFloat)
	fmt.Printf("int_as_string (\"25\") -> int: %d\n", intFromString)
	fmt.Printf("bool_as_string (\"true\") -> bool: %t\n", boolFromString)
	fmt.Printf("nested.value (\"100\") -> int: %d\n", nestedInt)
}

// Auxiliar function to modify files
func modifyFile(path, content string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}
