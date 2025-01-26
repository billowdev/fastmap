package fastmap_test

import (
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkHashMapPut(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < b.N; i++ {
		h.Put("key", i)
	}
}

func BenchmarkHashMapGet(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("key", 123)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Get("key")
	}
}

func BenchmarkHashMapRemove(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < b.N; i++ {
		h.Put("key", i)
		h.Remove("key")
	}
}
