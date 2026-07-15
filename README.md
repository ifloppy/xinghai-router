# Xinghai Router

Phase 1 implementation of an OpenAI-compatible LLM gateway. It gives operators a small admin API for users, API keys, and upstream channels, then proxies `/v1/models` and `/v1/chat/completions` through an authenticated API key.

## Included

- PostgreSQL migrations for users, hashed API keys, encrypted channel credentials, and request logs.
- Admin-token-protected APIs to create and list users, API keys, channels, and request logs.
- OpenAI-compatible `GET /v1/models` and `POST /v1/chat/completions` endpoints.
- Transparent SSE streaming, upstream timeouts, a per-key in-memory per-minute rate limit, request IDs, and model-based priority routing.
- Request logs with model, selected channel, HTTP status, elapsed time, and non-streaming token usage.

This is intentionally a Phase 1 gateway, not a billing-ready production service. The per-process rate limit must move to Redis before running multiple instances. Streaming responses record zero token usage because the service transparently forwards SSE and does not require a provider-specific final usage event; do not use those logs for billing yet.

## Run locally

1. Create local infrastructure: `docker compose up -d`.
2. Create configuration: `cp .env.example .env`, then replace both secrets with unique random values.
3. Export the environment variables in `.env` using your shell or an environment loader.
4. Run: `go run ./cmd/router`.
5. Check: `curl http://localhost:8080/healthz`.

### Admin web console

The Vue 3 management console is in `web/`. Start the Go service first, then run:

```sh
cd web
npm install
npm run dev
```

Open `http://localhost:5173` and enter `ADMIN_TOKEN`. The token is retained only in browser session storage. Vite proxies `/admin` calls to `http://localhost:8080`, so this development setup does not require a CORS policy. Create a production deployment by running `npm run build`; serve the generated `web/dist` directory behind the router or a reverse proxy.

The service performs migrations automatically at startup. `base_url` for a channel must be an HTTPS origin or path prefix without `/v1`; for example, `https://api.openai.com`. Provider secrets are encrypted in the database using `ENCRYPTION_KEY`, so keep this value stable and securely backed up.

## Admin API

All `/admin` endpoints require `Authorization: Bearer $ADMIN_TOKEN`.

Create a user:

```sh
curl -X POST http://localhost:8080/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","name":"Example User"}'
```

Create an API key. The full `key` in the response is displayed only at creation time:

```sh
curl -X POST http://localhost:8080/admin/keys \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"USER_UUID","name":"development"}'
```

Create an OpenAI-compatible upstream channel:

```sh
curl -X POST http://localhost:8080/admin/channels \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"openai","base_url":"https://api.openai.com","api_key":"PROVIDER_KEY","models":["gpt-4o-mini"],"priority":100}'
```

List management data with `GET /admin/users`, `GET /admin/keys`, `GET /admin/channels`, and `GET /admin/request-logs`. Revoke a user key with `POST /admin/keys/{id}/revoke`; enable or disable a channel with `POST /admin/channels/{id}/status` and `{"enabled":true}` or `{"enabled":false}`.

## Gateway API

Call the gateway with the API key returned by `/admin/keys`:

```sh
curl http://localhost:8080/v1/models -H "Authorization: Bearer $XINGHAI_API_KEY"

curl -N http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer $XINGHAI_API_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hello"}],"stream":true}'
```

The router selects the enabled channel with the lowest numeric priority that advertises the requested model. It does not yet retry or fall back to another channel.

## Verify

Run `go test ./...` and `go vet ./...`.
