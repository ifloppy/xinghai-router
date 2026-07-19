<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Search } from 'lucide-vue-next'
import type { CatalogGroup, CatalogModel } from '~/src/api'
import { FILTER_ALL, PAGE_SIZE, filterAndSort, toSquareModel, type SortOption, type SquareModel, type TokenUnit, type ViewMode } from '~/src/marketplace'
import SquareSearch from './SquareSearch.vue'
import SquareSidebar from './SquareSidebar.vue'
import SquareToolbar from './SquareToolbar.vue'
import SquareCard from './SquareCard.vue'
import SquareTable from './SquareTable.vue'
import SquarePagination from './SquarePagination.vue'
import SquareDrawer from './SquareDrawer.vue'

const props = defineProps<{
  catalog: CatalogModel[]
  groups: CatalogGroup[]
  loaded: boolean
}>()
const { t } = useI18n()

const models = computed(() => props.catalog.map(toSquareModel))

const searchInput = ref('')
const sortBy = ref<SortOption>('name')
const vendorFilter = ref(FILTER_ALL)
const groupFilter = ref(FILTER_ALL)
const quotaTypeFilter = ref('all')
const tokenUnit = ref<TokenUnit>('M')
const viewMode = ref<ViewMode>('card')
const page = ref(1)
const mobileFiltersOpen = ref(false)
const selectedModel = ref<SquareModel | null>(null)

const filteredModels = computed(() =>
  filterAndSort(models.value, {
    search: searchInput.value,
    vendor: vendorFilter.value,
    group: groupFilter.value,
    quotaType: quotaTypeFilter.value,
    sortBy: sortBy.value,
  }),
)

const hasActiveFilters = computed(() => vendorFilter.value !== FILTER_ALL || groupFilter.value !== FILTER_ALL || quotaTypeFilter.value !== 'all')
const activeFilterCount = computed(() => (vendorFilter.value !== FILTER_ALL ? 1 : 0) + (groupFilter.value !== FILTER_ALL ? 1 : 0) + (quotaTypeFilter.value !== 'all' ? 1 : 0))

const totalPages = computed(() => Math.max(1, Math.ceil(filteredModels.value.length / PAGE_SIZE)))
const currentPage = computed(() => Math.min(page.value, totalPages.value))
const pagedModels = computed(() => filteredModels.value.slice((currentPage.value - 1) * PAGE_SIZE, currentPage.value * PAGE_SIZE))

watch([searchInput, vendorFilter, groupFilter, quotaTypeFilter, sortBy], () => { page.value = 1 })

function clearFilters() {
  vendorFilter.value = FILTER_ALL
  groupFilter.value = FILTER_ALL
  quotaTypeFilter.value = 'all'
}
function clearAll() {
  clearFilters()
  searchInput.value = ''
}

const origin = computed(() => (import.meta.client ? window.location.origin : ''))
const skeletonCards = Array.from({ length: 6 }, (_, index) => index)
</script>

<template>
  <div class="msq-root">
    <div class="msq-gradient" aria-hidden="true"/>

    <div class="msq-container">
      <header class="msq-header">
        <h1>{{ t('marketplace') }}</h1>
        <p class="msq-sub1">{{ t('msSubtitle1a') }}{{ props.catalog.length }}{{ t('msSubtitle1b') }}</p>
        <p class="msq-sub2">{{ t('msSubtitle2') }}</p>
        <SquareSearch v-model="searchInput" class="msq-header-search" @clear="searchInput = ''" />
      </header>

      <div v-if="!props.loaded" class="msq-skeleton-wrap">
        <div class="msq-skeleton-toolbar"/>
        <div class="msq-grid">
          <div v-for="item in skeletonCards" :key="item" class="msq-skeleton-card">
            <div class="row"><i/><div><b/><span/></div></div>
            <p/><p class="short"/>
          </div>
        </div>
      </div>

      <div v-else class="msq-layout">
        <SquareSidebar
          v-model:vendor-filter="vendorFilter"
          v-model:group-filter="groupFilter"
          v-model:quota-type-filter="quotaTypeFilter"
          :models="models"
          :groups="props.groups"
          :has-active-filters="hasActiveFilters"
          class="msq-sidebar-desktop"
          @clear="clearFilters"
        />

        <main class="msq-main">
          <SquareToolbar
            v-model:sort-by="sortBy"
            v-model:token-unit="tokenUnit"
            v-model:view-mode="viewMode"
            :filtered-count="filteredModels.length"
            :total-count="models.length"
            :has-active-filters="hasActiveFilters"
            :active-filter-count="activeFilterCount"
            @open-filters="mobileFiltersOpen = true"
          />

          <div v-if="!filteredModels.length" class="msq-empty">
            <Search :size="38" class="msq-empty-icon" />
            <h3>{{ t('msEmptyTitle') }}</h3>
            <p>{{ searchInput.trim() ? `${t('msEmptySearch1')}${searchInput}${t('msEmptySearch2')}` : t('msEmptyFilter') }}</p>
            <button v-if="hasActiveFilters || searchInput.trim()" type="button" class="msq-empty-clear" @click="clearAll">
              {{ t('msClearFilters') }}
            </button>
          </div>

          <template v-else-if="viewMode === 'card'">
            <div class="msq-grid">
              <SquareCard
                v-for="model in pagedModels"
                :key="model.model"
                :model="model"
                :token-unit="tokenUnit"
                :selected-group="groupFilter"
                @open="selectedModel = $event"
              />
            </div>
            <SquarePagination v-model:page="page" :total-pages="totalPages" />
          </template>

          <SquareTable
            v-else
            v-model:page="page"
            :models="pagedModels"
            :token-unit="tokenUnit"
            :selected-group="groupFilter"
            :total-pages="totalPages"
            @open="selectedModel = $event"
          />
        </main>
      </div>
    </div>

    <Teleport to="body">
      <template v-if="mobileFiltersOpen">
        <div class="msq-sheet-backdrop" @click="mobileFiltersOpen = false"/>
        <div class="msq-sheet">
          <div class="msq-sheet-head">
            <h2>{{ t('msFilter') }}</h2>
            <p>{{ t('msFilterDesc') }}</p>
          </div>
          <SquareSidebar
            v-model:vendor-filter="vendorFilter"
            v-model:group-filter="groupFilter"
            v-model:quota-type-filter="quotaTypeFilter"
            :models="models"
            :groups="props.groups"
            :has-active-filters="hasActiveFilters"
            bare
            @clear="clearFilters"
          />
        </div>
      </template>
    </Teleport>

    <SquareDrawer
      v-if="selectedModel"
      :model="selectedModel"
      :token-unit="tokenUnit"
      :origin="origin"
      @close="selectedModel = null"
    />
  </div>
</template>
