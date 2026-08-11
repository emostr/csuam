<script setup lang="ts">
import { computed } from 'vue'

import type { BarPoint } from '@/lib/ui'

const props = withDefaults(
  defineProps<{
    data?: BarPoint[]
    height?: number
  }>(),
  { data: () => [], height: 200 },
)

const max = computed(() => Math.max(...props.data.map((d) => d.value), 1))
</script>

<template>
  <div class="flex items-end gap-2 sm:gap-3" :style="{ height: height + 'px' }">
    <div v-for="(d, i) in data" :key="i" class="flex-1 flex flex-col items-center justify-end h-full gap-2 group">
      <div class="text-[11px] font-bold text-muted opacity-0 group-hover:opacity-100 transition-opacity tabular-nums">{{ d.value }}</div>
      <div
        class="w-full bg-surface-3 group-hover:bg-accent transition-[height,background-color] duration-300"
        :class="d.active ? '!bg-accent' : ''"
        :style="{ height: Math.max(4, (d.value / max) * (height - 32)) + 'px' }"
      />
      <div class="text-[11px] text-faint truncate w-full text-center">{{ d.label }}</div>
    </div>
  </div>
</template>
