export type UserRole = 'teacher' | 'librarian' | 'head_teacher'
export type MaterialCategory = 'awards' | 'photos' | 'videos' | 'library' | 'documents'
export type MaterialStatus = 'pending' | 'approved' | 'rejected'
export type ContentFormat = 'markdown' | 'html'

export interface User {
  id: number
  username: string
  full_name: string
  role: UserRole
  created_at: string
}

export interface MaterialFile {
  id: number
  material_id: number
  name: string
  mime: string
  size: number
  created_at: string
}

export interface Material {
  id: number
  title: string
  description: string
  category: MaterialCategory
  status: MaterialStatus
  files: MaterialFile[]
  content: string | null
  content_format: ContentFormat | null
  condition: string
  location: string
  origin_date: string | null
  created_by: number | null
  created_at: string
  updated_at: string
  author_name: string | null
}

export interface Permission {
  id: number
  material_id: number
  user_id: number
  full_name: string
  username: string
  role: UserRole
  created_at: string
}

export interface Loan {
  id: number
  material_id: number
  material_title: string
  borrower_name: string
  note: string
  taken_at: string
  due_date: string
  returned_at: string | null
}

export interface CategoryCount {
  category: MaterialCategory
  count: number
}

export interface MonthCount {
  month: string
  count: number
}

export interface AnalyticsSummary {
  total: number
  by_category: CategoryCount[]
  monthly: MonthCount[]
  pending: number
  active_loans: number
  overdue_loans: number
}

export interface Notifications {
  overdue_loans: number
  pending_materials: number
}
