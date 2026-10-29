<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  BracesIcon,
  CheckIcon,
  ChevronDownIcon,
  ChevronLeftIcon,
  ClockIcon,
  CodeIcon,
  DownloadIcon,
  ExternalLinkIcon,
  PencilIcon,
  PrinterIcon,
  QrCodeIcon,
  ShareIcon,
  Trash2Icon,
  TriangleAlertIcon,
  UserPlusIcon,
  XIcon,
} from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty'
import { Field, FieldLabel } from '@/components/ui/field'
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemMedia,
  ItemTitle,
} from '@/components/ui/item'
import { Skeleton } from '@/components/ui/skeleton'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import SearchSelect from '@/components/SearchSelect.vue'
import MaterialViewer from '@/components/MaterialViewer.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user, isModerator, isHeadTeacher } from '@/lib/auth'
import {
  categoryMeta,
  STATUS_LABELS,
  STATUS_TONES,
  ROLE_LABELS,
  formatDate,
  formatDateTime,
  formatBytes,
} from '@/lib/catalog'
import { fileIcon, fileKind, fileUrl } from '@/lib/files'
import type { Material, MaterialFile, Permission, User } from '@/lib/types'
import type { SelectOption } from '@/lib/ui'

const route = useRoute()
const router = useRouter()

const material = ref<Material | null>(null)
const error = ref('')
const qrOpen = ref(false)

const permissions = ref<Permission[]>([])
const users = ref<User[]>([])
const grantUserId = ref('')

const id = computed(() => String(route.params.id))
const meta = computed(() => categoryMeta(material.value?.category ?? 'photos'))
const qrUrl = computed(() => `/api/materials/${id.value}/qr`)
const files = computed<MaterialFile[]>(() => material.value?.files ?? [])

const canManage = computed(() => {
  if (!material.value || !user.value) return false
  if (user.value.role === 'head_teacher') return true
  if (user.value.role === 'librarian') return true
  return material.value.created_by === user.value.id
})

const userOptions = computed<SelectOption[]>(() =>
  users.value
    .filter((u) => u.id !== user.value?.id && !permissions.value.some((p) => p.user_id === u.id))
    .map((u) => ({ value: String(u.id), label: `${u.full_name} (${ROLE_LABELS[u.role]})` })),
)

const details = computed(() => {
  const m = material.value
  if (!m) return []
  const rows = [
    { label: 'Состояние', value: m.condition || '—' },
    { label: 'Физический экземпляр', value: m.location || '—' },
    { label: 'Время происхождения', value: formatDate(m.origin_date) },
    { label: 'Добавлен в базу', value: formatDateTime(m.created_at) },
    { label: 'Добавил', value: m.author_name || '—' },
  ]
  if (files.value.length) rows.push({ label: 'Файлов', value: String(files.value.length) })
  return rows
})

async function load() {
  try {
    material.value = await api.get<Material>(`/materials/${id.value}`)
    if (isHeadTeacher.value) {
      permissions.value = await api.get<Permission[]>(`/materials/${id.value}/permissions`)
      users.value = await api.get<User[]>('/users')
    }
  } catch (e) {
    error.value = errorMessage(e)
  }
}

onMounted(load)

function download(file: MaterialFile) {
  if (!material.value) return
  window.open(fileUrl(material.value.id, file) + '?download=1', '_blank')
}

function exportCard(format: 'json' | 'xml') {
  window.open(`/api/materials/${id.value}/export?format=${format}`, '_blank')
}

async function moderate(action: 'approve' | 'reject') {
  const verb = action === 'approve' ? 'принят' : 'отклонён'
  try {
    material.value = await api.post<Material>(`/materials/${id.value}/${action}`)
    notify.success(`Материал ${verb}`)
  } catch (e) {
    notify.error('Не удалось изменить статус', { text: errorMessage(e) })
  }
}

async function remove() {
  const ok = await notify.confirm({
    title: 'Удалить экспонат?',
    text: 'Материал и его файлы будут удалены безвозвратно',
    confirmText: 'Удалить',
    danger: true,
  })
  if (!ok) return
  try {
    await api.del(`/materials/${id.value}`)
    notify.success('Экспонат удалён')
    router.push('/')
  } catch (e) {
    notify.error('Не удалось удалить', { text: errorMessage(e) })
  }
}

async function grant() {
  if (!grantUserId.value) return
  try {
    await api.post(`/materials/${id.value}/permissions`, { user_id: Number(grantUserId.value) })
    grantUserId.value = ''
    permissions.value = await api.get<Permission[]>(`/materials/${id.value}/permissions`)
    notify.success('Доступ выдан')
  } catch (e) {
    notify.error('Не удалось выдать доступ', { text: errorMessage(e) })
  }
}

async function revoke(p: Permission) {
  const ok = await notify.confirm({
    title: 'Отозвать доступ?',
    text: `${p.full_name} больше не сможет просматривать этот материал`,
    confirmText: 'Отозвать',
    danger: true,
  })
  if (!ok) return
  try {
    await api.del(`/materials/${id.value}/permissions/${p.user_id}`)
    permissions.value = permissions.value.filter((x) => x.user_id !== p.user_id)
    notify.success('Доступ отозван')
  } catch (e) {
    notify.error('Не удалось отозвать доступ', { text: errorMessage(e) })
  }
}

function openQrPng() {
  window.open(qrUrl.value, '_blank')
}

function printQr() {
  const w = window.open('', '_blank')
  if (!w) return
  w.document.write(
    `<html><head><title>Экспонат №${id.value}</title></head><body style="display:flex;flex-direction:column;align-items:center;font-family:sans-serif"><img src="${qrUrl.value}" style="width:280px;height:280px" onload="window.print()"><p>Экспонат №${id.value} — ${material.value?.title || ''}</p></body></html>`,
  )
  w.document.close()
}
</script>

<template>
  <div>
    <Empty v-if="error" class="border border-dashed">
      <EmptyHeader>
        <EmptyMedia variant="icon" class="bg-destructive/10 text-destructive">
          <TriangleAlertIcon />
        </EmptyMedia>
        <EmptyTitle>Не удалось открыть экспонат</EmptyTitle>
        <EmptyDescription>{{ error }}</EmptyDescription>
      </EmptyHeader>
      <EmptyContent>
        <Button variant="outline" @click="router.back()">
          <ChevronLeftIcon />
          Назад
        </Button>
      </EmptyContent>
    </Empty>

    <template v-else-if="material">
      <PageHeader :title="material.title" :description="`Экспонат №${material.id} · ${meta.label}`">
        <template #actions>
          <Button v-if="files.length === 1" variant="outline" @click="download(files[0])">
            <DownloadIcon />
            Скачать
          </Button>
          <Button variant="outline" @click="qrOpen = true">
            <QrCodeIcon />
            QR-код
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <Button variant="outline">
                <ShareIcon />
                Экспорт
                <ChevronDownIcon class="text-muted-foreground" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem @select="exportCard('json')">
                <BracesIcon />
                JSON
              </DropdownMenuItem>
              <DropdownMenuItem @select="exportCard('xml')">
                <CodeIcon />
                XML
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <Button v-if="canManage" variant="outline" @click="router.push(`/material/${material.id}/edit`)">
            <PencilIcon />
            Изменить
          </Button>
          <Button v-if="canManage" variant="destructive" @click="remove">
            <Trash2Icon />
            Удалить
          </Button>
        </template>
      </PageHeader>

      <Alert v-if="isModerator && material.status === 'pending'" class="mb-6">
        <ClockIcon />
        <AlertTitle>Материал ожидает модерации</AlertTitle>
        <AlertDescription>
          <p>Добавлен: {{ material.author_name || 'неизвестно' }}</p>
          <div class="mt-2 flex gap-2">
            <Button size="sm" @click="moderate('approve')">
              <CheckIcon />
              Принять
            </Button>
            <Button size="sm" variant="destructive" @click="moderate('reject')">
              <XIcon />
              Отклонить
            </Button>
          </div>
        </AlertDescription>
      </Alert>

      <div class="grid grid-cols-1 gap-6 xl:grid-cols-3">
        <div class="space-y-6 xl:col-span-2">
          <MaterialViewer :material="material" @download="download" />

          <Card v-if="files.length">
            <CardHeader>
              <CardTitle>Файлы</CardTitle>
              <CardDescription>Прикреплено: {{ files.length }}</CardDescription>
            </CardHeader>
            <CardContent>
              <ItemGroup class="gap-2">
                <Item v-for="f in files" :key="f.id" variant="outline" size="sm">
                  <ItemMedia variant="icon">
                    <component :is="fileIcon(fileKind(f.name, f.mime))" class="text-primary" />
                  </ItemMedia>
                  <ItemContent class="min-w-0">
                    <ItemTitle class="break-all">{{ f.name }}</ItemTitle>
                    <ItemDescription>{{ formatBytes(f.size) }}</ItemDescription>
                  </ItemContent>
                  <ItemActions>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button variant="ghost" size="icon-sm" @click="download(f)">
                          <DownloadIcon />
                          <span class="sr-only">Скачать</span>
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Скачать</TooltipContent>
                    </Tooltip>
                  </ItemActions>
                </Item>
              </ItemGroup>
            </CardContent>
          </Card>
        </div>

        <div class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Подробная информация</CardTitle>
            </CardHeader>
            <CardContent>
              <dl class="divide-y text-sm">
                <div class="flex justify-between gap-4 py-2.5 first:pt-0">
                  <dt class="text-muted-foreground">Категория</dt>
                  <dd class="flex items-center gap-1.5 font-medium">
                    <component :is="meta.icon" class="size-4 text-muted-foreground" />{{ meta.label }}
                  </dd>
                </div>
                <div class="flex items-center justify-between gap-4 py-2.5">
                  <dt class="text-muted-foreground">Статус</dt>
                  <dd>
                    <StatusBadge :tone="STATUS_TONES[material.status]">{{ STATUS_LABELS[material.status] }}</StatusBadge>
                  </dd>
                </div>
                <div v-for="row in details" :key="row.label" class="flex justify-between gap-4 py-2.5 last:pb-0">
                  <dt class="text-muted-foreground">{{ row.label }}</dt>
                  <dd class="text-right">{{ row.value }}</dd>
                </div>
              </dl>
            </CardContent>
          </Card>

          <Card v-if="material.description">
            <CardHeader>
              <CardTitle>Описание</CardTitle>
            </CardHeader>
            <CardContent>
              <p class="whitespace-pre-line text-muted-foreground">{{ material.description }}</p>
            </CardContent>
          </Card>

          <Card v-if="isHeadTeacher">
            <CardHeader>
              <CardTitle>Права доступа</CardTitle>
              <CardDescription>Кому разрешён просмотр материала</CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <p v-if="!permissions.length" class="text-muted-foreground">Индивидуальные разрешения не выданы</p>
              <ItemGroup v-else class="gap-2">
                <Item v-for="p in permissions" :key="p.id" variant="outline" size="sm">
                  <ItemContent class="min-w-0">
                    <ItemTitle class="truncate">{{ p.full_name }}</ItemTitle>
                    <ItemDescription>{{ ROLE_LABELS[p.role] || p.role }}</ItemDescription>
                  </ItemContent>
                  <ItemActions>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button variant="ghost" size="icon-sm" @click="revoke(p)">
                          <XIcon />
                          <span class="sr-only">Отозвать доступ</span>
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Отозвать доступ</TooltipContent>
                    </Tooltip>
                  </ItemActions>
                </Item>
              </ItemGroup>
              <Field>
                <FieldLabel for="grant-user">Выдать доступ</FieldLabel>
                <div class="flex gap-2">
                  <SearchSelect
                    id="grant-user"
                    v-model="grantUserId"
                    :options="userOptions"
                    placeholder="Выберите пользователя…"
                    search-placeholder="Поиск по имени…"
                    empty-text="Пользователи не найдены"
                  />
                  <Button :disabled="!grantUserId" @click="grant">
                    <UserPlusIcon />
                    Выдать
                  </Button>
                </div>
              </Field>
            </CardContent>
          </Card>
        </div>
      </div>

      <Dialog v-model:open="qrOpen">
        <DialogContent class="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>QR-код экспоната</DialogTitle>
            <DialogDescription>Наклейте его на физический экземпляр — код ведёт на эту карточку.</DialogDescription>
          </DialogHeader>
          <div class="flex flex-col items-center gap-3">
            <img :src="qrUrl" alt="QR-код" class="size-56 rounded-lg bg-white p-2 ring-1 ring-foreground/10" />
            <p class="text-center text-xs text-muted-foreground">Экспонат №{{ material.id }} · {{ material.title }}</p>
          </div>
          <DialogFooter>
            <Button variant="outline" @click="openQrPng">
              <ExternalLinkIcon />
              Открыть PNG
            </Button>
            <Button @click="printQr">
              <PrinterIcon />
              Печать
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </template>

    <div v-else class="space-y-6">
      <div class="space-y-2">
        <Skeleton class="h-8 w-80 max-w-full" />
        <Skeleton class="h-4 w-48" />
      </div>
      <div class="grid grid-cols-1 gap-6 xl:grid-cols-3">
        <Skeleton class="h-96 rounded-xl xl:col-span-2" />
        <Skeleton class="h-72 rounded-xl" />
      </div>
    </div>
  </div>
</template>
