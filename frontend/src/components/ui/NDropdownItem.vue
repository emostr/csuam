<script setup>
import NIcon from './NIcon.vue'

defineProps({
  icon: { type: String, default: '' },
  danger: { type: Boolean, default: false },
})
</script>

<template>
  <button
    class="w-full flex items-center gap-2.5 px-3.5 py-2 text-sm text-left transition-colors cursor-pointer"
    :class="danger ? 'text-danger hover:bg-danger/10' : 'text-ink hover:bg-surface-2'"
  >
    <NIcon v-if="icon" :name="icon" :size="16" :class="danger ? '' : 'text-muted'" />
    <slot />
  </button>
</template>
