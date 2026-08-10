<script setup>
import NIcon from './NIcon.vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  label: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

function toggle() {
  if (!props.disabled) emit('update:modelValue', !props.modelValue)
}
</script>

<template>
  <label
    class="flex w-fit items-center gap-2.5 select-none"
    :class="disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'"
    @click.prevent="toggle"
  >
    <span
      class="w-5 h-5 border flex items-center justify-center transition-colors shrink-0"
      :class="modelValue ? 'bg-accent border-accent text-on-accent' : 'bg-surface-2 border-line-strong'"
    >
      <NIcon v-if="modelValue" name="check" :size="14" :stroke="3" />
    </span>
    <span v-if="label" class="text-sm text-ink">{{ label }}</span>
  </label>
</template>
