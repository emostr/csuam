<script setup lang="ts">
import NIcon from './NIcon.vue'

withDefaults(
  defineProps<{
    icon?: string
    danger?: boolean
  }>(),
  { icon: '', danger: false },
)
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
