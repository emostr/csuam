<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import {
  AwardIcon,
  BookOpenIcon,
  ClockIcon,
  ImageIcon,
  PackageIcon,
  PlusIcon,
  SearchIcon,
  TriangleAlertIcon,
  UploadIcon,
} from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import PageHeader from '@/components/PageHeader.vue'
import StatCard from '@/components/StatCard.vue'
import DatePicker from '@/components/DatePicker.vue'
import MaterialGrid from '@/components/MaterialGrid.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user, isModerator } from '@/lib/auth'
import { CATEGORIES } from '@/lib/catalog'
import type { AnalyticsSummary, Material, Notifications } from '@/lib/types'

const route = useRoute()
const router = useRouter()

const queryString = (v: unknown): string => (typeof v === 'string' ? v : '')
const fromSelect = (v: unknown): string => (typeof v === 'string' && v !== 'all' ? v : '')

const materials = ref<Material[]>([])
const loading = ref(true)
const summary = ref<AnalyticsSummary | null>(null)

const q = ref(queryString(route.query.q))
const category = ref(queryString(route.query.category))
const status = ref(queryString(route.query.status))
const from = ref('')
const to = ref('')

const statusOptions = [
  { value: 'approved', label: 'Опубликованные' },
  { value: 'pending', label: 'На модерации' },
  { value: 'rejected', label: 'Отклонённые' },
]

const importInput = ref<HTMLInputElement | null>(null)

async function load() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (q.value) params.set('q', q.value)
    if (category.value) params.set('category', category.value)
    if (status.value) params.set('status', status.value)
    if (from.value) params.set('from', from.value)
    if (to.value) params.set('to', to.value)
    materials.value = await api.get<Material[]>('/materials?' + params.toString())
  } catch (e) {
    notify.error('Не удалось загрузить материалы', { text: errorMessage(e) })
  } finally {
    loading.value = false
  }
}

let syncing = false

watch([category, status, from, to], () => {
  if (!syncing) load()
})
watch(
  () => route.query,
  async (nq) => {
    syncing = true
    q.value = queryString(nq.q)
    status.value = queryString(nq.status)
    await nextTick()
    syncing = false
    load()
  },
)

const counts = computed<Partial<Record<string, number>>>(() => {
  const map: Partial<Record<string, number>> = {}
  for (const c of summary.value?.by_category || []) map[c.category] = c.count
  return map
})

onMounted(async () => {
  load()
  try {
    summary.value = await api.get<AnalyticsSummary>('/analytics/summary')
  } catch {}
  if (isModerator.value && !sessionStorage.getItem('csuam-overdue-warned')) {
    try {
      const n = await api.get<Notifications>('/notifications')
      if (n.overdue_loans > 0) {
        sessionStorage.setItem('csuam-overdue-warned', '1')
        notify.warning('Просрочен возврат книг', {
          text: `Читателей с просроченным сроком возврата: ${n.overdue_loans}. Откройте раздел «Прокат».`,
          timer: 6000,
        })
      }
    } catch {}
  }
})

function startImport() {
  importInput.value?.click()
}

async function onImportFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  const fd = new FormData()
  fd.append('file', file)
  try {
    const m = await api.postForm<Material>('/materials/import', fd)
    notify.success('Карточка импортирована')
    router.push(`/material/${m.id}`)
  } catch (err) {
    notify.error('Импорт не удался', { text: errorMessage(err) })
  }
}
</script>

<template>
  <div>
    <PageHeader title="Главная" :description="`Добро пожаловать, ${user?.full_name || ''}!`">
      <template #actions>
        <Button variant="outline" @click="startImport">
          <UploadIcon />
          Импорт карточки
        </Button>
        <Button @click="router.push('/add')">
          <PlusIcon />
          Добавить экспонат
        </Button>
        <input
          ref="importInput"
          type="file"
          accept=".json,.xml,application/json,text/xml,application/xml"
          class="hidden"
          @change="onImportFile"
        />
      </template>
    </PageHeader>

    <div v-if="summary" class="mb-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
      <StatCard label="Всего единиц" :value="summary.total" :icon="PackageIcon" />
      <StatCard label="Фотографии" :value="counts.photos || 0" :icon="ImageIcon" />
      <StatCard label="Библиотека" :value="counts.library || 0" :icon="BookOpenIcon" />
      <StatCard v-if="isModerator" label="На модерации" :value="summary.pending" :icon="ClockIcon" />
      <StatCard v-else label="Награды" :value="counts.awards || 0" :icon="AwardIcon" />
    </div>

    <Alert v-if="isModerator && summary && summary.overdue_loans > 0" variant="destructive" class="mb-6">
      <TriangleAlertIcon />
      <AlertTitle>Просроченный прокат</AlertTitle>
      <AlertDescription>
        <p>
          Не возвращено вовремя книг: {{ summary.overdue_loans }}.
          <RouterLink to="/loans" class="font-medium underline underline-offset-4">Перейти в «Прокат»</RouterLink>
        </p>
      </AlertDescription>
    </Alert>

    <div
      class="mb-6 grid grid-cols-1 gap-3 sm:grid-cols-2"
      :class="isModerator ? 'lg:grid-cols-[minmax(0,2fr)_repeat(4,minmax(0,1fr))]' : 'lg:grid-cols-[minmax(0,2fr)_repeat(3,minmax(0,1fr))]'"
    >
      <form class="sm:col-span-2 lg:col-span-1" role="search" @submit.prevent="load">
        <InputGroup>
          <InputGroupAddon>
            <SearchIcon />
          </InputGroupAddon>
          <InputGroupInput v-model="q" placeholder="Поиск по названию и описанию…" aria-label="Поиск" />
        </InputGroup>
      </form>
      <Select :model-value="category || 'all'" @update:model-value="category = fromSelect($event)">
        <SelectTrigger class="w-full" aria-label="Категория">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Все категории</SelectItem>
          <SelectItem v-for="c in CATEGORIES" :key="c.id" :value="c.id">{{ c.label }}</SelectItem>
        </SelectContent>
      </Select>
      <DatePicker v-model="from" placeholder="Происхождение с" />
      <DatePicker v-model="to" placeholder="Происхождение по" />
      <Select v-if="isModerator" :model-value="status || 'all'" @update:model-value="status = fromSelect($event)">
        <SelectTrigger class="w-full" aria-label="Статус">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Все статусы</SelectItem>
          <SelectItem v-for="s in statusOptions" :key="s.value" :value="s.value">{{ s.label }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <MaterialGrid
      :materials="materials"
      :loading="loading"
      empty-text="Попробуйте изменить условия поиска или сбросить фильтры."
    />
  </div>
</template>
