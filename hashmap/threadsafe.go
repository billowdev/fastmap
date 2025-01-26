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

// Contains checks if a key exists in the ThreadSafeHashMap with read lock
// Example:
//
//	if safeMap.Contains("user123") {
//	    fmt.Println("User exists")
//	}
func (t *ThreadSafeHashMap[K, V]) Contains(key K) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Contains(key)
}

// IsEmpty returns true if the ThreadSafeHashMap has no elements with read lock
// Example:
//
//	if safeMap.IsEmpty() {
//	    fmt.Println("HashMap is empty")
//	}
func (t *ThreadSafeHashMap[K, V]) IsEmpty() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.IsEmpty()
}

// Keys returns a slice of all keys in the ThreadSafeHashMap with read lock
// Example:
//
//	keys := safeMap.Keys()
//	for _, key := range keys {
//	    fmt.Printf("Key: %v\n", key)
//	}
func (t *ThreadSafeHashMap[K, V]) Keys() []K {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Keys()
}

// Values returns a slice of all values in the ThreadSafeHashMap with read lock
// Example:
//
//	values := safeMap.Values()
//	for _, value := range values {
//	    fmt.Printf("Value: %v\n", value)
//	}
func (t *ThreadSafeHashMap[K, V]) Values() []V {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.Values()
}

// ForEach executes a callback function for each key-value pair with read lock
// Example:
//
//	safeMap.ForEach(func(key string, value User) {
//	    fmt.Printf("User %s: %v\n", key, value)
//	})
func (t *ThreadSafeHashMap[K, V]) ForEach(callback ValueCallback[K, V]) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	t.data.ForEach(callback)
}

// Filter returns a new ThreadSafeHashMap containing only the elements that satisfy the predicate with read lock
// Example:
//
//	activeUsers := safeMap.Filter(func(key string, user User) bool {
//	    return user.Active
//	})
func (t *ThreadSafeHashMap[K, V]) Filter(predicate func(K, V) bool) *ThreadSafeHashMap[K, V] {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	result := NewThreadSafeHashMap[K, V]()
	result.data = t.data.Filter(predicate)
	return result
}

// Map transforms values using the provided function and returns a new ThreadSafeHashMap with read lock
// Example:
//
//	upperNames := safeMap.Map(func(key string, user User) User {
//	    user.Name = strings.ToUpper(user.Name)
//	    return user
//	})
func (t *ThreadSafeHashMap[K, V]) Map(transform func(K, V) V) *ThreadSafeHashMap[K, V] {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	result := NewThreadSafeHashMap[K, V]()
	result.data = t.data.Map(transform)
	return result
}

// UpdateValue updates an existing value by key with write lock, returns false if key doesn't exist
// Example:
//
//	if safeMap.UpdateValue("user123", updatedUser) {
//	    fmt.Println("User updated successfully")
//	}
func (t *ThreadSafeHashMap[K, V]) UpdateValue(key K, newValue V) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.data.UpdateValue(key, newValue)
}

// PutAll adds all key-value pairs from another ThreadSafeHashMap with write lock
// Example:
//
//	otherMap := NewThreadSafeHashMap[string, User]()
//	otherMap.Put("user456", newUser)
//	safeMap.PutAll(otherMap)
func (t *ThreadSafeHashMap[K, V]) PutAll(other *ThreadSafeHashMap[K, V]) {
	other.mutex.RLock()
	t.mutex.Lock()
	defer other.mutex.RUnlock()
	defer t.mutex.Unlock()
	t.data.PutAll(other.data)
}

// ToMap returns the underlying map with read lock
// Example:
//
//	standardMap := safeMap.ToMap()
//	for k, v := range standardMap {
//	    fmt.Printf("%v: %v\n", k, v)
//	}
func (t *ThreadSafeHashMap[K, V]) ToMap() map[K]V {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.ToMap()
}

// FromThreadSafeMap creates a new ThreadSafeHashMap from a regular map
// Example:
//
//	regularMap := map[string]int{"one": 1, "two": 2}
//	safeMap := FromThreadSafeMap(regularMap)
func FromThreadSafeMap[K comparable, V any](m map[K]V) *ThreadSafeHashMap[K, V] {
	result := NewThreadSafeHashMap[K, V]()
	result.data = FromMap(m)
	return result
}

// HandleFieldConfigs processes data using field configurations and returns results
// Example:
//
//	configs := map[string]FieldConfig[int]{
//	    "field1": {
//	        Handler: func(data map[string]interface{}) *int {
//	            if val, ok := data["value"].(int); ok {
//	                return &val
//	            }
//	            return nil
//	        },
//	    },
//	}
//	results := safeMap.HandleFieldConfigs(data, configs, "field1")
func (t *ThreadSafeHashMap[K, V]) HandleFieldConfigs(
	data []map[string]interface{},
	configs map[K]FieldConfig[V],
	fieldKey K,
) []V {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.HandleFieldConfigs(data, configs, fieldKey)
}

// ApplyFieldConfig applies a single field configuration to data
// Example:
//
//	config := FieldConfig[int]{
//	    Handler: func(data map[string]interface{}) *int {
//	        if val, ok := data["value"].(int); ok {
//	            return &val
//	        }
//	        return nil
//	    },
//	}
//	success := safeMap.ApplyFieldConfig("field1", config, data)
func (t *ThreadSafeHashMap[K, V]) ApplyFieldConfig(
	key K,
	config FieldConfig[V],
	data map[string]interface{},
) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.data.ApplyFieldConfig(key, config, data)
}

// ProcessFieldConfigs processes data using field configurations with a callback
// Example:
//
//	configs := map[string]FieldConfig[int]{
//	    "field1": {
//	        Handler: func(data map[string]interface{}) *int {
//	            if val, ok := data["value"].(int); ok {
//	                return &val
//	            }
//	            return nil
//	        },
//	    },
//	}
//	safeMap.ProcessFieldConfigs(configs, data, func(key string, value int, index int) {
//	    fmt.Printf("Processed: %s = %d at index %d\n", key, value, index)
//	})
func (t *ThreadSafeHashMap[K, V]) ProcessFieldConfigs(
	configs map[K]FieldConfig[V],
	data []map[string]interface{},
	processor func(key K, value V, index int),
) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data.ProcessFieldConfigs(configs, data, processor)
}
