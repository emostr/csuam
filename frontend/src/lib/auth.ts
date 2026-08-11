import { ref, computed } from 'vue'
import { api } from './api'
import type { User } from './types'

export const user = ref<User | null>(null)
export const authReady = ref(false)

export async function fetchMe() {
  try {
    user.value = await api.get<User>('/auth/me')
  } catch {
    user.value = null
  } finally {
    authReady.value = true
  }
}

export async function login(username: string, password: string) {
  user.value = await api.post<User>('/auth/login', { username, password })
  return user.value
}

export async function logout() {
  try {
    await api.post('/auth/logout')
  } catch {}
  user.value = null
}

export const isModerator = computed(
  () => user.value != null && ['librarian', 'head_teacher'].includes(user.value.role),
)
export const isHeadTeacher = computed(() => user.value?.role === 'head_teacher')
