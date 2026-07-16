<script setup lang="ts">
import { computed } from 'vue'

const route = useRoute()
const { locale, t } = useI18n()
const { initializeTheme } = useTheme()
const consoleTitles: Record<string, string> = {
  overview: 'overview', users: 'users', groups: 'groups', keys: 'keys', channels: 'channels', logs: 'logs', account: 'account', profile: 'profile', wallet: 'wallet', usage: 'usage', 'usage-overview': 'usageOverview', ledger: 'ledger', pricing: 'pricing', audit: 'audit',
}

const title = computed(() => {
  if (route.path === '/') return 'Xinghai Router'
  if (route.path === '/auth') return `${t('titleLogin')} | Xinghai Router`
  if (route.path === '/models') return `${t('titleMarketplace')} | Xinghai Router`
  if (route.path === '/rankings') return `${t('titleRankings')} | Xinghai Router`
  if (route.path === '/terms') return `${t('termsTitle')} | Xinghai Router`
  if (route.path === '/privacy') return `${t('privacyTitle')} | Xinghai Router`

  const queryView = route.query.view
  const view = route.path === '/console'
    ? typeof queryView === 'string' ? queryView : 'overview'
    : typeof route.params.view === 'string' ? route.params.view : route.path.slice('/console/'.length)
  return `${t((consoleTitles[view] ?? 'overview') as Parameters<typeof t>[0])} | Xinghai Router`
})

useHead({ title })
onMounted(initializeTheme)
</script>

<template>
  <NuxtPage />
</template>
