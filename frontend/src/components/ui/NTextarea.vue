<script setup>
import { computed, useId } from 'vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  hint: { type: String, default: '' },
  error: { type: String, default: '' },
  rows: { type: Number, default: 4 },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue'])

const uid = useId()
const fieldClass = computed(() => [
  'w-full bg-surface-2 text-ink text-sm border px-3 py-2.5 resize-y transition-colors',
  'placeholder:text-faint outline-none',
  props.error ? 'border-danger focus:border-danger' : 'border-line focus:border-accent',
  props.disabled ? 'opacity-50 pointer-events-none' : '',
])
</script>

<template>
  <div>
    <label v-if="label" :for="uid" class="ng-label text-muted block mb-1.5">{{ label }}</label>
    <textarea
      :id="uid"
      :value="modelValue"
      :rows="rows"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="fieldClass"
      @input="$emit('update:modelValue', $event.target.value)"
    />
    <p v-if="error" class="text-xs text-danger mt-1.5">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-faint mt-1.5">{{ hint }}</p>
  </div>
</template>
