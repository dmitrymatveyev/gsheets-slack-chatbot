package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config provides configuration details
type Config struct {
	dict map[string]string
}

// NewConfig creates new instance of Config
func NewConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dict map[string]string
	err = json.Unmarshal(raw, &dict)
	if err != nil {
		return nil, err
	}

	return &Config{dict: dict}, nil
}

// Get returns a config value by key
func (c *Config) Get(key string) (string, error) {
	val, ok := c.dict[key]
	if !ok {
		return "", fmt.Errorf("the key \"%s\" is not presented", key)
	}
	return val, nil
}
