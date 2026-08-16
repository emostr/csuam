<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NCard from '@/components/ui/NCard.vue'
import NButton from '@/components/ui/NButton.vue'
import NInput from '@/components/ui/NInput.vue'
import NTextarea from '@/components/ui/NTextarea.vue'
import NSelect from '@/components/ui/NSelect.vue'
import NDateInput from '@/components/ui/NDateInput.vue'
import NAlert from '@/components/ui/NAlert.vue'
import NIcon from '@/components/ui/NIcon.vue'
import TextEditor from '@/components/TextEditor.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user } from '@/lib/auth'
import { CATEGORIES, formatBytes } from '@/lib/catalog'
import { ACCEPT, fileIcon, fileKind, isAcceptedName } from '@/lib/files'
import type { ContentFormat, Material, MaterialCategory, MaterialFile } from '@/lib/types'
import type { SelectOption } from '@/lib/ui'

const route = useRoute()
const router = useRouter()

const editId = computed<string | null>(() => {
  const id = route.params.id
  return typeof id === 'string' && id ? id : null
})
const isEdit = computed(() => !!editId.value)

interface MaterialForm {
  title: string
  category: MaterialCategory
  description: string
  condition: string
  location: string
  origin_date: string
}

function queryCategory(v: unknown): MaterialCategory {
  return CATEGORIES.find((c) => c.id === v)?.id ?? 'photos'
}

const form = ref<MaterialForm>({
  title: '',
  category: queryCategory(route.query.category),
  description: '',
  condition: '',
  location: '',
  origin_date: '',
})

const sourceType = ref<'file' | 'text'>('file')
const existingFiles = ref<MaterialFile[]>([])
const removedIds = ref<number[]>([])
const pendingFiles = ref<File[]>([])
const fileInput = ref<HTMLInputElement | null>(null)
const dragOver = ref(false)
const content = ref('')
const contentFormat = ref<ContentFormat>('html')
const hasContent = ref(false)
const saving = ref(false)
const loaded = ref(false)

const categoryOptions = computed<SelectOption[]>(() =>
  CATEGORIES.filter((c) => c.id !== 'documents' || user.value?.role === 'head_teacher').map(
    (c) => ({ value: c.id, label: c.label }),
  ),
)

const isTeacher = computed(() => user.value?.role === 'teacher')
const showFiles = computed(() => isEdit.value || sourceType.value === 'file')
const showEditor = computed(() => (isEdit.value ? hasContent.value : sourceType.value === 'text'))

onMounted(async () => {
  if (!isEdit.value) {
    if (!categoryOptions.value.some((o) => o.value === form.value.category)) {
      form.value.category = 'photos'
    }
    loaded.value = true
    return
  }
  try {
    const m = await api.get<Material>(`/materials/${editId.value}`)
    form.value = {
      title: m.title,
      category: m.category,
      description: m.description,
      condition: m.condition,
      location: m.location,
      origin_date: m.origin_date ? m.origin_date.slice(0, 10) : '',
    }
    existingFiles.value = m.files
    if (m.content != null) {
      hasContent.value = true
      content.value = m.content
      contentFormat.value = m.content_format || 'html'
    }
    loaded.value = true
  } catch (e) {
    notify.error('Не удалось загрузить материал', { text: errorMessage(e) })
    router.back()
  }
})

function pickFiles() {
  fileInput.value?.click()
}

function addFiles(list: FileList | null | undefined) {
  if (!list?.length) return
  const rejected: string[] = []
  for (const f of Array.from(list)) {
    if (isAcceptedName(f.name)) {
      pendingFiles.value.push(f)
    } else {
      rejected.push(f.name)
    }
  }
  if (rejected.length) {
    notify.error('Формат не поддерживается', {
      text: `Не добавлены: ${rejected.join(', ')}. Допустимые форматы: PDF, DOC(X), PPT(X), ODT, ODP, PNG, JPG, GIF, SVG, BMP, WEBP, MP4, MP3, WAV, WMV, WEBM, HTML, MD.`,
    })
  }
}

function onFilesChange(e: Event) {
  const input = e.target as HTMLInputElement
  addFiles(input.files)
  input.value = ''
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  addFiles(e.dataTransfer?.files)
}

function removePending(index: number) {
  pendingFiles.value.splice(index, 1)
}

function removeExisting(file: MaterialFile) {
  removedIds.value.push(file.id)
}

function restoreExisting(file: MaterialFile) {
  removedIds.value = removedIds.value.filter((id) => id !== file.id)
}

function isRemoved(file: MaterialFile): boolean {
  return removedIds.value.includes(file.id)
}

async function uploadPending(materialId: number) {
  if (!pendingFiles.value.length) return
  const fd = new FormData()
  for (const f of pendingFiles.value) fd.append('files', f)
  await api.postForm<Material>(`/materials/${materialId}/files`, fd)
}

async function saveExisting(): Promise<Material> {
  const body: Record<string, string> = { ...form.value }
  if (hasContent.value) {
    body.content = content.value
    body.content_format = contentFormat.value
  }
  const m = await api.put<Material>(`/materials/${editId.value}`, body)
  for (const id of removedIds.value) {
    await api.del(`/materials/${m.id}/files/${id}`)
  }
  await uploadPending(m.id)
  return m
}

async function createNew(): Promise<Material> {
  if (sourceType.value === 'text') {
    return api.post<Material>('/materials', {
      ...form.value,
      content: content.value,
      content_format: contentFormat.value,
    })
  }
  const fd = new FormData()
  for (const [k, v] of Object.entries(form.value)) fd.append(k, v ?? '')
  for (const f of pendingFiles.value) fd.append('files', f)
  return api.postForm<Material>('/materials', fd)
}

async function submit() {
  if (!form.value.title.trim()) {
    notify.warning('Укажите название экспоната')
    return
  }
  saving.value = true
  try {
    const m = isEdit.value ? await saveExisting() : await createNew()
    notify.success(
      isEdit.value
        ? 'Изменения сохранены'
        : isTeacher.value
          ? 'Отправлено на модерацию'
          : 'Экспонат добавлен',
    )
    router.push(`/material/${m.id}`)
  } catch (e) {
    notify.error('Не удалось сохранить', { text: errorMessage(e) })
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div v-if="loaded">
    <NPageHeader
      :title="isEdit ? 'Редактирование экспоната' : 'Новый экспонат'"
      :subtitle="isEdit ? 'Обновите сведения о материале' : 'Добавьте материал в цифровой архив школы'"
    >
      <template #actions>
        <NButton variant="ghost" icon="chevronLeft" @click="router.back()">Назад</NButton>
      </template>
    </NPageHeader>

    <NAlert v-if="isTeacher && !isEdit" variant="info" title="Премодерация" class="mb-6">
      Материал будет опубликован после проверки библиотекарем или завучем.
    </NAlert>

    <form class="grid grid-cols-1 xl:grid-cols-3 gap-6" @submit.prevent="submit">
      <div class="xl:col-span-2 space-y-6">
        <NCard title="Основные сведения" accent>
          <div class="space-y-4">
            <NInput v-model="form.title" label="Название" placeholder="Например: Кубок за победу в олимпиаде, 1998 г." required />
            <div class="grid sm:grid-cols-2 gap-4">
              <NSelect v-model="form.category" label="Категория" :options="categoryOptions" :disabled="isEdit" />
              <NDateInput v-model="form.origin_date" label="Время происхождения" hint="Если известно (у фотографий, наград и т.д.)" />
            </div>
            <NTextarea v-model="form.description" label="Описание" placeholder="Подробное описание материала, ключевые слова для поиска…" :rows="4" />
          </div>
        </NCard>

        <NCard
          :title="isEdit && !hasContent ? 'Файлы' : 'Содержимое'"
          :subtitle="isEdit ? 'Добавьте новые файлы или удалите лишние' : 'Загрузите файлы или создайте текстовый документ'"
        >
          <div v-if="!isEdit" class="flex gap-2 mb-4">
            <button
              type="button"
              class="flex-1 border p-3 text-sm font-semibold transition-colors ng-tile-press cursor-pointer flex items-center justify-center gap-2"
              :class="sourceType === 'file' ? 'border-accent text-ink bg-surface-2' : 'border-line text-muted hover:border-line-strong'"
              @click="sourceType = 'file'"
            >
              <NIcon name="upload" :size="16" /> Файлы
            </button>
            <button
              type="button"
              class="flex-1 border p-3 text-sm font-semibold transition-colors ng-tile-press cursor-pointer flex items-center justify-center gap-2"
              :class="sourceType === 'text' ? 'border-accent text-ink bg-surface-2' : 'border-line text-muted hover:border-line-strong'"
              @click="sourceType = 'text'"
            >
              <NIcon name="edit" :size="16" /> Текстовый документ
            </button>
          </div>

          <template v-if="showFiles">
            <div
              class="border-2 border-dashed p-8 text-center transition-colors cursor-pointer"
              :class="dragOver ? 'border-accent bg-surface-2' : 'border-line hover:border-line-strong'"
              @click="pickFiles"
              @dragover.prevent="dragOver = true"
              @dragleave="dragOver = false"
              @drop.prevent="onDrop"
            >
              <input ref="fileInput" type="file" multiple :accept="ACCEPT" class="hidden" @change="onFilesChange" />
              <NIcon name="upload" :size="32" class="text-faint inline-block" />
              <div class="font-bold text-ink mt-2">Перетащите файлы сюда или нажмите</div>
              <div class="text-xs text-muted mt-1">
                Можно выбрать сразу несколько: PDF, DOC(X), PPT(X), ODT, ODP, изображения, видео, аудио, HTML, MD
              </div>
            </div>

            <div v-if="existingFiles.length || pendingFiles.length" class="space-y-2 mt-4">
              <div
                v-for="f in existingFiles"
                :key="f.id"
                class="flex items-center gap-3 border px-3 py-2"
                :class="isRemoved(f) ? 'border-line bg-surface opacity-60' : 'border-line bg-surface-2'"
              >
                <NIcon :name="fileIcon(fileKind(f.name, f.mime))" :size="18" class="text-accent shrink-0" />
                <div class="min-w-0 flex-1">
                  <div class="text-sm font-semibold text-ink break-all" :class="isRemoved(f) ? 'line-through' : ''">
                    {{ f.name }}
                  </div>
                  <div class="text-xs text-muted">
                    {{ formatBytes(f.size) }}{{ isRemoved(f) ? ' · будет удалён при сохранении' : '' }}
                  </div>
                </div>
                <button
                  v-if="isRemoved(f)"
                  type="button"
                  class="text-muted hover:text-accent transition-colors cursor-pointer shrink-0"
                  title="Вернуть файл"
                  @click="restoreExisting(f)"
                >
                  <NIcon name="refresh" :size="16" />
                </button>
                <button
                  v-else
                  type="button"
                  class="text-muted hover:text-danger transition-colors cursor-pointer shrink-0"
                  title="Удалить файл"
                  @click="removeExisting(f)"
                >
                  <NIcon name="trash" :size="16" />
                </button>
              </div>

              <div
                v-for="(f, i) in pendingFiles"
                :key="`new-${i}-${f.name}`"
                class="flex items-center gap-3 bg-surface-2 border border-line border-l-[3px] border-l-success px-3 py-2"
              >
                <NIcon name="checkCircle" :size="18" class="text-success shrink-0" />
                <div class="min-w-0 flex-1">
                  <div class="text-sm font-semibold text-ink break-all">{{ f.name }}</div>
                  <div class="text-xs text-muted">{{ formatBytes(f.size) }} · новый файл</div>
                </div>
                <button
                  type="button"
                  class="text-muted hover:text-danger transition-colors cursor-pointer shrink-0"
                  title="Убрать из списка"
                  @click="removePending(i)"
                >
                  <NIcon name="close" :size="16" />
                </button>
              </div>
            </div>
          </template>

          <template v-if="showEditor">
            <TextEditor v-model:content="content" v-model:format="contentFormat" :class="showFiles ? 'mt-6' : ''" />
            <p v-if="!isEdit" class="text-xs text-muted mt-2">
              Текстовые документы сохраняются в базе данных и просматриваются прямо на портале. Они относятся к «Библиотеке».
            </p>
          </template>
        </NCard>
      </div>

      <div class="space-y-6">
        <NCard title="Физический экземпляр" subtitle="Для учёта оригинала">
          <div class="space-y-4">
            <NInput v-model="form.condition" label="Состояние" placeholder="Например: хорошее, требует реставрации" />
            <NInput v-model="form.location" label="Где находится" placeholder="Например: музей школы, шкаф 3, полка 2" />
          </div>
        </NCard>

        <NCard title="Сохранение">
          <p class="text-xs text-muted mb-4">
            После сохранения можно будет сгенерировать QR-код со ссылкой на карточку и наклеить его на физический экземпляр.
          </p>
          <div v-if="isEdit && (pendingFiles.length || removedIds.length)" class="text-xs text-muted mb-4 space-y-1">
            <div v-if="pendingFiles.length">Будет загружено файлов: {{ pendingFiles.length }}</div>
            <div v-if="removedIds.length">Будет удалено файлов: {{ removedIds.length }}</div>
          </div>
          <NButton type="submit" block size="lg" :loading="saving" icon="check">
            {{ isEdit ? 'Сохранить изменения' : 'Добавить в архив' }}
          </NButton>
        </NCard>
      </div>
    </form>
  </div>
  <div v-else class="text-center text-muted py-16">Загрузка…</div>
</template>
