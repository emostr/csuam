<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue'
import NButton from './NButton.vue'
import NIcon from './NIcon.vue'
import NInput from './NInput.vue'
import { activeDialog } from '@/lib/notify'
import type { Dialog, NotifyIcon } from '@/lib/notify'

const icons: Record<NotifyIcon, { name: string; color: string; tile: string }> = {
  success: { name: 'check', color: 'text-success', tile: 'bg-success/10 border-success/40' },
  error: { name: 'close', color: 'text-danger', tile: 'bg-danger/10 border-danger/40' },
  warning: { name: 'alert', color: 'text-warning', tile: 'bg-warning/10 border-warning/40' },
  info: { name: 'info', color: 'text-info', tile: 'bg-info/10 border-info/40' },
  question: { name: 'help', color: 'text-accent', tile: 'bg-accent/10 border-accent/40' },
}

const modal = computed(() => (activeDialog.value?.kind === 'toast' ? null : activeDialog.value))
const toast = computed(() => (activeDialog.value?.kind === 'toast' ? activeDialog.value : null))
const hasActions = computed(() => modal.value?.kind === 'confirm' || modal.value?.kind === 'prompt')

const value = ref('')
const panel = ref<HTMLElement | null>(null)
const uid = useId()

let timer = 0
let returnFocus: HTMLElement | null = null
let bodyOverflow = ''

function dismiss(d: Dialog) {
  d.close({ confirmed: false, value: '' })
}

function accept(d: Dialog) {
  d.close({ confirmed: true, value: value.value })
}

watch(
  () => activeDialog.value?.id,
  () => {
    window.clearTimeout(timer)
    const d = activeDialog.value
    if (!d) return
    value.value = d.inputValue
    if (d.timer > 0) timer = window.setTimeout(() => dismiss(d), d.timer)
  },
)

watch(
  () => modal.value !== null,
  (open) => {
    if (open) {
      returnFocus = document.activeElement instanceof HTMLElement ? document.activeElement : null
      bodyOverflow = document.body.style.overflow
      document.body.style.overflow = 'hidden'
      return
    }
    document.body.style.overflow = bodyOverflow
    if (returnFocus?.isConnected) returnFocus.focus()
    returnFocus = null
  },
)

function onEnter(el: Element) {
  const d = modal.value
  if (!d) return
  const target =
    d.kind === 'prompt'
      ? el.querySelector<HTMLElement>('input')
      : d.kind === 'confirm'
        ? el.querySelector<HTMLElement>('[data-autofocus]')
        : panel.value
  target?.focus()
  if (target instanceof HTMLInputElement) target.select()
}

function trapFocus(e: KeyboardEvent) {
  const root = panel.value
  if (!root) return
  const items = Array.from(root.querySelectorAll<HTMLElement>('button, input')).filter(
    (el) => !el.hasAttribute('disabled'),
  )
  const first = items.at(0)
  const last = items.at(-1)
  const active = document.activeElement
  if (!first || !last) {
    e.preventDefault()
    root.focus()
  } else if (e.shiftKey && (active === first || !root.contains(active))) {
    e.preventDefault()
    last.focus()
  } else if (!e.shiftKey && (active === last || !root.contains(active))) {
    e.preventDefault()
    first.focus()
  }
}

function onKeydown(e: KeyboardEvent) {
  const d = modal.value
  if (!d) return
  if (e.key === 'Escape') {
    e.preventDefault()
    dismiss(d)
  } else if (e.key === 'Tab') {
    trapFocus(e)
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  window.clearTimeout(timer)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="notify" mode="out-in" @enter="onEnter">
      <div v-if="modal" :key="modal.id" class="fixed inset-0 z-[200] flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-[2px]" @click="dismiss(modal)" />
        <div
          ref="panel"
          class="notify-panel relative w-full max-w-md overflow-hidden bg-surface border border-line border-l-[3px] border-l-accent shadow-2xl outline-none"
          :role="modal.kind === 'confirm' ? 'alertdialog' : 'dialog'"
          aria-modal="true"
          :aria-labelledby="`${uid}-${modal.id}-title`"
          :aria-describedby="modal.text ? `${uid}-${modal.id}-text` : undefined"
          tabindex="-1"
        >
          <button
            class="absolute right-3 top-3 p-1 text-muted hover:text-danger transition-colors cursor-pointer"
            aria-label="Закрыть"
            @click="dismiss(modal)"
          >
            <NIcon name="close" :size="18" />
          </button>

          <div class="flex items-start gap-4 px-6 pt-6 pb-5 pr-12">
            <span
              v-if="modal.icon"
              class="grid size-11 shrink-0 place-items-center border"
              :class="[icons[modal.icon].color, icons[modal.icon].tile]"
            >
              <NIcon :name="icons[modal.icon].name" :size="22" />
            </span>
            <div class="min-w-0 flex-1 pt-0.5">
              <h2 :id="`${uid}-${modal.id}-title`" class="text-lg font-bold text-ink leading-snug">{{ modal.title }}</h2>
              <p v-if="modal.text" :id="`${uid}-${modal.id}-text`" class="mt-1.5 text-sm text-muted break-words">{{ modal.text }}</p>
            </div>
          </div>

          <div v-if="modal.kind === 'prompt'" class="px-6 pb-5">
            <NInput
              v-model="value"
              :label="modal.inputLabel"
              :placeholder="modal.placeholder"
              @keydown.enter.prevent="accept(modal)"
            />
          </div>

          <footer
            v-if="hasActions"
            class="flex items-center justify-end gap-2 px-6 py-4 border-t border-line bg-surface-2/40"
          >
            <NButton variant="secondary" data-autofocus @click="dismiss(modal)">{{ modal.cancelText }}</NButton>
            <NButton :variant="modal.danger ? 'danger' : 'primary'" @click="accept(modal)">
              {{ modal.confirmText }}
            </NButton>
          </footer>

          <div
            v-if="modal.timer > 0"
            class="notify-timer absolute inset-x-0 bottom-0 h-[3px] bg-accent"
            :style="{ animationDuration: `${modal.timer}ms` }"
          />
        </div>
      </div>
    </Transition>

    <Transition name="toast" mode="out-in">
      <div
        v-if="toast"
        :key="toast.id"
        role="status"
        aria-live="polite"
        class="fixed bottom-4 right-4 z-[200] w-[min(22rem,calc(100vw-2rem))] overflow-hidden bg-surface-2 border border-line border-l-[3px] border-l-accent shadow-2xl"
      >
        <div class="flex items-start gap-3 px-4 py-3 pr-10">
          <NIcon
            v-if="toast.icon"
            :name="icons[toast.icon].name"
            :size="18"
            class="shrink-0 mt-0.5"
            :class="icons[toast.icon].color"
          />
          <div class="min-w-0">
            <p class="text-sm font-bold text-ink">{{ toast.title }}</p>
            <p v-if="toast.text" class="mt-0.5 text-sm text-muted break-words">{{ toast.text }}</p>
          </div>
        </div>
        <button
          class="absolute right-2 top-2 p-1 text-muted hover:text-danger transition-colors cursor-pointer"
          aria-label="Закрыть"
          @click="dismiss(toast)"
        >
          <NIcon name="close" :size="16" />
        </button>
        <div
          v-if="toast.timer > 0"
          class="notify-timer absolute inset-x-0 bottom-0 h-[3px] bg-accent"
          :style="{ animationDuration: `${toast.timer}ms` }"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.notify-timer {
  transform-origin: left;
  animation-name: notify-timer;
  animation-timing-function: linear;
  animation-fill-mode: forwards;
}

@keyframes notify-timer {
  from {
    transform: scaleX(1);
  }
  to {
    transform: scaleX(0);
  }
}

.notify-enter-active,
.notify-leave-active {
  transition: opacity 0.2s ease;
}
.notify-enter-active .notify-panel,
.notify-leave-active .notify-panel {
  transition:
    transform 0.24s cubic-bezier(0.16, 1, 0.3, 1),
    opacity 0.2s ease;
}
.notify-enter-from,
.notify-leave-to {
  opacity: 0;
}
.notify-enter-from .notify-panel,
.notify-leave-to .notify-panel {
  transform: translateY(12px) scale(0.98);
  opacity: 0;
}

.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.24s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>
