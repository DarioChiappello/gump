package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DarioChiappello/gump/config"
)

func main() {
	// Get the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Error getting working directory: %v", err))
	}

	// It is assumed that when running from the 'example' folder the JSON files are located in '../testdata'
	baseConfigPath := filepath.Join(wd, "../..", "testdata", "base_config.json")
	emergencyConfigPath := filepath.Join(wd, "../..", "testdata", "emergency.json")
	overrideConfigPath := filepath.Join(wd, "../..", "testdata", "override.json")

	// --- Example 1: Load base configuration ---
	fmt.Println("Example 1: Load base configuration")
	cfg := config.NewConfig()
	if err := cfg.LoadFromJSON(baseConfigPath); err != nil {
		panic(fmt.Sprintf("Error loading base config: %v", err))
	}
	dbHost, _ := cfg.GetString("db.host")
	dbPort, _ := cfg.GetInt("db.port")
	dbSSL, _ := cfg.GetBool("db.ssl")
	fmt.Println("db.host =", dbHost)
	fmt.Println("db.port =", dbPort)
	fmt.Println("db.ssl  =", dbSSL)
	fmt.Println()

	// --- Example 2: Merge emergency configuration ---
	fmt.Println("Example 2: Merge emergency configuration")
	emergencyCfg := config.NewConfig()
	if err := emergencyCfg.LoadFromJSON(emergencyConfigPath); err != nil {
		panic(fmt.Sprintf("Error loading emergency config: %v", err))
	}
	cfg.Merge(emergencyCfg)
	dbSSL, _ = cfg.GetBool("db.ssl")
	fmt.Println("  Después de fusionar emergency config, db.ssl =", dbSSL)
	fmt.Println()

	// --- Example 3: Merge override configuration ---
	fmt.Println("Example 3: Merge override configuration")
	overrideCfg := config.NewConfig()
	if err := overrideCfg.LoadFromJSON(overrideConfigPath); err != nil {
		panic(fmt.Sprintf("Error loading override config: %v", err))
	}
	cfg.Merge(overrideCfg)
	dbHost, _ = cfg.GetString("db.host")
	logLevel, err := cfg.GetString("logging.level")
	if err != nil {
		logLevel = "not defined"
	}
	fmt.Println("After merging override config:")
	fmt.Println("db.host =", dbHost)
	fmt.Println("logging.level =", logLevel)
	fmt.Println()

	// --- Example 4: Validate required keys ---
	fmt.Println("Example 4: Validate required keys")
	err = cfg.Validate([]string{"db.host", "db.port", "app.name"})
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful: all required keys are present")
	}
}
