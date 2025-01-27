package fastmap_test

import (
	"fmt"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkHashMapFilterLarge(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 10000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Filter(func(k string, v int) bool {
			return v%2 == 0 && v%3 == 0
		})
	}
}

func BenchmarkHashMapMapLarge(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 10000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Map(func(k string, v int) int {
			return v*2 + v%3
		})
	}
}
func BenchmarkHashMapForEachLarge(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 10000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		if err := h.ForEach(func(k string, v int) error {
			sum += v % 10
			return nil
		}); err != nil {
			b.Fatalf("ForEach failed: %v", err)
		}
	}
}

func BenchmarkHashMapChainedOperations(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := h.Filter(func(k string, v int) bool {
			return v%2 == 0
		})
		filtered.Map(func(k string, v int) int {
			return v * 2
		})
	}
}
