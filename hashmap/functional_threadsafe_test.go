package fastmap_test

import (
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

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
