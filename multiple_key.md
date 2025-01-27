# MultiKeyHashMap and ThreadSafeMultiKeyHashMap

## MultiKeyHashMap

A hash map implementation that supports multiple keys (aliases) for accessing the same value.

### Basic Usage

```go
// Create new map
m := fastmap.NewMultiKeyHashMap[string, User]()

// Add with multiple keys (first key is primary)
m.Put([]string{"user1", "john", "john.doe"}, user)

// Get using any key
user, exists := m.Get("john.doe")

// Get primary key
primaryKey, exists := m.GetPrimaryKey("john")

// Get all associated keys
keys := m.GetAllKeys("john")

// Add new alias
m.AddAlias("user1", "jdoe")

// Remove specific key
m.Remove("john.doe")

// Remove key and all connected aliases
m.RemoveWithCascade("user1")

// Get unique entry count
size := m.Size()

// Clear all entries
m.Clear()
```

## ThreadSafeMultiKeyHashMap

Thread-safe version with mutex protection for concurrent access.

### Thread-Safe Usage

```go
// Create thread-safe map
m := fastmap.NewThreadSafeMultiKeyHashMap[string, User]()

// Operations are same as MultiKeyHashMap but thread-safe
m.Put([]string{"user1", "john", "john.doe"}, user)
user, exists := m.Get("john.doe")
m.AddAlias("user1", "jdoe")
m.RemoveWithCascade("user1")

// Concurrent usage example
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        key := fmt.Sprintf("user%d", id)
        m.Put([]string{key, fmt.Sprintf("alias%d", id)}, 
            User{ID: id})
    }(i)
}
wg.Wait()
```

### Performance Notes

- Read operations use RLock for concurrent access
- Write operations use exclusive Lock
- Primary key lookups are O(1)
- Alias lookups require map iteration in worst case
- Memory scales with number of aliases