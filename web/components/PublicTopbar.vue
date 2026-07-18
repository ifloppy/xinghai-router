<script setup lang="ts">
import { Bot, ChevronRight } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { api, getToken, type SiteSettings } from '~/src/api'

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
      const value = await api<SiteSettings>('/site-settings')
      fetchedName.value = value.name
    } catch { /* fall back to default name */ }
  }
  if (props.authenticated === undefined && getToken()) {
    try {
      await api('/account/me')
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
  <nav class="landing-nav">
    <a class="landing-logo" href="/"><span class="brand-mark small"><Bot :size="19" /></span><span>{{ displayName }}</span></a>
    <div class="landing-links">
      <a :href="featuresHref">{{ t('landingFeatures') }}</a>
      <a href="/rankings">{{ t('rankings') }}</a>
      <a :href="quickstartHref">{{ t('quickStart') }}</a>
      <a href="/models">{{ t('marketplace') }}</a>
      <ThemeCustomizer :locale="locale" />
      <select v-model="locale" class="language-select" :aria-label="t('switchLanguage')">
        <option value="zh-CN">{{ t('chinese') }}</option>
        <option value="en-US">{{ t('english') }}</option>
      </select>
      <button class="button ghost" @click="openConsoleOrAuth">{{ isAuthenticated ? t('console') : t('login') }} <ChevronRight :size="15" /></button>
    </div>
  </nav>
</template>
