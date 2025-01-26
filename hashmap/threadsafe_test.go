package fastmap_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

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

func TestThreadSafeConcurrentOperations(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.Put(fmt.Sprintf("key%d", val), val)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.Get(fmt.Sprintf("key%d", val))
		}(i)
	}

	wg.Wait()

	if size := m.Size(); size != numGoroutines {
		t.Errorf("Size after concurrent operations: got %d, want %d", size, numGoroutines)
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

func TestThreadSafeFunctionalOperations(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// Test Filter
	filtered := m.Filter(func(k string, v int) bool {
		return v > 1
	})
	if filtered.Size() != 2 {
		t.Error("Filter failed: wrong number of elements")
	}

	// Test Map
	doubled := m.Map(func(k string, v int) int {
		return v * 2
	})
	if val, _ := doubled.Get("one"); val != 2 {
		t.Error("Map failed: wrong transformation")
	}

	// Test ForEach
	sum := 0
	m.ForEach(func(k string, v int) {
		sum += v
	})
	if sum != 6 {
		t.Errorf("ForEach failed: got sum %d, want 6", sum)
	}
}

func TestThreadSafeUpdateAndPutAll(t *testing.T) {
	m1 := fastmap.NewThreadSafeHashMap[string, int]()
	m2 := fastmap.NewThreadSafeHashMap[string, int]()

	// Test UpdateValue
	m1.Put("key1", 100)
	success := m1.UpdateValue("key1", 200)
	if !success {
		t.Error("UpdateValue failed to update existing key")
	}
	if val, _ := m1.Get("key1"); val != 200 {
		t.Error("UpdateValue: wrong value after update")
	}

	// Test PutAll
	m2.Put("key2", 300)
	m2.Put("key3", 400)
	m1.PutAll(m2)
	if m1.Size() != 3 {
		t.Error("PutAll failed: wrong size after merge")
	}
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

func TestThreadSafeDeadlockPrevention(t *testing.T) {
	m1 := fastmap.NewThreadSafeHashMap[string, int]()
	m2 := fastmap.NewThreadSafeHashMap[string, int]()

	done := make(chan bool)
	go func() {
		m1.Put("key", 1)
		m2.Put("key", 2)
		done <- true
	}()

	go func() {
		m2.Put("key", 3)
		m1.Put("key", 4)
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("Potential deadlock detected")
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

func TestThreadSafeHighConcurrency(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[int, int]()
	numCPU := runtime.NumCPU()
	numOps := 10000
	var wg sync.WaitGroup

	for cpu := 0; cpu < numCPU; cpu++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			for i := 0; i < numOps; i++ {
				key := threadID*numOps + i
				m.Put(key, i)
				m.Get(key)
				m.Remove(key)
			}
		}(cpu)
	}
	wg.Wait()
}

func TestThreadSafeEmptyOperations(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()

	if keys := m.Keys(); len(keys) != 0 {
		t.Error("Keys() should return empty slice for empty map")
	}

	if values := m.Values(); len(values) != 0 {
		t.Error("Values() should return empty slice for empty map")
	}

	m.ForEach(func(k string, v int) {
		t.Error("ForEach should not execute on empty map")
	})
}

func TestThreadSafeConcurrentFilterMap(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	var wg sync.WaitGroup

	// Fill map with unique keys
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	// Rest of the test remains the same
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			filtered := m.Filter(func(k string, v int) bool {
				return v%2 == 0
			})
			if filtered.Size() == 0 {
				t.Error("Filter returned empty result")
			}
		}()

		go func() {
			defer wg.Done()
			mapped := m.Map(func(k string, v int) int {
				return v * 2
			})
			if mapped.Size() == 0 {
				t.Error("Map returned empty result")
			}
		}()
	}
	wg.Wait()
}

func TestThreadSafeZeroValues(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[int, int]()

	m.Put(0, 0)
	if val, exists := m.Get(0); !exists || val != 0 {
		t.Error("Failed to handle zero values")
	}

	if !m.Contains(0) {
		t.Error("Contains failed for zero key")
	}

	filtered := m.Filter(func(k, v int) bool {
		return k == 0 && v == 0
	})
	if filtered.Size() != 1 {
		t.Error("Filter failed for zero values")
	}
}

func TestThreadSafeConcurrentUpdateValue(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	m.Put("key", 0)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.UpdateValue("key", val)
		}(i)
	}
	wg.Wait()

	if val, _ := m.Get("key"); val == 0 {
		t.Error("UpdateValue failed under concurrency")
	}
}

func TestThreadSafeMapRace(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, *int]()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(3)
		val := i
		go func() {
			defer wg.Done()
			m.Put("key", &val)
		}()
		go func() {
			defer wg.Done()
			if v, exists := m.Get("key"); exists {
				_ = *v // Dereference to check for race
			}
		}()
		go func() {
			defer wg.Done()
			m.Remove("key")
		}()
	}
	wg.Wait()
}

func TestThreadSafeHandleFieldConfigs(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	rowIndex := 0
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			RowIndex: &rowIndex,
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"value": 1},
		{"value": 2},
	}

	h.Put("field1", 0)
	result := h.HandleFieldConfigs(data, configs, "field1")

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}

	if result[0] != 1 || result[1] != 2 {
		t.Errorf("Expected values [1,2], got %v", result)
	}
}

func TestThreadSafeApplyFieldConfig(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	config := fastmap.FieldConfig[int]{
		Handler: func(data map[string]interface{}) *int {
			if val, ok := data["value"].(int); ok {
				return &val
			}
			return nil
		},
	}

	data := map[string]interface{}{"value": 42}

	h.Put("field1", 0)
	result := h.ApplyFieldConfig("field1", config, data)

	if !result {
		t.Error("Expected true, got false")
	}

	if val, _ := h.Get("field1"); val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

func TestThreadSafeProcessFieldConfigs(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	rowIndex := 0
	processedCount := 0
	var mu sync.Mutex

	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			RowIndex: &rowIndex,
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"value": 1},
		{"value": 2},
	}

	processor := func(key string, value int, index int) {
		mu.Lock()
		processedCount++
		mu.Unlock()
	}

	h.Put("field1", 0)
	h.ProcessFieldConfigs(configs, data, processor)

	if processedCount != 2 {
		t.Errorf("Expected 2 processed items, got %d", processedCount)
	}
}

func TestProcessFieldConfigs_ConcurrentAccess(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	processedCount := 0
	var mu sync.Mutex

	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]interface{}{"value": i}
	}

	processor := func(key string, value int, index int) {
		mu.Lock()
		processedCount++
		mu.Unlock()
	}

	h.Put("field1", 0)
	h.ProcessFieldConfigs(configs, data, processor)

	if processedCount != 100 {
		t.Errorf("Expected 100 processed items, got %d", processedCount)
	}
}
