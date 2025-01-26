package fastmap

import (
	"fmt"
	"testing"
)

func BenchmarkHashMapToMap(b *testing.B) {
	h := NewHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		h.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ToMap()
	}
}

func BenchmarkHashMapFromMap(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < 1000; i++ {
		m[fmt.Sprintf("key%d", i)] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FromMap(m)
	}
}
