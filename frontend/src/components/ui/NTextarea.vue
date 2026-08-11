<script setup lang="ts">
import { computed, useId } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    label?: string
    placeholder?: string
    hint?: string
    error?: string
    rows?: number
    disabled?: boolean
  }>(),
  { modelValue: '', label: '', placeholder: '', hint: '', error: '', rows: 4, disabled: false },
)
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}

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
      @input="onInput"
    />
    <p v-if="error" class="text-xs text-danger mt-1.5">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-faint mt-1.5">{{ hint }}</p>
  </div>
</template>
