package config

import (
	"os"
	"path/filepath"
	"testing"
)

// getTestFilePath tries to get the path of the test file, first
// searching the current directory and, if it does not exist, the parent directory.
func getTestFilePath(t *testing.T, fileName string) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("Error getting working directory:", err)
	}

	// Try first at: <wd>/testdata/<fileName>
	path := filepath.Join(wd, "testdata", fileName)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	// If not found, try: <wd>/../testdata/<fileName>
	path = filepath.Join(wd, "..", "testdata", fileName)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	t.Fatalf("Test file not found %s", fileName)
	return ""
}

func TestLoadAndValidate(t *testing.T) {
	cfg := NewConfig()
	basePath := getTestFilePath(t, "base_config.json")
	if err := cfg.LoadFromJSON(basePath); err != nil {
		t.Fatal("Error loading base_config.json:", err)
	}

	// Validate existing keys
	if err := cfg.Validate([]string{"db.host", "db.port", "app.name"}); err != nil {
		t.Fatal("Validation error:", err)
	}

	// Test to obtain values
	if host, err := cfg.GetString("db.host"); err != nil || host != "localhost" {
		t.Fatalf("Expected host 'localhost', got '%v' (err: %v)", host, err)
	}

	if port, err := cfg.GetInt("db.port"); err != nil || port != 5432 {
		t.Fatalf("Expected port 5432, obtained '%v' (err: %v)", port, err)
	}
}

func TestMergeConfig(t *testing.T) {
	cfg := NewConfig()
	basePath := getTestFilePath(t, "base_config.json")
	overridePath := getTestFilePath(t, "override.json")

	if err := cfg.LoadFromJSON(basePath); err != nil {
		t.Fatal("Error loading base_config.json:", err)
	}

	overrideCfg := NewConfig()
	if err := overrideCfg.LoadFromJSON(overridePath); err != nil {
		t.Fatal("Error loading override.json:", err)
	}

	cfg.Merge(overrideCfg)

	// After the merge, db.host must be overwritten
	if host, err := cfg.GetString("db.host"); err != nil || host != "192.168.1.100" {
		t.Fatalf("Expected host '192.168.1.100', obtained '%v' (err: %v)", host, err)
	}

	// Additionally, the logging.level key must exist.
	if level, err := cfg.GetString("logging.level"); err != nil || level != "debug" {
		t.Fatalf("Expected logging.level 'debug', got '%v' (err: %v)", level, err)
	}
}

func TestMissingKey(t *testing.T) {
	cfg := NewConfig()
	missingPath := getTestFilePath(t, "missing.json")
	if err := cfg.LoadFromJSON(missingPath); err != nil {
		t.Fatal("Error loading missing.json:", err)
	}

	// The key "db.host" is missing
	if err := cfg.Validate([]string{"db.host"}); err == nil {
		t.Fatal("Error expected due to missing key 'db.host'")
	}
}

func TestGetStringFunction(t *testing.T) {
	cfg := NewConfig()
	cfg.data["simple"] = "hello"
	cfg.data["number"] = 123

	s, err := cfg.GetString("simple")
	if err != nil || s != "hello" {
		t.Fatalf("Expected 'hello', obtained '%v' (err: %v)", s, err)
	}

	s, err = cfg.GetString("number")
	if err != nil || s != "123" {
		t.Fatalf("Expected '123', obtained '%v' (err: %v)", s, err)
	}
}

func TestGetIntFunction(t *testing.T) {
	cfg := NewConfig()
	cfg.data["intValue"] = 100
	cfg.data["floatValue"] = 99.0
	cfg.data["stringValue"] = "42"

	i, err := cfg.GetInt("intValue")
	if err != nil || i != 100 {
		t.Fatalf("Expected 100, obtained %v (err: %v)", i, err)
	}

	i, err = cfg.GetInt("floatValue")
	if err != nil || i != 99 {
		t.Fatalf("Expected 99, obtained %v (err: %v)", i, err)
	}

	i, err = cfg.GetInt("stringValue")
	if err != nil || i != 42 {
		t.Fatalf("Expected 42, obtained %v (err: %v)", i, err)
	}
}

func TestGetBoolFunction(t *testing.T) {
	cfg := NewConfig()
	cfg.data["boolTrue"] = true
	cfg.data["boolString"] = "true"

	b, err := cfg.GetBool("boolTrue")
	if err != nil || !b {
		t.Fatalf("Expected true, got %v (err: %v)", b, err)
	}

	b, err = cfg.GetBool("boolString")
	if err != nil || !b {
		t.Fatalf("Expected true, got %v (err: %v)", b, err)
	}
}
