<script setup lang="ts">
import { useConsoleStore } from '~/composables/useConsoleStore'
import Empty from '~/components/console/Empty.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const store = useConsoleStore()
const { t, busy, subscriptionPlans, adminSubscriptions, extendForm, formatDate, extendSubscriptions } = store

const statusBadge = (status: string) => ({ pending: t('subscriptionPending'), active: t('subscriptionActive'), expired: t('subscriptionExpired'), cancelled: t('subscriptionCancelled') }[status] ?? status)
const statusVariant = (status: string) => ({ pending: 'outline', active: 'secondary', expired: 'destructive', cancelled: 'destructive' }[status] ?? 'secondary')
const selectClass = 'flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring'
</script>

<template>
  <section class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-lg font-semibold">{{ t('adminSubscriptions') }}</h2>
      <p class="text-sm text-muted-foreground">{{ t('adminSubscriptionsDesc') }}</p>
    </div>
  </section>
  <section class="mt-4 rounded-lg border border-border bg-card p-4">
    <h3 class="mb-3 text-sm font-medium">{{ t('batchExtend') }}</h3>
    <p class="mb-3 text-xs text-muted-foreground">{{ t('batchExtendDesc') }}</p>
    <form class="flex flex-wrap items-end gap-3" @submit.prevent="extendSubscriptions">
      <div class="flex flex-col gap-1.5">
        <Label class="text-xs">{{ t('selectPlan') }}</Label>
        <select v-model="extendForm.plan_id" :class="selectClass" required>
          <option value="">{{ t('extendAllPlans') }}</option>
          <option v-for="plan in subscriptionPlans" :key="plan.id" :value="plan.id">{{ plan.name }}</option>
        </select>
      </div>
      <div class="flex flex-col gap-1.5">
        <Label class="text-xs">{{ t('extendDays') }}</Label>
        <Input v-model.number="extendForm.days" type="number" min="1" max="3650" required :placeholder="t('extendDaysPlaceholder')" class="w-32" />
      </div>
      <Button type="submit" :disabled="busy || !extendForm.days" size="sm">{{ t('confirmExtend') }}</Button>
    </form>
  </section>
  <section class="mt-4 overflow-hidden rounded-lg border border-border bg-card">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{{ t('userLabel') }}</TableHead>
          <TableHead>{{ t('planLabel') }}</TableHead>
          <TableHead>{{ t('accountStatus') }}</TableHead>
          <TableHead>{{ t('currentPeriod') }}</TableHead>
          <TableHead>{{ t('autoRenew') }}</TableHead>
          <TableHead>{{ t('createdAt') }}</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="sub in adminSubscriptions" :key="sub.id">
          <TableCell>
            <div class="font-medium">{{ sub.user_name }}</div>
            <div class="text-xs text-muted-foreground">{{ sub.email }}</div>
          </TableCell>
          <TableCell>{{ sub.plan_name }}</TableCell>
          <TableCell><Badge :variant="statusVariant(sub.status)">{{ statusBadge(sub.status) }}</Badge></TableCell>
          <TableCell class="text-xs">
            <span v-if="sub.current_period_start">{{ formatDate(sub.current_period_start) }} → {{ formatDate(sub.current_period_end) }}</span>
            <span v-else>—</span>
          </TableCell>
          <TableCell>
            <Badge :variant="sub.auto_renew ? 'secondary' : 'destructive'">{{ sub.auto_renew ? t('autoRenewOn') : t('autoRenewOff') }}</Badge>
          </TableCell>
          <TableCell class="text-xs">{{ formatDate(sub.created_at) }}</TableCell>
        </TableRow>
      </TableBody>
    </Table>
    <Empty v-if="!adminSubscriptions.length" :text="t('noSubscriptions')" />
  </section>
</template>
