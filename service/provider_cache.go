package service

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

const defaultTTL = 30 * time.Minute

type cacheEntry struct {
	provider  any
	createdAt time.Time
}

type ProviderCache struct {
	cache sync.Map
	ttl   time.Duration
}

type CacheKey struct {
	Type        string
	AccessKeyID string
	SKHash      string
	Region      string
}

func NewProviderCache() *ProviderCache {
	return &ProviderCache{ttl: defaultTTL}
}

func NewProviderCacheWithTTL(ttl time.Duration) *ProviderCache {
	return &ProviderCache{ttl: ttl}
}

func (c *ProviderCache) Get(key CacheKey) (any, bool) {
	val, ok := c.cache.Load(key)
	if !ok {
		return nil, false
	}
	entry := val.(cacheEntry)
	if c.ttl > 0 && time.Since(entry.createdAt) > c.ttl {
		c.cache.Delete(key)
		return nil, false
	}
	return entry.provider, true
}

func (c *ProviderCache) Set(key CacheKey, provider any) {
	c.cache.Store(key, cacheEntry{provider: provider, createdAt: time.Now()})
	c.cleanup()
}

func (c *ProviderCache) Delete(key CacheKey) {
	c.cache.Delete(key)
}

func (c *ProviderCache) Clear() {
	c.cache.Range(func(key, _ any) bool {
		c.cache.Delete(key)
		return true
	})
}

func (c *ProviderCache) Size() int {
	n := 0
	c.cache.Range(func(_, _ any) bool {
		n++
		return true
	})
	return n
}

func (c *ProviderCache) cleanup() {
	if c.ttl <= 0 {
		return
	}
	c.cache.Range(func(key, val any) bool {
		entry := val.(cacheEntry)
		if time.Since(entry.createdAt) > c.ttl {
			c.cache.Delete(key)
		}
		return true
	})
}

var globalProviderCache = NewProviderCache()

func GetProviderCache() *ProviderCache {
	return globalProviderCache
}

func hashSK(sk string) string {
	h := sha256.Sum256([]byte(sk))
	return hex.EncodeToString(h[:8])
}

func CachedFactory[T any](providerType string, factory func(ak, sk, region string) (T, error)) func(ak, sk, region string) (T, error) {
	return func(ak, sk, region string) (T, error) {
		key := CacheKey{Type: providerType, AccessKeyID: ak, SKHash: hashSK(sk), Region: region}
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
