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
	Type        string // Provider type (e.g., "ecs", "cms", "sls", "oss", "vpc", "slb")
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

// CachedFactory wraps a provider factory with caching.
// If a cached provider exists for the given key, it returns the cached one.
// Otherwise, it calls the underlying factory and caches the result.
func CachedFactory[T any](providerType string, factory func(ak, sk, region string) (T, error)) func(ak, sk, region string) (T, error) {
	return func(ak, sk, region string) (T, error) {
		key := CacheKey{Type: providerType, AccessKeyID: ak, Region: region}
		if cached, ok := globalProviderCache.Get(key); ok {
			if typed, ok := cached.(T); ok {
				return typed, nil
			}
		}
		provider, err := factory(ak, sk, region)
		if err != nil {
			return provider, err
		}
		globalProviderCache.Set(key, provider)
		return provider, nil
	}
}
