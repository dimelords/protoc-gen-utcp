# Proto → UTCP Mapping Reference

This document describes exactly how protobuf definitions are translated into UTCP tool definitions.

## Tool Name

RPC method names are converted from `PascalCase` to `snake_case`.

| Proto RPC | UTCP tool name |
|-----------|----------------|
| `SayHello` | `say_hello` |
| `GetUserByID` | `get_user_by_id` |
| `ListDocuments` | `list_documents` |

## Tool Description

The description is taken from the leading comment on the RPC method. If no comment exists, a default is generated.

```protobuf
// SayHello returns a greeting message   ← becomes description
rpc SayHello(HelloRequest) returns (HelloResponse);
```

If no comment is present, the fallback is:
```
"SayHello method in GreetingService service"
```

## Input Schema

The request message is mapped to a JSON Schema object under `inputs`.

Each proto field becomes a property. The field's JSON name (camelCase) is used as the property key.

### Scalar type mapping

| Proto type | JSON Schema type |
|------------|-----------------|
| `bool` | `boolean` |
| `int32`, `sint32`, `uint32`, `int64`, `sint64`, `uint64`, `fixed32`, `sfixed32`, `fixed64`, `sfixed64` | `integer` |
| `float`, `double` | `number` |
| `string` | `string` |
| `bytes` | `string` (base64 encoded) |
| `message` | `object` |
| `enum` | `string` |

### Repeated fields

Repeated fields become `array` type with an `items` schema:

```protobuf
repeated string tags = 4;
```

```json
"tags": {
  "type": "array",
  "items": { "type": "string" }
}
```

### Map fields

Map fields become `object` type (without typed value schema).

### Enum fields

Enum fields include an `enum` array with all value names:

```protobuf
enum Role {
  ROLE_UNSPECIFIED = 0;
  ROLE_ADMIN = 1;
  ROLE_USER = 2;
}
```

```json
"role": {
  "type": "string",
  "enum": ["ROLE_UNSPECIFIED", "ROLE_ADMIN", "ROLE_USER"]
}
```

### Nested messages

Nested message fields are inlined as `object` with their own `properties`:

```protobuf
message Metadata {
  int64 created_at = 1;
  string author = 2;
}
```

```json
"metadata": {
  "type": "object",
  "properties": {
    "createdAt": { "type": "integer" },
    "author": { "type": "string" }
  }
}
```

### Required fields

In proto3, fields are optional by default. Only fields explicitly marked `required` (proto2 syntax) are added to the `required` array. In practice, for proto3 services the `required` array will be empty.

## Output Schema

The response message is mapped identically to the input schema, but placed under `outputs`. The same type mapping rules apply.

## Field Descriptions

Leading comments on message fields become the `description` property:

```protobuf
// Name of the person to greet (required)
string name = 1;
```

```json
"name": {
  "type": "string",
  "description": "Name of the person to greet (required)"
}
```

## Tool Provider

The `tool_provider` block is generated from the plugin options. See [configuration.md](configuration.md) for details on how the URL is constructed per provider type.

## Multiple Services in One File

If a `.proto` file contains multiple services, all their methods are collected into a single `tools` array in one output file. The output file is named after the `.proto` file (e.g. `service.utcp.json`).

```protobuf
service UserService {
  rpc GetUser(...) returns (...);
}

service OrderService {
  rpc CreateOrder(...) returns (...);
}
```

Produces one file with both `get_user` and `create_order` as tools.

## Streaming RPCs

Client-streaming, server-streaming, and bidirectional-streaming RPCs are currently generated the same way as unary RPCs. The `tool_provider` URL follows the same convention. There is no special streaming metadata in the output at this time.
