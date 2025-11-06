package cache

import (
	"fmt"
	"testing"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := NewLRUCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Set(key, i)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewLRUCache(1000)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Get(key)
	}
}

func BenchmarkCacheSetWithEviction(b *testing.B) {
	cache := NewLRUCache(100) // Small cache = constant eviction
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, i)
	}
}
