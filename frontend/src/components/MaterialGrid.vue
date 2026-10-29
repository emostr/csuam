<script setup lang="ts">
import { SearchXIcon } from '@lucide/vue'
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty'
import { Skeleton } from '@/components/ui/skeleton'
import MaterialCard from '@/components/MaterialCard.vue'
import type { Material } from '@/lib/types'

withDefaults(
  defineProps<{
    materials: Material[]
    loading: boolean
    emptyTitle?: string
    emptyText?: string
  }>(),
  { emptyTitle: 'Ничего не найдено', emptyText: '' },
)
</script>

<template>
  <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
    <div v-for="i in 8" :key="i" class="overflow-hidden rounded-xl ring-1 ring-foreground/10">
      <Skeleton class="aspect-video rounded-none" />
      <div class="space-y-2 p-4">
        <Skeleton class="h-4 w-3/4" />
        <Skeleton class="h-3 w-full" />
        <Skeleton class="h-3 w-1/2" />
      </div>
    </div>
  </div>

  <Empty v-else-if="!materials.length" class="border border-dashed">
    <EmptyHeader>
      <EmptyMedia variant="icon">
        <SearchXIcon />
      </EmptyMedia>
      <EmptyTitle>{{ emptyTitle }}</EmptyTitle>
      <EmptyDescription v-if="emptyText">{{ emptyText }}</EmptyDescription>
    </EmptyHeader>
    <slot name="empty" />
  </Empty>

  <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
    <MaterialCard v-for="m in materials" :key="m.id" :material="m" />
  </div>
</template>
