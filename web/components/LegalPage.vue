<script setup lang="ts">
import { Bot, ChevronLeft, Moon, Sun } from 'lucide-vue-next'

const props = defineProps<{ kind: 'terms' | 'privacy' }>()
const { locale, t, toggleLocale, initializeLocale } = useI18n()
const theme = ref<'light' | 'dark'>('light')

const isTerms = computed(() => props.kind === 'terms')
const title = computed(() => isTerms.value ? t('termsTitle') : t('privacyTitle'))
const updatedAt = '2026-07-16'

function setTheme(next: 'light' | 'dark') {
  theme.value = next
  document.documentElement.dataset.theme = next
  localStorage.setItem('xinghai-router-theme', next)
}

onMounted(() => {
  initializeLocale()
  const saved = localStorage.getItem('xinghai-router-theme')
  setTheme(saved === 'dark' || saved === 'light' ? saved : window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
})
</script>

<template>
  <main class="legal-shell">
    <nav class="legal-nav">
      <a class="landing-logo" href="/"><span class="brand-mark small"><Bot :size="19" /></span><span>Xinghai Router</span></a>
      <div class="legal-nav-actions">
        <button class="theme-toggle" :aria-label="theme === 'dark' ? t('lightMode') : t('darkMode')" @click="setTheme(theme === 'dark' ? 'light' : 'dark')"><Sun v-if="theme === 'dark'" :size="16" /><Moon v-else :size="16" /></button>
        <button class="language-toggle" :aria-label="t('switchLanguage')" @click="toggleLocale">{{ locale === 'zh-CN' ? t('english') : t('chinese') }}</button>
        <a class="button ghost" href="/"><ChevronLeft :size="15" />{{ t('backHome') }}</a>
      </div>
    </nav>

    <article class="legal-document">
      <header class="legal-header">
        <p class="eyebrow">XINGHAI ROUTER / {{ isTerms ? 'TERMS' : 'PRIVACY' }}</p>
        <h1>{{ title }}</h1>
        <p>{{ t('updatedAtLabel') }}{{ updatedAt }}</p>
      </header>

      <template v-if="locale === 'zh-CN' && isTerms">
        <p class="legal-lead">{{ t('termsServiceLead') }}</p>
        <section><h2>{{ t('termsServiceTitle1') }}</h2><p>{{ t('termsServiceBody1') }}</p></section>
        <section><h2>{{ t('termsServiceTitle2') }}</h2><p>{{ t('termsServiceBody2') }}</p></section>
        <section><h2>{{ t('termsServiceTitle3') }}</h2><p>{{ t('termsServiceBody3') }}</p></section>
        <section><h2>{{ t('termsServiceTitle4') }}</h2><p>{{ t('termsServiceBody4') }}</p></section>
        <section><h2>{{ t('termsServiceTitle5') }}</h2><p>{{ t('termsServiceBody5') }}</p></section>
        <section><h2>{{ t('termsServiceTitle6') }}</h2><p>{{ t('termsServiceBody6') }}</p></section>
        <section><h2>{{ t('termsServiceTitle7') }}</h2><p>{{ t('termsServiceBody7') }}</p></section>
        <section><h2>{{ t('termsServiceTitle8') }}</h2><p>{{ t('termsServiceBody8') }}</p></section>
      </template>

      <template v-else-if="locale === 'zh-CN'">
        <p class="legal-lead">{{ t('privacyPolicyLead') }}</p>
        <section><h2>{{ t('privacyTitle1') }}</h2><p>{{ t('privacyBody1') }}</p></section>
        <section><h2>{{ t('privacyTitle2') }}</h2><p>{{ t('privacyBody2') }}</p></section>
        <section><h2>{{ t('privacyTitle3') }}</h2><p>{{ t('privacyBody3') }}</p></section>
        <section><h2>{{ t('privacyTitle4') }}</h2><p>{{ t('privacyBody4') }}</p></section>
        <section><h2>{{ t('privacyTitle5') }}</h2><p>{{ t('privacyBody5') }}</p></section>
        <section><h2>{{ t('privacyTitle6') }}</h2><p>{{ t('privacyBody6') }}</p></section>
        <section><h2>{{ t('privacyTitle7') }}</h2><p>{{ t('privacyBody7') }}</p></section>
      </template>

      <template v-else-if="isTerms">
        <p class="legal-lead">{{ t('termsServiceLead') }}</p>
        <section><h2>{{ t('termsServiceTitle1') }}</h2><p>{{ t('termsServiceBody1') }}</p></section>
        <section><h2>{{ t('termsServiceTitle2') }}</h2><p>{{ t('termsServiceBody2') }}</p></section>
        <section><h2>{{ t('termsServiceTitle3') }}</h2><p>{{ t('termsServiceBody3') }}</p></section>
        <section><h2>{{ t('termsServiceTitle4') }}</h2><p>{{ t('termsServiceBody4') }}</p></section>
        <section><h2>{{ t('termsServiceTitle5') }}</h2><p>{{ t('termsServiceBody5') }}</p></section>
        <section><h2>{{ t('termsServiceTitle6') }}</h2><p>{{ t('termsServiceBody6') }}</p></section>
        <section><h2>{{ t('termsServiceTitle7') }}</h2><p>{{ t('termsServiceBody7') }}</p></section>
        <section><h2>{{ t('termsServiceTitle8') }}</h2><p>{{ t('termsServiceBody8') }}</p></section>
      </template>

      <template v-else>
        <p class="legal-lead">{{ t('privacyPolicyLead') }}</p>
        <section><h2>{{ t('privacyTitle1') }}</h2><p>{{ t('privacyBody1') }}</p></section>
        <section><h2>{{ t('privacyTitle2') }}</h2><p>{{ t('privacyBody2') }}</p></section>
        <section><h2>{{ t('privacyTitle3') }}</h2><p>{{ t('privacyBody3') }}</p></section>
        <section><h2>{{ t('privacyTitle4') }}</h2><p>{{ t('privacyBody4') }}</p></section>
        <section><h2>{{ t('privacyTitle5') }}</h2><p>{{ t('privacyBody5') }}</p></section>
        <section><h2>{{ t('privacyTitle6') }}</h2><p>{{ t('privacyBody6') }}</p></section>
        <section><h2>{{ t('privacyTitle7') }}</h2><p>{{ t('privacyBody7') }}</p></section>
      </template>
    </article>
  </main>
</template>
