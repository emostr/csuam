<script setup lang="ts">
import { computed, useId } from 'vue'
import NIcon from './NIcon.vue'

import type { SelectOption } from '@/lib/ui'

const props = withDefaults(
  defineProps<{
    modelValue?: string | number
    label?: string
    options?: Array<SelectOption | string | number>
    placeholder?: string
    hint?: string
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    label: '',
    options: () => [],
    placeholder: 'Выберите…',
    hint: '',
    disabled: false,
  },
)
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

function onChange(e: Event) {
  emit('update:modelValue', (e.target as HTMLSelectElement).value)
}

const uid = useId()
const normalized = computed<SelectOption[]>(() =>
  props.options.map((o) => (typeof o === 'object' ? o : { value: o, label: String(o) })),
)
</script>

<template>
  <div>
    <label v-if="label" :for="uid" class="ng-label text-muted block mb-1.5">{{ label }}</label>
    <div class="relative">
      <select
        :id="uid"
        :value="modelValue"
        :disabled="disabled"
        class="w-full h-11 bg-surface-2 text-ink text-sm border border-line focus:border-accent outline-none px-3 pr-9 appearance-none cursor-pointer transition-colors disabled:opacity-50"
        @change="onChange"
      >
        <option value="" disabled>{{ placeholder }}</option>
        <option v-for="o in normalized" :key="o.value" :value="o.value">{{ o.label }}</option>
      </select>
      <NIcon
        name="chevronDown"
        :size="16"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-muted pointer-events-none"
      />
    </div>
    <p v-if="hint" class="text-xs text-faint mt-1.5">{{ hint }}</p>
  </div>
</template>
