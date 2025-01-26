# `fastmap` Generic HashMap Package

A type-safe, generic HashMap implementation in Go with both thread-safe and non-thread-safe variants.

## Installation

```bash
go get github.com/billowdev/fastmap
```

## Features

- Generic type support for keys and values
- Thread-safe implementation available
- Type-safe operations
- Functional programming utilities
- Zero dependencies
- Comprehensive test coverage

## Usage

### Basic Operations (Non-Thread-Safe)

```go
// Initialize
hashMap := fastmap.NewHashMap[string, YourType]()

// Add/Update
hashMap.Put("key", value)

// Get
value, exists := hashMap.Get("key")

// Remove
hashMap.Remove("key")

// Get size
size := hashMap.Size()

// Update existing value
success := hashMap.UpdateValue("key", newValue)
```

### Thread-Safe Operations

```go
// Initialize thread-safe map
safeMap := fastmap.NewThreadSafeHashMap[string, YourType]()

// Basic Operations
safeMap.Put("key", value)
value, exists := safeMap.Get("key")
safeMap.Remove("key")

// Safe iteration
safeMap.ForEach(func(key string, value YourType) {
    // Process each key-value pair safely
})

// Check existence
if safeMap.Contains("key") {
    // Key exists, safe for concurrent access
}

// Get all keys safely
keys := safeMap.Keys()
for _, key := range keys {
    // Process keys
}

// Get all values safely
values := safeMap.Values()
for _, value := range values {
    // Process values
}

// Clear all entries safely
safeMap.Clear()

// Check if empty
if safeMap.IsEmpty() {
    // Map is empty
}

// Merge two thread-safe maps
otherMap := fastmap.NewThreadSafeHashMap[string, YourType]()
otherMap.Put("other", value)
safeMap.PutAll(otherMap)

// Convert to regular map
regularMap := safeMap.ToMap()

// Create from regular map
traditional := map[string]YourType{"key": value}
safeMap = fastmap.FromThreadSafeMap(traditional)
```

### Thread-Safe Functional Operations

```go
// Filter with thread safety
activeUsers := safeMap.Filter(func(key string, user User) bool {
    return user.Active
})

// Transform with thread safety
processedUsers := safeMap.Map(func(key string, user User) User {
    user.LastProcessed = time.Now()
    return user
})

// Conditional updates with thread safety
if safeMap.UpdateValue("user1", updatedUser) {
    // Update successful
}
```

### Real-World Example (Thread-Safe)

```go
// User management system with concurrent access
type UserSystem struct {
    users *fastmap.ThreadSafeHashMap[string, User]
}

func NewUserSystem() *UserSystem {
    return &UserSystem{
        users: fastmap.NewThreadSafeHashMap[string, User](),
    }
}

// Safe concurrent operations
func (s *UserSystem) AddUser(id string, user User) {
    s.users.Put(id, user)
}

func (s *UserSystem) GetActiveUsers() []User {
    activeUsers := s.users.Filter(func(id string, user User) bool {
        return user.Active && !user.Deleted
    })
    return activeUsers.Values()
}

func (s *UserSystem) UpdateUserStatus(id string, active bool) bool {
    if user, exists := s.users.Get(id); exists {
        user.Active = active
        return s.users.UpdateValue(id, user)
    }
    return false
}

func (s *UserSystem) ProcessUsers() {
    s.users.ForEach(func(id string, user User) {
        // Safe concurrent processing
        log.Printf("Processing user: %s", id)
    })
}

// Usage in concurrent environment
func main() {
    system := NewUserSystem()
    
    // Concurrent operations
    go func() {
        system.AddUser("1", User{Name: "John", Active: true})
    }()
    
    go func() {
        system.UpdateUserStatus("1", false)
    }()
    
    go func() {
        activeUsers := system.GetActiveUsers()
        for _, user := range activeUsers {
            log.Printf("Active user: %s", user.Name)
        }
    }()
}
```

### Other
```go
hashMap := fastmap.NewHashMap[string, string]()

configs := map[string]fastmap.FieldConfig[string]{
    "minimum_polarisation": {
        Handler: func(m map[string]interface{}) *string {
            return utils.ToPtr(utils.GetStringValueFromMap(m, "minimum_polarisation"))
        },
    },
}

hashMap.ProcessFieldConfigs(configs, productSpecifications, func(key string, value string, index int) {
    // Handle processed field
    element := models.DocumentFieldValue{
        FieldID:   key,
        Value:     &value,
        VersionID: versionID,
    }
    bulkCreateValue = append(bulkCreateValue, element)
})
```

## Migration Guide

### Traditional Map vs HashMap

```go
// Traditional
traditional := make(map[string]YourType)
traditional[key] = value

// Non-Thread-Safe HashMap
hashMap := fastmap.NewHashMap[string, YourType]()
hashMap.Put(key, value)

// Thread-Safe HashMap
safeMap := fastmap.NewThreadSafeHashMap[string, YourType]()
safeMap.Put(key, value)
```

## Performance

- Built on Go's native map implementation
- Minimal method call overhead
- O(1) average case for basic operations
- Thread-safe operations use sync.RWMutex
- Read operations allow concurrent access
- Write operations ensure exclusive access

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License