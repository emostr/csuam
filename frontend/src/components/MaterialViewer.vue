<script setup lang="ts">
import { computed, ref, watchEffect } from 'vue'
import NIcon from '@/components/ui/NIcon.vue'
import NButton from '@/components/ui/NButton.vue'
import { renderMarkdown } from '@/lib/markdown'
import type { Material } from '@/lib/types'

const props = defineProps<{ material: Material }>()
defineEmits<{ download: [] }>()

const fileUrl = computed(() => `/api/materials/${props.material.id}/file`)
const mime = computed(() => props.material.file_mime || '')
const ext = computed(() => {
  const name = props.material.file_name || ''
  const i = name.lastIndexOf('.')
  return i >= 0 ? name.slice(i + 1).toLowerCase() : ''
})

type ViewKind =
  | 'content'
  | 'empty'
  | 'image'
  | 'video'
  | 'audio'
  | 'pdf'
  | 'html'
  | 'markdown'
  | 'other'

const kind = computed<ViewKind>(() => {
  if (props.material.content != null) return 'content'
  if (!props.material.file_name) return 'empty'
  if (mime.value.startsWith('image/')) return 'image'
  if (mime.value.startsWith('video/')) return 'video'
  if (mime.value.startsWith('audio/')) return 'audio'
  if (mime.value === 'application/pdf' || ext.value === 'pdf') return 'pdf'
  if (mime.value === 'text/html' || ext.value === 'html') return 'html'
  if (ext.value === 'md' || mime.value === 'text/markdown') return 'markdown'
  return 'other'
})

const contentHtml = computed(() => {
  if (props.material.content == null) return ''
  if (props.material.content_format === 'markdown') {
    return renderMarkdown(props.material.content)
  }
  return props.material.content
})

const mdFileHtml = ref('')
watchEffect(async () => {
  if (kind.value !== 'markdown') return
  try {
    const res = await fetch(fileUrl.value, { credentials: 'include' })
    mdFileHtml.value = renderMarkdown(await res.text())
  } catch {
    mdFileHtml.value = '<p>Не удалось загрузить файл.</p>'
  }
})
</script>

<template>
  <div class="bg-surface border border-line">
    <div v-if="kind === 'content'" class="p-6 prose-archive" v-html="contentHtml" />

    <div v-else-if="kind === 'image'" class="flex items-center justify-center bg-surface-2 p-4">
      <img :src="fileUrl" :alt="material.title" class="max-w-full max-h-[70vh] object-contain" />
    </div>

    <video v-else-if="kind === 'video'" :src="fileUrl" controls class="w-full max-h-[70vh] bg-black" />

    <div v-else-if="kind === 'audio'" class="p-6">
      <audio :src="fileUrl" controls class="w-full" />
    </div>

    <iframe
      v-else-if="kind === 'pdf'"
      :src="fileUrl"
      class="w-full h-[75vh] border-0"
      :title="material.title"
    />

    <iframe
      v-else-if="kind === 'html'"
      :src="fileUrl"
      sandbox=""
      class="w-full h-[70vh] border-0 bg-white"
      :title="material.title"
    />

    <div v-else-if="kind === 'markdown'" class="p-6 prose-archive" v-html="mdFileHtml" />

    <div v-else-if="kind === 'other'" class="p-10 flex flex-col items-center text-center gap-3">
      <NIcon name="fileText" :size="44" class="text-faint" />
      <div class="font-bold text-ink">{{ material.file_name }}</div>
      <p class="text-sm text-muted max-w-sm">
        Предпросмотр этого формата в браузере недоступен. Скачайте файл, чтобы открыть его на устройстве.
      </p>
      <NButton icon="download" @click="$emit('download')">Скачать файл</NButton>
    </div>

    <div v-else class="p-10 flex flex-col items-center text-center gap-3">
      <NIcon name="box" :size="44" class="text-faint" />
      <p class="text-sm text-muted">У этого экспоната нет прикреплённого файла.</p>
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
