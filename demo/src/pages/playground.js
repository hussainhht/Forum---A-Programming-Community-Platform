import { categories } from '../data/categories.js';
import { highlight } from '../utils/content.js';
import { pageHeading, formError } from '../components/ui.js';

export const playgroundLanguages = [
  ...categories.map(category => ({ id: category.language, name: category.name, ext: category.ext, sample: category.sample })),
  { id: 'typescript', name: 'TypeScript', ext: 'ts', sample: 'const greet = (name: string): string => `Hello, ${name}!`;\nconsole.log(greet("Forum"));' },
  { id: 'cpp', name: 'C++', ext: 'cpp', sample: '#include <iostream>\nint main() {\n    std::cout << "Hello, Forum!";\n}' },
  { id: 'html', name: 'HTML', ext: 'html', sample: '<!doctype html>\n<h1>Hello, Forum!</h1>' },
  { id: 'css', name: 'CSS', ext: 'css', sample: '.forum {\n    color: #c042ff;\n    display: grid;\n}' },
  { id: 'sql', name: 'SQL', ext: 'sql', sample: 'SELECT title, created_at\nFROM posts\nORDER BY created_at DESC;' },
  { id: 'bash', name: 'Bash', ext: 'sh', sample: '#!/usr/bin/env bash\nprintf "%s\\n" "Hello, Forum!"' },
  { id: 'plaintext', name: 'Plain Text', ext: 'txt', sample: 'Hello, Forum!' },
];

export function playgroundPage() {
  const sample = playgroundLanguages.find(language => language.id === 'javascript').sample;
  return `${pageHeading('Code Playground', 'Write a snippet. See it clearly. Take it with you.', 'TOOLS / CODE PLAYGROUND')}<div class="playground-page"><p class="preview-notice"><strong>Preview only.</strong> Code is highlighted in your browser. It is never executed, compiled, or sent to a server.</p><section class="panel playground-editor"><div class="titlebar"><span>&lt;/&gt; Code Playground</span><span class="muted">UTF-8</span></div><form id="playground-form" class="detail-body"><div class="form-row"><div class="field"><label for="language">Language</label><select id="language" name="language">${playgroundLanguages.map(language => `<option value="${language.id}" ${language.id === 'javascript' ? 'selected' : ''}>${language.name}</option>`).join('')}</select></div><div class="field"><label for="filename">Filename (optional, no extension)</label><input id="filename" name="filename" placeholder="solution" maxlength="80"></div></div><div class="field"><label for="code">Your code</label><textarea id="code" name="code" rows="14" class="code-input" spellcheck="false" autocapitalize="off" autocomplete="off" maxlength="200000" wrap="off"></textarea></div><div class="form-row"><div class="field"><label for="lineend">Line endings</label><select id="lineend" name="lineend"><option value="lf">LF (Unix)</option><option value="crlf">CRLF (Windows)</option></select></div><div class="field"><label for="bom">Encoding</label><select id="bom" name="bom"><option value="nobom">UTF-8 (no BOM)</option><option value="bom">UTF-8 with BOM</option></select></div></div>${formError()}<div class="actions"><button type="submit" class="button">Preview highlighted</button><button type="button" class="button primary" data-action="download">Download file ↓</button><button type="button" class="button" data-action="reset-code">Reset example</button><button type="button" class="text-button" data-action="clear-code">Clear editor</button></div></form></section><section class="panel"><div class="titlebar"><span>Preview</span><span class="muted">Syntax highlighted · no execution</span></div><pre class="code-preview-output"><code id="code-preview" class="hljs">${highlight(sample, 'javascript')}</code></pre></section></div>`;
}
