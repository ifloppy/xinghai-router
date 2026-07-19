<script setup lang="ts">
const props = defineProps<{ kind: 'terms' | 'privacy' }>()
const { locale, t, initializeLocale } = useI18n()

const isTerms = computed(() => props.kind === 'terms')
const title = computed(() => isTerms.value ? t('termsTitle') : t('privacyTitle'))
const updatedAt = '2026-07-16'

onMounted(() => {
  initializeLocale()
})
</script>

<template>
  <main class="legal-shell">
    <PublicTopbar />

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
