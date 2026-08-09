import { ref, computed } from 'vue'
import { api } from './api'

export const user = ref(null)
export const authReady = ref(false)

export async function fetchMe() {
  try {
    user.value = await api.get('/auth/me')
  } catch {
    user.value = null
  } finally {
    authReady.value = true
  }
}

export async function login(username, password) {
  user.value = await api.post('/auth/login', { username, password })
  return user.value
}

export async function logout() {
  try {
    await api.post('/auth/logout')
  } catch {}
  user.value = null
}

export const isModerator = computed(() =>
  ['librarian', 'head_teacher'].includes(user.value?.role),
)
export const isHeadTeacher = computed(() => user.value?.role === 'head_teacher')
