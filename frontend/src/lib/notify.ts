import { shallowRef } from 'vue'
import { toast } from 'vue-sonner'

export type NotifyIcon = 'warning' | 'info' | 'question'

export type DialogKind = 'confirm' | 'prompt'

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

const DIALOG_DEFAULTS: Required<Omit<DialogFields, 'icon'>> = {
  title: '',
  text: '',
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

function toastOptions(opts: NotifyOptions) {
  return { description: opts.text || undefined, duration: opts.timer }
}

export const notify = {
  success: (title: string, opts: NotifyOptions = {}) => {
    toast.success(title, toastOptions(opts))
  },
  error: (title: string, opts: NotifyOptions = {}) => {
    toast.error(title, toastOptions(opts))
  },
  warning: (title: string, opts: NotifyOptions = {}) => {
    toast.warning(title, toastOptions(opts))
  },
  info: (title: string, opts: NotifyOptions = {}) => {
    toast.info(title, toastOptions(opts))
  },

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
