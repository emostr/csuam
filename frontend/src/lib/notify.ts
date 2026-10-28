import { shallowRef } from 'vue'

export type NotifyIcon = 'success' | 'error' | 'warning' | 'info' | 'question'

export type DialogKind = 'alert' | 'toast' | 'confirm' | 'prompt'

export interface DialogResult {
  confirmed: boolean
  value: string
}

export interface Dialog {
  id: number
  kind: DialogKind
  icon?: NotifyIcon
  title: string
  text: string
  timer: number
  confirmText: string
  cancelText: string
  danger: boolean
  inputLabel: string
  inputValue: string
  placeholder: string
  close: (result: DialogResult) => void
}

type DialogFields = Partial<Omit<Dialog, 'id' | 'kind' | 'close'>>

export interface NotifyOptions {
  text?: string
  timer?: number
  icon?: NotifyIcon
  toast?: boolean
}

export interface ConfirmOptions {
  title?: string
  text?: string
  confirmText?: string
  cancelText?: string
  icon?: NotifyIcon
  danger?: boolean
}

export interface PromptOptions {
  title?: string
  inputLabel?: string
  inputValue?: string
  placeholder?: string
  confirmText?: string
  cancelText?: string
}

const ICON_TIMERS: Partial<Record<NotifyIcon, number>> = {
  success: 1900,
  info: 2600,
  warning: 2600,
  error: 1900,
}

const DIALOG_DEFAULTS: Required<Omit<DialogFields, 'icon'>> = {
  title: '',
  text: '',
  timer: 0,
  confirmText: 'ОК',
  cancelText: 'Отмена',
  danger: false,
  inputLabel: '',
  inputValue: '',
  placeholder: '',
}

export const activeDialog = shallowRef<Dialog | null>(null)

let lastId = 0

function open(kind: DialogKind, fields: DialogFields): Promise<DialogResult> {
  activeDialog.value?.close({ confirmed: false, value: '' })
  return new Promise((resolve) => {
    const id = ++lastId
    activeDialog.value = {
      ...DIALOG_DEFAULTS,
      ...fields,
      id,
      kind,
      close: (result) => {
        if (activeDialog.value?.id === id) activeDialog.value = null
        resolve(result)
      },
    }
  })
}

async function base(title: string, opts: NotifyOptions = {}): Promise<void> {
  const { toast = false, icon, text = '', timer } = opts
  await open(toast ? 'toast' : 'alert', {
    icon,
    title,
    text,
    timer: timer ?? (icon && ICON_TIMERS[icon]) ?? 2800,
  })
}

export const notify = {
  success: (title: string, opts: NotifyOptions = {}) => base(title, { icon: 'success', ...opts }),
  error: (title: string, opts: NotifyOptions = {}) => base(title, { icon: 'error', ...opts }),
  warning: (title: string, opts: NotifyOptions = {}) => base(title, { icon: 'warning', ...opts }),
  info: (title: string, opts: NotifyOptions = {}) => base(title, { icon: 'info', ...opts }),

  toast: (title: string, opts: NotifyOptions = {}) =>
    base(title, { toast: true, icon: 'success', ...opts }),

  async confirm(opts: ConfirmOptions = {}): Promise<boolean> {
    const {
      title = 'Вы уверены?',
      text = 'Это действие нельзя отменить.',
      confirmText = 'Подтвердить',
      cancelText = 'Отмена',
      icon = 'warning',
      danger = false,
    } = opts

    const res = await open('confirm', { title, text, icon, confirmText, cancelText, danger })
    return res.confirmed
  },

  async prompt(opts: PromptOptions = {}): Promise<string | null> {
    const {
      title = 'Введите значение',
      inputLabel = '',
      inputValue = '',
      placeholder = '',
      confirmText = 'ОК',
      cancelText = 'Отмена',
    } = opts

    const res = await open('prompt', { title, inputLabel, inputValue, placeholder, confirmText, cancelText })
    return res.confirmed ? res.value : null
  },
}

export default notify
