<script setup>
import { computed } from 'vue'
import NIcon from './NIcon.vue'

const props = defineProps({
  label: { type: String, default: '' },
  value: { type: [String, Number], default: '' },
  icon: { type: String, default: '' },
  delta: { type: [String, Number], default: null },
  trend: { type: String, default: '' },
  solid: { type: Boolean, default: false },
})

const trendUp = computed(() => props.trend === 'up')
const trendClass = computed(() =>
  props.solid
    ? 'text-on-accent/85'
    : trendUp.value
      ? 'text-success'
      : 'text-danger',
)
</script>

<template>
  <div
    class="relative overflow-hidden ng-tile-press border p-5 flex flex-col justify-between min-h-[128px]"
    :class="solid
      ? 'bg-accent text-on-accent border-transparent'
      : 'bg-surface border-line hover:border-line-strong'"
  >
    <div class="flex items-start justify-between">
      <span
        class="ng-label"
        :class="solid ? 'text-on-accent/80' : 'text-muted'"
      >{{ label }}</span>
      <NIcon
        v-if="icon"
        :name="icon"
        :size="22"
        :class="solid ? 'text-on-accent/70' : 'text-accent'"
      />
    </div>

    <div>
      <div class="text-3xl font-extrabold tracking-normal leading-none">{{ value }}</div>
      <div v-if="delta !== null" class="flex items-center gap-1.5 mt-2 text-xs font-semibold" :class="trendClass">
        <NIcon :name="trendUp ? 'trendUp' : 'trendDown'" :size="14" />
        <span>{{ delta }}</span>
        <span :class="solid ? 'text-on-accent/60' : 'text-faint'" class="font-normal">за 7 дней</span>
      </div>
    </div>
  </div>
</template>
