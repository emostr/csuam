<script setup lang="ts">
import { STACK } from '@/lib/stack'
import type { StackPath } from '@/lib/stack'
import { theme } from '@/lib/theme'

function brand(p: StackPath): string {
  return theme.value === 'dark' ? (p.colorDark ?? p.color) : p.color
}
</script>

<template>
  <div class="flex items-center gap-3 text-xs text-faint">
    <span>powered by</span>
    <ul class="flex items-center gap-2.5">
      <li v-for="item in STACK" :key="item.label">
        <svg
          class="stack-icon"
          :viewBox="item.viewBox"
          role="img"
          :aria-label="item.label"
          width="18"
          height="18"
        >
          <title>{{ item.label }}</title>
          <path
            v-for="(p, i) in item.paths"
            :key="i"
            :d="p.d"
            :fill-rule="p.rule"
            :clip-rule="p.rule"
            :style="{ '--brand': brand(p) }"
          />
        </svg>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.stack-icon {
  display: block;
  opacity: 0.75;
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}

.stack-icon path {
  fill: currentColor;
  transition: fill 0.15s ease;
}

.stack-icon:hover {
  opacity: 1;
  transform: translateY(-1px);
}

.stack-icon:hover path {
  fill: var(--brand);
}
</style>
