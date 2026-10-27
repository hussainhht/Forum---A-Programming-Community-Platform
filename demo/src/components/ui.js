import { categories, categoryById } from '../data/categories.js';
import { escapeHTML as h, formatContent, plainPreview } from '../utils/content.js';
import { reactionCount, commentCount } from '../services/feed.js';

export const asset = path => `${import.meta.env.BASE_URL}assets/${path}`;
export const dateLabel = value => new Date(value).toLocaleDateString('en', { month: 'short', day: 'numeric', year: 'numeric' });
export const icon = (name, className = '') => `<img class="${className}" src="${asset(`icons/${name}.svg`)}" alt="" width="22" height="22">`;

export function avatar(user, large = false) {
  const initials = user.displayName.split(/\s+/).filter(Boolean).slice(0, 2).map(part => part[0]).join('').toUpperCase();
  return `<span class="avatar ${large ? 'avatar-large' : ''}" style="--avatar-color:${user.avatar}" role="img" aria-label="${h(user.displayName)} avatar">${h(initials)}</span>`;
}

export function author(user, createdAt) {
  return `<a class="author" href="#/profile/${h(user.id)}">${avatar(user)}<span><strong>${h(user.username)}</strong><time datetime="${h(createdAt)}">${dateLabel(createdAt)}</time></span></a>`;
}

export function badges(ids) {
  return `<div class="badges">${ids.map(id => {
    const category = categoryById(id);
    return `<a class="badge" href="#/category/${id}" style="--category-color:${category.color}">${icon(`categories/${id}`)}${category.name}</a>`;
  }).join('')}</div>`;
}

export function categoryChoices(selected = [], name = 'category') {
  return `<div class="category-choices">${categories.map(category => `<label class="category-choice"><input type="checkbox" name="${name}" value="${category.id}" ${selected.includes(category.id) ? 'checked' : ''}>${icon(`categories/${category.id}`)}<span>${category.name}</span></label>`).join('')}</div>`;
}

export function reactions(item, kind, userId) {
  return `<div class="reactions" data-reaction-group="${kind}-${h(item.id)}">${['like', 'dislike'].map(type => `<button class="reaction ${item.reactions[userId] === type ? 'active' : ''}" data-action="react" data-kind="${kind}" data-id="${h(item.id)}" data-reaction="${type}" aria-pressed="${item.reactions[userId] === type}" aria-label="${type === 'like' ? 'Like' : 'Dislike'} ${kind}"><span aria-hidden="true">${type === 'like' ? '↑' : '↓'}</span><span class="reaction-count">${reactionCount(item, type)}</span></button>`).join('')}</div>`;
}

export function postCard(post, state) {
  const user = state.users.find(user => user.id === post.authorId);
  return `<article class="forum-card" data-card="${h(post.id)}">
    <div class="titlebar"><span class="window-label">discussion.${categoryById(post.categories[0]).ext}</span><span class="window-controls">
      <button type="button" data-action="minimize" aria-label="Minimize ${h(post.title)}" aria-expanded="true">−</button>
      <a href="#/post/${h(post.id)}" aria-label="Open ${h(post.title)}">□</a>
      <button type="button" data-action="dismiss" aria-label="Hide ${h(post.title)}">×</button>
    </span></div>
    <div class="card-content"><div class="card-meta">${author(user, post.createdAt)}</div>
      <h3><a href="#/post/${h(post.id)}">${h(post.title)}</a></h3>
      ${badges(post.categories)}<p class="card-preview">${h(plainPreview(post.content))}</p>
      <div class="card-bottom">${reactions(post, 'post', state.currentUserId)}<a class="comment-link" href="#/post/${h(post.id)}">${commentCount(state, post.id)} comments →</a></div>
    </div></article>`;
}

export const postGrid = (posts, state) => posts.length ? `<div class="post-grid">${posts.map(post => postCard(post, state)).join('')}</div>` : `<div class="empty-state"><span aria-hidden="true">{ }</span><h3>No discussions yet</h3><p>Try another search or start a conversation.</p><a class="button" href="#/create-post">Create Post</a></div>`;
export const pageHeading = (title, subtitle, eyebrow = 'NULL POINTER / COMMUNITY') => `<header class="page-heading"><p class="eyebrow">${h(eyebrow)}</p><h1>${h(title)}</h1><p class="subtitle">${h(subtitle)}</p></header>`;
export const formError = () => '<p class="form-error" role="alert" hidden></p>';
export const content = (text, language) => `<div class="formatted-content">${formatContent(text, language)}</div>`;

export function authGate(next, title = 'Join the conversation') {
  return `${pageHeading(title, 'Choose a fictional profile to post, comment, and react.')}<section class="panel auth-gate"><h2>Your next idea belongs here.</h2><p>No real account is needed. Demo identity is simulated and stored in this browser.</p><button class="button primary" data-action="demo" data-next="${h(next)}">Continue as Demo User</button><a class="button" href="#/login?next=${encodeURIComponent(next)}">Choose another profile</a></section>`;
}
