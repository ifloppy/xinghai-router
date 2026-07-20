<script setup lang="ts">
import { useConsoleStore } from '~/composables/useConsoleStore'
import Empty from '~/components/console/Empty.vue'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

const store = useConsoleStore()
const { t, users, groups, activityLogs, activityModels, activityFilters, activityTypeLabel, busy, can, formatDate, actionLabel, activityDetail, loadActivity, resetActivityFilters } = store

const selectClass = 'flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring'
const inputClass = 'flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring'
</script>

<template>
  <Card>
    <CardContent class="pt-6">
      <form class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4" @submit.prevent="loadActivity(true)">
        <div v-if="can('users.read')" class="flex flex-col gap-2">
          <Label>{{ t('userLabel') }}</Label>
          <select v-model="activityFilters.user_id" :class="selectClass">
            <option value="">{{ t('allUsers') }}</option>
            <option v-for="user in users" :key="user.id" :value="user.id">{{ user.name }} · {{ user.email }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <Label>{{ t('modelLabel') }}</Label>
          <select v-model="activityFilters.model" :class="selectClass">
            <option value="">{{ t('allModels') }}</option>
            <option v-for="model in activityModels" :key="model" :value="model">{{ model }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <Label>{{ t('groupLabel') }}</Label>
          <select v-model="activityFilters.group_id" :class="selectClass">
            <option value="">{{ t('allGroups') }}</option>
            <option v-for="group in groups" :key="group.id" :value="group.id">{{ group.name }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <Label>{{ t('typeLabel') }}</Label>
          <select v-model="activityFilters.type" :class="selectClass">
            <option value="">{{ t('allTypes') }}</option>
            <option value="request">{{ activityTypeLabel['request'] }}</option>
            <option value="login">{{ activityTypeLabel['login'] }}</option>
            <option value="register">{{ activityTypeLabel['register'] }}</option>
            <option value="logout">{{ activityTypeLabel['logout'] }}</option>
            <option value="topup">{{ activityTypeLabel['topup'] }}</option>
            <option value="operation">{{ activityTypeLabel['operation'] }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-2">
          <Label>{{ t('startTime') }}</Label>
          <input v-model="activityFilters.start" type="datetime-local" :class="inputClass">
        </div>
        <div class="flex flex-col gap-2">
          <Label>{{ t('endTime') }}</Label>
          <input v-model="activityFilters.end" type="datetime-local" :class="inputClass">
        </div>
        <div class="flex items-end gap-2">
          <Button type="submit" :disabled="busy">{{ t('filterLabel') }}</Button>
          <Button type="button" variant="outline" :disabled="busy" @click="resetActivityFilters">{{ t('resetFiltersLabel') }}</Button>
        </div>
      </form>
    </CardContent>
  </Card>

  <section class="mt-4 overflow-hidden rounded-lg border border-border bg-card">
    <div class="border-b border-border px-4 py-3">
      <h2 class="text-sm font-semibold">{{ t('usageLogs') }}</h2>
      <p class="text-xs text-muted-foreground">{{ t('usageLogsDesc') }}</p>
    </div>
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>{{ t('createdAt') }}</TableHead>
          <TableHead>{{ t('typeLabel') }}</TableHead>
          <TableHead>{{ t('userLabel') }}</TableHead>
          <TableHead>{{ t('modelLabel') }} / Action</TableHead>
          <TableHead>{{ t('groupLabel') }}</TableHead>
          <TableHead>Status / Duration</TableHead>
          <TableHead>Usage / Details</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="item in activityLogs" :key="`${item.type}-${item.id}`">
          <TableCell class="text-xs">{{ formatDate(item.created_at) }}</TableCell>
          <TableCell><Badge variant="outline">{{ activityTypeLabel[item.type] }}</Badge></TableCell>
          <TableCell>{{ item.user_name }}</TableCell>
          <TableCell>
            <code v-if="item.model" class="font-mono text-xs">{{ item.model }}</code>
            <span v-else class="text-sm">{{ actionLabel(item) }}</span>
          </TableCell>
          <TableCell>{{ item.group_name || '-' }}</TableCell>
          <TableCell>
            <Badge v-if="item.status_code != null" :variant="item.status_code < 400 ? 'secondary' : 'destructive'">{{ item.status_code }}</Badge>
            <span v-if="item.duration_ms != null" class="ml-1 text-xs text-muted-foreground">{{ item.duration_ms }} ms</span>
            <span v-if="item.status_code == null" class="text-sm text-green-600 dark:text-green-500">{{ t('success') }}</span>
          </TableCell>
          <TableCell><code class="font-mono text-xs">{{ activityDetail(item) }}</code></TableCell>
        </TableRow>
      </TableBody>
    </Table>
    <Empty v-if="!activityLogs.length" :text="t('noMatchingLogs')" />
  </section>
</template>
