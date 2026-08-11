<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    label?: string
    disabled?: boolean
  }>(),
  { modelValue: false, label: '', disabled: false },
)
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()

function toggle() {
  if (!props.disabled) emit('update:modelValue', !props.modelValue)
}
</script>

<template>
  <label
    class="flex w-fit items-center gap-3 select-none"
    :class="disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'"
    @click.prevent="toggle"
  >
    <span
      class="relative inline-flex items-center w-12 h-6 px-1 border transition-colors duration-150 shrink-0"
      :class="modelValue ? 'bg-accent border-accent' : 'bg-surface-2 border-line-strong'"
    >
      <span
        class="w-4 h-4 transition-transform duration-200 ease-out"
        :class="modelValue ? 'translate-x-6 bg-on-accent' : 'translate-x-0 bg-muted'"
      />
    </span>
    <span v-if="label" class="text-sm text-ink">{{ label }}</span>
  </label>
</template>
