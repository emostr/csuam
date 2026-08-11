<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import NIcon from '@/components/ui/NIcon.vue'
import NBadge from '@/components/ui/NBadge.vue'
import { categoryMeta, STATUS_LABELS, STATUS_VARIANTS, formatDate } from '@/lib/catalog'
import type { Material } from '@/lib/types'

const props = defineProps<{ material: Material }>()

const router = useRouter()
const meta = computed(() => categoryMeta(props.material.category))
const isImage = computed(() => (props.material.file_mime || '').startsWith('image/'))
const previewUrl = computed(() => `/api/materials/${props.material.id}/file?thumb=1`)
</script>

<template>
  <button
    class="text-left bg-surface border border-line hover:border-accent transition-colors ng-tile-press cursor-pointer flex flex-col overflow-hidden"
    @click="router.push(`/material/${material.id}`)"
  >
    <div class="h-36 bg-surface-2 flex items-center justify-center overflow-hidden shrink-0">
      <img
        v-if="isImage"
        :src="previewUrl"
        :alt="material.title"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <NIcon v-else :name="meta.icon" :size="40" class="text-faint" />
    </div>
    <div class="p-4 flex-1 flex flex-col gap-2">
      <div class="flex items-start justify-between gap-2">
        <div class="font-bold text-ink text-sm leading-snug line-clamp-2">{{ material.title }}</div>
        <NBadge v-if="material.status !== 'approved'" :variant="STATUS_VARIANTS[material.status]">
          {{ STATUS_LABELS[material.status] }}
        </NBadge>
      </div>
      <p v-if="material.description" class="text-xs text-muted line-clamp-2">{{ material.description }}</p>
      <div class="mt-auto flex items-center gap-2 text-[11px] text-faint pt-1">
        <NIcon :name="meta.icon" :size="13" />
        <span>{{ meta.label }}</span>
        <span class="ml-auto">{{ formatDate(material.origin_date || material.created_at) }}</span>
      </div>
    </div>
  </button>
</template>
