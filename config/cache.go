package config

import "sync"

// cachedValue represents a value in cache
type cachedValue struct {
	value interface{}
	valid bool
}

// ConfigWithCache extends Config with capabilities of cache
type ConfigWithCache struct {
	*Config
	cache map[string]cachedValue
	mu    sync.RWMutex
}

// NewConfigWithCache create a config with cache
func NewConfigWithCache(base *Config) *ConfigWithCache {
	return &ConfigWithCache{
		Config: base,
		cache:  make(map[string]cachedValue),
	}
}

// InvalidateCache clean cache
func (c *ConfigWithCache) InvalidateCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]cachedValue)
}

// InvalidateKey clean an specific key from cache
func (c *ConfigWithCache) InvalidateKey(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

// GetString with cache
func (c *ConfigWithCache) GetString(key string) (string, error) {
	if val, ok := c.getFromCache(key); ok {
		return val.(string), nil
	}

	val, err := c.Config.GetString(key)
	if err != nil {
		return "", err
	}

	c.setCache(key, val)
	return val, nil
}

// GetInt with cache
func (c *ConfigWithCache) GetInt(key string) (int, error) {
	if val, ok := c.getFromCache(key); ok {
		return val.(int), nil
	}

	val, err := c.Config.GetInt(key)
	if err != nil {
		return 0, err
	}

	c.setCache(key, val)
	return val, nil
}

// GetBool with cache
func (c *ConfigWithCache) GetBool(key string) (bool, error) {
	if val, ok := c.getFromCache(key); ok {
		return val.(bool), nil
	}

	val, err := c.Config.GetBool(key)
	if err != nil {
		return false, err
	}

	c.setCache(key, val)
	return val, nil
}

func (c *ConfigWithCache) getFromCache(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, exists := c.cache[key]
	if exists && cached.valid {
		return cached.value, true
	}
	return nil, false
}

func (c *ConfigWithCache) setCache(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cachedValue{value: value, valid: true}
}
