<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import type { Tone } from '@/lib/catalog'

withDefaults(defineProps<{ tone?: Tone }>(), { tone: 'neutral' })

const dots: Record<Tone, string> = {
  neutral: 'bg-muted-foreground',
  primary: 'bg-primary',
  success: 'bg-success',
  warning: 'bg-warning',
  danger: 'bg-destructive',
  info: 'bg-info',
}
</script>

<template>
  <Badge variant="outline" class="gap-1.5">
    <span class="size-1.5 rounded-full" :class="dots[tone]" />
    <slot />
  </Badge>
</template>
