package fastmap_test

import (
	"fmt"
	"sync"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestThreadSafeBasicOperations(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()

	// Test Put and Get
	m.Put("key1", 100)
	if val, exists := m.Get("key1"); !exists || val != 100 {
		t.Errorf("Put/Get failed: got (%v, %v), want (100, true)", val, exists)
	}

	// Test Remove
	m.Remove("key1")
	if _, exists := m.Get("key1"); exists {
		t.Error("Remove failed: key still exists")
	}

	// Test Size
	m.Put("key2", 200)
	if size := m.Size(); size != 1 {
		t.Errorf("Size incorrect: got %d, want 1", size)
	}
}

func TestThreadSafeNilValueHandling(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, *string]()
	var nilStr *string

	m.Put("nilKey", nilStr)
	if val, exists := m.Get("nilKey"); !exists || val != nil {
		t.Error("Failed to handle nil value")
	}

	str := "value"
	m.Put("nilKey", &str)
	if val, exists := m.Get("nilKey"); !exists || val == nil {
		t.Error("Failed to update from nil to non-nil")
	}
}

func TestThreadSafeConcurrentMapClear(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	var wg sync.WaitGroup

	// Concurrent writes while clearing
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			m.Put("key", val)
		}(i)
		go func() {
			defer wg.Done()
			m.Clear()
		}()
	}
	wg.Wait()
}

func BenchmarkThreadSafeConcurrentOperations(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	var wg sync.WaitGroup

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				key := fmt.Sprintf("key%d", val)
				m.Put(key, val)
				m.Get(key)
				m.Contains(key)
			}(i)
			i++
		}
	})

	wg.Wait()
}

func TestThreadSafeConversion(t *testing.T) {
	original := map[string]int{"one": 1, "two": 2}

	// Test FromThreadSafeMap
	m := fastmap.FromThreadSafeMap(original)
	if m.Size() != len(original) {
		t.Error("FromThreadSafeMap: size mismatch")
	}

	// Test ToMap
	converted := m.ToMap()
	if len(converted) != len(original) {
		t.Error("ToMap: size mismatch")
	}
	for k, v := range original {
		if cv, exists := converted[k]; !exists || cv != v {
			t.Errorf("ToMap: value mismatch for key %s", k)
		}
	}
}

func TestThreadSafeMapOperations(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()

	// Test Clear
	m.Put("key1", 100)
	m.Clear()
	if !m.IsEmpty() {
		t.Error("Clear failed: map not empty")
	}

	// Test Contains
	m.Put("key2", 200)
	if !m.Contains("key2") {
		t.Error("Contains failed: key not found")
	}

	// Test Keys and Values
	m.Put("key3", 300)
	keys := m.Keys()
	values := m.Values()
	if len(keys) != 2 || len(values) != 2 {
		t.Error("Keys/Values length mismatch")
	}
}
