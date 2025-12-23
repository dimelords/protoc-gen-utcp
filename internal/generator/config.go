package generator

// Config holds the configuration for UTCP generation
type Config struct {
	// BaseURL is the base URL for API endpoints
	// Example: "https://api.example.com"
	BaseURL string

	// ProviderType is the type of provider (http, grpc, connectrpc, mcp)
	ProviderType string

	// AuthType is the authentication type (bearer, api_key, oauth2)
	AuthType string
}
