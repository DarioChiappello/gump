package config

type Validator interface {
	Validate(keys []string) error
}

func (c *Config) Validate(keys []string) error {
	for _, key := range keys {
		if _, err := c.GetValue(key); err != nil {
			return err
			//return &KeyError{Key: key}
		}
	}
	return nil
}
