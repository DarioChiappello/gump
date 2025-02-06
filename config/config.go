package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the loaded configuration.
type Config struct {
	data map[string]interface{}
}

// NewConfig creates a new Config instance.
func NewConfig() *Config {
	return &Config{data: make(map[string]interface{})}
}

// LoadFromJSON loads the configuration from a JSON file.
func (c *Config) LoadFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	var tempData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&tempData); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	mergeMaps(c.data, tempData)
	return nil
}

// Merge combines another configuration into the current one.
func (c *Config) Merge(other *Config) {
	mergeMaps(c.data, other.data)
}

// Validate verifies the existence of required keys.
func (c *Config) Validate(keys []string) error {
	for _, key := range keys {
		if _, err := c.getValue(key); err != nil {
			return fmt.Errorf("missing required key: %s", key)
		}
	}
	return nil
}

// GetString gets a value of type string.
func (c *Config) GetString(key string) (string, error) {
	val, err := c.getValue(key)
	if err != nil {
		return "", err
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// GetInt gets a value of type int.
func (c *Config) GetInt(key string) (int, error) {
	val, err := c.getValue(key)
	if err != nil {
		return 0, err
	}

	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("invalid type for key %s", key)
	}
}

// GetBool gets a value of type bool.
func (c *Config) GetBool(key string) (bool, error) {
	val, err := c.getValue(key)
	if err != nil {
		return false, err
	}

	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("invalid type for key %s", key)
	}
}

// getValue gets a configuration value given a key in "dot" notation.
func (c *Config) getValue(key string) (interface{}, error) {
	parts := strings.Split(key, ".")
	current := c.data

	for i, part := range parts {
		val, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("key not found: %s", key)
		}

		if i == len(parts)-1 {
			return val, nil
		}

		nextMap, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid path segment: %s", part)
		}
		current = nextMap
	}
	return nil, fmt.Errorf("invalid key path: %s", key)
}

// mergeMaps merges two maps recursively.
func mergeMaps(dest, src map[string]interface{}) {
	for key, srcVal := range src {
		if destVal, exists := dest[key]; exists {
			// If both values ​​are maps, combine them recursively.
			destMap, destIsMap := destVal.(map[string]interface{})
			srcMap, srcIsMap := srcVal.(map[string]interface{})
			if destIsMap && srcIsMap {
				mergeMaps(destMap, srcMap)
				continue
			}
		}
		dest[key] = srcVal
	}
}
