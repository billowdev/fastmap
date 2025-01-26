package fastmap_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

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
func TestThreadSafeEmptyOperations(t *testing.T) {
	m := fastmap.NewThreadSafeHashMap[string, int]()

	if keys := m.Keys(); len(keys) != 0 {
		t.Error("Keys() should return empty slice for empty map")
	}

	if values := m.Values(); len(values) != 0 {
		t.Error("Values() should return empty slice for empty map")
	}

	err := m.ForEach(func(k string, v int) error {
		t.Error("ForEach should not execute on empty map")
		return nil
	})
	if err != nil {
		t.Errorf("ForEach on empty map returned error: %v", err)
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
	err := m.ForEach(func(k string, v int) error {
		sum += v
		return nil
	})
	if err != nil {
		t.Errorf("ForEach failed with error: %v", err)
	}
	if sum != 6 {
		t.Errorf("ForEach failed: got sum %d, want 6", sum)
	}
}
