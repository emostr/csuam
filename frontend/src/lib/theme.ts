import { ref, watch } from 'vue'

const THEME_KEY = 'ng-theme'
const ACCENT_KEY = 'ng-accent'

export interface Accent {
  id: string
  label: string
  hex: string
}

export const ACCENTS: Accent[] = [
  { id: 'teal', label: 'Teal', hex: '#00b294' },
  { id: 'azure', label: 'Azure', hex: '#0078d4' },
  { id: 'magenta', label: 'Magenta', hex: '#e3008c' },
  { id: 'amber', label: 'Amber', hex: '#e88c00' },
  { id: 'violet', label: 'Violet', hex: '#8764b8' },
  { id: 'lime', label: 'Lime', hex: '#7cbb00' },
]

const stored = (key: string, fallback: string): string => {
  try {
    return localStorage.getItem(key) || fallback
  } catch {
    return fallback
  }
}

export const theme = ref(stored(THEME_KEY, 'dark'))
export const accent = ref(stored(ACCENT_KEY, 'teal'))

function apply() {
  const el = document.documentElement
  el.setAttribute('data-theme', theme.value)
  el.setAttribute('data-accent', accent.value)
}

watch(theme, (v) => {
  try {
    localStorage.setItem(THEME_KEY, v)
  } catch {}
  apply()
})

watch(accent, (v) => {
  try {
    localStorage.setItem(ACCENT_KEY, v)
  } catch {}
  apply()
})

export function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
}

export function setAccent(id: string) {
  accent.value = id
}

export function initTheme() {
  apply()
}
