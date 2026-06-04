package service

import (
	"sync"
)

// ProviderCache caches provider instances for reuse.
type ProviderCache struct {
	cache sync.Map
}

// CacheKey is the key for caching a provider.
type CacheKey struct {
	AccessKeyID string
	Region      string
}

// NewProviderCache creates a new ProviderCache.
func NewProviderCache() *ProviderCache {
	return &ProviderCache{}
}

// Get gets a provider from the cache.
func (c *ProviderCache) Get(key CacheKey) (interface{}, bool) {
	return c.cache.Load(key)
}

// Set sets a provider in the cache.
func (c *ProviderCache) Set(key CacheKey, provider interface{}) {
	c.cache.Store(key, provider)
}

// Delete deletes a provider from the cache.
func (c *ProviderCache) Delete(key CacheKey) {
	c.cache.Delete(key)
}

// Clear clears the cache.
func (c *ProviderCache) Clear() {
	c.cache.Clear()
}

// globalProviderCache is the global provider cache.
var globalProviderCache = NewProviderCache()

// GetProviderCache returns the global provider cache.
func GetProviderCache() *ProviderCache {
	return globalProviderCache
}
