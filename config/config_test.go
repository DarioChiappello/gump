package config

import (
	"os"
	"testing"
)

func TestLoadFromJSON(t *testing.T) {
	t.Run("Valid JSON File", func(t *testing.T) {
		cfg := NewConfig()
		err := cfg.LoadFromJSON("testdata/base_config.json")
		if err != nil {
			t.Fatalf("Failed to load JSON: %v", err)
		}

		host, err := cfg.GetString("db.host")
		if err != nil || host != "localhost" {
			t.Errorf("Expected 'localhost', got '%s'", host)
		}

		port, err := cfg.GetInt("db.port")
		if err != nil || port != 5432 {
			t.Errorf("Expected 5432, got %d", port)
		}
	})

	t.Run("Invalid File Path", func(t *testing.T) {
		cfg := NewConfig()
		err := cfg.LoadFromJSON("nonexistent.json")
		if err == nil {
			t.Error("Expected error for missing file")
		}
	})
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set test environment variables
	os.Setenv("APP_DB__HOST", "env-host")
	os.Setenv("APP_DB__PORT", "8080")
	os.Setenv("APP_LOG_LEVEL", "debug")
	os.Setenv("APP_FEATURES__NEWUI", "true")
	defer func() {
		os.Unsetenv("APP_DB__HOST")
		os.Unsetenv("APP_DB__PORT")
		os.Unsetenv("APP_LOG_LEVEL")
		os.Unsetenv("APP_FEATURES__NEWUI")
	}()

	t.Run("Valid Environment Loading", func(t *testing.T) {
		cfg := NewConfig()
		err := cfg.LoadFromEnvironment("APP_")
		if err != nil {
			t.Fatalf("Failed to load environment: %v", err)
		}

		host, err := cfg.GetString("db.host")
		if err != nil || host != "env-host" {
			t.Errorf("Expected 'env-host', got '%s'", host)
		}

		port, err := cfg.GetInt("db.port")
		if err != nil || port != 8080 {
			t.Errorf("Expected 8080, got %d", port)
		}
	})

	t.Run("Environment Type Conflict", func(t *testing.T) {
		cfg := NewConfig()
		cfg.data["db"] = "invalid" // Set non-map value
		err := cfg.LoadFromEnvironment("APP_")
		if err == nil {
			t.Error("Expected error for type conflict")
		}
	})
}

func TestMerge(t *testing.T) {
	cfg1 := NewConfig()
	cfg1.LoadFromJSON("testdata/base_config.json")

	cfg2 := NewConfig()
	cfg2.LoadFromJSON("testdata/emergency.json")

	cfg1.Merge(cfg2)

	host, _ := cfg1.GetString("db.host")
	if host != "backup-server" {
		t.Errorf("Expected merged host 'backup-server', got '%s'", host)
	}

	port, _ := cfg1.GetInt("db.port")
	if port != 5432 {
		t.Errorf("Expected original port 5432, got %d", port)
	}
}

func TestGetMethods(t *testing.T) {
	cfg := NewConfig()
	cfg.data = map[string]interface{}{
		"string_val": "hello",
		"int_val":    42,
		"bool_val":   true,
		"nested": map[string]interface{}{
			"number": "123",
		},
	}

	tests := []struct {
		name     string
		key      string
		expected interface{}
		err      bool
	}{
		{"Valid String", "string_val", "hello", false},
		{"Valid Int", "int_val", 42, false},
		{"Valid Bool", "bool_val", true, false},
		{"Convert String to Int", "nested.number", 123, false},
		{"Missing Key", "missing", nil, true},
		{"Invalid Conversion", "string_val", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.expected.(type) {
			case string:
				val, err := cfg.GetString(tt.key)
				if (err != nil) != tt.err {
					t.Errorf("Unexpected error: %v", err)
				}
				if !tt.err && val != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, val)
				}
			case int:
				val, err := cfg.GetInt(tt.key)
				if (err != nil) != tt.err {
					t.Errorf("Unexpected error: %v", err)
				}
				if !tt.err && val != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, val)
				}
			case bool:
				val, err := cfg.GetBool(tt.key)
				if (err != nil) != tt.err {
					t.Errorf("Unexpected error: %v", err)
				}
				if !tt.err && val != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, val)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	cfg := NewConfig()
	cfg.LoadFromJSON("testdata/base_config.json")

	t.Run("All Keys Present", func(t *testing.T) {
		err := cfg.Validate([]string{"db.host", "log_level"})
		if err != nil {
			t.Errorf("Validation failed unexpectedly: %v", err)
		}
	})

	t.Run("Missing Key", func(t *testing.T) {
		err := cfg.Validate([]string{"db.host", "missing.key"})
		if err == nil {
			t.Error("Expected validation error for missing key")
		}
	})
}

func TestComplexScenario(t *testing.T) {
	// Set environment variables
	os.Setenv("PROD_DB__PORT", "6432")
	os.Setenv("PROD_LOG_LEVEL", "warn")
	defer os.Unsetenv("PROD_DB__PORT")
	defer os.Unsetenv("PROD_LOG_LEVEL")

	cfg := NewConfig()
	// Load base config
	cfg.LoadFromJSON("testdata/base_config.json")
	// Load environment overrides
	cfg.LoadFromEnvironment("PROD_")
	// Merge emergency config
	emergencyCfg := NewConfig()
	emergencyCfg.LoadFromJSON("testdata/emergency.json")
	cfg.Merge(emergencyCfg)

	// Verify merged values
	host, _ := cfg.GetString("db.host")
	if host != "backup-server" {
		t.Errorf("Expected host from emergency config, got %s", host)
	}

	port, _ := cfg.GetInt("db.port")
	if port != 6432 {
		t.Errorf("Expected port from environment, got %d", port)
	}

	logLevel, _ := cfg.GetString("log_level")
	if logLevel != "warn" {
		t.Errorf("Expected log level from environment, got %s", logLevel)
	}
}
