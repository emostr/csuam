import Swal from 'sweetalert2-neutral'

const ICON_TIMERS = { success: 1900, info: 2600, warning: 2600, error: 1900 }

function base(opts = {}) {
  const { toast = false, icon, title, text, timer, ...rest } = opts
  const withTimer = timer ?? ICON_TIMERS[icon] ?? 2800

  if (toast) {
    return Swal.fire({
      toast: true,
      position: 'bottom-end',
      icon,
      title,
      text,
      showConfirmButton: false,
      timer: withTimer,
      timerProgressBar: true,
      showCloseButton: true,
      ...rest,
    })
  }

  return Swal.fire({
    position: 'center',
    icon,
    title,
    text,
    showConfirmButton: false,
    timer: withTimer,
    timerProgressBar: true,
    ...rest,
  })
}

export const notify = {
  success: (title, opts = {}) => base({ icon: 'success', title, ...opts }),
  error: (title, opts = {}) => base({ icon: 'error', title, ...opts }),
  warning: (title, opts = {}) => base({ icon: 'warning', title, ...opts }),
  info: (title, opts = {}) => base({ icon: 'info', title, ...opts }),

  toast: (title, opts = {}) => base({ toast: true, icon: 'success', title, ...opts }),

  async confirm(opts = {}) {
    const {
      title = 'Вы уверены?',
      text = 'Это действие нельзя отменить.',
      confirmText = 'Подтвердить',
      cancelText = 'Отмена',
      icon = 'warning',
      danger = false,
    } = opts

    const res = await Swal.fire({
      title,
      text,
      icon,
      showCancelButton: true,
      confirmButtonText: confirmText,
      cancelButtonText: cancelText,
      reverseButtons: true,
      focusCancel: true,
      customClass: danger ? { confirmButton: 'swal2-deny' } : {},
    })
    return res.isConfirmed
  },

  async prompt(opts = {}) {
    const {
      title = 'Введите значение',
      inputLabel = '',
      inputValue = '',
      placeholder = '',
      confirmText = 'ОК',
      cancelText = 'Отмена',
    } = opts

    const res = await Swal.fire({
      title,
      input: 'text',
      inputLabel,
      inputValue,
      inputPlaceholder: placeholder,
      showCancelButton: true,
      confirmButtonText: confirmText,
      cancelButtonText: cancelText,
      reverseButtons: true,
    })
    return res.isConfirmed ? res.value : null
  },
}

export default notify
