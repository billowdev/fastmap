// conversion_test.go
package fastmap

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestToMap(t *testing.T) {
	h := NewHashMap[string, int]()
	h.Put("one", 1)
	h.Put("two", 2)

	want := map[string]int{"one": 1, "two": 2}
	got := h.ToMap()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToMap() = %v, want %v", got, want)
	}
}

func TestFromMap(t *testing.T) {
	original := map[string]int{"one": 1, "two": 2}
	h := FromMap(original)

	if h.Size() != len(original) {
		t.Errorf("FromMap() size = %v, want %v", h.Size(), len(original))
	}

	for k, v := range original {
		if got, exists := h.Get(k); !exists || got != v {
			t.Errorf("FromMap() value for key %v = %v, want %v", k, got, v)
		}
	}
}
func TestConversionEdgeCases(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		h := FromMap(map[string]int{})
		if h.Size() != 0 {
			t.Error("FromMap empty map failed")
		}
		if len(h.ToMap()) != 0 {
			t.Error("ToMap empty map failed")
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[string]int
		h := FromMap(m)
		if h.Size() != 0 {
			t.Error("FromMap nil map failed")
		}
	})

	t.Run("large map", func(t *testing.T) {
		m := make(map[string]int)
		for i := 0; i < 10000; i++ {
			m[fmt.Sprintf("key%d", i)] = i
		}
		h := FromMap(m)
		if !reflect.DeepEqual(h.ToMap(), m) {
			t.Error("Large map conversion failed")
		}
	})

	t.Run("concurrent access", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				m := map[string]int{fmt.Sprintf("key%d", i): i}
				temp := FromMap(m)
				temp.ToMap()
			}(i)
		}
		wg.Wait()
	})

	t.Run("non-string keys", func(t *testing.T) {
		m := map[int]string{1: "one", 2: "two"}
		h := FromMap(m)
		if !reflect.DeepEqual(h.ToMap(), m) {
			t.Error("Integer key conversion failed")
		}
	})

	t.Run("pointer values", func(t *testing.T) {
		type Data struct {
			value int // Changed field name to be used
		}
		m := map[string]*Data{"a": {1}, "b": {2}}
		h := FromMap(m)
		convertedMap := h.ToMap()

		// Compare values explicitly
		for k, v := range m {
			if got, exists := convertedMap[k]; !exists || got.value != v.value {
				t.Errorf("For key %s, got %v, want %v", k, got, v)
			}
		}
	})
}
