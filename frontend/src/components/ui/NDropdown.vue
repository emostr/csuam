<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import NIcon from './NIcon.vue'

const props = defineProps({
  align: { type: String, default: 'right' },
  width: { type: Number, default: 200 },
})

const open = ref(false)
const trigger = ref(null)
const menu = ref(null)
const pos = ref({ top: 0, left: 0, origin: 'top' })

function place() {
  const el = trigger.value
  if (!el) return
  const r = el.getBoundingClientRect()
  const vh = window.innerHeight
  const vw = window.innerWidth
  const gap = 6
  const menuH = menu.value?.offsetHeight || 0
  const below = vh - r.bottom
  const flip = below < menuH + gap && r.top > below

  let left = props.align === 'right' ? r.right - props.width : r.left
  left = Math.max(8, Math.min(left, vw - props.width - 8))

  pos.value = {
    top: flip ? r.top - gap - menuH : r.bottom + gap,
    left,
    origin: flip ? 'bottom' : 'top',
  }
}

async function toggle() {
  open.value = !open.value
  if (open.value) {
    await nextTick()
    place()
    await nextTick()
    place()
  }
}

function close() {
  open.value = false
}

function onDocClick(e) {
  if (
    trigger.value?.contains(e.target) ||
    menu.value?.contains(e.target)
  )
    return
  close()
}

function onScroll() {
  if (open.value) close()
}

function onKey(e) {
  if (e.key === 'Escape') close()
}

onMounted(() => {
  document.addEventListener('click', onDocClick)
  window.addEventListener('scroll', onScroll, true)
  window.addEventListener('resize', close)
  document.addEventListener('keydown', onKey)
})
onBeforeUnmount(() => {
  document.removeEventListener('click', onDocClick)
  window.removeEventListener('scroll', onScroll, true)
  window.removeEventListener('resize', close)
  document.removeEventListener('keydown', onKey)
})
</script>

<template>
  <div ref="trigger" class="inline-block">
    <div @click="toggle">
      <slot name="trigger" :open="open">
        <button class="p-2 text-muted hover:text-ink hover:bg-surface-2 transition-colors cursor-pointer">
          <NIcon name="more" :size="18" />
        </button>
      </slot>
    </div>

    <Teleport to="body">
      <Transition name="dd">
        <div
          v-if="open"
          ref="menu"
          class="fixed z-120 bg-surface border border-line shadow-2xl py-1"
          :style="{
            top: pos.top + 'px',
            left: pos.left + 'px',
            width: width + 'px',
            transformOrigin: pos.origin,
          }"
          @click="close"
        >
          <slot />
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.dd-enter-active,
.dd-leave-active {
  transition:
    opacity 0.14s ease,
    transform 0.14s cubic-bezier(0.16, 1, 0.3, 1);
}
.dd-enter-from,
.dd-leave-to {
  opacity: 0;
  transform: scale(0.97) translateY(-4px);
}
</style>
