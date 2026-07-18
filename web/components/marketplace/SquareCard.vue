<script setup lang="ts">
import { computed, ref } from 'vue'
import { Check, ChevronRight, Copy } from 'lucide-vue-next'
import { effectivePrice, formatSquarePrice, getDisplayGroup, vendorIconUrl, type SquareModel, type TokenUnit } from '~/src/marketplace'

const props = defineProps<{
  model: SquareModel
  tokenUnit: TokenUnit
  selectedGroup: string
}>()
const emit = defineEmits<{ open: [model: SquareModel] }>()
const { t } = useI18n()

const iconError = ref(false)
const copied = ref(false)
const unitLabel = computed(() => (props.tokenUnit === 'K' ? '1K' : '1M'))
const primaryGroup = computed(() => getDisplayGroup(props.model, props.selectedGroup))

const inputPrice = computed(() => formatSquarePrice(effectivePrice(props.model, 'input', props.selectedGroup), props.tokenUnit))
const outputPrice = computed(() => formatSquarePrice(effectivePrice(props.model, 'output', props.selectedGroup), props.tokenUnit))
const cachePrice = computed(() => formatSquarePrice(effectivePrice(props.model, 'cache', props.selectedGroup), props.tokenUnit))

const hiddenCount = computed(() => Math.max(props.model.groups.length - 1, 0))

async function copyName(event: MouseEvent) {
  event.stopPropagation()
  try {
    await navigator.clipboard.writeText(props.model.model)
  } catch {
    const textarea = document.createElement('textarea')
    textarea.value = props.model.model
    document.body.append(textarea)
    textarea.select()
    document.execCommand('copy')
    textarea.remove()
  }
  copied.value = true
  window.setTimeout(() => { copied.value = false }, 1500)
}
</script>

<template>
  <div class="msq-card">
    <div class="msq-card-top">
      <div class="msq-card-id">
        <div class="msq-iconbox">
          <img v-if="props.model.vendor_slug && !iconError" :src="vendorIconUrl(props.model.vendor_slug)" :alt="props.model.vendor_name" loading="lazy" @error="iconError = true" />
          <span v-else>{{ props.model.model.slice(0, 1).toUpperCase() }}</span>
        </div>
        <div class="msq-card-title">
          <h3>{{ props.model.model }}</h3>
          <div class="msq-card-prices">
            <span class="msq-price-item">
              {{ t('inputLabel') }}
              <b v-if="inputPrice">{{ inputPrice }}</b>
              <b v-else class="msq-price-pending">{{ t('pendingConfig') }}</b>
            </span>
            <span class="msq-price-item">
              {{ t('outputLabel') }}
              <b v-if="outputPrice">{{ outputPrice }}</b>
              <b v-else class="msq-price-pending">{{ t('pendingConfig') }}</b>
            </span>
            <span v-if="cachePrice" class="msq-price-item">
              {{ t('msCached') }}
              <b>{{ cachePrice }}</b>
            </span>
          </div>
        </div>
      </div>
      <div class="msq-card-actions">
        <button type="button" class="msq-details-btn" @click="emit('open', props.model)">
          {{ t('msDetails') }}<ChevronRight :size="13" />
        </button>
        <button type="button" class="msq-copy-btn" :title="copied ? t('msCopied') : t('msCopy')" @click="copyName">
          <Check v-if="copied" :size="13" />
          <Copy v-else :size="13" />
        </button>
      </div>
    </div>

    <p class="msq-card-desc">{{ props.model.vendor_name || t('msNoDesc') }}</p>

    <div class="msq-card-footer">
      <div class="msq-card-footer-left">
        <span v-if="primaryGroup" class="msq-group-name">{{ primaryGroup.name }}</span>
        <span class="msq-billing-badge">{{ t('msTokenBased') }}</span>
      </div>
      <div class="msq-card-footer-bottom">
        <span class="msq-unit">{{ unitLabel }}</span>
        <span v-if="hiddenCount > 0" class="msq-hidden-count">+{{ hiddenCount }}</span>
      </div>
    </div>
  </div>
</template>
