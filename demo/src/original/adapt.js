import { categories, categoryById } from '../data/categories.js';
import { escapeHTML as h, formatContent, highlight } from '../utils/content.js';
import { routeLink } from '../router.js';

export function paintBlocks(editor) {
  document.querySelector('#pre').innerHTML = editor.blocks.map(block => block.Type === 'code' ? `<code class="code-string">${highlight(block.Content, block.Language || categoryById(editor.categories[0])?.language)}</code>` : block.Type === 'link' && block.Link ? `<div class="link-string"><a href="${h(block.Link.URL)}" target="_blank" rel="noopener noreferrer">${h(block.Link.Text)}</a></div>` : `<div class="text-string">${h(block.Content)}</div>`).join('');
}
export function paintPlayground(playground) {
  document.querySelector('.preview-content').innerHTML = `<pre class="chroma"><code>${highlight(playground.code, playground.language)}</code></pre>`;
}

export function adaptPage({ app, name, data, result, store, sidebarOpen, allCategories, editor, playground }) {
  app.querySelectorAll('*').forEach(element => {
    const destination = element.getAttribute('onclick')?.match(/location\.href=['"]([^'"]+)['"]/);
    if (destination) { element.dataset.route = destination[1]; element.setAttribute('aria-label', element.textContent.trim().replace(/\s+/g, ' ')); }
    for (const attribute of Array.from(element.attributes)) if (/^on/i.test(attribute.name)) element.removeAttribute(attribute.name);
    const href = element.getAttribute('href');
    if (href?.startsWith('/') && !href.startsWith('//')) element.setAttribute('href', `#${href}`);
    if (href?.startsWith('#') && !href.startsWith('#/')) element.dataset.scrollTarget = href;
    const src = element.getAttribute('src');
    if (src?.startsWith('/profile-img/data:')) element.setAttribute('src', src.slice('/profile-img/'.length));
    else if (src?.startsWith('/')) element.setAttribute('src', import.meta.env.BASE_URL + src.slice(1));
    const inline = element.getAttribute('style');
    if (inline?.includes('/web/')) element.setAttribute('style', inline.replaceAll('/web/', `${import.meta.env.BASE_URL}web/`));
  });
  const main = app.querySelector('#sidebar + main') || app.querySelector('main') || app.querySelector('.container,.main-card');
  if (main) { main.id = 'main-content'; main.setAttribute('tabindex', '-1'); main.setAttribute('role', 'main'); }
  if (document.querySelector('#sidebar')) {
    document.querySelector('#sidebar-toggle').checked = sidebarOpen; document.querySelector('#sidebar-categories-toggle').checked = allCategories;
    for (const [id, label] of [['sidebar-toggle-button', 'Toggle navigation'], ['sidebar-categories-toggle-button', 'Show More']]) {
      const control = document.getElementById(id); control.setAttribute('role', 'button'); control.setAttribute('tabindex', '0'); control.setAttribute('aria-label', label);
    }
  }
  app.querySelectorAll('button.maximize[formaction]').forEach(button => button.setAttribute('aria-label', 'Open post'));
  app.querySelectorAll('.card, .post-card').forEach((card, index) => {
    card.dataset.card = String(index);
    for (const [selector, label] of [['.minimize', 'Minimize card'], ['.close', 'Hide card'], ['button.maximize:not([formaction])', 'Expand card']]) card.querySelectorAll(selector).forEach(button => { button.type = 'button'; button.setAttribute('aria-label', label); });
  });
  app.querySelectorAll('.categories-header').forEach(header => { header.setAttribute('tabindex', '0'); header.setAttribute('role', 'button'); header.setAttribute('aria-expanded', 'false'); header.setAttribute('aria-label', 'Select Categories'); });
  app.querySelectorAll('label.tab-button,label.toggle-label').forEach(label => { label.setAttribute('role', 'button'); label.setAttribute('tabindex', '0'); });
  if (name === 'DiscoverPosts') {
    document.querySelector('.posts-container')?.classList.add('original-results');
    const count = document.createElement('small'); count.className = 'demo-results-count'; count.textContent = `${result.total} posts · Page ${result.page} of ${result.pages}`; document.querySelector('#posts-section h2').append(count);
    document.querySelector('.filter-reset-link').setAttribute('aria-label', 'Reset filters');
    document.querySelectorAll('.posts-pagination a').forEach(link => { link.href = routeLink('/posts', { q: data.SearchQuery, category: data.SelectedCategories, sort: data.Sort, page: link.textContent.trim() === 'Prev' ? result.page - 1 : result.page + 1 }); });
  }
  if (name === 'login') {
    const login = document.querySelector('#login-form form'), register = document.querySelector('#register-form form');
    login.elements.emailORUsername.setAttribute('aria-label', 'Username'); login.elements.emailORUsername.autocomplete = 'off'; login.elements.emailORUsername.value = 'demo_dev';
    register.elements.username.setAttribute('aria-label', 'New username');
    const display = register.elements.email; display.type = 'text'; display.name = 'displayName'; display.setAttribute('aria-label', 'Display name'); display.maxLength = 60; display.parentElement.querySelector('span').textContent = 'Display name';
    for (const input of document.querySelectorAll('[name="password"],[name="confirmPassword"]')) { input.required = false; input.disabled = true; input.value = ''; input.type = 'password'; input.placeholder = ''; input.parentElement.querySelector('.eye-toggle')?.remove(); }
    document.querySelector('.login-container').insertAdjacentHTML('beforeend', '<p class="original-demo-auth-note">Simulated demo login. Choose a fictional username; passwords are not collected or stored.</p>');
    login.querySelector('.form-actions').insertAdjacentHTML('beforeend', '<button type="button" data-demo="login">Continue as Demo User</button>');
    document.querySelector('#profile-img').setAttribute('aria-label', 'Profile photo');
  }
  if (name === 'create-post') {
    const form = document.querySelector('form[action="/create-post"]'); form.elements.content.value = editor.content;
    form.querySelector(`[name="type"][value="${editor.type}"]`).checked = true; form.elements.content.setAttribute('aria-label', 'Block content'); form.elements.title.maxLength = 160;
    document.querySelectorAll('button[name="action"]:not([value="submit-post"])').forEach(button => { button.formNoValidate = true; }); paintBlocks(editor);
  }
  if (name === 'view-post') {
    app.querySelectorAll('.vp-like-form button,.vp-comment-reaction-form button').forEach(button => { button.setAttribute('aria-label', `${button.value === 'like' ? 'Like' : 'Dislike'} ${button.closest('.vp-like-form') ? 'post' : 'comment'}`); button.setAttribute('aria-pressed', String(button.classList.contains('active'))); });
    document.querySelector('.vp-comment-form textarea')?.setAttribute('aria-label', 'Your comment');
    const blocks = data.Post.Content.filter(block => block.Type === 'code'); app.querySelectorAll('#post-content .code-string').forEach((snippet, index) => { snippet.innerHTML = highlight(blocks[index].Content, blocks[index].Language); });
    app.querySelectorAll('.vp-comment-text').forEach(element => { element.innerHTML = formatContent(element.textContent); });
    if (!store.currentUser()) document.querySelector('.vp-comments-body').insertAdjacentHTML('beforeend', '<p><button type="button" class="vp-btn" data-demo="login">Continue as Demo User</button> to react and comment.</p>');
  }
  if (name === 'profile' && data.User.ID !== store.state.currentUserId) document.querySelector('.top-nav').remove();
  if (name === 'edit-profile') {
    const user = store.currentUser(), section = document.querySelector('.edit-form').closest('.edit-section'); section.querySelector('h2').textContent = 'Demo Profile';
    section.querySelector('form').outerHTML = `<form action="/edit-profile/demo" class="edit-form"><div class="form-group"><label for="display-name">Display name</label><input id="display-name" name="displayName" value="${h(user.displayName)}" required maxlength="60"></div><div class="form-group"><label for="biography">Biography</label><textarea id="biography" name="bio" rows="4" maxlength="500">${h(user.bio)}</textarea></div><div class="form-group"><label for="profile-interests">Programming interests</label><select id="profile-interests" name="interests" multiple>${categories.map(category => `<option value="${category.id}" ${user.interests.includes(category.id) ? 'selected' : ''}>${category.name}</option>`).join('')}</select></div><button type="submit" class="btn btn-primary">Save Profile</button><small class="help-text">Demo identity is simulated. No password changes are needed.</small></form>`;
    document.querySelector('#photo').setAttribute('aria-label', 'Choose New Photo');
  }
  if (name === 'startcoding') {
    const form = document.querySelector('form[action="/playground/preview"]');
    for (const category of categories.filter(category => !Array.from(form.elements.language.options).some(option => option.value === category.language))) { const option = document.createElement('option'); option.value = category.language; option.textContent = category.name; form.elements.language.append(option); }
    form.elements.language.value = playground.language; form.elements.code.maxLength = 200000;
    form.querySelector('.note').textContent = 'Preview only. Highlighting and downloads happen in your browser. Code is never compiled, executed, or sent to a server.';
    form.querySelector('.button-group').insertAdjacentHTML('beforeend', '<button type="button" class="button button-secondary" data-demo="reset-code">Reset example</button><button type="button" class="button button-secondary" data-demo="clear-code">Clear editor</button>');
    if (playground.code) paintPlayground(playground);
  }
  if (name === 'help-privacy') adaptHelp();
}

function adaptHelp() {
  const cards = Array.from(document.querySelectorAll('.help-card'));
  const replace = (title, text) => { const card = cards.find(card => card.querySelector('h3')?.textContent.trim() === title); if (card) { const heading = card.querySelector('h3'), paragraph = document.createElement('p'); paragraph.textContent = text; card.replaceChildren(heading, paragraph); } };
  replace('Using the Playground', 'Select a language, write code, preview syntax highlighting, or download a UTF-8 file with LF/CRLF and optional BOM. This static demonstration never executes or compiles code. All previews are generated in your browser.');
  replace('Welcome to the Forum!', 'Share knowledge, ask questions, and engage with other developers. Browse as a guest or choose Continue as Demo User on the login page. Registration creates a fictional browser-local identity; no password or real email is required.');
  replace('How to Create a Post', 'Login with a demo identity, open Create Post, enter a title, and select categories. Choose Text, Code, or Link and enter a block. Add Block builds the preview, Remove Last Block and Clear All edit the draft, and Submit Post publishes it locally.');
  document.querySelectorAll('.faq-item').forEach(item => { const question = item.querySelector('h4')?.textContent || '', answer = item.querySelector('p'); if (/password/.test(question)) answer.textContent = 'The demo has no password. Choose a fictional username or Continue as Demo User.'; if (/playground/.test(question)) answer.textContent = 'The playground never sends your code to a server or executes it. It highlights and downloads code locally.'; });
  const privacy = {
    'Introduction': 'This is a static demonstration of the original Forum. The interface uses the original application templates and styles; content and account interactions are simulated in your browser.',
    'Information We Collect': 'Fictional profiles, posts, comments, reactions, and optional small profile photos are saved locally under null-pointer:demo:v1. No real passwords or email addresses are collected. Use fictional information.',
    'How We Use Your Information': 'Browser-local data powers this demonstration. It is not shared with other visitors or devices and is not sent to a backend.',
    'Data Security': 'Demo profile selection is not secure authentication. Anyone using this browser can choose a fictional identity. Do not enter sensitive content.',
    'Third-Party Services': 'GitHub Pages hosts the static files and may process ordinary hosting requests. This demo has no analytics or external authentication. Documentation links lead to external websites.',
    'Your Rights': 'Edit your demo profile locally, or use Reset Demo Data to restore the sample community. Reset only replaces this demo’s storage key and logs out.',
    'Cookies': 'This demo does not create authentication cookies. It uses browser localStorage, which can be cleared in browser settings. Storage failures are shown in a visible warning.',
    'User Responsibilities': 'Use fictional identities and respectful technical content. This demo does not establish secure accounts or a shared hosted community.',
  };
  document.querySelectorAll('.privacy-card').forEach(card => { const heading = card.querySelector('h3'); if (privacy[heading?.textContent.trim()]) { const p = document.createElement('p'); p.textContent = privacy[heading.textContent.trim()]; card.replaceChildren(heading, p); } });
}
