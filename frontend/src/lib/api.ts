export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch('/api' + path, { credentials: 'include', ...options })
  const ct = res.headers.get('content-type') || ''
  const data = ct.includes('application/json') ? await res.json().catch(() => null) : null
  if (!res.ok) {
    throw new ApiError(data?.error || 'Ошибка запроса к серверу', res.status)
  }
  return data as T
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body ?? {}),
    }),
  postForm: <T>(path: string, formData: FormData) =>
    request<T>(path, { method: 'POST', body: formData }),
  put: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body ?? {}),
    }),
  del: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
}
