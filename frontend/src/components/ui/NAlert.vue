<script setup lang="ts">
import { computed } from 'vue'
import NIcon from './NIcon.vue'

type Variant = 'info' | 'success' | 'warning' | 'danger'

const props = withDefaults(
  defineProps<{
    variant?: Variant
    title?: string
    closable?: boolean
  }>(),
  { variant: 'info', title: '', closable: false },
)
defineEmits<{ close: [] }>()

const map: Record<Variant, { border: string; icon: string; tint: string }> = {
  info: { border: 'border-l-info', icon: 'info', tint: 'text-info' },
  success: { border: 'border-l-success', icon: 'checkCircle', tint: 'text-success' },
  warning: { border: 'border-l-warning', icon: 'alert', tint: 'text-warning' },
  danger: { border: 'border-l-danger', icon: 'alert', tint: 'text-danger' },
}

const conf = computed(() => map[props.variant])
</script>

<template>
  <div class="flex items-start gap-3 bg-surface border border-line border-l-[3px] px-4 py-3" :class="conf.border">
    <NIcon :name="conf.icon" :size="20" :class="conf.tint" class="shrink-0 mt-0.5" />
    <div class="min-w-0 flex-1">
      <p v-if="title" class="font-bold text-ink text-sm">{{ title }}</p>
      <div class="text-sm text-muted" :class="title ? 'mt-0.5' : ''"><slot /></div>
    </div>
    <button
      v-if="closable"
      class="shrink-0 text-faint hover:text-danger transition-colors cursor-pointer"
      @click="$emit('close')"
    >
      <NIcon name="close" :size="16" />
    </button>
  </div>
</template>
