<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { BellIcon, ClockIcon, FileClockIcon, MoonIcon, SearchIcon, SunIcon } from '@lucide/vue'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group'
import { Separator } from '@/components/ui/separator'
import { SidebarTrigger } from '@/components/ui/sidebar'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { api } from '@/lib/api'
import { isModerator } from '@/lib/auth'
import { theme, toggleTheme } from '@/lib/theme'
import type { Notifications } from '@/lib/types'

const route = useRoute()
const router = useRouter()
const search = ref('')
const counters = ref<Notifications>({ overdue_loans: 0, pending_materials: 0 })

const title = computed(() => route.meta.title ?? '')
const hasNotifications = computed(
  () => counters.value.overdue_loans > 0 || counters.value.pending_materials > 0,
)

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
</script>

<template>
  <header
    class="sticky top-0 z-20 flex h-14 shrink-0 items-center gap-2 border-b bg-background/80 px-4 backdrop-blur supports-backdrop-filter:bg-background/60 md:px-6"
  >
    <SidebarTrigger class="-ml-1.5" />
    <Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4 data-[orientation=vertical]:self-center" />

    <Breadcrumb class="min-w-0">
      <BreadcrumbList class="flex-nowrap">
        <template v-if="route.path !== '/'">
          <BreadcrumbItem class="hidden md:inline-flex">
            <BreadcrumbLink as-child>
              <RouterLink to="/">Архивли</RouterLink>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator class="hidden md:block" />
        </template>
        <BreadcrumbItem class="min-w-0">
          <BreadcrumbPage class="truncate">{{ title }}</BreadcrumbPage>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>

    <div class="ml-auto flex items-center gap-1">
      <form class="mr-1 hidden w-64 lg:block" role="search" @submit.prevent="submitSearch">
        <InputGroup>
          <InputGroupAddon>
            <SearchIcon />
          </InputGroupAddon>
          <InputGroupInput v-model="search" placeholder="Поиск по архиву…" aria-label="Поиск по архиву" />
        </InputGroup>
      </form>

      <Tooltip>
        <TooltipTrigger as-child>
          <Button variant="ghost" size="icon-sm" @click="toggleTheme">
            <SunIcon v-if="theme === 'dark'" />
            <MoonIcon v-else />
            <span class="sr-only">Сменить тему</span>
          </Button>
        </TooltipTrigger>
        <TooltipContent>{{ theme === 'dark' ? 'Светлая тема' : 'Тёмная тема' }}</TooltipContent>
      </Tooltip>

      <DropdownMenu v-if="isModerator">
        <DropdownMenuTrigger as-child>
          <Button variant="ghost" size="icon-sm" class="relative">
            <BellIcon />
            <span
              v-if="hasNotifications"
              class="absolute top-1.5 right-1.5 size-2 rounded-full bg-destructive ring-2 ring-background"
            />
            <span class="sr-only">Уведомления</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-72">
          <DropdownMenuLabel>Уведомления</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem v-if="counters.overdue_loans > 0" class="items-start" @select="router.push('/loans')">
            <ClockIcon class="mt-0.5 text-destructive" />
            <div class="grid gap-0.5">
              <span class="font-medium">Просрочен возврат книг</span>
              <span class="text-xs text-muted-foreground">Не возвращено вовремя: {{ counters.overdue_loans }}</span>
            </div>
          </DropdownMenuItem>
          <DropdownMenuItem
            v-if="counters.pending_materials > 0"
            class="items-start"
            @select="router.push({ path: '/', query: { status: 'pending' } })"
          >
            <FileClockIcon class="mt-0.5 text-warning" />
            <div class="grid gap-0.5">
              <span class="font-medium">Материалы на модерации</span>
              <span class="text-xs text-muted-foreground">Ожидают проверки: {{ counters.pending_materials }}</span>
            </div>
          </DropdownMenuItem>
          <div v-if="!hasNotifications" class="px-2 py-6 text-center text-sm text-muted-foreground">
            Новых уведомлений нет
          </div>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </header>
</template>
