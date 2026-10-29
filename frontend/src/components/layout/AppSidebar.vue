<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import {
  ArchiveIcon,
  AwardIcon,
  BookOpenIcon,
  ChartColumnIcon,
  ChevronsUpDownIcon,
  ClockIcon,
  FileTextIcon,
  FilmIcon,
  HouseIcon,
  ImageIcon,
  LogOutIcon,
  MoonIcon,
  SettingsIcon,
  SunIcon,
  UsersIcon,
} from '@lucide/vue'
import type { LucideIcon } from '@lucide/vue'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  useSidebar,
} from '@/components/ui/sidebar'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { user, logout, isModerator, isHeadTeacher } from '@/lib/auth'
import { ROLE_LABELS, initials } from '@/lib/catalog'
import { notify } from '@/lib/notify'
import { theme, toggleTheme } from '@/lib/theme'

interface NavLink {
  to: string
  label: string
  icon: LucideIcon
}

interface NavSection {
  title: string
  links: NavLink[]
}

const route = useRoute()
const router = useRouter()
const { isMobile, setOpenMobile } = useSidebar()

const sections = computed<NavSection[]>(() => {
  const archive: NavSection = {
    title: 'Архив',
    links: [
      { to: '/', label: 'Главная', icon: HouseIcon },
      { to: '/photos', label: 'Фотографии', icon: ImageIcon },
      { to: '/videos', label: 'Видео', icon: FilmIcon },
      { to: '/library', label: 'Библиотека', icon: BookOpenIcon },
      { to: '/awards', label: 'Награды', icon: AwardIcon },
      { to: '/documents', label: 'Документы', icon: FileTextIcon },
    ],
  }
  const manage: NavSection = {
    title: 'Управление',
    links: [{ to: '/analytics', label: 'Аналитика', icon: ChartColumnIcon }],
  }
  if (isModerator.value) {
    manage.links.push({ to: '/loans', label: 'Прокат', icon: ClockIcon })
  }
  if (isHeadTeacher.value) {
    manage.links.push({ to: '/users', label: 'Пользователи', icon: UsersIcon })
  }
  manage.links.push({ to: '/settings', label: 'Настройки', icon: SettingsIcon })
  return [archive, manage]
})

function closeMobile() {
  if (isMobile.value) setOpenMobile(false)
}

function openSettings() {
  closeMobile()
  router.push('/settings')
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
  <Sidebar collapsible="icon">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <RouterLink to="/" @click="closeMobile">
              <span
                class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
              >
                <ArchiveIcon class="size-4" />
              </span>
              <span class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">Архивли</span>
                <span class="truncate text-xs text-muted-foreground">ЦСУАМ</span>
              </span>
            </RouterLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup v-for="sec in sections" :key="sec.title">
        <SidebarGroupLabel>{{ sec.title }}</SidebarGroupLabel>
        <SidebarMenu>
          <SidebarMenuItem v-for="link in sec.links" :key="link.to">
            <SidebarMenuButton as-child :is-active="route.path === link.to" :tooltip="link.label">
              <RouterLink :to="link.to" @click="closeMobile">
                <component :is="link.icon" />
                <span>{{ link.label }}</span>
              </RouterLink>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarGroup>
    </SidebarContent>

    <SidebarFooter v-if="user">
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <SidebarMenuButton
                size="lg"
                class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              >
                <Avatar class="size-8 rounded-lg">
                  <AvatarImage src="/ProfileIcon.svg" :alt="user.full_name" />
                  <AvatarFallback class="rounded-lg">{{ initials(user.full_name) }}</AvatarFallback>
                </Avatar>
                <span class="grid flex-1 text-left text-sm leading-tight">
                  <span class="truncate font-medium">{{ user.full_name }}</span>
                  <span class="truncate text-xs text-muted-foreground">{{ ROLE_LABELS[user.role] }}</span>
                </span>
                <ChevronsUpDownIcon class="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              class="w-(--reka-dropdown-menu-trigger-width) min-w-56"
              :side="isMobile ? 'bottom' : 'right'"
              align="end"
              :side-offset="4"
            >
              <DropdownMenuLabel class="p-0 font-normal">
                <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                  <Avatar class="size-8 rounded-lg">
                    <AvatarImage src="/ProfileIcon.svg" :alt="user.full_name" />
                    <AvatarFallback class="rounded-lg">{{ initials(user.full_name) }}</AvatarFallback>
                  </Avatar>
                  <div class="grid flex-1 leading-tight">
                    <span class="truncate font-medium">{{ user.full_name }}</span>
                    <span class="truncate text-xs text-muted-foreground">@{{ user.username }}</span>
                  </div>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuGroup>
                <DropdownMenuItem @select="openSettings">
                  <SettingsIcon />
                  Настройки
                </DropdownMenuItem>
                <DropdownMenuItem @select="toggleTheme">
                  <SunIcon v-if="theme === 'dark'" />
                  <MoonIcon v-else />
                  {{ theme === 'dark' ? 'Светлая тема' : 'Тёмная тема' }}
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuItem variant="destructive" @select="doLogout">
                <LogOutIcon />
                Выйти
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>

    <SidebarRail />
  </Sidebar>
</template>
