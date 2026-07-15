<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { Activity, Bot, Check, ChevronRight, CircleAlert, Copy, KeyRound, LayoutDashboard, LogOut, Plus, RadioTower, RefreshCw, TerminalSquare, Users } from 'lucide-vue-next'
import { api, clearToken, getToken, setToken } from './api'
import type { ApiKey, Channel, RequestLog, User } from './api'

type View = 'overview' | 'users' | 'keys' | 'channels' | 'logs'
const view = ref<View>('overview')
const authenticated = ref(Boolean(getToken()))
const adminToken = ref('')
const error = ref('')
const busy = ref(false)
const users = ref<User[]>([])
const keys = ref<ApiKey[]>([])
const channels = ref<Channel[]>([])
const logs = ref<RequestLog[]>([])
const createdKey = ref('')
const showUser = ref(false)
const showKey = ref(false)
const showChannel = ref(false)
const userForm = reactive({ email: '', name: '', role: 'user' })
const keyForm = reactive({ user_id: '', name: '', expires_at: '' })
const channelForm = reactive({ name: '', base_url: 'https://api.openai.com', api_key: '', models: '', priority: 100 })

const nav = [
  ['overview', '概览', LayoutDashboard], ['users', '用户', Users], ['keys', 'API 密钥', KeyRound], ['channels', '渠道', RadioTower], ['logs', '请求日志', TerminalSquare],
] as const
const activeChannels = computed(() => channels.value.filter((channel) => channel.enabled).length)
const successRate = computed(() => logs.value.length ? Math.round(logs.value.filter((log) => log.status_code < 400).length / logs.value.length * 100) : 100)
const totalTokens = computed(() => logs.value.reduce((sum, log) => sum + (log.total_tokens ?? 0), 0))
const userName = (id: string | null) => users.value.find((user) => user.id === id)?.name ?? '已删除用户'
const formatDate = (value: string | null) => value ? new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '从未'
const short = (value: string | null) => value ? `${value.slice(0, 8)}...` : '---'

async function load() {
  busy.value = true; error.value = ''
  try {
    const [userResponse, keyResponse, channelResponse, logResponse] = await Promise.all([
      api<{ data: User[] }>('/admin/users'), api<{ data: ApiKey[] }>('/admin/keys'), api<{ data: Channel[] }>('/admin/channels'), api<{ data: RequestLog[] }>('/admin/request-logs'),
    ])
    users.value = userResponse.data; keys.value = keyResponse.data; channels.value = channelResponse.data; logs.value = logResponse.data
  } catch (cause) { error.value = cause instanceof Error ? cause.message : '加载失败' } finally { busy.value = false }
}
async function signIn() { setToken(adminToken.value); authenticated.value = true; await load(); if (error.value) { clearToken(); authenticated.value = false } }
function signOut() { clearToken(); authenticated.value = false; adminToken.value = ''; error.value = '' }
async function createUser() { await action(async () => { await api('/admin/users', { method: 'POST', body: JSON.stringify(userForm) }); showUser.value = false; Object.assign(userForm, { email: '', name: '', role: 'user' }); await load() }) }
async function createKey() { await action(async () => { const response = await api<{ key: string }>('/admin/keys', { method: 'POST', body: JSON.stringify({ ...keyForm, expires_at: keyForm.expires_at ? new Date(keyForm.expires_at).toISOString() : '' }) }); createdKey.value = response.key; showKey.value = false; Object.assign(keyForm, { user_id: '', name: '', expires_at: '' }); await load() }) }
async function createChannel() { await action(async () => { await api('/admin/channels', { method: 'POST', body: JSON.stringify({ ...channelForm, models: channelForm.models.split(',').map((value) => value.trim()).filter(Boolean) }) }); showChannel.value = false; Object.assign(channelForm, { name: '', base_url: 'https://api.openai.com', api_key: '', models: '', priority: 100 }); await load() }) }
async function toggleChannel(channel: Channel) { await action(async () => { await api(`/admin/channels/${channel.id}/status`, { method: 'POST', body: JSON.stringify({ enabled: !channel.enabled }) }); await load() }) }
async function revokeKey(key: ApiKey) { if (!confirm(`吊销 ${key.key_prefix} 的访问权限？`)) return; await action(async () => { await api(`/admin/keys/${key.id}/revoke`, { method: 'POST' }); await load() }) }
async function action(work: () => Promise<void>) { busy.value = true; error.value = ''; try { await work() } catch (cause) { error.value = cause instanceof Error ? cause.message : '操作失败' } finally { busy.value = false } }
async function copyKey() { await navigator.clipboard.writeText(createdKey.value) }
onMounted(() => { if (authenticated.value) load() })
</script>

<template>
  <main v-if="!authenticated" class="login-shell">
    <section class="login-card">
      <div class="brand-mark"><Bot :size="29" /></div>
      <p class="eyebrow">XINGHAI ROUTER</p>
      <h1>控制模型流量。</h1>
      <p class="muted">输入部署时配置的管理员令牌，进入网关控制台。</p>
      <form @submit.prevent="signIn">
        <label>管理员令牌<input v-model="adminToken" type="password" autocomplete="current-password" placeholder="ADMIN_TOKEN" required /></label>
        <button class="button primary full" :disabled="busy">进入控制台 <ChevronRight :size="16" /></button>
      </form>
      <p v-if="error" class="error"><CircleAlert :size="16" />{{ error }}</p>
    </section>
  </main>

  <main v-else class="app-shell">
    <aside class="sidebar">
      <div class="logo"><span class="brand-mark small"><Bot :size="19" /></span><span>Xinghai</span><i>Router</i></div>
      <nav><button v-for="[id, label, Icon] in nav" :key="id" :class="{ active: view === id }" @click="view = id"><component :is="Icon" :size="18" />{{ label }}</button></nav>
      <div class="sidebar-footer"><span><span class="live-dot"></span>网关在线</span><button @click="signOut"><LogOut :size="16" />退出</button></div>
    </aside>
    <section class="content">
      <header><div><p class="eyebrow">控制台</p><h1>{{ nav.find((item) => item[0] === view)?.[1] }}</h1></div><button class="button ghost" @click="load" :disabled="busy"><RefreshCw :size="16" :class="{ spinning: busy }" />刷新</button></header>
      <p v-if="error" class="error banner"><CircleAlert :size="16" />{{ error }}</p>

      <template v-if="view === 'overview'">
        <div class="metrics"><article><span>活跃渠道</span><strong>{{ activeChannels }}<em>/{{ channels.length }}</em></strong><p><RadioTower :size="15" />可用于路由</p></article><article><span>近 100 请求</span><strong>{{ logs.length }}</strong><p><Activity :size="15" />最近记录</p></article><article><span>成功率</span><strong>{{ successRate }}<em>%</em></strong><p><Check :size="15" />HTTP 2xx / 3xx</p></article><article><span>计量 Token</span><strong>{{ totalTokens.toLocaleString() }}</strong><p><KeyRound :size="15" />非流式请求</p></article></div>
        <div class="grid-two"><section class="panel"><div class="panel-title"><div><h2>渠道状态</h2><p>按优先级选取首个可用渠道</p></div><button class="text-button" @click="view = 'channels'">管理</button></div><div v-if="channels.length" class="channel-list"><div v-for="channel in channels" :key="channel.id"><span :class="['status-dot', { off: !channel.enabled }]"></span><div><b>{{ channel.name }}</b><small>{{ channel.models.join(', ') }}</small></div><span class="priority">P{{ channel.priority }}</span></div></div><Empty v-else text="尚未配置上游渠道" /></section><section class="panel"><div class="panel-title"><div><h2>最近请求</h2><p>最近 100 条网关请求</p></div><button class="text-button" @click="view = 'logs'">查看全部</button></div><div v-if="logs.length" class="compact-list"><div v-for="log in logs.slice(0, 5)" :key="log.request_id"><code>{{ log.model }}</code><span>{{ log.duration_ms }} ms</span><b :class="log.status_code < 400 ? 'success' : 'danger'">{{ log.status_code }}</b></div></div><Empty v-else text="等待第一条模型请求" /></section></div>
      </template>

      <template v-if="view === 'users'"><section class="toolbar"><div><h2>用户</h2><p>管理可创建 API 密钥的账户。</p></div><button class="button primary" @click="showUser = true"><Plus :size="16" />创建用户</button></section><section class="panel table-panel"><table><thead><tr><th>用户</th><th>角色</th><th>状态</th><th>创建时间</th></tr></thead><tbody><tr v-for="user in users" :key="user.id"><td><b>{{ user.name }}</b><small>{{ user.email }}</small></td><td><span class="pill">{{ user.role }}</span></td><td><span class="state good">已启用</span></td><td>{{ formatDate(user.created_at) }}</td></tr></tbody></table><Empty v-if="!users.length" text="还没有用户" /></section></template>
      <template v-if="view === 'keys'"><section class="toolbar"><div><h2>API 密钥</h2><p>仅在创建时显示一次完整密钥。</p></div><button class="button primary" :disabled="!users.length" @click="showKey = true"><Plus :size="16" />创建密钥</button></section><section class="panel table-panel"><table><thead><tr><th>名称</th><th>所属用户</th><th>前缀</th><th>最后使用</th><th>状态</th><th></th></tr></thead><tbody><tr v-for="key in keys" :key="key.id"><td><b>{{ key.name }}</b></td><td>{{ userName(key.user_id) }}</td><td><code>{{ key.key_prefix }}...</code></td><td>{{ formatDate(key.last_used_at) }}</td><td><span :class="['state', key.revoked_at ? 'bad' : 'good']">{{ key.revoked_at ? '已吊销' : '有效' }}</span></td><td><button v-if="!key.revoked_at" class="text-button danger" @click="revokeKey(key)">吊销</button></td></tr></tbody></table><Empty v-if="!keys.length" text="创建用户后，即可签发 API 密钥" /></section></template>
      <template v-if="view === 'channels'"><section class="toolbar"><div><h2>上游渠道</h2><p>模型请求按渠道优先级进行选择。</p></div><button class="button primary" @click="showChannel = true"><Plus :size="16" />添加渠道</button></section><div class="channel-cards"><article v-for="channel in channels" :key="channel.id" class="panel channel-card"><div class="card-top"><span :class="['status-dot', { off: !channel.enabled }]"></span><span>优先级 {{ channel.priority }}</span><button class="toggle" :class="{ on: channel.enabled }" @click="toggleChannel(channel)"><i></i></button></div><h3>{{ channel.name }}</h3><p class="url">{{ channel.base_url }}</p><div class="model-tags"><span v-for="model in channel.models" :key="model">{{ model }}</span></div></article><Empty v-if="!channels.length" text="添加 OpenAI-compatible 上游开始路由" /></div></template>
      <template v-if="view === 'logs'"><section class="toolbar"><div><h2>请求日志</h2><p>最多显示最新 100 条记录。</p></div></section><section class="panel table-panel"><table><thead><tr><th>时间</th><th>模型</th><th>状态</th><th>耗时</th><th>Token</th><th>请求 ID</th></tr></thead><tbody><tr v-for="log in logs" :key="log.request_id"><td>{{ formatDate(log.created_at) }}</td><td><code>{{ log.model }}</code></td><td><span :class="['state', log.status_code < 400 ? 'good' : 'bad']">{{ log.status_code }}</span></td><td>{{ log.duration_ms }} ms</td><td>{{ log.total_tokens ?? 0 }}</td><td><code>{{ short(log.request_id) }}</code></td></tr></tbody></table><Empty v-if="!logs.length" text="暂无请求日志" /></section></template>
    </section>

    <div v-if="showUser || showKey || showChannel || createdKey" class="modal-backdrop" @click.self="showUser = showKey = showChannel = false">
      <form v-if="showUser" class="modal" @submit.prevent="createUser"><div class="modal-title"><h2>创建用户</h2><button type="button" @click="showUser = false">×</button></div><label>姓名<input v-model="userForm.name" required placeholder="例如：李雷" /></label><label>邮箱<input v-model="userForm.email" required type="email" placeholder="name@example.com" /></label><label>角色<select v-model="userForm.role"><option value="user">用户</option><option value="operator">运营</option><option value="admin">管理员</option></select></label><button class="button primary full" :disabled="busy">创建用户</button></form>
      <form v-if="showKey" class="modal" @submit.prevent="createKey"><div class="modal-title"><h2>创建 API 密钥</h2><button type="button" @click="showKey = false">×</button></div><label>用户<select v-model="keyForm.user_id" required><option disabled value="">选择用户</option><option v-for="user in users" :key="user.id" :value="user.id">{{ user.name }} · {{ user.email }}</option></select></label><label>密钥名称<input v-model="keyForm.name" required placeholder="例如：生产环境" /></label><label>过期时间 <small>可选</small><input v-model="keyForm.expires_at" type="datetime-local" /></label><button class="button primary full" :disabled="busy">签发密钥</button></form>
      <form v-if="showChannel" class="modal" @submit.prevent="createChannel"><div class="modal-title"><h2>添加上游渠道</h2><button type="button" @click="showChannel = false">×</button></div><label>名称<input v-model="channelForm.name" required placeholder="例如：OpenAI 主线路" /></label><label>Base URL<input v-model="channelForm.base_url" required type="url" /></label><label>上游 API Key<input v-model="channelForm.api_key" required type="password" /></label><label>模型 <small>逗号分隔</small><input v-model="channelForm.models" required placeholder="gpt-4o-mini, gpt-4o" /></label><label>优先级 <input v-model.number="channelForm.priority" type="number" min="0" /></label><button class="button primary full" :disabled="busy">保存渠道</button></form>
      <section v-if="createdKey" class="modal secret"><div class="modal-title"><h2>保存 API 密钥</h2><button @click="createdKey = ''">×</button></div><p>完整密钥只显示这一次，请立即保存。</p><code>{{ createdKey }}</code><button class="button primary full" @click="copyKey"><Copy :size="16" />复制密钥</button><button class="button ghost full" @click="createdKey = ''">我已保存</button></section>
    </div>
  </main>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({ components: { Empty: { props: { text: { type: String, required: true } }, template: '<div class="empty">{{ text }}</div>' } } })
</script>
