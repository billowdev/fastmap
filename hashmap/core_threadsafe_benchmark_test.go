package fastmap_test

import (
	"fmt"
	"sync"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkThreadSafePut(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Put(fmt.Sprintf("key%d", i), i)
			i++
		}
	})
}

func BenchmarkThreadSafeGet(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(fmt.Sprintf("key%d", i%1000))
			i++
		}
	})
}

func BenchmarkThreadSafeRemove(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			m.Put(key, i)
			m.Remove(key)
			i++
		}
	})
}

func BenchmarkThreadSafeClear(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Clear()
	}
}

func BenchmarkThreadSafeConcurrentAccess(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(4)

		// Writer goroutine
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.Put(fmt.Sprintf("key%d", j), j)
			}
		}()

		// Reader goroutine
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.Get(fmt.Sprintf("key%d", j))
			}
		}()

		// Updater goroutine
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.UpdateValue(fmt.Sprintf("key%d", j), j*2)
			}
		}()

		// Remover goroutine
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.Remove(fmt.Sprintf("key%d", j))
			}
		}()

		wg.Wait()
	}
}

func BenchmarkThreadSafeToMap(b *testing.B) {
	m := fastmap.NewThreadSafeHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.ToMap()
		}
	})
}

func BenchmarkThreadSafeFromMap(b *testing.B) {
	standardMap := make(map[string]int)
	for i := 0; i < 1000; i++ {
		standardMap[fmt.Sprintf("key%d", i)] = i
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fastmap.FromThreadSafeMap(standardMap)
		}
	})
}
