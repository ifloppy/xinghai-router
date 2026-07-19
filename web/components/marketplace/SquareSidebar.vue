<script setup lang="ts">
import { computed, ref } from 'vue'
import { ChevronDown, RotateCcw } from 'lucide-vue-next'
import type { CatalogGroup } from '~/src/api'
import { FILTER_ALL, extractVendors, formatRatio, vendorIconUrl, type SquareModel } from '~/src/marketplace'

const props = defineProps<{
  models: SquareModel[]
  groups: CatalogGroup[]
  vendorFilter: string
  groupFilter: string
  quotaTypeFilter: string
  hasActiveFilters: boolean
  bare?: boolean
}>()
const emit = defineEmits<{
  'update:vendorFilter': [value: string]
  'update:groupFilter': [value: string]
  'update:quotaTypeFilter': [value: string]
  clear: []
}>()
const { t } = useI18n()

const vendors = computed(() => extractVendors(props.models))
const iconErrors = ref(new Set<string>())
function iconFailed(slug: string) {
  iconErrors.value = new Set(iconErrors.value).add(slug)
}

const collapsed = ref<Record<string, boolean>>({})
function toggleSection(key: string) {
  collapsed.value = { ...collapsed.value, [key]: !collapsed.value[key] }
}

const quotaOptions = computed(() => [
  { value: 'all', label: t('msAllModels'), count: props.models.length },
  { value: 'token', label: t('msTokenBased'), count: props.models.length },
])
</script>

<template>
  <aside :class="['msq-sidebar', { bare: props.bare }]">
    <div class="msq-sidebar-head">
      <div>
        <h2>{{ t('msFilter') }}</h2>
        <p>{{ t('msFilterDesc') }}</p>
      </div>
      <button type="button" class="msq-reset" :disabled="!props.hasActiveFilters" @click="emit('clear')">
        <RotateCcw :size="13" />{{ t('msReset') }}
      </button>
    </div>
    <span v-if="props.hasActiveFilters" class="msq-active-badge">{{ t('msFiltersActive') }}</span>

    <div class="msq-sections">
      <section class="msq-section" :class="{ collapsed: collapsed.group }">
        <button type="button" class="msq-section-head" @click="toggleSection('group')">
          <span>{{ t('msGroups') }}</span>
          <ChevronDown :size="15" class="msq-chevron" />
        </button>
        <div class="msq-chips">
          <button type="button" :class="['msq-chip', { active: props.groupFilter === FILTER_ALL }]" @click="emit('update:groupFilter', FILTER_ALL)">
            <span>{{ t('msAllGroups') }}</span>
          </button>
          <button v-for="group in props.groups" :key="group.id" type="button" :class="['msq-chip', { active: props.groupFilter === group.id }]" :title="group.name" @click="emit('update:groupFilter', group.id)">
            <span class="msq-chip-label">{{ group.name }}</span>
            <span class="msq-chip-meta">{{ formatRatio(Number(group.multiplier)) }}</span>
          </button>
        </div>
      </section>

      <section class="msq-section" :class="{ collapsed: collapsed.vendor }">
        <button type="button" class="msq-section-head" @click="toggleSection('vendor')">
          <span>{{ t('msVendors') }}</span>
          <ChevronDown :size="15" class="msq-chevron" />
        </button>
        <div class="msq-chips">
          <button type="button" :class="['msq-chip', { active: props.vendorFilter === FILTER_ALL }]" @click="emit('update:vendorFilter', FILTER_ALL)">
            <span>{{ t('msAllVendors') }}</span>
            <span class="msq-chip-meta">{{ props.models.length }}</span>
          </button>
          <button v-for="vendor in vendors" :key="vendor.name" type="button" :class="['msq-chip', { active: props.vendorFilter === vendor.name }]" :title="vendor.name" @click="emit('update:vendorFilter', vendor.name)">
            <img v-if="vendor.slug && !iconErrors.has(vendor.slug)" class="msq-chip-icon" :src="vendorIconUrl(vendor.slug)" :alt="vendor.name" loading="lazy" @error="iconFailed(vendor.slug)" >
            <span class="msq-chip-label">{{ vendor.name }}</span>
            <span class="msq-chip-meta">{{ vendor.count }}</span>
          </button>
        </div>
      </section>

      <section class="msq-section" :class="{ collapsed: collapsed.quota }">
        <button type="button" class="msq-section-head" @click="toggleSection('quota')">
          <span>{{ t('msPricingType') }}</span>
          <ChevronDown :size="15" class="msq-chevron" />
        </button>
        <div class="msq-chips">
          <button v-for="option in quotaOptions" :key="option.value" type="button" :class="['msq-chip', { active: props.quotaTypeFilter === option.value }]" @click="emit('update:quotaTypeFilter', option.value)">
            <span>{{ option.label }}</span>
            <span class="msq-chip-meta">{{ option.count }}</span>
          </button>
        </div>
      </section>
    </div>
  </aside>
</template>
