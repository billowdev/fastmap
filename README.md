# FastMap - A Type-Safe Generic HashMap Implementation in Go

FastMap currently provides an efficient, type-safe HashMap implementation in Go, offering both thread-safe and non-thread-safe variants. Future releases will expand this into a comprehensive utilities package.

Current Feature Implementation Status:

✅ Core HashMap Implementation
- Basic HashMap operations (Put, Get, Remove)
- Type-safe implementation using Go generics
- Full error handling
- Comprehensive unit tests
- Performance-optimized operations

✅ Thread-Safe HashMap
- Mutex-protected operations
- Concurrent access support
- Thread-safe variants of all core operations
- Deadlock prevention
- Performance benchmarks for concurrent operations

✅ AppendableHashMap
- Specialized slice handling
- Type-safe append operations
- Thread-safe variant available
- Optimized memory management
- Comprehensive testing coverage

✅ Functional Operations
- Map transformations
- Filtering capabilities
- ForEach operations with error handling
- Chainable operations
- Performance benchmarks

<!-- Planned Future Enhancements:
	- RobinHood HashMap Implementation
	- Binary Search Tree Implementation
	- Struct-to-Struct Conversion
	- Case Conversion Utilities
	- Additional Data Structure Support
	- Extended Utility Functions
	- Advanced Type Conversion Tools
	- Deep Copy Functionality -->
 
The current version focuses on providing a robust and efficient HashMap implementation. Future releases will expand the package's capabilities to include additional data structures and utility functions.


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/billowdev/fastmap.svg?style=for-the-badge
[contributors-url]: https://github.com/billowdev/fastmap/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/billowdev/fastmap.svg?style=for-the-badge
[forks-url]: https://github.com/billowdev/fastmap/network/members
[stars-shield]: https://img.shields.io/github/stars/billowdev/fastmap.svg?style=for-the-badge
[stars-url]: https://github.com/billowdev/fastmap/stargazers
[issues-shield]: https://img.shields.io/github/issues/billowdev/fastmap.svg?style=for-the-badge
[issues-url]: https://github.com/billowdev/fastmap/issues
[license-shield]: https://img.shields.io/github/license/billowdev/fastmap.svg?style=for-the-badge
[license-url]: https://github.com/billowdev/fastmap/blob/main/LICENSE


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

- AppendValues
```go
// processSections handles the conversion and organization of section data into specialized
// hash maps for both body content and layout components. It demonstrates the use of
// AppendableHashMap for managing collections of PDF components per section.
//
// Parameters:
//   - listSections: Slice of Section objects containing section data and components
//
// Example:
//
//	sections := []Section{
//	    {
//	        ID: "section1",
//	        PDFListComponents: []models.PDFListComponent{comp1, comp2},
//	    },
//	}
//	processSections(sections)
func processSections(listSections []Section) {
    // Initialize specialized hash maps for different data types
    bodyHashMap := fastmap.NewHashMap[string, domain.ResSection]()
    layoutHashMap := fastmap.NewAppendableHashMap[string, models.PDFListComponent]()
    
    // Process each section and organize its data into appropriate maps
    for _, section := range listSections {
        // Store section metadata in body hash map
        bodyHashMap.Put(section.ID, domain.ResSection{
            SectionID: section.ID,
            Section:   section.Section,
            Priority:  section.Priority,
            Title:     string(section.Title),
            Elements:  nil,
        })
        
        // Append PDF components to layout hash map using spread operator equivalent
        layoutHashMap.AppendValues(section.ID, section.PDFListComponents...)
    }
    
    // Example of thread-safe implementation if needed for concurrent access
    safeLayoutHashMap := fastmap.NewThreadSafeAppendableHashMap[string, models.PDFListComponent]()
    for _, section := range listSections {
        safeLayoutHashMap.AppendValues(section.ID, section.PDFListComponents...)
    }
}
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
safeMap.ForEach(func(key string, value YourType) error {
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
    s.users.ForEach(func(id string, user User) error {
        // Safe concurrent processing
        log.Printf("Processing user: %s", id)
        return nil
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

- Append Values
```go
package main

import (
    "log"
    "sync"
    "time"

    "github.com/billowdev/fastmap/hashmap"
)

// PDFProcessor represents a system that processes PDF components concurrently
type PDFProcessor struct {
    layoutMap *fastmap.ThreadSafeAppendableHashMap[string, PDFComponent]
    wg        sync.WaitGroup
}

type PDFComponent struct {
    ID        string
    Content   string
    Timestamp time.Time
}

// NewPDFProcessor initializes a new PDF processor with thread-safe storage
func NewPDFProcessor() *PDFProcessor {
    return &PDFProcessor{
        layoutMap: fastmap.NewThreadSafeAppendableHashMap[string, PDFComponent](),
    }
}

// ProcessSection handles concurrent processing of PDF components for a section
func (p *PDFProcessor) ProcessSection(sectionID string, components []PDFComponent) {
    batchSize := 5
    for i := 0; i < len(components); i += batchSize {
        end := i + batchSize
        if end > len(components) {
            end = len(components)
        }

        batch := components[i:end]
        p.wg.Add(1)
        go p.processBatch(sectionID, batch)
    }
}

// processBatch handles a batch of components concurrently
func (p *PDFProcessor) processBatch(sectionID string, components []PDFComponent) {
    defer p.wg.Done()

    // Simulate processing time for each component
    for _, component := range components {
        // Simulate some processing work
        time.Sleep(100 * time.Millisecond)
        
        // Safely append the processed component
        p.layoutMap.AppendValues(sectionID, component)
        log.Printf("Processed component %s for section %s", component.ID, sectionID)
    }
}

// GetProcessedComponents safely retrieves all components for a section
func (p *PDFProcessor) GetProcessedComponents(sectionID string) []PDFComponent {
    components, exists := p.layoutMap.Get(sectionID)
    if !exists {
        return []PDFComponent{}
    }
    return components
}

// WaitForCompletion waits for all processing to complete
func (p *PDFProcessor) WaitForCompletion() {
    p.wg.Wait()
}

// Usage example
func main() {
    processor := NewPDFProcessor()

    // Simulate incoming PDF components for multiple sections
    sections := map[string][]PDFComponent{
        "section1": generateComponents("section1", 15),
        "section2": generateComponents("section2", 10),
        "section3": generateComponents("section3", 20),
    }

    // Process sections concurrently
    startTime := time.Now()
    for sectionID, components := range sections {
        processor.ProcessSection(sectionID, components)
    }

    // Wait for all processing to complete
    processor.WaitForCompletion()
    log.Printf("Processing completed in %v", time.Since(startTime))

    // Verify results
    for sectionID := range sections {
        processed := processor.GetProcessedComponents(sectionID)
        log.Printf("Section %s has %d processed components", sectionID, len(processed))
    }
}

// generateComponents creates test PDF components
func generateComponents(sectionID string, count int) []PDFComponent {
    components := make([]PDFComponent, count)
    for i := 0; i < count; i++ {
        components[i] = PDFComponent{
            ID:        fmt.Sprintf("%s-comp%d", sectionID, i),
            Content:   fmt.Sprintf("Content %d", i),
            Timestamp: time.Now(),
        }
    }
    return components
}

// Example of error handling and recovery
func (p *PDFProcessor) ProcessSectionWithRecovery(sectionID string, components []PDFComponent) error {
    errorChan := make(chan error, 1)
    
    go func() {
        defer func() {
            if r := recover(); r != nil {
                errorChan <- fmt.Errorf("processing panic: %v", r)
            }
            close(errorChan)
        }()

        p.ProcessSection(sectionID, components)
    }()

    // Wait for completion or error
    p.wg.Wait()
    if err := <-errorChan; err != nil {
        return fmt.Errorf("section %s processing failed: %w", sectionID, err)
    }

    return nil
}
```

### Config Field Processing

The fastmap package provides robust field configuration processing capabilities for handling dynamic data transformations. This feature is particularly useful when dealing with structured data that needs type-safe conversion and validation.

#### Basic Usage

```go
// Initialize HashMap
hashMap := fastmap.NewHashMap[string, int]()
hashMap.Put("age", 0) // Initialize field

// Define field configurations
configs := map[string]fastmap.FieldConfig[int]{
    "age": {
        Handler: func(data map[string]interface{}) *int {
            if val, ok := data["age"].(float64); ok {
                intVal := int(val)
                return &intVal
            }
            return nil
        },
    },
}

// Process data
data := []map[string]interface{}{
    {"age": 25.0},
    {"age": 30.0},
}

// Method 1: Handle single field
results := hashMap.HandleFieldConfigs(data, configs, "age")
// results = []int{25, 30}

// Method 2: Apply single config
success := hashMap.ApplyFieldConfig("age", configs["age"], data[0])
// success = true, hashMap["age"] = 25

// Method 3: Process all configs with callback
hashMap.ProcessFieldConfigs(configs, data, func(key string, value int, index int) {
    fmt.Printf("Processed %s: %d at index %d\n", key, value, index)
})
```

#### Thread-Safe Processing

1. Basic Thread-Safe Operations
```go
// Initialize thread-safe map
safeMap := fastmap.NewThreadSafeHashMap[string, float64]()
safeMap.Put("temperature", 0.0)

// Configure field handlers
configs := map[string]fastmap.FieldConfig[float64]{
    "temperature": {
        Handler: func(data map[string]interface{}) *float64 {
            if val, ok := data["temp"].(float64); ok {
                return &val
            }
            return nil
        },
    },
}

// Safe concurrent processing
safeMap.ProcessFieldConfigs(configs, data, func(key string, value float64, index int) {
    log.Printf("Temperature reading %f at index %d", value, index)
})
```

2. Concurrent Data Processing
```go
// Initialize thread-safe map with complex config
type Measurement struct {
    Value     float64
    Timestamp time.Time
    Valid     bool
}

safeMap := fastmap.NewThreadSafeHashMap[string, Measurement]()
rowIndex := 0

configs := map[string]fastmap.FieldConfig[Measurement]{
    "sensor_data": {
        RowIndex: &rowIndex,
        Handler: func(data map[string]interface{}) *Measurement {
            if val, ok := data["value"].(float64); ok {
                return &Measurement{
                    Value:     val,
                    Timestamp: time.Now(),
                    Valid:     val >= 0 && val <= 100,
                }
            }
            return nil
        },
    },
}

// Concurrent processing with error handling
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        sensorData := []map[string]interface{}{
            {"value": float64(index * 10)},
        }
        safeMap.ProcessFieldConfigs(configs, sensorData, func(key string, m Measurement, idx int) {
            if m.Valid {
                log.Printf("Valid measurement %f at time %v", m.Value, m.Timestamp)
            }
        })
    }(i)
}
wg.Wait()
```

3. Batch Processing with Multiple Fields
```go
type ProductData struct {
    Price    float64
    Quantity int
    Total    float64
}

safeMap := fastmap.NewThreadSafeHashMap[string, ProductData]()
configs := map[string]fastmap.FieldConfig[ProductData]{
    "product": {
        Handler: func(data map[string]interface{}) *ProductData {
            price, ok1 := data["price"].(float64)
            qty, ok2 := data["quantity"].(float64)
            if !ok1 || !ok2 {
                return nil
            }
            return &ProductData{
                Price:    price,
                Quantity: int(qty),
                Total:    price * float64(int(qty)),
            }
        },
    },
}

// Process batch data concurrently
batchData := []map[string]interface{}{
    {"price": 10.5, "quantity": 2.0},
    {"price": 20.0, "quantity": 3.0},
}

safeMap.ProcessFieldConfigs(configs, batchData, func(key string, pd ProductData, index int) {
    log.Printf("Processed product at index %d: Total = %.2f", index, pd.Total)
})
```

4. Error Handling in Thread-Safe Context
```go
safeMap := fastmap.NewThreadSafeHashMap[string, int]()
rowIndex := 0

configs := map[string]fastmap.FieldConfig[int]{
    "quantity": {
        RowIndex: &rowIndex,
        Handler: func(data map[string]interface{}) *int {
            val, ok := data["quantity"]
            if !ok {
                log.Printf("Missing quantity field at row %d", *rowIndex)
                return nil
            }
            if floatVal, ok := val.(float64); ok {
                intVal := int(floatVal)
                if intVal < 0 {
                    log.Printf("Invalid negative quantity at row %d", *rowIndex)
                    return nil
                }
                return &intVal
            }
            return nil
        },
    },
}

// Process with validation
safeMap.ProcessFieldConfigs(configs, data, func(key string, quantity int, index int) {
    log.Printf("Processed quantity %d at index %d", quantity, index)
})
```

#### Advanced Features

1. Row Index Tracking
```go
rowIndex := 0
configs := map[string]fastmap.FieldConfig[string]{
    "name": {
        RowIndex: &rowIndex,
        Handler: func(data map[string]interface{}) *string {
            if val, ok := data["name"].(string); ok {
                return &val
            }
            return nil
        },
    },
}
```

2. Complex Type Handling
```go
type UserData struct {
    Name  string
    Age   int
    Score float64
}

configs := map[string]fastmap.FieldConfig[UserData]{
    "user_info": {
        Handler: func(data map[string]interface{}) *UserData {
            if name, ok := data["name"].(string); ok {
                if age, ok := data["age"].(float64); ok {
                    if score, ok := data["score"].(float64); ok {
                        return &UserData{
                            Name:  name,
                            Age:   int(age),
                            Score: score,
                        }
                    }
                }
            }
            return nil
        },
    },
}
```

3. Batch Processing with Multiple Configs
```go
configs := map[string]fastmap.FieldConfig[string]{
    "name": {Handler: nameHandler},
    "email": {Handler: emailHandler},
    "phone": {Handler: phoneHandler},
}

hashMap.ProcessFieldConfigs(configs, bulkData, func(key string, value string, index int) {
    switch key {
    case "name":
        processName(value, index)
    case "email":
        processEmail(value, index)
    case "phone":
        processPhone(value, index)
    }
})
```

4. Field Dependencies
```go
type ProductSpec struct {
    Price    float64
    Quantity int
    Total    float64
}

configs := map[string]fastmap.FieldConfig[ProductSpec]{
    "product": {
        Handler: func(data map[string]interface{}) *ProductSpec {
            price, ok1 := data["price"].(float64)
            qty, ok2 := data["quantity"].(float64)
            if !ok1 || !ok2 {
                return nil
            }
            return &ProductSpec{
                Price:    price,
                Quantity: int(qty),
                Total:    price * qty,
            }
        },
    },
}
```
## Testing

### Unit Tests

Run all unit tests:
```bash
go test ./...
```

Run specific test:
```bash
go test -run TestHashMap_Put
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Key Test Categories

1. Basic Operations
```go
func TestHashMap_Put(t *testing.T) {
    h := fastmap.NewHashMap[string, int]()
    h.Put("key", 100)
    if val, exists := h.Get("key"); !exists || val != 100 {
        t.Errorf("Put failed, got (%v, %v), want (100, true)", val, exists)
    }
}
```

2. Thread-Safe Operations
```go
func TestThreadSafeConcurrentOperations(t *testing.T) {
    m := fastmap.NewThreadSafeHashMap[string, int]()
    var wg sync.WaitGroup
    numGoroutines := 100

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(val int) {
            defer wg.Done()
            m.Put(fmt.Sprintf("key%d", val), val)
        }(i)
    }
    wg.Wait()

    if m.Size() != numGoroutines {
        t.Errorf("Expected size %d, got %d", numGoroutines, m.Size())
    }
}
```

3. Functional Operations
```go
func TestFilter(t *testing.T) {
    h := fastmap.NewHashMap[string, int]()
    h.Put("one", 1)
    h.Put("two", 2)
    
    filtered := h.Filter(func(k string, v int) bool {
        return v%2 == 0
    })
    
    if filtered.Size() != 1 {
        t.Error("Filter failed")
    }
}
```

4. Edge Cases
```go
func TestEdgeCases(t *testing.T) {
    h := fastmap.NewHashMap[string, *string]()
    var nilStr *string
    
    h.Put("nilKey", nilStr)
    if val, exists := h.Get("nilKey"); !exists || val != nil {
        t.Error("Failed to handle nil value")
    }
}
```

### Benchmarks

Run all benchmarks:
```bash
go test -bench=. ./...
```

Run specific benchmark:
```bash
go test -bench=BenchmarkHashMapPut
```

Run benchmarks with memory allocation statistics:
```bash
go test -bench=. -benchmem ./...
```

#### Key Benchmark Categories

1. Basic Operations
```go
func BenchmarkHashMapPut(b *testing.B) {
    h := fastmap.NewHashMap[string, int]()
    for i := 0; i < b.N; i++ {
        h.Put("key", i)
    }
}
```

2. Thread-Safe Operations
```go
func BenchmarkThreadSafePut(b *testing.B) {
    m := fastmap.NewThreadSafeHashMap[string, int]()
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            m.Put(fmt.Sprintf("key%d", i), i)
            i++
        }
    })
}
```

3. Functional Operations
```go
func BenchmarkHashMapFilter(b *testing.B) {
    h := fastmap.NewHashMap[string, int]()
    for i := 0; i < 1000; i++ {
        h.Put(fmt.Sprintf("key%d", i), i)
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        h.Filter(func(k string, v int) bool {
            return v%2 == 0
        })
    }
}
```

#### Benchmark Tips

1. Use realistic data sizes
2. Include parallel benchmarks for thread-safe operations
3. Test with different data types and structures
4. Measure memory allocations for memory-sensitive operations
5. Compare performance with standard library map operations

#### Common Benchmark Flags

```bash
-bench=.                # Run all benchmarks
-benchmem              # Print memory allocation statistics
-benchtime=10s         # Run each benchmark for 10 seconds
-count=5               # Run each benchmark 5 times
-cpu=1,2,4            # Run benchmarks with different GOMAXPROCS values
```

#### Profile Benchmarks

Generate CPU profile:
```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

Generate memory profile:
```bash
go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof
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

### Top contributors:

<a href="https://github.com/billowdev/fastmap/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=billowdev/fastmap" alt="contrib.rocks image" />
</a>

## License

This project is licensed under the MIT License
