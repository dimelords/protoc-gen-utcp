# Annotations and Metadata

## What is supported today

`protoc-gen-utcp` currently derives all metadata from standard proto comments and the plugin options passed at generation time. There is no custom proto annotation support.

### Descriptions

Tool and field descriptions come from leading comments in the `.proto` file:

```protobuf
// GetUser retrieves a user by their ID.
// Returns a 404 error if the user does not exist.
rpc GetUser(GetUserRequest) returns (User);
```

```protobuf
message GetUserRequest {
  // The unique identifier of the user.
  string user_id = 1;
}
```

Both the RPC comment and the field comments are picked up automatically.

### Tags, authentication, and other metadata

These are set globally via plugin options (`--utcp_opt`) and apply to every tool in the generated file. See [configuration.md](configuration.md) for the full list.

There is currently no way to set per-method authentication or per-method tags via annotations.

## What is not supported (yet)

The following require custom proto annotations that are not yet implemented:

- Per-method `base_url` override
- Per-method `auth_type` override
- Custom tags on individual tools
- Marking a method as requiring user confirmation
- Hiding a method from the generated output

## Workaround: use comments for human-readable hints

Until annotation support is added, you can embed structured hints in comments. UTCP clients that parse descriptions may pick these up:

```protobuf
// DeleteDocument removes a document permanently.
// WARNING: This operation is destructive and cannot be undone.
// Requires user confirmation before executing.
rpc DeleteDocument(DeleteDocumentRequest) returns (DeleteDocumentResponse);
```

## Contributing annotation support

If you need per-method metadata, the recommended approach is to define a custom proto extension and add handling in `internal/generator/generator.go`. See [CONTRIBUTING.md](../CONTRIBUTING.md) for how to get started.
