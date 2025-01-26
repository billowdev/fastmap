package fastmap

import "fmt"

// Clear removes all elements from the HashMap
// Example:
//
//	hashMap.Clear()
//	fmt.Printf("Size after clear: %d\n", hashMap.Size())
func (h *HashMap[K, V]) Clear() {
	h.data = make(map[K]V)
}

// Contains checks if a key exists in the HashMap
// Example:
//
//	if hashMap.Contains("user123") {
//	    fmt.Println("User exists")
//	}
func (h *HashMap[K, V]) Contains(key K) bool {
	_, exists := h.data[key]
	return exists
}

// IsEmpty returns true if the HashMap has no elements
// Example:
//
//	if hashMap.IsEmpty() {
//	    fmt.Println("HashMap is empty")
//	}
func (h *HashMap[K, V]) IsEmpty() bool {
	return len(h.data) == 0
}

// Keys returns a slice of all keys in the HashMap
// Example:
//
//	keys := hashMap.Keys()
//	for _, key := range keys {
//	    fmt.Printf("Key: %v\n", key)
//	}
func (h *HashMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(h.data))
	for k := range h.data {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values in the HashMap
// Example:
//
//	values := hashMap.Values()
//	for _, value := range values {
//	    fmt.Printf("Value: %v\n", value)
//	}
func (h *HashMap[K, V]) Values() []V {
	values := make([]V, 0, len(h.data))
	for _, v := range h.data {
		values = append(values, v)
	}
	return values
}

// ForEach executes a callback function for each key-value pair and returns an error if the callback fails
// Example:
//
//	err := hashMap.ForEach(func(key string, value User) error {
//	    if value.IsInvalid() {
//	        return fmt.Errorf("invalid user data for key %s", key)
//	    }
//	    fmt.Printf("User %s: %v\n", key, value)
//	    return nil
//	})
func (h *HashMap[K, V]) ForEach(callback func(K, V) error) error {
	for k, v := range h.data {
		if err := callback(k, v); err != nil {
			return fmt.Errorf("ForEach operation failed at key %v: %w", k, err)
		}
	}
	return nil
}

// UpdateValue updates an existing value by key, returns false if key doesn't exist
// Example:
//
//	if hashMap.UpdateValue("user123", updatedUser) {
//	    fmt.Println("User updated successfully")
//	}
func (h *HashMap[K, V]) UpdateValue(id K, newValue V) bool {
	if _, exists := h.data[id]; exists {
		h.data[id] = newValue
		return true
	}
	return false
}

// PutAll adds all key-value pairs from another HashMap
// Example:
//
//	otherMap := NewHashMap[string, User]()
//	otherMap.Put("user456", newUser)
//	hashMap.PutAll(otherMap)
func (h *HashMap[K, V]) PutAll(other *HashMap[K, V]) {
	for k, v := range other.data {
		h.data[k] = v
	}
}
