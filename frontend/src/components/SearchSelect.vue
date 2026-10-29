<script setup lang="ts">
import { computed } from 'vue'
import { CheckIcon, ChevronsUpDownIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  Combobox,
  ComboboxAnchor,
  ComboboxEmpty,
  ComboboxGroup,
  ComboboxInput,
  ComboboxItem,
  ComboboxItemIndicator,
  ComboboxList,
  ComboboxTrigger,
} from '@/components/ui/combobox'
import { cn } from '@/lib/utils'
import type { SelectOption } from '@/lib/ui'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    options: SelectOption[]
    placeholder?: string
    searchPlaceholder?: string
    emptyText?: string
    id?: string
  }>(),
  {
    modelValue: '',
    placeholder: 'Выберите…',
    searchPlaceholder: 'Поиск…',
    emptyText: 'Ничего не найдено',
    id: undefined,
  },
)
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const selected = computed(() => props.options.find((o) => o.value === props.modelValue) ?? null)

function onUpdate(v: unknown) {
  emit('update:modelValue', typeof v === 'string' ? v : '')
}
</script>

<template>
  <Combobox :model-value="modelValue || undefined" @update:model-value="onUpdate">
    <ComboboxAnchor as-child>
      <ComboboxTrigger as-child>
        <Button
          :id="id"
          variant="outline"
          role="combobox"
          :class="cn('w-full justify-between font-normal', !selected && 'text-muted-foreground')"
        >
          <span class="truncate">{{ selected?.label ?? placeholder }}</span>
          <ChevronsUpDownIcon class="opacity-50" />
        </Button>
      </ComboboxTrigger>
    </ComboboxAnchor>
    <ComboboxList align="start">
      <ComboboxInput :placeholder="searchPlaceholder" :display-value="() => ''" />
      <ComboboxEmpty>{{ emptyText }}</ComboboxEmpty>
      <ComboboxGroup class="max-h-60 overflow-y-auto p-1">
        <ComboboxItem v-for="o in options" :key="o.value" :value="o.value">
          <span class="truncate">{{ o.label }}</span>
          <ComboboxItemIndicator>
            <CheckIcon />
          </ComboboxItemIndicator>
        </ComboboxItem>
      </ComboboxGroup>
    </ComboboxList>
  </Combobox>
</template>
