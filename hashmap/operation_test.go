package fastmap_test

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestClear(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("key1", 1)
	h.Put("key2", 2)
	h.Clear()
	if !h.IsEmpty() {
		t.Error("Clear failed, map not empty")
	}
}

func TestContains(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("key", 100)
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"existing key", "key", true},
		{"non-existing key", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.Contains(tt.key); got != tt.want {
				t.Errorf("Contains(%v) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestUpdateValue(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("key", 100)

	tests := []struct {
		name    string
		key     string
		value   int
		want    bool
		wantVal int
	}{
		{"existing key", "key", 200, true, 200},
		{"non-existing key", "nonexistent", 300, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.UpdateValue(tt.key, tt.value); got != tt.want {
				t.Errorf("UpdateValue(%v, %v) = %v, want %v", tt.key, tt.value, got, tt.want)
			}
			if val, _ := h.Get(tt.key); val != tt.wantVal && tt.want {
				t.Errorf("After update, Get(%v) = %v, want %v", tt.key, val, tt.wantVal)
			}
		})
	}
}

func TestClearWithConcurrentOperations(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	wg := sync.WaitGroup{}

	// Add items concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Add(-1)
			h.Put(fmt.Sprintf("key%d", val), val)
		}(i)
	}
	wg.Wait()

	h.Clear()
	if !h.IsEmpty() {
		t.Error("Clear failed under concurrent operations")
	}
}

func TestContainsWithNilValues(t *testing.T) {
	h := fastmap.NewHashMap[string, *int]()
	var nilValue *int
	h.Put("nilKey", nilValue)

	if !h.Contains("nilKey") {
		t.Error("Contains should return true for keys with nil values")
	}
}

func TestKeysWithDuplicateValues(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("key1", 100)
	h.Put("key2", 100)
	h.Put("key3", 100)

	keys := h.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys for duplicate values, got %d", len(keys))
	}
}

func TestValuesOrder(t *testing.T) {
	h := fastmap.NewHashMap[int, string]()
	expected := []string{"first", "second", "third"}

	for i, v := range expected {
		h.Put(i, v)
	}

	values := h.Values()
	if len(values) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(values))
	}
}
func TestForEachWithModification(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("key1", 1)
	h.Put("key2", 2)

	// Create copy of map to iterate safely
	mapCopy := make(map[string]int)
	h.ForEach(func(k string, v int) {
		mapCopy[k] = v
	})

	// Modify using the copy
	for k, v := range mapCopy {
		h.Put(k+"_new", v*2)
	}

	if h.Size() != 4 {
		t.Errorf("Expected size 4 after ForEach modification, got %d", h.Size())
	}

	// Verify the content
	expected := map[string]int{
		"key1":     1,
		"key2":     2,
		"key1_new": 2,
		"key2_new": 4,
	}

	for k, v := range expected {
		if val, exists := h.Get(k); !exists || val != v {
			t.Errorf("Value mismatch for key %s: got %d, want %d", k, val, v)
		}
	}
}

func TestUpdateValueEdgeCases(t *testing.T) {
	h := fastmap.NewHashMap[string, interface{}]()

	tests := []struct {
		name    string
		key     string
		value   interface{}
		setup   func()
		want    bool
		wantVal interface{}
	}{
		{
			name:    "update nil to value",
			key:     "key1",
			value:   100,
			setup:   func() { h.Put("key1", nil) },
			want:    true,
			wantVal: 100,
		},
		{
			name:    "update value to nil",
			key:     "key2",
			value:   nil,
			setup:   func() { h.Put("key2", 200) },
			want:    true,
			wantVal: nil,
		},
		{
			name:    "update non-existent key",
			key:     "nonexistent",
			value:   300,
			setup:   func() {},
			want:    false,
			wantVal: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Clear()
			tt.setup()

			got := h.UpdateValue(tt.key, tt.value)
			if got != tt.want {
				t.Errorf("UpdateValue() = %v, want %v", got, tt.want)
			}

			if val, _ := h.Get(tt.key); val != tt.wantVal && tt.want {
				t.Errorf("After update, value = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

func TestPutAllWithEmptyMaps(t *testing.T) {
	h1 := fastmap.NewHashMap[string, int]()
	h2 := fastmap.NewHashMap[string, int]()

	h1.PutAll(h2)
	if !h1.IsEmpty() {
		t.Error("PutAll with empty source should maintain empty destination")
	}

	h2.Put("key", 100)
	h1.PutAll(h2)
	if h1.Size() != 1 {
		t.Error("PutAll failed to copy from non-empty to empty map")
	}
}

func TestConcurrentOperations(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	wg := sync.WaitGroup{}
	numGoroutines := 100

	// Concurrent puts
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Add(-1)
			h.Put(fmt.Sprintf("key%d", val), val)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			h.Keys()
			h.Values()
		}()
	}

	// Concurrent updates
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Add(-1)
			h.UpdateValue(fmt.Sprintf("key%d", val), val*2)
		}(i)
	}

	wg.Wait()

	if h.Size() != numGoroutines {
		t.Errorf("Expected size %d after concurrent operations, got %d", numGoroutines, h.Size())
	}
}

func TestLargeDataSetOperations(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	numItems := 10000

	// Test large-scale insertions
	for i := 0; i < numItems; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}

	if h.Size() != numItems {
		t.Errorf("Expected size %d for large dataset, got %d", numItems, h.Size())
	}

	// Test ForEach on large dataset
	count := 0
	h.ForEach(func(k string, v int) {
		count++
	})

	if count != numItems {
		t.Errorf("ForEach processed %d items, expected %d", count, numItems)
	}

	// Test Keys/Values on large dataset
	keys := h.Keys()
	values := h.Values()

	if len(keys) != numItems || len(values) != numItems {
		t.Errorf("Keys/Values length mismatch: keys=%d, values=%d, expected=%d",
			len(keys), len(values), numItems)
	}
}

func TestTypeEdgeCases(t *testing.T) {
	// Test with complex types
	type complexKey struct {
		id   int
		name string
	}

	h := fastmap.NewHashMap[complexKey, interface{}]()
	key1 := complexKey{1, "test1"}
	key2 := complexKey{1, "test2"}

	h.Put(key1, "value1")
	h.Put(key2, "value2")

	if h.Size() != 2 {
		t.Error("Failed to handle complex key types correctly")
	}

	// Test with function values
	funcMap := fastmap.NewHashMap[string, func() int]()
	f1 := func() int { return 1 }
	f2 := func() int { return 2 }

	funcMap.Put("func1", f1)
	funcMap.Put("func2", f2)

	if funcMap.Size() != 2 {
		t.Error("Failed to handle function value types")
	}
}

func TestEdgeCaseEmptyMaps(t *testing.T) {
	h := fastmap.NewHashMap[string, struct{}]()
	h.Put("key", struct{}{})
	h.Remove("key")

	// Operations on empty map
	if !h.IsEmpty() {
		t.Error("IsEmpty failed after removing last element")
	}
	h.ForEach(func(k string, v struct{}) {
		t.Error("ForEach should not execute on empty map")
	})
	if h.Filter(func(k string, v struct{}) bool { return true }).Size() != 0 {
		t.Error("Filter on empty map should return empty map")
	}
}
func TestEdgeCaseKeyTypes(t *testing.T) {
	type complexKey struct {
		f float64
		s string
	}

	h := fastmap.NewHashMap[complexKey, int]()

	// Test only Inf values since NaN != NaN in Go
	k1 := complexKey{f: math.Inf(1), s: "inf"}
	k2 := complexKey{f: math.Inf(-1), s: "neginf"}
	k3 := complexKey{f: 0.0, s: "zero"}

	h.Put(k1, 1)
	h.Put(k2, 2)
	h.Put(k3, 3)

	if !h.Contains(k1) || !h.Contains(k2) || !h.Contains(k3) {
		t.Error("Failed to handle special float values in struct keys")
	}

	// Verify values
	if val, exists := h.Get(k1); !exists || val != 1 {
		t.Error("Failed to retrieve Inf value")
	}
	if val, exists := h.Get(k2); !exists || val != 2 {
		t.Error("Failed to retrieve -Inf value")
	}
	if val, exists := h.Get(k3); !exists || val != 3 {
		t.Error("Failed to retrieve zero value")
	}
}

func TestEdgeCaseNestedTypes(t *testing.T) {
	type nested struct {
		m *fastmap.HashMap[int, *fastmap.HashMap[string, []int]]
	}

	h := fastmap.NewHashMap[string, nested]()
	innerMap := fastmap.NewHashMap[string, []int]()
	innerMap.Put("nums", []int{1, 2, 3})

	midMap := fastmap.NewHashMap[int, *fastmap.HashMap[string, []int]]()
	midMap.Put(1, innerMap)

	h.Put("nested", nested{m: midMap})

	if val, exists := h.Get("nested"); !exists || val.m == nil {
		t.Error("Failed to handle deeply nested types")
	}
}

func TestEdgeCaseMemoryPressure(t *testing.T) {
	h := fastmap.NewHashMap[int, []byte]()
	const itemSize = 1 << 10 // 1KB
	const itemCount = 1000

	// Initial garbage collection
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	beforeAlloc := m.Alloc

	// Add items and keep references to prevent GC
	items := make([][]byte, itemCount)
	for i := 0; i < itemCount; i++ {
		items[i] = make([]byte, itemSize)
		h.Put(i, items[i])
	}

	runtime.ReadMemStats(&m)
	afterAlloc := m.Alloc
	if afterAlloc <= beforeAlloc {
		t.Error("No memory allocation detected")
	}

	// Clear and allow GC
	h.Clear()
	items = nil
	runtime.GC()
	runtime.GC() // Multiple GC calls to ensure cleanup

	runtime.ReadMemStats(&m)
	finalAlloc := m.Alloc
	if finalAlloc >= afterAlloc {
		t.Error("Memory not properly released after Clear")
	}
}

func TestEdgeCaseDataRaces(t *testing.T) {
	const mapCount = 100
	maps := make([]*fastmap.HashMap[int, int], mapCount)

	for i := range maps {
		maps[i] = fastmap.NewHashMap[int, int]()
		maps[i].Put(i, i)
	}

	for i := 0; i < mapCount-1; i++ {
		maps[i].PutAll(maps[i+1])
		if maps[i].Size() != 2 {
			t.Errorf("Map %d: expected size 2, got %d", i, maps[i].Size())
		}
	}
}

func TestEdgeCaseValueModification(t *testing.T) {
	type mutable struct {
		value int
	}

	h := fastmap.NewHashMap[string, *mutable]()
	m := &mutable{value: 1}
	h.Put("key", m)

	// Modify value through map reference
	if val, exists := h.Get("key"); exists {
		val.value = 2
	}

	// Verify modification persists
	if val, exists := h.Get("key"); !exists || val.value != 2 {
		t.Error("Value modification through map reference failed")
	}
}

func TestEdgeCaseZeroValues(t *testing.T) {
	h := fastmap.NewHashMap[int, int]()

	// Test zero key
	h.Put(0, 1)
	if val, exists := h.Get(0); !exists || val != 1 {
		t.Error("Failed to handle zero key")
	}

	// Test zero value
	h.Put(1, 0)
	if val, exists := h.Get(1); !exists || val != 0 {
		t.Error("Failed to handle zero value")
	}

	// Test zero key with zero value
	h.Put(0, 0)
	if val, exists := h.Get(0); !exists || val != 0 {
		t.Error("Failed to handle zero key with zero value")
	}
}

func TestEdgeCaseMapOperations(t *testing.T) {
	h := fastmap.NewHashMap[string, *fastmap.HashMap[string, int]]()

	// Put a map as value
	innerMap := fastmap.NewHashMap[string, int]()
	innerMap.Put("inner", 1)
	h.Put("outer", innerMap)

	// Modify inner map
	if val, exists := h.Get("outer"); exists {
		val.Put("inner2", 2)
	}

	// Verify modifications
	if val, exists := h.Get("outer"); !exists || val.Size() != 2 {
		t.Error("Failed to handle map as value")
	}
}
