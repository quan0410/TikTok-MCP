# tiktok-mcp

Go-based Model Context Protocol (MCP) tool-server for **TikTok Shop Open API**, adhering to the DatumBridge MCP Tool Reference Architecture (BR1).

## Features

- **Streamable HTTP MCP** transport on `/mcp` (Protocol version `2024-11-05`).
- **Comprehensive Scope Coverage (24 Scopes)**: Products, Orders, Fulfillment, Logistics, Package Split/Combine, Delivery Status, Returns & Refunds, Promotions, Finance, Shop Info, and Analytics.
- **Shop Cipher Pinning**: Shop identifier and cipher are pinned from the decrypted credential vault bundle to prevent tenant cross-talk.
- **HMAC-SHA256 Request Signing**: Automatic calculation of TikTok API signatures for all calls.
- **Credential Injection**: Full support for DatumBridge credential vault injection via `credentials_json` / `credentials_path`, with local dev fallback.

## Running Locally

```bash
# Copy environment template
cp .env.example .env

# Set local test credentials
export TIKTOK_SHOP_APP_KEY="your_app_key"
export TIKTOK_SHOP_APP_SECRET="your_app_secret"
export TIKTOK_SHOP_ACCESS_TOKEN="your_token"
export TIKTOK_SHOP_CIPHER="your_shop_cipher"

# Run server
go run ./cmd/api
```

The server listens on `http://localhost:8015`.

## Health Check

```bash
curl -i http://localhost:8015/health
```

Returns:
```json
{"status":"ok","service":"tiktok-mcp"}
```

## Running with Docker

```bash
docker build -t tiktok-mcp .
docker run -p 8015:8015 tiktok-mcp
```

## Running Tests

```bash
go test ./...
```
