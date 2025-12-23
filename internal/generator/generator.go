package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/dimelords/protoc-gen-utcp/internal/utcp"
)

// GenerateFile generates a UTCP JSON file from a proto file
func GenerateFile(gen *protogen.Plugin, file *protogen.File, config *Config) error {
	if len(file.Services) == 0 {
		return nil // No services to generate
	}

	filename := file.GeneratedFilenamePrefix + ".utcp.json"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	tools := &utcp.ToolCollection{
		Tools: make([]utcp.Tool, 0),
	}

	// Process each service in the file
	for _, service := range file.Services {
		serviceTools := generateServiceTools(service, file, config)
		tools.Tools = append(tools.Tools, serviceTools...)
	}

	// Marshal to pretty JSON
	data, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal UTCP tools: %w", err)
	}

	g.P(string(data))

	return nil
}

// generateServiceTools generates UTCP tools for a single service
func generateServiceTools(service *protogen.Service, file *protogen.File, config *Config) []utcp.Tool {
	tools := make([]utcp.Tool, 0, len(service.Methods))

	for _, method := range service.Methods {
		tool := generateTool(service, method, file, config)
		tools = append(tools, tool)
	}

	return tools
}

// generateTool generates a single UTCP tool from an RPC method
func generateTool(service *protogen.Service, method *protogen.Method, file *protogen.File, config *Config) utcp.Tool {
	toolName := toSnakeCase(string(method.Desc.Name()))
	description := extractComment(method.Comments.Leading)
	if description == "" {
		description = fmt.Sprintf("%s method in %s service", method.Desc.Name(), service.Desc.Name())
	}

	tool := utcp.Tool{
		Name:        toolName,
		Description: description,
		Inputs: utcp.InputSchema{
			Type:       "object",
			Properties: make(map[string]utcp.Property),
			Required:   []string{},
		},
	}

	// Generate input schema from request message
	if method.Input != nil {
		tool.Inputs = generateInputSchema(method.Input)
	}

	// Generate output schema from response message
	if method.Output != nil {
		tool.Outputs = generateOutputSchema(method.Output)
	}

	// Generate tool provider based on config
	tool.ToolProvider = generateToolProvider(service, method, file, config)

	return tool
}

// generateInputSchema generates the input schema from a message
func generateInputSchema(message *protogen.Message) utcp.InputSchema {
	schema := utcp.InputSchema{
		Type:       "object",
		Properties: make(map[string]utcp.Property),
		Required:   []string{},
	}

	for _, field := range message.Fields {
		prop := messageFieldToProperty(field)
		schema.Properties[field.Desc.JSONName()] = prop

		// Mark required fields (proto3: all fields are optional by default)
		// In proto3, we can't distinguish required vs optional without annotations
		// For now, we'll leave required empty unless the field is explicitly required
		if field.Desc.Cardinality() == protoreflect.Required {
			schema.Required = append(schema.Required, field.Desc.JSONName())
		}
	}

	return schema
}

// generateOutputSchema generates the output schema from a message
func generateOutputSchema(message *protogen.Message) *utcp.OutputSchema {
	schema := &utcp.OutputSchema{
		Type:       "object",
		Properties: make(map[string]utcp.Property),
	}

	for _, field := range message.Fields {
		prop := messageFieldToProperty(field)
		schema.Properties[field.Desc.JSONName()] = prop
	}

	return schema
}

// messageFieldToProperty converts a proto field to a UTCP property
func messageFieldToProperty(field *protogen.Field) utcp.Property {
	prop := utcp.Property{
		Description: extractComment(field.Comments.Leading),
	}

	// Handle repeated fields (arrays)
	if field.Desc.IsList() {
		prop.Type = "array"
		prop.Items = &utcp.Items{
			Type:        protoTypeToJSONType(field.Desc.Kind()),
			Description: prop.Description,
		}
		return prop
	}

	// Handle map fields
	if field.Desc.IsMap() {
		prop.Type = "object"
		return prop
	}

	// Handle scalar and message types
	prop.Type = protoTypeToJSONType(field.Desc.Kind())

	// Handle enum fields
	if field.Desc.Kind() == protoreflect.EnumKind && field.Enum != nil {
		enumValues := make([]string, 0, len(field.Enum.Values))
		for _, val := range field.Enum.Values {
			enumValues = append(enumValues, string(val.Desc.Name()))
		}
		prop.Enum = enumValues
	}

	// Handle nested message fields
	if field.Desc.Kind() == protoreflect.MessageKind && field.Message != nil {
		prop.Type = "object"
		prop.Properties = make(map[string]utcp.Property)
		for _, nestedField := range field.Message.Fields {
			prop.Properties[nestedField.Desc.JSONName()] = messageFieldToProperty(nestedField)
		}
	}

	return prop
}

// generateToolProvider generates the tool provider configuration
func generateToolProvider(service *protogen.Service, method *protogen.Method, file *protogen.File, config *Config) *utcp.ToolProvider {
	provider := &utcp.ToolProvider{
		ProviderType: config.ProviderType,
	}

	// Generate URL based on provider type
	switch config.ProviderType {
	case "http":
		provider.URL = generateHTTPURL(service, method, file, config)
		provider.HTTPMethod = "POST" // Default for most RPC frameworks
		provider.Headers = map[string]string{
			"Content-Type": "application/json",
		}

	case "grpc", "connectrpc":
		provider.URL = generateGRPCURL(service, method, file, config)

	default:
		// For custom provider types, use a generic URL
		provider.URL = fmt.Sprintf("%s/%s.%s/%s", config.BaseURL, file.Desc.Package(), service.Desc.Name(), method.Desc.Name())
	}

	// Add authentication
	if config.AuthType != "" {
		provider.Auth = &utcp.AuthConfig{
			AuthType: config.AuthType,
		}

		if config.AuthType == "bearer" {
			provider.Auth.Token = "${auth_token}"
		}
	}

	return provider
}

// generateHTTPURL generates an HTTP URL (Twirp-style by default)
func generateHTTPURL(service *protogen.Service, method *protogen.Method, file *protogen.File, config *Config) string {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.example.com"
	}

	// Twirp convention: /twirp/{package}.{Service}/{Method}
	return fmt.Sprintf("%s/twirp/%s.%s/%s",
		strings.TrimSuffix(config.BaseURL, "/"),
		file.Desc.Package(),
		service.Desc.Name(),
		method.Desc.Name(),
	)
}

// generateGRPCURL generates a gRPC-style URL
func generateGRPCURL(service *protogen.Service, method *protogen.Method, file *protogen.File, config *Config) string {
	if config.BaseURL == "" {
		config.BaseURL = "grpc://api.example.com"
	}

	// gRPC convention: {package}.{Service}/{Method}
	return fmt.Sprintf("%s/%s.%s/%s",
		strings.TrimSuffix(config.BaseURL, "/"),
		file.Desc.Package(),
		service.Desc.Name(),
		method.Desc.Name(),
	)
}

// protoTypeToJSONType converts proto field types to JSON Schema types
func protoTypeToJSONType(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.BoolKind:
		return "boolean"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		return "integer"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "number"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "string" // base64 encoded
	case protoreflect.MessageKind:
		return "object"
	case protoreflect.EnumKind:
		return "string"
	default:
		return "string"
	}
}

// extractComment extracts and cleans documentation comments
func extractComment(comments protogen.Comments) string {
	if comments == "" {
		return ""
	}

	text := string(comments)
	text = strings.TrimSpace(text)

	// Remove leading // or /* */ style comments
	lines := strings.Split(text, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "*/")
		line = strings.TrimSpace(line)

		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, " ")
}

// toSnakeCase converts PascalCase or camelCase to snake_case
func toSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) + 5)

	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prevIsLower := s[i-1] >= 'a' && s[i-1] <= 'z'
			nextIsLower := i < len(s)-1 && s[i+1] >= 'a' && s[i+1] <= 'z'

			if prevIsLower || nextIsLower {
				result.WriteRune('_')
			}
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}
