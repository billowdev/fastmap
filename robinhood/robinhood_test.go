package fastmap_test

import (
	"fmt"
	"testing"

	robinhood "github.com/billowdev/fastmap/robinhood"
)

func TestNewRobinHoodMap(t *testing.T) {
	m := robinhood.NewRobinHoodMap[string, int]()
	if m == nil {
		t.Error("NewRobinHoodMap returned nil")
	}
	if m.Size() != 0 {
		t.Errorf("New map should have size 0, got %d", m.Size())
	}
}

func TestPutAndGet(t *testing.T) {
	m := robinhood.NewRobinHoodMap[string, int]()

	// Test single put and get
	m.Put("key1", 100)
	value, exists := m.Get("key1")
	if !exists {
		t.Error("Expected to find key1")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %d", value)
	}

	// Test overwriting existing key
	m.Put("key1", 200)
	value, exists = m.Get("key1")
	if !exists {
		t.Error("Expected to find key1 after overwrite")
	}
	if value != 200 {
		t.Errorf("Expected value 200 after overwrite, got %d", value)
	}
}

func TestRemove(t *testing.T) {
	m := robinhood.NewRobinHoodMap[string, int]()

	// Put multiple items
	m.Put("key1", 100)
	m.Put("key2", 200)
	m.Put("key3", 300)

	// Remove existing key
	removed := m.Remove("key2")
	if !removed {
		t.Error("Failed to remove existing key")
	}

	// Verify removal
	_, exists := m.Get("key2")
	if exists {
		t.Error("Removed key should not exist")
	}

	// Verify other keys remain
	value, exists := m.Get("key1")
	if !exists || value != 100 {
		t.Error("Other keys should remain after removal")
	}

	// Remove non-existing key
	removed = m.Remove("key4")
	if removed {
		t.Error("Should not remove non-existing key")
	}
}

func TestSize(t *testing.T) {
	m := robinhood.NewRobinHoodMap[string, int]()

	// Initial size
	if m.Size() != 0 {
		t.Errorf("Initial size should be 0, got %d", m.Size())
	}

	// Add items
	m.Put("key1", 100)
	if m.Size() != 1 {
		t.Errorf("Size should be 1, got %d", m.Size())
	}

	m.Put("key2", 200)
	if m.Size() != 2 {
		t.Errorf("Size should be 2, got %d", m.Size())
	}

	// Remove item
	m.Remove("key1")
	if m.Size() != 1 {
		t.Errorf("Size should be 1 after removal, got %d", m.Size())
	}
}

func TestClear(t *testing.T) {
	m := robinhood.NewRobinHoodMap[string, int]()

	// Add multiple items
	m.Put("key1", 100)
	m.Put("key2", 200)
	m.Put("key3", 300)

	// Clear map
	m.Clear()

	// Verify cleared state
	if m.Size() != 0 {
		t.Errorf("Size should be 0 after clear, got %d", m.Size())
	}

	// Verify no items can be retrieved
	_, exists := m.Get("key1")
	if exists {
		t.Error("No items should exist after clear")
	}
}

func TestResizing(t *testing.T) {
	m := robinhood.NewRobinHoodMap[string, int]()

	// Add many items to trigger resizing
	for i := 0; i < 100; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	// Verify all items are still retrievable
	for i := 0; i < 100; i++ {
		value, exists := m.Get(fmt.Sprintf("key%d", i))
		if !exists {
			t.Errorf("Failed to retrieve key%d after resizing", i)
		}
		if value != i {
			t.Errorf("Incorrect value for key%d after resizing", i)
		}
	}

	// Verify size
	if m.Size() != 100 {
		t.Errorf("Size should be 100, got %d", m.Size())
	}
}

func TestComplexKeyTypes(t *testing.T) {
	// Test with struct keys
	type ComplexKey struct {
		ID   int
		Name string
	}

	m := robinhood.NewRobinHoodMap[ComplexKey, string]()

	key1 := ComplexKey{ID: 1, Name: "Test"}
	m.Put(key1, "value1")

	value, exists := m.Get(key1)
	if !exists {
		t.Error("Failed to retrieve complex key")
	}
	if value != "value1" {
		t.Error("Incorrect value for complex key")
	}
}

func BenchmarkPut(b *testing.B) {
	m := robinhood.NewRobinHoodMap[string, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}
}

func BenchmarkGet(b *testing.B) {
	m := robinhood.NewRobinHoodMap[string, int]()

	// Populate map
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(fmt.Sprintf("key%d", i%1000))
	}
}
