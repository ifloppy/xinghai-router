<script setup lang="ts">
import { computed, ref } from 'vue'
import { effectivePrice, formatSquarePrice, vendorColor, vendorIconUrl, type SquareModel, type TokenUnit } from '~/src/marketplace'
import SquarePagination from './SquarePagination.vue'

const props = defineProps<{
  models: SquareModel[]
  tokenUnit: TokenUnit
  selectedGroup: string
  page: number
  totalPages: number
}>()
const emit = defineEmits<{ open: [model: SquareModel]; 'update:page': [value: number] }>()
const { t } = useI18n()

const unitLabel = computed(() => (props.tokenUnit === 'K' ? '1K' : '1M'))
const iconErrors = ref(new Set<string>())
function iconFailed(slug: string) {
  iconErrors.value = new Set(iconErrors.value).add(slug)
}

function price(model: SquareModel, kind: 'input' | 'output' | 'cache') {
  return formatSquarePrice(effectivePrice(model, kind, props.selectedGroup), props.tokenUnit)
}
</script>

<template>
  <div class="msq-table-wrap">
    <table class="msq-table">
      <thead>
        <tr>
          <th>{{ t('msModel') }}</th>
          <th>{{ t('msType') }}</th>
          <th>{{ t('msPrice') }}</th>
          <th>{{ t('msCached') }}</th>
          <th>{{ t('msVendor') }}</th>
          <th>{{ t('msGroups') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="model in props.models" :key="model.model" @click="emit('open', model)">
          <td>
            <div class="msq-cell-model">
              <img v-if="model.vendor_slug && !iconErrors.has(model.vendor_slug)" :src="vendorIconUrl(model.vendor_slug)" :alt="model.vendor_name" loading="lazy" @error="iconFailed(model.vendor_slug)" />
              <span class="msq-mono">{{ model.model }}</span>
            </div>
          </td>
          <td><span class="msq-billing-badge">{{ t('msTokenBased') }}</span></td>
          <td>
            <div class="msq-cell-price">
              <span class="msq-mono">
                <template v-if="price(model, 'input')">{{ price(model, 'input') }}<i>/</i>{{ price(model, 'output') }}</template>
                <template v-else>{{ t('pendingConfig') }}</template>
              </span>
              <small>/ {{ unitLabel }} tokens</small>
            </div>
          </td>
          <td>
            <div v-if="price(model, 'cache')" class="msq-cell-price">
              <span class="msq-mono">{{ price(model, 'cache') }}</span>
              <small>/ {{ unitLabel }}</small>
            </div>
            <span v-else class="msq-dash">—</span>
          </td>
          <td>
            <span class="msq-vendor-badge" :style="{ background: vendorColor(model.vendor_name).bg, color: vendorColor(model.vendor_name).fg }">
              {{ model.vendor_name }}
            </span>
          </td>
          <td>
            <div class="msq-cell-groups">
              <span v-for="group in model.groups" :key="group.id" class="msq-group-badge">{{ group.name }}</span>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <SquarePagination :page="props.page" :total-pages="props.totalPages" @update:page="emit('update:page', $event)" />
  </div>
</template>
