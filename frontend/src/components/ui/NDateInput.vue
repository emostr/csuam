<script setup lang="ts">
import { ref, useId } from 'vue'
import NIcon from './NIcon.vue'

withDefaults(
  defineProps<{
    modelValue?: string
    label?: string
    type?: string
    hint?: string
    disabled?: boolean
  }>(),
  { modelValue: '', label: '', type: 'date', hint: '', disabled: false },
)
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLInputElement).value)
}

const uid = useId()
const field = ref<HTMLInputElement | null>(null)

function openPicker() {
  try {
    field.value?.showPicker?.()
  } catch {}
}
</script>

<template>
  <div>
    <label v-if="label" :for="uid" class="ng-label text-muted block mb-1.5">{{ label }}</label>
    <div class="relative group">
      <input
        :id="uid"
        ref="field"
        :type="type"
        :value="modelValue"
        :disabled="disabled"
        class="ng-date w-full h-11 bg-surface-2 text-ink text-sm border border-line focus:border-accent outline-none pl-3 pr-10 transition-colors disabled:opacity-50 cursor-pointer"
        @input="onInput"
      />
      <button
        type="button"
        class="absolute right-0 top-0 h-11 w-10 flex items-center justify-center text-muted group-focus-within:text-accent hover:text-accent transition-colors cursor-pointer"
        tabindex="-1"
        @click="openPicker"
      >
        <NIcon name="calendar" :size="17" />
      </button>
    </div>
    <p v-if="hint" class="text-xs text-faint mt-1.5">{{ hint }}</p>
  </div>
</template>

<style scoped>
.ng-date::-webkit-calendar-picker-indicator {
  opacity: 0;
  position: absolute;
  right: 0;
  width: 2.5rem;
  height: 100%;
  cursor: pointer;
}
.ng-date::-webkit-datetime-edit {
  padding: 0;
}
</style>
