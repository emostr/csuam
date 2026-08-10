<script setup>
import { ref } from 'vue'
import { theme, toggleTheme, accent, setAccent, ACCENTS } from '@/lib/theme'
import NPageHeader from '@/components/ui/NPageHeader.vue'
import NCard from '@/components/ui/NCard.vue'
import NTabs from '@/components/ui/NTabs.vue'
import NIcon from '@/components/ui/NIcon.vue'
import NAvatar from '@/components/ui/NAvatar.vue'
import NBadge from '@/components/ui/NBadge.vue'
import NInput from '@/components/ui/NInput.vue'
import NButton from '@/components/ui/NButton.vue'
import { api } from '@/lib/api'
import { notify } from '@/lib/notify'
import { user } from '@/lib/auth'
import { ROLE_LABELS, formatDate } from '@/lib/catalog'

const tab = ref('appearance')
const tabs = [
  { value: 'appearance', label: 'Оформление', icon: 'palette' },
  { value: 'account', label: 'Аккаунт', icon: 'user' },
]

const pw = ref({ current: '', next: '', repeat: '' })
const pwSaving = ref(false)

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
    notify.error('Не удалось сменить пароль', { text: e.message })
  } finally {
    pwSaving.value = false
  }
}

const roleHints = {
  head_teacher: 'Полный доступ',
  librarian: 'Частичный доступ',
  teacher: 'Частичный доступ',
}
</script>

<template>
  <div>
    <NPageHeader title="Настройки" />

    <NTabs v-model="tab" :tabs="tabs" class="mb-6" />

    <div v-if="tab === 'appearance'" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <NCard title="Тема">
        <div class="grid grid-cols-2 gap-3">
          <button
            class="border p-4 text-left transition-colors ng-tile-press cursor-pointer"
            :class="theme === 'dark' ? 'border-accent' : 'border-line hover:border-line-strong'"
            @click="theme !== 'dark' && toggleTheme()"
          >
            <div class="flex items-center justify-between mb-3">
              <NIcon name="moon" :size="20" class="text-accent" />
              <NIcon v-if="theme === 'dark'" name="checkCircle" :size="18" class="text-accent" />
            </div>
            <div class="font-bold text-ink text-sm">Тёмная</div>
          </button>
          <button
            class="border p-4 text-left transition-colors ng-tile-press cursor-pointer"
            :class="theme === 'light' ? 'border-accent' : 'border-line hover:border-line-strong'"
            @click="theme !== 'light' && toggleTheme()"
          >
            <div class="flex items-center justify-between mb-3">
              <NIcon name="sun" :size="20" class="text-accent" />
              <NIcon v-if="theme === 'light'" name="checkCircle" :size="18" class="text-accent" />
            </div>
            <div class="font-bold text-ink text-sm">Светлая</div>
          </button>
        </div>
      </NCard>

      <NCard title="Акцентный цвет">
        <div class="grid grid-cols-3 gap-3">
          <button
            v-for="a in ACCENTS"
            :key="a.id"
            class="relative h-16 border-2 transition-all ng-tile-press cursor-pointer flex items-end p-2"
            :class="accent === a.id ? 'border-ink' : 'border-transparent hover:border-line-strong'"
            :style="{ background: a.hex }"
            @click="setAccent(a.id)"
          >
            <span class="text-[11px] font-bold text-white mix-blend-difference">{{ a.label }}</span>
            <NIcon v-if="accent === a.id" name="check" :size="18" :stroke="3" class="absolute top-2 right-2 text-white mix-blend-difference" />
          </button>
        </div>
      </NCard>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <NCard title="Профиль" class="lg:col-span-1">
        <div class="flex flex-col items-center text-center gap-4 py-2">
          <NAvatar :name="user?.full_name || ''" src="/ProfileIcon.svg" :size="88" />
          <div>
            <div class="font-bold text-ink">{{ user?.full_name }}</div>
            <div class="text-xs text-muted mt-0.5">@{{ user?.username }}</div>
            <NBadge variant="accent" class="mt-2">{{ ROLE_LABELS[user?.role] || user?.role }}</NBadge>
          </div>
        </div>
      </NCard>

      <NCard title="Учётная запись" class="lg:col-span-2" accent>
        <dl class="divide-y divide-line text-sm">
          <div class="flex justify-between gap-4 py-2.5">
            <dt class="text-muted">Логин</dt>
            <dd class="text-ink font-semibold">{{ user?.username }}</dd>
          </div>
          <div class="flex justify-between gap-4 py-2.5">
            <dt class="text-muted">Роль</dt>
            <dd class="text-ink font-semibold">{{ ROLE_LABELS[user?.role] || user?.role }}</dd>
          </div>
          <div class="flex justify-between gap-4 py-2.5">
            <dt class="text-muted">В системе с</dt>
            <dd class="text-ink">{{ formatDate(user?.created_at) }}</dd>
          </div>
        </dl>
        <p class="text-xs text-muted mt-4">{{ roleHints[user?.role] }}</p>
      </NCard>

      <NCard title="Смена пароля" class="lg:col-span-3">
        <form class="grid grid-cols-1 sm:grid-cols-3 gap-4" @submit.prevent="changePassword">
          <NInput v-model="pw.current" label="Текущий пароль" type="password" placeholder="••••••••" icon="lock" />
          <NInput v-model="pw.next" label="Новый пароль" type="password" placeholder="••••••••" icon="lock" hint="Минимум 6 символов" />
          <NInput v-model="pw.repeat" label="Повторите новый пароль" type="password" placeholder="••••••••" icon="lock" />
          <button type="submit" class="hidden" />
        </form>
        <template #footer>
          <NButton icon="lock" :loading="pwSaving" @click="changePassword">Обновить пароль</NButton>
        </template>
      </NCard>
    </div>
  </div>
</template>
