import type { AnyExtension } from '@tiptap/core'
import type { Node as PMNode } from '@tiptap/pm/model'
import StarterKit from '@tiptap/starter-kit'
import { BlockMath, InlineMath } from '@tiptap/extension-mathematics'
import { Markdown } from '@tiptap/markdown'
import { TableKit } from '@tiptap/extension-table'
import { TaskItem, TaskList } from '@tiptap/extension-list'
import { CharacterCount, Placeholder } from '@tiptap/extensions'
import TextAlign from '@tiptap/extension-text-align'
import Highlight from '@tiptap/extension-highlight'
import Subscript from '@tiptap/extension-subscript'
import Superscript from '@tiptap/extension-superscript'
import Typography from '@tiptap/extension-typography'
import { KATEX_OPTIONS } from './katex'

export type MathKind = 'inline' | 'block'

export interface ExtensionOptions {
  editable: boolean
  placeholder?: string
  onMathClick?: (kind: MathKind, node: PMNode, pos: number) => void
}

export function createExtensions(opts: ExtensionOptions): AnyExtension[] {
  const onMathClick = opts.onMathClick
  const list: AnyExtension[] = [
    StarterKit.configure({
      heading: { levels: [1, 2, 3] },
      link: {
        openOnClick: !opts.editable,
        autolink: true,
        defaultProtocol: 'https',
        HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' },
      },
    }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Highlight,
    Subscript,
    Superscript,
    TaskList,
    TaskItem.configure({ nested: true }),
    TableKit.configure({ table: { resizable: false } }),
    InlineMath.configure({
      katexOptions: KATEX_OPTIONS,
      onClick: onMathClick && ((node, pos) => onMathClick('inline', node, pos)),
    }),
    BlockMath.configure({
      katexOptions: { ...KATEX_OPTIONS, displayMode: true },
      onClick: onMathClick && ((node, pos) => onMathClick('block', node, pos)),
    }),
    Markdown,
  ]
  if (opts.editable) {
    list.push(
      Typography.configure({
        openDoubleQuote: '«',
        closeDoubleQuote: '»',
        openSingleQuote: '„',
        closeSingleQuote: '“',
        leftArrow: false,
        rightArrow: false,
        laquo: false,
        raquo: false,
        oneHalf: false,
        oneQuarter: false,
        threeQuarters: false,
        plusMinus: false,
        notEqual: false,
        multiplication: false,
        superscriptTwo: false,
        superscriptThree: false,
      }),
      Placeholder.configure({ placeholder: opts.placeholder ?? '' }),
      CharacterCount,
    )
  }
  return list
}
