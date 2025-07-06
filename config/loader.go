package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func (c *Config) LoadFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var tempData map[string]interface{}
	if err := decoder.Decode(&tempData); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	c.Merge(&Config{Data: tempData})
	return nil
}
