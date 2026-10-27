import { test } from 'node:test';
import assert from 'node:assert/strict';
import { createSeed } from '../src/data/seed.js';
import { categories, avatarColors } from '../src/data/categories.js';
import { createStore, validateState, STORAGE_KEY } from '../src/services/storage.js';
import { selectPosts, paginate, reactionCount } from '../src/services/feed.js';
import { escapeHTML, formatContent, downloadBytes, safeFilename } from '../src/utils/content.js';
import { parseRoute, routeLink } from '../src/router.js';

const memory = (initial = {}) => {
  const data = new Map(Object.entries(initial));
  return { data, getItem: key => data.get(key) ?? null, setItem: (key, value) => data.set(key, value) };
};

test('seed is deterministic, populated, and referentially valid', () => {
  const seed = createSeed();
  assert.deepEqual(seed, createSeed());
  assert.equal(seed.users.length, 18); assert.equal(seed.posts.length, 40); assert.equal(seed.comments.length, 100);
  assert.ok(validateState(seed));
  for (const category of categories) assert.ok(seed.posts.filter(post => post.categories.includes(category.id)).length >= 5);
  assert.equal(new Set(seed.posts.map(post => post.content)).size, 40);
});

test('data shape and foreign keys are validated before restoration', () => {
  for (const mutate of [
    state => { state.posts[0].authorId = 'missing'; },
    state => { state.comments[0].postId = 'missing'; },
    state => { state.posts[0].categories = ['unknown']; },
    state => { state.users[0].username = null; },
    state => { state.users[0].avatar = 'url(javascript:bad)'; },
    state => { state.currentUserId = 'missing'; },
    state => { state.posts[0].reactions = { unknown: 'like' }; },
    state => { state.users[1].id = state.users[0].id; },
    state => { state.users[1].username = state.users[0].username.toUpperCase(); },
  ]) { const state = createSeed(); mutate(state); assert.equal(validateState(state), false); }
});

test('first visit seeds once and refresh preserves entity IDs and edits', () => {
  const storage = memory(); const store = createStore(storage);
  assert.equal(storage.data.size, 1); store.login();
  const post = store.createPost({ title: 'A new Go question', content: 'How can I test this channel?', categories: ['golang'] });
  const comment = store.addComment(post.id, 'Use a deterministic synchronization point.');
  const restored = createStore(storage);
  assert.equal(restored.state.posts.length, 41); assert.equal(restored.state.comments.length, 101);
  assert.ok(restored.state.posts.some(item => item.id === post.id));
  assert.ok(restored.state.comments.some(item => item.id === comment.id));
  assert.equal(restored.currentUser().username, 'demo_dev');
  assert.ok(validateState(restored.state));
});

test('reactions toggle off, switch type, and remain isolated per identity', () => {
  const storage = memory(); const store = createStore(storage); store.login();
  for (const [kind, entity] of [['post', store.state.posts[0]], ['comment', store.state.comments[0]]]) {
    const initialLikes = reactionCount(entity, 'like'); const initialDislikes = reactionCount(entity, 'dislike');
    store.react(kind, entity.id, 'like'); assert.equal(reactionCount(entity, 'like'), initialLikes + 1);
    store.react(kind, entity.id, 'like'); assert.equal(reactionCount(entity, 'like'), initialLikes);
    store.react(kind, entity.id, 'dislike'); store.react(kind, entity.id, 'like');
    assert.equal(reactionCount(entity, 'dislike'), initialDislikes); assert.equal(reactionCount(entity, 'like'), initialLikes + 1);
  }
  store.logout(); assert.throws(() => store.react('post', '1', 'like'));
  store.login('channel_chaser'); assert.notEqual(store.currentUser().id, 'user-1');
  assert.equal(createStore(storage).state.posts[0].reactions['user-1'], 'like');
});

test('registration, profile edits, and logout persist without passwords', () => {
  const storage = memory(); const store = createStore(storage);
  const user = store.register({ username: 'test_developer', displayName: 'Test Dev' });
  assert.throws(() => store.register({ username: 'TEST_DEVELOPER', displayName: 'Duplicate' }));
  assert.throws(() => store.register({ username: '<img>', displayName: 'Invalid' }));
  store.updateProfile({ displayName: 'Edited Name', bio: 'Enjoying matrices.', interests: ['matlab'], avatar: avatarColors[2] });
  assert.equal(createStore(storage).currentUser().displayName, 'Edited Name');
  assert.ok(!Object.hasOwn(user, 'password')); assert.ok(!Object.hasOwn(user, 'email'));
  store.logout(); assert.equal(createStore(storage).currentUser(), undefined);
});

test('local profile photos are bounded, non-executable, and persisted per profile', () => {
  const storage = memory(), store = createStore(storage);
  assert.throws(() => store.setPhoto('data:image/jpeg;base64,YQ=='));
  store.login(); store.setPhoto('data:image/jpeg;base64,YQ==');
  assert.equal(createStore(storage).currentUser().photo, 'data:image/jpeg;base64,YQ==');
  assert.throws(() => store.setPhoto('data:image/svg+xml,<svg onload=alert(1)>'));
  assert.throws(() => store.setPhoto('data:image/jpeg;base64,' + 'A'.repeat(50000)));
  assert.ok(validateState(store.state));
});

test('validation blocks empty content, invalid categories, and guest writes', () => {
  const store = createStore(memory());
  assert.throws(() => store.createPost({ title: 'Valid title', content: 'Some useful content', categories: ['golang'] }));
  store.login();
  for (const values of [
    { title: 'bad', content: 'Useful description', categories: ['golang'] },
    { title: 'Good title', content: 'short', categories: ['golang'] },
    { title: 'Good title', content: 'Useful description', categories: [] },
    { title: 'Good title', content: 'Useful description', categories: ['invalid'] },
  ]) assert.throws(() => store.createPost(values));
  assert.throws(() => store.addComment('missing', 'Useful reply'));
  assert.throws(() => store.addComment('1', '   '));
});

test('corrupted storage and storage access failures recover gracefully', () => {
  for (const saved of ['{invalid', 'null', '{"version":1}', JSON.stringify({ ...createSeed(), users: [null] })]) {
    const messages = []; const store = createStore(memory({ [STORAGE_KEY]: saved }), message => messages.push(message));
    assert.equal(store.state.posts.length, 40); assert.match(messages.at(-1), /restored/);
  }
  const messages = [];
  const broken = { getItem: () => null, setItem: () => { throw new Error('QuotaExceededError'); } };
  const store = createStore(broken, message => messages.push(message)); store.login();
  store.addComment('1', 'This works in memory.');
  assert.equal(store.state.comments.length, 101); assert.match(messages.at(-1), /may be lost/);
  assert.equal(createStore(null).state.posts.length, 40);
});

test('reset restores only demo state and retains unrelated keys', () => {
  const storage = memory({ 'unrelated-setting': 'keep-me' }); const store = createStore(storage); store.login(); store.addComment('1', 'Temporary reply');
  store.reset(); assert.deepEqual(store.state, createSeed()); assert.equal(storage.data.get('unrelated-setting'), 'keep-me');
  assert.deepEqual(createStore(storage).state, createSeed());
});

test('search, multi-category union, sorting, and pagination operate on actual data', () => {
  const state = createSeed();
  assert.equal(selectPosts(state, { query: 'gOrOuTiNeS' })[0].id, '1');
  assert.ok(selectPosts(state, { query: 'WaitGroup' }).some(post => post.id === '1') === false); // comments are not searched
  assert.ok(selectPosts(state, { query: 'Hash.new' }).some(post => post.id === '17'));
  const results = selectPosts(state, { categories: ['java', 'ruby'] }); assert.equal(results.length, 10);
  const popular = selectPosts(state, { sort: 'popular' }); assert.ok(reactionCount(popular[0], 'like') >= reactionCount(popular.at(-1), 'like'));
  const comments = selectPosts(state, { sort: 'comments' }); assert.equal(state.comments.filter(comment => comment.postId === comments[0].id).length, 3);
  assert.equal(paginate(state.posts, 100).page, 5); assert.equal(paginate(state.posts, -7).page, 1);
  assert.equal(paginate(state.posts, 'invalid').page, 1); assert.equal(paginate([], 2).pages, 1);
  assert.equal(paginate(state.posts, 2).items.length, 9);
});

test('formatting escapes HTML and never creates user-controlled elements', () => {
  const malicious = '<img src=x onerror=alert(1)><script>alert(1)</script>';
  assert.ok(escapeHTML(malicious).includes('&lt;img'));
  for (const result of [formatContent(malicious), formatContent(`\`\`\`html\n${malicious}\n\`\`\``)]) {
    assert.doesNotMatch(result, /<img|<script|<iframe/);
  }
  assert.match(formatContent('```go\nfunc main() {}\n```'), /hljs-keyword/);
  assert.match(formatContent('```bad-language\n<img>\n```'), /&lt;img&gt;/);
});

test('downloads produce exact UTF-8 bytes, line endings, BOM, and safe filenames', () => {
  assert.deepEqual([...downloadBytes('a\r\nb\rc\n', 'crlf', true)], [239,187,191,97,13,10,98,13,10,99,13,10]);
  assert.equal(new TextDecoder().decode(downloadBytes('λ\r\n雪\r', 'lf')), 'λ\n雪\n');
  assert.equal(safeFilename('../../secret', 'go'), '.._.._secret.go');
  assert.equal(safeFilename('CON', 'txt'), '_CON.txt');
  assert.equal(safeFilename('hello.js', 'js'), 'hello.js');
  assert.equal(safeFilename(' ', 'py'), 'solution.py');
});

test('hash links preserve search and multi-category filters on direct navigation', () => {
  const hash = routeLink('/posts', { q: 'x & y?', category: ['golang', 'java'], page: 2 });
  const route = parseRoute(hash); assert.equal(route.path, '/posts'); assert.equal(route.params.get('q'), 'x & y?');
  assert.deepEqual(route.params.getAll('category'), ['golang', 'java']); assert.equal(route.params.get('page'), '2');
  assert.equal(parseRoute('').path, '/');
});
