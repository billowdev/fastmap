package fastmap

// HandleFieldConfigs processes data using field configurations and returns results
// Example:
//
//	configs := map[string]FieldConfig[int]{
//	    "field1": {
//	        Handler: func(data map[string]interface{}) *int {
//	            if val, ok := data["value"].(int); ok {
//	                return &val
//	            }
//	            return nil
//	        },
//	    },
//	}
//	results := safeMap.HandleFieldConfigs(data, configs, "field1")
func (t *ThreadSafeHashMap[K, V]) HandleFieldConfigs(
	data []map[string]interface{},
	configs map[K]FieldConfig[V],
	fieldKey K,
) []V {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.HandleFieldConfigs(data, configs, fieldKey)
}

// ApplyFieldConfig applies a single field configuration to data
// Example:
//
//	config := FieldConfig[int]{
//	    Handler: func(data map[string]interface{}) *int {
//	        if val, ok := data["value"].(int); ok {
//	            return &val
//	        }
//	        return nil
//	    },
//	}
//	success := safeMap.ApplyFieldConfig("field1", config, data)
func (t *ThreadSafeHashMap[K, V]) ApplyFieldConfig(
	key K,
	config FieldConfig[V],
	data map[string]interface{},
) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.data.ApplyFieldConfig(key, config, data)
}

// ProcessFieldConfigs processes data using field configurations with a callback
// Example:
//
//	configs := map[string]FieldConfig[int]{
//	    "field1": {
//	        Handler: func(data map[string]interface{}) *int {
//	            if val, ok := data["value"].(int); ok {
//	                return &val
//	            }
//	            return nil
//	        },
//	    },
//	}
//	safeMap.ProcessFieldConfigs(configs, data, func(key string, value int, index int) {
//	    fmt.Printf("Processed: %s = %d at index %d\n", key, value, index)
//	})
func (t *ThreadSafeHashMap[K, V]) ProcessFieldConfigs(
	configs map[K]FieldConfig[V],
	data []map[string]interface{},
	processor func(key K, value V, index int),
) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.ProcessFieldConfigs(configs, data, processor)
}
