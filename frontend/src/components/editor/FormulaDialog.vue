<script setup lang="ts">
import { computed, nextTick, ref, useId, watch } from 'vue'
import { Trash2Icon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Field, FieldDescription, FieldLabel } from '@/components/ui/field'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import KatexView from './KatexView.vue'
import { FORMULA_GROUPS } from '@/lib/editor/formulas'
import type { MathKind } from '@/lib/editor/extensions'

const props = defineProps<{
  open: boolean
  latex: string
  kind: MathKind
  editing: boolean
}>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  submit: [value: { latex: string; kind: MathKind }]
  remove: []
  closed: []
}>()

const value = ref('')
const mode = ref<MathKind>('inline')
const group = ref(FORMULA_GROUPS[0]?.id ?? '')
const inputId = useId()

watch(
  () => props.open,
  (open) => {
    if (!open) return
    value.value = props.latex
    mode.value = props.kind
  },
  { immediate: true },
)

const canSubmit = computed(() => value.value.trim() !== '')

function input(): HTMLTextAreaElement | null {
  const el = document.getElementById(inputId)
  return el instanceof HTMLTextAreaElement ? el : null
}

async function insertTemplate(latex: string) {
  const el = input()
  const start = el?.selectionStart ?? value.value.length
  const end = el?.selectionEnd ?? value.value.length
  const before = value.value.slice(0, start)
  const gap = before && !before.endsWith(' ') ? ' ' : ''
  value.value = before + gap + latex + value.value.slice(end)
  await nextTick()
  const caret = start + gap.length + latex.length
  el?.focus()
  el?.setSelectionRange(caret, caret)
}

function setMode(v: unknown) {
  if (v === 'inline' || v === 'block') mode.value = v
}

function setGroup(v: unknown) {
  if (typeof v === 'string') group.value = v
}

function submit() {
  if (!canSubmit.value) return
  emit('submit', { latex: value.value.trim(), kind: mode.value })
  emit('update:open', false)
}

function remove() {
  emit('remove')
  emit('update:open', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
    e.preventDefault()
    submit()
  }
}

function onCloseAutoFocus(e: Event) {
  e.preventDefault()
  emit('closed')
}

function focusInput(e: Event) {
  const el = input()
  if (!el) return
  e.preventDefault()
  el.focus()
  el.setSelectionRange(el.value.length, el.value.length)
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="max-h-[calc(100svh-2rem)] overflow-y-auto sm:max-w-3xl" @open-auto-focus="focusInput" @close-auto-focus="onCloseAutoFocus">
      <DialogHeader>
        <DialogTitle>{{ editing ? 'Изменить формулу' : 'Вставить формулу' }}</DialogTitle>
        <DialogDescription>
          Выберите заготовку или введите формулу в синтаксисе LaTeX. Химия — через <code>\ce{…}</code>.
        </DialogDescription>
      </DialogHeader>

      <Tabs :model-value="group" class="min-w-0" @update:model-value="setGroup">
        <div class="no-scrollbar -mx-1 overflow-x-auto px-1">
          <TabsList>
            <TabsTrigger v-for="g in FORMULA_GROUPS" :key="g.id" :value="g.id" class="flex-none">
              {{ g.label }}
            </TabsTrigger>
          </TabsList>
        </div>
        <TabsContent v-for="g in FORMULA_GROUPS" :key="g.id" :value="g.id">
          <div class="grid max-h-56 grid-cols-2 gap-2 overflow-y-auto rounded-lg border bg-muted/30 p-2 sm:grid-cols-4">
            <Tooltip v-for="f in g.items" :key="f.label">
              <TooltipTrigger as-child>
                <button
                  type="button"
                  class="flex h-14 min-w-0 items-center justify-center overflow-hidden rounded-md border bg-background px-2 text-sm transition-colors outline-none hover:border-primary/50 hover:bg-primary/5 focus-visible:ring-3 focus-visible:ring-ring/50"
                  :aria-label="f.label"
                  @click="insertTemplate(f.latex)"
                >
                  <KatexView :latex="f.latex" class="pointer-events-none max-w-full truncate" />
                </button>
              </TooltipTrigger>
              <TooltipContent>{{ f.label }}</TooltipContent>
            </Tooltip>
          </div>
        </TabsContent>
      </Tabs>

      <div class="grid gap-4 md:grid-cols-2">
        <Field>
          <FieldLabel :for="inputId">LaTeX</FieldLabel>
          <Textarea
            :id="inputId"
            v-model="value"
            class="min-h-28 font-mono text-sm"
            placeholder="\frac{a}{b}"
            spellcheck="false"
            autocomplete="off"
            @keydown="onKeydown"
          />
          <FieldDescription>Ctrl+Enter — вставить</FieldDescription>
        </Field>
        <Field>
          <FieldLabel>Предпросмотр</FieldLabel>
          <div
            class="flex min-h-28 items-center justify-center overflow-x-auto rounded-md border bg-muted/30 p-3"
            aria-live="polite"
          >
            <KatexView v-if="value.trim()" :latex="value" :display="mode === 'block'" />
            <span v-else class="text-sm text-muted-foreground">Здесь появится формула</span>
          </div>
          <Tabs v-if="!editing" :model-value="mode" @update:model-value="setMode">
            <TabsList class="grid w-full grid-cols-2">
              <TabsTrigger value="inline">В строке</TabsTrigger>
              <TabsTrigger value="block">Отдельной строкой</TabsTrigger>
            </TabsList>
          </Tabs>
        </Field>
      </div>

      <DialogFooter class="sm:justify-between">
        <Button v-if="editing" type="button" variant="ghost" class="text-destructive" @click="remove">
          <Trash2Icon />
          Удалить формулу
        </Button>
        <span v-else />
        <div class="flex flex-col-reverse gap-2 sm:flex-row">
          <Button type="button" variant="outline" @click="emit('update:open', false)">Отмена</Button>
          <Button type="button" :disabled="!canSubmit" @click="submit">
            {{ editing ? 'Сохранить' : 'Вставить' }}
          </Button>
        </div>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
