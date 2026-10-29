<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { PaperclipIcon } from '@lucide/vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { categoryMeta, STATUS_LABELS, STATUS_TONES, formatDate } from '@/lib/catalog'
import type { Material } from '@/lib/types'

const props = defineProps<{ material: Material }>()

const meta = computed(() => categoryMeta(props.material.category))
const preview = computed(() => props.material.files.find((f) => f.mime.startsWith('image/')) ?? null)
const previewUrl = computed(() =>
  preview.value ? `/api/materials/${props.material.id}/files/${preview.value.id}?thumb=1` : '',
)
const fileCount = computed(() => props.material.files.length)
</script>

<template>
  <RouterLink
    :to="`/material/${material.id}`"
    class="group flex flex-col overflow-hidden rounded-xl bg-card text-card-foreground shadow-xs ring-1 ring-foreground/10 transition-shadow outline-none hover:shadow-md hover:ring-foreground/20 focus-visible:ring-3 focus-visible:ring-ring/50"
  >
    <div class="relative flex aspect-video items-center justify-center overflow-hidden bg-muted">
      <img
        v-if="preview"
        :src="previewUrl"
        :alt="material.title"
        class="size-full object-cover transition-transform duration-300 group-hover:scale-[1.03]"
        loading="lazy"
      />
      <component :is="meta.icon" v-else class="size-10 text-muted-foreground/60" />
      <span
        v-if="fileCount > 1"
        class="absolute right-2 bottom-2 flex items-center gap-1 rounded-md bg-background/90 px-1.5 py-0.5 text-xs font-medium shadow-xs backdrop-blur"
      >
        <PaperclipIcon class="size-3" />{{ fileCount }}
      </span>
    </div>
    <div class="flex flex-1 flex-col gap-2 p-4">
      <div class="flex items-start justify-between gap-2">
        <h3 class="line-clamp-2 leading-snug font-medium">{{ material.title }}</h3>
        <StatusBadge v-if="material.status !== 'approved'" :tone="STATUS_TONES[material.status]" class="shrink-0">
          {{ STATUS_LABELS[material.status] }}
        </StatusBadge>
      </div>
      <p v-if="material.description" class="line-clamp-2 text-sm text-muted-foreground">{{ material.description }}</p>
      <div class="mt-auto flex items-center gap-1.5 pt-1 text-xs text-muted-foreground">
        <component :is="meta.icon" class="size-3.5" />
        <span>{{ meta.label }}</span>
        <span class="ml-auto">{{ formatDate(material.origin_date || material.created_at) }}</span>
      </div>
    </div>
  </RouterLink>
</template>
