package config

import (
	"os"
	"strings"
)

// EnvLoader load config from env vars
type EnvLoader struct {
	prefix string
}

// NewEnvLoader create a new env loader
func NewEnvLoader(prefix string) *EnvLoader {
	return &EnvLoader{prefix: prefix}
}

// Load env vars into config
func (e *EnvLoader) Load(c *Config) error {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) < 2 {
			continue
		}

		key, value := pair[0], pair[1]

		// Apply prefix if is defined
		if e.prefix != "" && !strings.HasPrefix(key, e.prefix) {
			continue
		}

		// Normalize key
		configKey := strings.TrimPrefix(key, e.prefix)
		configKey = strings.ToLower(configKey)
		configKey = strings.ReplaceAll(configKey, "_", ".")

		// Asign value
		c.Data[configKey] = value
	}
	return nil
}

// LoadFromEnv carga variables de entorno en la configuración
func (c *Config) LoadFromEnv(prefix string) error {
	loader := NewEnvLoader(prefix)
	return loader.Load(c)
}
