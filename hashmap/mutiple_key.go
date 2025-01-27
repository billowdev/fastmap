package fastmap

// MultiKeyHashMap extends HashMap to support multiple keys (aliases) for accessing values.
// It maintains a primary key and optional aliases for each value.
type MultiKeyHashMap[K comparable, V any] struct {
	data    map[K]V
	aliases map[K][]K // Maps primary keys to their aliases
}

// NewMultiKeyHashMap creates a new MultiKeyHashMap instance
// Example:
//
//	map := NewMultiKeyHashMap[string, User]()
func NewMultiKeyHashMap[K comparable, V any]() *MultiKeyHashMap[K, V] {
	return &MultiKeyHashMap[K, V]{
		data:    make(map[K]V),
		aliases: make(map[K][]K),
	}
}

// Put adds or updates a value with multiple keys where first key is primary
// Example:
//
//	map.Put([]string{"main", "alias1", "alias2"}, value)
func (m *MultiKeyHashMap[K, V]) Put(keys []K, value V) {
	if len(keys) == 0 {
		return
	}

	// First, if the primary key exists, remove all its old aliases
	if oldAliases, exists := m.aliases[keys[0]]; exists {
		for _, alias := range oldAliases {
			delete(m.data, alias)
		}
	}

	primaryKey := keys[0]
	m.data[primaryKey] = value

	// Create a set of unique aliases
	uniqueAliases := make(map[K]bool)
	for _, alias := range keys[1:] {
		if alias != primaryKey { // Skip if alias is same as primary key
			uniqueAliases[alias] = true
		}
	}

	// Convert unique aliases to slice
	aliases := make([]K, 0, len(uniqueAliases))
	for alias := range uniqueAliases {
		aliases = append(aliases, alias)
		m.data[alias] = value
	}

	if len(aliases) > 0 {
		m.aliases[primaryKey] = aliases
	}
}

// Get retrieves a value using any associated key (primary or alias)
// Example:
//
//	if value, exists := map.Get("alias1"); exists {
//	    fmt.Printf("Found: %v\n", value)
//	}
func (m *MultiKeyHashMap[K, V]) Get(key K) (V, bool) {
	value, exists := m.data[key]
	return value, exists
}

// GetPrimaryKey returns the primary key for any given key (alias or primary)
// Example:
//
//	if primaryKey, exists := map.GetPrimaryKey("alias1"); exists {
//	    fmt.Printf("Primary: %v\n", primaryKey)
//	}
func (m *MultiKeyHashMap[K, V]) GetPrimaryKey(key K) (K, bool) {
	if _, exists := m.data[key]; !exists {
		var zero K
		return zero, false
	}

	// First check if it's a primary key
	if _, isPrimary := m.aliases[key]; isPrimary {
		return key, true
	}

	// Check if it's an alias
	for primaryKey, aliases := range m.aliases {
		for _, alias := range aliases {
			if alias == key {
				return primaryKey, true
			}
		}
	}

	// If the key exists but isn't in aliases map as either primary or alias,
	// then it must be a standalone primary key
	return key, true
}

// GetAllKeys returns all keys (primary and aliases) associated with a given key
// Example:
//
//	keys := map.GetAllKeys("alias1")
//	fmt.Printf("All keys: %v\n", keys)
func (m *MultiKeyHashMap[K, V]) GetAllKeys(key K) []K {
	// First check if the key exists
	if _, exists := m.data[key]; !exists {
		return nil
	}

	// Get the primary key
	primaryKey, exists := m.GetPrimaryKey(key)
	if !exists {
		return nil
	}

	// Start with the primary key
	result := make([]K, 0)
	result = append(result, primaryKey)

	// Add all current aliases
	if aliases, exists := m.aliases[primaryKey]; exists {
		result = append(result, aliases...)
	}

	return result
}

// AddAlias adds a new alias to an existing key
// Example:
//
//	success := map.AddAlias("main", "newAlias")
func (m *MultiKeyHashMap[K, V]) AddAlias(existingKey K, newAlias K) bool {
	primaryKey, exists := m.GetPrimaryKey(existingKey)
	if !exists {
		return false
	}

	// Check if the new alias is the same as primary key
	if newAlias == primaryKey {
		return false
	}

	// Add the value mapping for the new alias
	if value, exists := m.data[primaryKey]; exists {
		m.data[newAlias] = value
		// Check if the alias already exists
		for _, alias := range m.aliases[primaryKey] {
			if alias == newAlias {
				return true
			}
		}
		m.aliases[primaryKey] = append(m.aliases[primaryKey], newAlias)
		return true
	}

	return false
}

// Size returns the number of unique primary keys
// Example:
//
//	count := map.Size()
//	fmt.Printf("Unique entries: %d\n", count)
func (m *MultiKeyHashMap[K, V]) Size() int {
	// Count primary keys only
	primaryKeys := make(map[K]bool)
	for key := range m.data {
		primaryKey, exists := m.GetPrimaryKey(key)
		if exists {
			primaryKeys[primaryKey] = true
		}
	}
	return len(primaryKeys)
}

// Clear removes all entries from the map
// Example:
//
//	map.Clear()
func (m *MultiKeyHashMap[K, V]) Clear() {
	m.data = make(map[K]V)
	m.aliases = make(map[K][]K)
}

// Remove removes a key and potentially its connected keys
// If removing primary key, all aliases are removed
// If removing alias, only that alias is removed
// Example:
//
//	map.Remove("main")  // Removes main key and aliases
//	map.Remove("alias") // Removes only the alias
func (m *MultiKeyHashMap[K, V]) Remove(key K) {
	// First check if this key exists
	if _, exists := m.data[key]; !exists {
		return
	}

	// Get the primary key for this key
	primaryKey, exists := m.GetPrimaryKey(key)
	if !exists {
		return
	}

	if key == primaryKey {
		// If removing primary key, remove all aliases
		if aliases, exists := m.aliases[primaryKey]; exists {
			for _, alias := range aliases {
				delete(m.data, alias)
			}
			delete(m.aliases, primaryKey)
		}
		delete(m.data, primaryKey)
	} else {
		// If removing an alias, just remove it and update the aliases list
		delete(m.data, key)
		if aliases, exists := m.aliases[primaryKey]; exists {
			newAliases := make([]K, 0, len(aliases))
			for _, alias := range aliases {
				if alias != key {
					newAliases = append(newAliases, alias)
				}
			}

			// Update or remove aliases list based on remaining aliases
			if len(newAliases) > 0 {
				m.aliases[primaryKey] = newAliases
			} else {
				delete(m.aliases, primaryKey)
			}
		}
	}
}

// RemoveWithCascade removes a key and all connected keys through alias relationships
// Example:
//
//	map.RemoveWithCascade("main") // Removes main key and all connected keys
func (m *MultiKeyHashMap[K, V]) RemoveWithCascade(key K) {
	// First check if this key exists
	if _, exists := m.data[key]; !exists {
		return
	}

	// Get all connected keys (includes the key itself and all related keys)
	allKeys := m.getAllConnectedKeys(key)

	// Remove all connected keys
	for _, k := range allKeys {
		delete(m.data, k)
		delete(m.aliases, k)
	}
}

// getAllConnectedKeys returns all keys connected through alias relationships
func (m *MultiKeyHashMap[K, V]) getAllConnectedKeys(key K) []K {
	visited := make(map[K]bool)
	result := make([]K, 0)
	m.traverseConnectedKeys(key, visited, &result)
	return result
}

// traverseConnectedKeys performs DFS to find all connected keys
func (m *MultiKeyHashMap[K, V]) traverseConnectedKeys(key K, visited map[K]bool, result *[]K) {
	if visited[key] {
		return
	}

	visited[key] = true
	*result = append(*result, key)

	// Get primary key for current key
	primaryKey, exists := m.GetPrimaryKey(key)
	if !exists {
		return
	}

	// Add primary key and check its aliases
	if !visited[primaryKey] {
		m.traverseConnectedKeys(primaryKey, visited, result)
	}

	// Check aliases of primary key
	if aliases, exists := m.aliases[primaryKey]; exists {
		for _, alias := range aliases {
			if !visited[alias] {
				m.traverseConnectedKeys(alias, visited, result)
			}
		}
	}
}
