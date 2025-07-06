package config

import "strings"

type Getter interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetBool(key string) (bool, error)
	GetValue(key string) (interface{}, error)
}

func (c *Config) GetString(key string) (string, error) {
	val, err := c.GetValue(key)
	if err != nil {
		return "", err
	}
	return ConvertToString(val)
}

func (c *Config) GetInt(key string) (int, error) {
	val, err := c.GetValue(key)
	if err != nil {
		return 0, err
	}
	return ConvertToInt(val, key)
}

func (c *Config) GetBool(key string) (bool, error) {
	val, err := c.GetValue(key)
	if err != nil {
		return false, err
	}
	return ConvertToBool(val, key)
}

func (c *Config) GetValue(key string) (interface{}, error) {
	current := c.Data
	parts := strings.Split(key, ".")

	for i, part := range parts {
		val, exists := current[part]
		if !exists {
			return nil, &KeyError{Key: key}
		}

		if i == len(parts)-1 {
			return val, nil
		}

		nextMap, ok := val.(map[string]interface{})
		if !ok {
			return nil, &PathError{Key: key, Segment: part}
		}
		current = nextMap
	}
	return nil, &KeyError{Key: key}
}
