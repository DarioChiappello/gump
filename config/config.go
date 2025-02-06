package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	data map[string]interface{}
}

func NewConfig() *Config {
	return &Config{data: make(map[string]interface{})}
}

// LoadFromJSON loads configuration from a JSON file
func (c *Config) LoadFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var tempData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&tempData); err != nil {
		return err
	}

	mergeMaps(c.data, tempData)
	return nil
}

// LoadFromEnvironment loads configuration from environment variables with prefix
func (c *Config) LoadFromEnvironment(prefix string) error {
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key, value := kv[0], kv[1]
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		keyParts := strings.Split(strings.TrimPrefix(key, prefix), "__")
		var cleanParts []string
		for _, part := range keyParts {
			if part != "" {
				cleanParts = append(cleanParts, part)
			}
		}

		if len(cleanParts) == 0 {
			continue
		}

		currentMap := c.data
		for i, part := range cleanParts[:len(cleanParts)-1] {
			if _, exists := currentMap[part]; !exists {
				currentMap[part] = make(map[string]interface{})
			}

			nextMap, ok := currentMap[part].(map[string]interface{})
			if !ok {
				return fmt.Errorf("config conflict at %s", strings.Join(cleanParts[:i+1], "."))
			}
			currentMap = nextMap
		}
		currentMap[cleanParts[len(cleanParts)-1]] = value
	}
	return nil
}

// Merge combines another configuration into current one
func (c *Config) Merge(other *Config) {
	mergeMaps(c.data, other.data)
}

// Get methods with type conversion
func (c *Config) GetString(key string) (string, error) {
	val, err := c.getValue(key)
	if err != nil {
		return "", err
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

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

// Validate checks existence of required keys
func (c *Config) Validate(keys []string) error {
	for _, key := range keys {
		if _, err := c.getValue(key); err != nil {
			return fmt.Errorf("missing required key: %s", key)
		}
	}
	return nil
}

// Helper functions
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

		_, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid path segment: %s", part)
		}
	}
	return nil, fmt.Errorf("invalid key path: %s", key)
}

func mergeMaps(dest, src map[string]interface{}) {
	for key, srcVal := range src {
		destVal, exists := dest[key]
		if exists && isMap(srcVal) && isMap(destVal) {
			mergeMaps(destVal.(map[string]interface{}), srcVal.(map[string]interface{}))
		} else {
			dest[key] = srcVal
		}
	}
}

func isMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}
