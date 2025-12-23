package generator

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"WriteDocument", "write_document"},
		{"GetDocument", "get_document"},
		{"DocumentExists", "document_exists"},
		{"GetNewsMLDocument", "get_news_ml_document"},
		{"CreateUpload", "create_upload"},
		{"ValidateNavigaDoc", "validate_naviga_doc"},
		{"SayHello", "say_hello"},
		{"", ""},
		{"ABC", "abc"},
		{"camelCase", "camel_case"},
		{"HTTPServer", "http_server"},
		{"XMLParser", "xml_parser"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty comment",
			input:    "",
			expected: "",
		},
		{
			name:     "single line",
			input:    "Get document retrieves a document from content repo",
			expected: "Get document retrieves a document from content repo",
		},
		{
			name:     "multi-line comment",
			input:    "WriteDocument is the method you would want to use when creating or updating documents.\n\nThe method also contains functionality to handle optimistic locking.",
			expected: "WriteDocument is the method you would want to use when creating or updating documents. The method also contains functionality to handle optimistic locking.",
		},
		{
			name:     "comment with extra whitespace",
			input:    "  \n  Method to check weather a document exists  \n  ",
			expected: "Method to check weather a document exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractComment(protogen.Comments(tt.input))
			if result != tt.expected {
				t.Errorf("extractComment() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Note: TestProtoTypeToJSONType would require creating actual protoreflect.Kind values
// which is complex in tests. We test this indirectly through integration tests.
