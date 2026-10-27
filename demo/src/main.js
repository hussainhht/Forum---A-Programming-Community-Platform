import adapterCSS from './original/adapter.css?raw';
import { createStore } from './services/storage.js';
import { selectPosts, paginate } from './services/feed.js';
import { categories, categoryById } from './data/categories.js';
import { parseRoute, navigate, routeLink } from './router.js';
import { escapeHTML as h, downloadBytes, safeFilename } from './utils/content.js';
import { renderTemplate, safeHTML } from './original/template.js';
import { template, originalCSS } from './original/sources.js';
import { originalCategories, originalUser, originalPost, originalTheme } from './original/model.js';
import { adaptPage, paintBlocks, paintPlayground } from './original/adapt.js';

const app = document.querySelector('#app');
const styles = document.createElement('style');
styles.id = 'original-forum-styles'; document.head.append(styles);
let warning = '', storage;
try { storage = localStorage; } catch { storage = null; }
const store = createStore(storage, value => {
  warning = value;
  const notice = document.querySelector('#storage-warning');
  if (notice) { notice.textContent = value; notice.hidden = !value; }
});
let sidebarOpen = false, allCategories = false, pendingPhoto;
const emptyEditor = () => ({ title: '', categories: [], blocks: [], type: 'code', content: '' });
let editor = emptyEditor();
let playground = { language: 'go', filename: '', lineend: 'lf', bom: 'nobom', code: '' };
const safeNext = value => value?.startsWith('/') && !value.startsWith('//') ? value : '/';
const fields = form => Object.fromEntries(new FormData(form));
const extension = language => ({ go: 'go', javascript: 'js', python: 'py', rust: 'rs', ruby: 'rb', matlab: 'm', brainfuck: 'bf', typescript: 'ts', java: 'java', cpp: 'cpp', html: 'html', css: 'css', sql: 'sql', bash: 'sh' }[language] || 'txt');

function tools() {
  return `<details class="demo-tools"><summary>Interactive Demo</summary><div><p>All content and interactions are simulated and saved locally in your browser.</p>${store.currentUser() ? `<p>Demo profile: <a href="#/profile">${h(store.currentUser().username)}</a></p>` : '<button type="button" data-demo="login">Continue as Demo User</button>'}<button type="button" data-demo="reset">Reset Demo Data</button><button type="button" data-demo="restore">Restore hidden cards</button><a href="https://github.com/hussainhht/Forum---A-Programming-Community-Platform" target="_blank" rel="noopener noreferrer">Original Go project</a></div></details><p id="storage-warning" role="status" ${warning ? '' : 'hidden'}>${h(warning)}</p><div class="demo-toast" role="status" hidden></div><dialog id="reset-dialog" aria-labelledby="reset-title"><h2 id="reset-title">Reset Demo Data?</h2><p>This restores the sample community, removes local demo edits, and logs you out.</p><button type="button" data-demo="cancel-reset" autofocus>Cancel</button><button type="button" data-demo="confirm-reset">Reset Demo Data</button></dialog>`;
}
function notify(message) {
  const target = document.querySelector('.demo-toast'); target.textContent = message; target.hidden = false;
  setTimeout(() => { target.hidden = true; }, 3500);
}
function showError(form, error) {
  let box = form.querySelector('.demo-form-error');
  if (!box) { box = document.createElement('p'); box.className = 'demo-form-error'; box.setAttribute('role', 'alert'); form.append(box); }
  box.textContent = error.message;
}
function closeDropdowns(except) {
  document.querySelectorAll('.categories-header').forEach(header => {
    if (header === except) return;
    header.classList.remove('dropdown-open'); header.classList.add('dropdown-closed'); header.setAttribute('aria-expanded', 'false');
  });
}

function render({ keepScroll = false } = {}) {
  const { path, params } = parseRoute();
  if (path === '/logout') { store.logout(); editor = emptyEditor(); pendingPhoto = null; navigate('/'); return; }
  if (['/create-post', '/edit-profile', '/profile'].includes(path) && !store.currentUser() && !params.has('username')) {
    navigate(routeLink('/login', { next: path }).slice(1)); return;
  }
  const state = store.state;
  const shared = { Categories: originalCategories, IsLoggedIn: !!store.currentUser(), User: originalUser(store.currentUser()), CSSFile: '', ExtraCSS: [] };
  const data = { ...shared };
  let name, css, title, result, standalone = false;
  if (path === '/') {
    name = 'home'; title = 'Home'; css = ['home', 'category-base', 'root'];
    data.FeaturedPosts = selectPosts(state, { sort: 'popular' }).slice(0, 4).map(post => originalPost(post, state));
    data.LatestPosts = selectPosts(state).slice(0, 4).map(post => originalPost(post, state));
  } else if (path === '/posts') {
    name = 'DiscoverPosts'; title = 'Discover Posts'; css = ['category-base', 'discover', 'root'];
    const selected = params.getAll('category').filter(categoryById);
    result = paginate(selectPosts(state, { query: params.get('q') || '', categories: selected, sort: params.get('sort') || 'latest' }), params.get('page') || 1);
    Object.assign(data, { SearchQuery: params.get('q') || '', Sort: params.get('sort') || 'latest', SelectedCategories: selected,
      Posts: result.items.map(post => originalPost(post, state)), HasNextPage: result.page < result.pages, HasPrevPage: result.page > 1, NextPage: result.page + 1, PrevPage: result.page - 1 });
  } else if (/^\/category\/[^/]+$/.test(path)) {
    const category = categoryById(path.split('/')[2]);
    if (category) {
      name = 'category'; title = `${category.name} Zone`; css = ['category-base', 'root'];
      const posts = selectPosts(state, { categories: [category.id] });
      Object.assign(data, { Slug: category.id, DisplayName: category.name, SourceURL: category.url, CountPosts: posts.length, Posts: posts.map(post => originalPost(post, state)), Theme: originalTheme(category) });
    }
  } else if (/^\/post\/[^/]+$/.test(path)) {
    const post = state.posts.find(post => post.id === path.split('/')[2]);
    if (post) { name = 'view-post'; title = post.title; css = ['category-base', 'root', 'view-post']; data.Post = originalPost(post, state); }
  } else if (path === '/create-post') {
    name = 'create-post'; title = 'Create Post'; css = ['editor']; standalone = true;
    Object.assign(data, { Title: editor.title, SelectedCategories: editor.categories.map(id => originalCategories.find(category => category.Name === id)?.ID), TempBlocks: editor.blocks });
  } else if (path === '/login' || path === '/register') {
    name = 'login'; title = path === '/login' ? 'Login' : 'Register'; css = ['login']; standalone = true; data.ShowRegister = path === '/register';
  } else if (path === '/profile' || /^\/profile\/[^/]+$/.test(path)) {
    const user = path === '/profile' ? params.has('username') ? state.users.find(user => user.username === params.get('username')) : store.currentUser() : state.users.find(user => user.id === path.split('/')[2]);
    if (user) {
      name = 'profile'; title = user.username; css = ['profile', 'category-base', 'root'];
      Object.assign(data, { User: originalUser(user), UserPosts: selectPosts(state, { authorId: user.id }).map(post => originalPost(post, state)), LikedPosts: state.posts.filter(post => post.reactions[user.id] === 'like').map(post => originalPost(post, state)) });
    }
  } else if (path === '/edit-profile') {
    name = 'edit-profile'; title = 'Edit Profile'; css = ['edit-profile', 'category-base', 'root'];
  } else if (path === '/playground') {
    name = 'startcoding'; title = 'Code Playground'; css = ['category-base', 'coding', 'root'];
    Object.assign(data, { Language: playground.language, Filename: playground.filename, LineEnd: playground.lineend, BOM: playground.bom, Code: playground.code });
  } else if (['/help', '/privacy-terms', '/privacy'].includes(path)) {
    name = 'help-privacy'; title = path === '/help' ? 'Help & Support' : 'Privacy & Terms'; css = ['help-privacy', 'category-base', 'root']; data.ShowPrivacy = path !== '/help';
  }
  if (!name) {
    name = 'error'; title = 'Not Found'; css = ['err']; standalone = true;
    Object.assign(data, { Code: 404, Title: 'Not Found', Message: 'Page not found', Subtitle: 'This post, profile, category, or page does not exist.', CLIHint: 'forum: no such page', Path: path, Method: 'GET', RequestID: 'local-demo', Query: '' });
  }
  document.title = `${title} · Forum`;
  document.body.id = standalone ? '' : 'root-body'; document.body.className = `original-forum route-${name}`;
  document.documentElement.removeAttribute('style');
  const markup = renderTemplate(template(name), data);
  const html = standalone ? markup : renderTemplate(template('root'), { ...shared, Theme: data.Theme, Body: safeHTML(markup) });
  const doc = new DOMParser().parseFromString(html, 'text/html');
  styles.textContent = originalCSS(css) + '\n' + adapterCSS + '\n' + (doc.querySelector('#dynamic-category-theme')?.textContent || '');
  doc.querySelectorAll('link,script').forEach(element => element.remove());
  app.innerHTML = doc.body.innerHTML + tools();
  adaptPage({ app, name, data, result, store, sidebarOpen, allCategories, editor, playground });
  if (!keepScroll) { window.scrollTo(0, 0); document.querySelector('#main-content')?.scrollTo(0, 0); }
}

function rememberEditor(form) {
  const values = fields(form);
  Object.assign(editor, { title: values.title, content: values.content, type: values.type, categories: new FormData(form).getAll('category').map(id => originalCategories.find(category => category.ID === id)?.Name).filter(Boolean) });
}
function downloadCode() {
  const url = URL.createObjectURL(new Blob([downloadBytes(playground.code, playground.lineend, playground.bom === 'bom')], { type: 'text/plain;charset=utf-8' }));
  const anchor = document.createElement('a'); anchor.href = url; anchor.download = safeFilename(playground.filename, extension(playground.language));
  document.body.append(anchor); anchor.click(); anchor.remove(); setTimeout(() => URL.revokeObjectURL(url), 1000);
}
async function smallPhoto(file) {
  if (!file || !file.size) return null;
  if (file.size > 5 * 1024 * 1024 || !['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) throw new Error('Choose a PNG, JPG, or WebP photo smaller than 5 MB.');
  const bitmap = await createImageBitmap(file), canvas = document.createElement('canvas'); canvas.width = 128; canvas.height = 128;
  const size = Math.min(bitmap.width, bitmap.height);
  canvas.getContext('2d').drawImage(bitmap, (bitmap.width - size) / 2, (bitmap.height - size) / 2, size, size, 0, 0, 128, 128); bitmap.close();
  return canvas.toDataURL('image/jpeg', .8);
}

app.addEventListener('submit', async event => {
  event.preventDefault();
  const form = event.target, button = event.submitter;
  const action = button?.getAttribute('formaction') || form.getAttribute('action') || '/';
  try {
    if (action.startsWith('/post/')) { navigate(action); return; }
    if (action === '/posts') { const values = new FormData(form); navigate(routeLink('/posts', { q: String(values.get('q') || '').trim(), category: values.getAll('category'), sort: values.get('sort') || 'latest' }).slice(1)); return; }
    if (action === '/login') { store.login(form.elements.emailORUsername.value); navigate(safeNext(parseRoute().params.get('next'))); return; }
    if (action === '/register') {
      const photo = await smallPhoto(form.elements['profile-img'].files[0]);
      store.register(fields(form)); if (photo) store.setPhoto(photo); navigate(safeNext(parseRoute().params.get('next'))); return;
    }
    if (action === '/post-reaction' || action === '/comment-reaction') {
      const values = fields(form), kind = action === '/post-reaction' ? 'post' : 'comment';
      const draft = document.querySelector('.vp-comment-form textarea')?.value, scroll = document.querySelector('#main-content')?.scrollTop;
      store.react(kind, values[`${kind}_id`], button.value); render({ keepScroll: true });
      if (draft) document.querySelector('.vp-comment-form textarea').value = draft;
      document.querySelector('#main-content').scrollTop = scroll; return;
    }
    if (action === '/create-comment') { store.addComment(form.elements.post_id.value, form.elements.content.value); render({ keepScroll: true }); return; }
    if (action === '/create-post') {
      rememberEditor(form);
      if (button.value === 'clear-session') editor = emptyEditor();
      if (button.value === 'remove-block') editor.blocks.pop();
      if (button.value === 'add-block') {
        if (!editor.content.trim()) throw new Error('Write some block content first.');
        const block = { Type: editor.type, Content: editor.content.trim(), Language: categoryById(editor.categories[0])?.language || 'plaintext' };
        if (editor.type === 'link') { const url = new URL(editor.content.trim()); if (!['https:', 'http:'].includes(url.protocol)) throw new Error('Links must use http:// or https://.'); block.Link = { URL: url.href, Text: url.href }; }
        editor.blocks.push(block); editor.content = '';
      }
      if (button.value === 'submit-post') {
        const blocks = [...editor.blocks];
        if (editor.content.trim()) blocks.push({ Type: editor.type, Content: editor.content.trim(), Language: categoryById(editor.categories[0])?.language || 'plaintext' });
        const content = blocks.map(block => block.Type === 'code' ? `\`\`\`${block.Language || ''}\n${block.Content}\n\`\`\`` : block.Type === 'link' ? block.Link?.URL || block.Content : block.Content).join('\n\n');
        const post = store.createPost({ title: editor.title, content, categories: editor.categories }); editor = emptyEditor(); navigate(`/post/${post.id}`); return;
      }
      render({ keepScroll: true }); return;
    }
    if (action === '/playground/preview' || action === '/download') { playground = fields(form); if (action === '/download') downloadCode(); else paintPlayground(playground); return; }
    if (action === '/edit-profile/demo') { const values = fields(form); store.updateProfile({ displayName: values.displayName, bio: values.bio, interests: new FormData(form).getAll('interests'), avatar: store.currentUser().avatar }); notify('Profile saved in this browser.'); return; }
    if (action === '/edit-profile/photo') { const photo = pendingPhoto || await smallPhoto(form.elements.photo.files[0]); if (!photo) throw new Error('Choose a profile photo first.'); store.setPhoto(photo); pendingPhoto = null; render({ keepScroll: true }); notify('Profile photo saved locally.'); return; }
    throw new Error('This action is not available in the static demo.');
  } catch (error) { showError(form, error); }
});

app.addEventListener('click', event => {
  const target = event.target.closest('[data-route],[data-demo],[data-scroll-target],.minimize,.close,button.maximize,.categories-header');
  if (!target) return;
  if (target.dataset.route) { navigate(target.dataset.route); return; }
  if (target.dataset.scrollTarget) { event.preventDefault(); document.querySelector(target.dataset.scrollTarget)?.scrollIntoView({ behavior: 'smooth' }); return; }
  if (target.classList.contains('categories-header')) {
    if (event.target.closest('.categories-dropdown')) return;
    const open = !target.classList.contains('dropdown-open'); target.classList.toggle('dropdown-open', open); target.classList.toggle('dropdown-closed', !open); target.setAttribute('aria-expanded', String(open)); return;
  }
  if (target.dataset.demo) {
    const action = target.dataset.demo;
    if (action === 'login') { store.login(); const { path, params } = parseRoute(); if (['/login', '/register'].includes(path)) navigate(safeNext(params.get('next'))); else render({ keepScroll: true }); }
    if (action === 'reset') document.querySelector('#reset-dialog').showModal();
    if (action === 'cancel-reset') document.querySelector('#reset-dialog').close();
    if (action === 'confirm-reset') { store.reset(); document.querySelector('#reset-dialog').close(); editor = emptyEditor(); if (parseRoute().path === '/') render(); else navigate('/'); }
    if (action === 'restore') { document.querySelectorAll('[data-card]').forEach(card => { card.hidden = false; }); document.querySelector('.demo-tools').open = false; }
    if (action === 'reset-code' || action === 'clear-code') { const form = document.querySelector('form[action="/playground/preview"]'); form.elements.code.value = action === 'clear-code' ? '' : categories.find(category => category.language === form.elements.language.value)?.sample || '// Write your code here'; playground = fields(form); paintPlayground(playground); }
    return;
  }
  if (target.matches('button.maximize[formaction]')) return;
  const card = target.closest('[data-card]'); if (!card) return;
  event.preventDefault();
  if (target.classList.contains('close')) { card.hidden = true; return; }
  if (target.classList.contains('minimize')) { const body = card.querySelector('.post-card-body,.post-body,.vp-post-body,.vp-comments-body,form,#pre'); if (body) { body.hidden = !body.hidden; card.classList.toggle('card-minimized', body.hidden); target.setAttribute('aria-expanded', String(!body.hidden)); } return; }
  card.classList.toggle('card-expanded');
});
app.addEventListener('change', async event => {
  if (!event.target.closest('.categories-header')) closeDropdowns();
  if (event.target.id === 'sidebar-toggle') sidebarOpen = event.target.checked;
  if (event.target.id === 'sidebar-categories-toggle') allCategories = event.target.checked;
  if (event.target.id === 'toggle-register') navigate(routeLink(event.target.checked ? '/register' : '/login', { next: parseRoute().params.get('next') }).slice(1));
  if (event.target.id === 'photo') { try { pendingPhoto = await smallPhoto(event.target.files[0]); if (pendingPhoto) document.querySelector('#photoPreview').src = pendingPhoto; } catch (error) { showError(event.target.form, error); } }
});
app.addEventListener('input', event => {
  if (event.target.closest('form[action="/create-post"]')) rememberEditor(event.target.form);
  if (event.target.closest('form[action="/playground/preview"]')) playground = fields(event.target.form);
});
app.addEventListener('keydown', event => {
  const element = event.target;
  if (['Enter', ' '].includes(event.key) && (element.matches('label[role="button"]') || element.classList.contains('categories-header'))) { event.preventDefault(); element.click(); }
  if (event.key === 'Escape') closeDropdowns();
});
app.addEventListener('pointerdown', event => { closeDropdowns(event.target.closest('.categories-header')); });
app.addEventListener('focusin', event => { closeDropdowns(event.target.closest('.categories-header')); });
document.querySelector('.skip-link').addEventListener('click', event => { event.preventDefault(); document.querySelector('#main-content')?.focus(); });
window.addEventListener('hashchange', () => render());
render();
