<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NButton from '@/components/ui/NButton.vue'
import NInput from '@/components/ui/NInput.vue'
import NSelect from '@/components/ui/NSelect.vue'
import NDateInput from '@/components/ui/NDateInput.vue'
import NTile from '@/components/ui/NTile.vue'
import NAlert from '@/components/ui/NAlert.vue'
import MaterialCard from '@/components/MaterialCard.vue'
import { api } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user, isModerator } from '@/lib/auth'
import { CATEGORIES } from '@/lib/catalog'

const route = useRoute()
const router = useRouter()

const materials = ref([])
const loading = ref(true)
const summary = ref(null)

const q = ref(route.query.q || '')
const category = ref(route.query.category || '')
const status = ref(route.query.status || '')
const from = ref('')
const to = ref('')

const categoryOptions = [
  { value: '', label: 'Все категории' },
  ...CATEGORIES.map((c) => ({ value: c.id, label: c.label })),
]
const statusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'approved', label: 'Опубликованные' },
  { value: 'pending', label: 'На модерации' },
  { value: 'rejected', label: 'Отклонённые' },
]

const importInput = ref(null)

async function load() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (q.value) params.set('q', q.value)
    if (category.value) params.set('category', category.value)
    if (status.value) params.set('status', status.value)
    if (from.value) params.set('from', from.value)
    if (to.value) params.set('to', to.value)
    materials.value = await api.get('/materials?' + params.toString())
  } catch (e) {
    notify.error('Не удалось загрузить материалы', { text: e.message })
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
    q.value = nq.q || ''
    status.value = nq.status || ''
    await nextTick()
    syncing = false
    load()
  },
)

const counts = computed(() => {
  const map = {}
  for (const c of summary.value?.by_category || []) map[c.category] = c.count
  return map
})

onMounted(async () => {
  load()
  try {
    summary.value = await api.get('/analytics/summary')
  } catch {}
  if (isModerator.value && !sessionStorage.getItem('csuam-overdue-warned')) {
    try {
      const n = await api.get('/notifications')
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

async function onImportFile(e) {
  const file = e.target.files?.[0]
  e.target.value = ''
  if (!file) return
  const fd = new FormData()
  fd.append('file', file)
  try {
    const m = await api.postForm('/materials/import', fd)
    notify.success('Карточка импортирована')
    router.push(`/material/${m.id}`)
  } catch (err) {
    notify.error('Импорт не удался', { text: err.message })
  }
}
</script>

<template>
  <div>
    <NPageHeader
      title="Главная"
      :subtitle="`Добро пожаловать, ${user?.full_name || ''}!`"
    >
      <template #actions>
        <NButton variant="secondary" icon="upload" @click="startImport">Импорт карточки</NButton>
        <NButton icon="plus" @click="router.push('/add')">Добавить экспонат</NButton>
        <input
          ref="importInput"
          type="file"
          accept=".json,.xml,application/json,text/xml,application/xml"
          class="hidden"
          @change="onImportFile"
        />
      </template>
    </NPageHeader>

    <div v-if="summary" class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <NTile label="Всего единиц" :value="summary.total" icon="box" solid />
      <NTile label="Фотографии" :value="counts.photos || 0" icon="image" />
      <NTile label="Библиотека" :value="counts.library || 0" icon="book" />
      <NTile
        v-if="isModerator"
        label="На модерации"
        :value="summary.pending"
        icon="clock"
      />
      <NTile v-else label="Награды" :value="counts.awards || 0" icon="award" />
    </div>

    <NAlert
      v-if="isModerator && summary && summary.overdue_loans > 0"
      variant="danger"
      title="Просроченный прокат"
      class="mb-6"
    >
      Не возвращено вовремя книг: {{ summary.overdue_loans }}.
      <RouterLink to="/loans" class="text-accent font-semibold hover:underline">Перейти в «Прокат»</RouterLink>
    </NAlert>

    <div class="bg-surface border border-line p-4 mb-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3 items-end">
      <form class="lg:col-span-2" @submit.prevent="load">
        <NInput v-model="q" label="Поиск" placeholder="Ключевые слова…" icon="search" @keyup.enter="load" />
      </form>
      <NSelect v-model="category" label="Категория" :options="categoryOptions" />
      <NDateInput v-model="from" label="Происхождение с" />
      <NDateInput v-model="to" label="Происхождение по" />
      <NSelect
        v-if="isModerator"
        v-model="status"
        label="Статус"
        :options="statusOptions"
        class="lg:col-start-5"
      />
    </div>

    <div v-if="loading" class="text-center text-muted py-16">Загрузка…</div>
    <div v-else-if="!materials.length" class="bg-surface border border-line py-16 text-center text-muted">
      Ничего не найдено. Попробуйте изменить условия поиска.
    </div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <MaterialCard v-for="m in materials" :key="m.id" :material="m" />
    </div>
  </div>
</template>
