<script setup>
import { computed, useId } from 'vue'
import NIcon from './NIcon.vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  label: { type: String, default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  hint: { type: String, default: '' },
  error: { type: String, default: '' },
  icon: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  required: { type: Boolean, default: false },
})
defineEmits(['update:modelValue'])

const uid = useId()
const fieldClass = computed(() => [
  'w-full h-11 bg-surface-2 text-ink text-sm border transition-colors',
  'placeholder:text-faint outline-none',
  props.icon ? 'pl-10 pr-3' : 'px-3',
  props.error
    ? 'border-danger focus:border-danger'
    : 'border-line focus:border-accent',
  props.disabled ? 'opacity-50 pointer-events-none' : '',
])
</script>

<template>
  <div>
    <label v-if="label" :for="uid" class="ng-label text-muted block mb-1.5">
      {{ label }} <span v-if="required" class="text-accent">*</span>
    </label>
    <div class="relative">
      <NIcon
        v-if="icon"
        :name="icon"
        :size="17"
        class="absolute left-3 top-1/2 -translate-y-1/2 text-faint pointer-events-none"
      />
      <input
        :id="uid"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="fieldClass"
        @input="$emit('update:modelValue', $event.target.value)"
      />
    </div>
    <p v-if="error" class="text-xs text-danger mt-1.5">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-faint mt-1.5">{{ hint }}</p>
  </div>
</template>
