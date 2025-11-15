package cache

import (
	"fmt"
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
)

const benchmarkCacheSize = 1000

// --- Benchmarks for MY LRU Cache ---

func BenchmarkMyCache_Set(b *testing.B) {
	cache := NewLRUCache(benchmarkCacheSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%benchmarkCacheSize)
		cache.Set(key, i, 0)
	}
}

func BenchmarkMyCache_Get(b *testing.B) {
	cache := NewLRUCache(benchmarkCacheSize)
	for i := 0; i < benchmarkCacheSize; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, 0)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%benchmarkCacheSize)
		cache.Get(key)
	}
}

// --- Benchmarks for Hashicorp's LRU Cache ---

func BenchmarkHashicorpLRU_Set(b *testing.B) {
	// Note: Hashicorp's New() can return an error, which we ignore in a benchmark
	cache, _ := lru.New[string, int](benchmarkCacheSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%benchmarkCacheSize)
		cache.Add(key, i) // Their method is called Add, not Set
	}
}

func BenchmarkHashicorpLRU_Get(b *testing.B) {
	cache, _ := lru.New[string, int](benchmarkCacheSize)
	for i := 0; i < benchmarkCacheSize; i++ {
		cache.Add(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%benchmarkCacheSize)
		cache.Get(key) // Their Get is the same
	}
}

// BenchmarkMyCache_Set-8                             	10133540	       117.6 ns/op	      21 B/op	       2 allocs/op
// BenchmarkMyCache_Get-8                             	11157621	       106.8 ns/op	      13 B/op	       1 allocs/op
// BenchmarkHashicorpLRU_Set-8                        	10124403	       118.7 ns/op	      13 B/op	       1 allocs/op
// BenchmarkHashicorpLRU_Get-8                        	10293021	       116.4 ns/op	      13 B/op	       1 allocs/op
