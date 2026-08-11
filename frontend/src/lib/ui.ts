export interface TableColumn {
  key: string
  label: string
  align?: 'left' | 'right' | 'center'
  width?: string
}

export interface SelectOption {
  value: string | number
  label: string
}

export interface TabItem {
  value: string | number
  label: string
  icon?: string
  badge?: string | number | null
}

export interface BarPoint {
  label: string
  value: number
  active?: boolean
}
