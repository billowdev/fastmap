package fastmap_test

import (
	"fmt"
	"testing"

	robinhood "github.com/billowdev/fastmap/robinhood"
)

func BenchmarkRobinHoodPut(b *testing.B) {
	m := robinhood.NewRobinHoodMap[string, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Put(key, i)
	}
}

func BenchmarkRobinHoodGet(b *testing.B) {
	m := robinhood.NewRobinHoodMap[string, int]()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Put(key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		m.Get(key)
	}
}

func BenchmarkRobinHoodRemove(b *testing.B) {
	m := robinhood.NewRobinHoodMap[string, int]()
	keys := make([]string, b.N)

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		keys[i] = key
		m.Put(key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(keys[i])
	}
}

func BenchmarkComparisonWithStandardMap(b *testing.B) {
	b.Run("RobinHood", func(b *testing.B) {
		m := robinhood.NewRobinHoodMap[string, int]()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i)
			m.Put(key, i)
			m.Get(key)
		}
	})

	b.Run("StandardMap", func(b *testing.B) {
		m := make(map[string]int)
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i)
			m[key] = i
			_ = m[key]
		}
	})
}

func BenchmarkHighLoadFactor(b *testing.B) {
	m := robinhood.NewRobinHoodMap[string, int]()
	// Fill to 90% capacity
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("existing%d", i)
		m.Put(key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		m.Put(key, i)
	}
}
