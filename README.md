# Xinghai Router

阶段 1、2 的 OpenAI 兼容 LLM 网关与运营后台。管理员可管理用户、密钥、渠道、路由和模型价格；用户通过一个 API Key 调用模型，并获取自己的用量、余额和账本。

## Included

- PostgreSQL migrations for users、哈希 API Key、加密渠道凭据、不可变钱包账本、用量、路由和审计记录。
- 管理员令牌保护的用户、API Key、渠道、价格、余额调整、日志和审计 API。
- OpenAI-compatible `GET /v1/models` and `POST /v1/chat/completions` endpoints.
- 透明 SSE、上游超时、每 Key 每分钟基础限流、请求 ID、模型别名和同优先级权重路由。
- 对可重试上游错误自动切换备用渠道；连续失败三次的渠道冷却一分钟。
- 请求日志关联用户、Key、模型和最终渠道；非流式请求记录 token、按定价结算，并在上游调用前预留余额以避免并发透支。

当前限流器仍是进程内实现，水平扩容前必须替换为 Redis 的滑动窗口。流式响应透明透传；由于不同上游的 SSE 用量事件不一致，流式请求目前不结算 token 费用，仅记录请求日志，不应据此执行生产计费。

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

List management data with `GET /admin/users`, `GET /admin/keys`, `GET /admin/channels`, `GET /admin/request-logs`, `GET /admin/pricing`, and `GET /admin/audit-logs`. Revoke a user key with `POST /admin/keys/{id}/revoke`; enable or disable a channel with `POST /admin/channels/{id}/status` and `{"enabled":true}` or `{"enabled":false}`.

Set a model price (currency units per million tokens), then top up or adjust a user's balance:

```sh
curl -X POST http://localhost:8080/admin/pricing \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"model":"gpt-4o-mini","input_per_million":0.15,"cached_input_per_million":0.075,"output_per_million":0.60,"multiplier":1}'

curl -X POST http://localhost:8080/admin/wallets/adjustments \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"user_id":"USER_UUID","amount":10,"note":"initial credit"}'
```

Create a public-model alias for a specific channel, or apply a request quota to a user/API key:

```sh
curl -X POST http://localhost:8080/admin/model-routes \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"public_model":"gpt-4o","upstream_model":"provider-gpt-4o","channel_id":"CHANNEL_UUID","priority":10,"weight":100}'

curl -X POST http://localhost:8080/admin/quota-limits \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H 'Content-Type: application/json' \
  -d '{"user_id":"USER_UUID","window":"day","max_requests":1000}'
```

## Gateway API

Call the gateway with the API key returned by `/admin/keys`:

```sh
curl http://localhost:8080/v1/models -H "Authorization: Bearer $XINGHAI_API_KEY"

curl -N http://localhost:8080/v1/chat/completions \
  -H "Authorization: Bearer $XINGHAI_API_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hello"}],"stream":true}'
```

The router selects an enabled channel advertising the requested model. It tries the lowest numeric priority first, distributes equal-priority traffic by weight, and retries a different eligible channel for connection errors, `408`, `425`, `429`, and `5xx` responses.

## Verify

Run `go test ./...` and `go vet ./...`.
