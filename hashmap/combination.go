package fastmap

import (
	"fmt"
)

// ComboKeyHashMap maps combinations of keys to a single value
type ComboKeyHashMap[K comparable, V any] struct {
	data      map[string]V        // stores values
	keyGroups map[K][]KeyGroup[K] // maps individual keys to their key groups
}

// NewComboKeyHashMap creates a new ComboKeyHashMap
func NewComboKeyHashMap[K comparable, V any]() *ComboKeyHashMap[K, V] {
	return &ComboKeyHashMap[K, V]{
		data:      make(map[string]V),
		keyGroups: make(map[K][]KeyGroup[K]),
	}
}

// generateKeyString creates a unique string for a key group
func (m *ComboKeyHashMap[K, V]) generateKeyString(keys []K) string {
	sorted := make([]K, len(keys))
	copy(sorted, keys)
	sortSlice(sorted)
	result := ""
	for i, key := range sorted {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%v", key)
	}
	return result
}

// Put associates a value with a combination of keys
func (m *ComboKeyHashMap[K, V]) Put(keys []K, value V) {
	if len(keys) == 0 {
		return
	}

	keyString := m.generateKeyString(keys)
	m.data[keyString] = value

	// Store key group reference for each key
	keyGroup := make(KeyGroup[K], len(keys))
	copy(keyGroup, keys)

	for _, key := range keys {
		m.keyGroups[key] = append(m.keyGroups[key], keyGroup)
	}
}

// Get retrieves a value by its exact key combination
func (m *ComboKeyHashMap[K, V]) Get(keys []K) (V, bool) {
	keyString := m.generateKeyString(keys)
	val, exists := m.data[keyString]
	return val, exists
}

// GetByKeys retrieves all values where all provided keys are part of the key group
func (m *ComboKeyHashMap[K, V]) GetByKeys(keys []K) []V {
	if len(keys) == 0 {
		return []V{}
	}

	// Track unique values using a map to avoid duplicates
	valueMap := make(map[string]V)

	// Check each key's groups
	for _, key := range keys {
		groups := m.keyGroups[key]
		for _, group := range groups {
			// Check if all provided keys are in this group
			allKeysFound := true
			for _, searchKey := range keys {
				if !contains(group, searchKey) {
					allKeysFound = false
					break
				}
			}

			if allKeysFound {
				keyString := m.generateKeyString(group)
				if val, exists := m.data[keyString]; exists {
					valueMap[keyString] = val
				}
			}
		}
	}

	// Convert map to slice
	result := make([]V, 0, len(valueMap))
	for _, v := range valueMap {
		result = append(result, v)
	}

	return result
}

// Helper functions

func contains[K comparable](slice []K, item K) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func sortSlice[K comparable](slice []K) {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if fmt.Sprintf("%v", slice[j]) > fmt.Sprintf("%v", slice[j+1]) {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
}

type KeyGroup[K comparable] []K
