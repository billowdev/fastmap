package fastmap_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestMultiKeyHashMap_EmptyKeys(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Test Put with empty keys slice
	m.Put([]string{}, 100)
	if size := m.Size(); size != 0 {
		t.Errorf("Size after putting empty keys = %v, want 0", size)
	}

	// Test Put with nil keys slice
	var nilKeys []string
	m.Put(nilKeys, 100)
	if size := m.Size(); size != 0 {
		t.Errorf("Size after putting nil keys = %v, want 0", size)
	}
}

func TestMultiKeyHashMap_ZeroValues(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[int, string]()

	// Test with zero value keys
	m.Put([]int{0, 1, 2}, "test")
	if val, exists := m.Get(0); !exists || val != "test" {
		t.Error("Failed to handle zero value key")
	}

	// Test with zero value as alias
	m.AddAlias(0, 3)
	if val, exists := m.Get(3); !exists || val != "test" {
		t.Error("Failed to handle alias to zero value key")
	}
}

func TestMultiKeyHashMap_DuplicateKeys(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Test Put with duplicate keys
	m.Put([]string{"key1", "key1", "alias1"}, 100)
	if size := m.Size(); size != 1 {
		t.Errorf("Size with duplicate keys = %v, want 1", size)
	}

	// Test AddAlias with existing alias
	m.AddAlias("key1", "alias1")
	allKeys := m.GetAllKeys("key1")
	if len(allKeys) != 2 { // should be [key1, alias1]
		t.Errorf("GetAllKeys with duplicate alias = %v, want 2 keys", len(allKeys))
	}
}

func TestMultiKeyHashMap_CircularAliases(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Setup initial key-value
	m.Put([]string{"key1", "alias1"}, 100)

	// Try to create circular reference
	m.AddAlias("alias1", "key1")

	// Verify primary key remains unchanged
	if primaryKey, exists := m.GetPrimaryKey("alias1"); exists && primaryKey != "key1" {
		t.Error("Primary key changed unexpectedly")
	}
}

func TestMultiKeyHashMap_CascadingRemoval(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Create a chain of aliases
	m.Put([]string{"key1", "alias1"}, 100)
	m.AddAlias("key1", "alias2")
	m.AddAlias("key1", "alias3") // Note: aliases can only be added to primary keys

	// Remove an alias and verify only that alias is removed
	m.Remove("alias2")

	// alias2 should be removed
	if _, exists := m.Get("alias2"); exists {
		t.Error("alias2 should have been removed")
	}

	// Other keys should still exist
	expectedExisting := []string{"key1", "alias1", "alias3"}
	for _, key := range expectedExisting {
		if _, exists := m.Get(key); !exists {
			t.Errorf("Key %v should still exist", key)
		}
	}
}

func TestMultiKeyHashMap_OverwriteValue(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Initial value
	m.Put([]string{"key1", "alias1"}, 100)

	// Overwrite with new value using primary key
	m.Put([]string{"key1", "alias2"}, 200)

	// Check all aliases have updated value
	tests := []struct {
		key  string
		want int
	}{
		{"key1", 200},
		{"alias1", 100}, // Should be removed as it's not in new aliases
		{"alias2", 200},
	}

	for _, tt := range tests {
		if val, exists := m.Get(tt.key); exists && val != tt.want {
			t.Errorf("After overwrite, Get(%v) = %v, want %v", tt.key, val, tt.want)
		}
	}
}

func TestMultiKeyHashMap_ConcurrentKeys(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Add value with multiple keys
	m.Put([]string{"key1", "alias1"}, 100)

	// Try to use the same key in another entry
	m.Put([]string{"alias1", "newAlias"}, 200)

	// Verify original keys maintained their associations
	if val, exists := m.Get("key1"); !exists || val != 100 {
		t.Error("Original key-value pair was modified")
	}
}

func TestMultiKeyHashMap_NilValue(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, *int]()
	var nilValue *int

	// Test storing nil value
	m.Put([]string{"key1", "alias1"}, nilValue)
	if val, exists := m.Get("key1"); !exists || val != nil {
		t.Error("Failed to store and retrieve nil value")
	}

	// Test overwriting nil value
	value := 42
	m.Put([]string{"key1", "alias1"}, &value)
	if val, exists := m.Get("alias1"); !exists || *val != 42 {
		t.Error("Failed to overwrite nil value")
	}
}

func TestMultiKeyHashMap_ComplexKeys(t *testing.T) {
	type ComplexKey struct {
		ID   int
		Name string
	}

	m := fastmap.NewMultiKeyHashMap[ComplexKey, string]()
	key1 := ComplexKey{1, "one"}
	key2 := ComplexKey{2, "two"}

	// Test with struct keys
	m.Put([]ComplexKey{key1, key2}, "value")
	if val, exists := m.Get(key2); !exists || val != "value" {
		t.Error("Failed to handle complex struct keys")
	}
}

func TestMultiKeyHashMap_GetAllKeysEdgeCases(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Test with non-existent key
	if keys := m.GetAllKeys("nonexistent"); keys != nil {
		t.Error("GetAllKeys for non-existent key should return nil")
	}

	// Test with single key (no aliases)
	m.Put([]string{"key1"}, 100)
	keys := m.GetAllKeys("key1")
	if len(keys) != 1 || keys[0] != "key1" {
		t.Error("GetAllKeys for single key without aliases failed")
	}

	// Test after removing some aliases
	m.Put([]string{"key2", "alias1", "alias2"}, 200)
	m.Remove("alias1")
	keys = m.GetAllKeys("key2")
	if len(keys) != 2 { // should contain key2 and alias2
		t.Errorf("GetAllKeys after partial alias removal = %v keys, want 2", len(keys))
	}
}

func TestMultiKeyHashMap_SizeEdgeCases(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Test size after adding and removing
	m.Put([]string{"key1", "alias1"}, 100)
	m.Remove("alias1")
	if size := m.Size(); size != 1 {
		t.Errorf("Size after removing alias = %v, want 1", size)
	}

	// Test size after clearing and adding new value
	m.Clear()
	m.Put([]string{"key2"}, 200)
	if size := m.Size(); size != 1 {
		t.Errorf("Size after clear and new addition = %v, want 1", size)
	}

	// Test size with multiple values sharing aliases
	m.Put([]string{"key3", "shared"}, 300)
	m.Put([]string{"key4", "shared"}, 400)
	if size := m.Size(); size != 3 {
		t.Errorf("Size with shared aliases = %v, want 3", size)
	}
}

func TestMultiKeyHashMap_ValueOverwriting(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, float64]()

	// Test special float values
	tests := []struct {
		name   string
		keys   []string
		value  float64
		check  string
		expect float64
	}{
		{"infinity", []string{"inf", "inf_alias"}, math.Inf(1), "inf_alias", math.Inf(1)},
		{"negative infinity", []string{"neg_inf", "neg_inf_alias"}, math.Inf(-1), "neg_inf", math.Inf(-1)},
		{"nan", []string{"nan", "nan_alias"}, math.NaN(), "nan", math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.Put(tt.keys, tt.value)
			if val, exists := m.Get(tt.check); !exists {
				t.Errorf("%s: value not found", tt.name)
			} else if tt.name == "nan" {
				if !math.IsNaN(val) {
					t.Errorf("%s: expected NaN, got %v", tt.name, val)
				}
			} else if val != tt.expect {
				t.Errorf("%s: got %v, want %v", tt.name, val, tt.expect)
			}
		})
	}
}

func TestMultiKeyHashMap_AliasChaining(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Create a chain of aliases
	m.Put([]string{"primary", "alias1"}, 100)
	m.AddAlias("primary", "alias2")
	m.AddAlias("alias2", "alias3") // Should associate with primary
	m.AddAlias("alias3", "alias4") // Should associate with primary

	// Test primary key resolution
	testCases := []struct {
		key           string
		expectedPrim  string
		shouldExist   bool
		expectedValue int
	}{
		{"primary", "primary", true, 100},
		{"alias1", "primary", true, 100},
		{"alias2", "primary", true, 100},
		{"alias3", "primary", true, 100},
		{"alias4", "primary", true, 100},
		{"nonexistent", "", false, 0},
	}

	for _, tc := range testCases {
		prim, exists := m.GetPrimaryKey(tc.key)
		if exists != tc.shouldExist {
			t.Errorf("GetPrimaryKey(%s): got exists=%v, want %v", tc.key, exists, tc.shouldExist)
		}
		if exists && prim != tc.expectedPrim {
			t.Errorf("GetPrimaryKey(%s): got %s, want %s", tc.key, prim, tc.expectedPrim)
		}
		if val, exists := m.Get(tc.key); exists != tc.shouldExist {
			t.Errorf("Get(%s): got exists=%v, want %v", tc.key, exists, tc.shouldExist)
		} else if exists && val != tc.expectedValue {
			t.Errorf("Get(%s): got %d, want %d", tc.key, val, tc.expectedValue)
		}
	}
}

// func TestMultiKeyHashMap_CascadingRemovalComplex(t *testing.T) {
// 	m := fastmap.NewMultiKeyHashMap[string, string]()

// 	// Create complex alias relationships
// 	m.Put([]string{"root", "alias1", "alias2"}, "value1")
// 	m.Put([]string{"branch1", "b1_alias1", "b1_alias2"}, "value2")
// 	m.AddAlias("root", "shared_alias")
// 	m.AddAlias("branch1", "shared_alias") // This should not work as shared_alias is already used

// 	// Test RemoveWithCascade
// 	m.RemoveWithCascade("alias1")

// 	// Check removals
// 	removedKeys := []string{"root", "alias1", "alias2", "shared_alias"}
// 	for _, key := range removedKeys {
// 		if _, exists := m.Get(key); exists {
// 			t.Errorf("Key %s should have been removed", key)
// 		}
// 	}

// 	// Check preserved keys
// 	preservedKeys := []string{"branch1", "b1_alias1", "b1_alias2"}
// 	for _, key := range preservedKeys {
// 		if _, exists := m.Get(key); !exists {
// 			t.Errorf("Key %s should have been preserved", key)
// 		}
// 	}
// }

func TestMultiKeyHashMap_ConcurrentModification(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[int, string]()

	// Setup initial state
	initialKeys := []int{1, 2, 3}
	m.Put(initialKeys, "initial")

	// Modify while iterating through aliases
	allKeys := m.GetAllKeys(1)
	for _, key := range allKeys {
		m.AddAlias(key, key+10)
		m.Remove(key + 5) // Remove non-existent keys
	}

	// Verify state
	if size := m.Size(); size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// Check all valid keys still work
	for _, key := range initialKeys {
		if val, exists := m.Get(key); !exists || val != "initial" {
			t.Errorf("Key %d should still exist with value 'initial'", key)
		}
	}
}

func TestMultiKeyHashMap_EmptyStringKeys(t *testing.T) {
	m := fastmap.NewMultiKeyHashMap[string, int]()

	// Test with empty string keys
	emptyKeys := []string{"", " ", "  "}
	m.Put(emptyKeys, 100)

	// Verify empty string handling
	testCases := []struct {
		key     string
		exists  bool
		wantVal int
	}{
		{"", true, 100},
		{" ", true, 100},
		{"  ", true, 100},
		{"   ", false, 0}, // Not in original keys
	}

	for _, tc := range testCases {
		val, exists := m.Get(tc.key)
		if exists != tc.exists {
			t.Errorf("Get(%q): got exists=%v, want %v", tc.key, exists, tc.exists)
		}
		if exists && val != tc.wantVal {
			t.Errorf("Get(%q): got %d, want %d", tc.key, val, tc.wantVal)
		}
	}

	// Test alias operations with empty strings
	m.AddAlias("", "new_empty")
	if val, exists := m.Get("new_empty"); !exists || val != 100 {
		t.Error("Failed to add alias to empty string key")
	}
}

func TestMultiKeyHashMap_KeyTypeEdgeCases(t *testing.T) {
	type CustomKey struct {
		ID        int
		Reference *string
	}

	m := fastmap.NewMultiKeyHashMap[CustomKey, int]()

	// Create keys with nil and non-nil references
	str1 := "ref1"
	str2 := "ref2"
	keys := []CustomKey{
		{ID: 1, Reference: nil},
		{ID: 2, Reference: &str1},
		{ID: 3, Reference: &str2},
	}

	// Test Put and Get with complex keys
	m.Put(keys, 100)

	// Test retrieval with recreated keys
	testCases := []struct {
		name    string
		key     CustomKey
		exists  bool
		wantVal int
	}{
		{"nil reference", CustomKey{ID: 1, Reference: nil}, true, 100},
		{"with reference", CustomKey{ID: 2, Reference: &str1}, true, 100},
		{"different reference same content", CustomKey{ID: 2, Reference: &str2}, false, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, exists := m.Get(tc.key)
			if exists != tc.exists {
				t.Errorf("Get(%v): got exists=%v, want %v", tc.key, exists, tc.exists)
			}
			if exists && val != tc.wantVal {
				t.Errorf("Get(%v): got %d, want %d", tc.key, val, tc.wantVal)
			}
		})
	}
}

func TestAppendableMultiKeyHashMap_GetValuesByKeys(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Setup test data
	m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2, 3)
	m.AppendValuesWithKeys([]string{"key2", "alias2"}, 4, 5, 6)

	// Test GetValuesByKeys with existing keys
	keys := []string{"key1", "alias1", "key2"}
	result := m.GetValuesByKeys(keys)

	// Verify results
	expected := map[string][]int{
		"key1":   {1, 2, 3},
		"alias1": {1, 2, 3},
		"key2":   {4, 5, 6},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetValuesByKeys = %v, want %v", result, expected)
	}

	// Test with non-existent keys
	nonExistentKeys := []string{"key1", "nonexistent"}
	result = m.GetValuesByKeys(nonExistentKeys)
	expected = map[string][]int{
		"key1": {1, 2, 3},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetValuesByKeys with non-existent key = %v, want %v", result, expected)
	}
}

func TestAppendableMultiKeyHashMap_GetValuesByKeysWithPrimary(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Setup test data
	m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2, 3)
	m.AppendValuesWithKeys([]string{"key2", "alias2"}, 4, 5, 6)

	// Test GetValuesByKeysWithPrimary using mix of primary keys and aliases
	keys := []string{"alias1", "alias2", "key2"}
	result := m.GetValuesByKeysWithPrimary(keys)

	// Verify results - all values should be mapped to primary keys
	expected := map[string][]int{
		"key1": {1, 2, 3},
		"key2": {4, 5, 6},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetValuesByKeysWithPrimary = %v, want %v", result, expected)
	}

	// Test with non-existent keys
	nonExistentKeys := []string{"key1", "nonexistent"}
	result = m.GetValuesByKeysWithPrimary(nonExistentKeys)
	expected = map[string][]int{
		"key1": {1, 2, 3},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetValuesByKeysWithPrimary with non-existent key = %v, want %v", result, expected)
	}
}

func TestAppendableMultiKeyHashMap_GetValuesByKeys_EmptyInput(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Test with empty keys slice
	result := m.GetValuesByKeys([]string{})
	if len(result) != 0 {
		t.Errorf("GetValuesByKeys with empty input should return empty map, got %v", result)
	}

	// Test with nil keys slice
	var nilKeys []string
	result = m.GetValuesByKeys(nilKeys)
	if len(result) != 0 {
		t.Errorf("GetValuesByKeys with nil input should return empty map, got %v", result)
	}
}

func TestAppendableMultiKeyHashMap_GetValuesByKeysWithPrimary_EmptyInput(t *testing.T) {
	m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

	// Test with empty keys slice
	result := m.GetValuesByKeysWithPrimary([]string{})
	if len(result) != 0 {
		t.Errorf("GetValuesByKeysWithPrimary with empty input should return empty map, got %v", result)
	}

	// Test with nil keys slice
	var nilKeys []string
	result = m.GetValuesByKeysWithPrimary(nilKeys)
	if len(result) != 0 {
		t.Errorf("GetValuesByKeysWithPrimary with nil input should return empty map, got %v", result)
	}
}
func TestAppendableMultiKeyHashMap_GetValuesByKeys_EdgeCases(t *testing.T) {
	t.Run("duplicate keys in input", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()
		m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2, 3)

		// Test with duplicate keys in input
		result := m.GetValuesByKeys([]string{"key1", "key1", "alias1", "alias1"})
		expected := map[string][]int{
			"key1":   {1, 2, 3},
			"alias1": {1, 2, 3},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("GetValuesByKeys with duplicate keys = %v, want %v", result, expected)
		}
	})

	t.Run("cyclic references", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

		// Create a situation with potential cyclic references
		m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2)
		m.AppendValuesWithKeys([]string{"key2", "alias1"}, 3, 4) // alias1 now points to key2

		result := m.GetValuesByKeys([]string{"key1", "alias1", "key2"})
		if len(result) != 3 {
			t.Errorf("Expected 3 entries in result, got %d", len(result))
		}
	})

	t.Run("zero value keys", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[int, string]()
		m.AppendValuesWithKeys([]int{0, 1}, "zero", "one")

		result := m.GetValuesByKeys([]int{0, 1, 2})
		expected := map[int][]string{
			0: {"zero", "one"},
			1: {"zero", "one"},
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("GetValuesByKeys with zero value key = %v, want %v", result, expected)
		}
	})

	t.Run("nil values in slice", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, *int]()
		val1 := 1
		var nilVal *int
		m.AppendValuesWithKeys([]string{"key1"}, &val1, nilVal, nil)

		result := m.GetValuesByKeys([]string{"key1"})
		if len(result["key1"]) != 3 {
			t.Errorf("Expected slice with 3 elements (including nils), got %d", len(result["key1"]))
		}
	})

	t.Run("large number of aliases", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()
		aliases := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			aliases[i] = fmt.Sprintf("alias%d", i)
		}
		aliases = append([]string{"primaryKey"}, aliases...)

		m.AppendValuesWithKeys(aliases, 1, 2, 3)

		// Try to get values using various combinations of aliases
		result := m.GetValuesByKeys([]string{
			"primaryKey",
			"alias0",
			"alias499",
			"alias999",
		})

		for _, values := range result {
			if !reflect.DeepEqual(values, []int{1, 2, 3}) {
				t.Errorf("Incorrect values for large number of aliases: %v", values)
			}
		}
	})

	t.Run("overwritten aliases", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

		// First assignment
		m.AppendValuesWithKeys([]string{"key1", "shared"}, 1, 2)

		// Second assignment with shared alias
		m.AppendValuesWithKeys([]string{"key2", "shared"}, 3, 4)

		result := m.GetValuesByKeys([]string{"key1", "key2", "shared"})
		// Verify that "shared" points to the latest assignment
		if !reflect.DeepEqual(result["shared"], []int{3, 4}) {
			t.Errorf("Overwritten alias should point to latest values, got %v", result["shared"])
		}
	})

	t.Run("mixed key types", func(t *testing.T) {
		type CustomKey struct {
			id   int
			name string
		}

		m := fastmap.NewAppendableMultiKeyHashMap[CustomKey, int]()
		key1 := CustomKey{1, "one"}
		alias1 := CustomKey{1, "alias"}

		m.AppendValuesWithKeys([]CustomKey{key1, alias1}, 1, 2)

		result := m.GetValuesByKeys([]CustomKey{key1, alias1})
		if len(result) != 2 {
			t.Errorf("Expected 2 entries for complex key type, got %d", len(result))
		}
	})

	t.Run("concurrent key modification", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()
		m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2)

		// Get values while modifying
		go func() {
			m.AppendValuesWithKeys([]string{"key1", "alias2"}, 3, 4)
		}()

		result := m.GetValuesByKeys([]string{"key1", "alias1"})
		// Verify we got valid data regardless of concurrent modification
		if len(result) == 0 {
			t.Error("Got empty result during concurrent modification")
		}
	})

	t.Run("empty slices", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()
		m.AppendValuesWithKeys([]string{"key1", "alias1"}, []int{}...)

		result := m.GetValuesByKeys([]string{"key1", "alias1"})
		if len(result) != 2 {
			t.Error("Should return entries even for empty slices")
		}
		for _, values := range result {
			if len(values) != 0 {
				t.Error("Expected empty slice values")
			}
		}
	})

	t.Run("removed primary keys", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()
		m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2)

		// Remove primary key but try to access via alias
		m.Remove("key1")

		result := m.GetValuesByKeysWithPrimary([]string{"alias1"})
		if len(result) != 0 {
			t.Error("Should not be able to access values through alias after primary key removal")
		}
	})
}

func TestAppendableMultiKeyHashMap_GetByExactKeys(t *testing.T) {
	t.Run("exact match found", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

		// Set up data with multiple keys pointing to same values
		m.AppendValuesWithKeys([]string{"key1", "alias1", "alias2"}, 1, 2, 3)

		// Test getting by subset of keys
		values, exists := m.GetByExactKeys([]string{"key1", "alias1"})
		if !exists {
			t.Error("Expected to find exact match")
		}
		if !reflect.DeepEqual(values, []int{1, 2, 3}) {
			t.Errorf("Got values %v, want [1,2,3]", values)
		}
	})

	t.Run("no match when keys point to different values", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

		// Set up data with keys pointing to different values
		m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2)
		m.AppendValuesWithKeys([]string{"key2", "alias2"}, 3, 4)

		// Test getting by keys that point to different values
		_, exists := m.GetByExactKeys([]string{"key1", "alias2"})
		if exists {
			t.Error("Should not find match for keys pointing to different values")
		}
	})

	t.Run("empty keys", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

		_, exists := m.GetByExactKeys([]string{})
		if exists {
			t.Error("Should return false for empty keys")
		}
	})

	t.Run("non-existent keys", func(t *testing.T) {
		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

		_, exists := m.GetByExactKeys([]string{"nonexistent1", "nonexistent2"})
		if exists {
			t.Error("Should return false for non-existent keys")
		}
	})
}

// func TestAppendableMultiKeyHashMap_GetByExactKeysAll(t *testing.T) {
// 	t.Run("multiple exact matches", func(t *testing.T) {
// 		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

// 		// Set up multiple groups of related keys
// 		m.AppendValuesWithKeys([]string{"key1", "alias1", "alias2"}, 1, 2)
// 		m.AppendValuesWithKeys([]string{"key2", "alias3", "alias4"}, 3, 4)

// 		keyGroups := [][]string{
// 			{"key1", "alias1"},
// 			{"key2", "alias3"},
// 		}

// 		result := m.GetByExactKeysAll(keyGroups)

// 		expected := map[string][]int{
// 			"key1": {1, 2},
// 			"key2": {3, 4},
// 		}

// 		if !reflect.DeepEqual(result, expected) {
// 			t.Errorf("Got result %v, want %v", result, expected)
// 		}
// 	})

// 	t.Run("some matches, some not", func(t *testing.T) {
// 		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

// 		m.AppendValuesWithKeys([]string{"key1", "alias1"}, 1, 2)

// 		keyGroups := [][]string{
// 			{"key1", "alias1"},         // should match
// 			{"nonexistent1", "alias2"}, // should not match
// 		}

// 		result := m.GetByExactKeysAll(keyGroups)

// 		expected := map[string][]int{
// 			"key1": {1, 2},
// 		}

// 		if !reflect.DeepEqual(result, expected) {
// 			t.Errorf("Got result %v, want %v", result, expected)
// 		}
// 	})

// 	t.Run("shared aliases", func(t *testing.T) {
// 		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

// 		// Set up data with a shared alias
// 		m.AppendValuesWithKeys([]string{"key1", "shared"}, 1, 2)
// 		m.AppendValuesWithKeys([]string{"key2", "shared"}, 3, 4)

// 		keyGroups := [][]string{
// 			{"key1", "shared"},
// 			{"key2", "shared"},
// 		}

// 		result := m.GetByExactKeysAll(keyGroups)

// 		// Only the latest assignment should match
// 		expected := map[string][]int{
// 			"key2": {3, 4},
// 		}

// 		if !reflect.DeepEqual(result, expected) {
// 			t.Errorf("Got result %v, want %v", result, expected)
// 		}
// 	})

// 	t.Run("empty key groups", func(t *testing.T) {
// 		m := fastmap.NewAppendableMultiKeyHashMap[string, int]()

// 		result := m.GetByExactKeysAll([][]string{})
// 		if len(result) != 0 {
// 			t.Error("Should return empty map for empty key groups")
// 		}
// 	})
// }
