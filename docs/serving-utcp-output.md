# Serving the Generated UTCP JSON

After running `protoc-gen-utcp` you get a `.utcp.json` file. This document covers common patterns for exposing it so UTCP clients (e.g. AI agents) can discover your tools.

## Output file location

The generated file is placed next to the `.proto` file, named `<proto-basename>.utcp.json`.

```
service.proto  →  service.utcp.json
```

## Option 1: Embed and serve via HTTP

Embed the file in your Go binary and expose it on a `/utcp` endpoint.

```go
import (
    _ "embed"
    "net/http"
)

//go:embed service.utcp.json
var utcpJSON []byte

func main() {
    http.HandleFunc("/utcp", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write(utcpJSON)
    })
    http.ListenAndServe(":8080", nil)
}
```

UTCP clients can then discover tools at `https://your-service.example.com/utcp`.

## Option 2: Serve as a static file

Place the `.utcp.json` file behind any static file server or CDN:

```
https://api.example.com/utcp/service.utcp.json
```

## Option 3: Commit to repository

For internal tooling or AI agent configuration, commit the generated file to your repository and reference it directly. This works well when the agent reads tool definitions from a local path.

## Multiple services / files

If you have multiple `.proto` files, each generates its own `.utcp.json`. You can:

- Serve each file on a separate endpoint (`/utcp/users`, `/utcp/orders`)
- Merge them at build time with a simple script and serve a single `/utcp` endpoint
- Let the UTCP client load multiple manifests

There is no built-in merge step in `protoc-gen-utcp` today. A minimal merge example:

```bash
jq -s '{ tools: [.[].tools[]] }' *.utcp.json > combined.utcp.json
```

## Regenerating on proto changes

Add generation to your build pipeline so the JSON stays in sync:

```makefile
generate:
    protoc --utcp_out=. --utcp_opt=base_url=https://api.example.com service.proto
```

Or with buf:

```bash
buf generate
```
