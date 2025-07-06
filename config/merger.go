package config

// Merge combine other config
func (c *Config) Merge(other *Config) {
	if other != nil {
		mergeMaps(c.Data, other.Data)
	}
}

func mergeMaps(dest, src map[string]interface{}) {
	for key, srcVal := range src {
		if destVal, exists := dest[key]; exists {
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
