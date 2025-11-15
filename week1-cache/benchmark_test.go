package cache

import (
	"fmt"
	"testing"
	"time"
)

// This simulates a "warm" cache where most writes are updates.
func BenchmarkCacheSet(b *testing.B) {
	cache := NewLRUCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Set(key, i, 0)
	}
}

// This tests the pure "cache hit" scenario.
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

// This measures the performance of our worst-case write path.
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

func BenchmarkConcurrentGet(b *testing.B) {
	cache := NewLRUCache(1000)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, 0)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			cache.Get(key)
			i++
		}
	})
}

// Measures throughput of the cache under high contention
func BenchmarkConcurrentSet(b *testing.B) {
	cache := NewLRUCache(1000)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			cache.Set(key, i, 0)
			i++
		}
	})
}

// Low contention (1000 different keys)
func BenchmarkLowContention(b *testing.B) {
	cache := NewLRUCache(1000)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i) // Different key each time
			cache.Set(key, i, 0)
			i++
		}
	})
}

// High contention (only 10 keys)
func BenchmarkHighContention(b *testing.B) {
	cache := NewLRUCache(1000)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%10) // Only 10 keys!
			cache.Set(key, i, 0)
			i++
		}
	})
}

/*
------------------------------------------------------------
Benchmark Results (Without TTL)
------------------------------------------------------------
goos: darwin
goarch: arm64
pkg: cache
cpu: Apple M1

BenchmarkCacheSet-8                      9,204,403   119.6 ns/op    21 B/op   2 allocs/op
BenchmarkCacheGet-8                     11,508,284   104.1 ns/op    13 B/op   1 allocs/op
BenchmarkCacheSetWithEviction-8          4,714,653   254.7 ns/op    79 B/op   3 allocs/op

PASS
ok    cache    4.630s

------------------------------------------------------------
Benchmark Results (With TTL)
------------------------------------------------------------
BenchmarkCacheSet-8{ttl=0}                               9,381,373   121.5 ns/op    21 B/op    2 allocs/op
BenchmarkCacheGet-8{ttl=0}                              10,910,235   109.8 ns/op    13 B/op    1 allocs/op
BenchmarkCacheSetWithEviction-8{ttl=0}                   4,373,330   261.6 ns/op   111 B/op    4 allocs/op

BenchmarkCacheSetWithTTL-8{ttl=5min}                     4,914,612   230.4 ns/op    21 B/op    2 allocs/op
BenchmarkCacheGetWithExpiredChecks-8{ttl=5min}           5,829,934   205.4 ns/op    13 B/op    1 allocs/op
BenchmarkCacheSetWithEvictionWithExpiredChecks-8
  {ttl=5min}                                             3,327,799   360.7 ns/op   111 B/op    4 allocs/op

PASS
ok    cache    9.187s
Throughput = 1 operation / 205.4 ns ≈ 4.87 Million ops/sec on a single core.

------------------------------------------------------------
Notes on time.Now() Performance
------------------------------------------------------------
time.Now() triggers a syscall, which is one of the slowest operations.
Retrieving the current time requires a context switch: the program pauses,
the OS kernel provides the time, and control returns. This is *orders of
magnitude* slower than a memory read.

• Typical cost: ~20–30 ns on modern CPUs
• In this implementation, time.Now() is called on EVERY Set operation

Optimization idea (used by Redis):
  - Batch time checks by caching the current time and reusing it for multiple operations.
  - Tradeoff: expiration becomes slightly less precise (±100ms)
  - Benefit: up to 10× faster expiration checks

------------------------------------------------------------
Benchmark Results (Lock contention)
------------------------------------------------------------
- RWMutex
BenchmarkConcurrentGet-8   	 3803226	       349.0 ns/op	      13 B/op	       1 allocs/op
BenchmarkConcurrentSet-8   	 3788625	       338.6 ns/op	      22 B/op	       2 allocs/op
PASS
ok  	cache	4.091s
Total Throughput ≈ 8 cores * (1 operation / 349.0 ns) ≈ 22.9 Million ops/sec across the whole system.

-Mutex
BenchmarkConcurrentGet-8   	 4390030	       287.8 ns/op	      13 B/op	       1 allocs/op
BenchmarkConcurrentSet-8   	 4170391	       283.8 ns/op	      22 B/op	       2 allocs/op
PASS
ok  	cache	3.890s
Total Throughput ≈ 8 cores * (1 operation / 287.8 ns) ≈ 27.79 Million ops/sec across the whole system.

BenchmarkLowContention-8    	 2366044	       538.5 ns/op	      89 B/op	       3 allocs/op
BenchmarkHighContention-8   	 4403100	       328.5 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	cache	4.360s

------------------------------------------------------------
CPU Profile
------------------------------------------------------------
Commands
Generate: go test -bench=ConcurrentGet -cpuprofile=cpu.prof
BenchmarkConcurrentGet-8   	 4167906	       280.3 ns/op
PASS
ok  	cache	2.470s

View: go tool pprof cpu.prof



------------------------------------------------------------
Memory Profile
------------------------------------------------------------
Commands
Generate: go test -bench=ConcurrentGet -memprofile=mem.prof
View: go tool pprof mem.prof
*/
