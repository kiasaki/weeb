package weeb

import (
	"fmt"
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

// Values returns a copy of all configured values
func (c *Config) Values() map[string]string {
	valuesCopy := map[string]string{}
	for k, v := range c.values {
		if v != "" {
			valuesCopy[k] = v
		}
	}
	return valuesCopy
}

// GetBool gets a boolean value for a config key. Where the only true value a
// config can have is '1' or 'true'
func (c *Config) GetBool(key string) bool {
	value := c.Get(key, "")
	return value == "1" || value == "true"
}

// Get gets a config value defaulting to `alternative` if non-present
func (c *Config) Get(key, alternative string) string {
	value, ok := c.values[key]
	if !ok || value == "" {
		return alternative
	}
	return value
}

// MustGet is the same as Get but panics when a key is missing
func (c *Config) MustGet(key string) string {
	if value, ok := c.values[key]; !ok || value == "" {
		panic(fmt.Sprintf("Config.MustGet: Key '%s' not found", key))
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
			splittedValue := strings.SplitN(value[4:], "=", 2)
			c.values[ToCamelCase(splittedValue[0])] = splittedValue[1]
		}
	}
}

// LoadValues loads a set of values all at once as config
func (c *Config) LoadValues(values map[string]string) {
	for k, v := range values {
		c.values[k] = v
	}
}
