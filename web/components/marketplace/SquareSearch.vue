<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { Search, X } from 'lucide-vue-next'

const props = defineProps<{ modelValue: string; placeholder?: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string]; clear: [] }>()
const { t } = useI18n()
const inputRef = ref<HTMLInputElement | null>(null)

function handleKeydown(event: KeyboardEvent) {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    inputRef.value?.focus()
  }
  if (event.key === 'Escape' && document.activeElement === inputRef.value) inputRef.value?.blur()
}

onMounted(() => document.addEventListener('keydown', handleKeydown))
onBeforeUnmount(() => document.removeEventListener('keydown', handleKeydown))
</script>

<template>
  <div class="msq-search">
    <Search :size="16" class="msq-search-icon" />
    <input
      ref="inputRef"
      type="text"
      :value="props.modelValue"
      :placeholder="props.placeholder || t('msSearchPlaceholder')"
      :aria-label="t('msSearchPlaceholder')"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    >
    <div class="msq-search-side">
      <button v-if="props.modelValue" type="button" class="msq-search-clear" :aria-label="t('msClearSearch')" @click="emit('clear')">
        <X :size="15" />
      </button>
      <kbd v-else>⌘K</kbd>
    </div>
  </div>
</template>
