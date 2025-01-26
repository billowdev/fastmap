package fastmap

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
