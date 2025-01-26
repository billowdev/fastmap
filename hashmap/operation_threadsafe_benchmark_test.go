package fastmap_test

import (
	"fmt"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkThreadSafeFilter(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Filter(func(k string, v int) bool {
				return v%2 == 0
			})
		}
	})
}

func BenchmarkThreadSafeMap(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Map(func(k string, v int) int {
				return v * 2
			})
		}
	})
}

func BenchmarkThreadSafeForEach(b *testing.B) {
    m := fastmap.NewThreadSafeHashMap[string, int]()
    for i := 0; i < 1000; i++ {
        m.Put(fmt.Sprintf("key%d", i), i)
    }
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            if err := m.ForEach(func(k string, v int) error {
                _ = v
                return nil
            }); err != nil {
                b.Fatalf("ForEach failed: %v", err)
            }
        }
    })
}

func BenchmarkThreadSafeUpdateValue(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	m.Put("testKey", 100)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.UpdateValue("testKey", i)
			i++
		}
	})
}

func BenchmarkThreadSafePutAll(b *testing.B) {
	source := fastmap.NewThreadSafeHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		source.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dest := fastmap.NewThreadSafeHashMap[string, int]()
			dest.PutAll(source)
		}
	})
}
