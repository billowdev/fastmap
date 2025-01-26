// hashmap_test.go
package fastmap

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewHashMap(t *testing.T) {
	h := NewHashMap[string, int]()
	if h == nil {
		t.Error("NewHashMap returned nil")
	}
	if h.Size() != 0 {
		t.Errorf("Expected empty map, got size %d", h.Size())
	}
}

func TestPut(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("key", 100)
	if val, exists := h.Get("key"); !exists || val != 100 {
		t.Errorf("Put failed, got (%v, %v), want (100, true)", val, exists)
	}
}

func TestGet(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("key", 100)

	tests := []struct {
		name       string
		key        string
		wantValue  int
		wantExists bool
	}{
		{"existing key", "key", 100, true},
		{"non-existing key", "nonexistent", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, exists := h.Get(tt.key)
			if exists != tt.wantExists || value != tt.wantValue {
				t.Errorf("Get(%v) = (%v, %v), want (%v, %v)",
					tt.key, value, exists, tt.wantValue, tt.wantExists)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("key", 100)
	h.Remove("key")
	if _, exists := h.Get("key"); exists {
		t.Error("Remove failed, key still exists")
	}
}

func TestSize(t *testing.T) {
	h := NewHashMap[string, int]()
	sizes := []struct {
		operation string
		change    func()
		want      int
	}{
		{"initial", func() {}, 0},
		{"after put", func() { h.Put("key1", 1) }, 1},
		{"after another put", func() { h.Put("key2", 2) }, 2},
		{"after remove", func() { h.Remove("key1") }, 1},
	}

	for _, tt := range sizes {
		t.Run(tt.operation, func(t *testing.T) {
			tt.change()
			if got := h.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestEdgeCases(t *testing.T) {
	h := NewThreadSafeHashMap[string, int]() // Use thread-safe version

	t.Run("zero value operations", func(t *testing.T) {
		h.Put("", 0)
		if val, exists := h.Get(""); !exists || val != 0 {
			t.Error("Failed to handle empty string key")
		}
	})

	t.Run("overwrites", func(t *testing.T) {
		h.Put("key", 1)
		h.Put("key", 2)
		if val, _ := h.Get("key"); val != 2 {
			t.Error("Failed to overwrite value")
		}
	})

	t.Run("large number of items", func(t *testing.T) {
		h := NewThreadSafeHashMap[string, int]() // Fresh map
		for i := 0; i < 1000; i++ {
			h.Put(fmt.Sprintf("key%d", i), i)
		}
		time.Sleep(10 * time.Millisecond) // Allow operations to complete
		if h.Size() != 1000 {
			t.Errorf("Expected size 1000, got %d", h.Size())
		}
	})

	t.Run("concurrent modifications", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := fmt.Sprintf("concurrent%d", i)
				h.Put(key, i)
				h.Get(key)
				h.Remove(key)
			}(i)
		}
		wg.Wait()
	})

	t.Run("memory leak check", func(t *testing.T) {
		h := NewThreadSafeHashMap[string, int]() // Fresh map
		for i := 0; i < 100; i++ {
			h.Put(fmt.Sprintf("leak%d", i), i)
			h.Remove(fmt.Sprintf("leak%d", i))
		}
		if h.Size() != 0 {
			t.Error("Possible memory leak detected")
		}
	})
}
