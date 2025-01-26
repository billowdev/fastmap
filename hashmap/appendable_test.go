package fastmap_test

import (
	"reflect"
	"testing"

	fastmap "github.com/billowdev/fastmap/hashmap"
)

func TestAppendableHashMap_NewAppendableHashMap(t *testing.T) {
	// Test initialization
	t.Run("initialization", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		if hashMap == nil {
			t.Error("NewAppendableHashMap returned nil")
		}
		if hashMap.HashMap == nil {
			t.Error("Underlying HashMap is nil")
		}
	})

	// Test initial state
	t.Run("initial state", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		if val, exists := hashMap.Get("nonexistent"); exists || val != nil {
			t.Error("New map should be empty")
		}
	})
}

func TestAppendableHashMap_AppendValues(t *testing.T) {
	// Test appending to non-existent key
	t.Run("append to new key", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		hashMap.AppendValues("key1", 1, 2, 3)

		values, exists := hashMap.Get("key1")
		if !exists {
			t.Error("Key should exist after AppendValues")
		}

		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(values, expected) {
			t.Errorf("Expected %v, got %v", expected, values)
		}
	})

	// Test appending to existing key
	t.Run("append to existing key", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		hashMap.AppendValues("key1", 1, 2)
		hashMap.AppendValues("key1", 3, 4)

		values, exists := hashMap.Get("key1")
		if !exists {
			t.Error("Key should exist after AppendValues")
		}

		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(values, expected) {
			t.Errorf("Expected %v, got %v", expected, values)
		}
	})

	// Test appending empty slice
	t.Run("append empty slice", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		hashMap.AppendValues("key1")

		values, exists := hashMap.Get("key1")
		if !exists {
			t.Error("Key should exist after AppendValues")
		}

		if len(values) != 0 {
			t.Errorf("Expected empty slice, got %v", values)
		}
	})

	// Test appending to multiple keys
	t.Run("append to multiple keys", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		hashMap.AppendValues("key1", 1, 2)
		hashMap.AppendValues("key2", 3, 4)

		values1, exists1 := hashMap.Get("key1")
		values2, exists2 := hashMap.Get("key2")

		if !exists1 || !exists2 {
			t.Error("Both keys should exist")
		}

		expected1 := []int{1, 2}
		expected2 := []int{3, 4}

		if !reflect.DeepEqual(values1, expected1) {
			t.Errorf("Key1: Expected %v, got %v", expected1, values1)
		}
		if !reflect.DeepEqual(values2, expected2) {
			t.Errorf("Key2: Expected %v, got %v", expected2, values2)
		}
	})

	// Test with custom struct type
	t.Run("append custom struct type", func(t *testing.T) {
		type CustomType struct {
			ID   int
			Name string
		}

		hashMap := fastmap.NewAppendableHashMap[string, CustomType]()
		items := []CustomType{
			{ID: 1, Name: "First"},
			{ID: 2, Name: "Second"},
		}

		hashMap.AppendValues("key1", items...)
		hashMap.AppendValues("key1", CustomType{ID: 3, Name: "Third"})

		values, exists := hashMap.Get("key1")
		if !exists {
			t.Error("Key should exist")
		}

		expected := []CustomType{
			{ID: 1, Name: "First"},
			{ID: 2, Name: "Second"},
			{ID: 3, Name: "Third"},
		}

		if !reflect.DeepEqual(values, expected) {
			t.Errorf("Expected %v, got %v", expected, values)
		}
	})

	// Test capacity handling
	t.Run("capacity handling", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()

		// Append large number of values
		largeSlice := make([]int, 1000)
		for i := range largeSlice {
			largeSlice[i] = i
		}

		hashMap.AppendValues("key1", largeSlice...)
		values, exists := hashMap.Get("key1")

		if !exists {
			t.Error("Key should exist")
		}
		if len(values) != len(largeSlice) {
			t.Errorf("Expected length %d, got %d", len(largeSlice), len(values))
		}
	})

	// Test nil handling
	t.Run("nil value handling", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, *int]()
		var nilValue *int
		value := 42
		ptrValue := &value

		hashMap.AppendValues("key1", nilValue, ptrValue, nilValue)

		values, exists := hashMap.Get("key1")
		if !exists {
			t.Error("Key should exist")
		}

		if len(values) != 3 {
			t.Errorf("Expected length 3, got %d", len(values))
		}

		if values[0] != nil {
			t.Error("First value should be nil")
		}
		if *values[1] != 42 {
			t.Errorf("Second value should be 42, got %v", *values[1])
		}
		if values[2] != nil {
			t.Error("Third value should be nil")
		}
	})

	// Test zero value handling
	t.Run("zero value handling", func(t *testing.T) {
		hashMap := fastmap.NewAppendableHashMap[string, int]()
		hashMap.AppendValues("key1", 0, 1, 0)

		values, exists := hashMap.Get("key1")
		if !exists {
			t.Error("Key should exist")
		}

		expected := []int{0, 1, 0}
		if !reflect.DeepEqual(values, expected) {
			t.Errorf("Expected %v, got %v", expected, values)
		}
	})
}
