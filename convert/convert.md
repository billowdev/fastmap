# FastMap Convert Package

## Installation
```bash
go get github.com/billowdev/fastmap
```

## Features

### Case Conversion

```go
import "github.com/billowdev/fastmap/convert"

// Available conversions
camelCase := convert.ToCamelCase("user_id")           // "userId"
pascalCase := convert.ToPascalCase("user_id")         // "UserId"
snakeCase := convert.ToSnakeCase("userId")            // "user_id"
kebabCase := convert.ToKebabCase("userId")            // "user-id"
screamingSnake := convert.ToScreamingSnakeCase("userId") // "USER_ID"
dotCase := convert.ToDotCase("userId")                // "user.id"
trainCase := convert.ToTrainCase("userId")            // "User-Id"
titleCase := convert.ToTitleCase("user_id")           // "User Id"
```

### Struct Conversion

Basic Usage:
```go
type Source struct {
    Name  string
    Email string
}

type Dest struct {
    Name    string
    Contact string
}

src := &Source{
    Name:  "John",
    Email: "john@example.com",
}
dst := &Dest{}

// Basic conversion
err := convert.ConvertStruct(src, dst, nil)

// With custom options
opts := &convert.FastmapConvertOptions{
    MatchCase: true,                                  // Enable case-sensitive matching
    IgnoreUnmatchedFields: false,                     // Fail on unmatched fields
    CustomMatchers: map[string]string{
        "Email": "Contact",                           // Map Email field to Contact
    },
}
err = convert.ConvertStruct(src, dst, opts)
```

Advanced Examples:

1. Case-Insensitive Matching:
```go
type Source struct {
    NAME string
    AGE  int
}

type Dest struct {
    Name string
    Age  int
}

opts := &convert.FastmapConvertOptions{
    MatchCase: false,
}
err := convert.ConvertStruct(src, dst, opts)
```

2. Ignore Unmatched Fields:
```go
type Source struct {
    Name    string
    Age     int
    Extra   string
}

type Dest struct {
    Name    string
    Age     int
}

opts := &convert.FastmapConvertOptions{
    IgnoreUnmatchedFields: true,
}
err := convert.ConvertStruct(src, dst, opts)
```

## Error Handling

```go
switch err {
case convert.FastmapErrNotPointer:
    // Source or destination not a pointer
case convert.FastmapErrNotStruct:
    // Source or destination not a struct
case convert.ErrFieldNotFound:
    // Required field not found in destination
default:
    // Other errors
}
```

## Best Practices

1. Always use pointers for source and destination structs
2. Enable case-sensitive matching when field names must match exactly
3. Use custom matchers for fields with different names
4. Set IgnoreUnmatchedFields for partial conversions
5. Handle all potential errors

## Thread Safety
All functions in the convert package are thread-safe and can be used concurrently.



