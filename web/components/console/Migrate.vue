<script setup lang="ts">
import { useConsoleStore } from '~/composables/useConsoleStore'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { TriangleAlert } from 'lucide-vue-next'

const store = useConsoleStore()
const { t, busy, migrateForm, migrateResult, migrateRunning, runMigration, error } = store
</script>

<template>
  <section class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-lg font-semibold">{{ t('migrate') }}</h2>
      <p class="text-sm text-muted-foreground">{{ t('migrateDesc') }}</p>
    </div>
  </section>

  <div class="mt-4 flex items-start gap-2 rounded-md border border-amber-500/40 bg-amber-500/5 px-4 py-3 text-sm text-amber-600 dark:text-amber-400">
    <TriangleAlert class="mt-0.5 h-4 w-4 shrink-0" />
    <span>{{ t('migrateWarning') }}</span>
  </div>

  <Card class="mt-4">
    <CardHeader>
      <CardTitle class="text-base">{{ t('migrate') }}</CardTitle>
      <CardDescription>{{ t('migrateDesc') }}</CardDescription>
    </CardHeader>
    <CardContent>
      <form class="flex flex-col gap-6" @submit.prevent="runMigration">
        <div class="flex flex-col gap-2">
          <Label>{{ t('migrateSourceDsn') }}</Label>
          <Input v-model="migrateForm.source_dsn" required :placeholder="t('migrateDsnPlaceholder')" class="max-w-xl font-mono text-sm" />
          <p class="text-xs text-muted-foreground">{{ t('migrateSourceDsn') }}: <code class="bg-muted px-1 rounded">mysql://user:pass@host:3306/db</code></p>
        </div>

        <div class="flex flex-col gap-2 max-w-xs">
          <Label>{{ t('migrateSourceDriver') }}</Label>
          <Select v-model="migrateForm.source_driver">
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="mysql">MySQL</SelectItem>
              <SelectItem value="postgres">PostgreSQL</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div v-if="error" class="text-sm text-red-500">{{ error }}</div>
        <div v-if="migrateResult" class="text-sm text-green-500">{{ t('migrateSuccess') }}: {{ migrateResult }}</div>

        <Button type="submit" :disabled="migrateRunning || busy" class="w-fit">
          {{ migrateRunning ? t('migrateRunning') : t('migrateStart') }}
        </Button>
      </form>
    </CardContent>
  </Card>
</template>
