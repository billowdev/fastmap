package fastmap_test

import (
	"fmt"
	"testing"

	"github.com/billowdev/fastmap"
)

func BenchmarkHashMapFilter(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Filter(func(k string, v int) bool {
			return v%2 == 0
		})
	}
}

func BenchmarkHashMapMap(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Map(func(k string, v int) int {
			return v * 2
		})
	}
}

func BenchmarkHashMapForEach(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		h.ForEach(func(k string, v int) {
			sum += v
		})
	}
}
