<script setup>
import { computed, useId } from 'vue'
import NIcon from './NIcon.vue'

const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  label: { type: String, default: '' },
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Выберите…' },
  hint: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue'])

const uid = useId()
const normalized = computed(() =>
  props.options.map((o) =>
    typeof o === 'object' ? o : { value: o, label: String(o) },
  ),
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
        @change="$emit('update:modelValue', $event.target.value)"
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
