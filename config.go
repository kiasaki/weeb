package main

import (
	"os"
	"strings"
)

// Config represents a config instance
type Config struct {
	values map[string]string
}

// NewConfig creates a new config instance
func NewConfig() *Config {
	return &Config{values: map[string]string{}}
}

// Get gets a config value defaulting to `alternative` if non-present
func (c *Config) Get(key, alternative string) string {
	if value, ok := c.values[key]; !ok || value == "" {
		return alternative
	} else {
		return value
	}
}

// Set sets a config value
func (c *Config) Set(key, value string) {
	c.values[key] = value
}

// LoadFromEnv loads all environment variables starting with APP_ into config.
// It's camel cases all the keys so that `APP_DATABASE_URL` becomes `databaseUrl`.
func (c *Config) LoadFromEnv() {
	for _, value := range os.Environ() {
		if value[0:4] == "APP_" {
			splittedValue := strings.SplitN(value, "=", 2)
			c.values[camelCase(splittedValue[0])] = splittedValue[1]
		}
	}
}

// LoadValues loads a set of values all at once as config
func (c *Config) LoadValues(values map[string]string) {
	for k, v := range values {
		c.values[k] = v
	}
}
