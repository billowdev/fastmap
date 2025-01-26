package fastmap

// Filter returns a new HashMap containing only the elements that satisfy the predicate
// Example:
//
//	activeUsers := hashMap.Filter(func(key string, user User) bool {
//	    return user.Active
//	})
func (h *HashMap[K, V]) Filter(predicate func(K, V) bool) *HashMap[K, V] {
	result := NewHashMap[K, V]()
	for k, v := range h.data {
		if predicate(k, v) {
			result.Put(k, v)
		}
	}
	return result
}

// Map transforms values using the provided function and returns a new HashMap
// Example:
//
//	upperNames := hashMap.Map(func(key string, user User) User {
//	    user.Name = strings.ToUpper(user.Name)
//	    return user
//	})
func (h *HashMap[K, V]) Map(transform func(K, V) V) *HashMap[K, V] {
	result := NewHashMap[K, V]()
	for k, v := range h.data {
		result.Put(k, transform(k, v))
	}
	return result
}

