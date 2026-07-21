<script setup lang="ts">
import { ref, watch, onBeforeUnmount, onMounted } from 'vue'
import { CircleAlert, Copy } from 'lucide-vue-next'
import { useConsoleStore } from '~/composables/useConsoleStore'

const store = useConsoleStore()
const { t } = store

const errorAlert = ref<HTMLElement | null>(null)
const errorHovered = ref(false)
const errorSelected = ref(false)
let errorTimer: ReturnType<typeof setTimeout> | undefined

function clearErrorTimer() {
  if (errorTimer) window.clearTimeout(errorTimer)
  errorTimer = undefined
}

function scheduleErrorDismissal() {
  clearErrorTimer()
  if (!store.error.value || errorHovered.value || errorSelected.value) return
  errorTimer = window.setTimeout(() => { store.error.value = '' }, 5000)
}

function updateErrorSelection() {
  const selection = window.getSelection()
  const anchor = selection?.anchorNode
  const focus = selection?.focusNode
  errorSelected.value = Boolean(selection?.toString() && anchor && focus && errorAlert.value?.contains(anchor) && errorAlert.value?.contains(focus))
  scheduleErrorDismissal()
}

function lockError() {
  errorHovered.value = true
  clearErrorTimer()
}

function releaseError() {
  errorHovered.value = false
  updateErrorSelection()
  if (!errorSelected.value) store.error.value = ''
}

async function copyError() {
  if (!store.error.value) return
  if (navigator.clipboard) {
    await navigator.clipboard.writeText(store.error.value)
    return
  }
  const textarea = document.createElement('textarea')
  textarea.value = store.error.value
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.append(textarea)
  textarea.select()
  document.execCommand('copy')
  textarea.remove()
}

watch(store.error, () => {
  errorSelected.value = false
  scheduleErrorDismissal()
})

onMounted(() => {
  document.addEventListener('selectionchange', updateErrorSelection)
})

onBeforeUnmount(() => {
  clearErrorTimer()
  document.removeEventListener('selectionchange', updateErrorSelection)
})
</script>

<template>
  <Transition name="error-alert">
    <div v-if="store.error.value" ref="errorAlert" class="fixed left-1/2 top-4 z-50 flex max-w-[min(92vw,720px)] -translate-x-1/2 items-start gap-2 rounded-md border border-destructive/40 bg-destructive/10 px-3 py-2 text-sm text-destructive shadow-md backdrop-blur" role="alert" tabindex="0" :title="t('clickToCopyError')" @mouseenter="lockError" @mouseleave="releaseError" @click="copyError" @keydown.enter.prevent="copyError" @keydown.space.prevent="copyError">
      <CircleAlert :size="17" class="mt-px shrink-0" /><span class="flex-1 overflow-wrap-anywhere">{{ store.error.value }}</span><Copy :size="14" aria-hidden="true" class="mt-0.5 shrink-0" />
    </div>
  </Transition>
</template>