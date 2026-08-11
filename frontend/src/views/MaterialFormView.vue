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
import type { ContentFormat, Material, MaterialCategory } from '@/lib/types'
import type { SelectOption } from '@/lib/ui'

const route = useRoute()
const router = useRouter()

const editId = computed<string | null>(() => {
  const id = route.params.id
  return typeof id === 'string' && id ? id : null
})
const isEdit = computed(() => !!editId.value)

const ACCEPT =
  '.pdf,.doc,.docx,.ppt,.pptx,.odt,.odp,.png,.jpg,.jpeg,.gif,.svg,.bmp,.webp,.mp4,.mp3,.wav,.wmv,.webm,.html,.md'

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

interface ExistingFile {
  name: string
  size: number | null
}

const sourceType = ref<'file' | 'text' | 'existing'>('file')
const file = ref<File | ExistingFile | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const dragOver = ref(false)
const content = ref('')
const contentFormat = ref<ContentFormat>('html')
const saving = ref(false)
const loaded = ref(false)

const categoryOptions = computed<SelectOption[]>(() =>
  CATEGORIES.filter((c) => c.id !== 'documents' || user.value?.role === 'head_teacher').map(
    (c) => ({ value: c.id, label: c.label }),
  ),
)

const isTeacher = computed(() => user.value?.role === 'teacher')

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
    if (m.content != null) {
      sourceType.value = 'text'
      content.value = m.content
      contentFormat.value = m.content_format || 'html'
    } else if (m.file_name) {
      sourceType.value = 'existing'
      file.value = { name: m.file_name, size: m.file_size }
    }
    loaded.value = true
  } catch (e) {
    notify.error('Не удалось загрузить материал', { text: errorMessage(e) })
    router.back()
  }
})

function pickFile() {
  fileInput.value?.click()
}

function validateFile(f: File): boolean {
  const ext = '.' + (f.name.split('.').pop() || '').toLowerCase()
  if (!ACCEPT.split(',').includes(ext)) {
    notify.error('Формат не поддерживается', {
      text: 'Допустимые форматы: PDF, DOC(X), PPT(X), ODT, ODP, PNG, JPG, GIF, SVG, BMP, WEBP, MP4, MP3, WAV, WMV, WEBM, HTML, MD.',
    })
    return false
  }
  return true
}

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files?.[0]
  input.value = ''
  if (f && validateFile(f)) file.value = f
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  const f = e.dataTransfer?.files?.[0]
  if (f && validateFile(f)) file.value = f
}

async function submit() {
  if (!form.value.title.trim()) {
    notify.warning('Укажите название экспоната')
    return
  }
  saving.value = true
  try {
    let m: Material
    if (isEdit.value) {
      const body: Record<string, string> = { ...form.value }
      if (sourceType.value === 'text') {
        body.content = content.value
        body.content_format = contentFormat.value
      }
      m = await api.put<Material>(`/materials/${editId.value}`, body)
      notify.success('Изменения сохранены')
    } else if (sourceType.value === 'text') {
      m = await api.post<Material>('/materials', {
        ...form.value,
        content: content.value,
        content_format: contentFormat.value,
      })
      notify.success(isTeacher.value ? 'Отправлено на модерацию' : 'Экспонат добавлен')
    } else {
      const fd = new FormData()
      for (const [k, v] of Object.entries(form.value)) fd.append(k, v ?? '')
      if (file.value instanceof File) fd.append('file', file.value)
      m = await api.postForm<Material>('/materials', fd)
      notify.success(isTeacher.value ? 'Отправлено на модерацию' : 'Экспонат добавлен')
    }
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

        <NCard v-if="!isEdit" title="Содержимое" subtitle="Загрузите файл или создайте текстовый документ">
          <div class="flex gap-2 mb-4">
            <button
              type="button"
              class="flex-1 border p-3 text-sm font-semibold transition-colors ng-tile-press cursor-pointer flex items-center justify-center gap-2"
              :class="sourceType === 'file' ? 'border-accent text-ink bg-surface-2' : 'border-line text-muted hover:border-line-strong'"
              @click="sourceType = 'file'"
            >
              <NIcon name="upload" :size="16" /> Файл
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

          <template v-if="sourceType === 'file'">
            <div
              class="border-2 border-dashed p-8 text-center transition-colors cursor-pointer"
              :class="dragOver ? 'border-accent bg-surface-2' : 'border-line hover:border-line-strong'"
              @click="pickFile"
              @dragover.prevent="dragOver = true"
              @dragleave="dragOver = false"
              @drop.prevent="onDrop"
            >
              <input ref="fileInput" type="file" :accept="ACCEPT" class="hidden" @change="onFileChange" />
              <template v-if="file">
                <NIcon name="checkCircle" :size="32" class="text-success inline-block" />
                <div class="font-bold text-ink mt-2 break-all">{{ file.name }}</div>
                <div class="text-xs text-muted mt-1">{{ formatBytes(file.size) }} · нажмите, чтобы заменить</div>
              </template>
              <template v-else>
                <NIcon name="upload" :size="32" class="text-faint inline-block" />
                <div class="font-bold text-ink mt-2">Перетащите файл сюда или нажмите</div>
                <div class="text-xs text-muted mt-1">
                  PDF, DOC(X), PPT(X), ODT, ODP, изображения, видео, аудио, HTML, MD
                </div>
              </template>
            </div>
          </template>

          <template v-else>
            <TextEditor v-model:content="content" v-model:format="contentFormat" />
            <p class="text-xs text-muted mt-2">
              Текстовые документы сохраняются в базе данных и просматриваются прямо на портале. Они относятся к «Библиотеке».
            </p>
          </template>
        </NCard>

        <NCard v-else-if="sourceType === 'text'" title="Текстовый документ">
          <TextEditor v-model:content="content" v-model:format="contentFormat" />
        </NCard>

        <NCard v-else-if="sourceType === 'existing'" title="Файл">
          <div class="flex items-center gap-3 text-sm">
            <NIcon name="fileText" :size="22" class="text-accent shrink-0" />
            <div>
              <div class="font-semibold text-ink break-all">{{ file?.name }}</div>
              <div class="text-xs text-muted">{{ formatBytes(file?.size) }} · замена файла не предусмотрена</div>
            </div>
          </div>
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
          <NButton type="submit" block size="lg" :loading="saving" icon="check">
            {{ isEdit ? 'Сохранить изменения' : 'Добавить в архив' }}
          </NButton>
        </NCard>
      </div>
    </form>
  </div>
  <div v-else class="text-center text-muted py-16">Загрузка…</div>
</template>
