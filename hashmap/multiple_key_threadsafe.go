package fastmap

import "sync"

// ThreadSafeMultiKeyHashMap provides thread-safe operations for MultiKeyHashMap through mutex synchronization
type ThreadSafeMultiKeyHashMap[K comparable, V any] struct {
	mutex sync.RWMutex
	data  *MultiKeyHashMap[K, V]
}

// NewThreadSafeMultiKeyHashMap creates a new thread-safe MultiKeyHashMap
// Example:
//
//	safeMap := NewThreadSafeMultiKeyHashMap[string, User]()
func NewThreadSafeMultiKeyHashMap[K comparable, V any]() *ThreadSafeMultiKeyHashMap[K, V] {
	return &ThreadSafeMultiKeyHashMap[K, V]{
		data: NewMultiKeyHashMap[K, V](),
	}
}

// Put adds a value with multiple keys with write lock
// Example:
//
//	safeMap.Put([]string{"main", "alias1", "alias2"}, user)
func (t *ThreadSafeMultiKeyHashMap[K, V]) Put(keys []K, value V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.Put(keys, value)
}

// Get retrieves a value by key with read lock
// Example:
//
//	if user, exists := safeMap.Get("alias1"); exists {
//	    fmt.Printf("Found user: %v\n", user)
//	}
func (t *ThreadSafeMultiKeyHashMap[K, V]) Get(key K) (V, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Get(key)
}

// GetPrimaryKey returns the primary key for any given key (alias or primary) with read lock
// Example:
//
//	if primaryKey, exists := safeMap.GetPrimaryKey("alias1"); exists {
//	    fmt.Printf("Primary key: %v\n", primaryKey)
//	}
func (t *ThreadSafeMultiKeyHashMap[K, V]) GetPrimaryKey(key K) (K, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.GetPrimaryKey(key)
}

// GetAllKeys returns all keys (primary and aliases) associated with a given key with read lock
// Example:
//
//	keys := safeMap.GetAllKeys("alias1")
//	fmt.Printf("All associated keys: %v\n", keys)
func (t *ThreadSafeMultiKeyHashMap[K, V]) GetAllKeys(key K) []K {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.GetAllKeys(key)
}

// AddAlias adds a new alias to an existing key with write lock
// Example:
//
//	if safeMap.AddAlias("main", "newAlias") {
//	    fmt.Println("Alias added successfully")
//	}
func (t *ThreadSafeMultiKeyHashMap[K, V]) AddAlias(existingKey K, newAlias K) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.data.AddAlias(existingKey, newAlias)
}

// Size returns the number of unique primary keys with read lock
// Example:
//
//	count := safeMap.Size()
//	fmt.Printf("Number of unique entries: %d\n", count)
func (t *ThreadSafeMultiKeyHashMap[K, V]) Size() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Size()
}

// Clear removes all entries from the map with write lock
// Example:
//
//	safeMap.Clear()
//	fmt.Printf("Size after clear: %d\n", safeMap.Size())
func (t *ThreadSafeMultiKeyHashMap[K, V]) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.Clear()
}

// Remove removes a key and its associated aliases with write lock
// Example:
//
//	safeMap.Remove("main")  // Removes main key and its aliases
//	safeMap.Remove("alias") // Removes only the alias
func (t *ThreadSafeMultiKeyHashMap[K, V]) Remove(key K) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.Remove(key)
}

// RemoveWithCascade removes a key and all connected keys with write lock
// Example:
//
//	safeMap.RemoveWithCascade("main") // Removes main key and all connected keys
func (t *ThreadSafeMultiKeyHashMap[K, V]) RemoveWithCascade(key K) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.RemoveWithCascade(key)
}
