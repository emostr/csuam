<script setup lang="ts">
import { computed } from 'vue'
import NIcon from './NIcon.vue'

type Variant = 'primary' | 'secondary' | 'subtle' | 'ghost' | 'danger'
type Size = 'sm' | 'md' | 'lg'

const props = withDefaults(
  defineProps<{
    variant?: Variant
    size?: Size
    icon?: string
    iconRight?: string
    block?: boolean
    loading?: boolean
    disabled?: boolean
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    variant: 'primary',
    size: 'md',
    icon: '',
    iconRight: '',
    block: false,
    loading: false,
    disabled: false,
    type: 'button',
  },
)

const variants: Record<Variant, string> = {
  primary:
    'bg-accent text-on-accent hover:brightness-110 active:brightness-95 border border-transparent',
  secondary:
    'bg-transparent text-ink border border-line-strong hover:border-accent hover:text-accent',
  subtle:
    'bg-surface-2 text-ink border border-transparent hover:bg-surface-3',
  ghost:
    'bg-transparent text-muted border border-transparent hover:text-ink hover:bg-surface-2',
  danger:
    'bg-danger text-white border border-transparent hover:brightness-110 active:brightness-95',
}

const sizes: Record<Size, string> = {
  sm: 'h-8 px-3 text-xs gap-1.5',
  md: 'h-10 px-4 text-sm gap-2',
  lg: 'h-12 px-6 text-base gap-2.5',
}

const iconSizes: Record<Size, number> = { sm: 15, md: 17, lg: 19 }

const classes = computed(() => [
  'inline-flex items-center justify-center font-bold tracking-normal select-none',
  'transition-[filter,background-color,border-color,color] duration-150',
  'disabled:opacity-45 disabled:pointer-events-none whitespace-nowrap',
  'ng-tile-press cursor-pointer',
  variants[props.variant],
  sizes[props.size],
  props.block ? 'w-full' : '',
])
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled || loading">
    <NIcon v-if="loading" name="refresh" :size="iconSizes[size]" class="animate-spin" />
    <NIcon v-else-if="icon" :name="icon" :size="iconSizes[size]" />
    <span v-if="$slots.default"><slot /></span>
    <NIcon v-if="iconRight && !loading" :name="iconRight" :size="iconSizes[size]" />
  </button>
</template>
