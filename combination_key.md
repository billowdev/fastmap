# ComboKeyHashMap

ComboKeyHashMap is a specialized hash map implementation that allows you to associate values with combinations of keys. This is particularly useful when you need to map multiple related keys to a single value, such as in permission systems, feature flags, or any scenario where a combination of conditions should yield a specific result.

## Basic Usage

```go
import "github.com/billowdev/fastmap/hashmap"

// Create a new ComboKeyHashMap
m := fastmap.NewComboKeyHashMap[string, string]()

// Store a value with a combination of keys
m.Put([]string{"department_1", "role_admin"}, "FULL_ACCESS")

// Retrieve using exact key combination (order doesn't matter)
value, exists := m.Get([]string{"department_1", "role_admin"})
value, exists = m.Get([]string{"role_admin", "department_1"}) // Same result

// Get all values that match certain keys
results := m.GetByKeys([]string{"department_1"})
```

## Advanced Examples

### Permission System

```go
// Create a permission map
perms := fastmap.NewComboKeyHashMap[string, []string]()

// Define permissions for different combinations
perms.Put([]string{"finance", "manager"}, []string{"read", "write", "approve"})
perms.Put([]string{"finance", "staff"}, []string{"read", "write"})
perms.Put([]string{"finance", "staff", "junior"}, []string{"read"})

// Check permissions
if permissions, exists := perms.Get([]string{"finance", "manager"}); exists {
    // Use permissions
}

// Find all permissions containing "finance"
financePerms := perms.GetByKeys([]string{"finance"})
```

### Feature Flags with Multiple Conditions

```go
// Create a feature flag map
flags := fastmap.NewComboKeyHashMap[string, bool]()

// Set feature availability based on conditions
flags.Put([]string{"beta", "premium", "region_us"}, true)
flags.Put([]string{"beta", "premium", "region_eu"}, false)

// Check if feature is available
isAvailable, exists := flags.Get([]string{"beta", "premium", "region_us"})
```

### Content Access Control

```go
// Create an access control map
access := fastmap.NewComboKeyHashMap[string, string]()

// Define access levels for different combinations
access.Put([]string{"content_123", "subscription_premium"}, "full")
access.Put([]string{"content_123", "subscription_basic"}, "preview")
access.Put([]string{"content_123", "region_restricted"}, "blocked")

// Check access level
if level, exists := access.Get([]string{"content_123", "subscription_premium"}); exists {
    // Handle access level
}
```

## Key Concepts

### 1. Key Combinations
- Keys are treated as unordered combinations
- The same combination in any order returns the same value
- Each unique combination can map to only one value

### 2. Partial Matching
- `GetByKeys` returns all values where the provided keys are a subset of the stored combinations
- Useful for finding all relevant entries matching certain criteria

### 3. Overlapping Combinations
```go
m.Put([]string{"key1", "key2"}, "VALUE1")
m.Put([]string{"key1", "key2", "key3"}, "VALUE2")

// This will return both "VALUE1" and "VALUE2"
results := m.GetByKeys([]string{"key1", "key2"})
```

## Best Practices

1. **Key Design**
   - Use meaningful, consistent key names
   - Consider using enum-like constants for keys
   - Keep key combinations focused and logical

2. **Value Management**
   - Use appropriate value types for your use case
   - Consider using structs for complex values
   - Handle the exists/not-exists cases appropriately

3. **Performance**
   - Minimize the number of keys in combinations when possible
   - Use GetByKeys judiciously with large datasets
   - Consider using pointer values for large structs
