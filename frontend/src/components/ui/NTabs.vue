<script setup lang="ts">
import { computed } from 'vue'
import NIcon from './NIcon.vue'

import type { TabItem } from '@/lib/ui'

const props = withDefaults(
  defineProps<{
    modelValue?: string | number
    tabs?: Array<TabItem | string | number>
  }>(),
  { modelValue: '', tabs: () => [] },
)
defineEmits<{ 'update:modelValue': [value: string | number] }>()

const normalized = computed<TabItem[]>(() =>
  props.tabs.map((t) => (typeof t === 'object' ? t : { value: t, label: String(t) })),
)
</script>

<template>
  <div class="flex items-stretch gap-1 border-b border-line overflow-x-auto">
    <button
      v-for="t in normalized"
      :key="t.value"
      class="relative inline-flex items-center gap-2 px-4 py-2.5 text-sm font-semibold whitespace-nowrap transition-colors cursor-pointer"
      :class="modelValue === t.value ? 'text-ink' : 'text-muted hover:text-ink'"
      @click="$emit('update:modelValue', t.value)"
    >
      <NIcon v-if="t.icon" :name="t.icon" :size="16" />
      {{ t.label }}
      <span v-if="t.badge != null" class="text-[10px] font-bold bg-surface-3 text-muted px-1.5 py-0.5">{{ t.badge }}</span>
      <span
        v-if="modelValue === t.value"
        class="absolute left-0 right-0 -bottom-px h-0.5 bg-accent"
      />
    </button>
  </div>
</template>
