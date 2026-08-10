<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NCard from '@/components/ui/NCard.vue'
import NButton from '@/components/ui/NButton.vue'
import NBadge from '@/components/ui/NBadge.vue'
import NModal from '@/components/ui/NModal.vue'
import NSelect from '@/components/ui/NSelect.vue'
import NIcon from '@/components/ui/NIcon.vue'
import MaterialViewer from '@/components/MaterialViewer.vue'
import { api } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user, isModerator, isHeadTeacher } from '@/lib/auth'
import {
  categoryMeta,
  STATUS_LABELS,
  STATUS_VARIANTS,
  ROLE_LABELS,
  formatDate,
  formatDateTime,
  formatBytes,
} from '@/lib/catalog'

const route = useRoute()
const router = useRouter()

const material = ref(null)
const error = ref('')
const qrOpen = ref(false)
const exportOpen = ref(false)

const permissions = ref([])
const users = ref([])
const grantUserId = ref('')

const id = computed(() => route.params.id)
const meta = computed(() => (material.value ? categoryMeta(material.value.category) : null))
const fileUrl = computed(() => `/api/materials/${id.value}/file`)
const qrUrl = computed(() => `/api/materials/${id.value}/qr`)

const canManage = computed(() => {
  if (!material.value || !user.value) return false
  if (user.value.role === 'head_teacher') return true
  if (user.value.role === 'librarian') return true
  return material.value.created_by === user.value.id
})

const userOptions = computed(() =>
  users.value
    .filter((u) => u.id !== user.value?.id && !permissions.value.some((p) => p.user_id === u.id))
    .map((u) => ({ value: String(u.id), label: `${u.full_name} (${ROLE_LABELS[u.role] || u.role})` })),
)

async function load() {
  try {
    material.value = await api.get(`/materials/${id.value}`)
    if (isHeadTeacher.value) {
      permissions.value = await api.get(`/materials/${id.value}/permissions`)
      users.value = await api.get('/users')
    }
  } catch (e) {
    error.value = e.message
  }
}

onMounted(load)

function download() {
  window.open(fileUrl.value + '?download=1', '_blank')
}

function exportCard(format) {
  exportOpen.value = false
  window.open(`/api/materials/${id.value}/export?format=${format}`, '_blank')
}

async function moderate(action) {
  const verb = action === 'approve' ? 'принят' : 'отклонён'
  try {
    material.value = await api.post(`/materials/${id.value}/${action}`)
    notify.success(`Материал ${verb}`)
  } catch (e) {
    notify.error('Не удалось изменить статус', { text: e.message })
  }
}

async function remove() {
  const ok = await notify.confirm({
    title: 'Удалить экспонат?',
    text: 'Материал и его файл будут удалены безвозвратно',
    confirmText: 'Удалить',
    danger: true,
  })
  if (!ok) return
  try {
    await api.del(`/materials/${id.value}`)
    notify.success('Экспонат удалён')
    router.push('/')
  } catch (e) {
    notify.error('Не удалось удалить', { text: e.message })
  }
}

async function grant() {
  if (!grantUserId.value) return
  try {
    await api.post(`/materials/${id.value}/permissions`, { user_id: Number(grantUserId.value) })
    grantUserId.value = ''
    permissions.value = await api.get(`/materials/${id.value}/permissions`)
    notify.success('Доступ выдан')
  } catch (e) {
    notify.error('Не удалось выдать доступ', { text: e.message })
  }
}

async function revoke(p) {
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
    notify.error('Не удалось отозвать доступ', { text: e.message })
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
    <div v-if="error" class="bg-surface border border-line py-16 px-6 text-center">
      <NIcon name="alert" :size="40" class="text-danger inline-block" />
      <p class="text-ink font-bold mt-3">{{ error }}</p>
      <NButton variant="secondary" class="mt-4" icon="chevronLeft" @click="router.back()">Назад</NButton>
    </div>

    <template v-else-if="material">
      <NPageHeader :title="material.title" :subtitle="`Экспонат №${material.id} · ${meta.label}`">
        <template #actions>
          <NButton v-if="material.file_name" variant="secondary" icon="download" @click="download">
            Скачать
          </NButton>
          <NButton variant="secondary" icon="qr" @click="qrOpen = true">QR-код</NButton>
          <NButton variant="secondary" icon="upload" @click="exportOpen = true">Экспорт</NButton>
          <NButton v-if="canManage" variant="secondary" icon="edit" @click="router.push(`/material/${material.id}/edit`)">
            Изменить
          </NButton>
          <NButton v-if="canManage" variant="danger" icon="trash" @click="remove">Удалить</NButton>
        </template>
      </NPageHeader>

      <div
        v-if="isModerator && material.status === 'pending'"
        class="mb-6 bg-surface border border-line border-l-[3px] border-l-warning px-4 py-3 flex flex-col sm:flex-row sm:items-center gap-3"
      >
        <div class="flex-1">
          <div class="font-bold text-ink text-sm">Материал ожидает модерации</div>
          <div class="text-xs text-muted">Добавлен: {{ material.author_name || 'неизвестно' }}</div>
        </div>
        <div class="flex gap-2">
          <NButton size="sm" icon="check" @click="moderate('approve')">Принять</NButton>
          <NButton size="sm" variant="danger" icon="close" @click="moderate('reject')">Отклонить</NButton>
        </div>
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">
        <div class="xl:col-span-2">
          <MaterialViewer :material="material" @download="download" />
        </div>

        <div class="space-y-6">
          <NCard title="Подробная информация" accent>
            <dl class="divide-y divide-line text-sm">
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Категория</dt>
                <dd class="text-ink font-semibold flex items-center gap-1.5">
                  <NIcon :name="meta.icon" :size="15" class="text-accent" />{{ meta.label }}
                </dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Статус</dt>
                <dd>
                  <NBadge :variant="STATUS_VARIANTS[material.status]" dot>{{ STATUS_LABELS[material.status] }}</NBadge>
                </dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Состояние</dt>
                <dd class="text-ink text-right">{{ material.condition || '—' }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Физический экземпляр</dt>
                <dd class="text-ink text-right">{{ material.location || '—' }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Время происхождения</dt>
                <dd class="text-ink text-right">{{ formatDate(material.origin_date) }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Добавлен в базу</dt>
                <dd class="text-ink text-right">{{ formatDateTime(material.created_at) }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Добавил</dt>
                <dd class="text-ink text-right">{{ material.author_name || '—' }}</dd>
              </div>
              <div v-if="material.file_name" class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted">Файл</dt>
                <dd class="text-ink text-right break-all">
                  {{ material.file_name }}
                  <span class="text-faint">({{ formatBytes(material.file_size) }})</span>
                </dd>
              </div>
            </dl>
          </NCard>

          <NCard v-if="material.description" title="Описание">
            <p class="text-sm text-muted whitespace-pre-line">{{ material.description }}</p>
          </NCard>

          <NCard v-if="isHeadTeacher" title="Права доступа" subtitle="Кому разрешён просмотр материала">
            <div class="space-y-3">
              <div v-if="!permissions.length" class="text-sm text-muted">
                Индивидуальные разрешения не выданы
              </div>
              <div
                v-for="p in permissions"
                :key="p.id"
                class="flex items-center justify-between gap-3 bg-surface-2 border border-line px-3 py-2"
              >
                <div class="min-w-0">
                  <div class="text-sm font-semibold text-ink truncate">{{ p.full_name }}</div>
                  <div class="text-xs text-muted">{{ ROLE_LABELS[p.role] || p.role }}</div>
                </div>
                <button
                  class="text-muted hover:text-danger transition-colors cursor-pointer"
                  title="Отозвать доступ"
                  @click="revoke(p)"
                >
                  <NIcon name="close" :size="16" />
                </button>
              </div>
              <div class="flex gap-2 items-end">
                <div class="flex-1">
                  <NSelect v-model="grantUserId" label="Выдать доступ" :options="userOptions" placeholder="Выберите пользователя…" />
                </div>
                <NButton icon="userPlus" :disabled="!grantUserId" @click="grant">Выдать</NButton>
              </div>
            </div>
          </NCard>
        </div>
      </div>

      <NModal v-model="qrOpen" title="QR-код экспоната" subtitle="Вы можете наклеить на физический экземпляр. Код ведёт на эту карточку">
        <div class="flex flex-col items-center gap-4">
          <img :src="qrUrl" alt="QR-код" class="w-64 h-64 bg-white p-2 border border-line" />
          <div class="text-xs text-muted">Экспонат №{{ material.id }} · {{ material.title }}</div>
        </div>
        <template #footer>
          <NButton variant="secondary" icon="download" @click="openQrPng">Открыть PNG</NButton>
          <NButton icon="printer" @click="printQr">Печать</NButton>
        </template>
      </NModal>

      <NModal v-model="exportOpen" title="Экспорт карточки" subtitle="Выберите формат файла">
        <div class="grid grid-cols-2 gap-3">
          <button
            class="border border-line hover:border-accent p-5 text-center transition-colors ng-tile-press cursor-pointer"
            @click="exportCard('json')"
          >
            <NIcon name="code" :size="26" class="text-accent inline-block" />
            <div class="font-bold text-ink mt-2">JSON</div>
          </button>
          <button
            class="border border-line hover:border-accent p-5 text-center transition-colors ng-tile-press cursor-pointer"
            @click="exportCard('xml')"
          >
            <NIcon name="fileText" :size="26" class="text-accent inline-block" />
            <div class="font-bold text-ink mt-2">XML</div>
          </button>
        </div>
      </NModal>
    </template>

    <div v-else class="text-center text-muted py-16">Загрузка…</div>
  </div>
</template>
