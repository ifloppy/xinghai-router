<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Check, Code2, Copy, Info, X } from 'lucide-vue-next'
import { formatRatio, formatSquarePrice, groupPrice, vendorColor, vendorIconUrl, type SquareModel, type TokenUnit } from '~/src/marketplace'

const props = defineProps<{
  model: SquareModel
  tokenUnit: TokenUnit
  origin: string
}>()
const emit = defineEmits<{ close: [] }>()
const { t, locale } = useI18n()

const tab = ref<'overview' | 'api'>('overview')
const iconError = ref(false)
const copiedName = ref(false)
const copiedCurl = ref(false)
const unitLabel = computed(() => (props.tokenUnit === 'K' ? '1K' : '1M'))
const hasCache = computed(() => props.model.cached_input_per_million != null && Number(props.model.cached_input_per_million) !== 0)

function price(kind: 'input' | 'output' | 'cache', group?: (props.model.groups)[number]) {
  const value = group
    ? groupPrice(props.model, kind, group)
    : groupPrice(props.model, kind, props.model.groups[0] ?? { id: '', name: '', multiplier: 1 })
  return formatSquarePrice(value, props.tokenUnit)
}

const baseGroupFree = computed(() => ({
  input: formatSquarePrice(props.model.input_per_million == null ? null : Number(props.model.input_per_million) * Number(props.model.multiplier ?? 1), props.tokenUnit),
  output: formatSquarePrice(props.model.output_per_million == null ? null : Number(props.model.output_per_million) * Number(props.model.multiplier ?? 1), props.tokenUnit),
  cache: formatSquarePrice(props.model.cached_input_per_million == null || Number(props.model.cached_input_per_million) === 0 ? null : Number(props.model.cached_input_per_million) * Number(props.model.multiplier ?? 1), props.tokenUnit),
}))

const curlExample = computed(() => `curl ${props.origin}/v1/chat/completions \\
  -H "Authorization: Bearer sk-xh-your-key" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"${props.model.model}","messages":[{"role":"user","content":"${locale.value === 'en-US' ? 'Hello' : '你好'}"}]}'`)

const endpoints = computed(() => [
  { method: 'POST', path: '/v1/chat/completions', label: 'OpenAI Chat Completions' },
  { method: 'GET', path: '/v1/models', label: 'OpenAI Models' },
  { method: 'POST', path: '/v1/messages', label: 'Anthropic Messages' },
])

async function copyText(value: string, flag: 'name' | 'curl') {
  try {
    await navigator.clipboard.writeText(value)
  } catch {
    const textarea = document.createElement('textarea')
    textarea.value = value
    document.body.append(textarea)
    textarea.select()
    document.execCommand('copy')
    textarea.remove()
  }
  if (flag === 'name') { copiedName.value = true; window.setTimeout(() => { copiedName.value = false }, 1500) }
  else { copiedCurl.value = true; window.setTimeout(() => { copiedCurl.value = false }, 1500) }
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') emit('close')
}
onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  document.body.style.overflow = 'hidden'
})
onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})
watch(() => props.model.model, () => { tab.value = 'overview' })
</script>

<template>
  <Teleport to="body">
    <div class="msq-drawer-backdrop" @click="emit('close')"></div>
    <div class="msq-drawer" role="dialog" :aria-label="props.model.model">
      <button type="button" class="msq-drawer-close" :aria-label="t('msClose')" @click="emit('close')"><X :size="16" /></button>
      <div class="msq-drawer-body">
        <header class="msq-drawer-header">
          <div class="msq-drawer-title">
            <span class="msq-iconbox small">
              <img v-if="props.model.vendor_slug && !iconError" :src="vendorIconUrl(props.model.vendor_slug)" :alt="props.model.vendor_name" @error="iconError = true" />
              <span v-else>{{ props.model.model.slice(0, 1).toUpperCase() }}</span>
            </span>
            <h1>{{ props.model.model }}</h1>
            <button type="button" class="msq-copy-btn" :title="copiedName ? t('msCopied') : t('msCopyModelName')" @click="copyText(props.model.model, 'name')">
              <Check v-if="copiedName" :size="12" />
              <Copy v-else :size="12" />
            </button>
          </div>
          <div class="msq-drawer-meta">
            <span v-if="props.model.vendor_name">{{ props.model.vendor_name }}</span>
            <i>·</i>
            <span class="msq-billing-badge">{{ t('msTokenBased') }}</span>
          </div>
          <p v-if="props.model.description" class="msq-drawer-desc">{{ props.model.description }}</p>
        </header>

        <div class="msq-tabs">
          <button type="button" :class="{ active: tab === 'overview' }" @click="tab = 'overview'"><Info :size="13" />{{ t('msOverview') }}</button>
          <button type="button" :class="{ active: tab === 'api' }" @click="tab = 'api'"><Code2 :size="13" />API</button>
        </div>

        <div v-if="tab === 'overview'" class="msq-tab-body">
          <section class="msq-panel">
            <h2 class="msq-section-title">{{ t('msPricing') }}</h2>

            <h3 class="msq-sub-title">{{ t('msBasePrice') }}</h3>
            <div class="msq-price-cards">
              <div class="msq-price-card">
                <span>{{ t('inputLabel') }}</span>
                <b v-if="baseGroupFree.input">{{ baseGroupFree.input }}<small>/ {{ unitLabel }}</small></b>
                <b v-else class="msq-price-pending">{{ t('pendingConfig') }}</b>
              </div>
              <div class="msq-price-card">
                <span>{{ t('outputLabel') }}</span>
                <b v-if="baseGroupFree.output">{{ baseGroupFree.output }}<small>/ {{ unitLabel }}</small></b>
                <b v-else class="msq-price-pending">{{ t('pendingConfig') }}</b>
              </div>
            </div>
            <div v-if="baseGroupFree.cache" class="msq-price-secondary">
              <div>
                <span>{{ t('cachedInputLabel') }}</span>
                <b class="msq-mono">{{ baseGroupFree.cache }}<small>/ {{ unitLabel }}</small></b>
              </div>
            </div>

            <h3 class="msq-sub-title">{{ t('msGroupPricing') }}</h3>
            <table v-if="props.model.groups.length" class="msq-group-table">
              <thead>
                <tr>
                  <th>{{ t('msGroup') }}</th>
                  <th>{{ t('msRatio') }}</th>
                  <th class="num">{{ t('inputLabel') }}</th>
                  <th class="num">{{ t('outputLabel') }}</th>
                  <th v-if="hasCache" class="num">{{ t('msCached') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="group in props.model.groups" :key="group.id">
                  <td><span class="msq-group-badge">{{ group.name }}</span></td>
                  <td class="msq-mono muted">{{ formatRatio(Number(group.multiplier)) }}</td>
                  <td class="num msq-mono">{{ price('input', group) || '—' }}</td>
                  <td class="num msq-mono">{{ price('output', group) || '—' }}</td>
                  <td v-if="hasCache" class="num msq-mono">{{ price('cache', group) || '—' }}</td>
                </tr>
              </tbody>
            </table>
            <p v-else class="msq-muted-text">{{ t('msNoGroups') }}</p>
            <p class="msq-footnote">{{ t('msPriceNoteA') }} {{ unitLabel }} {{ t('msPriceNoteB') }}</p>
          </section>

          <section class="msq-panel">
            <h2 class="msq-section-title">{{ t('msModel') }}</h2>
            <div class="msq-info-grid">
              <div class="msq-info-cell">
                <span>{{ t('msVendor') }}</span>
                <b>
                  <em class="msq-vendor-badge" :style="{ background: vendorColor(props.model.vendor_name).bg, color: vendorColor(props.model.vendor_name).fg }">{{ props.model.vendor_name }}</em>
                </b>
              </div>
              <div class="msq-info-cell">
                <span>{{ t('msType') }}</span>
                <b><em class="msq-billing-badge">{{ t('msTokenBased') }}</em></b>
              </div>
              <div class="msq-info-cell">
                <span>{{ t('msGroups') }}</span>
                <b class="msq-info-groups"><em v-for="group in props.model.groups" :key="group.id" class="msq-group-badge">{{ group.name }}</em></b>
              </div>
              <div class="msq-info-cell">
                <span>{{ t('msEndpoints') }}</span>
                <b class="msq-info-groups"><em class="msq-group-badge">openai</em><em class="msq-group-badge">anthropic</em></b>
              </div>
            </div>
          </section>
        </div>

        <div v-else class="msq-tab-body">
          <section class="msq-panel">
            <h2 class="msq-section-title">{{ t('msAvailableEndpoints') }}</h2>
            <div class="msq-endpoint-list">
              <div v-for="endpoint in endpoints" :key="endpoint.path" class="msq-endpoint">
                <span :class="['msq-method', endpoint.method.toLowerCase()]">{{ endpoint.method }}</span>
                <code>{{ endpoint.path }}</code>
                <small>{{ endpoint.label }}</small>
              </div>
            </div>
          </section>
          <section class="msq-panel">
            <div class="msq-code-head">
              <h2 class="msq-section-title">{{ t('msRequestExample') }}</h2>
              <button type="button" class="msq-copy-btn bordered" @click="copyText(curlExample, 'curl')">
                <Check v-if="copiedCurl" :size="12" />
                <Copy v-else :size="12" />
                {{ copiedCurl ? t('msCopied') : t('msCopy') }}
              </button>
            </div>
            <pre class="msq-code"><code>{{ curlExample }}</code></pre>
          </section>
        </div>
      </div>
    </div>
  </Teleport>
</template>
