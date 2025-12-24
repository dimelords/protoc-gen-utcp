# protoc-gen-utcp

[![CI](https://github.com/dimelords/protoc-gen-utcp/actions/workflows/ci.yml/badge.svg)](https://github.com/dimelords/protoc-gen-utcp/actions/workflows/ci.yml)
[![CodeQL](https://github.com/dimelords/protoc-gen-utcp/actions/workflows/codeql.yml/badge.svg)](https://github.com/dimelords/protoc-gen-utcp/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dimelords/protoc-gen-utcp)](https://goreportcard.com/report/github.com/dimelords/protoc-gen-utcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Protocol Buffer compiler plugin that generates [UTCP (Universal Tool Calling Protocol)](https://github.com/universal-tool-calling-protocol/utcp-specification) definitions from `.proto` files.

## Overview

`protoc-gen-utcp` enables AI agents to discover and execute RPC methods by converting Protocol Buffer service definitions into UTCP-compliant tool specifications. This eliminates the need for runtime proto parsing or OpenAPI intermediate steps.

### Key Features

- ✅ **Build-time generation** - Generate UTCP at compile time, not runtime
- ✅ **Full proto3 support** - Messages, enums, nested types, repeated fields
- ✅ **Multiple providers** - HTTP, gRPC, ConnectRPC support
- ✅ **Twirp compatible** - Perfect for Twirp-based APIs
- ✅ **Rich schemas** - Generates complete JSON Schema from proto messages
- ✅ **Customizable** - Configure base URLs, auth types, provider types

## Installation

```bash
go install github.com/dimelords/protoc-gen-utcp/cmd/protoc-gen-utcp@latest
```

Or build from source:

```bash
git clone https://github.com/dimelords/protoc-gen-utcp.git
cd protoc-gen-utcp
make install
```

## Quick Start

### 1. Define your service

```protobuf
syntax = "proto3";

package example.v1;

service GreetingService {
  // SayHello returns a greeting message
  rpc SayHello(HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

### 2. Generate UTCP

```bash
protoc --utcp_out=. \
  --utcp_opt=base_url=https://api.example.com \
  --utcp_opt=provider_type=http \
  --utcp_opt=auth_type=bearer \
  service.proto
```

### 3. Use the generated JSON

```json
{
  "tools": [
    {
      "name": "say_hello",
      "description": "SayHello returns a greeting message",
      "inputs": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string"
          }
        }
      },
      "outputs": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string"
          }
        }
      },
      "tool_provider": {
        "provider_type": "http",
        "url": "https://api.example.com/twirp/example.v1.GreetingService/SayHello",
        "http_method": "POST",
        "auth": {
          "auth_type": "bearer",
          "token": "${auth_token}"
        },
        "headers": {
          "Content-Type": "application/json"
        }
      }
    }
  ]
}
```

## Configuration Options

### Plugin Options

Configure via `--utcp_opt`:

| Option | Description | Default |
|--------|-------------|---------|
| `base_url` | Base URL for API endpoints | `https://api.example.com` |
| `provider_type` | Provider type (`http`, `grpc`, `connectrpc`) | `http` |
| `auth_type` | Authentication type (`bearer`, `api_key`, `oauth2`) | `bearer` |

### Examples

**Twirp API:**
```bash
protoc --utcp_out=. \
  --utcp_opt=base_url=https://api.example.com \
  --utcp_opt=provider_type=http \
  --utcp_opt=auth_type=bearer \
  service.proto
```

**gRPC API:**
```bash
protoc --utcp_out=. \
  --utcp_opt=base_url=grpc://api.example.com \
  --utcp_opt=provider_type=grpc \
  --utcp_opt=auth_type=bearer \
  service.proto
```

**ConnectRPC API:**
```bash
protoc --utcp_out=. \
  --utcp_opt=base_url=https://api.example.com \
  --utcp_opt=provider_type=connectrpc \
  --utcp_opt=auth_type=bearer \
  service.proto
```


## URL Conventions

### HTTP/Twirp (Default)

```
POST {base_url}/twirp/{package}.{Service}/{Method}
```

Example:
```
POST https://api.example.com/twirp/example.v1.GreetingService/SayHello
```

### gRPC

```
{base_url}/{package}.{Service}/{Method}
```

Example:
```
grpc://api.example.com/example.v1.GreetingService/SayHello
```

## Features

### Rich Input Schemas

The plugin generates complete JSON Schema from proto messages:

```protobuf
message CreateUserRequest {
  string email = 1;  // User's email address
  int32 age = 2;     // User's age
  Role role = 3;     // User's role
  repeated string tags = 4;  // User tags
}

enum Role {
  ROLE_UNSPECIFIED = 0;
  ROLE_ADMIN = 1;
  ROLE_USER = 2;
}
```

Generates:

```json
{
  "inputs": {
    "type": "object",
    "properties": {
      "email": {
        "type": "string",
        "description": "User's email address"
      },
      "age": {
        "type": "integer",
        "description": "User's age"
      },
      "role": {
        "type": "string",
        "description": "User's role",
        "enum": ["ROLE_UNSPECIFIED", "ROLE_ADMIN", "ROLE_USER"]
      },
      "tags": {
        "type": "array",
        "description": "User tags",
        "items": {
          "type": "string"
        }
      }
    }
  }
}
```

### Nested Types

Handles nested messages, enums, and complex types:

```protobuf
message Document {
  string uuid = 1;
  Metadata metadata = 2;

  message Metadata {
    int64 created_at = 1;
    string author = 2;
  }
}
```

Generates nested property schemas.

## Integration Examples

### Using with Embedded Files

**1. Generate UTCP at build time:**

```bash
# In your Makefile or build script
protoc --utcp_out=. \
  --utcp_opt=base_url=https://api.example.com \
  --utcp_opt=provider_type=http \
  --utcp_opt=auth_type=bearer \
  service.proto
```

**2. Embed the generated JSON in your Go application:**

```go
import (
    _ "embed"
    "encoding/json"
)

//go:embed service.utcp.json
var utcpToolsJSON []byte

func getTools() (*ToolCollection, error) {
    var tools ToolCollection
    err := json.Unmarshal(utcpToolsJSON, &tools)
    return &tools, err
}
```

**3. Use in your application:**

```go
func main() {
    tools, err := getTools()
    if err != nil {
        log.Fatal(err)
    }

    // Serve via HTTP, use with UTCP client, etc.
    serveTools(tools)
}
```

### Benefits

- ✅ **No runtime parsing** - JSON pre-generated at build time
- ✅ **Fast startup** - No proto parser dependency
- ✅ **Small binaries** - Zero runtime dependencies
- ✅ **Type safe** - Generated at compile time with validation
- ✅ **Consistent** - Same output every build

## Examples

See the `examples/` directory for complete examples:

### Simple Example

```bash
make examples
cat examples/simple/service.utcp.json
```

### Twirp Example

```bash
make examples
cat examples/twirp/documents.utcp.json
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Generating Examples

```bash
make examples
```

### Installing Locally

```bash
make install
```

## UTCP Specification

This plugin generates UTCP-compliant JSON following the official specification:
https://github.com/universal-tool-calling-protocol/utcp-specification

### Generated Structure

```json
{
  "tools": [
    {
      "name": "tool_name",
      "description": "Tool description",
      "inputs": { /* JSON Schema */ },
      "outputs": { /* JSON Schema */ },
      "tool_provider": {
        "provider_type": "http",
        "url": "https://...",
        "http_method": "POST",
        "auth": {
          "auth_type": "bearer",
          "token": "${auth_token}"
        },
        "headers": {
          "Content-Type": "application/json"
        }
      }
    }
  ]
}
```

## Comparison with Alternatives

### vs OpenAPI Generation

| Aspect | protoc-gen-utcp | Proto→OpenAPI→UTCP |
|--------|-----------------|---------------------|
| Steps | 1 (Direct) | 2 (Indirect) |
| Dependencies | None | protoc-gen-openapi |
| Output | UTCP JSON | OpenAPI YAML → UTCP |
| Twirp Support | Native | Requires annotations |
| Build Complexity | Simple | Complex |

### vs Runtime Parsing

| Aspect | Build-time (this) | Runtime |
|--------|-------------------|---------|
| Cold Start | Fast | Slow (parsing) |
| Binary Size | Small | Large (parser) |
| Dependencies | None | proto parser libs |
| Consistency | Guaranteed | Variable |

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `make test` passes
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Credits

- UTCP Specification: https://github.com/universal-tool-calling-protocol/utcp-specification
- Protocol Buffers: https://protobuf.dev/
- Inspired by protoc plugins like `protoc-gen-go` and `protoc-gen-grpc-gateway`

## Related Projects

- [UTCP Specification](https://github.com/universal-tool-calling-protocol/utcp-specification)
- [Twirp](https://twitchtv.github.io/twirp/) - RPC framework
- [ConnectRPC](https://connectrpc.com/) - Modern gRPC alternative
