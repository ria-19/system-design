package cache

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := NewLRUCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Set(key, i, 0)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewLRUCache(1000)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, 0)
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
		cache.Set(key, i, 0)
	}
}

func BenchmarkCacheSetWithTTL(b *testing.B) {
	cache := NewLRUCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Set(key, i, 5*time.Minute)
	}
}

func BenchmarkCacheGetWithExpiredChecks(b *testing.B) {
	cache := NewLRUCache(1000)

	// Pre-populate with TTL
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, 5*time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Get(key)
	}
}

func BenchmarkCacheSetWithEvictionWithExpiredChecks(b *testing.B) {
	cache := NewLRUCache(100) // Small cache = constant eviction
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, i, 5*time.Minute)
	}
}

// Without TTL
// goos: darwin
// goarch: arm64
// pkg: cache
// cpu: Apple M1
// BenchmarkCacheSet-8               	 9204403	       119.6 ns/op	      21 B/op	       2 allocs/op
// BenchmarkCacheGet-8               	11508284	       104.1 ns/op	      13 B/op	       1 allocs/op
// BenchmarkCacheSetWithEviction-8   	 4714653	       254.7 ns/op	      79 B/op	       3 allocs/op
// PASS
// ok  	cache	4.630s

// WithTTL
// BenchmarkCacheSet-8{ttl=0} 			               	 		9381373	       121.5 ns/op	      21 B/op	       2 allocs/op
// BenchmarkCacheGet-8{ttl=0}                                	10910235	       109.8 ns/op	      13 B/op	       1 allocs/op
// BenchmarkCacheSetWithEviction-8{ttl=0}                    	 4373330	       261.6 ns/op	     111 B/op	       4 allocs/op
// BenchmarkCacheSetWithTTL-8{ttl=5min}                         	 4914612	       230.4 ns/op	      21 B/op	       2 allocs/op
// BenchmarkCacheGetWithExpiredChecks-8{ttl=5min}               	 5829934	       205.4 ns/op	      13 B/op	       1 allocs/op
// BenchmarkCacheSetWithEvictionWithExpiredChecks-8 {ttl=5min}  	 3327799	       360.7 ns/op	     111 B/op	       4 allocs/op
// PASS
// ok  	cache	9.187s

// time.Now() is a syscall! It's one of the slowest operations:
// Get current time from OS
// ~20-30ns on modern CPUs
// Called on EVERY Set
// Batch time checks:This is what Redis does! Cache the current time and reuse it for multiple checks.
// Tradeoff:
// Expiration timing less precise (±100ms)
// But 10x faster checks
