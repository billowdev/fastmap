package fastmap_test

import (
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func BenchmarkHandleFieldConfigs(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{"value": i}
	}

	h.Put("field1", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.HandleFieldConfigs(data, configs, "field1")
	}
}

func BenchmarkApplyFieldConfig(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	config := fastmap.FieldConfig[int]{
		Handler: func(data map[string]interface{}) *int {
			if val, ok := data["value"].(int); ok {
				return &val
			}
			return nil
		},
	}

	data := map[string]interface{}{"value": 42}
	h.Put("field1", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ApplyFieldConfig("field1", config, data)
	}
}

func BenchmarkProcessFieldConfigs(b *testing.B) {
	h := fastmap.NewHashMap[string, int]()
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := make([]map[string]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{"value": i}
	}

	h.Put("field1", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.ProcessFieldConfigs(configs, data, func(k string, v int, idx int) {})
	}
}
