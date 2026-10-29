<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { EllipsisIcon, KeyRoundIcon, ShieldAlertIcon, Trash2Icon, UserPlusIcon } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import PageHeader from '@/components/PageHeader.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { api, errorMessage } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user, isHeadTeacher } from '@/lib/auth'
import { ROLE_LABELS, ROLE_TONES, formatDate, initials } from '@/lib/catalog'
import type { User, UserRole } from '@/lib/types'

interface NewUser {
  full_name: string
  username: string
  password: string
  role: UserRole
}

const users = ref<User[]>([])
const loading = ref(true)

const createOpen = ref(false)
const newUser = ref<NewUser>({ full_name: '', username: '', password: '', role: 'teacher' })
const saving = ref(false)

const roles: UserRole[] = ['teacher', 'librarian', 'head_teacher']

function setRole(v: unknown) {
  const r = roles.find((x) => x === v)
  if (r) newUser.value.role = r
}

async function load() {
  loading.value = true
  try {
    users.value = await api.get<User[]>('/users')
  } catch (e) {
    notify.error('Не удалось загрузить пользователей', { text: errorMessage(e) })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (isHeadTeacher.value) load()
})

function openCreate() {
  newUser.value = { full_name: '', username: '', password: '', role: 'teacher' }
  createOpen.value = true
}

async function createUser() {
  const u = newUser.value
  if (!u.full_name.trim() || !u.username.trim()) {
    notify.warning('Укажите ФИО и логин')
    return
  }
  if (u.password.length < 6) {
    notify.warning('Пароль должен быть не короче 6 символов')
    return
  }
  saving.value = true
  try {
    await api.post('/users', u)
    createOpen.value = false
    notify.success('Пользователь создан', {
      text: `${u.full_name} может входить с логином «${u.username}»`,
    })
    load()
  } catch (e) {
    notify.error('Не удалось создать пользователя', { text: errorMessage(e) })
  } finally {
    saving.value = false
  }
}

async function resetPassword(row: User) {
  const password = await notify.prompt({
    title: 'Новый пароль',
    inputLabel: `Для пользователя ${row.full_name}`,
    placeholder: 'Минимум 6 символов',
  })
  if (password == null) return
  if (password.length < 6) {
    notify.warning('Пароль должен быть не короче 6 символов')
    return
  }
  try {
    await api.post(`/users/${row.id}/password`, { password })
    notify.success('Пароль изменён')
  } catch (e) {
    notify.error('Не удалось сменить пароль', { text: errorMessage(e) })
  }
}

async function removeUser(row: User) {
  const ok = await notify.confirm({
    title: 'Удалить пользователя?',
    text: `${row.full_name} (@${row.username}) потеряет доступ к платформе. Добавленные им материалы останутся в архиве.`,
    confirmText: 'Удалить',
    danger: true,
  })
  if (!ok) return
  try {
    await api.del(`/users/${row.id}`)
    notify.success('Пользователь удалён')
    load()
  } catch (e) {
    notify.error('Не удалось удалить', { text: errorMessage(e) })
  }
}
</script>

<template>
  <div>
    <PageHeader title="Пользователи" description="Учётные записи сотрудников школы">
      <template #actions>
        <Button v-if="isHeadTeacher" @click="openCreate">
          <UserPlusIcon />
          Создать пользователя
        </Button>
      </template>
    </PageHeader>

    <Alert v-if="!isHeadTeacher">
      <ShieldAlertIcon />
      <AlertTitle>Недостаточно прав</AlertTitle>
      <AlertDescription>Управление пользователями доступно только завучу</AlertDescription>
    </Alert>

    <template v-else>
      <Card class="py-0">
        <Table>
          <TableHeader>
            <TableRow class="hover:bg-transparent">
              <TableHead class="pl-4">Пользователь</TableHead>
              <TableHead>Логин</TableHead>
              <TableHead>Роль</TableHead>
              <TableHead>В системе с</TableHead>
              <TableHead class="w-12 pr-4"><span class="sr-only">Действия</span></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="loading">
              <TableRow v-for="i in 4" :key="i">
                <TableCell class="pl-4">
                  <div class="flex items-center gap-3">
                    <Skeleton class="size-8 rounded-full" />
                    <Skeleton class="h-4 w-40" />
                  </div>
                </TableCell>
                <TableCell><Skeleton class="h-4 w-24" /></TableCell>
                <TableCell><Skeleton class="h-5 w-20 rounded-full" /></TableCell>
                <TableCell><Skeleton class="h-4 w-28" /></TableCell>
                <TableCell class="pr-4" />
              </TableRow>
            </template>
            <TableRow v-for="row in users" v-else :key="row.id">
              <TableCell class="pl-4">
                <div class="flex items-center gap-3">
                  <Avatar class="size-8">
                    <AvatarImage src="/ProfileIcon.svg" :alt="row.full_name" />
                    <AvatarFallback>{{ initials(row.full_name) }}</AvatarFallback>
                  </Avatar>
                  <span class="font-medium">{{ row.full_name }}</span>
                  <Badge v-if="row.id === user?.id" variant="secondary">Это вы</Badge>
                </div>
              </TableCell>
              <TableCell class="text-muted-foreground">@{{ row.username }}</TableCell>
              <TableCell>
                <StatusBadge :tone="ROLE_TONES[row.role]">{{ ROLE_LABELS[row.role] }}</StatusBadge>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ formatDate(row.created_at) }}</TableCell>
              <TableCell class="pr-4 text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon-sm">
                      <EllipsisIcon />
                      <span class="sr-only">Действия</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" class="w-48">
                    <DropdownMenuItem @select="resetPassword(row)">
                      <KeyRoundIcon />
                      Сменить пароль
                    </DropdownMenuItem>
                    <template v-if="row.id !== user?.id">
                      <DropdownMenuSeparator />
                      <DropdownMenuItem variant="destructive" @select="removeUser(row)">
                        <Trash2Icon />
                        Удалить
                      </DropdownMenuItem>
                    </template>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
            <TableRow v-if="!loading && !users.length">
              <TableCell colspan="5" class="h-24 text-center text-muted-foreground">Нет данных</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Card>

      <Dialog v-model:open="createOpen">
        <DialogContent>
          <form class="grid gap-6" @submit.prevent="createUser">
            <DialogHeader>
              <DialogTitle>Новый пользователь</DialogTitle>
              <DialogDescription>Учётная запись сотрудника школы</DialogDescription>
            </DialogHeader>
            <FieldGroup>
              <Field>
                <FieldLabel for="new-user-name">ФИО</FieldLabel>
                <Input id="new-user-name" v-model="newUser.full_name" placeholder="Иванова Мария Петровна" />
              </Field>
              <div class="grid gap-6 sm:grid-cols-2">
                <Field>
                  <FieldLabel for="new-user-login">Логин</FieldLabel>
                  <Input id="new-user-login" v-model="newUser.username" placeholder="m.ivanova" autocomplete="off" />
                </Field>
                <Field>
                  <FieldLabel for="new-user-password">Пароль</FieldLabel>
                  <Input
                    id="new-user-password"
                    v-model="newUser.password"
                    type="password"
                    placeholder="Минимум 6 символов"
                    autocomplete="new-password"
                  />
                </Field>
              </div>
              <Field>
                <FieldLabel for="new-user-role">Роль</FieldLabel>
                <Select :model-value="newUser.role" @update:model-value="setRole">
                  <SelectTrigger id="new-user-role" class="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem v-for="r in roles" :key="r" :value="r">{{ ROLE_LABELS[r] }}</SelectItem>
                  </SelectContent>
                </Select>
              </Field>
            </FieldGroup>
            <DialogFooter>
              <DialogClose as-child>
                <Button type="button" variant="outline">Отмена</Button>
              </DialogClose>
              <Button type="submit" :disabled="saving">
                <Spinner v-if="saving" />
                Создать
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </template>
  </div>
</template>
