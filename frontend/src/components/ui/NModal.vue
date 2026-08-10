<script setup>
import { watch } from 'vue'
import NIcon from './NIcon.vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  size: { type: String, default: 'md' },
})
const emit = defineEmits(['update:modelValue'])

const sizes = { sm: 'max-w-md', md: 'max-w-lg', lg: 'max-w-2xl', xl: 'max-w-4xl' }

function close() {
  emit('update:modelValue', false)
}

watch(
  () => props.modelValue,
  (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
  },
)
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[100] flex items-start justify-center p-4 sm:p-8 overflow-y-auto"
      >
        <div class="fixed inset-0 bg-black/60 backdrop-blur-[2px]" @click="close" />
        <div
          class="relative w-full bg-surface border border-line border-l-[3px] border-l-accent shadow-2xl mt-8 sm:mt-16"
          :class="sizes[size] || sizes.md"
          role="dialog"
          aria-modal="true"
        >
          <header class="flex items-start justify-between gap-4 px-6 py-4 border-b border-line">
            <div class="min-w-0">
              <slot name="header">
                <h3 class="text-lg font-bold text-ink truncate">{{ title }}</h3>
                <p v-if="subtitle" class="text-sm text-muted mt-0.5">{{ subtitle }}</p>
              </slot>
            </div>
            <button
              class="shrink-0 -mr-1 p-1 text-muted hover:text-danger transition-colors cursor-pointer"
              @click="close"
            >
              <NIcon name="close" :size="20" />
            </button>
          </header>

          <div class="px-6 py-5">
            <slot />
          </div>

          <footer v-if="$slots.footer" class="flex items-center justify-end gap-2 px-6 py-4 border-t border-line bg-surface-2/40">
            <slot name="footer" :close="close" />
          </footer>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-active .relative,
.modal-leave-active .relative {
  transition:
    transform 0.24s cubic-bezier(0.16, 1, 0.3, 1),
    opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: translateY(16px) scale(0.98);
  opacity: 0;
}
</style>
