package fastmap

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
//	err := safeMap.ForEach(func(key string, value User) error {
//	    if value.IsInvalid() {
//	        return fmt.Errorf("invalid user data for key %s", key)
//	    }
//	    fmt.Printf("User %s: %v\n", key, value)
//	    return nil
//	})
func (t *ThreadSafeHashMap[K, V]) ForEach(callback func(K, V) error) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.data.ForEach(callback)
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
