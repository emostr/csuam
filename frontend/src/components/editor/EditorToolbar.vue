<script setup lang="ts">
import { computed } from 'vue'
import type { Editor } from '@tiptap/core'
import {
  AlignCenterIcon,
  AlignJustifyIcon,
  AlignLeftIcon,
  AlignRightIcon,
  BoldIcon,
  ChevronDownIcon,
  CodeIcon,
  CodeXmlIcon,
  HighlighterIcon,
  ItalicIcon,
  LayoutTemplateIcon,
  LinkIcon,
  ListChecksIcon,
  ListIcon,
  ListOrderedIcon,
  MinusIcon,
  PilcrowIcon,
  QuoteIcon,
  Redo2Icon,
  SigmaIcon,
  SquareSigmaIcon,
  StrikethroughIcon,
  SubscriptIcon,
  SuperscriptIcon,
  TableIcon,
  UnderlineIcon,
  Undo2Icon,
  Heading1Icon,
  Heading2Icon,
  Heading3Icon,
} from '@lucide/vue'
import type { LucideIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import ToolbarButton from './ToolbarButton.vue'
import { notify } from '@/lib/notify'
import type { MathKind } from '@/lib/editor/extensions'

const props = defineProps<{ editor: Editor }>()
const emit = defineEmits<{
  formula: [kind: MathKind]
  templates: []
}>()

type Align = 'left' | 'center' | 'right' | 'justify'

interface BlockType {
  id: string
  label: string
  icon: LucideIcon
  active: () => boolean
  run: () => void
}

const chain = () => props.editor.chain().focus()

const blockTypes: BlockType[] = [
  {
    id: 'paragraph',
    label: 'Обычный текст',
    icon: PilcrowIcon,
    active: () => props.editor.isActive('paragraph'),
    run: () => chain().setParagraph().run(),
  },
  {
    id: 'h1',
    label: 'Заголовок 1',
    icon: Heading1Icon,
    active: () => props.editor.isActive('heading', { level: 1 }),
    run: () => chain().toggleHeading({ level: 1 }).run(),
  },
  {
    id: 'h2',
    label: 'Заголовок 2',
    icon: Heading2Icon,
    active: () => props.editor.isActive('heading', { level: 2 }),
    run: () => chain().toggleHeading({ level: 2 }).run(),
  },
  {
    id: 'h3',
    label: 'Заголовок 3',
    icon: Heading3Icon,
    active: () => props.editor.isActive('heading', { level: 3 }),
    run: () => chain().toggleHeading({ level: 3 }).run(),
  },
  {
    id: 'code',
    label: 'Блок кода',
    icon: CodeXmlIcon,
    active: () => props.editor.isActive('codeBlock'),
    run: () => chain().toggleCodeBlock().run(),
  },
]

const aligns: { id: Align; label: string; icon: LucideIcon }[] = [
  { id: 'left', label: 'По левому краю', icon: AlignLeftIcon },
  { id: 'center', label: 'По центру', icon: AlignCenterIcon },
  { id: 'right', label: 'По правому краю', icon: AlignRightIcon },
  { id: 'justify', label: 'По ширине', icon: AlignJustifyIcon },
]

const currentBlock = computed(() => blockTypes.find((b) => b.active()) ?? blockTypes[0])
const currentAlign = computed(
  () => aligns.find((a) => props.editor.isActive({ textAlign: a.id })) ?? aligns[0],
)
const inTable = computed(() => props.editor.isActive('table'))

const tableActions = computed(() => {
  const can = () => props.editor.can().chain().focus()
  return [
    { label: 'Строка выше', run: () => chain().addRowBefore().run(), enabled: can().addRowBefore().run() },
    { label: 'Строка ниже', run: () => chain().addRowAfter().run(), enabled: can().addRowAfter().run() },
    { label: 'Столбец слева', run: () => chain().addColumnBefore().run(), enabled: can().addColumnBefore().run() },
    { label: 'Столбец справа', run: () => chain().addColumnAfter().run(), enabled: can().addColumnAfter().run() },
    { label: 'Объединить ячейки', run: () => chain().mergeCells().run(), enabled: can().mergeCells().run() },
    { label: 'Разделить ячейку', run: () => chain().splitCell().run(), enabled: can().splitCell().run() },
    { label: 'Строка-заголовок', run: () => chain().toggleHeaderRow().run(), enabled: can().toggleHeaderRow().run() },
  ]
})

const tableRemovals = computed(() => {
  const can = () => props.editor.can().chain().focus()
  return [
    { label: 'Удалить строку', run: () => chain().deleteRow().run(), enabled: can().deleteRow().run() },
    { label: 'Удалить столбец', run: () => chain().deleteColumn().run(), enabled: can().deleteColumn().run() },
    { label: 'Удалить таблицу', run: () => chain().deleteTable().run(), enabled: can().deleteTable().run() },
  ]
})

function insertTable() {
  chain().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()
}

async function editLink() {
  const previous = props.editor.getAttributes('link').href
  const url = await notify.prompt({
    title: previous ? 'Изменить ссылку' : 'Вставить ссылку',
    inputLabel: 'Адрес',
    inputValue: typeof previous === 'string' ? previous : '',
    placeholder: 'https://…',
    confirmText: 'Применить',
  })
  if (url == null) return
  if (!url.trim()) {
    chain().extendMarkRange('link').unsetLink().run()
    return
  }
  if (props.editor.state.selection.empty && !props.editor.isActive('link')) {
    chain().insertContent({ type: 'text', text: url.trim(), marks: [{ type: 'link', attrs: { href: url.trim() } }] }).run()
    return
  }
  chain().extendMarkRange('link').setLink({ href: url.trim() }).run()
}
</script>

<template>
  <div class="flex flex-wrap items-center gap-0.5" role="toolbar" aria-label="Форматирование" @mousedown.prevent>
    <ToolbarButton
      :icon="Undo2Icon"
      label="Отменить"
      shortcut="Ctrl+Z"
      :disabled="!editor.can().undo()"
      @click="chain().undo().run()"
    />
    <ToolbarButton
      :icon="Redo2Icon"
      label="Повторить"
      shortcut="Ctrl+Shift+Z"
      :disabled="!editor.can().redo()"
      @click="chain().redo().run()"
    />

    <Separator orientation="vertical" class="mx-1 data-[orientation=vertical]:h-5 data-[orientation=vertical]:self-center" />

    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button type="button" variant="ghost" size="sm" class="w-44 justify-between font-normal">
          <span class="flex items-center gap-1.5 truncate">
            <component :is="currentBlock.icon" />
            {{ currentBlock.label }}
          </span>
          <ChevronDownIcon class="text-muted-foreground" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" @close-auto-focus.prevent class="w-48">
        <DropdownMenuItem v-for="b in blockTypes" :key="b.id" @select="b.run()">
          <component :is="b.icon" />
          {{ b.label }}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <Separator orientation="vertical" class="mx-1 data-[orientation=vertical]:h-5 data-[orientation=vertical]:self-center" />

    <ToolbarButton
      :icon="BoldIcon"
      label="Жирный"
      shortcut="Ctrl+B"
      :active="editor.isActive('bold')"
      @click="chain().toggleBold().run()"
    />
    <ToolbarButton
      :icon="ItalicIcon"
      label="Курсив"
      shortcut="Ctrl+I"
      :active="editor.isActive('italic')"
      @click="chain().toggleItalic().run()"
    />
    <ToolbarButton
      :icon="UnderlineIcon"
      label="Подчёркнутый"
      shortcut="Ctrl+U"
      :active="editor.isActive('underline')"
      @click="chain().toggleUnderline().run()"
    />
    <ToolbarButton
      :icon="StrikethroughIcon"
      label="Зачёркнутый"
      :active="editor.isActive('strike')"
      @click="chain().toggleStrike().run()"
    />
    <ToolbarButton
      :icon="HighlighterIcon"
      label="Выделение маркером"
      :active="editor.isActive('highlight')"
      @click="chain().toggleHighlight().run()"
    />
    <ToolbarButton
      :icon="CodeIcon"
      label="Код в строке"
      :active="editor.isActive('code')"
      @click="chain().toggleCode().run()"
    />
    <ToolbarButton
      :icon="SubscriptIcon"
      label="Нижний индекс"
      :active="editor.isActive('subscript')"
      @click="chain().toggleSubscript().run()"
    />
    <ToolbarButton
      :icon="SuperscriptIcon"
      label="Верхний индекс"
      :active="editor.isActive('superscript')"
      @click="chain().toggleSuperscript().run()"
    />

    <Separator orientation="vertical" class="mx-1 data-[orientation=vertical]:h-5 data-[orientation=vertical]:self-center" />

    <ToolbarButton
      :icon="ListIcon"
      label="Маркированный список"
      :active="editor.isActive('bulletList')"
      @click="chain().toggleBulletList().run()"
    />
    <ToolbarButton
      :icon="ListOrderedIcon"
      label="Нумерованный список"
      :active="editor.isActive('orderedList')"
      @click="chain().toggleOrderedList().run()"
    />
    <ToolbarButton
      :icon="ListChecksIcon"
      label="Список задач"
      :active="editor.isActive('taskList')"
      @click="chain().toggleTaskList().run()"
    />
    <ToolbarButton
      :icon="QuoteIcon"
      label="Цитата"
      :active="editor.isActive('blockquote')"
      @click="chain().toggleBlockquote().run()"
    />

    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button type="button" variant="ghost" size="sm" class="px-2" aria-label="Выравнивание">
          <component :is="currentAlign.icon" />
          <ChevronDownIcon class="size-3 text-muted-foreground" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" @close-auto-focus.prevent class="w-48">
        <DropdownMenuItem v-for="a in aligns" :key="a.id" @select="chain().setTextAlign(a.id).run()">
          <component :is="a.icon" />
          {{ a.label }}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <Separator orientation="vertical" class="mx-1 data-[orientation=vertical]:h-5 data-[orientation=vertical]:self-center" />

    <ToolbarButton :icon="LinkIcon" label="Ссылка" :active="editor.isActive('link')" @click="editLink" />
    <ToolbarButton :icon="SigmaIcon" label="Формула в строке" @click="emit('formula', 'inline')" />
    <ToolbarButton :icon="SquareSigmaIcon" label="Формула отдельной строкой" @click="emit('formula', 'block')" />
    <ToolbarButton :icon="MinusIcon" label="Разделитель" @click="chain().setHorizontalRule().run()" />

    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          class="px-2"
          :class="inTable ? 'bg-muted' : ''"
          aria-label="Таблица"
        >
          <TableIcon />
          <ChevronDownIcon class="size-3 text-muted-foreground" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" @close-auto-focus.prevent class="w-52">
        <DropdownMenuItem @select="insertTable">
          <TableIcon />
          Вставить таблицу 3×3
        </DropdownMenuItem>
        <template v-if="inTable">
          <DropdownMenuSeparator />
          <DropdownMenuLabel class="text-xs text-muted-foreground">Текущая таблица</DropdownMenuLabel>
          <DropdownMenuItem v-for="a in tableActions" :key="a.label" :disabled="!a.enabled" @select="a.run()">
            {{ a.label }}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            v-for="a in tableRemovals"
            :key="a.label"
            variant="destructive"
            :disabled="!a.enabled"
            @select="a.run()"
          >
            {{ a.label }}
          </DropdownMenuItem>
        </template>
      </DropdownMenuContent>
    </DropdownMenu>

    <Separator orientation="vertical" class="mx-1 data-[orientation=vertical]:h-5 data-[orientation=vertical]:self-center" />

    <Button type="button" variant="ghost" size="sm" @click="emit('templates')">
      <LayoutTemplateIcon />
      Шаблоны
    </Button>
  </div>
</template>
