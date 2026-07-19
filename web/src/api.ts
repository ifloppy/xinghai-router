export interface User { id: string; email: string; name: string; role: string; enabled: boolean; balance: number; reserved: number; permissions: string[]; groups: string[]; created_at: string }
export interface ApiKey { id: string; user_id: string; name: string; key_prefix: string; group_id: string; group_name: string; expires_at: string | null; revoked_at: string | null; last_used_at: string | null; created_at: string }
export interface Channel { id: string; name: string; base_url: string; provider: 'openai' | 'ollama' | 'kimi' | 'opencode_go' | 'anthropic'; models: string[]; enabled: boolean; auto_disabled: boolean; disabled_reason: string; priority: number; groups: string[]; created_at: string }
export interface Group { id: string; name: string; multiplier: number; created_at: string }
export interface RequestLog { request_id: string; user_id: string | null; api_key_id: string | null; channel_id: string | null; model: string; status_code: number; prompt_tokens: number | null; completion_tokens: number | null; total_tokens: number | null; duration_ms: number; error_code: string | null; created_at: string }
export interface Account { id: string; email: string; name: string; role: string; avatar_url: string; permissions: string[]; balance: number; reserved: number; leaderboard_opt_in: boolean; leaderboard_mask_name: boolean }
export interface Pricing { id: string; model: string; input_per_million: number; cached_input_per_million: number; output_per_million: number; multiplier: number; enabled: boolean; updated_at: string }
export interface CatalogGroup { id: string; name: string; multiplier: number }
export interface CatalogModel { id: string; model: string; provider: string; provider_slug: string; input_per_million: number | null; cached_input_per_million: number | null; output_per_million: number | null; multiplier: number | null; groups: CatalogGroup[] }
export interface ModelProvider { id: string; name: string; slug: string; prefixes: string[]; priority: number }
export interface UsageRecord { request_id: string; model: string; prompt_tokens: number; cached_prompt_tokens: number; completion_tokens: number; cost: number; status: string; created_at: string }
export interface ActivityLog { id: string; type: 'request' | 'login' | 'register' | 'logout' | 'topup' | 'operation'; action: string; user_id: string; user_name: string; model: string; group_id: string; group_name: string; status_code: number | null; duration_ms: number | null; prompt_tokens: number; completion_tokens: number; total_tokens: number; cost: number; details: Record<string, unknown>; created_at: string }
export interface LedgerEntry { id: string; amount: number; balance_after: number; kind: string; request_id: string | null; note: string | null; created_at: string }
export interface PaymentOrder { order_no: string; payment_type: string; amount: string; status: 'pending' | 'paid' | 'failed' | 'expired'; provider_trade_no?: string; paid_at: string | null; created_at: string }
export interface PaymentMethod { id: string; code: string; name: string; enabled: boolean; created_at: string }
export interface PaymentSettings { enabled: boolean; base_url: string; merchant_id: string; has_merchant_key: boolean; public_base_url: string; methods: PaymentMethod[] }
export interface ModelRanking { rank: number; previous_rank?: number; model_name: string; vendor: string; total_tokens: number; share: number; growth_pct: number }
export interface VendorRanking { rank: number; vendor: string; total_tokens: number; share: number; growth_pct: number; models_count: number; top_model: string }
export interface RankingMover { model_name: string; vendor: string; rank_delta: number; current_rank: number; growth_pct: number }
export interface UserRanking { rank: number; name: string; total_tokens: number; total_cost: number; share: number; growth_pct: number; requests: number; top_model: string }
export interface Rankings { period: string; models: ModelRanking[]; vendors: VendorRanking[]; top_movers: RankingMover[]; top_droppers: RankingMover[]; users: UserRanking[]; total_tokens: number; updated_at: string }
export interface SiteSettings { name: string; icon_url: string; auto_disable_failed_channels: boolean; geetest_enabled?: boolean; geetest_captcha_id?: string; email_verification_enabled?: boolean }
export interface AdminSiteSettings { name: string; icon_url: string; auto_disable_failed_channels: boolean; geetest_captcha_id: string; has_geetest_captcha_key: boolean; smtp_host: string; smtp_port: string; smtp_username: string; has_smtp_password: boolean; smtp_from: string }
export interface ReliabilitySettings { retry_count: number; retry_status_codes: string; health_check_mode: 'off' | 'scheduled_all' | 'passive_recovery'; health_check_interval_minutes: number; health_check_auto_recover: boolean; health_check_channel_ids: string; auto_disable_on_test_failure: boolean; auto_disable_slow_seconds: number; auto_disable_status_codes: string; auto_disable_keywords: string }

export interface SubscriptionPlan {
  id: string
  name: string
  description: string
  price: string
  currency: string
  billing_period: 'month' | 'year'
  credit_amount: string
  group_id: string
  group_name: string
  model_whitelist: string[]
  max_requests_per_period: number | null
  max_tokens_per_period: number | null
  sort_order: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface PublicSubscriptionPlan {
  id: string
  name: string
  description: string
  price: string
  currency: string
  billing_period: 'month' | 'year'
  credit_amount: string
  group_name: string
  model_whitelist: string[]
  sort_order: number
}

export interface UserSubscription {
  id: string
  user_id: string
  plan_id: string
  plan_name: string
  status: 'pending' | 'active' | 'expired' | 'cancelled'
  current_period_start: string | null
  current_period_end: string | null
  auto_renew: boolean
  cancelled_at: string | null
  created_at: string
  updated_at: string
  price: string
  billing_period: 'month' | 'year'
  credit_amount: string
  group_id: string
  group_name: string
  model_whitelist: string[]
  max_requests_per_period: number | null
  max_tokens_per_period: number | null
}

export interface SubscriptionOrder {
  id: string
  order_no: string
  subscription_id: string
  plan_id: string
  plan_name: string
  provider: string
  payment_type: string
  amount: string
  status: 'pending' | 'paid' | 'failed' | 'expired'
  provider_trade_no?: string
  period_kind: 'new' | 'renewal'
  paid_at: string | null
  created_at: string
}

export interface AdminSubscription {
  id: string
  user_id: string
  email: string
  user_name: string
  plan_id: string
  plan_name: string
  status: 'pending' | 'active' | 'expired' | 'cancelled'
  current_period_start: string | null
  current_period_end: string | null
  auto_renew: boolean
  cancelled_at: string | null
  created_at: string
  updated_at: string
}

let token = import.meta.client ? sessionStorage.getItem('xinghai.admin-token') ?? '' : ''
export const getToken = () => token
export const setToken = (value: string) => { token = value.trim(); sessionStorage.setItem('xinghai.admin-token', token) }
export const clearToken = () => { token = ''; sessionStorage.removeItem('xinghai.admin-token') }

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`/api${path}`, { ...init, headers: { Authorization: `Bearer ${token}`, ...(init.body ? { 'Content-Type': 'application/json' } : {}), ...init.headers } })
  if (!response.ok) {
    const body = await response.json().catch(() => null)
    throw new Error(body?.error?.message ?? `请求失败 (${response.status})`)
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}
