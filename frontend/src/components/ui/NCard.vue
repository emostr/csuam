<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    title?: string
    subtitle?: string
    accent?: boolean
    padding?: boolean
  }>(),
  { title: '', subtitle: '', accent: false, padding: true },
)

const rootClass = computed(() => [
  'bg-surface border border-line',
  props.accent ? 'border-l-[3px] border-l-accent' : '',
])
</script>

<template>
  <section :class="rootClass">
    <header
      v-if="title || subtitle || $slots.header || $slots.actions"
      class="flex items-start justify-between gap-4 px-5 py-4 border-b border-line"
    >
      <div v-if="title || subtitle || $slots.header" class="min-w-0">
        <slot name="header">
          <h3 class="text-[15px] font-bold text-ink truncate">{{ title }}</h3>
          <p v-if="subtitle" class="text-xs text-muted mt-0.5">{{ subtitle }}</p>
        </slot>
      </div>
      <div v-if="$slots.actions" class="shrink-0 flex items-center gap-2">
        <slot name="actions" />
      </div>
    </header>
    <div :class="padding ? 'p-5' : ''">
      <slot />
    </div>
    <footer v-if="$slots.footer" class="px-5 py-3 border-t border-line bg-surface-2/40">
      <slot name="footer" />
    </footer>
  </section>
</template>
