<script setup lang="ts">
import { computed } from 'vue'
import { renderLatex } from '@/lib/editor/katex'

const props = withDefaults(defineProps<{ latex: string; display?: boolean }>(), { display: false })

const html = computed(() => renderLatex(props.latex, props.display))
</script>

<template>
  <span v-html="html" />
</template>
