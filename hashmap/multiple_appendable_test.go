package fastmap_test

import (
	"reflect"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

// func TestAppendableMultiKeyHashMap_Basic(t *testing.T) {
// 	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

// 	// Test appending to new key
// 	m.AppendValues("key1", 1, 2, 3)
// 	if values, exists := m.GetSlice("key1"); !exists || !reflect.DeepEqual(values, []int{1, 2, 3}) {
// 		t.Error("Failed to append values to new key")
// 	}

// 	// Test appending to existing key
// 	m.AppendValues("key1", 4, 5)
// 	if values, exists := m.GetSlice("key1"); !exists || !reflect.DeepEqual(values, []int{1, 2, 3, 4, 5}) {
// 		t.Error("Failed to append values to existing key")
// 	}
// }

func TestAppendableMultiKeyHashMap_Aliases(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Add values with aliases
	m.AppendValuesWithKeys([]string{"main", "alias1", "alias2"}, 1, 2)

	// Test accessing through different keys
	tests := []string{"main", "alias1", "alias2"}
	expected := []int{1, 2}

	for _, key := range tests {
		if values, exists := m.GetSlice(key); !exists || !reflect.DeepEqual(values, expected) {
			t.Errorf("Failed to get correct values through key %s", key)
		}
	}

	// Append through alias
	m.AppendValues("alias1", 3, 4)
	expected = []int{1, 2, 3, 4}

	// Verify update propagated to all keys
	for _, key := range tests {
		if values, exists := m.GetSlice(key); !exists || !reflect.DeepEqual(values, expected) {
			t.Errorf("Failed to propagate update through key %s", key)
		}
	}
}

func TestAppendableMultiKeyHashMap_UpdateSlice(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Setup initial data
	m.AppendValuesWithKeys([]string{"main", "alias1"}, 1, 2)

	// Update through alias
	newValues := []int{5, 6, 7}
	if !m.UpdateSlice("alias1", newValues) {
		t.Error("UpdateSlice failed")
	}

	// Verify update through all keys
	tests := []string{"main", "alias1"}
	for _, key := range tests {
		if values, exists := m.GetSlice(key); !exists || !reflect.DeepEqual(values, newValues) {
			t.Errorf("Failed to update values through key %s", key)
		}
	}
}

func TestAppendableMultiKeyHashMap_EdgeCases(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Test empty keys
	if m.AppendValuesWithKeys([]string{}, 1, 2) {
		t.Error("AppendValuesWithKeys should fail with empty keys")
	}

	// Test non-existent key
	if m.AppendValues("nonexistent", 1, 2) {
		t.Error("AppendValues should fail with non-existent key")
	}

	// Test empty values
	m.AppendValuesWithKeys([]string{"key1"})
	if values, exists := m.GetSlice("key1"); !exists || len(values) != 0 {
		t.Error("Failed to handle empty values")
	}
}
