package fastmap

// ThreadSafeAppendableHashMap provides thread-safe operations for handling slice values
// in a concurrent environment. It uses mutex locks to ensure safe access and modification
// of the underlying data structure.
//
// Example:
//
//	safeMap := NewThreadSafeAppendableHashMap[string, Component]()
//	// Safe for concurrent access
//	go func() { safeMap.AppendValues("section1", component1) }()
//	go func() { safeMap.AppendValues("section1", component2) }()
type ThreadSafeAppendableHashMap[K comparable, V any] struct {
	*ThreadSafeHashMap[K, []V]
}

// NewThreadSafeAppendableHashMap creates a new instance of ThreadSafeAppendableHashMap
// that provides synchronized access to slice operations. It's suitable for concurrent
// environments where multiple goroutines might append values simultaneously.
//
// Example:
//
//	safeMap := NewThreadSafeAppendableHashMap[string, int]()
//	// Safe for concurrent operations
//	go func() { safeMap.AppendValues("key1", 1, 2) }()
func NewThreadSafeAppendableHashMap[K comparable, V any]() *ThreadSafeAppendableHashMap[K, V] {
	return &ThreadSafeAppendableHashMap[K, V]{
		ThreadSafeHashMap: NewThreadSafeHashMap[K, []V](),
	}
}

// AppendValues safely appends multiple values to an existing slice for a given key
// in a thread-safe manner. It uses mutex locks to ensure concurrent safety during
// the append operation.
//
// Parameters:
//   - key: The key to associate the values with
//   - values: Variadic parameter of values to append
//
// Example:
//
//	safeMap.AppendValues("users", user1, user2)
//	// Concurrent access is safe
//	go safeMap.AppendValues("users", user3)
func (t *ThreadSafeAppendableHashMap[K, V]) AppendValues(key K, values ...V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if existing, exists := t.data.Get(key); exists {
		t.data.Put(key, append(existing, values...))
	} else {
		t.data.Put(key, values)
	}
}
