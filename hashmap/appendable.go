package fastmap

// AppendableHashMap extends HashMap to provide specialized functionality for handling slice operations.
// It maintains type safety through generics while providing convenient methods for appending values
// to existing slices within the map.
//
// Example:
//
//	layoutMap := NewAppendableHashMap[string, Component]()
//	layoutMap.AppendValues("section1", component1, component2)
//	layoutMap.AppendValues("section1", component3) // Appends to existing slice
type AppendableHashMap[K comparable, V any] struct {
	*HashMap[K, []V]
}

// NewAppendableHashMap creates a new instance of AppendableHashMap that safely manages
// slices of values associated with keys. It initializes an underlying HashMap to store
// the key-value pairs where values are slices.
//
// Example:
//
//	map := NewAppendableHashMap[string, int]()
//	map.AppendValues("key1", 1, 2, 3)
func NewAppendableHashMap[K comparable, V any]() *AppendableHashMap[K, V] {
	return &AppendableHashMap[K, V]{
		HashMap: NewHashMap[K, []V](),
	}
}

// AppendValues appends multiple values to an existing slice for a given key.
// If the key doesn't exist, it creates a new slice with the provided values.
// This method provides a safe way to handle the spread operator equivalent in Go.
//
// Parameters:
//   - key: The key to associate the values with
//   - values: Variadic parameter of values to append
//
// Example:
//
//	map.AppendValues("components", component1, component2)
//	map.AppendValues("components", component3) // Appends to existing slice
func (h *AppendableHashMap[K, V]) AppendValues(key K, values ...V) {
	if existing, exists := h.Get(key); exists {
		h.Put(key, append(existing, values...))
	} else {
		h.Put(key, values)
	}
}
