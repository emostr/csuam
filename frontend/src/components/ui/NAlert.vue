<script setup>
import { computed } from 'vue'
import NIcon from './NIcon.vue'

const props = defineProps({
  variant: { type: String, default: 'info' },
  title: { type: String, default: '' },
  closable: { type: Boolean, default: false },
})
defineEmits(['close'])

const map = {
  info: { border: 'border-l-info', icon: 'info', tint: 'text-info' },
  success: { border: 'border-l-success', icon: 'checkCircle', tint: 'text-success' },
  warning: { border: 'border-l-warning', icon: 'alert', tint: 'text-warning' },
  danger: { border: 'border-l-danger', icon: 'alert', tint: 'text-danger' },
}

const conf = computed(() => map[props.variant] || map.info)
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
