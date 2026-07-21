<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { Bot, ChevronRight, RadioTower, Tags, Activity, ShieldCheck } from 'lucide-vue-next'
import { endpoints, setToken } from '~/src/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useConsoleStore } from '~/composables/useConsoleStore'

const store = useConsoleStore()
const { t, locale, router, siteSettings, authenticated, busy, action, load } = store

const loginMode = ref<'token' | 'login' | 'register'>('login')
const accountForm = reactive({ name: '', email: '', password: '' })

// Geetest v4 CAPTCHA — loaded lazily the first time sign-in requires it.
declare global { interface Window { initGeetest4?: (options: Record<string, unknown>, callback: (captcha: GeetestCaptcha) => void) => void } }
interface GeetestCaptcha { onReady(fn: () => void): void; onSuccess(fn: () => void): void; onClose(fn: () => void): void; onError(fn: (cause: unknown) => void): void; showCaptcha(): void; getValidate(): Record<string, string> | null }
let geetestScript: Promise<void> | null = null
let geetestInstance: GeetestCaptcha | null = null
let geetestReady: Promise<void> | null = null
function loadGeetestScript() { geetestScript ??= new Promise<void>((resolve, reject) => { const script = document.createElement('script'); script.src = 'https://static.geetest.com/v4/gt4.js'; script.onload = () => resolve(); script.onerror = () => reject(new Error('captcha script failed')); document.head.appendChild(script) }); return geetestScript }
function ensureGeetest() { if (!siteSettings.value.geetest_enabled || !siteSettings.value.geetest_captcha_id) return Promise.resolve(); geetestReady ??= loadGeetestScript().then(() => new Promise<void>((resolve, reject) => { window.initGeetest4?.({ captchaId: siteSettings.value.geetest_captcha_id, product: 'float', language: locale.value === 'en-US' ? 'eng' : 'zho' }, (captcha) => { geetestInstance = captcha; captcha.onReady(() => resolve()); captcha.onError(reject) }) })); return geetestReady }
function runGeetest(): Promise<Record<string, string>> { return new Promise((resolve, reject) => { const captcha = geetestInstance; if (!captcha) return resolve({}); const cleanup = () => { captcha.onSuccess(() => {}); captcha.onClose(() => {}) }; captcha.onSuccess(() => { const result = captcha.getValidate(); cleanup(); if (result) resolve(result); else reject(new Error('captcha failed')) }); captcha.onClose(() => { cleanup(); reject(new Error('captcha closed')) }); captcha.showCaptcha() }) }

// Registration email verification code.
const emailCode = ref('')
const codeSending = ref(false)
const codeCountdown = ref(0)
const codeSentHint = ref('')
let codeTimer = 0
const emailLooksValid = computed(() => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(accountForm.email.trim()))
async function sendEmailCode() {
  if (!emailLooksValid.value || codeSending.value || codeCountdown.value > 0) return
  let captcha: Record<string, string> = {}
  if (siteSettings.value.geetest_enabled) { try { await ensureGeetest(); captcha = await runGeetest() } catch { return } }
  codeSending.value = true
  try {
    await endpoints.sendEmailCode(accountForm.email.trim(), captcha)
    codeSentHint.value = t('codeSent')
    codeCountdown.value = 60
    window.clearInterval(codeTimer)
    codeTimer = window.setInterval(() => { codeCountdown.value -= 1; if (codeCountdown.value <= 0) window.clearInterval(codeTimer) }, 1000)
  } catch (cause) {
    codeSentHint.value = cause instanceof Error ? cause.message : t('operationFailed')
  } finally {
    codeSending.value = false
  }
}

async function accountSignIn(register: boolean) {
  let captcha: Record<string, string> = {}
  const emailVerify = register && siteSettings.value.email_verification_enabled
  if (siteSettings.value.geetest_enabled && !emailVerify) { try { await ensureGeetest(); captcha = await runGeetest() } catch { return } }
  await action(async () => {
    const body = register ? { ...accountForm, ...(emailVerify ? { code: emailCode.value.trim() } : {}), ...captcha } : { email: accountForm.email, password: accountForm.password, ...(emailVerify ? { code: emailCode.value.trim() } : {}), ...captcha }
    const result = register ? await endpoints.register(body) : await endpoints.login(body)
    setToken(result.token); authenticated.value = true; await load(); await router.replace({ path: '/console', query: { view: store.managementNav.value.length ? 'overview' : 'account' } })
  })
}
</script>

<template>
  <div class="public-wrap">
    <PublicTopbar :site-name="siteSettings.name" :authenticated="authenticated" />
    <main class="login-shell">
      <section class="login-card"><aside class="login-aside"><div class="login-aside-glow"/><div class="login-aside-inner"><span class="brand-mark"><Bot :size="28" /></span><h2>{{ t('loginBrandTitle') }}</h2><p>{{ t('loginBrandDesc') }}</p><ul><li><i><RadioTower :size="14" /></i><span>{{ t('brandPoint1') }}</span></li><li><i><Tags :size="14" /></i><span>{{ t('brandPoint2') }}</span></li><li><i><Activity :size="14" /></i><span>{{ t('brandPoint3') }}</span></li><li><i><ShieldCheck :size="14" /></i><span>{{ t('brandPoint4') }}</span></li></ul></div></aside><div class="login-pane"><h1>{{ loginMode === 'register' ? t('createAccountTab') : t('signInTab') }}</h1><p class="login-sub">{{ loginMode === 'register' ? t('registerSub') : t('loginSub') }}</p><div class="flex gap-1 rounded-md bg-muted p-1"><button class="flex-1 rounded px-3 py-1.5 text-sm font-medium transition-colors" :class="loginMode === 'login' ? 'bg-background shadow-sm' : 'text-muted-foreground'" @click="loginMode = 'login'">{{ t('signInTab') }}</button><button class="flex-1 rounded px-3 py-1.5 text-sm font-medium transition-colors" :class="loginMode === 'register' ? 'bg-background shadow-sm' : 'text-muted-foreground'" @click="loginMode = 'register'">{{ t('createAccountTab') }}</button></div><form class="flex flex-col gap-3" @submit.prevent="accountSignIn(loginMode === 'register')"><div v-if="loginMode === 'register'" class="flex flex-col gap-1.5"><Label class="text-xs">{{ t('nameLabel') }}</Label><Input v-model="accountForm.name" autocomplete="name" required maxlength="100" :placeholder="t('namePlaceholder')" /></div><div class="flex flex-col gap-1.5"><Label class="text-xs">{{ t('emailLabel') }}</Label><Input v-model="accountForm.email" type="email" autocomplete="email" required placeholder="name@example.com" /></div><div v-if="loginMode === 'register' && siteSettings.email_verification_enabled" class="flex flex-col gap-1.5"><Label class="text-xs">{{ t('emailCodeLabel') }}</Label><div class="flex gap-2"><Input v-model="emailCode" inputmode="numeric" autocomplete="one-time-code" maxlength="6" required :placeholder="t('emailCodePlaceholder')" class="flex-1" /><Button variant="outline" type="button" :disabled="!emailLooksValid || codeSending || codeCountdown > 0" @click="sendEmailCode">{{ codeCountdown > 0 ? t('codeCountdown').replace('{n}', String(codeCountdown)) : codeSending ? t('codeSending') : t('sendCode') }}</Button></div><small v-if="codeSentHint" class="text-xs text-muted-foreground">{{ codeSentHint }}</small></div><div class="flex flex-col gap-1.5"><Label class="text-xs">{{ t('passwordLabel') }}</Label><Input v-model="accountForm.password" type="password" :autocomplete="loginMode === 'register' ? 'new-password' : 'current-password'" required minlength="8" :placeholder="t('passwordMinLength')" /></div><Button type="submit" class="mt-1 w-full" :disabled="busy">{{ loginMode === 'register' ? t('createAndOpenConsole') : t('signInConsole') }} <ChevronRight :size="16" /></Button></form><p class="auth-legal">{{ t('agreeText') }} <a href="/terms">{{ t('termsShort') }}</a> {{ t('andConnector') }} <a href="/privacy">{{ t('privacyShort') }}</a>.</p></div></section>
    </main>
  </div>
</template>