package fastmap_test

import (
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestComboKeyHashMap(t *testing.T) {
	t.Run("basic key combination lookup", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Store value with key combination
		m.Put([]string{"key1", "key2"}, "TEST")

		// Test retrieving with exact keys
		result, exists := m.Get([]string{"key1", "key2"})
		if !exists || result != "TEST" {
			t.Errorf("Expected TEST, got %v", result)
		}

		// Test retrieving with reversed key order
		result, exists = m.Get([]string{"key2", "key1"})
		if !exists || result != "TEST" {
			t.Errorf("Expected TEST, got %v", result)
		}
	})

	t.Run("get by partial keys", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Store multiple combinations
		m.Put([]string{"key1", "key2"}, "TEST1")
		m.Put([]string{"key2", "key3"}, "TEST2")

		// Test getting by single key
		result, exists := m.GetByKeys([]string{"key2"})
		if !exists {
			t.Error("Expected to find a result for key2")
		}
		if result != "TEST1" && result != "TEST2" {
			t.Errorf("Expected either TEST1 or TEST2, got %v", result)
		}

		// Test getting by two keys
		result, exists = m.GetByKeys([]string{"key1", "key2"})
		if !exists || result != "TEST1" {
			t.Errorf("Expected TEST1, got %v", result)
		}
	})

	t.Run("overlapping key combinations", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		m.Put([]string{"key1", "key2"}, "TEST1")
		m.Put([]string{"key1", "key2", "key3"}, "TEST2")

		result, exists := m.GetByKeys([]string{"key1", "key2"})
		if !exists {
			t.Error("Expected to find a result for key1,key2")
		}
		if result != "TEST1" && result != "TEST2" {
			t.Errorf("Expected either TEST1 or TEST2, got %v", result)
		}
	})

	t.Run("empty and non-existent keys", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Test empty keys
		_, exists := m.GetByKeys([]string{})
		if exists {
			t.Error("Expected no result for empty keys")
		}

		// Test non-existent keys
		_, exists = m.GetByKeys([]string{"nonexistent"})
		if exists {
			t.Error("Expected no result for non-existent keys")
		}
	})

	t.Run("partial key matching", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()

		// Store test data
		m.Put([]string{"001", "A", "B", "C"}, "PACKING LIST")
		m.Put([]string{"002", "A", "B", "C"}, "HEALTH CERTIFICATE")

		// Test partial key matching
		result, exists := m.GetByKeys([]string{"001", "A"})
		if !exists {
			t.Error("Expected to find result for partial keys 001,A")
		}
		if result != "PACKING LIST" {
			t.Errorf("Expected PACKING LIST, got %v", result)
		}

		// Test different partial key combination
		result, exists = m.GetByKeys([]string{"002", "B"})
		if !exists {
			t.Error("Expected to find result for partial keys 002,B")
		}
		if result != "HEALTH CERTIFICATE" {
			t.Errorf("Expected HEALTH CERTIFICATE, got %v", result)
		}
	})
}

func TestComboKeyHashMap_Basic(t *testing.T) {
	t.Run("initialization", func(t *testing.T) {
		m := fastmap.NewComboKeyHashMap[string, string]()
		if m == nil {
			t.Error("NewComboKeyHashMap returned nil")
		}
	})
}

func TestComboKeyHashMap_Put(t *testing.T) {
	m := fastmap.NewComboKeyHashMap[string, string]()

	t.Run("put with empty keys", func(t *testing.T) {
		m.Put([]string{}, "TEST")
		_, exists := m.Get([]string{})
		if exists {
			t.Error("Expected no value for empty keys")
		}
	})

	t.Run("put with single key", func(t *testing.T) {
		m.Put([]string{"key1"}, "TEST1")
		val, exists := m.Get([]string{"key1"})
		if !exists || val != "TEST1" {
			t.Errorf("Expected TEST1, got %v", val)
		}
	})

	t.Run("put with multiple keys", func(t *testing.T) {
		m.Put([]string{"key1", "key2", "key3"}, "TEST2")
		val, exists := m.Get([]string{"key1", "key2", "key3"})
		if !exists || val != "TEST2" {
			t.Errorf("Expected TEST2, got %v", val)
		}
	})

	t.Run("put overwrite", func(t *testing.T) {
		m.Put([]string{"key1"}, "TEST1")
		m.Put([]string{"key1"}, "TEST2")
		val, exists := m.Get([]string{"key1"})
		if !exists || val != "TEST2" {
			t.Errorf("Expected TEST2 after overwrite, got %v", val)
		}
	})
}

func TestComboKeyHashMap_Get(t *testing.T) {
	m := fastmap.NewComboKeyHashMap[string, string]()
	m.Put([]string{"key1", "key2"}, "TEST")

	t.Run("get existing combination", func(t *testing.T) {
		val, exists := m.Get([]string{"key1", "key2"})
		if !exists || val != "TEST" {
			t.Errorf("Expected TEST, got %v", val)
		}
	})

	t.Run("get with different order", func(t *testing.T) {
		val, exists := m.Get([]string{"key2", "key1"})
		if !exists || val != "TEST" {
			t.Errorf("Expected TEST with different order, got %v", val)
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		val, exists := m.Get([]string{"nonexistent"})
		if exists {
			t.Errorf("Expected no value for non-existent key, got %v", val)
		}
	})
}

func TestComboKeyHashMap_GetByKeys(t *testing.T) {
	m := fastmap.NewComboKeyHashMap[string, string]()

	t.Run("empty keys", func(t *testing.T) {
		_, exists := m.GetByKeys([]string{})
		if exists {
			t.Error("Expected no result for empty keys")
		}
	})

	t.Run("partial match", func(t *testing.T) {
		m.Put([]string{"001", "A", "B", "C"}, "PACKING LIST")
		val, exists := m.GetByKeys([]string{"001", "A"})
		if !exists || val != "PACKING LIST" {
			t.Errorf("Expected PACKING LIST, got %v", val)
		}
	})

	t.Run("multiple matches", func(t *testing.T) {
		m.Put([]string{"key1", "key2", "key3"}, "TEST1")
		m.Put([]string{"key1", "key2", "key4"}, "TEST2")
		val, exists := m.GetByKeys([]string{"key1", "key2"})
		if !exists {
			t.Error("Expected to find a result for partial match")
		}
		// Should return either TEST1 or TEST2
		if val != "TEST1" && val != "TEST2" {
			t.Errorf("Expected either TEST1 or TEST2, got %v", val)
		}
	})

	t.Run("non-existent keys", func(t *testing.T) {
		val, exists := m.GetByKeys([]string{"nonexistent"})
		if exists {
			t.Errorf("Expected no value for non-existent key, got %v", val)
		}
	})
}

func TestComboKeyHashMap_GetByExactKeys(t *testing.T) {
	m := fastmap.NewComboKeyHashMap[string, string]()
	m.Put([]string{"001", "A", "B", "C"}, "PACKING LIST")

	t.Run("exact match", func(t *testing.T) {
		val, exists := m.GetByExactKeys([]string{"001", "A", "B", "C"})
		if !exists || val != "PACKING LIST" {
			t.Errorf("Expected PACKING LIST for exact match, got %v", val)
		}
	})

	t.Run("partial match should fail", func(t *testing.T) {
		val, exists := m.GetByExactKeys([]string{"001", "A"})
		if exists {
			t.Errorf("Expected no match for partial keys, got %v", val)
		}
	})

	t.Run("different order", func(t *testing.T) {
		val, exists := m.GetByExactKeys([]string{"C", "B", "A", "001"})
		if !exists || val != "PACKING LIST" {
			t.Errorf("Expected PACKING LIST for different order, got %v", val)
		}
	})

	t.Run("empty keys", func(t *testing.T) {
		val, exists := m.GetByExactKeys([]string{})
		if exists {
			t.Errorf("Expected no result for empty keys, got %v", val)
		}
	})

	t.Run("too many keys", func(t *testing.T) {
		val, exists := m.GetByExactKeys([]string{"001", "A", "B", "C", "D"})
		if exists {
			t.Errorf("Expected no result for too many keys, got %v", val)
		}
	})
}

func TestComboKeyHashMap_ComplexTypes(t *testing.T) {
	type ComplexKey struct {
		ID   int
		Name string
	}

	m := fastmap.NewComboKeyHashMap[ComplexKey, string]()
	key1 := ComplexKey{1, "one"}
	key2 := ComplexKey{2, "two"}

	t.Run("complex key types", func(t *testing.T) {
		m.Put([]ComplexKey{key1, key2}, "TEST")
		val, exists := m.Get([]ComplexKey{key1, key2})
		if !exists || val != "TEST" {
			t.Errorf("Expected TEST for complex keys, got %v", val)
		}
	})

	t.Run("complex key partial match", func(t *testing.T) {
		val, exists := m.GetByKeys([]ComplexKey{key1})
		if !exists || val != "TEST" {
			t.Errorf("Expected TEST for partial complex key match, got %v", val)
		}
	})
}

// func TestComboKeyHashMap_EdgeCases(t *testing.T) {
// 	m := fastmap.NewComboKeyHashMap[string, *string]()

// 	t.Run("nil value", func(t *testing.T) {
// 		var nilValue *string
// 		m.Put([]string{"key1"}, nilValue)
// 		val, exists := m.Get([]string{"key1"})
// 		if !exists || val != nil {
// 			t.Error("Failed to handle nil value")
// 		}
// 	})

// 	t.Run("zero value keys", func(t *testing.T) {
// 		m := fastmap.NewComboKeyHashMap[int, string]()
// 		m.Put([]int{0, 1, 2}, "TEST")
// 		val, exists := m.GetByKeys([]int{0, 1})
// 		if !exists || val != "TEST" {
// 			t.Error("Failed to handle zero value keys")
// 		}
// 	})

// 	t.Run("duplicate keys in input", func(t *testing.T) {
// 		m := fastmap.NewComboKeyHashMap[string, string]()
// 		m.Put([]string{"key1", "key1", "key2"}, "TEST")
// 		val, exists := m.Get([]string{"key1", "key2"})
// 		if !exists || val != "TEST" {
// 			t.Error("Failed to handle duplicate keys in input")
// 		}
// 	})
// }
