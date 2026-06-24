package service

import (
	"testing"
	"time"
)

func TestProviderCache_GetSet(t *testing.T) {
	c := NewProviderCache()
	key := CacheKey{Type: "ecs", AccessKeyID: "ak1", SKHash: "hash1", Region: "cn-hangzhou"}

	_, ok := c.Get(key)
	if ok {
		t.Fatal("expected miss on empty cache")
	}

	c.Set(key, "provider1")
	val, ok := c.Get(key)
	if !ok {
		t.Fatal("expected hit after set")
	}
	if val != "provider1" {
		t.Fatalf("expected provider1, got %v", val)
	}
}

func TestProviderCache_SKIsolation(t *testing.T) {
	c := NewProviderCache()
	key1 := CacheKey{Type: "ecs", AccessKeyID: "ak1", SKHash: hashSK("sk1"), Region: "cn-hangzhou"}
	key2 := CacheKey{Type: "ecs", AccessKeyID: "ak1", SKHash: hashSK("sk2"), Region: "cn-hangzhou"}

	c.Set(key1, "provider-sk1")
	c.Set(key2, "provider-sk2")

	val1, ok1 := c.Get(key1)
	if !ok1 || val1 != "provider-sk1" {
		t.Fatalf("expected provider-sk1, got %v", val1)
	}
	val2, ok2 := c.Get(key2)
	if !ok2 || val2 != "provider-sk2" {
		t.Fatalf("expected provider-sk2, got %v", val2)
	}
}

func TestProviderCache_TTLExpiry(t *testing.T) {
	c := NewProviderCacheWithTTL(50 * time.Millisecond)
	key := CacheKey{Type: "ecs", AccessKeyID: "ak1", SKHash: "hash1", Region: "cn-hangzhou"}

	c.Set(key, "provider1")
	val, ok := c.Get(key)
	if !ok || val != "provider1" {
		t.Fatal("expected hit before TTL")
	}

	time.Sleep(60 * time.Millisecond)

	_, ok = c.Get(key)
	if ok {
		t.Fatal("expected miss after TTL")
	}
}

func TestProviderCache_Cleanup(t *testing.T) {
	c := NewProviderCacheWithTTL(50 * time.Millisecond)

	for i := range 5 {
		key := CacheKey{Type: "ecs", AccessKeyID: "ak", SKHash: "h", Region: string(rune('a' + i))}
		c.Set(key, i)
	}
	if c.Size() != 5 {
		t.Fatalf("expected 5 entries, got %d", c.Size())
	}

	time.Sleep(60 * time.Millisecond)

	freshKey := CacheKey{Type: "ecs", AccessKeyID: "ak", SKHash: "h", Region: "fresh"}
	c.Set(freshKey, "fresh")

	if c.Size() != 1 {
		t.Fatalf("expected 1 entry after cleanup, got %d", c.Size())
	}
}

func TestProviderCache_Clear(t *testing.T) {
	c := NewProviderCache()
	for i := range 3 {
		key := CacheKey{Type: "ecs", AccessKeyID: "ak", SKHash: "h", Region: string(rune('a' + i))}
		c.Set(key, i)
	}
	if c.Size() != 3 {
		t.Fatalf("expected 3, got %d", c.Size())
	}

	c.Clear()
	if c.Size() != 0 {
		t.Fatalf("expected 0 after clear, got %d", c.Size())
	}
}

func TestProviderCache_Delete(t *testing.T) {
	c := NewProviderCache()
	key := CacheKey{Type: "ecs", AccessKeyID: "ak1", SKHash: "h", Region: "cn-hangzhou"}
	c.Set(key, "provider1")
	c.Delete(key)

	_, ok := c.Get(key)
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestHashSK_DifferentInputs(t *testing.T) {
	h1 := hashSK("secret1")
	h2 := hashSK("secret2")
	if h1 == h2 {
		t.Fatal("different SKs should produce different hashes")
	}
	if len(h1) != 16 {
		t.Fatalf("expected 16 char hex, got %d", len(h1))
	}
}

func TestCachedFactory_ReusesClient(t *testing.T) {
	globalProviderCache.Clear()
	calls := 0
	factory := CachedFactory("test", func(ak, sk, region string) (string, error) {
		calls++
		return "provider", nil
	})

	factory("ak", "sk", "cn-hangzhou")
	factory("ak", "sk", "cn-hangzhou")

	if calls != 1 {
		t.Fatalf("expected 1 factory call (cached), got %d", calls)
	}
	globalProviderCache.Clear()
}

func TestCachedFactory_DifferentSKCreatesNew(t *testing.T) {
	globalProviderCache.Clear()
	calls := 0
	factory := CachedFactory("test", func(ak, sk, region string) (string, error) {
		calls++
		return "provider-" + sk, nil
	})

	r1, _ := factory("ak", "sk1", "cn-hangzhou")
	r2, _ := factory("ak", "sk2", "cn-hangzhou")

	if calls != 2 {
		t.Fatalf("expected 2 factory calls (different SK), got %d", calls)
	}
	if r1 == r2 {
		t.Fatal("different SKs should produce different providers")
	}
	globalProviderCache.Clear()
}
