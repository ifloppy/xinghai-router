<script setup lang="ts">
import { Bot, ChevronRight, RadioTower, Activity, WalletCards } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

defineProps<{ siteName: string; authenticated: boolean }>()
const emit = defineEmits<{ openConsoleOrAuth: [] }>()
const { t } = useI18n()
</script>

<template>
  <main class="relative min-h-screen overflow-hidden bg-background text-foreground">
    <div class="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_80%_50%_at_20%_-20%,hsl(var(--primary)/0.08),transparent),radial-gradient(ellipse_60%_40%_at_80%_10%,hsl(var(--accent)/0.06),transparent)]" />
    <PublicTopbar :site-name="siteName" :authenticated="authenticated" />

    <section class="relative z-10 mx-auto grid w-[min(1160px,calc(100%-48px))] min-h-[590px] grid-cols-1 items-center gap-16 lg:grid-cols-2">
      <div class="max-w-[560px]">
        <p class="mb-2 mt-5 text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">OPENAI-COMPATIBLE MODEL GATEWAY</p>
        <h1 class="m-0 max-w-[600px] bg-gradient-to-br from-foreground to-foreground/70 bg-clip-text text-[clamp(42px,5vw,70px)] font-bold leading-[1.03] tracking-[-0.075em] text-transparent">{{ t('heroTagline') }}</h1>
        <p class="mt-6 max-w-[460px] text-[15px] leading-[1.85] text-muted-foreground">{{ t('heroDescription') }}</p>
        <div class="mt-8 flex items-center gap-6">
          <Button size="lg" @click="emit('openConsoleOrAuth')">{{ t('openConsole') }} <ChevronRight :size="16" /></Button>
          <a class="text-xs font-bold text-muted-foreground transition-colors hover:text-foreground" href="#quickstart">{{ t('viewRequestExample') }}<span class="ml-1 text-[17px]">→</span></a>
        </div>
      </div>

      <div class="relative grid min-h-[440px] place-items-center">
        <div class="absolute h-[390px] w-[390px] rounded-full bg-[radial-gradient(circle_at_32%_27%,hsl(var(--primary)/0.15),hsl(var(--accent)/0.1)_55%,transparent_72%)] opacity-85 blur-[10px]" style="animation: glow-pulse 4s ease-in-out infinite" />
        <Card class="relative w-[min(100%,420px)] rotate-[-2deg] gap-0 p-[22px] shadow-xl backdrop-blur-sm transition-transform duration-200 hover:rotate-0 hover:scale-[1.02]">
          <div class="flex items-center justify-between font-mono text-[10px] text-muted-foreground">
            <span class="flex items-center"><span class="mr-1.5 inline-block h-[7px] w-[7px] rounded-full bg-success shadow-[0_0_0_3px_hsl(var(--success)/0.2)]" style="animation: pulse-dot 2s ease-in-out infinite" />ROUTER ONLINE</span>
            <code>POST /v1/chat/completions</code>
          </div>
          <div class="my-4 flex items-center gap-2.5 rounded-[var(--radius-lg)] bg-muted p-3">
            <Bot :size="18" />
            <strong class="text-foreground">kimi-k3</strong>
            <span class="text-muted-foreground">{{ t('routingActive') }}</span>
          </div>
          <div class="flex items-center gap-2.5 border-b border-border py-2.5">
            <span class="h-2 w-2 flex-none rounded-full bg-success shadow-[0_0_0_3px_hsl(var(--success)/0.2)]" />
            <div>
              <b class="block text-[12px] text-foreground">{{ t('openaiMainRoute') }}</b>
              <small class="text-[10px] text-muted-foreground">{{ t('priorityLabel') }} P10 &middot; 42ms</small>
            </div>
            <span class="grid h-5 w-5 place-items-center rounded-full bg-success/10 text-[10px] font-bold text-success">✓</span>
          </div>
          <div class="flex items-center gap-2.5 py-2.5 opacity-60">
            <span class="h-2 w-2 flex-none rounded-full bg-input" />
            <div>
              <b class="block text-[12px] text-foreground">{{ t('backupChannelText') }}</b>
              <small class="text-[10px] text-muted-foreground">{{ t('waitingForFailover') }}</small>
            </div>
          </div>
          <div class="my-4 -mx-[22px] -mb-[22px] flex items-center gap-2 border-t border-border bg-muted/50 px-[22px] py-3.5 text-[11px]">
            <span>{{ t('successRate') }}</span>
            <strong class="font-mono text-foreground">99.98%</strong>
            <span class="h-3.5 w-px bg-border" />
            <span>{{ t('avgLatency') }}</span>
            <strong class="font-mono text-foreground">186ms</strong>
          </div>
        </Card>
      </div>
    </section>

    <section id="features" class="mx-auto my-20 w-[min(1160px,calc(100%-48px))] border-t border-border py-16">
      <div class="mb-10">
        <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">BUILT FOR CONTROL</p>
        <h2 class="m-0 text-[clamp(28px,3vw,42px)] font-bold tracking-[-0.06em] text-foreground">{{ t('featuresTagline') }}</h2>
      </div>
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-3">
        <Card class="min-h-[200px] gap-0 bg-[linear-gradient(145deg,hsl(var(--card)),hsl(var(--muted)/0.3))] p-6 transition-all duration-150 hover:-translate-y-1 hover:border-ring/30 hover:shadow-lg">
          <span class="mb-4 inline-block rounded-[var(--radius-lg)] bg-primary/10 px-2 py-1 font-mono text-[10px] font-bold text-primary">01</span>
          <RadioTower :size="21" class="mb-3 block text-primary" />
          <h3 class="mb-2 text-base font-semibold text-foreground">{{ t('smartRouting') }}</h3>
          <p class="text-[13px] leading-relaxed text-muted-foreground">{{ t('smartRoutingDesc') }}</p>
        </Card>
        <Card class="min-h-[200px] gap-0 bg-[linear-gradient(145deg,hsl(var(--card)),hsl(var(--muted)/0.3))] p-6 transition-all duration-150 hover:-translate-y-1 hover:border-ring/30 hover:shadow-lg">
          <span class="mb-4 inline-block rounded-[var(--radius-lg)] bg-primary/10 px-2 py-1 font-mono text-[10px] font-bold text-primary">02</span>
          <Activity :size="21" class="mb-3 block text-primary" />
          <h3 class="mb-2 text-base font-semibold text-foreground">{{ t('fullObservability') }}</h3>
          <p class="text-[13px] leading-relaxed text-muted-foreground">{{ t('fullObservabilityDesc') }}</p>
        </Card>
        <Card class="min-h-[200px] gap-0 bg-[linear-gradient(145deg,hsl(var(--card)),hsl(var(--muted)/0.3))] p-6 transition-all duration-150 hover:-translate-y-1 hover:border-ring/30 hover:shadow-lg">
          <span class="mb-4 inline-block rounded-[var(--radius-lg)] bg-primary/10 px-2 py-1 font-mono text-[10px] font-bold text-primary">03</span>
          <WalletCards :size="21" class="mb-3 block text-primary" />
          <h3 class="mb-2 text-base font-semibold text-foreground">{{ t('usageAndCost') }}</h3>
          <p class="text-[13px] leading-relaxed text-muted-foreground">{{ t('usageAndCostDesc') }}</p>
        </Card>
      </div>
    </section>

    <section id="quickstart" class="mx-auto mb-20 w-[min(1160px,calc(100%-48px))] border-t border-border py-16">
      <p class="mb-2 text-[10px] font-bold uppercase tracking-[0.14em] text-muted-foreground">ONE ENDPOINT</p>
      <h2 class="mb-4 text-[clamp(28px,3vw,42px)] font-bold tracking-[-0.06em] text-foreground">{{ t('quickstartTagline') }}</h2>
      <p class="max-w-[640px] text-[15px] leading-relaxed text-muted-foreground">{{ t('quickstartDesc') }}</p>
      <pre class="mt-6 max-w-[720px] overflow-x-auto rounded-[var(--radius-xl)] border border-primary/30 bg-[linear-gradient(140deg,#111827,#172033)] p-6 font-mono text-[12px] leading-[1.8] text-[#e5e7eb] shadow-xl"><span>// Using OpenAI SDK</span>
<span><b class="text-[#93c5fd]">const</b> client = <b class="text-[#93c5fd]">new</b> OpenAI({</span>
<span>  apiKey: <i class="text-[#6ee7b7]">'sk-xh-your-key'</i>,</span>
<span>  baseURL: <i class="text-[#6ee7b7]">'http://localhost:8080/v1'</i></span>
<span>})</span>
<span>&nbsp;</span>
<span><b class="text-[#93c5fd]">await</b> client.chat.completions.create({</span>
<span>  model: <i class="text-[#6ee7b7]">'kimi-k3'</i>,</span>
<span>  messages: [{ role: <i class="text-[#6ee7b7]">'user'</i>, content: <i class="text-[#6ee7b7]">'你好'</i> }]</span>
<span>})</span></pre>
    </section>

    <footer class="mx-auto flex w-[min(1160px,calc(100%-48px))] items-center justify-between border-t border-border py-6 text-xs text-muted-foreground">
      <span>© 2026 Xinghai Router</span>
      <span>{{ t('footerSlogan') }}</span>
      <span class="flex gap-3.5">
        <Button variant="link" as-child><a href="/terms">{{ t('termsShort') }}</a></Button>
        <Button variant="link" as-child><a href="/privacy">{{ t('privacyShort') }}</a></Button>
      </span>
    </footer>
  </main>
</template>