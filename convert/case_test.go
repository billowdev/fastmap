package convert_test

import (
	"testing"

	"github.com/billowdev/fastmap/convert"
)

type testCase struct {
	input    string
	expected string
}

func TestToCamelCase(t *testing.T) {
	tests := []testCase{
		{"user_id", "userId"},
		{"UserID", "userId"},
		{"USER_ID", "userId"},
		{"user-id", "userId"},
		{"userId", "userId"},
		{"APIResponse", "apiResponse"},
		{"", ""},
		{"a", "a"},
		{"A", "a"},
	}

	for _, test := range tests {
		if got := convert.ToCamelCase(test.input); got != test.expected {
			t.Errorf("ToCamelCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []testCase{
		{"user_id", "UserId"},
		{"UserID", "UserId"},
		{"USER_ID", "UserId"},
		{"user-id", "UserId"},
		{"userId", "UserId"},
		{"APIResponse", "ApiResponse"},
		{"", ""},
		{"a", "A"},
		{"A", "A"},
	}

	for _, test := range tests {
		if got := convert.ToPascalCase(test.input); got != test.expected {
			t.Errorf("ToPascalCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []testCase{
		{"userId", "user_id"},
		{"UserID", "user_id"},
		{"USER_ID", "user_id"},
		{"user-id", "user_id"},
		{"user_id", "user_id"},
		{"APIResponse", "api_response"},
		{"", ""},
		{"a", "a"},
		{"A", "a"},
	}

	for _, test := range tests {
		if got := convert.ToSnakeCase(test.input); got != test.expected {
			t.Errorf("ToSnakeCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []testCase{
		{"userId", "user-id"},
		{"UserID", "user-id"},
		{"USER_ID", "user-id"},
		{"user-id", "user-id"},
		{"user_id", "user-id"},
		{"APIResponse", "api-response"},
		{"", ""},
		{"a", "a"},
		{"A", "a"},
	}

	for _, test := range tests {
		if got := convert.ToKebabCase(test.input); got != test.expected {
			t.Errorf("ToKebabCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToScreamingSnakeCase(t *testing.T) {
	tests := []testCase{
		{"userId", "USER_ID"},
		{"UserID", "USER_ID"},
		{"USER_ID", "USER_ID"},
		{"user-id", "USER_ID"},
		{"user_id", "USER_ID"},
		{"APIResponse", "API_RESPONSE"},
		{"", ""},
		{"a", "A"},
		{"A", "A"},
	}

	for _, test := range tests {
		if got := convert.ToScreamingSnakeCase(test.input); got != test.expected {
			t.Errorf("ToScreamingSnakeCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToDotCase(t *testing.T) {
	tests := []testCase{
		{"userId", "user.id"},
		{"UserID", "user.id"},
		{"USER_ID", "user.id"},
		{"user-id", "user.id"},
		{"user_id", "user.id"},
		{"APIResponse", "api.response"},
		{"", ""},
		{"a", "a"},
		{"A", "a"},
	}

	for _, test := range tests {
		if got := convert.ToDotCase(test.input); got != test.expected {
			t.Errorf("ToDotCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToTrainCase(t *testing.T) {
	tests := []testCase{
		{"userId", "User-Id"},
		{"UserID", "User-Id"},
		{"USER_ID", "User-Id"},
		{"user-id", "User-Id"},
		{"user_id", "User-Id"},
		{"APIResponse", "Api-Response"},
		{"", ""},
		{"a", "A"},
		{"A", "A"},
	}

	for _, test := range tests {
		if got := convert.ToTrainCase(test.input); got != test.expected {
			t.Errorf("ToTrainCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}

func TestToTitleCase(t *testing.T) {
	tests := []testCase{
		{"userId", "User Id"},
		{"UserID", "User Id"},
		{"USER_ID", "User Id"},
		{"user-id", "User Id"},
		{"user_id", "User Id"},
		{"APIResponse", "Api Response"},
		{"", ""},
		{"a", "A"},
		{"A", "A"},
	}

	for _, test := range tests {
		if got := convert.ToTitleCase(test.input); got != test.expected {
			t.Errorf("ToTitleCase(%q) = %q, expected %q", test.input, got, test.expected)
		}
	}
}
