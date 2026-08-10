<script setup>
import { ref, onMounted } from 'vue'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NButton from '@/components/ui/NButton.vue'
import NTable from '@/components/ui/NTable.vue'
import NBadge from '@/components/ui/NBadge.vue'
import NModal from '@/components/ui/NModal.vue'
import NInput from '@/components/ui/NInput.vue'
import NSelect from '@/components/ui/NSelect.vue'
import NAvatar from '@/components/ui/NAvatar.vue'
import NAlert from '@/components/ui/NAlert.vue'
import NIcon from '@/components/ui/NIcon.vue'
import { api } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user, isHeadTeacher } from '@/lib/auth'
import { ROLE_LABELS, formatDate } from '@/lib/catalog'

const users = ref([])
const loading = ref(true)

const createOpen = ref(false)
const newUser = ref({ full_name: '', username: '', password: '', role: 'teacher' })
const saving = ref(false)

const columns = [
  { key: 'full_name', label: 'Пользователь' },
  { key: 'username', label: 'Логин' },
  { key: 'role', label: 'Роль' },
  { key: 'created_at', label: 'В системе с' },
  { key: 'actions', label: '', align: 'right' },
]

const roleOptions = [
  { value: 'teacher', label: 'Учитель' },
  { value: 'librarian', label: 'Библиотекарь' },
  { value: 'head_teacher', label: 'Завуч' },
]

const roleVariants = {
  head_teacher: 'accent',
  librarian: 'info',
  teacher: 'neutral',
}

async function load() {
  loading.value = true
  try {
    users.value = await api.get('/users')
  } catch (e) {
    notify.error('Не удалось загрузить пользователей', { text: e.message })
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
    notify.error('Не удалось создать пользователя', { text: e.message })
  } finally {
    saving.value = false
  }
}

async function resetPassword(row) {
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
    notify.error('Не удалось сменить пароль', { text: e.message })
  }
}

async function removeUser(row) {
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
    notify.error('Не удалось удалить', { text: e.message })
  }
}
</script>

<template>
  <div>
    <NPageHeader title="Пользователи">
      <template #actions>
        <NButton v-if="isHeadTeacher" icon="userPlus" @click="openCreate">Создать пользователя</NButton>
      </template>
    </NPageHeader>

    <NAlert v-if="!isHeadTeacher" variant="warning" title="Недостаточно прав">
      Управление пользователями доступно только завучу
    </NAlert>

    <template v-else>
      <div v-if="loading" class="text-center text-muted py-16">Загрузка…</div>
      <NTable v-else :columns="columns" :rows="users">
        <template #full_name="{ row }">
          <div class="flex items-center gap-3">
            <NAvatar :name="row.full_name" src="/ProfileIcon.svg" :size="34" />
            <div>
              <div class="font-semibold text-ink">{{ row.full_name }}</div>
              <div v-if="row.id === user?.id" class="text-[11px] text-accent font-semibold">Это Вы</div>
            </div>
          </div>
        </template>
        <template #username="{ value }">
          <span class="text-muted">@{{ value }}</span>
        </template>
        <template #role="{ value }">
          <NBadge :variant="roleVariants[value]" dot>{{ ROLE_LABELS[value] || value }}</NBadge>
        </template>
        <template #created_at="{ value }">
          <span class="text-muted">{{ formatDate(value) }}</span>
        </template>
        <template #actions="{ row }">
          <div class="flex items-center justify-end gap-1">
            <button
              class="h-8 w-8 flex items-center justify-center text-muted hover:text-ink hover:bg-surface-2 transition-colors cursor-pointer"
              title="Сменить пароль"
              @click="resetPassword(row)"
            >
              <NIcon name="lock" :size="16" />
            </button>
            <button
              v-if="row.id !== user?.id"
              class="h-8 w-8 flex items-center justify-center text-muted hover:text-danger hover:bg-surface-2 transition-colors cursor-pointer"
              title="Удалить"
              @click="removeUser(row)"
            >
              <NIcon name="trash" :size="16" />
            </button>
          </div>
        </template>
      </NTable>

      <NModal v-model="createOpen" title="Новый пользователь" subtitle="Учётная запись сотрудника школы">
        <div class="space-y-4">
          <NInput v-model="newUser.full_name" label="ФИО" placeholder="Иванова Мария Петровна" icon="user" />
          <div class="grid sm:grid-cols-2 gap-4">
            <NInput v-model="newUser.username" label="Логин" placeholder="m.ivanova" />
            <NInput v-model="newUser.password" label="Пароль" type="password" placeholder="Минимум 6 символов" icon="lock" />
          </div>
          <NSelect v-model="newUser.role" label="Роль" :options="roleOptions" />
        </div>
        <template #footer="{ close }">
          <NButton variant="ghost" @click="close">Отмена</NButton>
          <NButton :loading="saving" icon="check" @click="createUser">Создать</NButton>
        </template>
      </NModal>
    </template>
  </div>
</template>
