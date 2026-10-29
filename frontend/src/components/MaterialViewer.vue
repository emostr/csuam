<script setup lang="ts">
import { computed, defineAsyncComponent, ref, watch } from 'vue'
import { DownloadIcon, FileTextIcon, PackageOpenIcon } from '@lucide/vue'
import type { LucideIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { fileIcon, fileKind, fileUrl } from '@/lib/files'
import type { FileKind } from '@/lib/files'
import type { Material, MaterialFile } from '@/lib/types'

const RichContent = defineAsyncComponent(() => import('@/components/editor/RichContent.vue'))

const props = defineProps<{ material: Material }>()
defineEmits<{ download: [file: MaterialFile] }>()

interface ViewItem {
  key: string
  label: string
  icon: LucideIcon
  kind: FileKind | 'content'
  file: MaterialFile | null
}

const items = computed<ViewItem[]>(() => {
  const list: ViewItem[] = []
  if (props.material.content != null) {
    list.push({ key: 'content', label: 'Текст', icon: FileTextIcon, kind: 'content', file: null })
  }
  for (const f of props.material.files) {
    const kind = fileKind(f.name, f.mime)
    list.push({ key: `file-${f.id}`, label: f.name, icon: fileIcon(kind), kind, file: f })
  }
  return list
})

const activeKey = ref('')

watch(
  items,
  (list) => {
    if (!list.some((i) => i.key === activeKey.value)) {
      activeKey.value = list[0]?.key ?? ''
    }
  },
  { immediate: true },
)

function selectTab(v: unknown) {
  if (typeof v === 'string') activeKey.value = v
}

const active = computed<ViewItem | null>(
  () => items.value.find((i) => i.key === activeKey.value) ?? null,
)

const activeUrl = computed(() =>
  active.value?.file ? fileUrl(props.material.id, active.value.file) : '',
)

const markdownText = ref<string | null>(null)
const markdownError = ref(false)

watch(
  active,
  async (item) => {
    markdownText.value = null
    markdownError.value = false
    if (!item || item.kind !== 'markdown' || !item.file) return
    try {
      const res = await fetch(fileUrl(props.material.id, item.file), { credentials: 'include' })
      if (!res.ok) throw new Error(String(res.status))
      markdownText.value = await res.text()
    } catch {
      markdownError.value = true
    }
  },
  { immediate: true },
)
</script>

<template>
  <Card class="gap-0 py-0">
    <Tabs v-if="items.length > 1" :model-value="activeKey" class="border-b" @update:model-value="selectTab">
      <div class="overflow-x-auto px-3 pt-2">
        <TabsList variant="line" class="h-10">
          <TabsTrigger v-for="i in items" :key="i.key" :value="i.key" class="max-w-56 flex-none">
            <component :is="i.icon" />
            <span class="truncate">{{ i.label }}</span>
          </TabsTrigger>
        </TabsList>
      </div>
    </Tabs>

    <RichContent
      v-if="active?.kind === 'content'"
      :content="material.content ?? ''"
      :format="material.content_format ?? 'html'"
      class="p-6"
    />

    <div v-else-if="active?.kind === 'image'" class="flex items-center justify-center bg-muted/40 p-4">
      <img :src="activeUrl" :alt="active.label" class="max-h-[70vh] max-w-full rounded-md object-contain" />
    </div>

    <video
      v-else-if="active?.kind === 'video'"
      :key="activeUrl"
      :src="activeUrl"
      controls
      class="max-h-[70vh] w-full bg-black"
    />

    <div v-else-if="active?.kind === 'audio'" class="p-6">
      <audio :key="activeUrl" :src="activeUrl" controls class="w-full" />
    </div>

    <iframe
      v-else-if="active?.kind === 'pdf'"
      :src="activeUrl"
      class="h-[75vh] w-full border-0"
      :title="active.label"
    />

    <iframe
      v-else-if="active?.kind === 'html'"
      :src="activeUrl"
      sandbox=""
      class="h-[70vh] w-full border-0 bg-white"
      :title="active.label"
    />

    <div v-else-if="active?.kind === 'markdown'" class="p-6">
      <p v-if="markdownError" class="text-sm text-destructive">Не удалось загрузить файл.</p>
      <RichContent v-else-if="markdownText != null" :content="markdownText" format="markdown" />
    </div>

    <Empty v-else-if="active">
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <FileTextIcon />
        </EmptyMedia>
        <EmptyTitle class="break-all">{{ active.label }}</EmptyTitle>
        <EmptyDescription>
          Предпросмотр этого формата в браузере недоступен. Скачайте файл, чтобы открыть его на устройстве.
        </EmptyDescription>
      </EmptyHeader>
      <EmptyContent>
        <Button @click="active.file && $emit('download', active.file)">
          <DownloadIcon />
          Скачать файл
        </Button>
      </EmptyContent>
    </Empty>

    <Empty v-else>
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <PackageOpenIcon />
        </EmptyMedia>
        <EmptyTitle>Нет файлов</EmptyTitle>
        <EmptyDescription>У этого экспоната нет прикреплённых файлов.</EmptyDescription>
      </EmptyHeader>
    </Empty>
  </Card>
</template>
