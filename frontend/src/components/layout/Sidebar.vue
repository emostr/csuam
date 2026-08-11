<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import NIcon from '../ui/NIcon.vue'
import NAvatar from '../ui/NAvatar.vue'
import { user, isModerator, isHeadTeacher } from '@/lib/auth'
import { ROLE_LABELS } from '@/lib/catalog'

interface NavLink {
  to: string
  label: string
  icon: string
}

interface NavSection {
  title: string
  links: NavLink[]
}

withDefaults(defineProps<{ open?: boolean }>(), { open: false })
defineEmits<{ close: [] }>()

const sections = computed<NavSection[]>(() => {
  const archive: NavSection = {
    title: 'Архив',
    links: [
      { to: '/', label: 'Главная', icon: 'home' },
      { to: '/photos', label: 'Фотографии', icon: 'image' },
      { to: '/videos', label: 'Видео', icon: 'film' },
      { to: '/library', label: 'Библиотека', icon: 'book' },
      { to: '/awards', label: 'Награды', icon: 'award' },
      { to: '/documents', label: 'Документы', icon: 'fileText' },
    ],
  }
  const manage: NavSection = {
    title: 'Управление',
    links: [{ to: '/analytics', label: 'Аналитика', icon: 'barChart' }],
  }
  if (isModerator.value) {
    manage.links.push({ to: '/loans', label: 'Прокат', icon: 'clock' })
  }
  if (isHeadTeacher.value) {
    manage.links.push({ to: '/users', label: 'Пользователи', icon: 'users' })
  }
  manage.links.push({ to: '/settings', label: 'Настройки', icon: 'settings' })
  return [archive, manage]
})
</script>

<template>
  <aside
    class="fixed lg:sticky top-0 left-0 z-40 h-screen w-64 shrink-0 bg-bg border-r border-line flex flex-col transition-transform duration-200 lg:translate-x-0"
    :class="open ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="h-16 flex items-center gap-2.5 px-5 border-b border-line shrink-0">
      <span class="w-8 h-8 bg-accent flex items-center justify-center shrink-0">
        <NIcon name="archive" :size="18" class="text-on-accent" />
      </span>
      <div class="leading-tight">
        <div class="font-extrabold text-ink tracking-normal">Архив<span class="text-accent">ли</span></div>
        <div class="text-[10px] text-faint uppercase font-bold">ЦСУАМ</div>
      </div>
    </div>

    <nav class="flex-1 overflow-y-auto px-3 py-4 space-y-6">
      <div v-for="sec in sections" :key="sec.title">
        <div class="ng-label text-faint px-3 mb-1.5">{{ sec.title }}</div>
        <RouterLink
          v-for="link in sec.links"
          :key="link.to"
          :to="link.to"
          class="group flex items-center gap-3 px-3 py-2.5 text-sm text-muted border-l-2 border-transparent hover:text-ink hover:bg-surface-2 transition-colors"
          exact-active-class="!text-ink !bg-surface-2 !border-accent font-semibold"
          @click="$emit('close')"
        >
          <NIcon :name="link.icon" :size="18" class="shrink-0" />
          <span class="flex-1">{{ link.label }}</span>
        </RouterLink>
      </div>
    </nav>

    <div v-if="user" class="border-t border-line p-3 shrink-0">
      <RouterLink
        to="/settings"
        class="flex items-center gap-3 p-2 hover:bg-surface-2 transition-colors"
        @click="$emit('close')"
      >
        <NAvatar :name="user.full_name" src="/ProfileIcon.svg" :size="38" />
        <div class="min-w-0 flex-1 leading-tight">
          <div class="text-sm font-bold text-ink truncate">{{ user.full_name }}</div>
          <div class="text-[11px] text-accent font-semibold uppercase tracking-normal">
            {{ ROLE_LABELS[user.role] || user.role }}
          </div>
        </div>
        <NIcon name="chevronRight" :size="16" class="text-faint" />
      </RouterLink>
    </div>
  </aside>
</template>
