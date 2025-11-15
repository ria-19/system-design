package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBasicOperations(t *testing.T) {
	cache := NewLRUCache(2)

	// Test Set and Get
	cache.Set("a", 1, 0) // use default TTL
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("Expected 1, got %v", val)
	}

	// Test missing key
	if _, ok := cache.Get("b"); ok {
		t.Error("Expected key 'b' to not exist")
	}
}

func TestLRUEviction(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", 1, 0)
	cache.Set("b", 2, 0)
	cache.Set("c", 3, 0) // Should evict "a"

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be evicted")
	}

	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("Expected 'b' to exist with value 2")
	}
}

func TestGetRefreshesRecency(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", 1, 0)
	cache.Set("b", 2, 0)
	cache.Get("a")       // "a" is now most recent
	cache.Set("c", 3, 0) // Should evict "b", not "a"

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Error("Expected 'a' to still exist")
	}

	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to be evicted")
	}
}

func TestUpdateRefreshesRecency(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", 1, 0)
	cache.Set("b", 2, 0)
	cache.Set("a", 3, 0) // Updates value of a AND refreshes recency
	cache.Set("c", 4, 0) // Should evict "b", not "a"

	if val, ok := cache.Get("a"); !ok || val != 3 {
		t.Error("Expected 'a' to still exist")
	}

	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to be evicted")
	}
}

func TestTTLExpiration(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key", "value", 100*time.Millisecond)

	// Should exist immediately
	if _, ok := cache.Get("key"); !ok {
		t.Error("Key should exist")
	}

	// Should expire after TTL
	time.Sleep(150 * time.Millisecond)
	if _, ok := cache.Get("key"); ok {
		t.Error("Key should be expired")
	}
}

func TestTTLRenewal(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("key", "value", 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	// Should update TTL (irrespective of expiration) and value
	cache.Set("key", "value_new", 100*time.Millisecond)

	// key should exist
	if _, ok := cache.Get("key"); !ok {
		t.Error("Key should exist")
	}
}
func TestExpiredNodesCountTowardCapacity(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", 1, 50*time.Millisecond)
	cache.Set("b", 2, 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond) // Both expired

	// Cache is "full" but both entries expired
	cache.Set("c", 3, 0)
	cache.Set("d", 4, 0) // This should work but might evict "b"

	// Both c and d should exist
	if _, ok := cache.Get("c"); !ok {
		t.Error("c should exist")
	}
	if _, ok := cache.Get("d"); !ok {
		t.Error("d should exist")
	}
}

// Simulates a busy web server where multiple requests (goroutines) are trying to read and write
// to the same shared cache at the exact same time.
func TestConcurrentAccess(t *testing.T) {
	cache := NewLRUCache(100)

	// Launch 10 goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				cache.Set(key, j, 0)
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
}

// Concurrency stress test
func TestHighContention(t *testing.T) {
	cache := NewLRUCache(100)

	// 100 goroutines hammering the SAME keys
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Everyone fights over keys 0-9
				key := fmt.Sprintf("key%d", j%10)
				cache.Set(key, j, 0)
				cache.Get(key)
			}
		}()
	}
	wg.Wait()
}

/* Before Thread Safety
 ==================
 WARNING: DATA RACE
Write at 0x00c00040c028 by goroutine 17:
  cache.(*LRUCache).addToTail()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache.go:110
  cache.(*LRUCache).Set()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache.go:88
  cache.TestConcurrentAccess.func1()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache_test.go:138

Previous read at 0x00c00040c028 by goroutine 9:
  cache.(*LRUCache).removeNode()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache.go:98
  cache.(*LRUCache).moveToTail()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache.go:116
  cache.(*LRUCache).Get()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache.go:53
  cache.TestConcurrentAccess.func1()
      /Users/riya/Desktop/dev/System-Design/week1-cache/cache_test.go:139
 ==================
*/
