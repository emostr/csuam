<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import {
  BoldIcon,
  CodeIcon,
  Heading2Icon,
  ItalicIcon,
  LinkIcon,
  ListIcon,
  ListOrderedIcon,
  MinusIcon,
  PencilIcon,
  QuoteIcon,
  StrikethroughIcon,
  UnderlineIcon,
} from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { notify } from '@/lib/notify'
import { renderMarkdown } from '@/lib/markdown'
import type { ContentFormat } from '@/lib/types'

const props = withDefaults(
  defineProps<{
    content?: string
    format?: ContentFormat
  }>(),
  { content: '', format: 'html' },
)
const emit = defineEmits<{
  'update:content': [value: string]
  'update:format': [value: ContentFormat]
}>()

const editor = ref<HTMLDivElement | null>(null)
const mdText = ref(props.format === 'markdown' ? props.content : '')

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

function onModeChange(v: unknown) {
  if (v === 'html' || v === 'markdown') switchMode(v)
}

async function switchMode(v: ContentFormat) {
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
  { icon: Heading2Icon, title: 'Заголовок', run: () => exec('formatBlock', 'h2') },
  { icon: BoldIcon, title: 'Жирный', run: () => exec('bold') },
  { icon: ItalicIcon, title: 'Курсив', run: () => exec('italic') },
  { icon: UnderlineIcon, title: 'Подчёркнутый', run: () => exec('underline') },
  { icon: StrikethroughIcon, title: 'Зачёркнутый', run: () => exec('strikeThrough') },
  { icon: ListIcon, title: 'Список', run: () => exec('insertUnorderedList') },
  { icon: ListOrderedIcon, title: 'Нумерованный список', run: () => exec('insertOrderedList') },
  { icon: QuoteIcon, title: 'Цитата', run: () => exec('formatBlock', 'blockquote') },
  { icon: LinkIcon, title: 'Ссылка', run: addLink },
  { icon: MinusIcon, title: 'Разделитель', run: () => exec('insertHorizontalRule') },
]

function exec(cmd: string, value?: string) {
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
  <div class="overflow-hidden rounded-lg border shadow-xs">
    <div class="flex flex-wrap items-center justify-between gap-2 border-b bg-muted/40 px-2 py-1.5">
      <div v-if="format === 'html'" class="flex flex-wrap items-center gap-0.5">
        <Tooltip v-for="btn in toolbar" :key="btn.title">
          <TooltipTrigger as-child>
            <Button type="button" variant="ghost" size="icon-sm" @mousedown.prevent @click="btn.run">
              <component :is="btn.icon" />
              <span class="sr-only">{{ btn.title }}</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>{{ btn.title }}</TooltipContent>
        </Tooltip>
      </div>
      <span v-else class="px-1.5 text-xs text-muted-foreground">Слева — разметка, справа — предпросмотр</span>

      <Tabs :model-value="format" class="ml-auto" @update:model-value="onModeChange">
        <TabsList>
          <TabsTrigger value="html">
            <PencilIcon />
            Визуальный
          </TabsTrigger>
          <TabsTrigger value="markdown">
            <CodeIcon />
            Markdown
          </TabsTrigger>
        </TabsList>
      </Tabs>
    </div>

    <div
      v-if="format === 'html'"
      ref="editor"
      contenteditable="true"
      class="prose-archive max-h-[60vh] min-h-[280px] overflow-y-auto p-4 text-sm outline-none"
      @input="onInput"
    />

    <div v-else class="grid divide-y md:grid-cols-2 md:divide-x md:divide-y-0">
      <textarea
        v-model="mdText"
        class="max-h-[60vh] min-h-[280px] resize-none bg-transparent p-4 font-mono text-sm outline-none placeholder:text-muted-foreground"
        placeholder="# Заголовок&#10;&#10;Текст в формате **Markdown**…"
        aria-label="Текст в формате Markdown"
      />
      <div class="prose-archive max-h-[60vh] min-h-[280px] overflow-y-auto p-4 text-sm" v-html="preview" />
    </div>
  </div>
</template>
