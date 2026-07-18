<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { ArrowUpDown, Check, Filter, Grid2X2, Table2 } from 'lucide-vue-next'
import type { SortOption, TokenUnit, ViewMode } from '~/src/marketplace'

const props = defineProps<{
  filteredCount: number
  totalCount: number
  hasActiveFilters: boolean
  activeFilterCount: number
  sortBy: SortOption
  tokenUnit: TokenUnit
  viewMode: ViewMode
}>()
const emit = defineEmits<{
  'update:sortBy': [value: SortOption]
  'update:tokenUnit': [value: TokenUnit]
  'update:viewMode': [value: ViewMode]
  openFilters: []
}>()
const { t } = useI18n()

const sortOpen = ref(false)
const sortRef = ref<HTMLElement | null>(null)
const sortLabels = computed<Record<SortOption, string>>(() => ({
  name: t('msSortName'),
  'price-low': t('msSortPriceLow'),
  'price-high': t('msSortPriceHigh'),
}))

function chooseSort(value: SortOption) {
  emit('update:sortBy', value)
  sortOpen.value = false
}

function handleOutside(event: MouseEvent) {
  if (sortRef.value && !sortRef.value.contains(event.target as Node)) sortOpen.value = false
}
onMounted(() => document.addEventListener('click', handleOutside))
onBeforeUnmount(() => document.removeEventListener('click', handleOutside))
</script>

<template>
  <div class="msq-toolbar">
    <div class="msq-toolbar-left">
      <button type="button" class="msq-filter-btn" @click="emit('openFilters')">
        <Filter :size="15" />{{ t('msFilter') }}
        <span v-if="props.activeFilterCount > 0" class="msq-filter-count">{{ props.activeFilterCount }}</span>
      </button>
      <div class="msq-count">
        <strong>{{ props.filteredCount.toLocaleString() }}</strong>
        <span>{{ t('msModels') }}</span>
        <small v-if="props.hasActiveFilters && props.totalCount">/ {{ props.totalCount.toLocaleString() }}</small>
      </div>
    </div>

    <div class="msq-toolbar-right">
      <div class="msq-segmented msq-hide-mobile" role="group" :aria-label="t('msTokenUnit')">
        <button type="button" :class="{ active: props.tokenUnit === 'M' }" :aria-pressed="props.tokenUnit === 'M'" @click="emit('update:tokenUnit', 'M')">/1M</button>
        <button type="button" :class="{ active: props.tokenUnit === 'K' }" :aria-pressed="props.tokenUnit === 'K'" @click="emit('update:tokenUnit', 'K')">/1K</button>
      </div>

      <div ref="sortRef" class="msq-dropdown">
        <button type="button" class="msq-sort-btn" @click.stop="sortOpen = !sortOpen">
          <ArrowUpDown :size="13" /><span>{{ sortLabels[props.sortBy] || t('msSort') }}</span>
        </button>
        <div v-if="sortOpen" class="msq-dropdown-menu">
          <button v-for="(label, value) in sortLabels" :key="value" type="button" @click="chooseSort(value as SortOption)">
            <Check :size="14" :style="{ opacity: props.sortBy === value ? 1 : 0 }" />{{ label }}
          </button>
        </div>
      </div>

      <div class="msq-segmented" role="group" :aria-label="t('msViewMode')">
        <button type="button" class="msq-icon-btn" :class="{ active: props.viewMode === 'card' }" :title="t('msCardView')" :aria-pressed="props.viewMode === 'card'" @click="emit('update:viewMode', 'card')">
          <Grid2X2 :size="13" />
        </button>
        <button type="button" class="msq-icon-btn" :class="{ active: props.viewMode === 'table' }" :title="t('msTableView')" :aria-pressed="props.viewMode === 'table'" @click="emit('update:viewMode', 'table')">
          <Table2 :size="13" />
        </button>
      </div>
    </div>
  </div>
</template>
