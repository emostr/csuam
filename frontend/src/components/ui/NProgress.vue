<script setup lang="ts">
import { computed } from 'vue'

type Variant = 'accent' | 'success' | 'warning' | 'danger'

const props = withDefaults(
  defineProps<{
    value?: number
    max?: number
    variant?: Variant
    label?: string
    showValue?: boolean
  }>(),
  { value: 0, max: 100, variant: 'accent', label: '', showValue: false },
)

const pct = computed(() => Math.max(0, Math.min(100, (props.value / props.max) * 100)))
const barColor: Record<Variant, string> = {
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
        :class="barColor[variant]"
        :style="{ width: pct + '%' }"
      />
    </div>
  </div>
</template>
