<script setup lang="ts">
import { ref } from 'vue'
import { CheckIcon, LockIcon, MoonIcon, PaletteIcon, SunIcon, UserIcon } from '@lucide/vue'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Spinner } from '@/components/ui/spinner'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { theme, setTheme, accent, setAccent, ACCENTS } from '@/lib/theme'
import type { Theme } from '@/lib/theme'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user } from '@/lib/auth'
import { ROLE_LABELS, ROLE_TONES, formatDate, initials } from '@/lib/catalog'
import { cn } from '@/lib/utils'
import type { UserRole } from '@/lib/types'

const tab = ref('appearance')

const themes: { id: Theme; label: string; icon: typeof SunIcon }[] = [
  { id: 'light', label: 'Светлая', icon: SunIcon },
  { id: 'dark', label: 'Тёмная', icon: MoonIcon },
]

const pw = ref({ current: '', next: '', repeat: '' })
const pwSaving = ref(false)

function setTab(v: unknown) {
  if (typeof v === 'string') tab.value = v
}

async function changePassword() {
  if (!pw.value.current) {
    notify.warning('Введите текущий пароль')
    return
  }
  if (pw.value.next.length < 6) {
    notify.warning('Новый пароль должен быть не короче 6 символов')
    return
  }
  if (pw.value.next !== pw.value.repeat) {
    notify.warning('Пароли не совпадают')
    return
  }
  pwSaving.value = true
  try {
    await api.post('/auth/password', {
      current_password: pw.value.current,
      new_password: pw.value.next,
    })
    pw.value = { current: '', next: '', repeat: '' }
    notify.success('Пароль изменён')
  } catch (e) {
    notify.error('Не удалось сменить пароль', { text: errorMessage(e) })
  } finally {
    pwSaving.value = false
  }
}

const roleHints: Record<UserRole, string> = {
  head_teacher: 'Полный доступ',
  librarian: 'Частичный доступ',
  teacher: 'Частичный доступ',
}
</script>

<template>
  <div>
    <PageHeader title="Настройки" description="Оформление интерфейса и учётная запись" />

    <Tabs :model-value="tab" @update:model-value="setTab">
      <TabsList class="mb-4">
        <TabsTrigger value="appearance">
          <PaletteIcon />
          Оформление
        </TabsTrigger>
        <TabsTrigger value="account">
          <UserIcon />
          Аккаунт
        </TabsTrigger>
      </TabsList>

      <TabsContent value="appearance" class="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Тема</CardTitle>
            <CardDescription>Светлое или тёмное оформление</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-2 gap-3" role="radiogroup" aria-label="Тема">
              <button
                v-for="t in themes"
                :key="t.id"
                type="button"
                role="radio"
                :aria-checked="theme === t.id"
                :class="
                  cn(
                    'relative flex flex-col items-start gap-3 rounded-lg border p-4 text-left text-sm font-medium transition-colors outline-none hover:bg-muted/50 focus-visible:ring-3 focus-visible:ring-ring/50',
                    theme === t.id && 'border-primary ring-1 ring-primary',
                  )
                "
                @click="setTheme(t.id)"
              >
                <component :is="t.icon" class="size-5 text-muted-foreground" />
                {{ t.label }}
                <CheckIcon v-if="theme === t.id" class="absolute top-3 right-3 size-4 text-primary" />
              </button>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Акцентный цвет</CardTitle>
            <CardDescription>Цвет кнопок, ссылок и выделений</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-3 gap-3" role="radiogroup" aria-label="Акцентный цвет">
              <button
                v-for="a in ACCENTS"
                :key="a.id"
                type="button"
                role="radio"
                :aria-checked="accent === a.id"
                :class="
                  cn(
                    'flex items-center gap-2.5 rounded-lg border p-3 text-sm font-medium transition-colors outline-none hover:bg-muted/50 focus-visible:ring-3 focus-visible:ring-ring/50',
                    accent === a.id && 'border-primary ring-1 ring-primary',
                  )
                "
                @click="setAccent(a.id)"
              >
                <span
                  class="flex size-5 shrink-0 items-center justify-center rounded-full"
                  :style="{ background: a.hex }"
                >
                  <CheckIcon v-if="accent === a.id" class="size-3 text-white" :stroke-width="3" />
                </span>
                {{ a.label }}
              </button>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="account" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card>
          <CardContent class="flex flex-col items-center gap-4 py-2 text-center">
            <Avatar class="size-20">
              <AvatarImage src="/ProfileIcon.svg" :alt="user?.full_name" />
              <AvatarFallback class="text-xl">{{ initials(user?.full_name) }}</AvatarFallback>
            </Avatar>
            <div class="space-y-1">
              <div class="font-semibold">{{ user?.full_name }}</div>
              <div class="text-muted-foreground">@{{ user?.username }}</div>
            </div>
            <StatusBadge v-if="user" :tone="ROLE_TONES[user.role]">{{ ROLE_LABELS[user.role] }}</StatusBadge>
          </CardContent>
        </Card>

        <Card class="lg:col-span-2">
          <CardHeader>
            <CardTitle>Учётная запись</CardTitle>
            <CardDescription>{{ user ? roleHints[user.role] : '' }}</CardDescription>
          </CardHeader>
          <CardContent>
            <dl class="divide-y">
              <div class="flex justify-between gap-4 py-2.5 first:pt-0">
                <dt class="text-muted-foreground">Логин</dt>
                <dd class="font-medium">{{ user?.username }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5">
                <dt class="text-muted-foreground">Роль</dt>
                <dd class="font-medium">{{ user ? ROLE_LABELS[user.role] : '' }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-2.5 last:pb-0">
                <dt class="text-muted-foreground">В системе с</dt>
                <dd>{{ formatDate(user?.created_at) }}</dd>
              </div>
            </dl>
          </CardContent>
        </Card>

        <Card class="lg:col-span-3">
          <CardHeader>
            <CardTitle>Смена пароля</CardTitle>
            <CardDescription>Для подтверждения понадобится текущий пароль</CardDescription>
          </CardHeader>
          <form @submit.prevent="changePassword">
            <CardContent class="grid grid-cols-1 gap-6 sm:grid-cols-3">
              <Field>
                <FieldLabel for="pw-current">Текущий пароль</FieldLabel>
                <Input
                  id="pw-current"
                  v-model="pw.current"
                  type="password"
                  placeholder="••••••••"
                  autocomplete="current-password"
                />
              </Field>
              <Field>
                <FieldLabel for="pw-next">Новый пароль</FieldLabel>
                <Input id="pw-next" v-model="pw.next" type="password" placeholder="••••••••" autocomplete="new-password" />
                <FieldDescription>Минимум 6 символов</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="pw-repeat">Повторите новый пароль</FieldLabel>
                <Input
                  id="pw-repeat"
                  v-model="pw.repeat"
                  type="password"
                  placeholder="••••••••"
                  autocomplete="new-password"
                />
              </Field>
            </CardContent>
            <CardFooter class="mt-6 border-t pt-6">
              <Button type="submit" :disabled="pwSaving">
                <Spinner v-if="pwSaving" />
                <LockIcon v-else />
                Обновить пароль
              </Button>
            </CardFooter>
          </form>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
