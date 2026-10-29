<script setup lang="ts">
import { watch } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import { createExtensions } from '@/lib/editor/extensions'
import type { ContentFormat } from '@/lib/types'

const props = withDefaults(defineProps<{ content: string; format?: ContentFormat }>(), { format: 'html' })

const contentType = () => (props.format === 'markdown' ? 'markdown' : 'html')

const editor = useEditor({
  editable: false,
  content: props.content,
  contentType: contentType(),
  extensions: createExtensions({ editable: false }),
  editorProps: { attributes: { class: 'prose-archive outline-none' } },
})

watch(
  () => [props.content, props.format],
  () => editor.value?.commands.setContent(props.content, { contentType: contentType(), emitUpdate: false }),
)
</script>

<template>
  <EditorContent :editor="editor" />
</template>
