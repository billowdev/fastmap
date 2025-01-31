# FastMap Case Converter

Convert between common string case formats with proper handling of acronyms, numbers, and special characters.

## Case Formats

1. **camelCase**
```go
convert.ToCamelCase("user_id")      // "userId"
convert.ToCamelCase("UserID")       // "userId"
convert.ToCamelCase("API-access")   // "apiAccess"
```

2. **PascalCase**
```go
convert.ToPascalCase("user_id")     // "UserId"
convert.ToPascalCase("userID")      // "UserId"
convert.ToPascalCase("API-access")  // "ApiAccess"
```

3. **snake_case**
```go
convert.ToSnakeCase("userId")       // "user_id"
convert.ToSnakeCase("UserID")       // "user_id"
convert.ToSnakeCase("APIAccess")    // "api_access"
```

4. **kebab-case**
```go
convert.ToKebabCase("userId")       // "user-id"
convert.ToKebabCase("UserID")       // "user-id"
convert.ToKebabCase("api_access")   // "api-access"
```

5. **SCREAMING_SNAKE_CASE**
```go
convert.ToScreamingSnakeCase("userId")      // "USER_ID"
convert.ToScreamingSnakeCase("api-access")  // "API_ACCESS"
```

6. **dot.case**
```go
convert.ToDotCase("userId")         // "user.id"
convert.ToDotCase("APIAccess")      // "api.access"
```

7. **Train-Case**
```go
convert.ToTrainCase("userId")       // "User-Id"
convert.ToTrainCase("api_access")   // "Api-Access"
```

8. **Title Case**
```go
convert.ToTitleCase("user_id")      // "User Id"
convert.ToTitleCase("APIAccess")    // "Api Access"
```

## Special Cases

### Acronym Handling
```go
convert.ToSnakeCase("APIResponse")    // "api_response"
convert.ToCamelCase("API_response")   // "apiResponse"
convert.ToPascalCase("api_response")  // "ApiResponse"
```

### Number Handling
```go
convert.ToSnakeCase("user123Name")    // "user_123_name"
convert.ToCamelCase("user_123_name")  // "user123Name"
```

### Mixed Cases
```go
convert.ToSnakeCase("mixedCASE_and-hyphen") // "mixed_case_and_hyphen"
convert.ToCamelCase("mixed_CASE-andHyphen") // "mixedCaseAndHyphen"
```

## Performance

All conversion functions:
- Use strings.Builder for efficient string manipulation
- Minimize allocations
- Pre-allocate buffer space when possible
- Are thread-safe