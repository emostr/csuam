<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    points?: number[]
    width?: number
    height?: number
    area?: boolean
  }>(),
  { points: () => [], width: 100, height: 32, area: true },
)

const geo = computed(() => {
  const pts = props.points
  if (pts.length < 2) return { line: '', fill: '' }
  const min = Math.min(...pts)
  const max = Math.max(...pts)
  const span = max - min || 1
  const step = props.width / (pts.length - 1)
  const coords = pts.map((p, i) => {
    const x = i * step
    const y = props.height - ((p - min) / span) * (props.height - 4) - 2
    return [x, y]
  })
  const line = coords.map(([x, y], i) => `${i ? 'L' : 'M'}${x.toFixed(1)},${y.toFixed(1)}`).join(' ')
  const fill = `${line} L${props.width},${props.height} L0,${props.height} Z`
  return { line, fill }
})
</script>

<template>
  <svg :viewBox="`0 0 ${width} ${height}`" :width="width" :height="height" preserveAspectRatio="none" class="overflow-visible">
    <path v-if="area" :d="geo.fill" fill="var(--ng-accent)" opacity="0.14" />
    <path :d="geo.line" fill="none" stroke="var(--ng-accent)" stroke-width="2" vector-effect="non-scaling-stroke" />
  </svg>
</template>
