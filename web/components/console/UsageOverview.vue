<script setup lang="ts">
import { useConsoleStore } from '~/composables/useConsoleStore'
import { Card, CardContent } from '@/components/ui/card'

const store = useConsoleStore()
const { t, personalTokens, personalCost, personalRequests, usageChart } = store
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-3">
    <Card>
      <CardContent class="pt-6">
        <span class="text-xs text-muted-foreground">{{ t('last7DaysTokens') }}</span>
        <div class="text-2xl font-semibold">{{ personalTokens.toLocaleString() }}</div>
        <p class="mt-1 text-xs text-muted-foreground">{{ t('inputOutputTotal') }}</p>
      </CardContent>
    </Card>
    <Card>
      <CardContent class="pt-6">
        <span class="text-xs text-muted-foreground">{{ t('last7DaysCost') }}</span>
        <div class="text-2xl font-semibold">{{ personalCost.toFixed(6) }}</div>
        <p class="mt-1 text-xs text-muted-foreground">{{ t('basedOnCurrentPricing') }}</p>
      </CardContent>
    </Card>
    <Card>
      <CardContent class="pt-6">
        <span class="text-xs text-muted-foreground">{{ t('callCount') }}</span>
        <div class="text-2xl font-semibold">{{ personalRequests }}</div>
        <p class="mt-1 text-xs text-muted-foreground">{{ t('recent100UsageRecords') }}</p>
      </CardContent>
    </Card>
  </div>

  <Card class="mt-4">
    <CardContent class="pt-6">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
        <div>
          <h2 class="text-sm font-semibold">{{ t('usageTrend') }}</h2>
          <p class="text-xs text-muted-foreground">{{ t('last7DaysTokenAndCost') }}</p>
        </div>
        <div class="flex items-center gap-4 text-xs text-muted-foreground">
          <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-green-600" />{{ t('tokenLabel') }}</span>
          <span class="flex items-center gap-1"><span class="h-2 w-2 rounded-full bg-orange-500" />{{ t('costLabel') }}</span>
        </div>
      </div>
      <div class="flex items-end justify-between gap-2">
        <div v-for="day in usageChart" :key="day.key" class="flex flex-1 flex-col items-center gap-2">
          <div class="flex h-40 w-full items-end justify-center gap-1">
            <span class="w-2 rounded-t bg-green-600/80 transition-all" :style="{ height: `${day.tokenHeight}%` }" :title="`${day.tokens.toLocaleString()} tokens`" />
            <span class="w-2 rounded-t bg-orange-500/80 transition-all" :style="{ height: `${day.costHeight}%` }" :title="`${t('costLabel')} ${day.cost.toFixed(6)}`" />
          </div>
          <b class="text-xs font-medium">{{ day.label }}</b>
          <small class="font-mono text-xs text-muted-foreground">{{ day.tokens ? day.tokens.toLocaleString() : '-' }}</small>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
