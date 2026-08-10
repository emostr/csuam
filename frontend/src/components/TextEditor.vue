<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import NIcon from '@/components/ui/NIcon.vue'
import NTabs from '@/components/ui/NTabs.vue'
import { notify } from '@/lib/notify'
import { renderMarkdown } from '@/lib/markdown'

const props = defineProps({
  content: { type: String, default: '' },
  format: { type: String, default: 'html' },
})
const emit = defineEmits(['update:content', 'update:format'])

const tabs = [
  { value: 'html', label: 'Визуальный', icon: 'edit' },
  { value: 'markdown', label: 'Markdown', icon: 'code' },
]

const editor = ref(null)
const mdText = ref(props.format === 'markdown' ? props.content : '')

const mode = computed({
  get: () => props.format,
  set: (v) => switchMode(v),
})

onMounted(() => {
  if (props.format === 'html' && editor.value) {
    editor.value.innerHTML = props.content
  }
})

watch(mdText, (v) => {
  if (props.format === 'markdown') emit('update:content', v)
})

function onInput() {
  emit('update:content', editor.value?.innerHTML || '')
}

async function switchMode(v) {
  if (v === props.format) return
  const hasContent = (props.content || '').trim() !== ''
  if (hasContent) {
    const ok = await notify.confirm({
      title: 'Сменить режим редактора?',
      text: 'Содержимое не конвертируется автоматически: текст останется как есть.',
      confirmText: 'Сменить',
      icon: 'info',
    })
    if (!ok) return
  }
  emit('update:format', v)
  if (v === 'markdown') {
    mdText.value = props.content
  } else {
    requestAnimationFrame(() => {
      if (editor.value) editor.value.innerHTML = props.content
    })
  }
}

const toolbar = [
  { icon: 'heading', title: 'Заголовок', run: () => exec('formatBlock', 'h2') },
  { icon: 'bold', title: 'Жирный', run: () => exec('bold') },
  { icon: 'italic', title: 'Курсив', run: () => exec('italic') },
  { icon: 'underline', title: 'Подчёркнутый', run: () => exec('underline') },
  { icon: 'strikethrough', title: 'Зачёркнутый', run: () => exec('strikeThrough') },
  { icon: 'list', title: 'Список', run: () => exec('insertUnorderedList') },
  { icon: 'listOrdered', title: 'Нумерованный список', run: () => exec('insertOrderedList') },
  { icon: 'quote', title: 'Цитата', run: () => exec('formatBlock', 'blockquote') },
  { icon: 'link', title: 'Ссылка', run: addLink },
  { icon: 'minus', title: 'Разделитель', run: () => exec('insertHorizontalRule') },
]

function exec(cmd, value = null) {
  editor.value?.focus()
  document.execCommand(cmd, false, value)
  onInput()
}

async function addLink() {
  const url = await notify.prompt({
    title: 'Вставить ссылку',
    placeholder: 'https://…',
  })
  if (url) exec('createLink', url)
}

const preview = computed(() => renderMarkdown(mdText.value))
</script>

<template>
  <div class="border border-line bg-surface">
    <div class="px-3 pt-2 border-b border-line bg-surface-2/40">
      <NTabs v-model="mode" :tabs="tabs" class="!border-0" />
    </div>

    <template v-if="mode === 'html'">
      <div class="flex flex-wrap items-center gap-0.5 px-2 py-1.5 border-b border-line">
        <button
          v-for="btn in toolbar"
          :key="btn.title"
          type="button"
          :title="btn.title"
          class="h-8 w-8 flex items-center justify-center text-muted hover:text-ink hover:bg-surface-2 transition-colors cursor-pointer"
          @mousedown.prevent
          @click="btn.run"
        >
          <NIcon :name="btn.icon" :size="16" />
        </button>
      </div>
      <div
        ref="editor"
        contenteditable="true"
        class="prose-archive min-h-[280px] max-h-[60vh] overflow-y-auto p-4 outline-none text-sm"
        @input="onInput"
      />
    </template>

    <template v-else>
      <div class="grid md:grid-cols-2 divide-y md:divide-y-0 md:divide-x divide-line">
        <textarea
          v-model="mdText"
          class="min-h-[280px] max-h-[60vh] p-4 bg-transparent text-sm text-ink font-mono outline-none resize-none placeholder:text-faint"
          placeholder="# Заголовок&#10;&#10;Текст в формате **Markdown**…"
        />
        <div class="prose-archive min-h-[280px] max-h-[60vh] overflow-y-auto p-4 text-sm" v-html="preview" />
      </div>
    </template>
  </div>
</template>
