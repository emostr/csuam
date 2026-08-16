<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import NIcon from '@/components/ui/NIcon.vue'
import NButton from '@/components/ui/NButton.vue'
import NTabs from '@/components/ui/NTabs.vue'
import { renderMarkdown } from '@/lib/markdown'
import { fileIcon, fileKind, fileUrl } from '@/lib/files'
import type { FileKind } from '@/lib/files'
import type { Material, MaterialFile } from '@/lib/types'
import type { TabItem } from '@/lib/ui'

const props = defineProps<{ material: Material }>()
defineEmits<{ download: [file: MaterialFile] }>()

interface ViewItem {
  key: string
  label: string
  icon: string
  kind: FileKind | 'content'
  file: MaterialFile | null
}

const items = computed<ViewItem[]>(() => {
  const list: ViewItem[] = []
  if (props.material.content != null) {
    list.push({ key: 'content', label: 'Текст', icon: 'fileText', kind: 'content', file: null })
  }
  for (const f of props.material.files) {
    const kind = fileKind(f.name, f.mime)
    list.push({ key: `file-${f.id}`, label: f.name, icon: fileIcon(kind), kind, file: f })
  }
  return list
})

const tabs = computed<TabItem[]>(() =>
  items.value.map((i) => ({ value: i.key, label: i.label, icon: i.icon })),
)

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

const active = computed<ViewItem | null>(
  () => items.value.find((i) => i.key === activeKey.value) ?? null,
)

const activeUrl = computed(() =>
  active.value?.file ? fileUrl(props.material.id, active.value.file) : '',
)

const contentHtml = computed(() => {
  if (props.material.content == null) return ''
  if (props.material.content_format === 'markdown') {
    return renderMarkdown(props.material.content)
  }
  return props.material.content
})

const markdownHtml = ref('')

watch(
  active,
  async (item) => {
    if (!item || item.kind !== 'markdown' || !item.file) {
      markdownHtml.value = ''
      return
    }
    const url = fileUrl(props.material.id, item.file)
    markdownHtml.value = ''
    try {
      const res = await fetch(url, { credentials: 'include' })
      markdownHtml.value = renderMarkdown(await res.text())
    } catch {
      markdownHtml.value = '<p>Не удалось загрузить файл.</p>'
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="bg-surface border border-line">
    <NTabs v-if="items.length > 1" v-model="activeKey" :tabs="tabs" />

    <div v-if="active?.kind === 'content'" class="p-6 prose-archive" v-html="contentHtml" />

    <div v-else-if="active?.kind === 'image'" class="flex items-center justify-center bg-surface-2 p-4">
      <img :src="activeUrl" :alt="active.label" class="max-w-full max-h-[70vh] object-contain" />
    </div>

    <video
      v-else-if="active?.kind === 'video'"
      :key="activeUrl"
      :src="activeUrl"
      controls
      class="w-full max-h-[70vh] bg-black"
    />

    <div v-else-if="active?.kind === 'audio'" class="p-6">
      <audio :key="activeUrl" :src="activeUrl" controls class="w-full" />
    </div>

    <iframe
      v-else-if="active?.kind === 'pdf'"
      :src="activeUrl"
      class="w-full h-[75vh] border-0"
      :title="active.label"
    />

    <iframe
      v-else-if="active?.kind === 'html'"
      :src="activeUrl"
      sandbox=""
      class="w-full h-[70vh] border-0 bg-white"
      :title="active.label"
    />

    <div v-else-if="active?.kind === 'markdown'" class="p-6 prose-archive" v-html="markdownHtml" />

    <div v-else-if="active" class="p-10 flex flex-col items-center text-center gap-3">
      <NIcon name="fileText" :size="44" class="text-faint" />
      <div class="font-bold text-ink break-all">{{ active.label }}</div>
      <p class="text-sm text-muted max-w-sm">
        Предпросмотр этого формата в браузере недоступен. Скачайте файл, чтобы открыть его на устройстве.
      </p>
      <NButton icon="download" @click="active.file && $emit('download', active.file)">Скачать файл</NButton>
    </div>

    <div v-else class="p-10 flex flex-col items-center text-center gap-3">
      <NIcon name="box" :size="44" class="text-faint" />
      <p class="text-sm text-muted">У этого экспоната нет прикреплённых файлов.</p>
    </div>
  </div>
</template>

<style>
.prose-archive {
  color: var(--ng-ink);
  line-height: 1.7;
  overflow-wrap: break-word;
}
.prose-archive h1,
.prose-archive h2,
.prose-archive h3,
.prose-archive h4 {
  font-weight: 800;
  margin: 1.2em 0 0.4em;
  line-height: 1.25;
}
.prose-archive h1 { font-size: 1.6rem; }
.prose-archive h2 { font-size: 1.35rem; }
.prose-archive h3 { font-size: 1.15rem; }
.prose-archive p { margin: 0.5em 0; }
.prose-archive ul,
.prose-archive ol {
  margin: 0.5em 0;
  padding-left: 1.5em;
}
.prose-archive ul { list-style: square; }
.prose-archive ol { list-style: decimal; }
.prose-archive a {
  color: var(--ng-accent);
  text-decoration: underline;
}
.prose-archive blockquote {
  border-left: 3px solid var(--ng-accent);
  background: var(--ng-surface-2);
  padding: 0.5em 1em;
  margin: 0.75em 0;
  color: var(--ng-muted);
}
.prose-archive code {
  background: var(--ng-surface-3);
  padding: 0.1em 0.35em;
  font-size: 0.9em;
}
.prose-archive pre {
  background: var(--ng-surface-2);
  border: 1px solid var(--ng-line);
  padding: 0.9em 1em;
  overflow-x: auto;
  margin: 0.75em 0;
}
.prose-archive pre code {
  background: transparent;
  padding: 0;
}
.prose-archive hr {
  border: 0;
  border-top: 1px solid var(--ng-line);
  margin: 1.25em 0;
}
.prose-archive img {
  max-width: 100%;
}
.prose-archive s { color: var(--ng-faint); }
</style>
