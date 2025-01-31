package convert

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a string to camelCase.
// Example: "user_id" -> "userId"
func ToCamelCase(s string) string {
	return toCamelInitCase(s, false)
}

// ToPascalCase converts a string to PascalCase.
// Example: "user_id" -> "UserId"
func ToPascalCase(s string) string {
	return toCamelInitCase(s, true)
}

// ToSnakeCase converts a string to snake_case.
// Example: "userId" -> "user_id"
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) * 2)

	for i, r := range s {
		if i == 0 {
			result.WriteRune(unicode.ToLower(r))
			continue
		}

		prevRune := rune(s[i-1])

		if r == '-' {
			result.WriteRune('_')
			continue
		}

		switch {
		case unicode.IsUpper(r):
			if !unicode.IsUpper(prevRune) && prevRune != '_' {
				result.WriteRune('_')
			} else if i+1 < len(s) && unicode.IsLower(rune(s[i+1])) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		case unicode.IsNumber(r):
			if !unicode.IsNumber(prevRune) && prevRune != '_' {
				result.WriteRune('_')
			}
			result.WriteRune(r)
		default:
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}

// ToKebabCase converts a string to kebab-case.
// Example: "userId" -> "user-id"
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// ToScreamingSnakeCase converts a string to SCREAMING_SNAKE_CASE.
// Example: "userId" -> "USER_ID"
func ToScreamingSnakeCase(s string) string {
	return strings.ToUpper(ToSnakeCase(s))
}

// ToDotCase converts a string to dot.case.
// Example: "userId" -> "user.id"
func ToDotCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", ".")
}

// ToConstantCase converts a string to CONSTANT_CASE (alias for SCREAMING_SNAKE_CASE).
// Example: "userId" -> "USER_ID"
func ToConstantCase(s string) string {
	return ToScreamingSnakeCase(s)
}

// ToTrainCase converts a string to Train-Case.
// Example: "userId" -> "User-Id"
func ToTrainCase(s string) string {
	words := strings.Split(ToKebabCase(s), "-")
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
	}
	return strings.Join(words, "-")
}

// ToTitleCase converts a string to Title Case.
// Example: "user_id" -> "User Id"
func ToTitleCase(s string) string {
	words := strings.Split(ToSnakeCase(s), "_")
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = string(unicode.ToUpper(rune(word[0]))) + word[1:]
	}
	return strings.Join(words, " ")
}

// toCamelInitCase is a helper function for camelCase and PascalCase conversion
func toCamelInitCase(s string, pascal bool) string {
	if s == "" {
		return s
	}

	s = ToSnakeCase(s)
	words := strings.Split(s, "_")
	var result strings.Builder
	result.Grow(len(s))

	for i, word := range words {
		if word == "" {
			continue
		}
		if i == 0 && !pascal {
			result.WriteString(strings.ToLower(word))
			continue
		}
		result.WriteString(string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:]))
	}

	return result.String()
}
