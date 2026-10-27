import { categories, categoryById } from '../data/categories.js';
import { selectPosts, paginate } from '../services/feed.js';
import { escapeHTML as h, highlight } from '../utils/content.js';
import { routeLink } from '../router.js';
import { pageHeading, postGrid, categoryChoices, asset, icon } from '../components/ui.js';

export function homePage(state) {
  return `${pageHeading('Welcome to Null Pointer', 'Discuss, share, and connect with others.', 'HELLO, WORLD. / WELCOME TO THE FORUM')}
  <div class="hero-actions"><a class="button primary" href="#/posts">Discover Posts →</a><a class="button" href="#/create-post">Start a discussion</a></div>
  <div class="community-stats"><span><strong>${state.users.length}</strong> fictional developers</span><span><strong>${state.posts.length}</strong> discussions</span><span><strong>8</strong> languages, endless ideas</span></div>
  <section class="home-categories" aria-label="Programming categories">${categories.map(category => `<a href="#/category/${category.id}" style="--category-color:${category.color}">${icon(`categories/${category.id}`)}<span>${category.name}</span></a>`).join('')}</section>
  <section class="feed-section"><div class="section-heading"><div><p class="eyebrow">WORTH A READ</p><h2>Featured Posts</h2></div><a href="#/posts?sort=popular">Explore popular →</a></div>${postGrid(selectPosts(state, { sort: 'popular' }).slice(0, 3), state)}</section>
  <section class="feed-section"><div class="section-heading"><div><p class="eyebrow">FRESH FROM THE COMMUNITY</p><h2>Latest Posts</h2></div><a href="#/posts">View all posts →</a></div>${postGrid(selectPosts(state).slice(0, 6), state)}</section>
  <section class="playground-callout panel"><div><p class="eyebrow">A LITTLE SPACE TO EXPERIMENT</p><h2>Have a snippet in mind?</h2><p>Write, preview, and download it in the Code Playground.</p></div><a class="button" href="#/playground">Open Playground ${icon('code')}</a></section>`;
}

export function feedPage(state, params, categoryId) {
  const category = categoryById(categoryId);
  const query = params.get('q') || '';
  const selected = category ? [category.id] : params.getAll('category').filter(id => categoryById(id));
  const sort = ['latest', 'popular', 'comments'].includes(params.get('sort')) ? params.get('sort') : 'latest';
  const result = paginate(selectPosts(state, { query, categories: selected, sort }), params.get('page') || 1);
  const path = category ? `/category/${category.id}` : '/posts';
  const link = page => routeLink(path, { q: query, category: category ? [] : selected, sort, page });
  const heading = category ? `<section class="lang-hero"><div class="lang-hero-content"><div class="lang-hero-text"><span class="hero-badge">${category.name} / DEVELOPER COMMUNITY</span><h1 class="lang-title"><span class="title-highlight">The</span><span class="title-main">${category.name} Zone</span></h1><p class="lang-subtitle">${category.description}</p><p class="lang-desc">From first questions to deeper discoveries. Learn ${category.name} with our programming community.</p><div class="lang-actions"><a class="button primary" href="${category.url}" target="_blank" rel="noopener noreferrer">Start learning now ↗</a><button class="button" data-action="browse-category">Browse posts ↓</button></div></div><div class="category-code-art"><div class="code-header"><span class="code-dots"><span></span><span></span><span></span></span><span>main.${category.ext}</span><img class="hero-language-icon" src="${asset(`icons/categories/${category.id}.svg`)}" alt="" width="36" height="36"></div><pre><code class="hljs">${highlight(category.sample, category.language)}</code></pre></div></div></section><nav class="category-switcher" aria-label="Switch category">${categories.map(c => `<a href="#/category/${c.id}" ${c.id === category.id ? 'aria-current="page"' : ''}>${c.name}</a>`).join('')}</nav>` : pageHeading('Discover Posts', 'Explore community knowledge, filter by category, and find what interests you.');
  return `${heading}<section class="panel filter-panel" aria-label="Post filters"><form id="filter-form"><div class="filter-main"><div class="field search-field"><label for="search-q">Search</label><input id="search-q" type="search" name="q" placeholder="Search titles and content…" maxlength="200" value="${h(query)}"></div><div class="field"><label for="filter-sort">Sort</label><select id="filter-sort" name="sort"><option value="latest" ${sort === 'latest' ? 'selected' : ''}>Latest</option><option value="popular" ${sort === 'popular' ? 'selected' : ''}>Most Popular</option><option value="comments" ${sort === 'comments' ? 'selected' : ''}>Most Commented</option></select></div><div class="filter-actions"><button type="submit" class="button primary">Apply filters</button><a class="button" href="#${path}">Reset filters</a></div></div>${category ? '' : `<fieldset><legend>Categories <span class="muted">· match any selected</span></legend>${categoryChoices(selected)}</fieldset>`}</form></section>
  <section id="posts-section" class="feed-section"><div class="section-heading"><h2>Results <span class="result-count">${result.total} discussions</span></h2><a href="#/create-post">+ Create Post</a></div>${postGrid(result.items, state)}</section>
  ${result.pages > 1 ? `<nav class="pagination" aria-label="Post pagination">${result.page > 1 ? `<a class="button" href="${link(result.page - 1)}">← Previous</a>` : ''}<span>Page ${result.page} of ${result.pages}</span>${result.page < result.pages ? `<a class="button" href="${link(result.page + 1)}">Next →</a>` : ''}</nav>` : ''}`;
}
