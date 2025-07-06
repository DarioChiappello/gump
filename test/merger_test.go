package config

import (
	"fmt"
	"testing"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	t.Run("Merge simple - claves no existentes", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
			"b": 2,
		}}
		src := &config.Config{Data: map[string]interface{}{
			"c": 3,
			"d": 4,
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
			"d": 4,
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Merge simple - sobrescritura", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
			"b": 2,
		}}
		src := &config.Config{Data: map[string]interface{}{
			"b": 20,
			"c": 30,
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"a": 1,
			"b": 20,
			"c": 30,
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Merge anidado - claves no existentes", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 10,
			},
		}}
		src := &config.Config{Data: map[string]interface{}{
			"b": map[string]interface{}{
				"y": 20,
			},
			"c": 3,
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 10,
				"y": 20,
			},
			"c": 3,
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Merge anidado - sobrescritura recursiva", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 10,
				"y": map[string]interface{}{
					"p": 100,
				},
			},
		}}
		src := &config.Config{Data: map[string]interface{}{
			"b": map[string]interface{}{
				"y": map[string]interface{}{
					"q": 200,
				},
				"z": 30,
			},
			"c": 3,
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 10,
				"y": map[string]interface{}{
					"p": 100,
					"q": 200,
				},
				"z": 30,
			},
			"c": 3,
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Sobrescritura de valor con mapa", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 10,
			},
		}}
		src := &config.Config{Data: map[string]interface{}{
			"b": 20, // Sobrescribe el mapa con un entero
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"a": 1,
			"b": 20,
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Sobrescritura de mapa con valor", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
			"b": 2,
		}}
		src := &config.Config{Data: map[string]interface{}{
			"b": map[string]interface{}{ // Sobrescribe el entero con un mapa
				"x": 10,
			},
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"x": 10,
			},
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Merge con nil", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
		}}
		dest.Merge(nil) // No debería causar pánico
		assert.Equal(t, map[string]interface{}{"a": 1}, dest.Data)
	})

	t.Run("Merge con config vacía", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"a": 1,
		}}
		src := &config.Config{Data: map[string]interface{}{}}

		dest.Merge(src)
		assert.Equal(t, map[string]interface{}{"a": 1}, dest.Data)
	})

	t.Run("Merge de config vacía con datos", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{}}
		src := &config.Config{Data: map[string]interface{}{
			"a": 1,
		}}

		dest.Merge(src)
		assert.Equal(t, map[string]interface{}{"a": 1}, dest.Data)
	})

	t.Run("Merge múltiple niveles", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"level1": map[string]interface{}{
				"level2a": map[string]interface{}{
					"value":  "original",
					"common": "original",
				},
			},
		}}
		src := &config.Config{Data: map[string]interface{}{
			"level1": map[string]interface{}{
				"level2a": map[string]interface{}{
					"common": "overridden",
					"new":    "new",
				},
				"level2b": map[string]interface{}{
					"value": "new",
				},
			},
		}}

		dest.Merge(src)

		expected := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2a": map[string]interface{}{
					"value":  "original",
					"common": "overridden",
					"new":    "new",
				},
				"level2b": map[string]interface{}{
					"value": "new",
				},
			},
		}
		assert.Equal(t, expected, dest.Data)
	})

	t.Run("Merge con tipos diferentes", func(t *testing.T) {
		dest := &config.Config{Data: map[string]interface{}{
			"key": "original",
		}}
		src := &config.Config{Data: map[string]interface{}{
			"key": 42, // Sobrescribe string con int
		}}

		dest.Merge(src)

		// Verificar que la sobrescritura funciona
		val, ok := dest.Data["key"]
		require.True(t, ok)
		assert.Equal(t, 42, val)

		// Verificar con los métodos Get
		result, err := dest.GetInt("key")
		require.NoError(t, err)
		assert.Equal(t, 42, result)
	})
}

func BenchmarkMergeMaps(b *testing.B) {
	// Crear configuraciones grandes
	createLargeConfig := func(size int) *config.Config {
		data := make(map[string]interface{})
		for i := 0; i < size; i++ {
			key := fmt.Sprintf("key%d", i)
			data[key] = map[string]interface{}{
				"sub1": "value",
				"sub2": 42,
				"sub3": map[string]interface{}{
					"deep": true,
				},
			}
		}
		return &config.Config{Data: data}
	}

	dest := createLargeConfig(1000)
	src := createLargeConfig(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dest.Merge(src)
	}
}
