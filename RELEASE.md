# FastMap v1.1.0 Release Notes

We are pleased to announce the release of FastMap v1.1.0, introducing two major enhancements that improve the library's functionality and error handling capabilities.

## What's Changed
* Add AppendableHashMap and ForEach Error Handling (v1.1.0) by @billowdev in https://github.com/billowdev/fastmap/pull/2


**Full Changelog**: https://github.com/billowdev/fastmap/compare/v1.0.0...v1.1.0

## New Features

### AppendableHashMap

This release introduces AppendableHashMap, a specialized HashMap implementation designed for managing slices of values. This new feature provides a more intuitive and type-safe way to handle collections of values associated with keys.

The AppendableHashMap offers:
- Type-safe operations for managing slices of values
- Convenient methods for appending multiple values to existing slices
- Thread-safe variant (ThreadSafeAppendableHashMap) for concurrent operations
- Memory-efficient handling of growing collections

Example usage:
```go
// Create a new AppendableHashMap
components := NewAppendableHashMap[string, Component]()

// Add values
components.AppendValues("header", headerComponent1, headerComponent2)
components.AppendValues("header", headerComponent3) // Appends to existing slice

// Thread-safe version
safeComponents := NewThreadSafeAppendableHashMap[string, Component]()
safeComponents.AppendValues("footer", footerComponent1, footerComponent2)
```

### Enhanced ForEach Error Handling

The ForEach method has been enhanced to support error propagation, allowing for better error handling during iterations. This improvement enables more robust error management in applications using FastMap.

Example usage:
```go
err := hashMap.ForEach(func(key string, value User) error {
    if value.IsInvalid() {
        return fmt.Errorf("invalid user data for key %s", key)
    }
    fmt.Printf("User %s: %v\n", key, value)
    return nil
})
```

## Migration Guide

### Updating ForEach Implementation

The ForEach method now includes error handling. Here's how to update your existing code:

Previous implementation:
```go
hashMap.ForEach(func(k string, v int) {
    fmt.Printf("%s: %d\n", k, v)
})
```

New implementation:
```go
err := hashMap.ForEach(func(k string, v int) error {
    fmt.Printf("%s: %d\n", k, v)
    return nil
})
if err != nil {
    // Handle error
}
```

## Bug Fixes and Improvements

- Enhanced thread safety in concurrent operations
- Improved memory management for slice operations
- Updated documentation with comprehensive examples
- Added new test coverage for AppendableHashMap features

## Breaking Changes

The ForEach method signature has been updated to include error handling. Existing implementations will need to be modified to handle error returns.

## Documentation

Complete documentation for the new features is available in the package documentation. Each new type and method includes detailed examples and usage guidelines.

## Acknowledgments

We thank our contributors and users for their valuable feedback and suggestions that helped shape these improvements.

For more detailed information about these changes, please refer to our GitHub repository and documentation.



# FastMap v1.0.0 Release Notes

## What's Changed
* Pull Request: Release fastmap v1.0.0 by @billowdev in https://github.com/billowdev/fastmap/pull/1

## New Contributors
* @billowdev made their first contribution in https://github.com/billowdev/fastmap/pull/1

**Full Changelog**: https://github.com/billowdev/fastmap/commits/1.0.0

# Release fastmap v1.0.0

## Overview
This pull request introduces version 1.0.0 of the `fastmap` package, representing the first stable release of our generic, type-safe HashMap implementation in Go. This release includes significant improvements in functionality, thread-safety, and configuration processing capabilities.

## Key Features and Improvements

### Core Functionality
- Implemented generic, type-safe HashMap with both thread-safe and non-thread-safe variants
- Comprehensive support for key-value operations in Go
- Zero external dependencies
- Full type inference and compile-time type checking

### Thread-Safe Operations
- Introduced `NewThreadSafeHashMap` for concurrent access
- Implemented safe concurrent operations with `sync.RWMutex`
- Added thread-safe functional operations like `Filter`, `Map`, and `ForEach`

### Advanced Configuration Processing
- Added robust field configuration processing capabilities
- Support for complex data transformations
- Flexible handler functions for dynamic data processing
- Enhanced error handling and validation

### Performance and Reliability
- Minimal method call overhead
- O(1) average case for basic operations
- Comprehensive test coverage
- Extensive benchmarking suite

## Breaking Changes
- None. This is a feature-complete initial release maintaining full backward compatibility

## Testing
- Added comprehensive unit tests covering all core functionality
- Implemented benchmarks for performance evaluation
- Achieved high test coverage across different use cases and scenarios

## Documentation
- Updated README with detailed usage examples
- Provided migration guide for transitioning from traditional maps
- Included best practices and performance considerations

## Examples of New Capabilities

### Thread-Safe User Management
```go
system := NewUserSystem()
go func() {
    system.AddUser("1", User{Name: "John", Active: true})
}()
go func() {
    system.UpdateUserStatus("1", false)
}()
```

### Configuration Field Processing
```go
safeMap := fastmap.NewThreadSafeHashMap[string, Measurement]()
configs := map[string]fastmap.FieldConfig[Measurement]{
    "sensor_data": {
        Handler: func(data map[string]interface{}) *Measurement {
            // Complex data transformation logic
        },
    },
}
```

## Next Steps
- Continued performance optimization
- Expanded test coverage
- Community feedback integration
- Potential additional functional programming utilities

## Contribution
Special thanks to all contributors who helped shape this initial release.

## License
Continues to be distributed under the MIT License.