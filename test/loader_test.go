package config

import (
	"os"
	"testing"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromJSON(t *testing.T) {
	// Config test routes files
	basePath := getTestFilePath(t, "base_config.json")
	emergencyPath := getTestFilePath(t, "emergency.json")
	missingPath := getTestFilePath(t, "missing.json")
	overridePath := getTestFilePath(t, "override.json")
	t.Run("Load valid file base", func(t *testing.T) {
		cfg := config.NewConfig()
		err := cfg.LoadFromJSON(basePath)
		require.NoError(t, err)

		// Verify values
		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "localhost", host)

		port, err := cfg.GetInt("db.port")
		require.NoError(t, err)
		assert.Equal(t, 5432, port)

		ssl, err := cfg.GetBool("db.ssl")
		require.NoError(t, err)
		assert.False(t, ssl)

		appName, err := cfg.GetString("app.name")
		require.NoError(t, err)
		assert.Equal(t, "GUMP App", appName)
	})

	t.Run("Fusion with emergency file", func(t *testing.T) {
		cfg := config.NewConfig()
		err := cfg.LoadFromJSON(basePath)
		require.NoError(t, err)
		err = cfg.LoadFromJSON(emergencyPath)
		require.NoError(t, err)

		// Modified value
		ssl, err := cfg.GetBool("db.ssl")
		require.NoError(t, err)
		assert.True(t, ssl)

		// Value without changes
		port, err := cfg.GetInt("db.port")
		require.NoError(t, err)
		assert.Equal(t, 5432, port)
	})

	t.Run("Fusion with file with missing keys", func(t *testing.T) {
		cfg := config.NewConfig()
		err := cfg.LoadFromJSON(basePath)
		require.NoError(t, err)
		err = cfg.LoadFromJSON(missingPath)
		require.NoError(t, err)

		// Modified Value
		port, err := cfg.GetInt("db.port")
		require.NoError(t, err)
		assert.Equal(t, 3306, port)

		// Value without changes
		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "localhost", host)
	})

	t.Run("Fusion with file that add new keys", func(t *testing.T) {
		cfg := config.NewConfig()
		err := cfg.LoadFromJSON(basePath)
		require.NoError(t, err)
		err = cfg.LoadFromJSON(overridePath)
		require.NoError(t, err)

		// Modified value
		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "192.168.1.100", host)

		// New key added
		logLevel, err := cfg.GetString("logging.level")
		require.NoError(t, err)
		assert.Equal(t, "debug", logLevel)
	})

	t.Run("Load multiple files", func(t *testing.T) {
		cfg := config.NewConfig()
		err := cfg.LoadFromJSON(basePath)
		require.NoError(t, err)
		err = cfg.LoadFromJSON(emergencyPath)
		require.NoError(t, err)
		err = cfg.LoadFromJSON(overridePath)
		require.NoError(t, err)

		// Verify values combination
		host, err := cfg.GetString("db.host")
		require.NoError(t, err)
		assert.Equal(t, "192.168.1.100", host)

		port, err := cfg.GetInt("db.port")
		require.NoError(t, err)
		assert.Equal(t, 5432, port) // of base_config

		ssl, err := cfg.GetBool("db.ssl")
		require.NoError(t, err)
		assert.True(t, ssl) // of emergency

		logLevel, err := cfg.GetString("logging.level")
		require.NoError(t, err)
		assert.Equal(t, "debug", logLevel) // of override
	})
}

func createInvalidJSON(t *testing.T, path string) {
	// Create invalid JSON file
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString(`{ "invalid": json }`)
	require.NoError(t, err)
}
