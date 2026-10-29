<script setup lang="ts">
import type { LucideIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Toggle } from '@/components/ui/toggle'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

withDefaults(
  defineProps<{
    icon: LucideIcon
    label: string
    shortcut?: string
    active?: boolean
    disabled?: boolean
  }>(),
  { shortcut: '', active: undefined, disabled: false },
)
const emit = defineEmits<{ click: [] }>()
</script>

<template>
  <Tooltip>
    <TooltipTrigger as-child>
      <Toggle
        v-if="active !== undefined"
        size="sm"
        :model-value="active"
        :disabled="disabled"
        :aria-label="label"
        @update:model-value="emit('click')"
      >
        <component :is="icon" />
      </Toggle>
      <Button
        v-else
        type="button"
        variant="ghost"
        size="icon-sm"
        :disabled="disabled"
        :aria-label="label"
        @click="emit('click')"
      >
        <component :is="icon" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      {{ label }}<span v-if="shortcut" class="ml-2 opacity-60">{{ shortcut }}</span>
    </TooltipContent>
  </Tooltip>
</template>
