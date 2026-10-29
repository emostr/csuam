<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { Selection } from '@tiptap/pm/state'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import type { Node as PMNode } from '@tiptap/pm/model'
import { CodeIcon, PencilIcon } from '@lucide/vue'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import EditorToolbar from './EditorToolbar.vue'
import FormulaDialog from './FormulaDialog.vue'
import TemplatesDialog from './TemplatesDialog.vue'
import RichContent from './RichContent.vue'
import { createExtensions } from '@/lib/editor/extensions'
import type { MathKind } from '@/lib/editor/extensions'
import type { DocTemplate } from '@/lib/editor/templates'
import { notify } from '@/lib/notify'
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

const formula = reactive({ open: false, latex: '', kind: 'inline' as MathKind, pos: null as number | null })
const templatesOpen = ref(false)
let focusAt: number | null = null
const mdText = ref(props.format === 'markdown' ? props.content : '')

function openMath(kind: MathKind, node: PMNode, pos: number) {
  formula.latex = String(node.attrs.latex ?? '')
  formula.kind = kind
  formula.pos = pos
  formula.open = true
}

const editor = useEditor({
  content: props.content,
  contentType: props.format === 'markdown' ? 'markdown' : 'html',
  extensions: createExtensions({
    editable: true,
    placeholder: 'Начните писать или выберите шаблон…',
    onMathClick: openMath,
  }),
  editorProps: {
    attributes: {
      class: 'prose-archive min-h-[320px] max-h-[70vh] overflow-y-auto px-4 py-3 text-sm outline-none',
      'aria-label': 'Текст документа',
    },
  },
  onUpdate: ({ editor: e }) => {
    if (props.format === 'html') emit('update:content', e.isEmpty ? '' : e.getHTML())
  },
})

watch(mdText, (v) => {
  if (props.format === 'markdown') emit('update:content', v)
})

const stats = computed(() => {
  const counter = editor.value?.storage.characterCount
  return counter ? { words: counter.words(), chars: counter.characters() } : { words: 0, chars: 0 }
})

function openFormula(kind: MathKind) {
  const e = editor.value
  if (!e) return
  const { from, to, empty } = e.state.selection
  formula.latex = empty ? '' : e.state.doc.textBetween(from, to, ' ')
  formula.kind = kind
  formula.pos = null
  formula.open = true
}

function submitFormula(v: { latex: string; kind: MathKind }) {
  const e = editor.value
  if (!e) return
  const pos = formula.pos
  if (pos != null) {
    focusAt = pos + (e.state.doc.nodeAt(pos)?.nodeSize ?? 1)
    if (formula.kind === 'inline') e.commands.updateInlineMath({ latex: v.latex, pos })
    else e.commands.updateBlockMath({ latex: v.latex, pos })
    return
  }
  if (v.kind === 'inline') e.commands.insertInlineMath({ latex: v.latex })
  else e.commands.insertBlockMath({ latex: v.latex })
}

function removeFormula() {
  const e = editor.value
  const pos = formula.pos
  if (!e || pos == null) return
  focusAt = pos
  if (formula.kind === 'inline') e.commands.deleteInlineMath({ pos })
  else e.commands.deleteBlockMath({ pos })
}

function applyTemplate(t: DocTemplate) {
  const e = editor.value
  if (!e) return
  if (e.isEmpty) e.commands.setContent(t.html, { emitUpdate: true })
  else e.commands.insertContent(t.html)
}

function focusEditor() {
  const e = editor.value
  if (!e) return
  const at = focusAt
  focusAt = null
  if (at == null) {
    e.commands.focus()
    return
  }
  e.chain()
    .command(({ tr }) => {
      tr.setSelection(Selection.near(tr.doc.resolve(Math.min(at, tr.doc.content.size))))
      return true
    })
    .focus()
    .run()
}

function onModeChange(v: unknown) {
  if (v === 'html' || v === 'markdown') switchMode(v)
}

async function switchMode(v: ContentFormat) {
  const e = editor.value
  if (!e || v === props.format) return
  if (v === 'markdown') {
    if (!e.isEmpty) {
      const ok = await notify.confirm({
        title: 'Перейти в Markdown?',
        text: 'Текст, формулы, таблицы и списки сохранятся. Выравнивание и индексы в Markdown не поддерживаются и будут потеряны.',
        confirmText: 'Перейти',
        icon: 'info',
      })
      if (!ok) return
    }
    mdText.value = e.isEmpty ? '' : e.getMarkdown()
    emit('update:format', 'markdown')
    emit('update:content', mdText.value)
    return
  }
  e.commands.setContent(mdText.value, { contentType: 'markdown', emitUpdate: false })
  emit('update:format', 'html')
  emit('update:content', e.isEmpty ? '' : e.getHTML())
}
</script>

<template>
  <div class="overflow-hidden rounded-lg border shadow-xs">
    <div class="border-b bg-muted/40 px-2 py-1.5">
      <EditorToolbar
        v-if="editor && format === 'html'"
        :editor="editor"
        @formula="openFormula"
        @templates="templatesOpen = true"
      />
      <p v-else class="px-1.5 py-1.5 text-xs text-muted-foreground">
        Слева — разметка Markdown, справа — предпросмотр. Формулы: <code>$…$</code> и <code>$$…$$</code>
      </p>
    </div>

    <EditorContent v-show="format === 'html'" :editor="editor" />

    <div v-if="format === 'markdown'" class="grid divide-y md:grid-cols-2 md:divide-x md:divide-y-0">
      <textarea
        v-model="mdText"
        class="max-h-[70vh] min-h-[320px] resize-none bg-transparent p-4 font-mono text-sm outline-none placeholder:text-muted-foreground"
        placeholder="# Заголовок&#10;&#10;Текст в формате **Markdown**, формула $E = mc^2$…"
        aria-label="Текст в формате Markdown"
        spellcheck="false"
      />
      <RichContent :content="mdText" format="markdown" class="max-h-[70vh] min-h-[320px] overflow-y-auto p-4 text-sm" />
    </div>

    <div class="flex flex-wrap items-center justify-between gap-2 border-t px-3 py-1.5 text-xs text-muted-foreground">
      <span v-if="format === 'html'" class="tabular-nums">
        Слов: {{ stats.words }} · Символов: {{ stats.chars }} · нажмите на формулу, чтобы изменить её
      </span>
      <span v-else>Режим Markdown</span>
      <Tabs :model-value="format" @update:model-value="onModeChange">
        <TabsList class="h-8">
          <TabsTrigger value="html" class="text-xs">
            <PencilIcon />
            Визуальный
          </TabsTrigger>
          <TabsTrigger value="markdown" class="text-xs">
            <CodeIcon />
            Markdown
          </TabsTrigger>
        </TabsList>
      </Tabs>
    </div>

    <FormulaDialog
      v-model:open="formula.open"
      :latex="formula.latex"
      :kind="formula.kind"
      :editing="formula.pos != null"
      @submit="submitFormula"
      @remove="removeFormula"
      @closed="focusEditor"
    />
    <TemplatesDialog v-model:open="templatesOpen" :empty="editor?.isEmpty ?? true" @apply="applyTemplate" @closed="focusEditor" />
  </div>
</template>
