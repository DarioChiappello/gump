package config

import (
	"testing"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	// Test config
	cfg := config.NewConfig()
	cfg.SetData(map[string]interface{}{
		"app": map[string]interface{}{
			"name":    "GUMP",
			"version": "1.0.0",
		},
		"db": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"credentials": map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
		},
		"debug": true,
	})

	t.Run("Successfully validation - all keys existing", func(t *testing.T) {
		keys := []string{
			"app.name",
			"db.host",
			"db.credentials.username",
			"debug",
		}
		err := cfg.Validate(keys)
		assert.NoError(t, err)
	})

	t.Run("Failed validation - key missing", func(t *testing.T) {
		keys := []string{
			"app.name",
			"db.ssl", // no exists
			"debug",
		}
		err := cfg.Validate(keys)
		require.Error(t, err)

		keyErr, ok := err.(*config.KeyError)
		require.True(t, ok, "Must be a KeyError")
		assert.Equal(t, "db.ssl", keyErr.Key)
		assert.Equal(t, "missing required key: db.ssl", keyErr.Error())
	})

	t.Run("Failed validation - multiples keys missing", func(t *testing.T) {
		keys := []string{
			"app.description",     // missing
			"db.credentials.cert", // missing
			"log_level",           // missing
		}
		err := cfg.Validate(keys)
		require.Error(t, err)

		// Only must report 1st missing key
		keyErr, ok := err.(*config.KeyError)
		require.True(t, ok, "Must be a KeyError")
		assert.Equal(t, "app.description", keyErr.Key)
	})

	t.Run("Successfully validation - nested deep keys", func(t *testing.T) {
		keys := []string{
			"db.credentials.password",
		}
		err := cfg.Validate(keys)
		assert.NoError(t, err)
	})

	t.Run("Failed validation - intermediate segment is not a map", func(t *testing.T) {
		keys := []string{
			"app.name.invalid", // "name" is string, no map
		}
		err := cfg.Validate(keys)
		require.Error(t, err)

		pathErr, ok := err.(*config.PathError)
		require.True(t, ok, "Must be a PathError")
		assert.Equal(t, "app.name.invalid", pathErr.Key)
		assert.Equal(t, "name", pathErr.Segment)
		assert.Equal(t, "invalid path segment 'name' in key: app.name.invalid", pathErr.Error())
	})

	t.Run("Successfully validation - empty keys list", func(t *testing.T) {
		err := cfg.Validate([]string{})
		assert.NoError(t, err)
	})

	t.Run("Validación con valor nil", func(t *testing.T) {
		// Config with nil value
		cfgWithNil := config.NewConfig()
		cfgWithNil.SetData(map[string]interface{}{
			"valid_key": "value",
			"nil_key":   nil,
		})

		t.Run("Key with nil value should be considered existing", func(t *testing.T) {
			err := cfgWithNil.Validate([]string{"nil_key"})
			assert.NoError(t, err)
		})

		t.Run("Key with nil value in nested route", func(t *testing.T) {
			cfgWithNestedNil := config.NewConfig()
			cfgWithNestedNil.SetData(map[string]interface{}{
				"parent": map[string]interface{}{
					"child": nil,
				},
			})

			err := cfgWithNestedNil.Validate([]string{"parent.child"})
			assert.NoError(t, err)
		})
	})

	t.Run("Validation with different errors types", func(t *testing.T) {
		t.Run("KeyError format", func(t *testing.T) {
			err := &config.KeyError{Key: "missing.key"}
			assert.Equal(t, "missing required key: missing.key", err.Error())
		})

		t.Run("PathError format", func(t *testing.T) {
			err := &config.PathError{Key: "parent.child.grandchild", Segment: "child"}
			assert.Equal(t, "invalid path segment 'child' in key: parent.child.grandchild", err.Error())
		})
	})

	t.Run("Combined validation with GetValue", func(t *testing.T) {
		keys := []string{"app.version", "db.port"}
		err := cfg.Validate(keys)
		require.NoError(t, err)

		// Verify that exists and are able to get
		version, err := cfg.GetString("app.version")
		require.NoError(t, err)
		assert.Equal(t, "1.0.0", version)

		port, err := cfg.GetInt("db.port")
		require.NoError(t, err)
		assert.Equal(t, 5432, port)
	})
}
