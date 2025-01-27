package fastmap

import "sort"

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

// SortValues sorts the slice stored at the given key using a custom comparison function.
// Returns true if the key exists and sorting was performed, false otherwise.
// The less function should return true if a comes before b in the desired order.
//
// Parameters:
//   - key: The key whose values should be sorted
//   - less: Comparison function that defines the sort order
//
// Example:
//
//	map.AppendValues("scores", 75, 82, 90, 65)
//	map.SortValues("scores", func(a, b int) bool { return a < b })
//	// Result: [65, 75, 82, 90]
func (h *AppendableHashMap[K, V]) SortValues(key K, less func(a, b V) bool) bool {
	if values, exists := h.Get(key); exists {
		sort.Slice(values, func(i, j int) bool {
			return less(values[i], values[j])
		})
		h.Put(key, values)
		return true
	}
	return false
}

// SortValuesByField sorts values using a field extractor function that retrieves
// comparable values (int, string, float64) from the slice elements.
// Returns true if sorting was performed, false if key doesn't exist.
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	map.AppendValues("users", User{"Bob", 30}, User{"Alice", 25})
//	map.SortValuesByField("users", func(u User) any { return u.Name })
//	// Result: [{Alice 25} {Bob 30}]
func (h *AppendableHashMap[K, V]) SortValuesByField(key K, extractor func(V) any) bool {
	if values, exists := h.Get(key); exists {
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
		h.Put(key, values)
		return true
	}
	return false
}

// ReverseValues reverses the order of values stored at the given key.
// Returns true if reversal was performed, false if key doesn't exist.
//
// Example:
//
//	map.AppendValues("items", "a", "b", "c")
//	map.ReverseValues("items")
//	// Result: ["c", "b", "a"]
func (h *AppendableHashMap[K, V]) ReverseValues(key K) bool {
	if values, exists := h.Get(key); exists {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
		h.Put(key, values)
		return true
	}
	return false
}

// FilterValues applies a predicate function to filter values at the given key.
// Keeps only values where predicate returns true.
// Returns true if filtering was performed, false if key doesn't exist.
//
// Example:
//
//	map.AppendValues("numbers", 1, 2, 3, 4, 5)
//	map.FilterValues("numbers", func(n int) bool { return n%2 == 0 })
//	// Result: [2, 4]
func (h *AppendableHashMap[K, V]) FilterValues(key K, predicate func(V) bool) bool {
	if values, exists := h.Get(key); exists {
		filtered := make([]V, 0, len(values))
		for _, v := range values {
			if predicate(v) {
				filtered = append(filtered, v)
			}
		}
		h.Put(key, filtered)
		return true
	}
	return false
}

// TransformValues applies a transformation function to all values at the given key.
// The transform function maps each value to a new value of the same type.
// Returns true if transformation was performed, false if key doesn't exist.
//
// Example:
//
//	map.AppendValues("numbers", 1, 2, 3)
//	map.TransformValues("numbers", func(n int) int { return n * 2 })
//	// Result: [2, 4, 6]
func (h *AppendableHashMap[K, V]) TransformValues(key K, transform func(V) V) bool {
	if values, exists := h.Get(key); exists {
		for i, v := range values {
			values[i] = transform(v)
		}
		h.Put(key, values)
		return true
	}
	return false
}

// RemoveDuplicates removes duplicate values at the given key while maintaining order.
// Uses the provided equals function to determine value equality.
// Returns true if deduplication was performed, false if key doesn't exist.
//
// Example:
//
//	map.AppendValues("tags", "go", "java", "go", "python", "java")
//	map.RemoveDuplicates("tags", func(a, b string) bool { return a == b })
//	// Result: ["go", "java", "python"]
func (h *AppendableHashMap[K, V]) RemoveDuplicates(key K, equals func(a, b V) bool) bool {
	if values, exists := h.Get(key); exists {
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

		h.Put(key, unique)
		return true
	}
	return false
}
