import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { renderTemplate, safeHTML } from '../src/original/template.js';
import { originalCategories, originalUser, originalPost, originalTheme, contentBlocks } from '../src/original/model.js';
import { categories } from '../src/data/categories.js';
import { createSeed } from '../src/data/seed.js';

const source = name => readFileSync(new URL(`../../web/templates/${name}.html`, import.meta.url), 'utf8');
test('original Go template conditions, ranges, variables, and methods render safely', () => {
  const result = renderTemplate('{{range $i, $entry := .Items}}{{if lt $i 1}}{{$entry.Name}}{{else if eq $entry.Name "B"}}second{{else}}other{{end}}{{end}}', { Items: [{ Name: '<A>' }, { Name: 'B' }, { Name: 'C' }] });
  assert.equal(result, '&lt;A&gt;secondother');
  assert.equal(renderTemplate('{{if or (not .Name) (eq .Name "go")}}yes{{else}}no{{end}}', { Name: '' }), 'yes');
  assert.equal(renderTemplate('{{range .Items}}{{.}}{{else}}empty{{end}}', { Items: [] }), 'empty');
  assert.equal(renderTemplate('{{.Date.Format "layout"}}', { Date: { Format: layout => `date-${layout}` } }), 'date-layout');
  assert.equal(renderTemplate('{{.User.constructor}}', { User: {} }), '');
});

test('all original application templates render from the original files', () => {
  const state = createSeed(); state.currentUserId = state.users[0].id;
  const post = originalPost(state.posts[0], state);
  const data = { User: originalUser(state.users[0]), IsLoggedIn: true, Categories: originalCategories,
    FeaturedPosts: [post], LatestPosts: [post], Posts: [post], Post: post, UserPosts: [post], LikedPosts: [post],
    SelectedCategories: ['golang'], TempBlocks: [{ Type: 'text', Content: 'A safe block' }],
    Slug: 'golang', DisplayName: 'Go', Theme: originalTheme(categories[0]), CountPosts: 1, SourceURL: categories[0].url,
    Code: '<script>neverRun()</script>', Language: 'go', LineEnd: 'lf', BOM: 'nobom', Sort: 'latest' };
  for (const name of ['home', 'DiscoverPosts', 'category', 'profile', 'edit-profile', 'login', 'startcoding', 'create-post', 'view-post', 'help-privacy', 'error']) {
    const result = renderTemplate(source(name), data);
    assert.doesNotMatch(result, /{{[\s\S]*?}}/, name);
    assert.ok(result.includes('<'), name);
  }
  assert.match(renderTemplate(source('startcoding'), data), /&lt;script&gt;neverRun\(\)&lt;\/script&gt;/);
  const root = renderTemplate(source('root'), { ...data, Body: safeHTML('<section>known template output</section>') });
  assert.ok(root.includes('<section>known template output</section>'));
});

test('original home retains its exact sections, window cards, and four-post collections', () => {
  const state = createSeed(), posts = state.posts.slice(0, 4).map(post => originalPost(post, state));
  const result = renderTemplate(source('home'), { FeaturedPosts: posts, LatestPosts: posts });
  assert.match(result, /id="hero-title" class="section-title">Welcome to Null Pointer/);
  assert.equal((result.match(/class="card"/g) || []).length, 8);
  assert.equal((result.match(/class="main-section"/g) || []).length, 3);
  assert.equal((result.match(/class="minimize"/g) || []).length, 8);
  assert.doesNotMatch(result, /community-stats|home-categories|topbar|playground-callout/);
  assert.equal(contentBlocks('https://go.dev/learn/')[0].Type, 'link');
  assert.equal(contentBlocks('javascript:alert(1)')[0].Type, 'text');
});
