package config

import (
	"testing"
	"time"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToString(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"String", "hello", "hello"},
		{"Int", 42, "42"},
		{"Int8", int8(8), "8"},
		{"Int16", int16(16), "16"},
		{"Int32", int32(32), "32"},
		{"Int64", int64(64), "64"},
		{"Float32", float32(3.14), "3.14"},
		{"Float64", 3.1415926535, "3.1415926535"},
		{"Bool true", true, "true"},
		{"Bool false", false, "false"},
		{"Uint", uint(42), "42"},
		{"Byte", byte('A'), "65"},
		{"Time", now, now.String()},
		{"Struct", struct{ Name string }{"test"}, "{test}"},
		{"Nil", nil, "<nil>"},
		{"Array", [3]int{1, 2, 3}, "[1 2 3]"},
		{"Slice", []string{"a", "b"}, "[a b]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := config.ConvertToString(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertToInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected int
		err      bool
	}{
		{"Int", 42, 42, false},
		{"Int8", int8(8), 8, false},
		{"Int16", int16(16), 16, false},
		{"Int32", int32(32), 32, false},
		{"Int64", int64(64), 64, false},
		{"Float32", float32(3.14), 3, false},
		{"Float64", 3.99, 3, false},
		{"String integer", "123", 123, false},
		{"String float", "3.14", 3, false},
		{"String with spaces", " 42 ", 42, false},
		{"Negative int", -42, -42, false},
		{"Negative float", -3.14, -3, false},
		{"Negative string", "-42", -42, false},
		{"Large number", 1_000_000, 1_000_000, false},
		{"Max int", 1<<31 - 1, 1<<31 - 1, false}, // 2147483647

		// Error cases
		{"String not number", "abc", 0, true},
		{"Bool true", true, 0, true},
		{"Bool false", false, 0, true},
		{"Nil", nil, 0, true},
		{"Time", time.Now(), 0, true},
		{"Struct", struct{}{}, 0, true},
		{"Array", [3]int{}, 0, true},
		{"Slice", []int{}, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := config.ConvertToInt(tc.input, "test_key")
			if tc.err {
				require.Error(t, err)
				assert.IsType(t, &config.TypeError{}, err)
				assert.Contains(t, err.Error(), "expected int")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertToBool(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected bool
		err      bool
	}{
		{"Bool true", true, true, false},
		{"Bool false", false, false, false},
		{"String true", "true", true, false},
		{"String TRUE", "TRUE", true, false},
		{"String True", "True", true, false},
		{"String t", "t", true, false},
		{"String T", "T", true, false},
		{"String yes", "yes", true, false},
		{"String y", "y", true, false},
		{"String on", "on", true, false},
		{"String 1", "1", true, false},
		{"String false", "false", false, false},
		{"String FALSE", "FALSE", false, false},
		{"String False", "False", false, false},
		{"String f", "f", false, false},
		{"String F", "F", false, false},
		{"String no", "no", false, false},
		{"String n", "n", false, false},
		{"String off", "off", false, false},
		{"String 0", "0", false, false},
		{"Int 1", 1, true, false},
		{"Int 0", 0, false, false},
		{"Int8 1", int8(1), true, false},
		{"Int8 0", int8(0), false, false},
		{"Int16 1", int16(1), true, false},
		{"Int16 0", int16(0), false, false},
		{"Int32 1", int32(1), true, false},
		{"Int32 0", int32(0), false, false},
		{"Int64 1", int64(1), true, false},
		{"Int64 0", int64(0), false, false},
		{"Uint 1", uint(1), true, false},
		{"Uint 0", uint(0), false, false},
		{"Float32 1.0", float32(1.0), true, false},
		{"Float32 0.0", float32(0.0), false, false},
		{"Float64 1.0", 1.0, true, false},
		{"Float64 0.0", 0.0, false, false},
		{"Float64 0.5", 0.5, true, false}, // Any value != 0 is true

		// Error cases
		{"String invalid", "invalid", false, true},
		{"Nil", nil, false, true},
		{"Time", time.Now(), false, true},
		{"Struct", struct{}{}, false, true},
		{"Array", [3]bool{}, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := config.ConvertToBool(tc.input, "test_key")
			if tc.err {
				require.Error(t, err)
				assert.IsType(t, &config.TypeError{}, err)
				assert.Contains(t, err.Error(), "expected bool")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, result,
				"Case '%s' with input %v", tc.name, tc.input)
		})
	}
}

func TestTypeError(t *testing.T) {
	t.Run("Int error message", func(t *testing.T) {
		err := &config.TypeError{
			Key:      "server.port",
			Expected: "int",
			Actual:   "string",
		}
		assert.Equal(t, "invalid type for key 'server.port': expected int, got string", err.Error())
	})

	t.Run("Bool error message", func(t *testing.T) {
		err := &config.TypeError{
			Key:      "debug",
			Expected: "bool",
			Actual:   "int",
		}
		assert.Equal(t, "invalid type for key 'debug': expected bool, got int", err.Error())
	})
}
