import { readonly, ref } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'
export type ThemeColor = 'neutral' | 'blue' | 'green' | 'orange' | 'rose' | 'violet'
export type ThemeRadius = 'none' | 'small' | 'medium' | 'large'
export type ThemePreset = 'custom' | 'a-site'

const STORAGE_KEY = 'xinghai-router-theme-config'
const mode = ref<ThemeMode>('system')
const color = ref<ThemeColor>('neutral')
const radius = ref<ThemeRadius>('medium')
const preset = ref<ThemePreset>('custom')
let mediaQuery: MediaQueryList | undefined
let initialized = false

function resolvedMode() {
  return mode.value === 'system' ? mediaQuery?.matches ? 'dark' : 'light' : mode.value
}

function applyTheme() {
  if (!import.meta.client) return
  const root = document.documentElement
  root.dataset.theme = resolvedMode()
  root.dataset.themeMode = mode.value
  root.dataset.themeColor = color.value
  root.dataset.themeRadius = radius.value
  root.dataset.themePreset = preset.value
  root.style.colorScheme = resolvedMode()
}

function persist() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify({ mode: mode.value, color: color.value, radius: radius.value, preset: preset.value }))
}

function initializeTheme() {
  if (!import.meta.client || initialized) return
  initialized = true
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const legacyMode = localStorage.getItem('xinghai-router-theme')
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
    if (['light', 'dark', 'system'].includes(saved.mode)) mode.value = saved.mode
    else if (legacyMode === 'light' || legacyMode === 'dark') mode.value = legacyMode
    if (['neutral', 'blue', 'green', 'orange', 'rose', 'violet'].includes(saved.color)) color.value = saved.color
    if (['none', 'small', 'medium', 'large'].includes(saved.radius)) radius.value = saved.radius
    if (['custom', 'a-site'].includes(saved.preset)) preset.value = saved.preset
  } catch {
    if (legacyMode === 'light' || legacyMode === 'dark') mode.value = legacyMode
  }
  mediaQuery.addEventListener('change', applyTheme)
  applyTheme()
}

function setMode(value: ThemeMode) { mode.value = value; applyTheme(); persist() }
function setColor(value: ThemeColor) { color.value = value; preset.value = 'custom'; applyTheme(); persist() }
function setRadius(value: ThemeRadius) { radius.value = value; preset.value = 'custom'; applyTheme(); persist() }
function setPreset(value: ThemePreset) {
  preset.value = value
  if (value === 'a-site') {
    mode.value = 'light'
    color.value = 'orange'
    radius.value = 'large'
  }
  applyTheme()
  persist()
}
function resetTheme() { mode.value = 'system'; color.value = 'neutral'; radius.value = 'medium'; preset.value = 'custom'; applyTheme(); persist() }

export function useTheme() {
  return { mode: readonly(mode), color: readonly(color), radius: readonly(radius), preset: readonly(preset), initializeTheme, setMode, setColor, setRadius, setPreset, resetTheme }
}
