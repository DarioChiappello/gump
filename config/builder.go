package config

import "fmt"

// ConfigBuilder facilitate fluent config contruction
type ConfigBuilder struct {
	config *Config
	errors []error
}

// NewConfigBuilder create a new ConfigBuilder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: NewConfig(),
	}
}

// WithJSON add config from JSON file
func (b *ConfigBuilder) WithJSON(filePath string) *ConfigBuilder {
	if err := b.config.LoadFromJSON(filePath); err != nil {
		b.errors = append(b.errors, fmt.Errorf("JSON load error: %w", err))
	}
	return b
}

// WithEnv add config from env vars
func (b *ConfigBuilder) WithEnv(prefix string) *ConfigBuilder {
	if err := b.config.LoadFromEnv(prefix); err != nil {
		b.errors = append(b.errors, fmt.Errorf("ENV load error: %w", err))
	}
	return b
}

// WithConfig add an existing config
func (b *ConfigBuilder) WithConfig(cfg *Config) *ConfigBuilder {
	b.config.Merge(cfg)
	return b
}

// Build final config
func (b *ConfigBuilder) Build() (*Config, error) {
	if len(b.errors) > 0 {
		return nil, MultiError{Errors: b.errors}
	}
	return b.config, nil
}

// MustBuild build config or get in panic
func (b *ConfigBuilder) MustBuild() *Config {
	cfg, err := b.Build()
	if err != nil {
		panic(err)
	}
	return cfg
}
