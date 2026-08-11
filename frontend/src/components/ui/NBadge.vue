<script setup lang="ts">
import { computed } from 'vue'

type Variant = 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info'

const props = withDefaults(
  defineProps<{
    variant?: Variant | string
    dot?: boolean
  }>(),
  { variant: 'neutral', dot: false },
)

const variants: Record<Variant, string> = {
  neutral: 'bg-surface-3 text-muted',
  accent: 'bg-accent/15 text-accent',
  success: 'bg-success/15 text-success',
  warning: 'bg-warning/15 text-warning',
  danger: 'bg-danger/15 text-danger',
  info: 'bg-info/15 text-info',
}

const dotColor: Record<Variant, string> = {
  neutral: 'bg-muted',
  accent: 'bg-accent',
  success: 'bg-success',
  warning: 'bg-warning',
  danger: 'bg-danger',
  info: 'bg-info',
}

const isVariant = (v: string): v is Variant => v in variants
const key = computed<Variant>(() => (isVariant(props.variant) ? props.variant : 'neutral'))
const cls = computed(() => variants[key.value])
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2 py-0.5 text-[11px] font-bold uppercase tracking-normal"
    :class="cls"
  >
    <span v-if="dot" class="w-1.5 h-1.5" :class="dotColor[key]" />
    <slot />
  </span>
</template>
