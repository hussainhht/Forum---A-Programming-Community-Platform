import { createSeed } from '../data/seed.js';
import { categories, avatarColors } from '../data/categories.js';

export const STORAGE_KEY = 'null-pointer:demo:v1';
const categoryIds = new Set(categories.map(category => category.id));
const text = (value, max) => typeof value === 'string' && value.length <= max;
const date = value => text(value, 40) && Number.isFinite(Date.parse(value));
const id = value => text(value, 100) && /^[\w-]+$/.test(value);
const validCategories = values => Array.isArray(values) && values.length <= 8 && values.every(value => categoryIds.has(value)) && new Set(values).size === values.length;
const uniqueIds = values => values.every(value => value && id(value.id)) && new Set(values.map(value => value.id)).size === values.length;

export function validateState(state) {
  if (!state || state.version !== 1 || !Array.isArray(state.users) || !Array.isArray(state.posts) || !Array.isArray(state.comments)) return false;
  if (!state.users.length || state.users.length > 1000 || state.posts.length > 5000 || state.comments.length > 15000) return false;
  if (![state.users, state.posts, state.comments].every(uniqueIds)) return false;
  const users = new Set(state.users.map(user => user.id));
  const posts = new Set(state.posts.map(post => post.id));
  const reactions = value => value && typeof value === 'object' && !Array.isArray(value)
    && Object.entries(value).every(([key, reaction]) => users.has(key) && ['like', 'dislike'].includes(reaction));
  return (state.currentUserId === null || users.has(state.currentUserId))
    && state.users.every(user => text(user.username, 30) && /^[a-zA-Z0-9_]{3,30}$/.test(user.username) && text(user.displayName, 60)
      && text(user.bio, 500) && validCategories(user.interests) && avatarColors.includes(user.avatar) && date(user.createdAt)
      && (user.photo === undefined || text(user.photo, 50000) && /^data:image\/jpeg;base64,[A-Za-z0-9+/=]+$/.test(user.photo)))
    && new Set(state.users.map(user => user.username.toLowerCase())).size === state.users.length
    && state.posts.every(post => users.has(post.authorId) && text(post.title, 160) && post.title.trim().length >= 5
      && text(post.content, 30000) && post.content.trim().length >= 10 && validCategories(post.categories) && post.categories.length > 0 && date(post.createdAt) && reactions(post.reactions))
    && state.comments.every(comment => users.has(comment.authorId) && posts.has(comment.postId)
      && text(comment.content, 10000) && comment.content.trim().length > 0 && date(comment.createdAt) && reactions(comment.reactions));
}

export function createStore(storage, onStatus = () => {}) {
  let state;
  function save() {
    try {
      storage.setItem(STORAGE_KEY, JSON.stringify(state));
      onStatus('');
    } catch {
      onStatus('Browser storage is unavailable or full. Your changes work in this tab but may be lost on refresh.');
    }
  }
  try {
    const saved = storage.getItem(STORAGE_KEY);
    if (saved) {
      state = JSON.parse(saved);
      if (!validateState(state)) throw new Error('Invalid demo data');
    } else {
      state = createSeed();
      save();
    }
  } catch {
    state = createSeed();
    onStatus('Stored demo data could not be loaded. The original sample community has been restored.');
    // Do not immediately overwrite unreadable data. The next successful
    // interaction or explicit reset writes a complete, valid snapshot.
  }
  const currentUser = () => state.users.find(user => user.id === state.currentUserId);
  const requireUser = () => {
    const user = currentUser();
    if (!user) throw new Error('Continue as Demo User or choose a fictional profile first.');
    return user;
  };
  const now = () => new Date().toISOString();
  const makeId = prefix => `${prefix}-${globalThis.crypto.randomUUID()}`;
  return {
    get state() { return state; },
    currentUser,
    login(username = 'demo_dev') {
      const user = state.users.find(user => user.username.toLowerCase() === username.trim().toLowerCase());
      if (!user) throw new Error('That fictional username was not found. Choose a profile below or register one.');
      state.currentUserId = user.id;
      save();
      return user;
    },
    logout() { state.currentUserId = null; save(); },
    setPhoto(photo) {
      const user = requireUser();
      if (!text(photo, 50000) || !/^data:image\/jpeg;base64,[A-Za-z0-9+/=]+$/.test(photo)) throw new Error('Choose a valid, small profile photo.');
      user.photo = photo; save();
    },
    register({ username, displayName }) {
      username = username.trim(); displayName = displayName.trim();
      if (!/^[a-zA-Z0-9_]{3,30}$/.test(username)) throw new Error('Use 3–30 letters, numbers, or underscores for your username.');
      if (!displayName || displayName.length > 60) throw new Error('Enter a fictional display name of up to 60 characters.');
      if (state.users.some(user => user.username.toLowerCase() === username.toLowerCase())) throw new Error('That username is already in use.');
      const user = { id: makeId('user'), username, displayName, bio: 'New to the Null Pointer community.', interests: [], avatar: avatarColors[0], createdAt: now() };
      state.users.push(user); state.currentUserId = user.id; save(); return user;
    },
    updateProfile({ displayName, bio, interests, avatar }) {
      const user = requireUser();
      if (!displayName.trim() || displayName.trim().length > 60 || bio.length > 500 || !validCategories(interests) || !avatarColors.includes(avatar)) throw new Error('Check your profile name, biography, interests, and avatar color.');
      Object.assign(user, { displayName: displayName.trim(), bio: bio.trim(), interests, avatar }); save();
    },
    createPost({ title, content, categories: selected }) {
      const user = requireUser(); title = title.trim(); content = content.trim();
      if (title.length < 5 || title.length > 160) throw new Error('Write a title between 5 and 160 characters.');
      if (content.length < 10 || content.length > 30000) throw new Error('Write between 10 and 30,000 characters of content.');
      if (!validCategories(selected) || !selected.length) throw new Error('Select at least one programming category.');
      const post = { id: makeId('post'), authorId: user.id, title, content, categories: selected, createdAt: now(), reactions: {} };
      state.posts.unshift(post); save(); return post;
    },
    addComment(postId, content) {
      const user = requireUser(); content = content.trim();
      if (!state.posts.some(post => post.id === postId)) throw new Error('This post could not be found.');
      if (!content || content.length > 10000) throw new Error('Write a comment between 1 and 10,000 characters.');
      const comment = { id: makeId('comment'), postId, authorId: user.id, content, createdAt: now(), reactions: {} };
      state.comments.push(comment); save(); return comment;
    },
    react(kind, entityId, reaction) {
      const user = requireUser();
      if (!['post', 'comment'].includes(kind) || !['like', 'dislike'].includes(reaction)) throw new Error('Unknown reaction.');
      const entity = state[kind === 'post' ? 'posts' : 'comments'].find(item => item.id === entityId);
      if (!entity) throw new Error('This item could not be found.');
      if (entity.reactions[user.id] === reaction) delete entity.reactions[user.id];
      else entity.reactions[user.id] = reaction;
      save();
    },
    reset() { state = createSeed(); save(); },
  };
}
