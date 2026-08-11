<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    name?: string
    src?: string
    size?: number
    square?: boolean
  }>(),
  { name: '', src: '', size: 40, square: true },
)

const initials = computed(() =>
  props.name
    .split(' ')
    .filter(Boolean)
    .slice(0, 2)
    .map((p) => p[0]?.toUpperCase())
    .join(''),
)
</script>

<template>
  <span
    class="inline-flex items-center justify-center overflow-hidden bg-accent text-on-accent font-bold shrink-0"
    :class="square ? '' : 'rounded-full'"
    :style="{ width: size + 'px', height: size + 'px', fontSize: size * 0.38 + 'px' }"
  >
    <img v-if="src" :src="src" :alt="name" class="w-full h-full object-cover" />
    <template v-else>{{ initials }}</template>
  </span>
</template>
