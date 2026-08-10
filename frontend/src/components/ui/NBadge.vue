<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: { type: String, default: 'neutral' },
  dot: { type: Boolean, default: false },
})

const variants = {
  neutral: 'bg-surface-3 text-muted',
  accent: 'bg-accent/15 text-accent',
  success: 'bg-success/15 text-success',
  warning: 'bg-warning/15 text-warning',
  danger: 'bg-danger/15 text-danger',
  info: 'bg-info/15 text-info',
}

const dotColor = {
  neutral: 'bg-muted',
  accent: 'bg-accent',
  success: 'bg-success',
  warning: 'bg-warning',
  danger: 'bg-danger',
  info: 'bg-info',
}

const cls = computed(() => variants[props.variant] || variants.neutral)
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 px-2 py-0.5 text-[11px] font-bold uppercase tracking-normal"
    :class="cls"
  >
    <span v-if="dot" class="w-1.5 h-1.5" :class="dotColor[variant]" />
    <slot />
  </span>
</template>
