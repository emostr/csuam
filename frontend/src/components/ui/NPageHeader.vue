<script setup>
defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
})
</script>

<template>
  <div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-6">
    <div>
      <div class="w-10 h-1 bg-accent mb-3" />
      <h1 class="text-2xl sm:text-3xl font-extrabold text-ink tracking-normal">{{ title }}</h1>
      <p v-if="subtitle" class="text-muted mt-1 text-sm">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="flex items-center gap-2 shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>
