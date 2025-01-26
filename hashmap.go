package fastmap

type ValueCallback[K comparable, V any] func(key K, value V)

// HashMap is a generic key-value store supporting comparable keys and any value type
// Example:
//
//	hashMap := NewHashMap[string, int]()
type HashMap[K comparable, V any] struct {
	data map[K]V
}

// NewHashMap creates a new empty HashMap
// Example:
//
//	hashMap := NewHashMap[string, User]()
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		data: make(map[K]V),
	}
}

// Put adds or updates a key-value pair in the HashMap
// Example:
//
//	hashMap.Put("user123", User{Name: "John"})
func (h *HashMap[K, V]) Put(key K, value V) {
	h.data[key] = value
}

// Get retrieves a value by key and returns whether it exists
// Example:
//
//	if user, exists := hashMap.Get("user123"); exists {
//	    fmt.Printf("Found user: %v\n", user)
//	}
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	value, exists := h.data[key]
	return value, exists
}

// Remove deletes a key-value pair from the HashMap
// Example:
//
//	hashMap.Remove("user123")
func (h *HashMap[K, V]) Remove(key K) {
	delete(h.data, key)
}

// Size returns the number of elements in the HashMap
// Example:
//
//	count := hashMap.Size()
//	fmt.Printf("HashMap contains %d elements\n", count)
func (h *HashMap[K, V]) Size() int {
	return len(h.data)
}
