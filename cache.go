package main

import (
	"encoding/json"
	"errors"
)

type Cache interface {
	Get(key string, value interface{}) error
	Set(key string, result interface{}) error
	Del(key string) error
}

var CacheKeyNotFoundError = errors.New("Cache key not found.")

type MemoryCache struct {
	contents map[string][]byte
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{contents: map[string][]byte{}}
}

func (c *MemoryCache) Get(key string, result interface{}) error {
	if bytes, ok := c.contents[key]; !ok {
		return CacheKeyNotFoundError
	} else {
		return json.Unmarshal(bytes, result)
	}
}
func (c *MemoryCache) Set(key string, value interface{}) error {
	if bytes, err := json.Marshal(value); err != nil {
		return err
	} else {
		c.contents[key] = bytes
		return nil
	}
}
func (c *MemoryCache) Del(key string) error {
	delete(c.contents, key)
	return nil
}
