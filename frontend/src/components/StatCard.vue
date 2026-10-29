<script setup lang="ts">
import type { LucideIcon } from '@lucide/vue'
import { Card, CardAction, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

defineProps<{
  label: string
  value: string | number
  icon: LucideIcon
}>()
</script>

<template>
  <Card size="sm">
    <CardHeader>
      <CardDescription>{{ label }}</CardDescription>
      <CardTitle class="text-2xl font-semibold tabular-nums group-data-[size=sm]/card:text-2xl">{{ value }}</CardTitle>
      <CardAction>
        <span class="flex size-8 items-center justify-center rounded-md bg-primary/10 text-primary">
          <component :is="icon" class="size-4" />
        </span>
      </CardAction>
    </CardHeader>
  </Card>
</template>
