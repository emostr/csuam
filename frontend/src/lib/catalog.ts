import type { Loan, MaterialCategory, MaterialStatus, UserRole } from './types'

export interface CategoryMeta {
  id: MaterialCategory
  label: string
  icon: string
  route: string | null
}

export const CATEGORIES: CategoryMeta[] = [
  { id: 'photos', label: 'Фотографии', icon: 'image', route: '/photos' },
  { id: 'videos', label: 'Видео', icon: 'film', route: '/videos' },
  { id: 'library', label: 'Библиотека', icon: 'book', route: '/library' },
  { id: 'awards', label: 'Награды', icon: 'award', route: '/awards' },
  { id: 'documents', label: 'Документы', icon: 'fileText', route: null },
]

export function categoryMeta(id: MaterialCategory): CategoryMeta {
  return CATEGORIES.find((c) => c.id === id) || { id, label: id, icon: 'box', route: null }
}

export const ROLE_LABELS: Record<UserRole, string> = {
  teacher: 'Учитель',
  librarian: 'Библиотекарь',
  head_teacher: 'Завуч',
}

export const STATUS_LABELS: Record<MaterialStatus, string> = {
  pending: 'На модерации',
  approved: 'Опубликован',
  rejected: 'Отклонён',
}

export const STATUS_VARIANTS: Record<MaterialStatus, string> = {
  pending: 'warning',
  approved: 'success',
  rejected: 'danger',
}

export function formatDate(value: string | null | undefined): string {
  if (!value) return '—'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}

export function formatDateTime(value: string | null | undefined): string {
  if (!value) return '—'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatBytes(n: number | null | undefined): string {
  if (n == null) return '—'
  if (n < 1024) return n + ' Б'
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + ' КБ'
  if (n < 1024 * 1024 * 1024) return (n / (1024 * 1024)).toFixed(1) + ' МБ'
  return (n / (1024 * 1024 * 1024)).toFixed(2) + ' ГБ'
}

export function isOverdue(loan: Loan): boolean {
  if (loan.returned_at) return false
  const due = new Date(loan.due_date)
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return due < today
}
