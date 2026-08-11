<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { theme, toggleTheme } from '@/lib/theme'
import { notify } from '@/lib/notify'
import { api } from '@/lib/api'
import { user, logout, isModerator } from '@/lib/auth'
import { ROLE_LABELS } from '@/lib/catalog'
import type { Notifications } from '@/lib/types'
import NIcon from '../ui/NIcon.vue'
import NAvatar from '../ui/NAvatar.vue'
import NDropdown from '../ui/NDropdown.vue'
import NDropdownItem from '../ui/NDropdownItem.vue'

defineEmits<{ 'toggle-sidebar': [] }>()

const router = useRouter()
const search = ref('')
const counters = ref<Notifications>({ overdue_loans: 0, pending_materials: 0 })

onMounted(async () => {
  if (!isModerator.value) return
  try {
    counters.value = await api.get<Notifications>('/notifications')
  } catch {}
})

function submitSearch() {
  const q = search.value.trim()
  router.push({ path: '/', query: q ? { q } : {} })
}

async function doLogout() {
  const ok = await notify.confirm({
    title: 'Выйти из аккаунта?',
    text: 'Текущая сессия будет завершена.',
    confirmText: 'Выйти',
    danger: true,
  })
  if (!ok) return
  await logout()
  router.push('/login')
}
</script>

<template>
  <header class="sticky top-0 z-30 h-16 bg-bg/85 backdrop-blur border-b border-line flex items-center gap-3 px-4 sm:px-6">
    <button
      class="lg:hidden h-9 w-9 flex items-center justify-center text-muted hover:text-ink cursor-pointer"
      @click="$emit('toggle-sidebar')"
    >
      <NIcon name="menu" :size="22" />
    </button>

    <form class="relative hidden sm:block w-full max-w-xs" @submit.prevent="submitSearch">
      <NIcon name="search" :size="17" class="absolute left-3 top-1/2 -translate-y-1/2 text-faint pointer-events-none" />
      <input
        v-model="search"
        type="text"
        placeholder="Поиск по архиву…"
        class="w-full h-9 bg-surface-2 border border-line focus:border-accent outline-none pl-9 pr-3 text-sm text-ink placeholder:text-faint transition-colors"
      />
    </form>

    <div class="flex-1" />

    <button
      class="h-9 w-9 flex items-center justify-center text-muted hover:text-ink hover:bg-surface-2 transition-colors cursor-pointer"
      :title="theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'"
      @click="toggleTheme"
    >
      <NIcon :name="theme === 'dark' ? 'sun' : 'moon'" :size="19" />
    </button>

    <NDropdown v-if="isModerator" align="right" :width="320">
      <template #trigger>
        <button class="relative h-9 w-9 flex items-center justify-center text-muted hover:text-ink hover:bg-surface-2 transition-colors cursor-pointer">
          <NIcon name="bell" :size="19" />
          <span
            v-if="counters.overdue_loans > 0 || counters.pending_materials > 0"
            class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-danger ring-2 ring-bg"
          />
        </button>
      </template>
      <div class="px-3.5 py-2 ng-label text-faint border-b border-line">Уведомления</div>
      <div class="max-h-72 overflow-y-auto">
        <div
          v-if="counters.overdue_loans > 0"
          class="px-3.5 py-3 hover:bg-surface-2 border-b border-line/60 cursor-pointer"
          @click="router.push('/loans')"
        >
          <div class="text-sm text-danger font-semibold">Просрочен возврат книг</div>
          <div class="text-xs text-muted mt-0.5">Не возвращено вовремя: {{ counters.overdue_loans }}</div>
        </div>
        <div
          v-if="counters.pending_materials > 0"
          class="px-3.5 py-3 hover:bg-surface-2 border-b border-line/60 cursor-pointer"
          @click="router.push({ path: '/', query: { status: 'pending' } })"
        >
          <div class="text-sm text-ink font-semibold">Материалы на модерации</div>
          <div class="text-xs text-muted mt-0.5">Ожидают проверки: {{ counters.pending_materials }}</div>
        </div>
        <div
          v-if="!counters.overdue_loans && !counters.pending_materials"
          class="px-3.5 py-6 text-center text-sm text-muted"
        >
          Новых уведомлений нет
        </div>
      </div>
    </NDropdown>

    <div class="w-px h-6 bg-line mx-1 hidden sm:block" />

    <NDropdown align="right" :width="240">
      <template #trigger>
        <button class="flex items-center gap-2 pl-1 pr-2 h-9 hover:bg-surface-2 transition-colors cursor-pointer">
          <NAvatar :name="user?.full_name || ''" src="/ProfileIcon.svg" :size="30" />
          <NIcon name="chevronDown" :size="15" class="text-faint hidden sm:block" />
        </button>
      </template>
      <div class="px-3.5 py-2.5 border-b border-line">
        <div class="text-sm font-bold text-ink">{{ user?.full_name }}</div>
        <div class="text-xs text-muted">{{ user ? ROLE_LABELS[user.role] : '' }}</div>
      </div>
      <NDropdownItem icon="settings" @click="router.push('/settings')">Настройки</NDropdownItem>
      <div class="my-1 border-t border-line" />
      <NDropdownItem icon="logout" danger @click="doLogout">Выйти</NDropdownItem>
    </NDropdown>
  </header>
</template>
