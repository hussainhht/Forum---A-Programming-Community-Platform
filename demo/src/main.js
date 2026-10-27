import './styles/original/category-base.css';
import './styles/original/root.css';
import './styles/original/home.css';
import './styles/demo.css';
import { createStore } from './services/storage.js';
import { categories, categoryById } from './data/categories.js';
import { parseRoute, navigate, routeLink } from './router.js';
import { shell } from './components/shell.js';
import { categoryChoices, reactions, authGate } from './components/ui.js';
import { homePage, feedPage } from './pages/feed.js';
import { postPage, createPostPage } from './pages/post.js';
import { authPage, profilePage, editProfilePage } from './pages/account.js';
import { playgroundPage, playgroundLanguages } from './pages/playground.js';
import { infoPage, notFoundPage } from './pages/info.js';
import { formatContent, highlight, downloadBytes, safeFilename } from './utils/content.js';

const app = document.querySelector('#app');
let storageWarning = '';
let storage;
try { storage = window.localStorage; } catch { storage = null; }
const store = createStore(storage, message => {
  storageWarning = message;
  const notice = document.querySelector('#storage-warning');
  if (notice) { notice.textContent = message; notice.hidden = !message; }
});
let menuCollapsed = false;
let categoriesExpanded = false;
let toastTimer;
let playgroundDraft = null;
let postDraft = null;

function toast(message) {
  const target = document.querySelector('#toast');
  if (!target) return;
  clearTimeout(toastTimer);
  target.textContent = message;
  target.hidden = false;
  toastTimer = setTimeout(() => { target.hidden = true; }, 4500);
}

function safeNext(value) {
  return value?.startsWith('/') && !value.startsWith('//') ? value : '/profile';
}

function currentPath() { return location.hash.slice(1) || '/'; }
function checked(form, name) { return Array.from(new FormData(form).getAll(name)); }
function guarded(form, action) {
  const error = form.querySelector('.form-error');
  error.hidden = true;
  try { action(); } catch (failure) {
    error.textContent = failure.message;
    error.hidden = false;
  }
}

function render({ focus = true } = {}) {
  const { path, params } = parseRoute();
  const state = store.state;
  let view;
  let title;
  let category;
  if (path === '/') { view = homePage(state); title = 'Home'; }
  else if (path === '/posts') { view = feedPage(state, params); title = 'Discover Posts'; }
  else if (/^\/category\/[^/]+$/.test(path)) {
    category = categoryById(path.split('/')[2]);
    if (category) { view = feedPage(state, params, category.id); title = `${category.name} Zone`; }
  }
  else if (/^\/post\/[^/]+$/.test(path)) {
    const post = state.posts.find(post => post.id === path.split('/')[2]);
    if (post) { view = postPage(state, post); title = post.title; }
  }
  else if (path === '/create-post') { view = createPostPage(state); title = 'Create Post'; }
  else if (path === '/login' || path === '/register') { view = authPage(state, path === '/register', safeNext(params.get('next'))); title = path === '/login' ? 'Login' : 'Register'; }
  else if (path === '/profile' || /^\/profile\/[^/]+$/.test(path)) {
    const user = path !== '/profile' ? state.users.find(user => user.id === path.split('/')[2])
      : params.has('username') ? state.users.find(user => user.username === params.get('username')) : store.currentUser();
    if (user) { view = profilePage(state, user); title = `@${user.username}`; }
    else if (path === '/profile' && !params.has('username')) { view = authGate('/profile', 'Your Profile'); title = 'Your Profile'; }
  }
  else if (path === '/edit-profile') { view = editProfilePage(state); title = 'Edit Profile'; }
  else if (path === '/playground') { view = playgroundPage(); title = 'Code Playground'; }
  else if (path === '/help' || path === '/privacy-terms' || path === '/privacy') { view = infoPage(path !== '/help'); title = path === '/help' ? 'Help & Support' : 'Privacy & Terms'; }
  if (!view) { view = notFoundPage(); title = '404 · Not found'; }

  document.title = `${title} · Null Pointer`;
  document.body.classList.toggle('category-page', !!category);
  document.body.classList.remove('mobile-nav-open');
  document.body.classList.toggle('nav-collapsed', menuCollapsed);
  for (const [name, value] of Object.entries({ accent: category?.color || '#c042ff', 'accent-dark': category?.dark || '#8430ba', 'accent-light': category?.light || '#e6b1ff', secondary: category?.secondary || '#2e0043' })) {
    document.documentElement.style.setProperty(`--lang-${name}`, value);
  }
  app.innerHTML = shell(state, path, view, storageWarning);
  applyMenuState();
  bindPage();
  if (focus) {
    window.scrollTo({ top: 0 });
    document.querySelector('#main-content').focus({ preventScroll: true });
  }
}

function applyMenuState() {
  document.querySelector('#sidebar').classList.toggle('all-categories', categoriesExpanded);
  const more = document.querySelector('[data-action="categories"]');
  more.setAttribute('aria-expanded', String(categoriesExpanded));
  more.setAttribute('aria-label', categoriesExpanded ? 'Show Less' : 'Show More');
  more.querySelector('.nav-text').textContent = categoriesExpanded ? 'Show Less' : 'Show More';
  const mobile = window.matchMedia('(max-width: 700px)').matches;
  const open = mobile ? document.body.classList.contains('mobile-nav-open') : !menuCollapsed;
  document.querySelectorAll('[data-action="menu"]').forEach(button => button.setAttribute('aria-expanded', String(open)));
  const sidebar = document.querySelector('#sidebar');
  sidebar.inert = mobile && !open;
  document.querySelector('.workspace').inert = mobile && open;
  if (!open && mobile) sidebar.setAttribute('aria-hidden', 'true');
  else sidebar.removeAttribute('aria-hidden');
}

function previewPost() {
  const form = document.querySelector('#create-form');
  const data = new FormData(form);
  postDraft = { title: data.get('title'), content: data.get('content'), categories: data.getAll('category') };
  document.querySelector('#post-preview-title').textContent = postDraft.title || 'Your post preview';
  document.querySelector('#post-preview').innerHTML = formatContent(postDraft.content, categoryById(postDraft.categories[0])?.language);
}

function rememberPlayground() {
  const form = document.querySelector('#playground-form');
  playgroundDraft = Object.fromEntries(new FormData(form));
  return playgroundDraft;
}

function previewCode() {
  const draft = rememberPlayground();
  document.querySelector('#code-preview').innerHTML = highlight(draft.code, draft.language);
}

function bindPage() {
  const filter = document.querySelector('#filter-form');
  filter?.addEventListener('submit', event => {
    event.preventDefault();
    const data = new FormData(filter);
    navigate(routeLink(parseRoute().path, { q: data.get('q').trim(), category: data.getAll('category'), sort: data.get('sort') }).slice(1));
  });
  const auth = document.querySelector('#auth-form');
  auth?.addEventListener('submit', event => {
    event.preventDefault();
    guarded(auth, () => {
      const values = Object.fromEntries(new FormData(auth));
      const user = auth.dataset.register === 'true' ? store.register(values) : store.login(values.username);
      navigate(safeNext(auth.dataset.next));
      if (parseRoute().path === safeNext(auth.dataset.next)) render();
      toast(`Welcome, ${user.displayName}. This is a simulated demo identity.`);
    });
  });
  document.querySelector('#profile-picker')?.addEventListener('change', event => {
    document.querySelector('#auth-username').value = event.target.value;
  });
  const comment = document.querySelector('#comment-form');
  comment?.addEventListener('submit', event => {
    event.preventDefault();
    guarded(comment, () => {
      const created = store.addComment(comment.dataset.postId, new FormData(comment).get('content'));
      render({ focus: false });
      document.getElementById(created.id).scrollIntoView({ block: 'center', behavior: 'smooth' });
      document.querySelector('#comment-content').focus({ preventScroll: true });
      toast('Comment added to this browser’s demo.');
    });
  });
  const create = document.querySelector('#create-form');
  if (create) {
    document.querySelector('#create-category-choices').innerHTML = categoryChoices(postDraft?.categories || []);
    if (postDraft) { create.elements.title.value = postDraft.title; create.elements.content.value = postDraft.content; previewPost(); }
    create.addEventListener('input', previewPost);
    create.addEventListener('submit', event => {
      event.preventDefault();
      guarded(create, () => {
        const data = new FormData(create);
        const post = store.createPost({ title: data.get('title'), content: data.get('content'), categories: data.getAll('category') });
        postDraft = null;
        navigate(`/post/${post.id}`);
      });
    });
  }
  const profile = document.querySelector('#profile-form');
  profile?.addEventListener('submit', event => {
    event.preventDefault();
    guarded(profile, () => {
      const values = Object.fromEntries(new FormData(profile));
      store.updateProfile({ ...values, interests: checked(profile, 'interest') });
      render({ focus: false });
      toast('Profile updated in this browser.');
    });
  });
  const playground = document.querySelector('#playground-form');
  if (playground) {
    if (playgroundDraft) {
      for (const [name, value] of Object.entries(playgroundDraft)) playground.elements[name].value = value;
    } else playground.elements.code.value = playgroundLanguages.find(language => language.id === 'javascript').sample;
    previewCode();
    playground.addEventListener('input', previewCode);
    playground.addEventListener('change', previewCode);
    playground.addEventListener('submit', event => { event.preventDefault(); previewCode(); toast('Syntax preview updated. Code was not executed.'); });
  }
}

app.addEventListener('click', event => {
  const actionElement = event.target.closest('[data-action]');
  if (!actionElement) {
    if (event.target.closest('a[href^="#/"]') && document.body.classList.contains('mobile-nav-open')) {
      document.body.classList.remove('mobile-nav-open'); applyMenuState();
    }
    return;
  }
  const action = actionElement.dataset.action;
  if (action === 'menu') {
    if (window.matchMedia('(max-width: 700px)').matches) document.body.classList.toggle('mobile-nav-open');
    else { menuCollapsed = !menuCollapsed; document.body.classList.toggle('nav-collapsed', menuCollapsed); }
    applyMenuState();
    if (document.body.classList.contains('mobile-nav-open')) document.querySelector('#sidebar .menu-button').focus();
    return;
  }
  if (action === 'menu-close') { document.body.classList.remove('mobile-nav-open'); applyMenuState(); document.querySelector('.mobile-menu').focus(); return; }
  if (action === 'categories') { categoriesExpanded = !categoriesExpanded; applyMenuState(); return; }
  if (action === 'reset') { document.querySelector('#reset-dialog').showModal(); return; }
  if (action === 'reset-cancel') { document.querySelector('#reset-dialog').close(); return; }
  if (action === 'reset-confirm') {
    store.reset(); postDraft = null; playgroundDraft = null;
    document.querySelector('#reset-dialog').close();
    if (parseRoute().path === '/') render(); else navigate('/');
    toast('The sample community has been restored.'); return;
  }
  if (action === 'demo') {
    store.login();
    const { path, params } = parseRoute();
    const next = actionElement.dataset.next ? safeNext(actionElement.dataset.next)
      : ['/login', '/register'].includes(path) ? safeNext(params.get('next')) : currentPath();
    if (currentPath() === next) render({ focus: false }); else navigate(next);
    toast('You are exploring as demo_dev. Authentication is simulated.'); return;
  }
  if (action === 'logout') {
    store.logout(); postDraft = null;
    if (parseRoute().path === '/') render(); else navigate('/');
    toast('Logged out of the simulated profile.'); return;
  }
  if (action === 'react') {
    if (!store.currentUser()) { navigate(routeLink('/login', { next: currentPath() }).slice(1)); return; }
    try {
      const { kind, id, reaction } = actionElement.dataset;
      store.react(kind, id, reaction);
      const entity = store.state[kind === 'post' ? 'posts' : 'comments'].find(item => item.id === id);
      document.querySelectorAll(`[data-reaction-group="${CSS.escape(`${kind}-${id}`)}"]`).forEach(group => { group.outerHTML = reactions(entity, kind, store.state.currentUserId); });
      const replacement = document.querySelector(`[data-reaction-group="${CSS.escape(`${kind}-${id}`)}"] [data-reaction="${reaction}"]`);
      replacement?.focus({ preventScroll: true });
    } catch (error) { toast(error.message); }
    return;
  }
  if (action === 'minimize') {
    const body = actionElement.closest('[data-card]').querySelector('.card-content');
    body.hidden = !body.hidden;
    actionElement.setAttribute('aria-expanded', String(!body.hidden));
    return;
  }
  if (action === 'dismiss') {
    actionElement.closest('[data-card]').hidden = true;
    document.querySelector('.restore-row').hidden = false;
    document.querySelector('[data-action="restore-cards"]').focus({ preventScroll: true });
    return;
  }
  if (action === 'restore-cards') {
    document.querySelectorAll('[data-card]').forEach(card => { card.hidden = false; });
    document.querySelector('.restore-row').hidden = true; return;
  }
  if (action === 'browse-category') { document.querySelector('#posts-section').scrollIntoView({ behavior: 'smooth' }); return; }
  if (action === 'insert-code') {
    const editor = document.querySelector('#post-content');
    const language = categoryById(checked(document.querySelector('#create-form'), 'category')[0])?.language || 'javascript';
    editor.setRangeText(`\n\n\`\`\`${language}\n// Your code here\n\`\`\`\n`, editor.selectionStart, editor.selectionEnd, 'end');
    previewPost(); editor.focus(); return;
  }
  if (action === 'clear-post') {
    document.querySelector('#create-form').reset(); postDraft = null; previewPost(); return;
  }
  if (action === 'reset-code' || action === 'clear-code') {
    const form = document.querySelector('#playground-form');
    form.elements.code.value = action === 'clear-code' ? '' : playgroundLanguages.find(language => language.id === form.elements.language.value).sample;
    previewCode(); return;
  }
  if (action === 'download') {
    const form = document.querySelector('#playground-form');
    guarded(form, () => {
      const draft = rememberPlayground();
      const lang = playgroundLanguages.find(language => language.id === draft.language);
      const bytes = downloadBytes(draft.code, draft.lineend, draft.bom === 'bom');
      const url = URL.createObjectURL(new Blob([bytes], { type: 'text/plain;charset=utf-8' }));
      const link = document.createElement('a');
      link.href = url; link.download = safeFilename(draft.filename, lang.ext);
      document.body.append(link); link.click(); link.remove();
      setTimeout(() => URL.revokeObjectURL(url), 1000);
      toast(`Downloaded ${safeFilename(draft.filename, lang.ext)}.`);
    });
  }
});

document.addEventListener('keydown', event => {
  if (event.key === 'Escape' && document.body.classList.contains('mobile-nav-open')) {
    document.body.classList.remove('mobile-nav-open'); applyMenuState(); document.querySelector('.mobile-menu').focus();
  }
});
document.querySelector('.skip-link').addEventListener('click', event => {
  event.preventDefault();
  const main = document.querySelector('#main-content');
  main.focus({ preventScroll: true });
  main.scrollIntoView({ block: 'start' });
});
window.addEventListener('hashchange', () => render());
window.addEventListener('resize', applyMenuState);
render({ focus: false });
