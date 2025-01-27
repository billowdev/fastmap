package fastmap_test

import (
	"fmt"
	"sync"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestThreadSafeMultiKeyMap_ConcurrentOperations(t *testing.T) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	var wg sync.WaitGroup
	numOps := 100

	// Concurrent writes
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			keys := []string{
				fmt.Sprintf("primary%d", val),
				fmt.Sprintf("alias1_%d", val),
				fmt.Sprintf("alias2_%d", val),
			}
			m.Put(keys, val)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.Get(fmt.Sprintf("primary%d", val))
			m.Get(fmt.Sprintf("alias1_%d", val))
		}(i)
	}

	// Concurrent alias additions
	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			m.AddAlias(fmt.Sprintf("primary%d", val), fmt.Sprintf("new_alias_%d", val))
		}(i)
	}

	wg.Wait()

	// Verify state
	if size := m.Size(); size != numOps {
		t.Errorf("Expected size %d, got %d", numOps, size)
	}
}

func TestThreadSafeMultiKeyMap_ConcurrentRemoval(t *testing.T) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	var wg sync.WaitGroup
	numOps := 100

	// Setup initial data
	for i := 0; i < numOps; i++ {
		keys := []string{
			fmt.Sprintf("primary%d", i),
			fmt.Sprintf("alias1_%d", i),
			fmt.Sprintf("alias2_%d", i),
		}
		m.Put(keys, i)
	}

	// Concurrent removals and reads
	for i := 0; i < numOps; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			m.RemoveWithCascade(fmt.Sprintf("primary%d", val))
		}(i)
		go func(val int) {
			defer wg.Done()
			m.GetAllKeys(fmt.Sprintf("primary%d", val))
		}(i)
	}

	wg.Wait()

	if size := m.Size(); size != 0 {
		t.Errorf("Expected size 0 after removals, got %d", size)
	}
}

func TestThreadSafeMultiKeyMap_ConcurrentAliasOperations(t *testing.T) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	var wg sync.WaitGroup
	numOps := 100

	// Setup initial data
	primaryKey := "primary"
	m.Put([]string{primaryKey}, 42)

	// Concurrent alias additions and removals
	for i := 0; i < numOps; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			alias := fmt.Sprintf("alias_%d", val)
			m.AddAlias(primaryKey, alias)
		}(i)
		go func(val int) {
			defer wg.Done()
			alias := fmt.Sprintf("alias_%d", val)
			m.Remove(alias)
		}(i)
	}

	wg.Wait()

	// Primary key should still exist
	if val, exists := m.Get(primaryKey); !exists || val != 42 {
		t.Error("Primary key should still exist with original value")
	}
}

func TestThreadSafeMultiKeyMap_ConcurrentClear(t *testing.T) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	var wg sync.WaitGroup
	numOps := 100

	// Concurrent writes and clears
	for i := 0; i < numOps; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			keys := []string{fmt.Sprintf("primary%d", val)}
			m.Put(keys, val)
		}(i)
		go func() {
			defer wg.Done()
			m.Clear()
		}()
	}

	wg.Wait()

	// Final clear to ensure consistent state
	m.Clear()
	if size := m.Size(); size != 0 {
		t.Errorf("Expected size 0 after clear, got %d", size)
	}
}