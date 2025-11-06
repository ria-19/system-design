package cache

import "testing"

func TestBasicOperations(t *testing.T) {
	cache := NewLRUCache(2)

	// Test Set and Get
	cache.Set("a", 1)
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

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3) // Should evict "a"

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be evicted")
	}

	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("Expected 'b' to exist with value 2")
	}
}

func TestGetRefreshesRecency(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Get("a")    // "a" is now most recent
	cache.Set("c", 3) // Should evict "b", not "a"

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Error("Expected 'a' to still exist")
	}

	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to be evicted")
	}
}

func TestUpdateRefreshesRecency(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("a", 3) // Updates value of a AND refreshes recency
	cache.Set("c", 4) // Should evict "b", not "a"

	if val, ok := cache.Get("a"); !ok || val != 3 {
		t.Error("Expected 'a' to still exist")
	}

	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to be evicted")
	}
}
