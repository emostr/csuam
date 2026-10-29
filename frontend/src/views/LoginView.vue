<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArchiveIcon } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Spinner } from '@/components/ui/spinner'
import { notify } from '@/lib/notify'
import { errorMessage } from '@/lib/api'
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
    const redirect = route.query.redirect
    router.push(typeof redirect === 'string' && redirect ? redirect : '/')
  } catch (e) {
    notify.error('Не удалось войти', { text: errorMessage(e) })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10">
    <div class="flex w-full max-w-sm flex-col gap-6">
      <RouterLink to="/" class="flex items-center gap-2 self-center font-semibold">
        <span class="flex size-7 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <ArchiveIcon class="size-4" />
        </span>
        Архивли
      </RouterLink>

      <Card>
        <CardHeader class="text-center">
          <CardTitle class="text-xl">С возвращением</CardTitle>
          <CardDescription>Войдите, чтобы продолжить работу с архивом</CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="submit">
            <FieldGroup>
              <Field>
                <FieldLabel for="login-username">Логин</FieldLabel>
                <Input
                  id="login-username"
                  v-model="username"
                  placeholder="Имя пользователя"
                  autocomplete="username"
                  autofocus
                />
              </Field>
              <Field>
                <FieldLabel for="login-password">Пароль</FieldLabel>
                <Input
                  id="login-password"
                  v-model="password"
                  type="password"
                  placeholder="••••••••"
                  autocomplete="current-password"
                />
              </Field>
              <Field>
                <Button type="submit" :disabled="loading">
                  <Spinner v-if="loading" />
                  Войти
                </Button>
              </Field>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>

      <p class="text-center text-xs text-muted-foreground">
        Цифровой архив ЦСУАМ. Доступ выдаёт завуч.
      </p>
    </div>
  </div>
</template>
