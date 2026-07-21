<script setup lang="ts">
import { ArrowDown, ArrowUp, RefreshCw, Trophy } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import type { Rankings, SiteSettings } from '~/src/api'
import { Button } from '@/components/ui/button'

type Period = 'today' | 'week' | 'month' | 'year'
const period = ref<Period>('week')
const rankings = ref<Rankings | null>(null)
const siteSettings = ref<SiteSettings>({ name: 'Xinghai Router', icon_url: '', auto_disable_failed_channels: false })
const loading = ref(true)
const error = ref('')
const { locale, t, initializeLocale } = useI18n()
const periods: { value: Period; label: string }[] = [{ value: 'today', label: t('today') }, { value: 'week', label: t('thisWeek') }, { value: 'month', label: t('thisMonth') }, { value: 'year', label: t('thisYear') }]

const compactNumber = (value: number) => new Intl.NumberFormat(locale.value, { notation: 'compact', maximumFractionDigits: 1 }).format(value)
const change = (value: number) => `${value >= 0 ? '+' : ''}${value.toFixed(1)}%`
const share = (value: number) => value > 0 && value < .001 ? '<0.1%' : `${(value * 100).toFixed(1)}%`

async function load(next = period.value) {
  period.value = next
  loading.value = true
  error.value = ''
  try {
    const response = await fetch(`/api/rankings?period=${next}`)
    if (!response.ok) throw new Error(t('rankingsUnavailable'))
    rankings.value = await response.json()
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('rankingsUnavailable')
  } finally {
    loading.value = false
  }
}
async function loadSiteSettings() {
  const response = await fetch('/api/site-settings')
  if (!response.ok) return
  siteSettings.value = await response.json() as SiteSettings
  document.title = `${siteSettings.value.name} · ${t('titleRankings')}`
  if (siteSettings.value.icon_url) {
    const link = document.querySelector<HTMLLinkElement>('link[rel="icon"]') ?? document.head.appendChild(Object.assign(document.createElement('link'), { rel: 'icon' }))
    link.href = siteSettings.value.icon_url
  }
}

function selectPeriod(next: Period) {
  history.replaceState({}, '', `/rankings?period=${next}`)
  load(next)
}

onMounted(() => {
  initializeLocale()
  const queryPeriod = new URLSearchParams(location.search).get('period') as Period | null
  if (queryPeriod && periods.some((item) => item.value === queryPeriod)) period.value = queryPeriod
  loadSiteSettings()
  load()
})
</script>

<template>
  <main class="rankings-page">
     <PublicTopbar />

    <section class="rankings-shell">
      <header class="rankings-hero">
        <div><span><Trophy :size="15" /> LIVE USAGE RANKINGS</span><h1>{{ t('modelRankingsTitle') }}</h1><p>{{ t('rankingsDesc') }}</p></div>
        <div v-if="rankings" class="ranking-total"><strong>{{ compactNumber(rankings.total_tokens) }}</strong><small>{{ t('periodTokens') }}</small></div>
      </header>
      <div class="rankings-controls"><div><span>{{ t('timeRange') }}</span><div class="flex gap-1 rounded-lg border border-border bg-card p-1"><button v-for="item in periods" :key="item.value" class="min-w-[58px] rounded-md px-3 py-1.5 text-[10px] font-bold transition-colors disabled:opacity-50" :class="period === item.value ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-accent'" :disabled="loading" @click="selectPeriod(item.value)">{{ item.label }}</button></div></div><Button variant="ghost" size="sm" class="ranking-refresh" :disabled="loading" @click="load()"><RefreshCw :size="15" :class="{ 'animate-spin': loading }" />{{ t('refreshData') }}</Button></div>

      <div v-if="loading && !rankings" class="rankings-loading"><div/><div/><div/></div>
      <section v-else-if="error && !rankings" class="rankings-error"><h2>{{ t('cannotLoadRankings') }}</h2><p>{{ error }}</p><Button variant="outline" @click="load()">{{ t('reload') }}</Button></section>
      <template v-else-if="rankings">
        <section class="ranking-panel">
          <div class="ranking-panel-title"><div><span>TOP MODELS</span><h2>{{ t('llmRankings') }}</h2><p>{{ t('rankingsSortByTokens') }}</p></div><b>{{ rankings.models.length }} {{ t('modelCount') }}</b></div>
          <div v-if="rankings.models.length" class="model-rankings"><article v-for="item in rankings.models" :key="item.model_name"><b :class="['rank-number', { podium: item.rank <= 3 }]">{{ String(item.rank).padStart(2, '0') }}</b><div class="rank-model"><i>{{ item.model_name.slice(0, 1).toUpperCase() }}</i><span><strong>{{ item.model_name }}</strong><small>{{ item.vendor }}</small></span></div><div class="rank-share"><i><span :style="{ width: `${Math.max(item.share * 100, 1)}%` }"/></i><small>{{ share(item.share) }}</small></div><div class="rank-tokens"><strong>{{ compactNumber(item.total_tokens) }}</strong><small>Token</small></div><em :class="item.growth_pct < 0 ? 'down' : 'up'"><ArrowDown v-if="item.growth_pct < 0" :size="11" /><ArrowUp v-else :size="11" />{{ change(item.growth_pct) }}</em></article></div>
          <div v-else class="ranking-empty">{{ t('noModelUsage') }}</div>
        </section>

        <section class="ranking-panel">
          <div class="ranking-panel-title"><div><span>MARKET SHARE</span><h2>{{ t('vendorShare') }}</h2><p>{{ t('vendorShareDesc') }}</p></div></div>
          <div v-if="rankings.vendors.length" class="vendor-list"><article v-for="item in rankings.vendors.slice(0, 12)" :key="item.vendor"><b>{{ String(item.rank).padStart(2, '0') }}</b><div><strong>{{ item.vendor }}</strong><small>{{ item.models_count }} {{ t('modelCount') }} · {{ t('topModelLabel') }} {{ item.top_model }}</small></div><i><span :style="{ width: `${item.share * 100}%` }"/></i><span>{{ share(item.share) }}</span><em :class="item.growth_pct < 0 ? 'down' : 'up'">{{ change(item.growth_pct) }}</em></article></div>
          <div v-else class="ranking-empty">{{ t('noVendorData') }}</div>
        </section>

        <div class="movers-grid"><section class="ranking-panel"><div class="ranking-panel-title compact"><div><span>TRENDING UP</span><h2>{{ t('trendingUp') }}</h2></div><ArrowUp :size="18" /></div><div class="mover-list"><article v-for="item in rankings.top_movers" :key="item.model_name"><div><strong>{{ item.model_name }}</strong><small>#{{ item.current_rank }} · {{ item.vendor }}</small></div><b class="up"><ArrowUp :size="13" />{{ item.rank_delta }}</b></article><div v-if="!rankings.top_movers.length" class="ranking-empty small">{{ t('noTrendingUpModels') }}</div></div></section><section class="ranking-panel"><div class="ranking-panel-title compact"><div><span>TRENDING DOWN</span><h2>{{ t('trendingDown') }}</h2></div><ArrowDown :size="18" /></div><div class="mover-list"><article v-for="item in rankings.top_droppers" :key="item.model_name"><div><strong>{{ item.model_name }}</strong><small>#{{ item.current_rank }} · {{ item.vendor }}</small></div><b class="down"><ArrowDown :size="13" />{{ Math.abs(item.rank_delta) }}</b></article><div v-if="!rankings.top_droppers.length" class="ranking-empty small">{{ t('noTrendingDownModels') }}</div></div></section></div>
        <p class="rankings-updated">{{ t('dataUpdatedAt') }} {{ new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(rankings.updated_at)) }}</p>
      </template>
    </section>
  </main>
</template>

