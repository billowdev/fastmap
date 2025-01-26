package fastmap_test

import (
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestHashMap_HandleFieldConfigs(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
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

	if rowIndex != 1 {
		t.Errorf("Expected rowIndex 1, got %d", rowIndex)
	}
}

func TestHashMap_ApplyFieldConfig(t *testing.T) {
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
	result := h.ApplyFieldConfig("field1", config, data)

	if !result {
		t.Error("Expected true, got false")
	}

	if val, _ := h.Get("field1"); val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

func TestHashMap_ProcessFieldConfigs(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	rowIndex := 0
	processedCount := 0

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
		processedCount++
	}

	h.Put("field1", 0)
	h.ProcessFieldConfigs(configs, data, processor)

	if processedCount != 2 {
		t.Errorf("Expected 2 processed items, got %d", processedCount)
	}

	if rowIndex != 1 {
		t.Errorf("Expected rowIndex 1, got %d", rowIndex)
	}
}

func TestHandleFieldConfigs_EmptyData(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			Handler: func(data map[string]interface{}) *int {
				val := 42
				return &val
			},
		},
	}

	h.Put("field1", 0)
	result := h.HandleFieldConfigs([]map[string]interface{}{}, configs, "field1")

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty data, got %d items", len(result))
	}
}

func TestHandleFieldConfigs_NilHandler(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			Handler: nil,
		},
	}

	data := []map[string]interface{}{{"value": 1}}
	h.Put("field1", 0)
	result := h.HandleFieldConfigs(data, configs, "field1")

	if len(result) != 0 {
		t.Errorf("Expected empty result for nil handler, got %d items", len(result))
	}
}

func TestHandleFieldConfigs_MultipleFields(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	index1, index2 := 0, 0
	configs := map[string]fastmap.FieldConfig[int]{
		"field1": {
			RowIndex: &index1,
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value1"].(int); ok {
					return &val
				}
				return nil
			},
		},
		"field2": {
			RowIndex: &index2,
			Handler: func(data map[string]interface{}) *int {
				if val, ok := data["value2"].(int); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"value1": 1, "value2": 10},
		{"value1": 2, "value2": 20},
	}

	h.Put("field1", 0)
	h.Put("field2", 0)

	// Test field1
	result1 := h.HandleFieldConfigs(data, configs, "field1")
	if len(result1) != 2 || result1[0] != 1 || result1[1] != 2 {
		t.Errorf("Field1 results incorrect, got %v", result1)
	}

	// Test field2
	result2 := h.HandleFieldConfigs(data, configs, "field2")
	if len(result2) != 2 || result2[0] != 10 || result2[1] != 20 {
		t.Errorf("Field2 results incorrect, got %v", result2)
	}
}

func TestProcessFieldConfigs_ComplexTypes(t *testing.T) {
	type Complex struct {
		ID    int
		Value string
	}

	h := fastmap.NewHashMap[string, Complex]()
	processedCount := 0

	configs := map[string]fastmap.FieldConfig[Complex]{
		"complex": {
			Handler: func(data map[string]interface{}) *Complex {
				if id, ok := data["id"].(int); ok {
					if val, ok := data["value"].(string); ok {
						return &Complex{ID: id, Value: val}
					}
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"id": 1, "value": "first"},
		{"id": 2, "value": "second"},
	}

	processor := func(key string, value Complex, index int) {
		processedCount++
		if index >= 0 && index < 2 {
			expected := data[index]
			if value.ID != expected["id"].(int) || value.Value != expected["value"].(string) {
				t.Errorf("Processor received incorrect value at index %d", index)
			}
		}
	}

	h.Put("complex", Complex{})
	h.ProcessFieldConfigs(configs, data, processor)

	if processedCount != 2 {
		t.Errorf("Expected 2 processed items, got %d", processedCount)
	}
}

func TestApplyFieldConfig_HandlerReturnsNil(t *testing.T) {
	h := fastmap.NewHashMap[string, int]()
	config := fastmap.FieldConfig[int]{
		Handler: func(data map[string]interface{}) *int {
			return nil
		},
	}

	data := map[string]interface{}{"value": "invalid"}
	h.Put("field1", 42)

	result := h.ApplyFieldConfig("field1", config, data)
	if result {
		t.Error("Expected false when handler returns nil")
	}

	if val, _ := h.Get("field1"); val != 42 {
		t.Error("Value should not change when handler returns nil")
	}
}

func TestHandleFieldConfigs_InvalidData(t *testing.T) {
	h := fastmap.NewHashMap[string, float64]()
	configs := map[string]fastmap.FieldConfig[float64]{
		"field1": {
			Handler: func(data map[string]interface{}) *float64 {
				if val, ok := data["value"].(float64); ok {
					return &val
				}
				return nil
			},
		},
	}

	data := []map[string]interface{}{
		{"value": "invalid"},
		{"value": 1.5},
		{"wrong_key": 2.0},
		{"value": nil},
	}

	h.Put("field1", 0.0)
	result := h.HandleFieldConfigs(data, configs, "field1")

	if len(result) != 1 || result[0] != 1.5 {
		t.Errorf("Expected one valid result (1.5), got %v", result)
	}
}
