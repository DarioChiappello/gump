package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigBuilder(t *testing.T) {
	// Config JSON test routes files
	basePath := getTestFilePath(t, "base_config.json")
	overridePath := getTestFilePath(t, "override.json")
	// invalidPath := getTestFilePath(t, "invalid.json")
	// nonexistentPath := getTestFilePath(t, "nonexistent.json")

	// Config env vars for testing
	os.Setenv("APP_DB_HOST", "env-host")
	os.Setenv("APP_DB_PORT", "8080")
	os.Setenv("APP_LOGGING_LEVEL", "info")
	defer func() {
		os.Unsetenv("APP_DB_HOST")
		os.Unsetenv("APP_DB_PORT")
		os.Unsetenv("APP_LOGGING_LEVEL")
	}()

	t.Run("Successfully construction only with JSON", func(t *testing.T) {
		builder := config.NewConfigBuilder()
		cfg, err := builder.WithJSON(basePath).Build()
		require.NoError(t, err)

		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "localhost", host)
	})

	t.Run("Successfully construction with multiples sources", func(t *testing.T) {
		builder := config.NewConfigBuilder()
		cfg, err := builder.
			WithJSON(basePath).
			WithEnv("APP_").
			WithJSON(overridePath).
			Build()
		require.NoError(t, err)

		// Verify sources combination
		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "192.168.1.100", host) // of override.json

		port, err := cfg.GetInt("db.port")
		require.NoError(t, err)
		assert.Equal(t, 5432, port) // of env vars

		logLevel, err := cfg.GetString("logging.level")
		require.NoError(t, err)
		assert.Equal(t, "debug", logLevel) // of override.json
	})

	t.Run("Fusion with existing configuration", func(t *testing.T) {
		// Create base config
		baseCfg := config.NewConfig()
		err := baseCfg.LoadFromJSON(basePath)
		require.NoError(t, err)

		builder := config.NewConfigBuilder()
		cfg, err := builder.
			WithConfig(baseCfg).
			WithEnv("APP_").
			Build()
		require.NoError(t, err)

		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "localhost", host) // of env vars
	})

	t.Run("MustBuild successfully", func(t *testing.T) {
		builder := config.NewConfigBuilder()
		cfg := builder.
			WithJSON(basePath).
			WithEnv("APP_").
			MustBuild() // No panic

		assert.NotNil(t, cfg)
		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "localhost", host) // of env vars
	})

	t.Run("Precedence Order", func(t *testing.T) {
		// Last source added must have higher precedence
		builder := config.NewConfigBuilder()
		cfg, err := builder.
			WithJSON(basePath).     // db.host = "localhost"
			WithEnv("APP_").        // db.host = "localhost" (overwrite)
			WithJSON(overridePath). // db.host = "192.168.1.100" (overwrite)
			Build()
		require.NoError(t, err)

		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "192.168.1.100", host) // last source added
	})

	t.Run("Construction without sources", func(t *testing.T) {
		builder := config.NewConfigBuilder()
		cfg, err := builder.Build()
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Empty(t, cfg.Data) // Must be empty
	})
}

func TestMultiError(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")
	multiErr := config.MultiError{Errors: []error{err1, err2}}

	assert.Equal(t, "multiple errors: [error 1, error 2]", multiErr.Error())
}

func BenchmarkBuilder(b *testing.B) {
	basePath := filepath.Join("testdata", "base_config.json")

	for i := 0; i < b.N; i++ {
		builder := config.NewConfigBuilder()
		_, _ = builder.WithJSON(basePath).Build()
	}
}
