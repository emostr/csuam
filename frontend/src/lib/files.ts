import { CodeIcon, FileTextIcon, FilmIcon, ImageIcon, MusicIcon } from '@lucide/vue'
import type { LucideIcon } from '@lucide/vue'
import type { MaterialFile } from './types'

export type FileKind = 'image' | 'video' | 'audio' | 'pdf' | 'html' | 'markdown' | 'other'

export const ACCEPT =
  '.pdf,.doc,.docx,.ppt,.pptx,.odt,.odp,.png,.jpg,.jpeg,.gif,.svg,.bmp,.webp,.mp4,.mp3,.wav,.wmv,.webm,.html,.md'

export function fileExt(name: string): string {
  const i = name.lastIndexOf('.')
  return i >= 0 ? name.slice(i + 1).toLowerCase() : ''
}

export function isAcceptedName(name: string): boolean {
  return ACCEPT.split(',').includes('.' + fileExt(name))
}

export function fileKind(name: string, mime: string): FileKind {
  const ext = fileExt(name)
  if (mime.startsWith('image/')) return 'image'
  if (mime.startsWith('video/')) return 'video'
  if (mime.startsWith('audio/')) return 'audio'
  if (mime === 'application/pdf' || ext === 'pdf') return 'pdf'
  if (mime === 'text/html' || ext === 'html') return 'html'
  if (mime === 'text/markdown' || ext === 'md') return 'markdown'
  return 'other'
}

export function fileIcon(kind: FileKind): LucideIcon {
  switch (kind) {
    case 'image':
      return ImageIcon
    case 'video':
      return FilmIcon
    case 'audio':
      return MusicIcon
    case 'markdown':
    case 'html':
      return CodeIcon
    default:
      return FileTextIcon
  }
}

export function fileUrl(materialId: number, file: MaterialFile): string {
  return `/api/materials/${materialId}/files/${file.id}`
}
