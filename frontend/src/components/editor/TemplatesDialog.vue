<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { DOC_TEMPLATES } from '@/lib/editor/templates'
import type { DocTemplate } from '@/lib/editor/templates'

defineProps<{ open: boolean; empty: boolean }>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  apply: [template: DocTemplate]
  closed: []
}>()

function onCloseAutoFocus(e: Event) {
  e.preventDefault()
  emit('closed')
}

function apply(t: DocTemplate) {
  emit('apply', t)
  emit('update:open', false)
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="max-h-[calc(100svh-2rem)] overflow-y-auto sm:max-w-3xl" @close-auto-focus="onCloseAutoFocus">
      <DialogHeader>
        <DialogTitle>Шаблоны документов</DialogTitle>
        <DialogDescription>
          {{
            empty
              ? 'Шаблон заполнит документ — останется вписать свои данные.'
              : 'Шаблон будет вставлен в позицию курсора, текущий текст сохранится.'
          }}
        </DialogDescription>
      </DialogHeader>
      <div class="grid gap-2 sm:grid-cols-2">
        <button
          v-for="t in DOC_TEMPLATES"
          :key="t.id"
          type="button"
          class="flex items-start gap-3 rounded-lg border p-3 text-left transition-colors outline-none hover:border-primary/50 hover:bg-primary/5 focus-visible:ring-3 focus-visible:ring-ring/50"
          @click="apply(t)"
        >
          <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
            <component :is="t.icon" class="size-4" />
          </span>
          <span class="grid gap-0.5">
            <span class="text-sm font-medium">{{ t.label }}</span>
            <span class="text-xs text-muted-foreground">{{ t.description }}</span>
          </span>
        </button>
      </div>
    </DialogContent>
  </Dialog>
</template>
