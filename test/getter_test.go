package config

import (
	"testing"
	"time"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetter(t *testing.T) {
	// Config for test
	cfg := config.NewConfig()
	cfg.SetData(map[string]interface{}{
		"string": "hello",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"time":   time.Now(),
		"nested": map[string]interface{}{
			"value": "world",
			"deep": map[string]interface{}{
				"value": 100,
			},
		},
		"invalid_int":  "not a number",
		"invalid_bool": "not a boolean",
	})

	t.Run("getValue - 1st level existing key", func(t *testing.T) {
		val, err := cfg.GetValue("string")
		require.NoError(t, err)
		assert.Equal(t, "hello", val)
	})

	t.Run("getValue - nested key", func(t *testing.T) {
		val, err := cfg.GetValue("nested.value")
		require.NoError(t, err)
		assert.Equal(t, "world", val)
	})

	t.Run("getValue - nested deep value", func(t *testing.T) {
		val, err := cfg.GetValue("nested.deep.value")
		require.NoError(t, err)
		assert.Equal(t, 100, val)
	})

	t.Run("getValue - no existing value", func(t *testing.T) {
		_, err := cfg.GetValue("nonexistent")
		assert.Error(t, err)
		assert.IsType(t, &config.KeyError{}, err)
		assert.Contains(t, err.Error(), "nonexistent")
	})

	t.Run("getValue - segment intermediate isnt mapa", func(t *testing.T) {
		_, err := cfg.GetValue("string.invalid")
		assert.Error(t, err)
		assert.IsType(t, &config.PathError{}, err)
		assert.Contains(t, err.Error(), "invalid path segment")
	})

	t.Run("GetString - string value", func(t *testing.T) {
		val, err := cfg.GetString("string")
		require.NoError(t, err)
		assert.Equal(t, "hello", val)
	})

	t.Run("GetString - int value", func(t *testing.T) {
		val, err := cfg.GetString("int")
		require.NoError(t, err)
		assert.Equal(t, "42", val)
	})

	t.Run("GetString - float value", func(t *testing.T) {
		val, err := cfg.GetString("float")
		require.NoError(t, err)
		assert.Equal(t, "3.14", val)
	})

	t.Run("GetString - bool value", func(t *testing.T) {
		val, err := cfg.GetString("bool")
		require.NoError(t, err)
		assert.Equal(t, "true", val)
	})

	t.Run("GetString - no existing key", func(t *testing.T) {
		_, err := cfg.GetString("nonexistent")
		assert.Error(t, err)
	})

	t.Run("GetInt - int value", func(t *testing.T) {
		val, err := cfg.GetInt("int")
		require.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("GetInt - float value", func(t *testing.T) {
		val, err := cfg.GetInt("float")
		require.NoError(t, err)
		assert.Equal(t, 3, val) // truncado
	})

	t.Run("GetInt - string value convertible", func(t *testing.T) {
		cfg.Data["int_string"] = "123"
		val, err := cfg.GetInt("int_string")
		require.NoError(t, err)
		assert.Equal(t, 123, val)
	})

	t.Run("GetInt - string value no convertible", func(t *testing.T) {
		_, err := cfg.GetInt("invalid_int")
		assert.Error(t, err)
		assert.IsType(t, &config.TypeError{}, err)
	})

	t.Run("GetInt - no convertible type", func(t *testing.T) {
		_, err := cfg.GetInt("time")
		assert.Error(t, err)
		assert.IsType(t, &config.TypeError{}, err)
	})

	t.Run("GetInt - no existing key", func(t *testing.T) {
		_, err := cfg.GetInt("nonexistent")
		assert.Error(t, err)
	})

	t.Run("GetBool - bool value", func(t *testing.T) {
		val, err := cfg.GetBool("bool")
		require.NoError(t, err)
		assert.True(t, val)
	})

	t.Run("GetBool - string value 'true'", func(t *testing.T) {
		cfg.Data["true_string"] = "true"
		val, err := cfg.GetBool("true_string")
		require.NoError(t, err)
		assert.True(t, val)
	})

	t.Run("GetBool - string value 'false'", func(t *testing.T) {
		cfg.Data["false_string"] = "false"
		val, err := cfg.GetBool("false_string")
		require.NoError(t, err)
		assert.False(t, val)
	})

	t.Run("GetBool - string value no convertible", func(t *testing.T) {
		_, err := cfg.GetBool("invalid_bool")
		assert.Error(t, err)
		assert.IsType(t, &config.TypeError{}, err)
	})

	t.Run("GetBool - no convertible type", func(t *testing.T) {
		_, err := cfg.GetBool("time")
		assert.Error(t, err)
		assert.IsType(t, &config.TypeError{}, err)
	})

	t.Run("GetBool - no existing key", func(t *testing.T) {
		_, err := cfg.GetBool("nonexistent")
		assert.Error(t, err)
	})

	t.Run("Test errors personalized", func(t *testing.T) {
		// KeyError
		keyErr := &config.KeyError{Key: "test"}
		assert.Equal(t, "missing required key: test", keyErr.Error())

		// PathError
		pathErr := &config.PathError{Key: "parent.child", Segment: "child"}
		assert.Equal(t, "invalid path segment 'child' in key: parent.child", pathErr.Error())

		// TypeError
		typeErr := &config.TypeError{Key: "age", Expected: "int", Actual: "string"}
		assert.Equal(t, "invalid type for key 'age': expected int, got string", typeErr.Error())
	})
}
