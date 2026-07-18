<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps<{ page: number; totalPages: number }>()
const emit = defineEmits<{ 'update:page': [value: number] }>()
const { t } = useI18n()
</script>

<template>
  <div v-if="props.totalPages > 1" class="msq-pagination">
    <p>{{ t('msPageA') }} {{ props.page }} {{ t('msPageB') }} {{ props.totalPages }} {{ t('msPageC') }}</p>
    <div>
      <button type="button" :disabled="props.page <= 1" @click="emit('update:page', props.page - 1)">
        <ChevronLeft :size="15" />{{ t('msPrev') }}
      </button>
      <button type="button" :disabled="props.page >= props.totalPages" @click="emit('update:page', props.page + 1)">
        {{ t('msNext') }}<ChevronRight :size="15" />
      </button>
    </div>
  </div>
</template>
