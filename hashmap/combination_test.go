package fastmap_test

import (
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestComboKeyHashMap(t *testing.T) {
	t.Run("basic key combination lookup", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Store value with key combination
		m.Put([]string{"key1", "key2"}, "TEST")

		// Test retrieving with exact keys
		result, exists := m.Get([]string{"key1", "key2"})
		if !exists || result != "TEST" {
			t.Errorf("Expected TEST, got %v", result)
		}

		// Test retrieving with reversed key order
		result, exists = m.Get([]string{"key2", "key1"})
		if !exists || result != "TEST" {
			t.Errorf("Expected TEST, got %v", result)
		}
	})

	t.Run("get by partial keys", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Store multiple combinations
		m.Put([]string{"key1", "key2"}, "TEST1")
		m.Put([]string{"key2", "key3"}, "TEST2")

		// Test getting by single key
		results := m.GetByKeys([]string{"key2"})
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		// Test getting by two keys
		results = m.GetByKeys([]string{"key1", "key2"})
		if len(results) != 1 || results[0] != "TEST1" {
			t.Errorf("Expected [TEST1], got %v", results)
		}
	})

	t.Run("overlapping key combinations", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		m.Put([]string{"key1", "key2"}, "TEST1")
		m.Put([]string{"key1", "key2", "key3"}, "TEST2")

		results := m.GetByKeys([]string{"key1", "key2"})
		if len(results) != 2 {
			t.Errorf("Expected 2 results for overlapping keys, got %d", len(results))
		}
	})

	t.Run("empty and non-existent keys", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Test empty keys
		results := m.GetByKeys([]string{})
		if len(results) != 0 {
			t.Error("Expected empty result for empty keys")
		}

		// Test non-existent keys
		results = m.GetByKeys([]string{"nonexistent"})
		if len(results) != 0 {
			t.Error("Expected empty result for non-existent keys")
		}
	})
}
