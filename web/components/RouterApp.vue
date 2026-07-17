<script setup lang="ts">
import { computed, h, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { Activity, Bot, Check, ChevronDown, ChevronRight, CircleAlert, Copy, KeyRound, Layers3, LayoutDashboard, LogOut, PanelLeftClose, PanelLeftOpen, Plus, RadioTower, RefreshCw, Search, ShieldCheck, Sparkles, TerminalSquare, UserRound, Users, WalletCards, ReceiptText, Tags, Settings } from 'lucide-vue-next'
import { api, clearToken, getToken, setToken } from '~/src/api'
import type { Account, ActivityLog, ApiKey, CatalogGroup, CatalogModel, Channel, Group, LedgerEntry, ModelProvider, PaymentMethod, PaymentOrder, PaymentSettings, Pricing, ReliabilitySettings, SiteSettings, UsageRecord, User } from '~/src/api'
const { locale, t, setLocale, toggleLocale, initializeLocale } = useI18n()

type View = 'overview' | 'users' | 'groups' | 'keys' | 'channels' | 'providers' | 'logs' | 'account' | 'profile' | 'wallet' | 'usage' | 'usage-overview' | 'ledger' | 'pricing' | 'site-settings' | 'payment-settings' | 'audit' | 'reliability'
const props = withDefaults(defineProps<{ activeView?: View }>(), { activeView: 'overview' })
const route = useRoute()
const router = useRouter()
const views: View[] = ['overview', 'users', 'groups', 'keys', 'channels', 'providers', 'logs', 'account', 'profile', 'wallet', 'usage', 'usage-overview', 'ledger', 'pricing', 'site-settings', 'payment-settings', 'audit', 'reliability']
const view = computed<View>(() => {
  const selected = views.includes(route.query.view as View) ? route.query.view as View : props.activeView && views.includes(props.activeView) ? props.activeView : views.includes(route.params.view as View) ? route.params.view as View : 'overview'
  return selected === 'logs' || selected === 'audit' || selected === 'usage-overview' ? 'usage' : selected
})
const authenticated = ref(false)
const error = ref('')
const errorAlert = ref<HTMLElement | null>(null)
const errorHovered = ref(false)
const errorSelected = ref(false)
let errorTimer: ReturnType<typeof setTimeout> | undefined
const busy = ref(false)
const sidebarCollapsed = ref(false)
const setupCollapsed = ref(false)
const users = ref<User[]>([])
const groups = ref<Group[]>([])
const ownGroups = ref<string[]>([])
const keys = ref<ApiKey[]>([])
const accountKeys = ref<ApiKey[]>([])
const channels = ref<Channel[]>([])
const providers = ref<ModelProvider[]>([])
const activityLogs = ref<ActivityLog[]>([])
const account = ref<Account | null>(null)
const usageRecords = ref<UsageRecord[]>([])
const ledger = ref<LedgerEntry[]>([])
const payments = ref<PaymentOrder[]>([])
const paymentMethods = ref<PaymentMethod[]>([])
const paymentsEnabled = ref(false)
const paymentMessage = ref('')
const paymentForm = reactive({ amount: 10, type: '' })
const paymentSettings = reactive<PaymentSettings>({ enabled: false, base_url: '', merchant_id: '', has_merchant_key: false, public_base_url: '', methods: [] })
const paymentSettingsForm = reactive({ enabled: false, base_url: '', merchant_id: '', merchant_key: '', public_base_url: '' })
const paymentMethodForm = reactive({ code: '', name: '', enabled: true })
const pricing = ref<Pricing[]>([])
const siteSettings = ref<SiteSettings>({ name: 'Xinghai Router', icon_url: '', auto_disable_failed_channels: false })
const catalog = ref<CatalogModel[]>([])
const catalogGroups = ref<CatalogGroup[]>([])
const catalogGroup = ref('all')
const catalogSearch = ref('')
const activityModels = ref<string[]>([])
const activityFilters = reactive({ user_id: '', model: '', group_id: '', start: '', end: '', type: '' })
const createdKey = ref('')
const showKey = ref(false)
const showAccountKey = ref(false)
const editingAccountKey = ref<ApiKey | null>(null)
const showChannel = ref(false)
const editingChannel = ref<Channel | null>(null)
const showProvider = ref(false)
const selectedUser = ref<User | null>(null)
const originalUser = ref<User | null>(null)
const selectedPermissions = ref<string[]>([])
const selectedGroups = ref<string[]>([])
const userPassword = ref('')
const userBalance = ref<number | null>(null)
const userBalanceNote = ref('')
const keyForm = reactive({ user_id: '', name: '', expires_at: '', group_id: '' })
const accountKeyForm = reactive({ name: '', expires_at: '', group_id: '' })
const channelForm = reactive({ name: '', provider: 'openai', base_url: 'https://api.openai.com', api_key: '', models: '', priority: 100, groups: [] as string[] })
const providerForm = reactive({ name: '', slug: '', prefixes: '', priority: 100 })
const groupForm = reactive({ name: '', multiplier: 1 })
const groupImportText = ref('')

const generalNav = computed(() => [['overview', t('overview'), LayoutDashboard], ['account', t('account'), KeyRound], ['usage-overview', t('usageOverview'), Activity], ['usage', t('usage'), TerminalSquare]] as const)
const billingNav = computed(() => [['wallet', t('wallet'), WalletCards], ['ledger', t('ledger'), ReceiptText]] as const)
const personalNav = computed(() => [['profile', t('profile'), UserRound]] as const)
const managementNavItems = [
  ['users', 'users', Users, 'users.read'], ['groups', 'groups', Layers3, 'system.manage'], ['keys', 'keys', KeyRound, 'keys.manage'], ['channels', 'channels', RadioTower, 'channels.read'], ['providers', 'providers', Tags, 'system.manage'],
] as const
const adminExtraNav = [['pricing', 'pricing', Tags, 'pricing.read'], ['reliability', 'reliability', ShieldCheck, 'system.manage'], ['site-settings', 'siteSettings', Settings, 'system.manage'], ['payment-settings', 'paymentSettings', WalletCards, 'system.manage']] as const
const localizedManagementNavItems = computed(() => managementNavItems.map(([id, key, icon, permission]) => [id, t(key as Parameters<typeof t>[0]), icon, permission] as const))
const localizedAdminExtraNav = computed(() => adminExtraNav.map(([id, key, icon, permission]) => [id, t(key as Parameters<typeof t>[0]), icon, permission] as const))
const permissions = ['users.read', 'users.manage', 'keys.manage', 'channels.read', 'channels.manage', 'logs.read', 'pricing.read', 'pricing.manage', 'audit.read', 'wallets.manage', 'routes.manage', 'quotas.manage', 'system.manage']
const pricingForm = reactive({ model: '', input_per_million: 0, cached_input_per_million: 0, output_per_million: 0, multiplier: 1 })
const siteSettingsForm = reactive<SiteSettings>({ name: '', icon_url: '', auto_disable_failed_channels: false })
const reliabilityForm = reactive<ReliabilitySettings>({ retry_count: 3, retry_status_codes: '', health_check_mode: 'off', health_check_interval_minutes: 5, health_check_auto_recover: true, health_check_channel_ids: '', auto_disable_on_test_failure: false, auto_disable_slow_seconds: 0, auto_disable_status_codes: '', auto_disable_keywords: '' })
const newAPIPricingForm = reactive({ base_url: '', api_key: '', price_per_quota_unit: 1 })
const loginMode = ref<'token' | 'login' | 'register'>('token')
const accountForm = reactive({ name: '', email: '', password: '' })
const avatarInput = ref<HTMLInputElement | null>(null)
const avatarUrlInput = ref('')
const isLanding = computed(() => route.path === '/')
const isAuthPage = computed(() => route.path === '/auth')
const isMarketplacePage = computed(() => route.path === '/models')
const activeChannels = computed(() => channels.value.filter((channel) => channel.enabled).length)
const isAdmin = computed(() => account.value?.role === 'admin')
const can = (permission: string) => isAdmin.value || Boolean(account.value?.permissions.includes(permission))
const managementNav = computed(() => [...localizedManagementNavItems.value, ...localizedAdminExtraNav.value].filter((item) => can(item[3])))
const personalRequests = computed(() => usageRecords.value.length)
const personalTokens = computed(() => usageRecords.value.reduce((sum, item) => sum + item.prompt_tokens + item.completion_tokens, 0))
const personalCost = computed(() => usageRecords.value.reduce((sum, item) => sum + Number(item.cost), 0))
const setupProgress = computed(() => [accountKeys.value.some((item) => !item.revoked_at), Number(account.value?.balance ?? 0) > 0, personalRequests.value > 0].filter(Boolean).length)
const filteredCatalog = computed(() => {
  const search = catalogSearch.value.trim().toLowerCase()
  return catalog.value.filter((item) => (!search || item.model.toLowerCase().includes(search)) && (catalogGroup.value === 'all' || item.groups.some((group) => group.id === catalogGroup.value)))
})
const apiEndpoint = computed(() => import.meta.client ? `${window.location.origin}/v1/chat/completions` : '/v1/chat/completions')
const usageChart = computed(() => {
  const days = Array.from({ length: 7 }, (_, index) => {
    const date = new Date()
    date.setHours(0, 0, 0, 0)
    date.setDate(date.getDate() - 6 + index)
    return { key: date.toISOString().slice(0, 10), label: new Intl.DateTimeFormat(locale.value === 'en-US' ? 'en-US' : 'zh-CN', { weekday: 'short' }).format(date), tokens: 0, cost: 0 }
  })
  const byDay = new Map(days.map((day) => [day.key, day]))
  for (const item of usageRecords.value) {
    const day = byDay.get(item.created_at.slice(0, 10))
    if (day) { day.tokens += item.prompt_tokens + item.completion_tokens; day.cost += Number(item.cost) }
  }
  const maxTokens = Math.max(...days.map((day) => day.tokens), 1)
  const maxCost = Math.max(...days.map((day) => day.cost), 1)
  return days.map((day) => ({ ...day, tokenHeight: Math.max(day.tokens ? 8 : 2, Math.round(day.tokens / maxTokens * 100)), costHeight: Math.max(day.cost ? 8 : 2, Math.round(day.cost / maxCost * 100)) }))
})
const usageLinePoints = computed(() => usageChart.value.map((day, index) => `${index * 100 / 6},${100 - day.tokenHeight}`).join(' '))
const userName = (id: string | null) => users.value.find((user) => user.id === id)?.name ?? t('deletedUser')
const formatDate = (value: string | null) => value ? new Intl.DateTimeFormat(locale.value === 'en-US' ? 'en-US' : 'zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : t('never')
const short = (value: string | null) => value ? `${value.slice(0, 8)}...` : '---'
const formatPrice = (value: number | null, multiplier = 1) => value == null ? t('pendingConfig') : `¥${(Number(value) * multiplier).toFixed(Number(value) * multiplier < 0.01 ? 4 : 2)}`
const providerIcon = (slug: string) => `https://unpkg.com/@lobehub/icons-static-svg@1.93.0/icons/${slug}.svg`
const modelProvider = (model: string) => catalog.value.find((item) => item.model === model)?.provider ?? t('other')
const selectedCatalogGroup = (item: CatalogModel) => item.groups.find((group) => group.id === catalogGroup.value) ?? item.groups[0]
const actualMultiplier = (item: CatalogModel) => Number(item.multiplier ?? 1) * Number(selectedCatalogGroup(item)?.multiplier ?? 1)
const Empty = (props: { text: string }) => h('div', { class: 'empty' }, props.text)
Empty.props = { text: { type: String, required: true } }

function clearErrorTimer() {
  if (errorTimer) window.clearTimeout(errorTimer)
  errorTimer = undefined
}

function scheduleErrorDismissal() {
  clearErrorTimer()
  if (!error.value || errorHovered.value || errorSelected.value) return
  errorTimer = window.setTimeout(() => { error.value = '' }, 5000)
}

function updateErrorSelection() {
  const selection = window.getSelection()
  const anchor = selection?.anchorNode
  const focus = selection?.focusNode
  errorSelected.value = Boolean(selection?.toString() && anchor && focus && errorAlert.value?.contains(anchor) && errorAlert.value?.contains(focus))
  scheduleErrorDismissal()
}

function lockError() {
  errorHovered.value = true
  clearErrorTimer()
}

function releaseError() {
  errorHovered.value = false
  updateErrorSelection()
  if (!errorSelected.value) error.value = ''
}

async function copyError() {
  if (!error.value) return
  if (navigator.clipboard) {
    await navigator.clipboard.writeText(error.value)
    return
  }
  const textarea = document.createElement('textarea')
  textarea.value = error.value
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.append(textarea)
  textarea.select()
  document.execCommand('copy')
  textarea.remove()
}

watch(error, () => {
  errorSelected.value = false
  scheduleErrorDismissal()
})

const activityTypeLabel = computed(() => ({ request: 'Request', login: 'Login', register: 'Register', logout: 'Logout', topup: 'Top-up', operation: t('otherOperation') } as Record<ActivityLog['type'], string>))
const actionLabel = (item: ActivityLog) => ({ 'account.logged_in': t('actionAccountLogin'), 'account.registered': t('actionAccountRegister'), 'account.logged_out': t('actionAccountLogout'), 'wallet.adjusted': t('actionWalletAdjusted') }[item.action] ?? item.action)
const activityDetail = (item: ActivityLog) => item.type === 'request' ? `${item.prompt_tokens} / ${item.completion_tokens} tokens · ${Number(item.cost).toFixed(6)}` : JSON.stringify(item.details)

async function loadActivity(filters = false) {
  const query = new URLSearchParams()
  if (filters) Object.entries(activityFilters).forEach(([key, value]) => { if (value) query.set(key, key === 'start' || key === 'end' ? new Date(value).toISOString() : value) })
  const value = await api<{ data: ActivityLog[] }>(`/activity-logs${query.size ? `?${query}` : ''}`)
  activityLogs.value = value.data
  if (!filters) activityModels.value = [...new Set(value.data.map((item) => item.model).filter(Boolean))].sort()
}
async function filterActivity() { await action(() => loadActivity(true)) }
async function resetActivityFilters() { Object.assign(activityFilters, { user_id: '', model: '', group_id: '', start: '', end: '', type: '' }); await action(() => loadActivity()) }

async function load() {
  busy.value = true; error.value = ''
  try {
    const [settings, me] = await Promise.all([api<SiteSettings>('/site-settings'), api<Account>('/account/me')])
    siteSettings.value = settings
    Object.assign(siteSettingsForm, settings)
    account.value = me
    const [ownKeys, ownUsage, ownLedger, ownGroupValue, ownPayments] = await Promise.all([
      api<{ data: ApiKey[] }>('/account/keys').catch(() => ({ data: [] })),
      api<{ data: UsageRecord[] }>('/account/usage').catch(() => ({ data: [] })),
      api<{ data: LedgerEntry[] }>('/account/ledger').catch(() => ({ data: [] })),
      api<{ data: string[]; groups: Group[] }>('/account/groups').catch(() => ({ data: [], groups: [] })),
      api<{ enabled: boolean; payment_methods: PaymentMethod[]; data: PaymentOrder[] }>('/account/payments').catch(() => ({ enabled: false, payment_methods: [], data: [] })),
    ])
    accountKeys.value = ownKeys.data; usageRecords.value = ownUsage.data; ledger.value = ownLedger.data; ownGroups.value = ownGroupValue.data
    paymentsEnabled.value = ownPayments.enabled; payments.value = ownPayments.data; paymentMethods.value = ownPayments.payment_methods ?? []
    if (!paymentMethods.value.some((method) => method.code === paymentForm.type)) paymentForm.type = paymentMethods.value[0]?.code ?? ''
    await loadActivity()
    if (!can('users.read')) groups.value = ownGroupValue.groups
    const requests: Promise<void>[] = []
    if (can('users.read')) requests.push(Promise.all([api<{ data: User[] }>('/admin/users'), api<{ data: Group[] }>('/admin/groups')]).then(([userValue, groupValue]) => { users.value = userValue.data; groups.value = groupValue.data }))
    if (can('keys.manage')) requests.push(api<{ data: ApiKey[] }>('/admin/keys').then((value) => { keys.value = value.data }))
    if (can('channels.read')) requests.push(api<{ data: Channel[] }>('/admin/channels').then((value) => { channels.value = value.data }))
    if (can('system.manage')) requests.push(api<{ data: ModelProvider[] }>('/admin/providers').then((value) => { providers.value = value.data }))
    if (can('system.manage')) requests.push(api<PaymentSettings>('/admin/payment-settings').then((value) => { Object.assign(paymentSettings, value); Object.assign(paymentSettingsForm, { enabled: value.enabled, base_url: value.base_url, merchant_id: value.merchant_id, merchant_key: '', public_base_url: value.public_base_url }) }))
    if (can('pricing.read')) requests.push(api<{ data: Pricing[] }>('/admin/pricing').then((value) => { pricing.value = value.data }))
    if (can('system.manage')) requests.push(api<ReliabilitySettings>('/admin/reliability-settings').then((value) => { Object.assign(reliabilityForm, value) }))
    await Promise.all(requests)
    } catch (cause) { error.value = cause instanceof Error ? cause.message : t('loadFailed') } finally { busy.value = false }
}
async function loadSiteSettings() { const value = await api<SiteSettings>('/site-settings'); siteSettings.value = value; Object.assign(siteSettingsForm, value); document.title = value.name; const link = document.querySelector<HTMLLinkElement>('link[rel="icon"]') ?? document.head.appendChild(Object.assign(document.createElement('link'), { rel: 'icon' })); if (value.icon_url) link.href = value.icon_url; else link.removeAttribute('href') }
async function saveSiteSettings() { await action(async () => { const value = await api<SiteSettings>('/admin/site-settings', { method: 'PUT', body: JSON.stringify(siteSettingsForm) }); siteSettings.value = value; await loadSiteSettings() }) }
async function saveReliabilitySettings() { await action(async () => { const value = await api<ReliabilitySettings>('/admin/reliability-settings', { method: 'PUT', body: JSON.stringify(reliabilityForm) }); Object.assign(reliabilityForm, value) }) }
async function loadCatalog() { const value = await api<{ data: CatalogModel[]; groups: CatalogGroup[] }>('/model-catalog'); catalog.value = value.data; catalogGroups.value = value.groups }
async function accountSignIn(register: boolean) { await action(async () => { const result = await api<{ token: string }>(register ? '/auth/register' : '/auth/login', { method: 'POST', body: JSON.stringify(register ? accountForm : { email: accountForm.email, password: accountForm.password }) }); setToken(result.token); authenticated.value = true; await load(); await router.replace({ path: '/console', query: { view: managementNav.value.length ? 'overview' : 'account' } }) }) }
async function signOut() { try { await api('/auth/logout', { method: 'POST' }) } catch { /* Local session removal is sufficient when the server is unreachable. */ } clearToken(); authenticated.value = false; error.value = ''; await router.replace('/') }
async function createKey() { await action(async () => { const response = await api<{ key: string }>('/admin/keys', { method: 'POST', body: JSON.stringify({ ...keyForm, expires_at: keyForm.expires_at ? new Date(keyForm.expires_at).toISOString() : '' }) }); createdKey.value = response.key; showKey.value = false; Object.assign(keyForm, { user_id: '', name: '', expires_at: '', group_id: '' }); await load() }) }
async function createAccountKey() { await action(async () => { const response = await api<{ key: string }>('/account/keys', { method: 'POST', body: JSON.stringify({ ...accountKeyForm, expires_at: accountKeyForm.expires_at ? new Date(accountKeyForm.expires_at).toISOString() : '' }) }); createdKey.value = response.key; showAccountKey.value = false; Object.assign(accountKeyForm, { name: '', expires_at: '', group_id: '' }); await load() }) }
function editAccountKey(key: ApiKey) { editingAccountKey.value = key; Object.assign(accountKeyForm, { name: key.name, expires_at: key.expires_at ? new Date(key.expires_at).toISOString().slice(0, 16) : '', group_id: key.group_id }) }
async function updateAccountKey() { if (!editingAccountKey.value) return; await action(async () => { await api(`/account/keys/${editingAccountKey.value.id}`, { method: 'PUT', body: JSON.stringify({ ...accountKeyForm, expires_at: accountKeyForm.expires_at ? new Date(accountKeyForm.expires_at).toISOString() : '' }) }); editingAccountKey.value = null; Object.assign(accountKeyForm, { name: '', expires_at: '', group_id: '' }); await load() }) }
async function fetchChannelModels() { await action(async () => { const response = await api<{ models: string[] }>('/admin/channels/models', { method: 'POST', body: JSON.stringify({ base_url: channelForm.base_url, api_key: channelForm.api_key }) }); channelForm.models = response.models.join(', ');     if (!response.models.length) throw new Error(t('upstreamNoModels')) }) }
async function createChannel() { await action(async () => { await api('/admin/channels', { method: 'POST', body: JSON.stringify({ ...channelForm, models: channelForm.models.split(',').map((value) => value.trim()).filter(Boolean) }) }); showChannel.value = false; Object.assign(channelForm, { name: '', provider: 'openai', base_url: 'https://api.openai.com', api_key: '', models: '', priority: 100, groups: [] }); await load() }) }
function editChannel(channel: Channel) { editingChannel.value = channel; Object.assign(channelForm, { name: channel.name, provider: channel.provider, base_url: channel.base_url, api_key: '', models: channel.models.join(', '), priority: channel.priority, groups: [...channel.groups] }) }
async function updateChannel() { if (!editingChannel.value) return; await action(async () => { const models = channelForm.models.split(',').map((value) => value.trim()).filter(Boolean); await api(`/admin/channels/${editingChannel.value.id}`, { method: 'PUT', body: JSON.stringify({ ...channelForm, models }) }); await api(`/admin/channels/${editingChannel.value.id}/groups`, { method: 'PUT', body: JSON.stringify({ groups: channelForm.groups }) }); editingChannel.value = null; Object.assign(channelForm, { name: '', provider: 'openai', base_url: 'https://api.openai.com', api_key: '', models: '', priority: 100, groups: [] }); await load() }) }
async function saveProvider() { await action(async () => { await api('/admin/providers', { method: 'POST', body: JSON.stringify({ ...providerForm, prefixes: providerForm.prefixes.split(',').map((value) => value.trim()).filter(Boolean) }) }); showProvider.value = false; Object.assign(providerForm, { name: '', slug: '', prefixes: '', priority: 100 }); await load() }) }
async function createGroup() { const name = groupForm.name.trim(); const multiplier = Number(groupForm.multiplier);     if (!name) { error.value = t('enterGroupName'); return } if (!Number.isFinite(multiplier) || multiplier <= 0) { error.value = t('multiplierMustBePositive'); return } await action(async () => { await api('/admin/groups', { method: 'POST', body: JSON.stringify({ name, multiplier }) }); Object.assign(groupForm, { name: '', multiplier: 1 }); await load() }) }
async function editGroupMultiplier(group: Group, event: Event) { const multiplier = Number(new FormData(event.currentTarget as HTMLFormElement).get('multiplier'));     if (!Number.isFinite(multiplier) || multiplier < 0) { error.value = t('multiplierMustBeNonNegative'); return } await action(async () => { await api(`/admin/groups/${group.id}`, { method: 'PUT', body: JSON.stringify({ multiplier }) }); await load() }) }
async function importGroups() { let values: Record<string, unknown>; try { values = JSON.parse(groupImportText.value) } catch { error.value = t('enterValidJSON'); return } if (!values || Array.isArray(values) || typeof values !== 'object') { error.value = t('importMustBeGroupMultiplierJSON'); return } const entries = Object.entries(values); if (!entries.length || entries.some(([name, multiplier]) => !name.trim() || typeof multiplier !== 'number' || !Number.isFinite(multiplier) || multiplier < 0)) { error.value = t('groupNameRequiredMultiplierNonNegative'); return } await action(async () => { await api('/admin/groups/import', { method: 'POST', body: JSON.stringify(values) }); groupImportText.value = ''; await load() }) }
async function toggleChannel(channel: Channel) { await action(async () => { await api(`/admin/channels/${channel.id}/status`, { method: 'POST', body: JSON.stringify({ enabled: !channel.enabled }) }); await load() }) }
async function revokeKey(key: ApiKey) { if (!confirm(t('revokeKeyConfirm').replace('{prefix}', key.key_prefix))) return; await action(async () => { await api(`/admin/keys/${key.id}/revoke`, { method: 'POST' }); await load() }) }
async function action(work: () => Promise<void>) { busy.value = true; error.value = ''; try { await work() } catch (cause) { error.value = cause instanceof Error ? cause.message : t('operationFailed') } finally { busy.value = false } }
async function createPayment() {
  const amount = Number(paymentForm.amount)
  if (!Number.isFinite(amount) || amount < 1 || amount > 100000) { error.value = t('paymentAmountRange'); return }
  await action(async () => {
    const result = await api<{ pay_url: string }>('/account/payments', { method: 'POST', body: JSON.stringify({ amount: amount.toFixed(2), type: paymentForm.type }) })
    window.location.assign(result.pay_url)
  })
}
async function savePaymentSettings() { await action(async () => { const value = await api<PaymentSettings>('/admin/payment-settings', { method: 'PUT', body: JSON.stringify(paymentSettingsForm) }); Object.assign(paymentSettings, value); paymentSettingsForm.merchant_key = ''; await load() }) }
async function createPaymentMethod() { await action(async () => { await api('/admin/payment-methods', { method: 'POST', body: JSON.stringify(paymentMethodForm) }); Object.assign(paymentMethodForm, { code: '', name: '', enabled: true }); await load() }) }
async function updatePaymentMethod(method: PaymentMethod) { await action(async () => { await api(`/admin/payment-methods/${method.id}`, { method: 'PUT', body: JSON.stringify({ code: method.code, name: method.name, enabled: method.enabled }) }); await load() }) }
async function deletePaymentMethod(method: PaymentMethod) { if (!confirm(t('deletePaymentMethodConfirm').replace('{name}', method.name))) return; await action(async () => { await api(`/admin/payment-methods/${method.id}`, { method: 'DELETE' }); await load() }) }
async function copyKey() { await navigator.clipboard.writeText(createdKey.value) }
async function savePricing() { await action(async () => { await api('/admin/pricing', { method: 'POST', body: JSON.stringify(pricingForm) }); Object.assign(pricingForm, { model: '', input_per_million: 0, cached_input_per_million: 0, output_per_million: 0, multiplier: 1 }); await load() }) }
async function syncNewAPIPricing() { await action(async () => { const result = await api<{ synced: number }>('/admin/pricing/newapi/sync', { method: 'POST', body: JSON.stringify(newAPIPricingForm) }); await load(); error.value = t('syncPricingResult').replace('{count}', String(result.synced)) }) }
function manageUser(user: User) { originalUser.value = user; selectedUser.value = { ...user }; selectedPermissions.value = [...user.permissions]; selectedGroups.value = [...(user.groups ?? [])]; userPassword.value = ''; userBalance.value = Number(user.balance ?? 0); userBalanceNote.value = '' }
async function saveUserAccess() {
  if (!selectedUser.value || !originalUser.value) return
  const current = selectedUser.value
  const original = originalUser.value
  const update: Record<string, unknown> = {}
  if (current.name !== original.name) update.name = current.name
  if (current.email !== original.email) update.email = current.email
  if (current.role !== original.role) update.role = current.role
  if (current.enabled !== original.enabled) update.enabled = current.enabled
  if (userPassword.value) update.password = userPassword.value
  if (Number(userBalance.value) !== Number(original.balance ?? 0)) { update.balance = userBalance.value; update.note = userBalanceNote.value }
  if ([...selectedPermissions.value].sort().join('\n') !== [...original.permissions].sort().join('\n')) update.permissions = selectedPermissions.value
  if ([...selectedGroups.value].sort().join('\n') !== [...(original.groups ?? [])].sort().join('\n')) update.groups = selectedGroups.value
  if (!Object.keys(update).length) { selectedUser.value = null; originalUser.value = null; return }
  await action(async () => { await api(`/admin/users/${current.id}`, { method: 'PUT', body: JSON.stringify(update) }); selectedUser.value = null; originalUser.value = null; await load() })
}
async function chooseAvatar(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  if (!['image/png', 'image/jpeg', 'image/gif', 'image/webp'].includes(file.type) || file.size > 1.5 * 1024 * 1024) { error.value = t('avatarFileRequirements'); return }
  await action(async () => {
    const avatarURL = await new Promise<string>((resolve, reject) => { const reader = new FileReader(); reader.onload = () => resolve(String(reader.result)); reader.onerror = () => reject(new Error(t('avatarReadFailed'))); reader.readAsDataURL(file) })
    await api('/account/profile', { method: 'PUT', body: JSON.stringify({ avatar_url: avatarURL }) })
    await load()
  })
  if (avatarInput.value) avatarInput.value.value = ''
}
async function removeAvatar() { await action(async () => { await api('/account/profile', { method: 'PUT', body: JSON.stringify({ avatar_url: '' }) }); await load() }) }
async function saveAvatarUrl() {
  const url = avatarUrlInput.value.trim()
  if (!url) return
  try {
    const parsed = new URL(url)
    if (!['http:', 'https:'].includes(parsed.protocol)) { error.value = t('avatarUrlInvalid'); return }
  } catch { error.value = t('avatarUrlInvalid'); return }
  await action(async () => {
    await api('/account/profile', { method: 'PUT', body: JSON.stringify({ avatar_url: url }) })
    avatarUrlInput.value = ''
    await load()
  })
}
function openAuth() { router.push('/auth') }
function openConsoleOrAuth() { router.push(authenticated.value ? '/console/overview' : '/auth') }
function closeAuth() { router.push('/') }
function openConsole(nextView: string) { if (views.includes(nextView as View)) router.push({ path: '/console', query: { view: nextView } }) }
onMounted(async () => {
  initializeLocale()
  document.addEventListener('selectionchange', updateErrorSelection)
  await loadSiteSettings().catch(() => undefined)
  authenticated.value = Boolean(getToken())
  if (isMarketplacePage.value) { await loadCatalog().catch((cause) => { error.value = cause instanceof Error ? cause.message : t('loadFailed') }); return }
  if (!authenticated.value) return
  await load()
  const returnedOrder = typeof route.query.payment_order === 'string' ? route.query.payment_order : ''
  if (returnedOrder) {
    const order = await api<PaymentOrder>(`/account/payments/${encodeURIComponent(returnedOrder)}`).catch(() => null)
    paymentMessage.value = order?.status === 'paid' ? t('paymentPaid') : t('paymentPending')
    if (order?.status === 'paid') await load()
  }
  if (route.meta.requiresAuth && error.value) {
    clearToken()
    authenticated.value = false
    await router.replace('/auth')
  }
})

onBeforeUnmount(() => {
  clearErrorTimer()
  document.removeEventListener('selectionchange', updateErrorSelection)
})
</script>

<template>
  <Transition name="error-alert">
    <div v-if="error" ref="errorAlert" class="error-alert" role="alert" tabindex="0" :title="t('clickToCopyError')" @mouseenter="lockError" @mouseleave="releaseError" @click="copyError" @keydown.enter.prevent="copyError" @keydown.space.prevent="copyError">
      <CircleAlert :size="17" /><span>{{ error }}</span><Copy :size="14" aria-hidden="true" />
    </div>
  </Transition>
  <main v-if="isLanding" class="landing-shell">
    <nav class="landing-nav">
       <a class="landing-logo" href="/"><span class="brand-mark small"><Bot :size="19" /></span><span>{{ siteSettings.name }}</span></a>
       <div class="landing-links"><a href="#features">{{ t('landingFeatures') }}</a><a href="/rankings">{{ t('rankings') }}</a><a href="#quickstart">{{ t('quickStart') }}</a><a href="/models">{{ t('marketplace') }}</a><ThemeCustomizer :locale="locale" /><select v-model="locale" class="language-select" :aria-label="t('switchLanguage')"><option value="zh-CN">{{ t('chinese') }}</option><option value="en-US">{{ t('english') }}</option></select><button class="button ghost" @click="openConsoleOrAuth">{{ t('console') }} <ChevronRight :size="15" /></button></div>
    </nav>
    <section class="hero-section">
      <div class="hero-copy"><p class="eyebrow">OPENAI-COMPATIBLE MODEL GATEWAY</p><h1>{{ t('heroTagline') }}</h1><p class="hero-description">{{ t('heroDescription') }}</p><div class="hero-actions"><button class="button primary hero-button" @click="openConsoleOrAuth">{{ t('openConsole') }} <ChevronRight :size="16" /></button><a class="text-link" href="#quickstart">{{ t('viewRequestExample') }}</a></div></div>
      <div class="hero-visual"><div class="visual-glow"></div><div class="route-card"><div class="route-card-top"><span><i class="live-dot"></i>ROUTER ONLINE</span><code>POST /v1/chat/completions</code></div><div class="route-model"><Bot :size="18" /><strong>kimi-k3</strong><span>{{ t('routingActive') }}</span></div><div class="route-line"><span class="route-node active"></span><div><b>{{ t('openaiMainRoute') }}</b><small>{{ t('priorityLabel') }} P10 &middot; 42ms</small></div><span class="route-check">✓</span></div><div class="route-line muted-route"><span class="route-node"></span><div><b>{{ t('backupChannelText') }}</b><small>{{ t('waitingForFailover') }}</small></div></div><div class="route-footer"><span>{{ t('successRate') }}</span><strong>99.98%</strong><span class="route-divider"></span><span>{{ t('avgLatency') }}</span><strong>186ms</strong></div></div></div>
    </section>
    <section id="features" class="feature-section"><div class="section-intro"><p class="eyebrow">BUILT FOR CONTROL</p><h2>{{ t('featuresTagline') }}</h2></div><div class="feature-grid"><article><span class="feature-number">01</span><RadioTower :size="21" /><h3>{{ t('smartRouting') }}</h3><p>{{ t('smartRoutingDesc') }}</p></article><article><span class="feature-number">02</span><Activity :size="21" /><h3>{{ t('fullObservability') }}</h3><p>{{ t('fullObservabilityDesc') }}</p></article><article><span class="feature-number">03</span><WalletCards :size="21" /><h3>{{ t('usageAndCost') }}</h3><p>{{ t('usageAndCostDesc') }}</p></article></div></section>
    <section id="quickstart" class="quickstart-section"><div><p class="eyebrow">ONE ENDPOINT</p><h2>{{ t('quickstartTagline') }}</h2><p>{{ t('quickstartDesc') }}</p></div><pre><span class="code-comment">// Using OpenAI SDK</span><span><b>const</b> client = <b>new</b> OpenAI({</span><span>  apiKey: <i>'sk-xh-your-key'</i>,</span><span>  baseURL: <i>'http://localhost:8080/v1'</i></span><span>})</span><span class="code-gap"></span><span><b>await</b> client.chat.completions.create({</span><span>  model: <i>'kimi-k3'</i>,</span><span>  messages: [{ role: <i>'user'</i>, content: <i>'你好'</i> }]</span><span>})</span></pre></section>
     <footer class="landing-footer"><span>© 2026 Xinghai Router</span><span>{{ t('footerSlogan') }}</span><span class="legal-links"><a href="/terms">{{ t('termsShort') }}</a><a href="/privacy">{{ t('privacyShort') }}</a></span></footer>
  </main>

  <main v-else-if="isMarketplacePage" class="public-marketplace">
       <nav class="marketplace-nav"><a class="landing-logo" href="/"><span class="brand-mark small"><Bot :size="19" /></span><span>{{ siteSettings.name }}</span></a><div><ThemeCustomizer :locale="locale" /><select v-model="locale" class="language-select" :aria-label="t('switchLanguage')"><option value="zh-CN">{{ t('chinese') }}</option><option value="en-US">{{ t('english') }}</option></select><button class="button ghost" @click="openConsoleOrAuth">{{ authenticated ? t('console') : t('login') }} <ChevronRight :size="15" /></button></div></nav>
    <section class="marketplace-page-content">
      <section class="marketplace-hero"><div><span class="marketplace-kicker"><Sparkles :size="13" /> MODEL CATALOG</span><h1>{{ t('marketplaceTitle') }}</h1><p>{{ t('marketplaceDesc') }}</p></div><div class="marketplace-count"><strong>{{ catalog.length }}</strong><span>{{ t('modelsAvailable') }}</span></div></section>
      <section class="marketplace-tools"><div class="marketplace-search"><Search :size="16" /><input v-model="catalogSearch" :aria-label="t('searchModelPlaceholder')" :placeholder="t('searchModelPlaceholder')" /></div><div class="group-filters"><button :class="{ active: catalogGroup === 'all' }" @click="catalogGroup = 'all'">{{ t('allGroups') }}</button><button v-for="group in catalogGroups" :key="group.id" :class="{ active: catalogGroup === group.id }" @click="catalogGroup = group.id">{{ group.name }} <small>{{ Number(group.multiplier).toFixed(2) }}x</small></button></div></section>
      <div class="model-market-grid"><article v-for="item in filteredCatalog" :key="item.model" class="model-market-card"><div class="model-card-heading"><span class="model-avatar">{{ item.model.slice(0, 1).toUpperCase() }}</span><div><h3>{{ item.model }}</h3><p>{{ modelProvider(item.model) }}</p></div><span :class="['pricing-state', { missing: item.input_per_million == null }]">{{ item.input_per_million == null ? t('pendingPrice') : t('available') }}</span></div><div class="model-price-grid"><div><span>{{ t('inputLabel') }}</span><strong>{{ formatPrice(item.input_per_million, actualMultiplier(item)) }}</strong><small>/ 1M tokens</small></div><div><span>{{ t('cachedInputLabel') }}</span><strong>{{ formatPrice(item.cached_input_per_million, actualMultiplier(item)) }}</strong><small>/ 1M tokens</small></div><div><span>{{ t('outputLabel') }}</span><strong>{{ formatPrice(item.output_per_million, actualMultiplier(item)) }}</strong><small>/ 1M tokens</small></div></div><footer><div class="model-groups"><span v-for="group in item.groups" :key="group.id" :class="{ selected: catalogGroup === group.id }">{{ group.name }}</span></div><span class="actual-rate">{{ t('actualRate') }} <b>{{ actualMultiplier(item).toFixed(2) }}x</b></span></footer></article><Empty v-if="!filteredCatalog.length" :text="catalog.length ? t('noMatchingModels') : error ? t('catalogUnavailable') : t('enableChannelsToShow')" /></div>
      <p class="pricing-note">{{ t('pricingNote') }}</p>
    </section>
  </main>

      <main v-else-if="isAuthPage" class="login-shell">
       <select v-model="locale" class="language-select login-language-select" :aria-label="t('switchLanguage')"><option value="zh-CN">{{ t('chinese') }}</option><option value="en-US">{{ t('english') }}</option></select>
      <section class="login-card"><div class="login-card-actions"><ThemeCustomizer :locale="locale" /><button class="language-toggle" :aria-label="t('switchLanguage')" @click="toggleLocale">{{ locale === 'zh-CN' ? 'EN' : '中' }}</button><button class="login-close" :aria-label="t('backHome')" @click="closeAuth">&times;</button></div><div class="brand-mark"><Bot :size="29" /></div><p class="eyebrow">XINGHAI ROUTER</p><h1>{{ t('loginTagline') }}</h1><div class="auth-tabs"><button :class="{ active: loginMode === 'login' }" @click="loginMode = 'login'">{{ t('signInTab') }}</button><button :class="{ active: loginMode === 'register' }" @click="loginMode = 'register'">{{ t('createAccountTab') }}</button></div><form @submit.prevent="accountSignIn(loginMode === 'register')"><label v-if="loginMode === 'register'">{{ t('nameLabel') }}<input v-model="accountForm.name" autocomplete="name" required maxlength="100" :placeholder="t('namePlaceholder')" /></label><label>{{ t('emailLabel') }}<input v-model="accountForm.email" type="email" autocomplete="email" required placeholder="name@example.com" /></label><label>{{ t('passwordLabel') }}<input v-model="accountForm.password" type="password" :autocomplete="loginMode === 'register' ? 'new-password' : 'current-password'" required minlength="8" :placeholder="t('passwordMinLength')" /></label><button class="button primary full" :disabled="busy">{{ loginMode === 'register' ? t('createAndOpenConsole') : t('signInConsole') }} <ChevronRight :size="16" /></button></form><p class="auth-legal">{{ t('agreeText') }} <a href="/terms">{{ t('termsShort') }}</a> {{ t('andConnector') }} <a href="/privacy">{{ t('privacyShort') }}</a>.</p></section>
  </main>

  <main v-else class="app-shell" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <aside class="sidebar">
        <div class="logo"><span class="brand-mark small"><Bot :size="19" /></span><span>{{ siteSettings.name }}</span></div>
        <nav>
           <div class="nav-group"><p class="nav-label">{{ t('general') }}</p><button v-for="[id, label, Icon] in generalNav" :key="id" :class="{ active: id === 'usage-overview' ? (route.query.view === id || route.params.view === id) : id === 'usage' ? view === id && route.query.view !== 'usage-overview' && route.params.view !== 'usage-overview' : view === id }" @click="openConsole(id)"><component :is="Icon" :size="17" /><span>{{ label }}</span></button></div>
          <div class="nav-group"><p class="nav-label">{{ t('billing') }}</p><button v-for="[id, label, Icon] in billingNav" :key="id" :class="{ active: view === id }" @click="openConsole(id)"><component :is="Icon" :size="17" /><span>{{ label }}</span></button></div>
          <div class="nav-group"><p class="nav-label">{{ t('personal') }}</p><button v-for="[id, label, Icon] in personalNav" :key="id" :class="{ active: view === id }" @click="openConsole(id)"><component :is="Icon" :size="17" /><span>{{ label }}</span></button></div>
          <div v-if="managementNav.length" class="nav-group management-group"><p class="nav-label">{{ t('management') }}</p><button v-for="[id, label, Icon] in managementNav" :key="id" :class="{ active: view === id }" @click="openConsole(id)"><component :is="Icon" :size="17" /><span>{{ label }}</span></button></div>
        </nav>
      <div class="sidebar-footer"><div class="gateway-status"><span class="live-dot"></span><span><b>{{ t('gatewayOnline') }}</b><small>{{ t('serviceRunning') }}</small></span></div><div class="sidebar-account"><i>{{ account?.name?.slice(0, 1) || '?' }}</i><span><b>{{ account?.name || t('loadingLabel') }}</b><small>{{ account?.role || t('accountLabel') }}</small></span><button :aria-label="t('actionAccountLogout')" :title="t('actionAccountLogout')" @click="signOut"><LogOut :size="16" /></button></div></div>
    </aside>
      <section class="content" :data-usage-page="route.query.view === 'usage-overview' || route.params.view === 'usage-overview' ? 'overview' : 'logs'">
        <header class="console-header"><div><p class="eyebrow">{{ managementNav.some((item) => item[0] === view) ? t('management') : personalNav.some((item) => item[0] === view) ? t('personal') : billingNav.some((item) => item[0] === view) ? t('billing') : t('general') }}</p><h1>{{ [...localizedManagementNavItems, ...generalNav, ...billingNav, ...personalNav, ...localizedAdminExtraNav].find((item) => item[0] === view)?.[1] }}</h1></div><div class="header-actions"><button class="theme-toggle sidebar-toggle" :aria-label="sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'" :title="sidebarCollapsed ? '展开侧边栏' : '收起侧边栏'" @click="sidebarCollapsed = !sidebarCollapsed"><PanelLeftOpen v-if="sidebarCollapsed" :size="16" /><PanelLeftClose v-else :size="16" /></button><a class="button ghost marketplace-link" href="/models"><Sparkles :size="15" />{{ t('marketplace') }}</a><span class="account-chip"><i>{{ account?.name?.slice(0, 1) || '?' }}</i>{{ account?.name || t('loadingLabel') }}</span><ThemeCustomizer :locale="locale" /><select v-model="locale" class="language-select" :aria-label="t('switchLanguage')"><option value="zh-CN">{{ t('chinese') }}</option><option value="en-US">{{ t('english') }}</option></select><button class="button ghost" @click="load" :disabled="busy"><RefreshCw :size="16" :class="{ spinning: busy }" />{{ t('refresh') }}</button></div></header>
        <template v-if="view === 'providers'"><section class="toolbar"><div><h2>{{ t('modelProviders') }}</h2><p>{{ t('providersDesc') }}</p></div><button v-if="can('system.manage')" class="button primary" @click="showProvider = true"><Plus :size="16" />{{ t('addProvider') }}</button></section><section class="panel table-panel"><table><thead><tr><th>{{ t('supplier') }}</th><th>{{ t('modelPrefix') }}</th><th>{{ t('iconSlug') }}</th><th>{{ t('matchPriority') }}</th></tr></thead><tbody><tr v-for="provider in providers" :key="provider.id"><td><b>{{ provider.name }}</b></td><td><span v-for="prefix in provider.prefixes" :key="prefix" class="pill">{{ prefix }}</span></td><td><code>{{ provider.slug }}</code></td><td>{{ provider.priority }}</td></tr></tbody></table><Empty v-if="!providers.length" :text="t('noProvidersYet')" /></section></template>
       <template v-if="view === 'overview'">
        <button class="setup-toggle" @click="setupCollapsed = !setupCollapsed"><ChevronDown :size="16" :class="{ rotated: !setupCollapsed }" /><span>{{ t('quickStartHeading') }} / {{ t('firstApiRequest') }}</span><span class="setup-progress">{{ setupProgress }} / 3</span></button>
        <section v-show="!setupCollapsed" class="setup-workspace">
           <div class="setup-guide">
              <div class="setup-heading"><div><span class="overview-kicker">{{ t('quickStartHeading') }}</span><h2>{{ account?.name || '' }}{{ t('sendModelRequest') }}</h2><p>{{ t('completeThreeSteps') }}</p></div></div>
             <div class="setup-steps">
               <button :class="{ complete: accountKeys.some((item) => !item.revoked_at) }" @click="openConsole('account')"><i><Check v-if="accountKeys.some((item) => !item.revoked_at)" :size="14" /><span v-else>1</span></i><span><b>{{ t('createApiKey') }}</b><small>{{ t('issueAccessCredentials') }}</small></span><ChevronRight :size="16" /></button>
               <button :class="{ complete: Number(account?.balance ?? 0) > 0 }" @click="openConsole('wallet')"><i><Check v-if="Number(account?.balance ?? 0) > 0" :size="14" /><span v-else>2</span></i><span><b>{{ t('checkAccountBalance') }}</b><small>{{ t('prepareAvailableQuota') }}</small></span><ChevronRight :size="16" /></button>
               <button :class="{ complete: personalRequests > 0 }" @click="openConsole('usage')"><i><Check v-if="personalRequests > 0" :size="14" /><span v-else>3</span></i><span><b>{{ t('sendModelRequest') }}</b><small>{{ t('confirmCallResults') }}</small></span><ChevronRight :size="16" /></button>
             </div>
           </div>
           <div class="request-preview"><div class="request-preview-title"><span><TerminalSquare :size="16" />{{ t('firstApiRequest') }}</span><code>curl</code></div><pre><span>curl {{ apiEndpoint }} \</span><span>  -H <i>"Authorization: Bearer sk-xh-..."</i> \</span><span>  -H <i>"Content-Type: application/json"</i> \</span><span>  -d <i>'{"model":"kimi-m3",</i></span><span><i>       "messages":[{"role":"user","content":"你好"}]}'</i></span></pre><div class="request-signals"><span><i></i>{{ t('gatewayOnline') }}</span><button @click="openConsole('account')">{{ t('viewKeys') }} <ChevronRight :size="14" /></button></div></div>
         </section>
        <div class="metrics"><article><span>{{ t('availableBalance') }}</span><strong>{{ Number(account?.balance ?? 0).toFixed(4) }}</strong><p><WalletCards :size="15" />{{ t('exclusiveOfReserved') }}</p></article><article><span>{{ t('validApiKeys') }}</span><strong>{{ accountKeys.filter((item) => !item.revoked_at).length }}</strong><p><KeyRound :size="15" />{{ t('currentAccountKeys') }}</p></article><article><span>{{ t('recentCalls') }}</span><strong>{{ personalRequests }}</strong><p><Activity :size="15" />{{ t('recent100Records') }}</p></article><article><span>{{ t('meteredTokens') }}</span><strong>{{ personalTokens.toLocaleString() }}</strong><p><ReceiptText :size="15" />{{ t('cumulativeCost') }} {{ personalCost.toFixed(6) }}</p></article></div>
        <div class="grid-two"><section class="panel usage-line-chart"><div class="panel-title"><div><h2>{{ t('usageTrend') }}</h2><p>{{ t('last7DaysTokenUsage') }}</p></div><button class="text-button" @click="openConsole('usage')">{{ t('viewAll') }}</button></div><div class="line-plot"><svg viewBox="0 0 100 100" preserveAspectRatio="none" :aria-label="t('last7DaysTokens')"><defs><linearGradient id="usageFill" x1="0" x2="0" y1="0" y2="1"><stop offset="0%" stop-color="#65a986" stop-opacity=".34" /><stop offset="100%" stop-color="#65a986" stop-opacity="0" /></linearGradient></defs><path :d="`M 0,100 L ${usageLinePoints} L 100,100 Z`" fill="url(#usageFill)" /><polyline :points="usageLinePoints" fill="none" stroke="#2d7657" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" vector-effect="non-scaling-stroke" /></svg><div class="line-labels"><span v-for="day in usageChart" :key="day.key">{{ day.label }}<b>{{ day.tokens ? day.tokens.toLocaleString() : '-' }}</b></span></div></div></section><section class="panel"><div class="panel-title"><div><h2>{{ t('accessKeys') }}</h2><p>{{ t('currentAccountKeysDesc') }}</p></div><button class="text-button" @click="openConsole('account')">{{ t('myAccount') }}</button></div><div v-if="accountKeys.length" class="compact-list"><div v-for="key in accountKeys.slice(0, 5)" :key="key.id"><code>{{ key.key_prefix }}...</code><span>{{ key.name }}</span><b :class="key.revoked_at ? 'danger' : 'success'">{{ key.revoked_at ? t('revoked') : t('valid') }}</b></div></div><Empty v-else :text="t('noApiKeysYet')" /></section></div>
      </template>

       <template v-if="view === 'users'"><section class="toolbar"><div><h2>{{ t('usersAndPermissions') }}</h2><p>{{ t('usersPermissionDesc') }}</p></div></section><section class="panel table-panel"><table><thead><tr><th>{{ t('userLabel') }}</th><th>{{ t('roleLabel') }}</th><th>{{ t('groupLabel') }}</th><th>{{ t('balanceLabel') }}</th><th>{{ t('permissionScope') }}</th><th>{{ t('accountStatus') }}</th><th></th></tr></thead><tbody><tr v-for="user in users" :key="user.id"><td><b>{{ user.name }}</b><small>{{ user.email }}</small></td><td><span class="pill">{{ user.role }}</span></td><td>{{ user.groups.length || t('none') }}</td><td>{{ Number(user.balance ?? 0).toFixed(4) }}<small v-if="Number(user.reserved ?? 0)">{{ t('reserved') }} {{ Number(user.reserved).toFixed(4) }}</small></td><td>{{ user.role === 'admin' ? t('allPermissions') : user.permissions.join(', ') || t('none') }}</td><td><span :class="['state', user.enabled ? 'good' : 'bad']">{{ user.enabled ? t('enabled') : t('disabled') }}</span></td><td><button v-if="can('system.manage')" class="text-button" @click="manageUser(user)">{{ t('edit') }}</button></td></tr></tbody></table><Empty v-if="!users.length" :text="t('noUsersYet')" /></section></template>
        <template v-if="view === 'groups'"><section class="toolbar"><div><h2>{{ t('accessGroups') }}</h2><p>{{ t('groupsDesc') }}</p></div></section><div class="group-page"><form class="panel group-import-form" @submit.prevent="importGroups"><div><h3>{{ t('batchImportGroups') }}</h3><p>{{ t('importGroupsDesc') }}</p></div><textarea v-model="groupImportText" required :placeholder="t('importJsonPlaceholder')"></textarea><button class="button primary" :disabled="busy"><Plus :size="16" />{{ t('importOneClick') }}</button></form><form class="panel group-create-form" @submit.prevent="createGroup"><div><h3>{{ t('createGroupTitle') }}</h3><p>{{ t('createGroupDesc') }}</p></div><label>{{ t('groupNameLabel') }}<input v-model="groupForm.name" required maxlength="100" :placeholder="t('groupNamePlaceholder')" /></label><label>{{ t('multiplierLabel') }}<input v-model.number="groupForm.multiplier" required type="number" min="0" step="0.01" /></label><button class="button primary" :disabled="busy"><Plus :size="16" />{{ t('createGroupButton') }}</button></form><section class="panel table-panel"><table><thead><tr><th>{{ t('groupNameLabel') }}</th><th>{{ t('createdAt') }}</th><th>{{ t('multiplierLabel') }}</th><th></th></tr></thead><tbody><tr v-for="group in groups" :key="group.id"><td><b>{{ group.name }}</b><small>{{ group.id }}</small></td><td>{{ formatDate(group.created_at) }}</td><td><form class="group-rate-form" @submit.prevent="editGroupMultiplier(group, $event)"><input name="multiplier" :value="Number(group.multiplier)" :aria-label="t('multiplierLabel')" required type="number" min="0" step="0.01" /><span>x</span><button class="button ghost" :disabled="busy">{{ t('saveLabel') }}</button></form></td><td></td></tr></tbody></table><Empty v-if="!groups.length" :text="t('noGroupsYet')" /></section></div></template>
      <template v-if="view === 'keys'"><section class="toolbar"><div><h2>{{ t('keys') }}</h2><p>{{ t('showOnceAfterCreation') }}</p></div><button class="button primary" :disabled="!users.length" @click="showKey = true"><Plus :size="16" />{{ t('createKeyButton') }}</button></section><section class="panel table-panel"><table><thead><tr><th>{{ t('keyName') }}</th><th>{{ t('userLabel') }}</th><th>{{ t('keyPrefix') }}</th><th>{{ t('lastUsed') }}</th><th>{{ t('accountStatus') }}</th><th></th></tr></thead><tbody><tr v-for="key in keys" :key="key.id"><td><b>{{ key.name }}</b></td><td>{{ userName(key.user_id) }}</td><td><code>{{ key.key_prefix }}...</code></td><td>{{ formatDate(key.last_used_at) }}</td><td><span :class="['state', key.revoked_at ? 'bad' : 'good']">{{ key.revoked_at ? t('revoked') : t('valid') }}</span></td><td><button v-if="!key.revoked_at" class="text-button danger" @click="revokeKey(key)">{{ t('revokeLabel') }}</button></td></tr></tbody></table><Empty v-if="!keys.length" :text="t('createUserThenIssueKey')" /></section></template>
        <template v-if="view === 'channels'"><section class="toolbar"><div><h2>{{ t('upstreamChannels') }}</h2><p>{{ t('channelsDesc') }}</p></div><button v-if="can('channels.manage')" class="button primary" @click="showChannel = true"><Plus :size="16" />{{ t('addChannel') }}</button></section><section class="panel table-panel channel-table-panel"><table><thead><tr><th>Status</th><th>Channel name</th><th>Upstream URL</th><th>{{ t('modelLabel') }}</th><th>{{ t('priorityLabel') }}</th><th v-if="can('channels.manage')">{{ t('enableChannelLabel') }}</th><th v-if="can('channels.manage')"></th></tr></thead><tbody><tr v-for="channel in channels" :key="channel.id"><td><span :class="['state', channel.enabled ? 'good' : 'bad']"><i :class="['status-dot', { off: !channel.enabled }]"></i>{{ channel.enabled ? t('enabled') : channel.auto_disabled ? t('autoDisabled') : t('disabled') }}</span><small v-if="channel.auto_disabled && channel.disabled_reason" class="muted" :title="channel.disabled_reason">{{ channel.disabled_reason }}</small></td><td><b>{{ channel.name }}</b><small>{{ channel.provider }}</small></td><td><code :title="channel.base_url">{{ channel.base_url }}</code></td><td><div class="model-tags"><span v-for="model in channel.models" :key="model">{{ model }}</span></div></td><td>{{ channel.priority }}</td><td v-if="can('channels.manage')"><button class="toggle" :class="{ on: channel.enabled }" :aria-label="channel.enabled ? t('disableChannelLabel') : t('enableChannelLabel')" @click="toggleChannel(channel)"><i></i></button></td><td v-if="can('channels.manage')"><button class="text-button" @click="editChannel(channel)">{{ t('edit') }}</button></td></tr></tbody></table><Empty v-if="!channels.length" :text="t('addOpenAICompatibleUpstream')" /></section></template>
            <template v-if="view === 'account'"><section class="toolbar"><div><h2>{{ t('account') }}</h2><p>{{ t('keyBelongsToAccount') }}</p></div><button class="button primary" @click="showAccountKey = true"><Plus :size="16" />{{ t('createKeyButton') }}</button></section><section class="panel table-panel"><table><thead><tr><th>{{ t('keyName') }}</th><th>{{ t('keyPrefix') }}</th><th>{{ t('createdAt') }}</th><th>{{ t('lastUsed') }}</th><th>{{ t('accountStatus') }}</th><th></th></tr></thead><tbody><tr v-for="key in accountKeys" :key="key.id"><td><b>{{ key.name }}</b></td><td><code>{{ key.key_prefix }}...</code></td><td>{{ formatDate(key.created_at) }}</td><td>{{ formatDate(key.last_used_at) }}</td><td><span :class="['state', key.revoked_at ? 'bad' : 'good']">{{ key.revoked_at ? t('revoked') : t('valid') }}</span></td><td><button v-if="!key.revoked_at" class="text-button" @click="editAccountKey(key)">{{ t('edit') }}</button></td></tr></tbody></table><Empty v-if="!accountKeys.length" :text="t('noApiKeysYet')" /></section></template>
            <template v-if="view === 'profile'"><section class="profile-layout"><section class="panel account-card"><div class="profile-avatar">{{ account?.name?.slice(0, 1) || '?' }}</div><div><span class="overview-kicker">{{ t('accountProfile') }}</span><h2>{{ account?.name }}</h2><p>{{ account?.email }}</p></div></section><section class="panel profile-details"><div><span>{{ t('accountRole') }}</span><b>{{ account?.role }}</b></div><div><span>{{ t('userGroups') }}</span><div v-if="ownGroups.length" class="profile-group-tags"><span v-for="group in ownGroups" :key="group">{{ group }}</span></div><b v-else>{{ t('notInAnyGroup') }}</b></div><div><span>{{ t('accountId') }}</span><code>{{ account?.id }}</code></div><div><span>{{ t('permissionScope') }}</span><b>{{ account?.role === 'admin' ? t('allPermissions') : account?.permissions.join(', ') || t('regularAccount') }}</b></div></section><section class="panel avatar-settings"><div><h2>{{ t('avatarSection') }}</h2><p>{{ t('avatarDesc') }}</p></div><div class="avatar-current"><img v-if="account?.avatar_url" class="profile-avatar" :src="account.avatar_url" :alt="t('avatarSection')" /><span v-else class="profile-avatar">{{ account?.name?.slice(0, 1) || '?' }}</span></div><div class="avatar-actions"><form class="avatar-url-form" @submit.prevent="saveAvatarUrl"><input v-model="avatarUrlInput" type="url" :placeholder="t('avatarUrlPlaceholder')" /><button class="button ghost" type="submit">{{ t('setAvatar') }}</button></form><div class="avatar-upload-row"><input ref="avatarInput" class="visually-hidden" type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="chooseAvatar" /><button class="button ghost" type="button" @click="avatarInput?.click()">{{ t('uploadAvatar') }}</button><button v-if="account?.avatar_url" class="text-button danger" type="button" @click="removeAvatar">{{ t('remove') }}</button></div></div></section></section></template>
         <template v-if="view === 'wallet'"><section class="wallet-hero"><div><span>{{ t('availableBalance') }}</span><strong>{{ Number(account?.balance ?? 0).toFixed(4) }}</strong><p>{{ t('balanceForModelCalls') }}</p></div><WalletCards :size="64" /></section><div class="metrics wallet-metrics"><article><span>{{ t('currentBalance') }}</span><strong>{{ Number(account?.balance ?? 0).toFixed(4) }}</strong><p><WalletCards :size="15" />{{ t('accountAvailableQuota') }}</p></article><article><span>{{ t('reservedAmount') }}</span><strong>{{ Number(account?.reserved ?? 0).toFixed(4) }}</strong><p>{{ t('reservedForConcurrent') }}</p></article><article><span>{{ t('cumulativeSpending') }}</span><strong>{{ personalCost.toFixed(6) }}</strong><p><ReceiptText :size="15" />{{ t('recent100Records') }}</p></article></div><form class="panel payment-form" @submit.prevent="createPayment"><div><h2>{{ t('onlineTopup') }}</h2><p>{{ paymentsEnabled ? t('onlineTopupDesc') : t('paymentNotConfigured') }}</p><strong v-if="paymentMessage" class="payment-message">{{ paymentMessage }}</strong></div><label>{{ t('topupAmount') }}<input v-model.number="paymentForm.amount" type="number" min="1" max="100000" step="0.01" required /></label><label>{{ t('paymentMethod') }}<select v-model="paymentForm.type" required><option v-for="method in paymentMethods" :key="method.id" :value="method.code">{{ method.name }}</option></select></label><button class="button primary" :disabled="busy || !paymentsEnabled || !paymentForm.type">{{ t('goToPay') }}</button></form><section v-if="payments.length" class="panel table-panel payment-orders"><div class="panel-title"><div><h2>{{ t('topupOrders') }}</h2><p>{{ t('topupOrdersDesc') }}</p></div></div><table><thead><tr><th>{{ t('createdAt') }}</th><th>{{ t('orderNumber') }}</th><th>{{ t('paymentMethod') }}</th><th>{{ t('topupAmount') }}</th><th>{{ t('accountStatus') }}</th></tr></thead><tbody><tr v-for="item in payments.slice(0, 10)" :key="item.order_no"><td>{{ formatDate(item.created_at) }}</td><td><code>{{ item.order_no }}</code></td><td>{{ paymentMethods.find((method) => method.code === item.payment_type)?.name ?? item.payment_type }}</td><td>{{ item.amount }}</td><td><span :class="['state', item.status === 'paid' ? 'good' : 'bad']">{{ item.status }}</span></td></tr></tbody></table></section><section class="panel table-panel"><div class="panel-title"><div><h2>{{ t('balanceLedger') }}</h2><p>{{ t('ledgerDesc') }}</p></div><button class="text-button" @click="openConsole('ledger')">{{ t('viewAll') }}</button></div><table><thead><tr><th>{{ t('createdAt') }}</th><th>{{ t('typeLabel') }}</th><th>Change</th><th>{{ t('balanceLabel') }}</th><th>Note</th></tr></thead><tbody><tr v-for="item in ledger.slice(0, 10)" :key="item.id"><td>{{ formatDate(item.created_at) }}</td><td>{{ item.kind }}</td><td :class="item.amount < 0 ? 'danger' : 'success'">{{ item.amount }}</td><td>{{ item.balance_after }}</td><td>{{ item.note || item.request_id }}</td></tr></tbody></table><Empty v-if="!ledger.length" :text="t('noLedgerEntries')" /></section></template>
         <template v-if="view === 'usage'"><section class="usage-summary"><article><span>{{ t('last7DaysTokens') }}</span><strong>{{ personalTokens.toLocaleString() }}</strong><small>{{ t('inputOutputTotal') }}</small></article><article><span>{{ t('last7DaysCost') }}</span><strong>{{ personalCost.toFixed(6) }}</strong><small>{{ t('basedOnCurrentPricing') }}</small></article><article><span>{{ t('callCount') }}</span><strong>{{ personalRequests }}</strong><small>{{ t('recent100UsageRecords') }}</small></article></section><section class="panel usage-chart"><div class="panel-title"><div><h2>{{ t('usageTrend') }}</h2><p>{{ t('last7DaysTokenAndCost') }}</p></div><div class="chart-legend"><span><i class="token-dot"></i>{{ t('tokenLabel') }}</span><span><i class="cost-dot"></i>{{ t('costLabel') }}</span></div></div><div class="chart-bars"><div v-for="day in usageChart" :key="day.key" class="chart-day"><div class="chart-values"><span :style="{ height: `${day.tokenHeight}%` }" :title="`${day.tokens.toLocaleString()} tokens`"></span><i :style="{ height: `${day.costHeight}%` }" :title="`${t('costLabel')} ${day.cost.toFixed(6)}`"></i></div><b>{{ day.label }}</b><small>{{ day.tokens ? day.tokens.toLocaleString() : '-' }}</small></div></div></section><form class="panel activity-filters" @submit.prevent="loadActivity(true)"><label v-if="can('users.read')">{{ t('userLabel') }}<select v-model="activityFilters.user_id"><option value="">{{ t('allUsers') }}</option><option v-for="user in users" :key="user.id" :value="user.id">{{ user.name }} · {{ user.email }}</option></select></label><label>{{ t('modelLabel') }}<select v-model="activityFilters.model"><option value="">{{ t('allModels') }}</option><option v-for="model in activityModels" :key="model" :value="model">{{ model }}</option></select></label><label>{{ t('groupLabel') }}<select v-model="activityFilters.group_id"><option value="">{{ t('allGroups') }}</option><option v-for="group in groups" :key="group.id" :value="group.id">{{ group.name }}</option></select></label><label>{{ t('typeLabel') }}<select v-model="activityFilters.type"><option value="">{{ t('allTypes') }}</option><option value="request">{{ activityTypeLabel['request'] }}</option><option value="login">{{ activityTypeLabel['login'] }}</option><option value="register">{{ activityTypeLabel['register'] }}</option><option value="logout">{{ activityTypeLabel['logout'] }}</option><option value="topup">{{ activityTypeLabel['topup'] }}</option><option value="operation">{{ activityTypeLabel['operation'] }}</option></select></label><label>{{ t('startTime') }}<input v-model="activityFilters.start" type="datetime-local" /></label><label>{{ t('endTime') }}<input v-model="activityFilters.end" type="datetime-local" /></label><div class="activity-filter-actions"><button class="button primary" :disabled="busy">{{ t('filterLabel') }}</button><button type="button" class="button ghost" :disabled="busy" @click="resetActivityFilters">{{ t('resetFiltersLabel') }}</button></div></form><section class="panel table-panel"><div class="panel-title"><div><h2>{{ t('usageLogs') }}</h2><p>{{ t('usageLogsDesc') }}</p></div></div><table><thead><tr><th>{{ t('createdAt') }}</th><th>{{ t('typeLabel') }}</th><th>{{ t('userLabel') }}</th><th>{{ t('modelLabel') }} / Action</th><th>{{ t('groupLabel') }}</th><th>Status / Duration</th><th>Usage / Details</th></tr></thead><tbody><tr v-for="item in activityLogs" :key="`${item.type}-${item.id}`"><td>{{ formatDate(item.created_at) }}</td><td><span class="pill">{{ activityTypeLabel[item.type] }}</span></td><td>{{ item.user_name }}</td><td><code v-if="item.model">{{ item.model }}</code><span v-else>{{ actionLabel(item) }}</span></td><td>{{ item.group_name || '-' }}</td><td><span v-if="item.status_code != null" :class="['state', item.status_code < 400 ? 'good' : 'bad']">{{ item.status_code }}</span><small v-if="item.duration_ms != null">{{ item.duration_ms }} ms</small><span v-if="item.status_code == null">{{ t('success') }}</span></td><td><code>{{ activityDetail(item) }}</code></td></tr></tbody></table><Empty v-if="!activityLogs.length" :text="t('noMatchingLogs')" /></section></template>
         <template v-if="view === 'ledger'"><section class="toolbar"><div><h2>{{ t('balanceLedger') }}</h2><p>{{ t('ledgerDesc') }}</p></div></section><section class="panel table-panel"><table><thead><tr><th>{{ t('createdAt') }}</th><th>{{ t('typeLabel') }}</th><th>Change</th><th>{{ t('balanceLabel') }}</th><th>Note</th></tr></thead><tbody><tr v-for="item in ledger" :key="item.id"><td>{{ formatDate(item.created_at) }}</td><td><span class="pill">{{ item.kind }}</span></td><td :class="item.amount < 0 ? 'danger' : 'success'">{{ item.amount }}</td><td>{{ item.balance_after }}</td><td>{{ item.note || item.request_id }}</td></tr></tbody></table><Empty v-if="!ledger.length" :text="t('noLedgerEntries')" /></section></template>
           <template v-if="view === 'site-settings'"><section class="toolbar"><div><h2>{{ t('siteSettings') }}</h2><p>{{ t('siteSettingsDesc') }}</p></div></section><form class="panel pricing-form" @submit.prevent="saveSiteSettings"><label>{{ t('siteName') }}<input v-model="siteSettingsForm.name" required maxlength="100" /></label><label>{{ t('siteIconUrl') }}<input v-model="siteSettingsForm.icon_url" type="url" placeholder="https://example.com/favicon.png" /></label><label><input v-model="siteSettingsForm.auto_disable_failed_channels" type="checkbox" /> Auto-disable after 3 consecutive failures</label><button class="button primary" :disabled="busy">{{ t('saveSettings') }}</button></form></template>
           <template v-if="view === 'reliability'"><section class="toolbar"><div><h2>{{ t('reliability') }}</h2><p>{{ t('reliabilityDesc') }}</p></div></section><form class="panel pricing-form reliability-form" @submit.prevent="saveReliabilitySettings"><h3 class="reliability-section-title">{{ t('requestRetry') }}</h3><label>{{ t('retryCount') }}<input v-model.number="reliabilityForm.retry_count" type="number" min="0" max="10" step="1" required /><small>{{ t('retryCountHint') }}</small></label><label>{{ t('retryStatusCodes') }}<input v-model="reliabilityForm.retry_status_codes" required placeholder="100-199,300-407,409-503,505-523,525-599" /><small>{{ t('statusCodesHint') }}</small></label><h3 class="reliability-section-title">{{ t('channelHealthCheck') }}</h3><label>{{ t('healthCheckMode') }}<select v-model="reliabilityForm.health_check_mode"><option value="off">{{ t('healthCheckOff') }}</option><option value="scheduled_all">{{ t('healthCheckScheduledAll') }}</option><option value="passive_recovery">{{ t('healthCheckPassiveRecovery') }}</option></select><small>{{ t('healthCheckModeHint') }}</small></label><label>{{ t('healthCheckInterval') }}<input v-model.number="reliabilityForm.health_check_interval_minutes" type="number" min="1" max="1440" step="1" required /><small>{{ t('healthCheckIntervalHint') }}</small></label><label class="payment-enabled"><input v-model="reliabilityForm.health_check_auto_recover" type="checkbox" />{{ t('healthCheckAutoRecover') }}</label><label>{{ t('healthCheckChannelIds') }}<textarea v-model="reliabilityForm.health_check_channel_ids" rows="3" :placeholder="t('healthCheckChannelIdsPlaceholder')"></textarea><small>{{ t('healthCheckChannelIdsHint') }}</small></label><h3 class="reliability-section-title">{{ t('autoDisableRules') }}</h3><label class="payment-enabled"><input v-model="reliabilityForm.auto_disable_on_test_failure" type="checkbox" />{{ t('autoDisableOnTestFailure') }}</label><label>{{ t('autoDisableSlowSeconds') }}<input v-model.number="reliabilityForm.auto_disable_slow_seconds" type="number" min="0" max="600" step="1" /><small>{{ t('autoDisableSlowSecondsHint') }}</small></label><label>{{ t('autoDisableStatusCodes') }}<input v-model="reliabilityForm.auto_disable_status_codes" placeholder="401,429,503" /><small>{{ t('statusCodesHint') }}</small></label><label>{{ t('autoDisableKeywords') }}<textarea v-model="reliabilityForm.auto_disable_keywords" rows="8" :placeholder="t('autoDisableKeywordsPlaceholder')"></textarea><small>{{ t('autoDisableKeywordsHint') }}</small></label><button class="button primary" :disabled="busy">{{ t('saveSettings') }}</button></form></template>
          <template v-if="view === 'payment-settings'"><section class="toolbar"><div><h2>{{ t('paymentSettings') }}</h2><p>{{ t('paymentSettingsDesc') }}</p></div></section><form class="panel payment-settings-form" @submit.prevent="savePaymentSettings"><label class="payment-enabled"><input v-model="paymentSettingsForm.enabled" type="checkbox" />{{ t('enableOnlinePayment') }}</label><label>{{ t('epayBaseUrl') }}<input v-model="paymentSettingsForm.base_url" type="url" required placeholder="https://pay.example.com" /></label><label>{{ t('publicBaseUrl') }}<input v-model="paymentSettingsForm.public_base_url" type="url" required placeholder="https://router.example.com" /></label><label>{{ t('merchantId') }}<input v-model="paymentSettingsForm.merchant_id" required /></label><label>{{ t('merchantKey') }}<input v-model="paymentSettingsForm.merchant_key" type="password" :required="!paymentSettings.has_merchant_key" :placeholder="paymentSettings.has_merchant_key ? t('leaveBlankUnchanged') : t('requiredField')" autocomplete="new-password" /></label><button class="button primary" :disabled="busy">{{ t('saveSettings') }}</button></form><form class="panel payment-method-create" @submit.prevent="createPaymentMethod"><div><h3>{{ t('addPaymentMethod') }}</h3><p>{{ t('paymentMethodCodeHint') }}</p></div><label>{{ t('paymentMethodCode') }}<input v-model="paymentMethodForm.code" required maxlength="50" placeholder="alipay" /></label><label>{{ t('paymentMethodName') }}<input v-model="paymentMethodForm.name" required maxlength="100" placeholder="支付宝" /></label><label class="method-enabled"><input v-model="paymentMethodForm.enabled" type="checkbox" />{{ t('enabled') }}</label><button class="button primary" :disabled="busy">{{ t('addPaymentMethod') }}</button></form><section class="panel table-panel"><div class="panel-title"><div><h2>{{ t('paymentMethods') }}</h2><p>{{ t('paymentMethodsDesc') }}</p></div></div><table><thead><tr><th>{{ t('paymentMethodCode') }}</th><th>{{ t('paymentMethodName') }}</th><th>{{ t('accountStatus') }}</th><th></th></tr></thead><tbody><tr v-for="method in paymentSettings.methods" :key="method.id"><td><input v-model="method.code" maxlength="50" /></td><td><input v-model="method.name" maxlength="100" /></td><td><label class="method-enabled"><input v-model="method.enabled" type="checkbox" />{{ method.enabled ? t('enabled') : t('disabled') }}</label></td><td><button class="text-button" @click="updatePaymentMethod(method)">{{ t('saveLabel') }}</button><button class="text-button danger" @click="deletePaymentMethod(method)">{{ t('remove') }}</button></td></tr></tbody></table><Empty v-if="!paymentSettings.methods.length" :text="t('noPaymentMethods')" /></section></template>
        <template v-if="view === 'pricing'"><section class="toolbar"><div><h2>{{ t('modelPricing') }}</h2><p>{{ t('pricingDesc') }}</p></div></section><form v-if="can('pricing.manage')" class="panel pricing-form" @submit.prevent="syncNewAPIPricing"><label>{{ t('newapiUrl') }}<input v-model="newAPIPricingForm.base_url" required type="url" placeholder="https://newapi.example.com" /></label><label>{{ t('loginToken') }}<input v-model="newAPIPricingForm.api_key" type="password" :placeholder="t('optional')" /></label><label>{{ t('perQuotaPrice') }}<input v-model.number="newAPIPricingForm.price_per_quota_unit" required type="number" min="0" step="any" /><small>{{ t('newapiPricingHint') }}</small></label><button class="button primary" :disabled="busy">{{ t('syncFromNewapi') }}</button></form><form v-if="can('pricing.manage')" class="panel pricing-form" @submit.prevent="savePricing"><label>{{ t('modelLabel') }}<input v-model="pricingForm.model" required placeholder="e.g. kimi-k3" /></label><label>{{ t('inputPrice') }}<input v-model.number="pricingForm.input_per_million" type="number" min="0" step="any" placeholder="0" /></label><label>{{ t('cachedInput') }}<input v-model.number="pricingForm.cached_input_per_million" type="number" min="0" step="any" placeholder="0" /></label><label>{{ t('outputPrice') }}<input v-model.number="pricingForm.output_per_million" type="number" min="0" step="any" placeholder="0" /></label><label>{{ t('multiplierLabel') }}<input v-model.number="pricingForm.multiplier" type="number" min="0.01" step="any" placeholder="1" /></label><button class="button primary">{{ t('saveRule') }}</button></form><section class="panel table-panel"><table><thead><tr><th>{{ t('modelLabel') }}</th><th>{{ t('inputPrice') }}</th><th>{{ t('cachedInput') }}</th><th>{{ t('outputPrice') }}</th><th>{{ t('multiplierLabel') }}</th></tr></thead><tbody><tr v-for="item in pricing" :key="item.id"><td><code>{{ item.model }}</code></td><td>{{ item.input_per_million }}</td><td>{{ item.cached_input_per_million }}</td><td>{{ item.output_per_million }}</td><td>{{ item.multiplier }}</td></tr></tbody></table><Empty v-if="!pricing.length" :text="t('noPricingRules')" /></section></template>
      </section>


        <div v-if="selectedUser || showKey || showAccountKey || editingAccountKey || showChannel || editingChannel || showProvider || createdKey" class="modal-backdrop" @click.self="selectedUser = null; editingAccountKey = null; editingChannel = null; showKey = showAccountKey = showChannel = showProvider = false">
      <form v-if="selectedUser" class="modal" @submit.prevent="saveUserAccess"><div class="modal-title"><h2>{{ t('editUser') }}</h2><button type="button" @click="selectedUser = null; originalUser = null">×</button></div><p class="muted">{{ selectedUser.id }}</p><label>{{ t('nameLabel') }}<input v-model="selectedUser.name" required maxlength="100" /></label><label>{{ t('emailLabel') }}<input v-model="selectedUser.email" required type="email" /></label><label>{{ t('newPassword') }} <small>{{ t('leaveEmptyToKeep') }}</small><input v-model="userPassword" type="password" minlength="8" autocomplete="new-password" /></label><label>{{ t('accountStatus') }}<select v-model="selectedUser.enabled"><option :value="true">{{ t('enabled') }}</option><option :value="false">{{ t('disabled') }}</option></select></label><label>{{ t('roleLabel') }}<select v-model="selectedUser.role"><option value="user">{{ t('userRole') }}</option><option value="operator">{{ t('operatorRole') }}</option><option value="admin">{{ t('adminRoleFull') }}</option></select></label><label>{{ t('balanceLabel') }}<input v-model.number="userBalance" required type="number" min="0" step="0.00000001" /></label><label>{{ t('balanceChangeNote') }}<input v-model="userBalanceNote" maxlength="200" :placeholder="t('balanceNotePlaceholder')" /></label><label>{{ t('userGroupsLabel') }}<select v-model="selectedGroups" multiple size="5"><option v-for="group in groups" :key="group.id" :value="group.id">{{ group.name }} · {{ Number(group.multiplier).toFixed(2) }}x</option></select></label><label v-if="selectedUser.role !== 'admin'">{{ t('granularPermissions') }}<select v-model="selectedPermissions" multiple size="8"><option v-for="permission in permissions" :key="permission" :value="permission">{{ permission }}</option></select></label><button class="button primary full" :disabled="busy">{{ t('saveChanges') }}</button></form>
      <form v-if="showKey" class="modal" @submit.prevent="createKey"><div class="modal-title"><h2>{{ t('createApiKeyTitle') }}</h2><button type="button" @click="showKey = false">×</button></div><label>{{ t('userLabel') }}<select v-model="keyForm.user_id" required><option disabled value="">{{ t('selectUser') }}</option><option v-for="user in users" :key="user.id" :value="user.id">{{ user.name }} · {{ user.email }}</option></select></label><label>{{ t('useGroup') }}<select v-model="keyForm.group_id"><option value="">{{ t('autoMatch') }}</option><option v-for="group in groups.filter((item) => users.find((user) => user.id === keyForm.user_id)?.groups.includes(item.id))" :key="group.id" :value="group.id">{{ group.name }} · {{ Number(group.multiplier).toFixed(2) }}x</option></select></label><label>{{ t('keyName') }}<input v-model="keyForm.name" required :placeholder="t('keyNamePlaceholder')" /></label><label>{{ t('expiresAt') }} <small>{{ t('optional') }}</small><input v-model="keyForm.expires_at" type="datetime-local" /></label><button class="button primary full" :disabled="busy">{{ t('issueKey') }}</button></form>
        <form v-if="showAccountKey || editingAccountKey" class="modal" @submit.prevent="editingAccountKey ? updateAccountKey() : createAccountKey()"><div class="modal-title"><h2>{{ editingAccountKey ? t('editApiKey') : t('createApiKeyTitle') }}</h2><button type="button" @click="showAccountKey = false; editingAccountKey = null">×</button></div><p v-if="!editingAccountKey" class="muted">{{ t('keyBelongsToAccount') }}</p><label>{{ t('useGroup') }}<select v-model="accountKeyForm.group_id"><option value="">{{ t('autoMatch') }}</option><option v-for="group in groups.filter((item) => ownGroups.includes(item.name))" :key="group.id" :value="group.id">{{ group.name }} · {{ Number(group.multiplier).toFixed(2) }}x</option></select></label><label>{{ t('keyName') }}<input v-model="accountKeyForm.name" required maxlength="100" :placeholder="t('keyNamePlaceholder2')" /></label><label>{{ t('expiresAt') }} <small>{{ t('optional') }}</small><input v-model="accountKeyForm.expires_at" type="datetime-local" /></label><button class="button primary full" :disabled="busy">{{ editingAccountKey ? t('saveChanges') : t('createKeyButton') }}</button></form>
       <form v-if="showChannel || editingChannel" class="modal" @submit.prevent="editingChannel ? updateChannel() : createChannel()"><div class="modal-title"><h2>{{ editingChannel ? t('editChannel') : t('addChannel') }}</h2><button type="button" @click="showChannel = false; editingChannel = null">×</button></div><label>{{ t('channelNameLabel') }}<input v-model="channelForm.name" required maxlength="100" /></label><label>{{ t('upstreamProtocol') }}<select v-model="channelForm.provider"><option value="openai">OpenAI</option><option value="anthropic">Anthropic</option><option value="ollama">Ollama</option><option value="kimi">Kimi</option><option value="opencode_go">OpenCode Go</option></select></label><label>{{ t('upstreamURL') }}<input v-model="channelForm.base_url" required type="url" /></label><label>{{ t('apiKeyLabel') }} <small>{{ editingChannel ? t('leaveBlankUnchanged') : t('requiredField') }}</small><input v-model="channelForm.api_key" :required="!editingChannel" type="password" autocomplete="new-password" /></label><label>{{ t('modelLabel') }} <small>{{ t('modelsCommaSeparated') }}</small><input v-model="channelForm.models" required /></label><button class="text-button" type="button" :disabled="busy || !channelForm.api_key" @click="fetchChannelModels">{{ t('fetchUpstreamModels') }}</button><label>{{ t('priorityLabel') }}<input v-model.number="channelForm.priority" required type="number" min="0" /></label><label>{{ t('availableGroups') }}<select v-model="channelForm.groups" multiple size="5"><option v-for="group in groups" :key="group.id" :value="group.id">{{ group.name }} · {{ Number(group.multiplier).toFixed(2) }}x</option></select></label><button class="button primary full" :disabled="busy">{{ editingChannel ? t('saveChanges') : t('addChannel') }}</button></form>
       <form v-if="showProvider" class="modal" @submit.prevent="saveProvider"><div class="modal-title"><h2>{{ t('configureProvider') }}</h2><button type="button" @click="showProvider = false">×</button></div><p class="muted">{{ t('providerDesc') }}</p><label>{{ t('providerName') }}<input v-model="providerForm.name" required :placeholder="t('providerNamePlaceholder')" /></label><label>{{ t('iconSlug') }}<input v-model="providerForm.slug" required :placeholder="t('iconSlugPlaceholder')" /></label><label>{{ t('modelPrefix') }}<input v-model="providerForm.prefixes" required :placeholder="t('modelPrefixPlaceholder')" /></label><label>{{ t('matchPriority') }}<input v-model.number="providerForm.priority" required type="number" min="0" /></label><button class="button primary full" :disabled="busy">{{ t('saveProviderButton') }}</button></form>
      <section v-if="createdKey" class="modal secret"><div class="modal-title"><h2>{{ t('saveApiKey') }}</h2><button @click="createdKey = ''">×</button></div><p>{{ t('saveKeyWarning') }}</p><code>{{ createdKey }}</code><button class="button primary full" @click="copyKey"><Copy :size="16" />{{ t('copyKeyButton') }}</button><button class="button ghost full" @click="createdKey = ''">{{ t('iHaveSaved') }}</button></section>
    </div>
  </main>
</template>
