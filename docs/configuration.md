# Configuration Reference

All options are passed to `protoc` via `--utcp_opt=key=value`.

## Options

### `base_url`

The base URL prepended to all generated endpoint URLs.

- Default: `https://api.example.com`
- Required: no (but strongly recommended)

```bash
--utcp_opt=base_url=https://api.mycompany.com
```

### `provider_type`

Controls how the `tool_provider` block is generated and how the URL is formatted.

- Default: `http`
- Allowed values: `http`, `grpc`, `connectrpc`

```bash
--utcp_opt=provider_type=grpc
```

### `auth_type`

Sets the authentication type in the generated `auth` block.

- Default: `bearer`
- Allowed values: `bearer`, `api_key`, `oauth2`
- Set to empty string (`auth_type=`) to omit the `auth` block entirely.

```bash
--utcp_opt=auth_type=api_key
```

## URL conventions per provider type

### `http` (Twirp-style)

```
POST {base_url}/twirp/{package}.{Service}/{Method}
```

Example:
```
POST https://api.example.com/twirp/example.v1.GreetingService/SayHello
```

Headers added automatically:
```json
{ "Content-Type": "application/json" }
```

### `grpc`

```
{base_url}/{package}.{Service}/{Method}
```

Example:
```
grpc://api.example.com/example.v1.GreetingService/SayHello
```

### `connectrpc`

Same URL format as `grpc`:
```
{base_url}/{package}.{Service}/{Method}
```

### Custom / unknown provider type

Falls back to:
```
{base_url}/{package}.{Service}/{Method}
```

## Authentication blocks

### `bearer`

```json
"auth": {
  "auth_type": "bearer",
  "token": "${auth_token}"
}
```

The placeholder `${auth_token}` is intended to be substituted at runtime by the UTCP client.

### `api_key`

```json
"auth": {
  "auth_type": "api_key"
}
```

### `oauth2`

```json
"auth": {
  "auth_type": "oauth2"
}
```

## Full example

```bash
protoc \
  --utcp_out=. \
  --utcp_opt=base_url=https://api.example.com \
  --utcp_opt=provider_type=http \
  --utcp_opt=auth_type=bearer \
  path/to/service.proto
```

## Using with buf

Add to your `buf.gen.yaml`:

```yaml
version: v2
plugins:
  - local: protoc-gen-utcp
    out: .
    opt:
      - base_url=https://api.example.com
      - provider_type=http
      - auth_type=bearer
```

Then run:

```bash
buf generate
```
