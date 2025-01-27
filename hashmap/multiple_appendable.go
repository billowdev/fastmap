package fastmap

// AppendableMultiKeyHashMap extends MultiKeyHashMap to support slice operations
type AppendableMultiKeyHashMap[K comparable, V any] struct {
	*MultiKeyHashMap[K, []V]
}

// NewAppendableMultiKeyHashMap creates a new AppendableMultiKeyHashMap instance
func NewAppendableMultiKeyHashMap[K comparable, V any]() *AppendableMultiKeyHashMap[K, V] {
	return &AppendableMultiKeyHashMap[K, V]{
		MultiKeyHashMap: NewMultiKeyHashMap[K, []V](),
	}
}

// AppendValues appends values to the slice associated with the primary key or any of its aliases
func (m *AppendableMultiKeyHashMap[K, V]) AppendValues(key K, values ...V) bool {
	// Get primary key if it exists
	primaryKey, exists := m.GetPrimaryKey(key)
	if !exists {
		return false
	}

	// Get existing values or create new slice
	if currentValues, exists := m.Get(primaryKey); exists {
		newValues := append(currentValues, values...)
		m.Put(m.GetAllKeys(primaryKey), newValues)
	} else {
		m.Put([]K{primaryKey}, values)
	}

	return true
}

// AppendValuesWithKeys appends values and associates them with multiple keys
func (m *AppendableMultiKeyHashMap[K, V]) AppendValuesWithKeys(keys []K, values ...V) bool {
	if len(keys) == 0 {
		return false
	}

	// Get existing values or create new slice
	if currentValues, exists := m.Get(keys[0]); exists {
		newValues := append(currentValues, values...)
		m.Put(keys, newValues)
	} else {
		m.Put(keys, values)
	}

	return true
}

// GetSlice returns the slice associated with any key (primary or alias)
func (m *AppendableMultiKeyHashMap[K, V]) GetSlice(key K) ([]V, bool) {
	return m.Get(key)
}

// UpdateSlice updates the entire slice for a given key and its aliases
func (m *AppendableMultiKeyHashMap[K, V]) UpdateSlice(key K, newValues []V) bool {
	if primaryKey, exists := m.GetPrimaryKey(key); exists {
		m.Put(m.GetAllKeys(primaryKey), newValues)
		return true
	}
	return false
}
