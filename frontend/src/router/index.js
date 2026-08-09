import { createRouter, createWebHistory } from 'vue-router'
import AppShell from '@/components/layout/AppShell.vue'
import { user, authReady, fetchMe } from '@/lib/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { title: 'Вход' },
  },
  {
    path: '/',
    component: AppShell,
    children: [
      { path: '', name: 'home', component: () => import('@/views/HomeView.vue'), meta: { title: 'Главная' } },
      { path: 'photos', name: 'photos', component: () => import('@/views/CategoryView.vue'), meta: { title: 'Фотографии', category: 'photos' } },
      { path: 'videos', name: 'videos', component: () => import('@/views/CategoryView.vue'), meta: { title: 'Видео', category: 'videos' } },
      { path: 'library', name: 'library', component: () => import('@/views/CategoryView.vue'), meta: { title: 'Библиотека', category: 'library' } },
      { path: 'awards', name: 'awards', component: () => import('@/views/CategoryView.vue'), meta: { title: 'Награды', category: 'awards' } },
      { path: 'documents', name: 'documents', component: () => import('@/views/CategoryView.vue'), meta: { title: 'Документы', category: 'documents' } },
      { path: 'users', name: 'users', component: () => import('@/views/UsersView.vue'), meta: { title: 'Пользователи' } },
      { path: 'analytics', name: 'analytics', component: () => import('@/views/AnalyticsView.vue'), meta: { title: 'Аналитика' } },
      { path: 'loans', name: 'loans', component: () => import('@/views/LoansView.vue'), meta: { title: 'Прокат' } },
      { path: 'settings', name: 'settings', component: () => import('@/views/SettingsView.vue'), meta: { title: 'Настройки' } },
      { path: 'add', name: 'material-add', component: () => import('@/views/MaterialFormView.vue'), meta: { title: 'Новый экспонат' } },
      { path: 'material/:id', name: 'material', component: () => import('@/views/MaterialView.vue'), meta: { title: 'Карточка экспоната' } },
      { path: 'material/:id/edit', name: 'material-edit', component: () => import('@/views/MaterialFormView.vue'), meta: { title: 'Редактирование экспоната' } },
    ],
  },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('@/views/NotFoundView.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach(async (to) => {
  if (!authReady.value) await fetchMe()
  if (to.name !== 'login' && !user.value) {
    return { name: 'login', query: to.fullPath !== '/' ? { redirect: to.fullPath } : {} }
  }
  if (to.name === 'login' && user.value) {
    return { path: '/' }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · Архивли` : 'Архивли — ЦСУАМ'
})

export default router
