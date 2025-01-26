package fastmap_test

import (
	"fmt"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkClear(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Clear()
	}
}

func BenchmarkContains(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("testKey", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Contains("testKey")
	}
}

func BenchmarkIsEmpty(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.IsEmpty()
	}
}

func BenchmarkKeys(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Keys()
	}
}

func BenchmarkValues(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Values()
	}
}

func BenchmarkForEach(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.ForEach(func(k string, v int) error {
			_ = v
			return nil
		})
		if err != nil {
			b.Fatalf("ForEach failed: %v", err)
		}
	}
}

func BenchmarkUpdateValue(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	h.Put("testKey", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.UpdateValue("testKey", i)
	}
}

func BenchmarkPutAll(b *testing.B) {
	h1 := fastmap.NewHashMap[string, int]()
	h2 := fastmap.NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h2.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h1.PutAll(h2)
	}
}
