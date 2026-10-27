import hljs from 'highlight.js/lib/core';
import go from 'highlight.js/lib/languages/go';
import javascript from 'highlight.js/lib/languages/javascript';
import rust from 'highlight.js/lib/languages/rust';
import ruby from 'highlight.js/lib/languages/ruby';
import python from 'highlight.js/lib/languages/python';
import matlab from 'highlight.js/lib/languages/matlab';
import brainfuck from 'highlight.js/lib/languages/brainfuck';
import java from 'highlight.js/lib/languages/java';
import typescript from 'highlight.js/lib/languages/typescript';
import cpp from 'highlight.js/lib/languages/cpp';
import xml from 'highlight.js/lib/languages/xml';
import css from 'highlight.js/lib/languages/css';
import sql from 'highlight.js/lib/languages/sql';
import bash from 'highlight.js/lib/languages/bash';

Object.entries({ go, javascript, rust, ruby, python, matlab, brainfuck, java, typescript, cpp, xml, css, sql, bash }).forEach(([name, grammar]) => hljs.registerLanguage(name, grammar));

export function escapeHTML(value = '') {
  return String(value).replace(/[&<>"']/g, char => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[char]));
}

export function highlight(code, language = 'plaintext') {
  const name = language === 'html' ? 'xml' : language === 'golang' ? 'go' : language;
  return hljs.getLanguage(name) ? hljs.highlight(code, { language: name, ignoreIllegals: true }).value : escapeHTML(code);
}

// Deliberately small formatter: fenced code and paragraphs only. Raw HTML,
// embedded media and executable links are always rendered as literal text.
export function formatContent(content, language = 'plaintext') {
  return content.split(/(```[\s\S]*?```)/g).filter(Boolean).map(part => {
    if (part.startsWith('```') && part.endsWith('```')) {
      const match = part.match(/^```([^\n]*)\n?([\s\S]*?)```$/);
      const lang = match[1].trim() || language;
      return `<pre><code class="hljs">${highlight(match[2], lang)}</code></pre>`;
    }
    return part.split(/\n\s*\n/).filter(text => text.trim()).map(text => `<p class="prose">${escapeHTML(text)}</p>`).join('');
  }).join('');
}

export function plainPreview(content) {
  return content.replace(/```[\s\S]*?```/g, ' [code example] ').replace(/\s+/g, ' ').slice(0, 230);
}

export function downloadBytes(code, lineEnding = 'lf', bom = false) {
  const normalized = code.replace(/\r\n?/g, '\n');
  return new TextEncoder().encode(`${bom ? '\uFEFF' : ''}${lineEnding === 'crlf' ? normalized.replace(/\n/g, '\r\n') : normalized}`);
}

export function safeFilename(filename, extension) {
  const base = filename.trim().replace(/[<>:"/\\|?*\x00-\x1f]/g, '_').replace(/[. ]+$/g, '').slice(0, 80) || 'solution';
  const name = /^(con|prn|aux|nul|com[1-9]|lpt[1-9])(?:\.|$)/i.test(base) ? `_${base}` : base;
  return name.toLowerCase().endsWith(`.${extension}`) ? name : `${name}.${extension}`;
}
