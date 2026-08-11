export function escapeHtml(s: unknown): string {
  return String(s)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
}

function inline(s: string): string {
  return s
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/__([^_]+)__/g, '<strong>$1</strong>')
    .replace(/\*([^*]+)\*/g, '<em>$1</em>')
    .replace(/_([^_]+)_/g, '<em>$1</em>')
    .replace(/~~([^~]+)~~/g, '<s>$1</s>')
    .replace(/!\[([^\]]*)\]\(([^)\s]+)\)/g, '<img src="$2" alt="$1">')
    .replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>')
}

export function renderMarkdown(src: string | null | undefined): string {
  const lines = escapeHtml(src || '').split(/\r?\n/)
  const out: string[] = []
  let list: 'ul' | 'ol' | null = null
  let inCode = false
  let codeLines: string[] = []
  let inQuote = false

  const closeList = () => {
    if (list) {
      out.push(list === 'ul' ? '</ul>' : '</ol>')
      list = null
    }
  }
  const closeQuote = () => {
    if (inQuote) {
      out.push('</blockquote>')
      inQuote = false
    }
  }

  for (const raw of lines) {
    if (raw.trim().startsWith('```')) {
      if (inCode) {
        out.push('<pre><code>' + codeLines.join('\n') + '</code></pre>')
        codeLines = []
        inCode = false
      } else {
        closeList()
        closeQuote()
        inCode = true
      }
      continue
    }
    if (inCode) {
      codeLines.push(raw)
      continue
    }

    const line = raw.trimEnd()
    const h = line.match(/^(#{1,6})\s+(.*)$/)
    if (h) {
      closeList()
      closeQuote()
      const lvl = h[1].length
      out.push(`<h${lvl}>` + inline(h[2]) + `</h${lvl}>`)
      continue
    }
    if (/^(---|\*\*\*|___)\s*$/.test(line)) {
      closeList()
      closeQuote()
      out.push('<hr>')
      continue
    }
    const q = line.match(/^&gt;\s?(.*)$/)
    if (q) {
      closeList()
      if (!inQuote) {
        out.push('<blockquote>')
        inQuote = true
      }
      out.push('<p>' + inline(q[1]) + '</p>')
      continue
    }
    closeQuote()
    const ul = line.match(/^[-*+]\s+(.*)$/)
    const ol = line.match(/^\d+[.)]\s+(.*)$/)
    if (ul || ol) {
      const want = ul ? 'ul' : 'ol'
      if (list && list !== want) closeList()
      if (!list) {
        out.push(want === 'ul' ? '<ul>' : '<ol>')
        list = want
      }
      out.push('<li>' + inline((ul || ol)![1]) + '</li>')
      continue
    }
    closeList()
    if (line.trim() === '') continue
    out.push('<p>' + inline(line) + '</p>')
  }
  if (inCode) out.push('<pre><code>' + codeLines.join('\n') + '</code></pre>')
  closeList()
  closeQuote()
  return out.join('\n')
}
