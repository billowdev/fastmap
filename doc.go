// Package fastmap provides a generic, type-safe HashMap implementation with both thread-safe and non-thread-safe options.
//
// Key features:
//   - Generic type support for keys and values
//   - Type-safe operations through Go's type system
//   - Functional programming utilities (Map, Filter, etc.)
//   - Zero external dependencies
//   - Efficient memory usage leveraging Go's built-in map
//   - Thread-safe implementation available via ThreadSafeHashMap
//
// Basic usage (non-thread-safe):
//
//	// Create a new HashMap
//	hashMap := fastmap.NewHashMap[string, User]()
//
//	// Add elements
//	hashMap.Put("user1", User{Name: "John"})
//
//	// Get elements
//	if user, exists := hashMap.Get("user1"); exists {
//	    fmt.Printf("Found user: %v\n", user)
//	}
//
// Thread-safe usage:
//
//	// Create a thread-safe HashMap
//	safeMap := fastmap.NewThreadSafeHashMap[string, User]()
//
//	// Safe concurrent operations
//	safeMap.Put("user1", User{Name: "John"})
//	
//	// Safe concurrent reads
//	if user, exists := safeMap.Get("user1"); exists {
//	    fmt.Printf("Found user: %v\n", user)
//	}
//
// Functional operations (both variants):
//
//	// Filter users
//	activeUsers := hashMap.Filter(func(key string, user User) bool {
//	    return user.Active
//	})
//
//	// Transform users
//	upperNames := hashMap.Map(func(key string, user User) User {
//	    user.Name = strings.ToUpper(user.Name)
//	    return user
//	})
//
// Converting from traditional maps:
//
//	// To thread-safe HashMap
//	traditional := map[string]int{"one": 1, "two": 2}
//	safeMap := fastmap.FromThreadSafeMap(traditional)
//
//	// To regular HashMap
//	hashMap := fastmap.FromMap(traditional)
//
// Performance considerations:
//   - Uses Go's built-in map implementation under the hood
//   - Method calls have minimal overhead
//   - Memory usage equivalent to standard Go maps plus synchronization primitives
//   - O(1) average case for basic operations (Get, Put, Remove)
//   - O(n) for operations that traverse all elements (ForEach, Filter, Map)
//   - Thread-safe operations may have additional overhead due to mutex locking
//
// Thread safety considerations:
//   - HashMap is not thread-safe and should not be used concurrently
//   - ThreadSafeHashMap provides full thread safety through sync.RWMutex
//   - Read operations use RLock for concurrent access
//   - Write operations use Lock for exclusive access
//   - Nested operations (like PutAll) handle multiple locks correctly
//
// For more examples and detailed API documentation, see the package tests
// and individual method documentation.
//
// Package repository: https://github.com/billowdev/fastmap
package fastmap