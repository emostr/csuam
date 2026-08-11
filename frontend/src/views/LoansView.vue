<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NButton from '@/components/ui/NButton.vue'
import NTable from '@/components/ui/NTable.vue'
import NTabs from '@/components/ui/NTabs.vue'
import NBadge from '@/components/ui/NBadge.vue'
import NModal from '@/components/ui/NModal.vue'
import NInput from '@/components/ui/NInput.vue'
import NTextarea from '@/components/ui/NTextarea.vue'
import NSelect from '@/components/ui/NSelect.vue'
import NDateInput from '@/components/ui/NDateInput.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { formatDate, formatDateTime, isOverdue } from '@/lib/catalog'
import type { Loan, Material } from '@/lib/types'
import type { SelectOption, TabItem, TableColumn } from '@/lib/ui'

const router = useRouter()

const tab = ref<string | number>('active')
const loans = ref<Loan[]>([])
const loading = ref(true)

const addOpen = ref(false)
const books = ref<Material[]>([])
const newLoan = ref({ material_id: '', borrower_name: '', due_date: '', note: '' })
const saving = ref(false)

const tabs = computed<TabItem[]>(() => [
  { value: 'active', label: 'Активные', icon: 'clock' },
  { value: 'archive', label: 'Архив', icon: 'archive' },
])

const columns = computed<TableColumn[]>(() => {
  const cols: TableColumn[] = [
    { key: 'material_title', label: 'Предмет' },
    { key: 'borrower_name', label: 'Кто взял' },
    { key: 'taken_at', label: 'Взято' },
    { key: 'due_date', label: 'Срок возврата' },
  ]
  if (tab.value === 'active') {
    cols.push({ key: 'actions', label: '', align: 'right' })
  } else {
    cols.push({ key: 'returned_at', label: 'Возвращено' })
  }
  return cols
})

const overdueCount = computed(() => loans.value.filter(isOverdue).length)

async function load() {
  loading.value = true
  try {
    loans.value = await api.get<Loan[]>('/loans' + (tab.value === 'archive' ? '?archive=1' : ''))
  } catch (e) {
    notify.error('Не удалось загрузить прокат', { text: errorMessage(e) })
  } finally {
    loading.value = false
  }
}

watch(tab, load)

onMounted(async () => {
  await load()
  if (tab.value === 'active' && overdueCount.value > 0 && !sessionStorage.getItem('csuam-loans-warned')) {
    sessionStorage.setItem('csuam-loans-warned', '1')
    notify.warning('Есть просроченные возвраты', {
      text: `Книг не возвращено вовремя: ${overdueCount.value}. Они выделены в списке.`,
      timer: 5000,
    })
  }
})

async function openAdd() {
  newLoan.value = { material_id: '', borrower_name: '', due_date: '', note: '' }
  addOpen.value = true
  if (!books.value.length) {
    try {
      books.value = await api.get<Material[]>('/materials?category=library')
    } catch {}
  }
}

const bookOptions = computed<SelectOption[]>(() =>
  books.value.map((b) => ({ value: String(b.id), label: b.title })),
)

async function createLoan() {
  if (!newLoan.value.material_id || !newLoan.value.borrower_name || !newLoan.value.due_date) {
    notify.warning('Заполните предмет, имя и срок возврата')
    return
  }
  saving.value = true
  try {
    await api.post('/loans', {
      material_id: Number(newLoan.value.material_id),
      borrower_name: newLoan.value.borrower_name,
      note: newLoan.value.note,
      due_date: newLoan.value.due_date,
    })
    addOpen.value = false
    notify.success('Запись добавлена')
    tab.value = 'active'
    load()
  } catch (e) {
    notify.error('Не удалось добавить запись', { text: errorMessage(e) })
  } finally {
    saving.value = false
  }
}

async function markReturned(loan: Loan) {
  const ok = await notify.confirm({
    title: 'Предмет возвращён?',
    text: `«${loan.material_title}» — ${loan.borrower_name}. Запись будет перенесена в архив.`,
    confirmText: 'Возвращён',
    icon: 'question',
  })
  if (!ok) return
  try {
    await api.post(`/loans/${loan.id}/return`)
    notify.success('Возврат отмечен')
    load()
  } catch (e) {
    notify.error('Не удалось отметить возврат', { text: errorMessage(e) })
  }
}
</script>

<template>
  <div>
    <NPageHeader title="Прокат">
      <template #actions>
        <NButton icon="plus" @click="openAdd">Выдать предмет</NButton>
      </template>
    </NPageHeader>

    <NTabs v-model="tab" :tabs="tabs" class="mb-6" />

    <div v-if="loading" class="text-center text-muted py-16">Загрузка…</div>
    <NTable v-else :columns="columns" :rows="loans">
      <template #material_title="{ row }">
        <button
          class="font-semibold text-ink hover:text-accent transition-colors cursor-pointer text-left"
          @click="router.push(`/material/${row.material_id}`)"
        >
          {{ row.material_title }}
        </button>
        <div v-if="row.note" class="text-xs text-muted mt-0.5">{{ row.note }}</div>
      </template>
      <template #taken_at="{ row }">
        <span class="text-muted">{{ formatDate(row.taken_at) }}</span>
      </template>
      <template #due_date="{ row }">
        <div class="flex items-center gap-2">
          <span :class="isOverdue(row) ? 'text-danger font-bold' : ''">{{ formatDate(row.due_date) }}</span>
          <NBadge v-if="isOverdue(row)" variant="danger" dot>Просрочено</NBadge>
        </div>
      </template>
      <template #returned_at="{ row }">
        <span class="text-muted">{{ formatDateTime(row.returned_at) }}</span>
      </template>
      <template #actions="{ row }">
        <NButton size="sm" variant="secondary" icon="check" @click="markReturned(row)">
          Предмет возвращён
        </NButton>
      </template>
    </NTable>

    <NModal v-model="addOpen" title="Выдать предмет" subtitle="Запись о выдаче книги из библиотеки">
      <div class="space-y-4">
        <NSelect
          v-model="newLoan.material_id"
          label="Предмет из библиотеки"
          :options="bookOptions"
          placeholder="Выберите книгу…"
        />
        <NInput v-model="newLoan.borrower_name" label="Кто взял" placeholder="ФИО и класс ученика" icon="user" />
        <NDateInput v-model="newLoan.due_date" label="Вернуть до" />
        <NTextarea v-model="newLoan.note" label="Заметка" placeholder="Дополнительная информация…" :rows="2" />
      </div>
      <template #footer="{ close }">
        <NButton variant="ghost" @click="close">Отмена</NButton>
        <NButton :loading="saving" icon="check" @click="createLoan">Записать</NButton>
      </template>
    </NModal>
  </div>
</template>
