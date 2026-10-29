<script setup lang="ts">
import { computed, ref } from 'vue'
import { CalendarIcon } from '@lucide/vue'
import { CalendarDate, getLocalTimeZone, parseDate, today } from '@internationalized/date'
import type { DateValue } from '@internationalized/date'
import { createYearRange } from 'reka-ui/date'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    placeholder?: string
    id?: string
  }>(),
  { modelValue: '', placeholder: 'Выберите дату', id: undefined },
)
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const open = ref(false)

const date = computed<DateValue | undefined>(() => {
  if (!props.modelValue) return undefined
  try {
    return parseDate(props.modelValue.slice(0, 10))
  } catch {
    return undefined
  }
})

const label = computed(() =>
  date.value
    ? date.value
        .toDate(getLocalTimeZone())
        .toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
    : '',
)

const yearRange = createYearRange({
  start: new CalendarDate(1900, 1, 1),
  end: today(getLocalTimeZone()).add({ years: 5 }),
})

function select(v: DateValue | DateValue[] | undefined) {
  if (Array.isArray(v)) return
  emit('update:modelValue', v ? v.toString() : '')
  open.value = false
}
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <Button
        :id="id"
        variant="outline"
        :class="cn('w-full justify-start font-normal', !date && 'text-muted-foreground')"
      >
        <CalendarIcon />
        <span class="truncate">{{ label || placeholder }}</span>
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-auto p-0" align="start">
      <Calendar
        :model-value="date"
        :default-placeholder="date"
        layout="month-and-year"
        locale="ru-RU"
        :week-starts-on="1"
        :year-range="yearRange"
        initial-focus
        @update:model-value="select"
      />
      <div v-if="date" class="border-t p-2">
        <Button variant="ghost" size="sm" class="w-full" @click="select(undefined)">Очистить</Button>
      </div>
    </PopoverContent>
  </Popover>
</template>
