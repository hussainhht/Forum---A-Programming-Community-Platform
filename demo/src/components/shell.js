import { categories } from '../data/categories.js';
import { escapeHTML as h } from '../utils/content.js';
import { icon, avatar } from './ui.js';

export function shell(state, path, view, storageWarning) {
  const user = state.users.find(user => user.id === state.currentUserId);
  const nav = (href, label, image) => `<a class="sidebar-button ${path === href ? 'selected' : ''}" href="#${href}" ${path === href ? 'aria-current="page"' : ''} title="${label}">${icon(image, 'sidebar-button-icon')}<span class="nav-text">${label}</span></a>`;
  return `<button class="sidebar-scrim" data-action="menu-close" aria-label="Close navigation" tabindex="-1"></button>
  <aside id="sidebar" aria-label="Forum sidebar">
    <div class="sidebar-brand"><button class="menu-button" data-action="menu" aria-label="Toggle navigation" aria-expanded="true" aria-controls="sidebar-navigation">${icon('menu')}</button><a class="nav-text" href="#/">Null Pointer<span>PROGRAMMING COMMUNITY</span></a></div>
    <nav id="sidebar-navigation" aria-label="Main navigation">
      <div class="nav-group">${nav('/', 'Home', 'home')}${nav('/create-post', 'Create Post', 'add_post')}${nav('/posts', 'Discover Posts', 'read_more')}${nav('/playground', 'Playground', 'code')}</div>
      <div class="nav-group"><p class="nav-caption nav-text">LANGUAGES</p>${categories.map((category, index) => `<a class="sidebar-button ${index > 2 ? 'extra-category' : ''} ${path === `/category/${category.id}` ? 'selected' : ''}" href="#/category/${category.id}" ${path === `/category/${category.id}` ? 'aria-current="page"' : ''} title="${category.name}">${icon(`categories/${category.id}`, 'category-icon')}<span class="nav-text">${category.name}</span></a>`).join('')}
      <button class="sidebar-button more-categories" data-action="categories" aria-expanded="false">${icon('arrow_down', 'sidebar-button-icon')}<span class="nav-text">Show More</span></button></div>
      <div class="nav-group">${nav('/edit-profile', 'Settings', 'settings')}${nav('/help', 'Help', 'help')}${nav('/privacy-terms', 'Privacy & Terms', 'info')}
      <button class="sidebar-button reset-button" data-action="reset" title="Reset Demo Data">${icon('settings', 'sidebar-button-icon')}<span class="nav-text">Reset Demo Data</span></button></div>
    </nav>
    <div class="sidebar-account">${user ? `<a class="sidebar-user" href="#/profile">${avatar(user)}<span class="nav-text"><strong>${h(user.username)}</strong><small>Demo profile</small></span></a><button class="sidebar-button" data-action="logout" title="Logout">${icon('logout', 'sidebar-button-icon')}<span class="nav-text">Logout</span></button>` : `${nav('/login', 'Login', 'login')}${nav('/register', 'Register', 'register')}`}</div>
  </aside>
  <div class="workspace"><header class="topbar"><button class="menu-button mobile-menu" data-action="menu" aria-label="Open navigation" aria-expanded="false" aria-controls="sidebar-navigation">${icon('menu')}</button><a class="topbar-brand" href="#/">&lt;/&gt; Forum</a><span class="topbar-context">A space to think in code.</span><div class="topbar-actions"><a href="https://github.com/hussainhht/Forum---A-Programming-Community-Platform" target="_blank" rel="noopener noreferrer">GitHub ↗</a>${user ? `<a href="#/profile">@${h(user.username)}</a>` : `<button class="button small" data-action="demo">Continue as Demo User</button>`}</div></header>
    <div class="demo-notice"><span class="status-dot" aria-hidden="true"></span><span><strong>Interactive Demo</strong> — All content and interactions are simulated and saved locally in your browser.</span></div>
    <p id="storage-warning" class="storage-warning" role="status" ${storageWarning ? '' : 'hidden'}>${h(storageWarning)}</p>
    <div class="restore-row" hidden><span>Some cards are hidden in this view.</span><button class="text-button" data-action="restore-cards">Restore cards</button></div>
    <main id="main-content" tabindex="-1">${view}</main>
    <footer class="site-footer"><span>Null Pointer · Built to share what you learn.</span><a href="#/privacy-terms">Local demo privacy</a><a href="https://github.com/hussainhht/Forum---A-Programming-Community-Platform" target="_blank" rel="noopener noreferrer">Original Go project ↗</a></footer>
  </div>
  <div id="toast" role="status" class="toast" hidden></div>
  <dialog id="reset-dialog" aria-labelledby="reset-title"><h2 id="reset-title">Reset Demo Data?</h2><p>This removes posts, comments, reactions, profile edits, and demo accounts you created in this browser. The original sample community will be restored and you will be logged out.</p><p>Other websites and browser data are unaffected.</p><div class="actions"><button class="button" data-action="reset-cancel" autofocus>Cancel</button><button class="button danger" data-action="reset-confirm">Reset Demo Data</button></div></dialog>`;
}
