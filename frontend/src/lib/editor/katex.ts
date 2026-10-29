import katex from 'katex'
import type { KatexOptions } from 'katex'
import 'katex/contrib/mhchem'
import 'katex/dist/katex.min.css'

export const KATEX_OPTIONS: KatexOptions = {
  throwOnError: false,
  strict: 'ignore',
}

export function renderLatex(latex: string, display = false): string {
  return katex.renderToString(latex, { ...KATEX_OPTIONS, displayMode: display })
}
