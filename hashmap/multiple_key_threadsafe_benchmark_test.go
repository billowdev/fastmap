package fastmap_test

import (
	"fmt"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkThreadSafeMultiKeyMap_Put(b *testing.B) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			keys := []string{fmt.Sprintf("primary%d", i), fmt.Sprintf("alias%d", i)}
			m.Put(keys, i)
			i++
		}
	})
}

func BenchmarkThreadSafeMultiKeyMap_Get(b *testing.B) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		keys := []string{fmt.Sprintf("primary%d", i), fmt.Sprintf("alias%d", i)}
		m.Put(keys, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(fmt.Sprintf("primary%d", i%1000))
			m.Get(fmt.Sprintf("alias%d", i%1000))
			i++
		}
	})
}

func BenchmarkThreadSafeMultiKeyMap_AddAlias(b *testing.B) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		keys := []string{fmt.Sprintf("primary%d", i)}
		m.Put(keys, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.AddAlias(fmt.Sprintf("primary%d", i%1000), fmt.Sprintf("newalias%d", i))
			i++
		}
	})
}

func BenchmarkThreadSafeMultiKeyMap_GetAllKeys(b *testing.B) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	for i := 0; i < 1000; i++ {
		keys := []string{
			fmt.Sprintf("primary%d", i),
			fmt.Sprintf("alias1_%d", i),
			fmt.Sprintf("alias2_%d", i),
		}
		m.Put(keys, i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.GetAllKeys(fmt.Sprintf("primary%d", i%1000))
			i++
		}
	})
}

func BenchmarkThreadSafeMultiKeyMap_RemoveWithCascade(b *testing.B) {
	m := fastmap.NewThreadSafeMultiKeyHashMap[string, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			keys := []string{
				fmt.Sprintf("primary%d", i),
				fmt.Sprintf("alias1_%d", i),
				fmt.Sprintf("alias2_%d", i),
			}
			m.Put(keys, i)
			m.RemoveWithCascade(fmt.Sprintf("primary%d", i))
			i++
		}
	})
}