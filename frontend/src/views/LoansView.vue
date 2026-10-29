<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { ArchiveIcon, CheckIcon, ClockIcon, PlusIcon } from '@lucide/vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import PageHeader from '@/components/PageHeader.vue'
import DatePicker from '@/components/DatePicker.vue'
import SearchSelect from '@/components/SearchSelect.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { formatDate, formatDateTime, isOverdue } from '@/lib/catalog'
import type { Loan, Material } from '@/lib/types'
import type { SelectOption } from '@/lib/ui'

type LoanTab = 'active' | 'archive'

const tab = ref<LoanTab>('active')
const loans = ref<Loan[]>([])
const loading = ref(true)

const addOpen = ref(false)
const books = ref<Material[]>([])
const newLoan = ref({ material_id: '', borrower_name: '', due_date: '', note: '' })
const saving = ref(false)

const overdueCount = computed(() => loans.value.filter(isOverdue).length)

function setTab(v: unknown) {
  if (v === 'active' || v === 'archive') tab.value = v
}

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
    if (tab.value === 'active') load()
    else tab.value = 'active'
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
    <PageHeader title="Прокат" description="Выдача книг из школьной библиотеки">
      <template #actions>
        <Button @click="openAdd">
          <PlusIcon />
          Выдать предмет
        </Button>
      </template>
    </PageHeader>

    <Tabs :model-value="tab" class="mb-4" @update:model-value="setTab">
      <TabsList>
        <TabsTrigger value="active">
          <ClockIcon />
          Активные
        </TabsTrigger>
        <TabsTrigger value="archive">
          <ArchiveIcon />
          Архив
        </TabsTrigger>
      </TabsList>
    </Tabs>

    <Card class="py-0">
      <Table>
        <TableHeader>
          <TableRow class="hover:bg-transparent">
            <TableHead class="pl-4">Предмет</TableHead>
            <TableHead>Кто взял</TableHead>
            <TableHead>Взято</TableHead>
            <TableHead>Срок возврата</TableHead>
            <TableHead v-if="tab === 'archive'" class="pr-4">Возвращено</TableHead>
            <TableHead v-else class="pr-4"><span class="sr-only">Действия</span></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="loading">
            <TableRow v-for="i in 4" :key="i">
              <TableCell class="pl-4"><Skeleton class="h-4 w-48" /></TableCell>
              <TableCell><Skeleton class="h-4 w-32" /></TableCell>
              <TableCell><Skeleton class="h-4 w-24" /></TableCell>
              <TableCell><Skeleton class="h-4 w-24" /></TableCell>
              <TableCell class="pr-4"><Skeleton class="ml-auto h-8 w-36" /></TableCell>
            </TableRow>
          </template>
          <TableRow
            v-for="row in loans"
            v-else
            :key="row.id"
            :class="isOverdue(row) ? 'bg-destructive/5 hover:bg-destructive/10' : ''"
          >
            <TableCell class="pl-4 whitespace-normal">
              <RouterLink :to="`/material/${row.material_id}`" class="font-medium hover:underline underline-offset-4">
                {{ row.material_title }}
              </RouterLink>
              <div v-if="row.note" class="mt-0.5 text-xs text-muted-foreground">{{ row.note }}</div>
            </TableCell>
            <TableCell>{{ row.borrower_name }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(row.taken_at) }}</TableCell>
            <TableCell>
              <div class="flex items-center gap-2">
                <span :class="isOverdue(row) ? 'font-medium text-destructive' : ''">{{ formatDate(row.due_date) }}</span>
                <Badge v-if="isOverdue(row)" variant="destructive">Просрочено</Badge>
              </div>
            </TableCell>
            <TableCell v-if="tab === 'archive'" class="pr-4 text-muted-foreground">
              {{ formatDateTime(row.returned_at) }}
            </TableCell>
            <TableCell v-else class="pr-4 text-right">
              <Button size="sm" variant="outline" @click="markReturned(row)">
                <CheckIcon />
                Предмет возвращён
              </Button>
            </TableCell>
          </TableRow>
          <TableRow v-if="!loading && !loans.length">
            <TableCell colspan="5" class="h-24 text-center text-muted-foreground">
              {{ tab === 'active' ? 'Сейчас ничего не выдано' : 'Архив пуст' }}
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>

    <Dialog v-model:open="addOpen">
      <DialogContent>
        <form class="grid gap-6" @submit.prevent="createLoan">
          <DialogHeader>
            <DialogTitle>Выдать предмет</DialogTitle>
            <DialogDescription>Запись о выдаче книги из библиотеки</DialogDescription>
          </DialogHeader>
          <FieldGroup>
            <Field>
              <FieldLabel for="loan-book">Предмет из библиотеки</FieldLabel>
              <SearchSelect
                id="loan-book"
                v-model="newLoan.material_id"
                :options="bookOptions"
                placeholder="Выберите книгу…"
                search-placeholder="Поиск по названию…"
                empty-text="Книги не найдены"
              />
            </Field>
            <Field>
              <FieldLabel for="loan-borrower">Кто взял</FieldLabel>
              <Input id="loan-borrower" v-model="newLoan.borrower_name" placeholder="ФИО и класс ученика" />
            </Field>
            <Field>
              <FieldLabel for="loan-due">Вернуть до</FieldLabel>
              <DatePicker id="loan-due" v-model="newLoan.due_date" />
            </Field>
            <Field>
              <FieldLabel for="loan-note">Заметка</FieldLabel>
              <Textarea id="loan-note" v-model="newLoan.note" placeholder="Дополнительная информация…" :rows="2" />
            </Field>
          </FieldGroup>
          <DialogFooter>
            <DialogClose as-child>
              <Button type="button" variant="outline">Отмена</Button>
            </DialogClose>
            <Button type="submit" :disabled="saving">
              <Spinner v-if="saving" />
              Записать
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </div>
</template>
