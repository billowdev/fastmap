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

// GetValuesByKeys retrieves values associated with an array of keys
// If a key doesn't exist, it will be skipped in the result
func (m *AppendableMultiKeyHashMap[K, V]) GetValuesByKeys(keys []K) map[K][]V {
	result := make(map[K][]V)

	for _, key := range keys {
		if values, exists := m.GetSlice(key); exists {
			result[key] = values
		}
	}

	return result
}

// GetValuesByKeysWithPrimary retrieves values associated with an array of keys
// The returned map uses primary keys instead of the provided aliases
func (m *AppendableMultiKeyHashMap[K, V]) GetValuesByKeysWithPrimary(keys []K) map[K][]V {
	result := make(map[K][]V)

	for _, key := range keys {
		if primaryKey, exists := m.GetPrimaryKey(key); exists {
			if values, valExists := m.GetSlice(primaryKey); valExists {
				result[primaryKey] = values
			}
		}
	}

	return result
}

// GetByExactKeys retrieves values only if all provided keys are associated with the same value set
// Returns the values and true if an exact match is found, empty slice and false otherwise
func (m *AppendableMultiKeyHashMap[K, V]) GetByExactKeys(keys []K) ([]V, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	// First find the latest primary key among all input keys
	var latestPrimaryKey K
	var foundPrimary bool

	for _, key := range keys {
		if currentPrimary, exists := m.GetPrimaryKey(key); exists {
			if !foundPrimary {
				latestPrimaryKey = currentPrimary
				foundPrimary = true
			} else if currentPrimary != latestPrimaryKey {
				// If we find a different primary key that has one of our keys as an alias,
				// it means this key was reassigned more recently
				if _, exists := m.GetSlice(currentPrimary); exists {
					latestPrimaryKey = currentPrimary
				}
			}
		}
	}

	if !foundPrimary {
		return nil, false
	}

	// Now verify all keys point to this latest primary key
	for _, key := range keys {
		if currentPrimary, exists := m.GetPrimaryKey(key); !exists || currentPrimary != latestPrimaryKey {
			return nil, false
		}
	}

	// Return the values associated with the latest primary key
	values, exists := m.GetSlice(latestPrimaryKey)
	return values, exists
}

// GetByExactKeysAll returns a map of all entries where the provided keys match exactly
// Each entry in the result map contains the primary key and its values
func (m *AppendableMultiKeyHashMap[K, V]) GetByExactKeysAll(keyGroups [][]K) map[K][]V {
	result := make(map[K][]V)

	for _, keys := range keyGroups {
		if values, exists := m.GetByExactKeys(keys); exists {
			// Use the primary key for storing the result
			if primaryKey, exists := m.GetPrimaryKey(keys[0]); exists {
				result[primaryKey] = values
			}
		}
	}

	return result
}
