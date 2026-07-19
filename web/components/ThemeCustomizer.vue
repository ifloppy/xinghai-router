<script setup lang="ts">
import { Check, Monitor, Moon, Paintbrush, RotateCcw, Sun, X } from 'lucide-vue-next'
import { onBeforeUnmount, ref } from 'vue'
import type { Locale } from '~/composables/useI18n'
import type { ThemeColor, ThemeMode, ThemeRadius } from '~/composables/useTheme'

const props = defineProps<{ locale: Locale }>()
const open = ref(false)
const panel = ref<HTMLElement | null>(null)
const { mode, color, radius, preset, setMode, setColor, setRadius, setPreset, resetTheme } = useTheme()
const modes: { value: ThemeMode; zh: string; en: string; icon: typeof Sun }[] = [
  { value: 'light', zh: '浅色', en: 'Light', icon: Sun }, { value: 'dark', zh: '深色', en: 'Dark', icon: Moon }, { value: 'system', zh: '跟随系统', en: 'System', icon: Monitor },
]
const colors: { value: ThemeColor; zh: string; en: string }[] = [
  { value: 'neutral', zh: '中性', en: 'Neutral' }, { value: 'blue', zh: '蓝色', en: 'Blue' }, { value: 'green', zh: '绿色', en: 'Green' }, { value: 'orange', zh: '橙色', en: 'Orange' }, { value: 'rose', zh: '玫红', en: 'Rose' }, { value: 'violet', zh: '紫色', en: 'Violet' },
]
const radii: { value: ThemeRadius; zh: string; en: string }[] = [
  { value: 'none', zh: '无', en: 'None' }, { value: 'small', zh: '小', en: 'Small' }, { value: 'medium', zh: '中', en: 'Medium' }, { value: 'large', zh: '大', en: 'Large' },
]
const label = (item: { zh: string; en: string }) => props.locale === 'zh-CN' ? item.zh : item.en
function closeOnOutsideClick(event: PointerEvent) { if (open.value && !panel.value?.contains(event.target as Node)) open.value = false }
if (import.meta.client) document.addEventListener('pointerdown', closeOnOutsideClick)
onBeforeUnmount(() => document.removeEventListener('pointerdown', closeOnOutsideClick))
</script>

<template>
  <div ref="panel" class="theme-customizer">
    <button class="theme-toggle" :aria-expanded="open" :aria-label="locale === 'zh-CN' ? '配置主题' : 'Customize theme'" :title="locale === 'zh-CN' ? '配置主题' : 'Customize theme'" @click="open = !open"><Paintbrush :size="16" /></button>
    <Transition name="theme-panel">
      <section v-if="open" class="theme-panel" role="dialog" :aria-label="locale === 'zh-CN' ? '主题配置' : 'Theme settings'">
        <header><div><h2>{{ locale === 'zh-CN' ? '自定义主题' : 'Customize' }}</h2><p>{{ locale === 'zh-CN' ? '实时预览界面外观' : 'Preview your interface in real time.' }}</p></div><button :aria-label="locale === 'zh-CN' ? '关闭' : 'Close'" @click="open = false"><X :size="16" /></button></header>
        <div class="theme-section"><span>{{ locale === 'zh-CN' ? '预设' : 'Preset' }}</span><button class="theme-preset-card" :class="{ active: preset === 'a-site' }" @click="setPreset('a-site')"><i/><span><b>Claude 官网</b><small>{{ locale === 'zh-CN' ? 'Anthropic 奶油白与 Claude 橙' : 'Anthropic cream with Claude orange' }}</small></span><Check v-if="preset === 'a-site'" :size="15" /></button></div>
        <div class="theme-section"><span>{{ locale === 'zh-CN' ? '模式' : 'Mode' }}</span><div class="theme-mode-grid"><button v-for="item in modes" :key="item.value" :class="{ active: mode === item.value }" @click="setMode(item.value)"><component :is="item.icon" :size="15" />{{ label(item) }}</button></div></div>
        <div class="theme-section"><span>{{ locale === 'zh-CN' ? '主色' : 'Color' }}</span><div class="theme-color-grid"><button v-for="item in colors" :key="item.value" :class="['theme-color', `theme-color-${item.value}`, { active: color === item.value }]" @click="setColor(item.value)"><i/>{{ label(item) }}<Check v-if="color === item.value" :size="14" /></button></div></div>
        <div class="theme-section"><span>{{ locale === 'zh-CN' ? '圆角' : 'Radius' }}</span><div class="theme-radius-grid"><button v-for="item in radii" :key="item.value" :class="{ active: radius === item.value }" @click="setRadius(item.value)">{{ label(item) }}</button></div></div>
        <button class="theme-reset" @click="resetTheme"><RotateCcw :size="14" />{{ locale === 'zh-CN' ? '恢复默认' : 'Reset theme' }}</button>
      </section>
    </Transition>
  </div>
</template>
