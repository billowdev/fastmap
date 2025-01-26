package fastmap

// ToMap returns the underlying map
// Example:
//
//	standardMap := hashMap.ToMap()
//	for k, v := range standardMap {
//	    fmt.Printf("%v: %v\n", k, v)
//	}
func (h *HashMap[K, V]) ToMap() map[K]V {
	result := make(map[K]V, len(h.data))
	for k, v := range h.data {
		result[k] = v
	}
	return result
}

// FromMap creates a new HashMap from a regular map
// Example:
//
//	regularMap := map[string]int{"one": 1, "two": 2}
//	hashMap := FromMap(regularMap)
func FromMap[K comparable, V any](m map[K]V) *HashMap[K, V] {
	h := NewHashMap[K, V]()
	for k, v := range m {
		h.Put(k, v)
	}
	return h
}
