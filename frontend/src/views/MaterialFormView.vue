<script setup lang="ts">
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  CheckIcon,
  ChevronLeftIcon,
  FilePlusIcon,
  InfoIcon,
  PencilIcon,
  RotateCcwIcon,
  Trash2Icon,
  UploadIcon,
  XIcon,
} from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemMedia,
  ItemTitle,
} from '@/components/ui/item'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import PageHeader from '@/components/PageHeader.vue'
import DatePicker from '@/components/DatePicker.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user } from '@/lib/auth'
import { CATEGORIES, formatBytes } from '@/lib/catalog'
import { ACCEPT, fileIcon, fileKind, isAcceptedName } from '@/lib/files'
import type { ContentFormat, Material, MaterialCategory, MaterialFile } from '@/lib/types'

const RichTextEditor = defineAsyncComponent(() => import('@/components/editor/RichTextEditor.vue'))

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

const categories = computed(() =>
  CATEGORIES.filter((c) => c.id !== 'documents' || user.value?.role === 'head_teacher'),
)

const isTeacher = computed(() => user.value?.role === 'teacher')
const showFiles = computed(() => isEdit.value || sourceType.value === 'file')
const showEditor = computed(() => (isEdit.value ? hasContent.value : sourceType.value === 'text'))

function setCategory(v: unknown) {
  const c = CATEGORIES.find((x) => x.id === v)
  if (c) form.value.category = c.id
}

function setSourceType(v: unknown) {
  if (v === 'file' || v === 'text') sourceType.value = v
}

onMounted(async () => {
  if (!isEdit.value) {
    if (!categories.value.some((c) => c.id === form.value.category)) {
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
    <PageHeader
      :title="isEdit ? 'Редактирование экспоната' : 'Новый экспонат'"
      :description="isEdit ? 'Обновите сведения о материале' : 'Добавьте материал в цифровой архив школы'"
    >
      <template #actions>
        <Button variant="ghost" @click="router.back()">
          <ChevronLeftIcon />
          Назад
        </Button>
      </template>
    </PageHeader>

    <Alert v-if="isTeacher && !isEdit" class="mb-6">
      <InfoIcon />
      <AlertTitle>Премодерация</AlertTitle>
      <AlertDescription>Материал будет опубликован после проверки библиотекарем или завучем.</AlertDescription>
    </Alert>

    <form class="grid grid-cols-1 gap-6 xl:grid-cols-3" @submit.prevent="submit">
      <div class="space-y-6 xl:col-span-2">
        <Card>
          <CardHeader>
            <CardTitle>Основные сведения</CardTitle>
          </CardHeader>
          <CardContent>
            <FieldGroup>
              <Field>
                <FieldLabel for="material-title">
                  Название <span class="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="material-title"
                  v-model="form.title"
                  placeholder="Например: Кубок за победу в олимпиаде, 1998 г."
                  required
                />
              </Field>
              <div class="grid gap-6 sm:grid-cols-2">
                <Field>
                  <FieldLabel for="material-category">Категория</FieldLabel>
                  <Select :model-value="form.category" :disabled="isEdit" @update:model-value="setCategory">
                    <SelectTrigger id="material-category" class="w-full">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="c in categories" :key="c.id" :value="c.id">
                        <component :is="c.icon" />
                        {{ c.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </Field>
                <Field>
                  <FieldLabel for="material-origin">Время происхождения</FieldLabel>
                  <DatePicker id="material-origin" v-model="form.origin_date" />
                  <FieldDescription>Если известно (у фотографий, наград и т.д.)</FieldDescription>
                </Field>
              </div>
              <Field>
                <FieldLabel for="material-description">Описание</FieldLabel>
                <Textarea
                  id="material-description"
                  v-model="form.description"
                  placeholder="Подробное описание материала, ключевые слова для поиска…"
                  :rows="4"
                />
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{{ isEdit && !hasContent ? 'Файлы' : 'Содержимое' }}</CardTitle>
            <CardDescription>
              {{ isEdit ? 'Добавьте новые файлы или удалите лишние' : 'Загрузите файлы или создайте текстовый документ' }}
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <Tabs v-if="!isEdit" :model-value="sourceType" @update:model-value="setSourceType">
              <TabsList class="grid w-full grid-cols-2">
                <TabsTrigger value="file">
                  <UploadIcon />
                  Файлы
                </TabsTrigger>
                <TabsTrigger value="text">
                  <PencilIcon />
                  Текстовый документ
                </TabsTrigger>
              </TabsList>
            </Tabs>

            <template v-if="showFiles">
              <button
                type="button"
                class="flex w-full flex-col items-center gap-2 rounded-lg border-2 border-dashed p-8 text-center transition-colors outline-none hover:bg-muted/40 focus-visible:ring-3 focus-visible:ring-ring/50"
                :class="dragOver ? 'border-primary bg-primary/5' : 'border-border'"
                @click="pickFiles"
                @dragover.prevent="dragOver = true"
                @dragleave="dragOver = false"
                @drop.prevent="onDrop"
              >
                <span class="flex size-10 items-center justify-center rounded-lg bg-muted">
                  <FilePlusIcon class="size-5 text-muted-foreground" />
                </span>
                <span class="font-medium">Перетащите файлы сюда или нажмите</span>
                <span class="text-xs text-muted-foreground">
                  Можно выбрать сразу несколько: PDF, DOC(X), PPT(X), ODT, ODP, изображения, видео, аудио, HTML, MD
                </span>
              </button>
              <input ref="fileInput" type="file" multiple :accept="ACCEPT" class="hidden" @change="onFilesChange" />

              <ItemGroup v-if="existingFiles.length || pendingFiles.length" class="gap-2">
                <Item
                  v-for="f in existingFiles"
                  :key="f.id"
                  variant="outline"
                  size="sm"
                  :class="isRemoved(f) ? 'opacity-60' : ''"
                >
                  <ItemMedia variant="icon">
                    <component :is="fileIcon(fileKind(f.name, f.mime))" class="text-primary" />
                  </ItemMedia>
                  <ItemContent class="min-w-0">
                    <ItemTitle class="break-all" :class="isRemoved(f) ? 'line-through' : ''">{{ f.name }}</ItemTitle>
                    <ItemDescription>
                      {{ formatBytes(f.size) }}{{ isRemoved(f) ? ' · будет удалён при сохранении' : '' }}
                    </ItemDescription>
                  </ItemContent>
                  <ItemActions>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon-sm"
                          @click="isRemoved(f) ? restoreExisting(f) : removeExisting(f)"
                        >
                          <RotateCcwIcon v-if="isRemoved(f)" />
                          <Trash2Icon v-else />
                          <span class="sr-only">{{ isRemoved(f) ? 'Вернуть файл' : 'Удалить файл' }}</span>
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{{ isRemoved(f) ? 'Вернуть файл' : 'Удалить файл' }}</TooltipContent>
                    </Tooltip>
                  </ItemActions>
                </Item>

                <Item
                  v-for="(f, i) in pendingFiles"
                  :key="`new-${i}-${f.name}`"
                  variant="outline"
                  size="sm"
                  class="border-success/40 bg-success/5"
                >
                  <ItemMedia variant="icon">
                    <component :is="fileIcon(fileKind(f.name, f.type))" class="text-success" />
                  </ItemMedia>
                  <ItemContent class="min-w-0">
                    <ItemTitle class="break-all">{{ f.name }}</ItemTitle>
                    <ItemDescription>{{ formatBytes(f.size) }} · новый файл</ItemDescription>
                  </ItemContent>
                  <ItemActions>
                    <Tooltip>
                      <TooltipTrigger as-child>
                        <Button type="button" variant="ghost" size="icon-sm" @click="removePending(i)">
                          <XIcon />
                          <span class="sr-only">Убрать из списка</span>
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>Убрать из списка</TooltipContent>
                    </Tooltip>
                  </ItemActions>
                </Item>
              </ItemGroup>
            </template>

            <template v-if="showEditor">
              <RichTextEditor v-model:content="content" v-model:format="contentFormat" :class="showFiles ? 'mt-2' : ''" />
              <p v-if="!isEdit" class="text-xs text-muted-foreground">
                Текстовые документы сохраняются в базе данных и просматриваются прямо на портале. Они относятся к «Библиотеке».
              </p>
            </template>
          </CardContent>
        </Card>
      </div>

      <div class="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Физический экземпляр</CardTitle>
            <CardDescription>Для учёта оригинала</CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup>
              <Field>
                <FieldLabel for="material-condition">Состояние</FieldLabel>
                <Input
                  id="material-condition"
                  v-model="form.condition"
                  placeholder="Например: хорошее, требует реставрации"
                />
              </Field>
              <Field>
                <FieldLabel for="material-location">Где находится</FieldLabel>
                <Input
                  id="material-location"
                  v-model="form.location"
                  placeholder="Например: музей школы, шкаф 3, полка 2"
                />
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>

        <Card class="xl:sticky xl:top-20">
          <CardHeader>
            <CardTitle>Сохранение</CardTitle>
            <CardDescription>
              После сохранения можно будет сгенерировать QR-код со ссылкой на карточку и наклеить его на физический
              экземпляр.
            </CardDescription>
          </CardHeader>
          <CardContent class="space-y-4">
            <div
              v-if="isEdit && (pendingFiles.length || removedIds.length)"
              class="space-y-1 text-xs text-muted-foreground"
            >
              <p v-if="pendingFiles.length">Будет загружено файлов: {{ pendingFiles.length }}</p>
              <p v-if="removedIds.length">Будет удалено файлов: {{ removedIds.length }}</p>
            </div>
            <Button type="submit" size="lg" class="w-full" :disabled="saving">
              <Spinner v-if="saving" />
              <CheckIcon v-else />
              {{ isEdit ? 'Сохранить изменения' : 'Добавить в архив' }}
            </Button>
          </CardContent>
        </Card>
      </div>
    </form>
  </div>

  <div v-else class="space-y-6">
    <div class="space-y-2">
      <Skeleton class="h-8 w-72 max-w-full" />
      <Skeleton class="h-4 w-56" />
    </div>
    <div class="grid grid-cols-1 gap-6 xl:grid-cols-3">
      <Skeleton class="h-96 rounded-xl xl:col-span-2" />
      <Skeleton class="h-64 rounded-xl" />
    </div>
  </div>
</template>
