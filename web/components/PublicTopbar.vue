<script setup lang="ts">
import { Bot, ChevronRight } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { getToken } from '~/src/api'
import { Button } from '@/components/ui/button'

const props = withDefaults(defineProps<{
  siteName?: string
  authenticated?: boolean
}>(), { siteName: '', authenticated: undefined })

const { locale, t } = useI18n()
const route = useRoute()
const router = useRouter()

const fetchedName = ref('')
const selfAuthenticated = ref(false)

onMounted(async () => {
  if (!props.siteName) {
    try {
      const value = await endpoints.getSiteSettings()
      fetchedName.value = value.name
    } catch { /* fall back to default name */ }
  }
  if (props.authenticated === undefined && getToken()) {
    try {
      await endpoints.getAccount()
      selfAuthenticated.value = true
    } catch {
      selfAuthenticated.value = false
    }
  }
})

const displayName = computed(() => props.siteName || fetchedName.value || 'Xinghai Router')
const isAuthenticated = computed(() => props.authenticated ?? selfAuthenticated.value)
const isHome = computed(() => route.path === '/')
const featuresHref = computed(() => (isHome.value ? '#features' : '/#features'))
const quickstartHref = computed(() => (isHome.value ? '#quickstart' : '/#quickstart'))

function openConsoleOrAuth() {
  router.push(isAuthenticated.value ? '/console/overview' : '/auth')
}
</script>

<template>
  <nav class="mx-auto flex h-16 w-[min(1160px,calc(100%-48px))] items-center justify-between gap-4 sm:h-20 max-[700px]:h-[70px] max-[700px]:w-[min(100%-32px,560px)]">
    <a class="flex min-w-0 shrink-0 items-center gap-2.5 font-extrabold tracking-tight text-foreground transition-opacity hover:opacity-80" href="/">
      <span class="grid h-9 w-9 shrink-0 place-items-center rounded-[var(--radius-xl)] bg-primary text-primary-foreground shadow-md ring-1 ring-primary/10 transition-transform hover:scale-105 sm:h-11 sm:w-11"><Bot :size="20" /></span>
      <span class="truncate">{{ displayName }}</span>
    </a>
    <div class="flex shrink-0 items-center gap-3 sm:gap-7">
      <a class="whitespace-nowrap text-xs text-muted-foreground transition-colors hover:text-foreground max-[700px]:hidden" :href="featuresHref">{{ t('landingFeatures') }}</a>
      <a class="whitespace-nowrap text-xs text-muted-foreground transition-colors hover:text-foreground max-[700px]:hidden" href="/rankings">{{ t('rankings') }}</a>
      <a class="whitespace-nowrap text-xs text-muted-foreground transition-colors hover:text-foreground max-[700px]:hidden" :href="quickstartHref">{{ t('quickStart') }}</a>
      <a class="whitespace-nowrap text-xs text-muted-foreground transition-colors hover:text-foreground max-[700px]:hidden" href="/models">{{ t('marketplace') }}</a>
      <ThemeCustomizer :locale="locale" />
      <select v-model="locale" class="h-8 w-[68px] min-w-[68px] rounded-md border border-border bg-card px-1 text-[10px] text-muted-foreground transition-colors hover:bg-accent focus:outline-none focus-visible:ring-1 focus-visible:ring-ring" :aria-label="t('switchLanguage')">
        <option value="zh-CN">{{ t('chinese') }}</option>
        <option value="en-US">{{ t('english') }}</option>
      </select>
      <Button variant="ghost" size="sm" class="shrink-0" @click="openConsoleOrAuth">{{ isAuthenticated ? t('console') : t('login') }} <ChevronRight :size="15" /></Button>
    </div>
  </nav>
</template>