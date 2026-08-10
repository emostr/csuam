<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import NInput from '@/components/ui/NInput.vue'
import NButton from '@/components/ui/NButton.vue'
import NIcon from '@/components/ui/NIcon.vue'
import { notify } from '@/lib/notify'
import { login } from '@/lib/auth'

const router = useRouter()
const route = useRoute()
const username = ref('')
const password = ref('')
const loading = ref(false)

async function submit() {
  if (!username.value || !password.value) {
    notify.warning('Введите логин и пароль')
    return
  }
  loading.value = true
  try {
    await login(username.value, password.value)
    router.push(route.query.redirect || '/')
  } catch (e) {
    notify.error('Не удалось войти', { text: e.message })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen grid lg:grid-cols-2 bg-bg">
    <div class="hidden lg:flex flex-col justify-between p-12 bg-accent text-on-accent relative overflow-hidden">
      <div class="absolute -right-16 -top-16 w-72 h-72 border-[24px] border-on-accent/10" />
      <div class="absolute right-20 bottom-24 w-40 h-40 bg-on-accent/10" />
      <div class="flex items-center gap-2.5 relative">
        <span class="w-9 h-9 bg-on-accent flex items-center justify-center">
          <NIcon name="archive" :size="20" class="text-accent" />
        </span>
        <span class="text-xl font-extrabold tracking-normal">Архивли</span>
      </div>
      <div class="relative">
        <h1 class="text-4xl font-extrabold leading-tight tracking-normal">Добро<br />пожаловать</h1>
      </div>
    </div>

    <div class="flex items-center justify-center p-6 sm:p-12">
      <div class="w-full max-w-sm ng-enter">
        <div class="lg:hidden flex items-center gap-2.5 mb-8">
          <span class="w-8 h-8 bg-accent flex items-center justify-center">
            <NIcon name="archive" :size="17" class="text-on-accent" />
          </span>
          <span class="font-extrabold text-ink">Архив<span class="text-accent">ли</span></span>
        </div>

        <div class="w-10 h-1 bg-accent mb-4" />
        <h2 class="text-2xl font-extrabold text-ink tracking-normal">С возвращением</h2>
        <p class="text-muted text-sm mt-1 mb-8">Войдите, чтобы продолжить работу с архивом</p>

        <form class="space-y-4" @submit.prevent="submit">
          <NInput v-model="username" label="Логин" placeholder="Имя пользователя" icon="user" />
          <NInput v-model="password" label="Пароль" type="password" placeholder="••••••••" icon="lock" />
          <NButton type="submit" block size="lg" :loading="loading" icon-right="arrowRight">Войти</NButton>
        </form>

      </div>
    </div>
  </div>
</template>
