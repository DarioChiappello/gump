package config

import "time"

// Configurer define operations config interface
type Configurer interface {
	Getter
	Validator
}

// Loader define load config interface
type Loader interface {
	LoadFromJSON(filePath string) error
	Merge(other *Config)
}

// Config implements interfaces
type Config struct {
	Data         map[string]interface{}
	LastModified time.Time
}

// NewConfig create a new instance
func NewConfig() *Config {
	return &Config{Data: make(map[string]interface{})}
}

func (c *Config) SetData(data map[string]interface{}) {
	c.Data = data
}
