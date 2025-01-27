package fastmap

import "sort"

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

// SortValues sorts the slice in a thread-safe manner using the custom comparison function.
// Returns true if sorting was performed, false if key doesn't exist.
//
// Example:
//
//	safeMap.AppendValues("scores", 75, 82, 90, 65)
//	safeMap.SortValues("scores", func(a, b int) bool { return a < b })
//	// Result: [65, 75, 82, 90]
func (t *ThreadSafeAppendableHashMap[K, V]) SortValues(key K, less func(a, b V) bool) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if values, exists := t.Get(key); exists {
		sort.Slice(values, func(i, j int) bool {
			return less(values[i], values[j])
		})
		t.Put(key, values)
		return true
	}
	return false
}

// SortValuesByField sorts values in a thread-safe manner using a field extractor.
// Returns true if sorting was performed, false if key doesn't exist.
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	safeMap.AppendValues("users", User{"Bob", 30}, User{"Alice", 25})
//	safeMap.SortValuesByField("users", func(u User) any { return u.Name })
//	// Result: [{Alice 25} {Bob 30}]
func (t *ThreadSafeAppendableHashMap[K, V]) SortValuesByField(key K, extractor func(V) any) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if values, exists := t.Get(key); exists {
		sort.Slice(values, func(i, j int) bool {
			a := extractor(values[i])
			b := extractor(values[j])
			switch v := a.(type) {
			case int:
				return v < b.(int)
			case string:
				return v < b.(string)
			case float64:
				return v < b.(float64)
			default:
				return false
			}
		})
		t.Put(key, values)
		return true
	}
	return false
}

// ReverseValues reverses values order in a thread-safe manner.
// Returns true if reversal was performed, false if key doesn't exist.
//
// Example:
//
//	safeMap.AppendValues("items", "a", "b", "c")
//	safeMap.ReverseValues("items")
//	// Result: ["c", "b", "a"]
func (t *ThreadSafeAppendableHashMap[K, V]) ReverseValues(key K) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if values, exists := t.Get(key); exists {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
		t.Put(key, values)
		return true
	}
	return false
}

// FilterValues filters values in a thread-safe manner using a predicate.
// Returns true if filtering was performed, false if key doesn't exist.
//
// Example:
//
//	safeMap.AppendValues("numbers", 1, 2, 3, 4, 5)
//	safeMap.FilterValues("numbers", func(n int) bool { return n%2 == 0 })
//	// Result: [2, 4]
func (t *ThreadSafeAppendableHashMap[K, V]) FilterValues(key K, predicate func(V) bool) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if values, exists := t.Get(key); exists {
		filtered := make([]V, 0, len(values))
		for _, v := range values {
			if predicate(v) {
				filtered = append(filtered, v)
			}
		}
		t.Put(key, filtered)
		return true
	}
	return false
}

// TransformValues transforms values in a thread-safe manner.
// Returns true if transformation was performed, false if key doesn't exist.
//
// Example:
//
//	safeMap.AppendValues("numbers", 1, 2, 3)
//	safeMap.TransformValues("numbers", func(n int) int { return n * 2 })
//	// Result: [2, 4, 6]
func (t *ThreadSafeAppendableHashMap[K, V]) TransformValues(key K, transform func(V) V) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if values, exists := t.Get(key); exists {
		for i, v := range values {
			values[i] = transform(v)
		}
		t.Put(key, values)
		return true
	}
	return false
}

// RemoveDuplicates removes duplicates in a thread-safe manner.
// Returns true if deduplication was performed, false if key doesn't exist.
//
// Example:
//
//	safeMap.AppendValues("tags", "go", "java", "go", "python", "java")
//	safeMap.RemoveDuplicates("tags", func(a, b string) bool { return a == b })
//	// Result: ["go", "java", "python"]
func (t *ThreadSafeAppendableHashMap[K, V]) RemoveDuplicates(key K, equals func(a, b V) bool) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if values, exists := t.Get(key); exists {
		if len(values) <= 1 {
			return true
		}

		unique := make([]V, 0, len(values))
		unique = append(unique, values[0])

		for i := 1; i < len(values); i++ {
			isDuplicate := false
			for j := 0; j < len(unique); j++ {
				if equals(values[i], unique[j]) {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				unique = append(unique, values[i])
			}
		}

		t.Put(key, unique)
		return true
	}
	return false
}
