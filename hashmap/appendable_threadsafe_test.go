package fastmap_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestThreadSafeAppendableHashMap_Basic(t *testing.T) {
	t.Run("initialization", func(t *testing.T) {
		m := fastmap.NewThreadSafeAppendableHashMap[string, int]()
		if m == nil {
			t.Error("NewThreadSafeAppendableHashMap returned nil")
		}
	})

	t.Run("append to new key", func(t *testing.T) {
		m := fastmap.NewThreadSafeAppendableHashMap[string, int]()
		m.AppendValues("key1", 1, 2, 3)

		values, exists := m.Get("key1")
		if !exists {
			t.Fatal("Key should exist after AppendValues")
		}

		expected := []int{1, 2, 3}
		if len(values) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(values))
		}

		for i, v := range expected {
			if values[i] != v {
				t.Errorf("At index %d: expected %d, got %d", i, v, values[i])
			}
		}
	})

	t.Run("append to existing key", func(t *testing.T) {
		m := fastmap.NewThreadSafeAppendableHashMap[string, int]()
		m.AppendValues("key1", 1, 2)
		m.AppendValues("key1", 3, 4)

		values, exists := m.Get("key1")
		if !exists {
			t.Fatal("Key should exist after AppendValues")
		}

		expected := []int{1, 2, 3, 4}
		if len(values) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(values))
		}

		for i, v := range expected {
			if values[i] != v {
				t.Errorf("At index %d: expected %d, got %d", i, v, values[i])
			}
		}
	})
}

func TestThreadSafeAppendableHashMap_Concurrent(t *testing.T) {
	t.Run("concurrent appends to same key", func(t *testing.T) {
		m := fastmap.NewThreadSafeAppendableHashMap[string, int]()
		var wg sync.WaitGroup
		iterations := 100
		goroutines := 10

		for g := 0; g < goroutines; g++ {
			wg.Add(1)
			go func(routine int) {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					m.AppendValues("shared", routine*iterations+i)
				}
			}(g)
		}

		wg.Wait()

		values, exists := m.Get("shared")
		if !exists {
			t.Fatal("Key should exist after concurrent AppendValues")
		}

		expectedLength := goroutines * iterations
		if len(values) != expectedLength {
			t.Errorf("Expected length %d, got %d", expectedLength, len(values))
		}

		// Verify no values were lost (using map for efficient lookup)
		seen := make(map[int]bool)
		for _, v := range values {
			seen[v] = true
		}

		for g := 0; g < goroutines; g++ {
			for i := 0; i < iterations; i++ {
				expected := g*iterations + i
				if !seen[expected] {
					t.Errorf("Missing value: %d", expected)
				}
			}
		}
	})

	t.Run("concurrent appends to different keys", func(t *testing.T) {
		m := fastmap.NewThreadSafeAppendableHashMap[string, int]()
		var wg sync.WaitGroup
		goroutines := 10

		for g := 0; g < goroutines; g++ {
			wg.Add(1)
			go func(routine int) {
				defer wg.Done()
				key := fmt.Sprintf("key%d", routine)
				m.AppendValues(key, routine, routine+1, routine+2)
			}(g)
		}

		wg.Wait()

		for g := 0; g < goroutines; g++ {
			key := fmt.Sprintf("key%d", g)
			values, exists := m.Get(key)
			if !exists {
				t.Errorf("Key %s should exist", key)
				continue
			}

			expected := []int{g, g + 1, g + 2}
			if len(values) != len(expected) {
				t.Errorf("For key %s: expected length %d, got %d", key, len(expected), len(values))
			}

			for i, v := range expected {
				if values[i] != v {
					t.Errorf("For key %s at index %d: expected %d, got %d", key, i, v, values[i])
				}
			}
		}
	})

	t.Run("concurrent read-write operations", func(t *testing.T) {
		m := fastmap.NewThreadSafeAppendableHashMap[string, int]()
		var wg sync.WaitGroup
		done := make(chan bool)
		iterations := 1000

		// Writer goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				m.AppendValues("key", i)
				time.Sleep(time.Microsecond) // Simulate work
			}
		}()

		// Reader goroutines
		for r := 0; r < 5; r++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-done:
						return
					default:
						values, _ := m.Get("key")
						if len(values) > 0 {
							// Verify slice integrity
							for i := 1; i < len(values); i++ {
								if values[i] < values[i-1] {
									t.Errorf("Slice ordering violation: index %d (%d) < index %d (%d)",
										i, values[i], i-1, values[i-1])
								}
							}
						}
						time.Sleep(time.Microsecond) // Simulate work
					}
				}
			}()
		}

		// First close the done channel to signal readers to stop
		time.Sleep(time.Millisecond * 100) // Allow some time for concurrent operations
		close(done)

		// Then wait for all goroutines to finish
		wg.Wait()

		values, exists := m.Get("key")
		if !exists {
			t.Fatal("Key should exist after operations")
		}

		if len(values) != iterations {
			t.Errorf("Expected length %d, got %d", iterations, len(values))
		}
	})
}
