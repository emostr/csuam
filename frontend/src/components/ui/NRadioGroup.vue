<script setup lang="ts">
import { computed } from 'vue'

interface RadioOption {
  value: string | number
  label: string
}

const props = withDefaults(
  defineProps<{
    modelValue?: string | number
    options?: Array<RadioOption | string | number>
    label?: string
    inline?: boolean
  }>(),
  { modelValue: '', options: () => [], label: '', inline: false },
)
const emit = defineEmits<{ 'update:modelValue': [value: string | number] }>()

const normalized = computed<RadioOption[]>(() =>
  props.options.map((o) => (typeof o === 'object' ? o : { value: o, label: String(o) })),
)
</script>

<template>
  <div>
    <span v-if="label" class="ng-label text-muted block mb-2">{{ label }}</span>
    <div :class="inline ? 'flex flex-wrap gap-x-6 gap-y-2' : 'flex flex-col gap-2'">
      <label
        v-for="o in normalized"
        :key="o.value"
        class="inline-flex items-center gap-2.5 cursor-pointer select-none"
        @click.prevent="emit('update:modelValue', o.value)"
      >
        <span
          class="w-5 h-5 border flex items-center justify-center transition-colors shrink-0"
          :class="modelValue === o.value ? 'border-accent' : 'border-line-strong'"
        >
          <span v-if="modelValue === o.value" class="w-2.5 h-2.5 bg-accent" />
        </span>
        <span class="text-sm text-ink">{{ o.label }}</span>
      </label>
    </div>
  </div>
</template>
