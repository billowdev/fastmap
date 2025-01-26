package fastmap

type FieldConfig[T any] struct {
	RowIndex *int
	Handler  func(data map[string]interface{}) *T
}

func (h *HashMap[K, V]) HandleFieldConfigs(
	data []map[string]interface{},
	configs map[K]FieldConfig[V],
	fieldKey K,
) []V {
	result := make([]V, 0, len(data)*len(configs))

	for index, temp := range data {
		if _, exists := h.Get(fieldKey); exists {
			if handler := configs[fieldKey].Handler; handler != nil {
				if configs[fieldKey].RowIndex != nil {
					*configs[fieldKey].RowIndex = index
				}
				if value := handler(temp); value != nil {
					result = append(result, *value)
				}
			}
		}
	}
	return result
}

func (h *HashMap[K, V]) ApplyFieldConfig(
	key K,
	config FieldConfig[V],
	data map[string]interface{},
) bool {
	if value := config.Handler(data); value != nil {
		return h.UpdateValue(key, *value)
	}
	return false
}

func (h *HashMap[K, V]) ProcessFieldConfigs(
	configs map[K]FieldConfig[V],
	data []map[string]interface{},
	processor func(key K, value V, index int),
) {
	for index, temp := range data {
		for key, config := range configs {
			if h.Contains(key) {
				if config.RowIndex != nil {
					*config.RowIndex = index
				}
				if value := config.Handler(temp); value != nil {
					if h.UpdateValue(key, *value) {
						if processor != nil {
							if val, exists := h.Get(key); exists {
								processor(key, val, index)
							}
						}
					}
				}
			}
		}
	}
}
