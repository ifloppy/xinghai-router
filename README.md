# Xinghai Router

支持 OpenAI 与 Anthropic 格式的 LLM 网关与运营后台。管理员可管理用户、密钥、渠道、路由和模型价格；用户通过一个 API Key 调用模型，并获取自己的用量、余额和账本。

## Included

- PostgreSQL migrations for users、哈希 API Key、加密渠道凭据、不可变钱包账本、用量、路由和审计记录。
- 基于用户会话、管理员角色和细粒度权限保护的管理 API。
- OpenAI-compatible `GET /v1/models`、`POST /v1/chat/completions`，以及 Anthropic-compatible `POST /v1/messages`。
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

Open `http://localhost:5173/auth` and create an account or sign in with email and password. The first registered account becomes an administrator; administrators can promote users or grant individual permissions. Browser sessions are retained only in session storage. Nuxt proxies browser requests from `/api/*` to `http://127.0.0.1:8080/*`, so this development setup does not require a CORS policy. `npm run generate` emits prerendered HTML for the public home and authentication pages; deploy the Nuxt `.output` directory for the full application.

The service performs migrations automatically at startup. `base_url` for a channel must be an HTTPS origin or path prefix without `/v1`; for example, `https://api.openai.com`. Loopback HTTP URLs are also accepted for local services such as Ollama, for example `http://127.0.0.1:11434`. Provider secrets are encrypted in the database using `ENCRYPTION_KEY`, so keep this value stable and securely backed up.

## Administration API

All `/admin` endpoints require an authenticated account session: `Authorization: Bearer $SESSION_TOKEN`. An `admin` user has every permission. Other users must be individually granted the permission required by each endpoint.

| Permission | Access |
| --- | --- |
| `users.read` | List users |
| `keys.manage` | Create, list, and revoke API keys |
| `channels.read`, `channels.manage` | View or manage upstream channels |
| `logs.read`, `audit.read` | View request or audit logs |
| `pricing.read`, `pricing.manage` | View or edit pricing |
| `wallets.manage`, `routes.manage`, `quotas.manage` | Manage balances, model routes, or quotas |
| `system.manage` | Promote users and assign permissions |

Promote a user or grant permissions using `POST /admin/users/{id}/role` with `{"role":"admin"}`, and `PUT /admin/users/{id}/permissions` with `{"permissions":["channels.read","logs.read"]}`. These operations require `system.manage`.

## Account API

Register an account. Passwords must have at least eight characters; the service stores only bcrypt password hashes. A successful registration or login returns a seven-day bearer session token.

```sh
curl -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","name":"Example User","password":"a-strong-password"}'
```

Log in with `POST /auth/login` using `{"email":"user@example.com","password":"a-strong-password"}`. Use the returned token with `Authorization: Bearer $SESSION_TOKEN` for `GET /account/me`, and revoke the current session using `POST /auth/logout`.

Create an API key. The full `key` in the response is displayed only at creation time:

```sh
curl -X POST http://localhost:8080/admin/keys \
  -H "Authorization: Bearer $SESSION_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"USER_UUID","name":"development"}'
```

Create an OpenAI-compatible upstream channel:

```sh
curl -X POST http://localhost:8080/admin/channels \
  -H "Authorization: Bearer $SESSION_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"openai","base_url":"https://api.openai.com","api_key":"PROVIDER_KEY","models":["gpt-4o-mini"],"priority":100}'
```

创建渠道时可选 `provider`：`openai`、`ollama`、`kimi`、`opencode_go` 或 `anthropic`。Ollama、Kimi 和 OpenCode Go 使用各自的 OpenAI-compatible 接口；Anthropic 渠道会转换为 Messages API。`base_url` 不要包含末尾的 `/v1`：

```sh
# 本机 Ollama，API Key 会被 Ollama 忽略
curl -X POST http://localhost:8080/admin/channels \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"ollama","provider":"ollama","base_url":"http://127.0.0.1:11434","api_key":"ollama","models":["qwen3-coder:30b"],"priority":100}'

# Kimi / Moonshot
curl -X POST http://localhost:8080/admin/channels \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"kimi","provider":"kimi","base_url":"https://api.moonshot.cn","api_key":"MOONSHOT_API_KEY","models":["kimi-k2.6"],"priority":100}'
```

OpenCode Go 使用相同方式创建渠道并设置 `"provider":"opencode_go"`，填写其 OpenAI-compatible API origin、订阅 API Key 和可用模型 ID。Anthropic 上游使用 `"provider":"anthropic"`、`https://api.anthropic.com` 和 Anthropic API Key。

List management data with `GET /admin/users`, `GET /admin/keys`, `GET /admin/channels`, `GET /admin/request-logs`, `GET /admin/pricing`, and `GET /admin/audit-logs`. Revoke a user key with `POST /admin/keys/{id}/revoke`; enable or disable a channel with `POST /admin/channels/{id}/status` and `{"enabled":true}` or `{"enabled":false}`.

Set a model price (currency units per million tokens), then top up or adjust a user's balance:

```sh
curl -X POST http://localhost:8080/admin/pricing \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"model":"gpt-4o-mini","input_per_million":0.15,"cached_input_per_million":0.075,"output_per_million":0.60,"multiplier":1}'

curl -X POST http://localhost:8080/admin/wallets/adjustments \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"user_id":"USER_UUID","amount":10,"note":"initial credit"}'
```

Create a public-model alias for a specific channel, or apply a request quota to a user/API key:

```sh
curl -X POST http://localhost:8080/admin/model-routes \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
  -d '{"public_model":"gpt-4o","upstream_model":"provider-gpt-4o","channel_id":"CHANNEL_UUID","priority":10,"weight":100}'

curl -X POST http://localhost:8080/admin/quota-limits \
  -H "Authorization: Bearer $SESSION_TOKEN" -H 'Content-Type: application/json' \
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

Anthropic 客户端（包括将 OpenCode 的 Anthropic provider 指向本服务）可使用 `x-api-key` 调用 `/v1/messages`。请求、非流式响应、SSE 和工具调用会在 Anthropic Messages 与上游 OpenAI Chat Completions 格式之间转换：

```sh
curl -N http://localhost:8080/v1/messages \
  -H "x-api-key: $XINGHAI_API_KEY" \
  -H 'anthropic-version: 2023-06-01' \
  -H 'Content-Type: application/json' \
  -d '{"model":"kimi-k2.6","max_tokens":1024,"messages":[{"role":"user","content":"Hello"}],"stream":true}'
```

OpenCode 配置示例：

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "xinghai": {
      "npm": "@ai-sdk/anthropic",
      "name": "Xinghai Router",
      "options": { "baseURL": "http://localhost:8080/v1", "apiKey": "sk-xh-your-key" },
      "models": { "kimi-k2.6": { "name": "Kimi K2.6" } }
    }
  }
}
```

The router selects an enabled channel advertising the requested model. It tries the lowest numeric priority first, distributes equal-priority traffic by weight, and retries a different eligible channel for connection errors, `408`, `425`, `429`, and `5xx` responses.

## Verify

Run `go test ./...` and `go vet ./...`.
