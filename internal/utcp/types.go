package utcp

// ToolCollection represents a collection of UTCP tools
type ToolCollection struct {
	Tools []Tool `json:"tools"`
}

// Tool represents a single UTCP tool definition
// Following the UTCP specification: https://github.com/universal-tool-calling-protocol/utcp-specification
type Tool struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Inputs       InputSchema   `json:"inputs"`                  // UTCP standard field name
	Outputs      *OutputSchema `json:"outputs,omitempty"`       // UTCP standard outputs
	ToolProvider *ToolProvider `json:"tool_provider,omitempty"` // UTCP standard provider info
}

// InputSchema defines the JSON Schema for tool input parameters
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a single parameter property
type Property struct {
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
	Items       *Items              `json:"items,omitempty"`       // For array types
	Properties  map[string]Property `json:"properties,omitempty"`  // For nested objects
	Format      string              `json:"format,omitempty"`      // e.g., "date-time", "email"
	Pattern     string              `json:"pattern,omitempty"`     // Regex pattern
	MinLength   *int                `json:"minLength,omitempty"`   // Min string length
	MaxLength   *int                `json:"maxLength,omitempty"`   // Max string length
	Minimum     *float64            `json:"minimum,omitempty"`     // Min number value
	Maximum     *float64            `json:"maximum,omitempty"`     // Max number value
	Default     interface{}         `json:"default,omitempty"`     // Default value
}

// Items defines the schema for array items
type Items struct {
	Type        string              `json:"type"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Description string              `json:"description,omitempty"`
}

// OutputSchema defines the JSON Schema for tool output
type OutputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// ToolProvider defines how to execute the tool (UTCP standard)
type ToolProvider struct {
	ProviderType string            `json:"provider_type"`         // "http", "grpc", "connectrpc", "mcp", etc.
	URL          string            `json:"url"`                   // Endpoint URL
	HTTPMethod   string            `json:"http_method,omitempty"` // GET, POST, PUT, DELETE (for HTTP)
	Auth         *AuthConfig       `json:"auth,omitempty"`        // Authentication configuration
	Headers      map[string]string `json:"headers,omitempty"`     // Additional headers
}

// AuthConfig defines authentication for tool providers
type AuthConfig struct {
	AuthType     string `json:"auth_type"`               // "bearer", "api_key", "oauth2"
	Token        string `json:"token,omitempty"`         // For bearer tokens (use ${auth_token} as placeholder)
	APIKey       string `json:"api_key,omitempty"`       // For API key auth
	VarName      string `json:"var_name,omitempty"`      // Header name for API key
	ClientID     string `json:"client_id,omitempty"`     // For OAuth2
	ClientSecret string `json:"client_secret,omitempty"` // For OAuth2
	TokenURL     string `json:"token_url,omitempty"`     // For OAuth2
}
