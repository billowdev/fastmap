package fastmap_test

import (
	"sync"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestThreadSafeHandleFieldConfigs(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	rowIndex := 0
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			RowIndex: &rowIndex,
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"value": 1},
		{"value": 2},
	}

	h.Put("field1", 0)
	result := h.HandleFieldConfigs(data, configs, "field1")

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}

	if result[0] != 1 || result[1] != 2 {
		t.Errorf("Expected values [1,2], got %v", result)
	}
}

func TestThreadSafeApplyFieldConfig(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
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
	result := h.ApplyFieldConfig("field1", config, data)

	if !result {
		t.Error("Expected true, got false")
	}

	if val, _ := h.Get("field1"); val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

func TestThreadSafeProcessFieldConfigs(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	rowIndex := 0
	processedCount := 0
	var mu sync.Mutex

	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			RowIndex: &rowIndex,
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"value": 1},
		{"value": 2},
	}

	processor := func(key string, value int, index int) {
		mu.Lock()
		processedCount++
		mu.Unlock()
	}

	h.Put("field1", 0)
	h.ProcessFieldConfigs(configs, data, processor)

	if processedCount != 2 {
		t.Errorf("Expected 2 processed items, got %d", processedCount)
	}
}

func TestProcessFieldConfigs_ConcurrentAccess(t *testing.T) {
	h := fastmap.NewThreadSafeHashMap[string, int]()
	processedCount := 0
	var mu sync.Mutex

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

	data := make([]map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]interface{}{"value": i}
	}

	processor := func(key string, value int, index int) {
		mu.Lock()
		processedCount++
		mu.Unlock()
	}

	h.Put("field1", 0)
	h.ProcessFieldConfigs(configs, data, processor)

	if processedCount != 100 {
		t.Errorf("Expected 100 processed items, got %d", processedCount)
	}
}
