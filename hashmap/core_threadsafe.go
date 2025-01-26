package fastmap

import (
	"sync"
)

// ThreadSafeHashMap provides thread-safe operations for HashMap through mutex synchronization
// Example:
//
//	safeMap := NewThreadSafeHashMap[string, User]()
//	safeMap.Put("user1", User{Name: "John"})
type ThreadSafeHashMap[K comparable, V any] struct {
	mutex sync.RWMutex
	data  *HashMap[K, V]
}

// NewThreadSafeHashMap creates a new thread-safe HashMap
// Example:
//
//	safeMap := NewThreadSafeHashMap[string, User]()
func NewThreadSafeHashMap[K comparable, V any]() *ThreadSafeHashMap[K, V] {
	return &ThreadSafeHashMap[K, V]{
		data: NewHashMap[K, V](),
	}
}

// Put adds or updates a key-value pair in the ThreadSafeHashMap with write lock
// Example:
//
//	safeMap.Put("user123", User{Name: "John"})
func (t *ThreadSafeHashMap[K, V]) Put(key K, value V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.Put(key, value)
}

// Get retrieves a value by key and returns whether it exists with read lock
// Example:
//
//	if user, exists := safeMap.Get("user123"); exists {
//	    fmt.Printf("Found user: %v\n", user)
//	}
func (t *ThreadSafeHashMap[K, V]) Get(key K) (V, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Get(key)
}

// Remove deletes a key-value pair from the ThreadSafeHashMap with write lock
// Example:
//
//	safeMap.Remove("user123")
func (t *ThreadSafeHashMap[K, V]) Remove(key K) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.Remove(key)
}

// Clear removes all elements from the ThreadSafeHashMap with write lock
// Example:
//
//	safeMap.Clear()
//	fmt.Printf("Size after clear: %d\n", safeMap.Size())
func (t *ThreadSafeHashMap[K, V]) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.Clear()
}

// Size returns the number of elements in the ThreadSafeHashMap with read lock
// Example:
//
//	count := safeMap.Size()
//	fmt.Printf("HashMap contains %d elements\n", count)
func (t *ThreadSafeHashMap[K, V]) Size() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Size()
}
