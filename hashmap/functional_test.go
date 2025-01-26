// functional_test.go
package fastmap

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestFilter(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("one", 1)
	h.Put("two", 2)
	h.Put("three", 3)

	filtered := h.Filter(func(k string, v int) bool {
		return v%2 == 0
	})

	want := map[string]int{"two": 2}
	if !reflect.DeepEqual(filtered.ToMap(), want) {
		t.Errorf("Filter() = %v, want %v", filtered.ToMap(), want)
	}
}

func TestMap(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("one", 1)
	h.Put("two", 2)

	doubled := h.Map(func(k string, v int) int {
		return v * 2
	})

	want := map[string]int{"one": 2, "two": 4}
	if !reflect.DeepEqual(doubled.ToMap(), want) {
		t.Errorf("Map() = %v, want %v", doubled.ToMap(), want)
	}
}

func TestForEach(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("one", 1)
	h.Put("two", 2)

	sum := 0
	err := h.ForEach(func(k string, v int) error {
		sum += v
		return nil
	})
	if err != nil {
		t.Errorf("ForEach failed: %v", err)
	}

	if sum != 3 {
		t.Errorf("ForEach sum = %v, want 3", sum)
	}
}

func TestFunctionalOperations(t *testing.T) {
	t.Run("filter empty map", func(t *testing.T) {
		h := NewHashMap[string, int]()
		filtered := h.Filter(func(k string, v int) bool {
			return v > 0
		})
		if filtered.Size() != 0 {
			t.Error("Filter on empty map should return empty map")
		}
	})

	t.Run("filter all elements", func(t *testing.T) {
		h := NewHashMap[string, int]()
		h.Put("a", 1)
		h.Put("b", 2)
		filtered := h.Filter(func(k string, v int) bool {
			return true
		})
		if filtered.Size() != 2 {
			t.Error("Filter with true predicate should return all elements")
		}
	})

	t.Run("map transformation", func(t *testing.T) {
		h := NewHashMap[string, int]()
		h.Put("x", -1)
		h.Put("y", -2)
		mapped := h.Map(func(k string, v int) int {
			return -v
		})
		want := map[string]int{"x": 1, "y": 2}
		if !reflect.DeepEqual(mapped.ToMap(), want) {
			t.Errorf("Map transformation failed")
		}
	})

	t.Run("forEach order", func(t *testing.T) {
		h := NewHashMap[string, string]()
		h.Put("1", "one")
		h.Put("2", "two")

		var keys []string
		err := h.ForEach(func(k string, v string) error {
			keys = append(keys, k)
			return nil
		})

		if err != nil {
			t.Errorf("ForEach failed: %v", err)
		}

		if len(keys) != 2 {
			t.Error("ForEach didn't iterate all elements")
		}
	})

	t.Run("chained operations", func(t *testing.T) {
		h := NewHashMap[string, int]()
		for i := 0; i < 10; i++ {
			h.Put(fmt.Sprintf("key%d", i), i)
		}

		result := h.Filter(func(k string, v int) bool {
			return v%2 == 0
		}).Map(func(k string, v int) int {
			return v * 2
		})

		if result.Size() != 5 {
			t.Error("Chained operations failed")
		}
	})

	t.Run("concurrent functional operations", func(t *testing.T) {
		h := NewThreadSafeHashMap[string, int]()
		for i := 0; i < 100; i++ {
			h.Put(fmt.Sprintf("key%d", i), i)
		}

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				h.Filter(func(k string, v int) bool {
					return v%2 == 0
				})
				h.Map(func(k string, v int) int {
					return v * 2
				})
			}()
		}
		wg.Wait()
	})
}
