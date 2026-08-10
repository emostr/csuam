<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: Number, default: 0 },
  max: { type: Number, default: 100 },
  variant: { type: String, default: 'accent' },
  label: { type: String, default: '' },
  showValue: { type: Boolean, default: false },
})

const pct = computed(() => Math.max(0, Math.min(100, (props.value / props.max) * 100)))
const barColor = {
  accent: 'bg-accent',
  success: 'bg-success',
  warning: 'bg-warning',
  danger: 'bg-danger',
}
</script>

<template>
  <div>
    <div v-if="label || showValue" class="flex items-center justify-between mb-1.5">
      <span v-if="label" class="text-xs font-semibold text-muted">{{ label }}</span>
      <span v-if="showValue" class="text-xs font-bold text-ink tabular-nums">{{ Math.round(pct) }}%</span>
    </div>
    <div class="h-2 w-full bg-surface-3 overflow-hidden">
      <div
        class="h-full transition-[width] duration-500 ease-out"
        :class="barColor[variant] || barColor.accent"
        :style="{ width: pct + '%' }"
      />
    </div>
  </div>
</template>
