import { test, expect } from '@playwright/test';
import { readFile } from 'node:fs/promises';

const key = 'null-pointer:demo:v1';
const issues = new WeakMap();
test.beforeEach(async ({ page }) => {
  const errors = []; issues.set(page, errors);
  page.on('pageerror', error => errors.push(error.message));
  page.on('console', message => { if (['error', 'warning'].includes(message.type())) errors.push(message.text()); });
  page.on('response', response => { if (response.status() >= 400) errors.push(`${response.status()} ${response.url()}`); });
});
test.afterEach(async ({ page }) => { expect(issues.get(page)).toEqual([]); });
async function login(page) {
  await page.goto('./#/login');
  await page.locator('#login-form').getByRole('button', { name: 'Continue as Demo User', exact: true }).click();
  await expect(page.locator('#hero-title')).toBeVisible();
}
async function navigation(page, route) {
  await page.locator(`#sidebar [data-route="${route}"]`).click();
}

test('original home identity, collapsed sidebar, original card controls and navigation', async ({ page }, info) => {
  await page.goto('./#/');
  await expect(page.locator('#hero-title')).toHaveText('Welcome to Null Pointer');
  await expect(page.locator('#sidebar-toggle')).not.toBeChecked();
  await expect(page.locator('#featured-section .card')).toHaveCount(4);
  await expect(page.locator('#latest-section .card')).toHaveCount(4);
  await expect(page.locator('.topbar,.community-stats,.home-categories,.playground-callout')).toHaveCount(0);
  const original = await page.locator('#hero-title').evaluate(element => ({ fontSize: getComputedStyle(element).fontSize, align: getComputedStyle(element).textAlign }));
  expect(original).toEqual({ fontSize: '48px', align: 'center' });
  const card = page.locator('#featured-section .card').first();
  expect(Math.round((await card.boundingBox()).width)).toBe(360);
  expect(Math.round((await card.boundingBox()).height)).toBe(300);
  await page.screenshot({ path: info.outputPath('original-home.png'), animations: 'disabled' });
  await card.getByRole('button', { name: 'Minimize card' }).click();
  await expect(card.locator('.post-card-body')).toBeHidden();
  await card.getByRole('button', { name: 'Minimize card' }).click();
  await card.getByRole('button', { name: 'Hide card' }).click(); await expect(card).toBeHidden();
  await page.locator('.demo-tools summary').click(); await page.getByRole('button', { name: 'Restore hidden cards' }).click();
  await expect(card).toBeVisible();
  await page.getByRole('button', { name: 'Toggle navigation', exact: true }).click();
  await expect(page.locator('#sidebar-toggle')).toBeChecked();
  await navigation(page, '/posts');
  await expect(page.getByRole('heading', { name: 'Discover Posts', exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'Show More', exact: true }).click();
  await navigation(page, '/category/ruby');
  await expect(page.locator('.title-main')).toHaveText('Ruby Zone');
  await expect(page.locator('.feature-card')).toHaveCount(4);
  await expect(page.locator('.float-item')).toHaveCount(4);
});

test('all original themed categories, help/privacy and error templates', async ({ page }, info) => {
  for (const [slug, name] of [['golang', 'Go'], ['javascript', 'JavaScript'], ['rust', 'Rust'], ['ruby', 'Ruby'], ['python', 'Python'], ['matlab', 'MATLAB'], ['brainfuck', 'Brainfuck'], ['java', 'Java']]) {
    await page.goto(`./#/category/${slug}`);
    await expect(page.locator('.title-main')).toHaveText(`${name} Zone`);
    await expect(page.locator('.feature-card')).toHaveCount(4);
    expect(await page.locator('.post-card').count()).toBeGreaterThanOrEqual(5);
    await page.getByRole('link', { name: /Browse posts/ }).click();
    await expect(page).toHaveURL(new RegExp(`#/category/${slug}$`));
    if (slug === 'golang') { await page.locator('#main-content').evaluate(element => element.scrollTo(0, 0)); await page.screenshot({ path: info.outputPath('original-go.png'), animations: 'disabled' }); }
  }
  await page.goto('./#/help'); await expect(page.locator('#help-panel h1')).toHaveText('Help & Support');
  await page.goto('./#/privacy-terms'); await expect(page.locator('#privacy-panel h1')).toBeVisible();
  await expect(page.locator('#privacy-panel')).toContainText('browser localStorage');
  for (const route of ['missing', 'post/missing', 'profile/missing', 'category/missing']) {
    await page.goto(`./#/${route}`); await expect(page.getByRole('heading', { name: 'Page not found' })).toBeVisible();
  }
});

test('original discover search, dropdown categories, sorting, pagination and history', async ({ page }, info) => {
  await page.goto('./#/posts');
  await expect(page.locator('.post-card')).toHaveCount(9);
  await page.getByLabel('Search', { exact: true }).fill('goroutines');
  await page.getByRole('button', { name: 'Apply', exact: true }).click();
  await expect(page.locator('.post-card')).toHaveCount(1); await expect(page.locator('.post-title')).toContainText('goroutines');
  await page.getByRole('link', { name: 'Reset filters' }).click();
  await page.getByRole('button', { name: 'Select Categories' }).click();
  await page.locator('#posts-filter-form .category-option:has(input[value="ruby"])').click();
  await page.locator('#posts-filter-form .category-option:has(input[value="java"])').click();
  await expect(page.getByRole('checkbox', { name: 'ruby', exact: true })).toBeChecked();
  await page.getByRole('button', { name: 'Select Categories' }).click();
  await page.getByLabel('Sort', { exact: true }).selectOption('comments');
  await page.getByRole('button', { name: 'Apply', exact: true }).click();
  await expect(page.locator('.demo-results-count')).toContainText('10 posts');
  await page.getByRole('link', { name: 'Next', exact: true }).click();
  await expect(page.locator('.post-card')).toHaveCount(1); await expect(page.locator('.demo-results-count')).toContainText('Page 2 of 2');
  await page.reload(); await expect(page.locator('.demo-results-count')).toContainText('Page 2 of 2');
  await page.goBack(); await expect(page.locator('.demo-results-count')).toContainText('Page 1 of 2');
  await page.getByRole('link', { name: 'Reset filters' }).click();
  await page.getByLabel('Sort', { exact: true }).selectOption('popular'); await page.getByRole('button', { name: 'Apply', exact: true }).click();
  await page.screenshot({ path: info.outputPath('original-discover.png'), animations: 'disabled' });
  await page.getByLabel('Search', { exact: true }).fill('nomatch123456'); await page.getByRole('button', { name: 'Apply', exact: true }).click();
  await expect(page.locator('#posts-section')).toContainText('No posts found.');
});

test('original post and comments cards, reactions, comment persistence and guest mode', async ({ page }, info) => {
  await login(page); await page.goto('./#/post/1');
  await expect(page.locator('.vp-post')).toBeVisible(); await expect(page.locator('.vp-comments')).toBeVisible();
  await expect(page.locator('#post-content .code-string')).toContainText('make(chan int)');
  const like = page.getByRole('button', { name: 'Like post', exact: true }), dislike = page.getByRole('button', { name: 'Dislike post', exact: true });
  await like.click(); await expect(like).toHaveAttribute('aria-pressed', 'true');
  await dislike.click(); await expect(like).toHaveAttribute('aria-pressed', 'false'); await expect(dislike).toHaveAttribute('aria-pressed', 'true');
  await dislike.click(); await expect(dislike).toHaveAttribute('aria-pressed', 'false'); await like.click();
  const commentLike = page.getByRole('button', { name: 'Like comment', exact: true }).first(); await commentLike.click();
  await page.getByLabel('Your comment').fill('The producer should close the channel after the last send.\n\n```go\nclose(values)\n```');
  await page.getByRole('button', { name: 'Comment', exact: true }).click(); await expect(page.locator('.vp-comment')).toHaveCount(4);
  await page.reload(); await expect(page.locator('.vp-comment')).toHaveCount(4); await expect(like).toHaveAttribute('aria-pressed', 'true');
  await expect(commentLike).toHaveAttribute('aria-pressed', 'true');
  await page.screenshot({ path: info.outputPath('original-post.png'), animations: 'disabled' });
  await page.goto('./#/logout'); await expect(page.locator('#hero-title')).toBeVisible();
  await page.goto('./#/post/1'); await expect(page.getByRole('button', { name: 'Continue as Demo User' }).first()).toBeVisible();
});

test('original block editor adds text, code, links, removes blocks, validates, and publishes safely', async ({ page }, info) => {
  await login(page); await page.goto('./#/create-post');
  await expect(page.locator('#sidebar')).toHaveCount(0);
  await page.getByLabel('Title:', { exact: true }).fill('My original block editor discussion');
  await page.getByRole('button', { name: 'Select Categories' }).click();
  await page.locator('.category-option:has(input[value="1"])').click();
  await expect(page.getByRole('checkbox', { name: 'golang', exact: true })).toBeChecked();
  await page.getByRole('button', { name: 'Select Categories' }).click();
  await page.getByLabel('Text', { exact: true }).check();
  const malicious = '<img src=x onerror="window.injected=true">';
  await page.getByLabel('Block content').fill(`Review this code please. ${malicious}`);
  await page.getByRole('button', { name: 'Add Block', exact: true }).click();
  await expect(page.locator('#pre .text-string')).toContainText(malicious); await expect(page.locator('#pre img')).toHaveCount(0);
  await page.getByLabel('Code', { exact: true }).check(); await page.getByLabel('Block content').fill('fmt.Println("hello")');
  await page.getByRole('button', { name: 'Add Block', exact: true }).click(); await expect(page.locator('#pre .code-string')).toContainText('fmt.Println');
  await page.getByLabel('Link', { exact: true }).check(); await page.getByLabel('Block content').fill('https://go.dev/learn/');
  await page.getByRole('button', { name: 'Add Block', exact: true }).click(); await expect(page.locator('#pre a')).toHaveAttribute('href', 'https://go.dev/learn/');
  await page.getByRole('button', { name: 'Remove Last Block' }).click(); await expect(page.locator('#pre a')).toHaveCount(0);
  await page.screenshot({ path: info.outputPath('original-editor.png'), animations: 'disabled' });
  await page.getByRole('button', { name: 'Submit Post', exact: true }).click();
  await expect(page.locator('.vp-title')).toHaveText('My original block editor discussion');
  await expect(page.locator('#post-content img')).toHaveCount(0); expect(await page.evaluate(() => window.injected)).toBeUndefined();
  await page.goto('./#/posts'); await expect(page.locator('.post-title').first()).toHaveText('My original block editor discussion');
});

test('original animated login/register, profile tabs, local profile and photo editing', async ({ page }, info) => {
  await page.goto('./#/register');
  await expect(page.locator('.bg-rotator')).toBeVisible(); await expect(page.locator('.logo-rotator img')).toHaveCount(12);
  await page.getByLabel('New username', { exact: true }).fill('original_tester'); await page.getByLabel('Display name', { exact: true }).fill('Original Dev');
  await expect(page.locator('#password-reg')).toBeDisabled();
  await expect.poll(() => page.locator('.logo-rotator .s1').evaluate(element => Number(getComputedStyle(element).opacity))).toBeGreaterThan(.5);
  await page.screenshot({ path: info.outputPath('original-login.png') });
  await page.locator('#register-form').getByRole('button', { name: 'Register', exact: true }).click(); await expect(page.locator('#hero-title')).toBeVisible();
  await page.goto('./#/profile'); await expect(page.locator('.profile-photo')).toBeVisible(); await expect(page.locator('.typing-effect')).toHaveText('original_tester');
  await page.getByRole('link', { name: 'Edit Profile', exact: true }).click();
  await page.getByLabel('Display name', { exact: true }).fill('Updated Original Dev'); await page.getByLabel('Biography').fill('I prefer the original forum.');
  await page.getByLabel('Programming interests').selectOption(['golang', 'rust']);
  await page.getByRole('button', { name: 'Save Profile' }).click(); await page.reload();
  await expect(page.getByLabel('Display name', { exact: true })).toHaveValue('Updated Original Dev');
  const png = await readFile(new URL('../../../web/static/img/default.png', import.meta.url));
  await page.getByLabel('Choose New Photo', { exact: true }).setInputFiles({ name: 'avatar.png', mimeType: 'image/png', buffer: png });
  await page.getByRole('button', { name: 'Upload Photo' }).click();
  await expect(page.locator('#photoPreview')).toHaveAttribute('src', /^data:image\/jpeg;base64,/);
  await page.goto('./#/profile'); await expect(page.locator('.profile-photo')).toHaveAttribute('src', /^data:image\/jpeg;base64,/);
  await page.getByText('Liked Posts', { exact: true }).first().click(); await expect(page.locator('#liked-posts')).toBeVisible();
  await page.screenshot({ path: info.outputPath('original-profile.png'), animations: 'disabled' });
});

test('original playground highlights and downloads precise UTF-8 BOM and CRLF', async ({ page }, info) => {
  await page.goto('./#/playground'); await expect(page.getByRole('heading', { name: 'Code Playground' })).toBeVisible();
  await expect(page.locator('.note')).toContainText('Preview only.');
  await page.getByLabel('Language', { exact: true }).selectOption('python'); await page.getByRole('button', { name: 'Reset example' }).click();
  await expect(page.getByLabel('Your code')).toHaveValue(/def greet/);
  await page.getByLabel('Your code').fill('print("λ")\nprint("雪")\n'); await page.getByRole('button', { name: 'Preview highlighted' }).click();
  await expect(page.locator('.preview-content .hljs-built_in').first()).toHaveText('print');
  await page.getByLabel('Filename (optional, no extension)').fill('original-example');
  await page.getByLabel('Line endings').selectOption('crlf'); await page.getByLabel('Encoding').selectOption('bom');
  const pending = page.waitForEvent('download'); await page.getByRole('button', { name: 'Download file' }).click(); const download = await pending;
  expect(download.suggestedFilename()).toBe('original-example.py'); const bytes = await readFile(await download.path());
  expect([...bytes.subarray(0, 3)]).toEqual([239,187,191]); expect(bytes.subarray(3).toString('utf8')).toBe('print("λ")\r\nprint("雪")\r\n');
  await page.locator('#main-content').evaluate(element => element.scrollTo(0,0)); await page.screenshot({ path: info.outputPath('original-playground.png'), animations: 'disabled' });
  await page.getByRole('button', { name: 'Clear editor' }).click(); await expect(page.getByLabel('Your code')).toHaveValue('');
});

test('reset confirmation, unrelated storage, corrupt storage and quota recovery', async ({ page }) => {
  await login(page); await page.goto('./#/post/1'); await page.getByLabel('Your comment').fill('A temporary comment for reset.'); await page.getByRole('button', { name: 'Comment', exact: true }).click();
  await page.evaluate(() => localStorage.setItem('unrelated', 'keep'));
  await page.locator('.demo-tools summary').click(); await page.getByRole('button', { name: 'Reset Demo Data', exact: true }).click();
  await page.getByRole('button', { name: 'Cancel', exact: true }).click(); await expect(page.locator('.vp-comment')).toHaveCount(4);
  await page.getByRole('button', { name: 'Reset Demo Data', exact: true }).click();
  await page.getByRole('dialog').getByRole('button', { name: 'Reset Demo Data', exact: true }).click();
  await expect(page.locator('#hero-title')).toBeVisible();
  const state = await page.evaluate(key => ({ state: JSON.parse(localStorage.getItem(key)), other: localStorage.getItem('unrelated') }), key);
  expect(state.state.comments).toHaveLength(100); expect(state.state.currentUserId).toBeNull(); expect(state.other).toBe('keep');
  await page.evaluate(key => localStorage.setItem(key, '{broken'), key); await page.reload();
  await expect(page.locator('#storage-warning')).toContainText('restored'); await expect(page.locator('#featured-section .card')).toHaveCount(4);
  await page.addInitScript(() => { Storage.prototype.setItem = function () { throw new DOMException('Full', 'QuotaExceededError'); }; });
  await page.goto('./#/login'); await page.reload(); await page.locator('#login-form').getByRole('button', { name: 'Continue as Demo User' }).click();
  await expect(page.locator('#storage-warning')).toContainText('may be lost');
});

test('original pages at desktop, tablet and phone widths; hash links and local-only data', async ({ page }, info) => {
  const requests = []; page.on('request', request => { if (['xhr','fetch'].includes(request.resourceType())) requests.push(request.url()); });
  for (const width of [1440, 820, 390, 320]) {
    await page.setViewportSize({ width, height: 900 });
    for (const route of ['', 'posts', 'post/1', 'playground', 'category/javascript', 'login']) {
      await page.goto(`./#/${route}`);
      await expect(page.locator('body')).toContainText(route === 'login' ? 'Welcome To Coding Forum' : route === 'playground' ? 'Code Playground' : route === '' ? 'Welcome to Null Pointer' : route === 'posts' ? 'Discover Posts' : route === 'post/1' ? 'Understanding goroutines' : 'JavaScript Zone');
      const sizes = await page.evaluate(() => ({ width: innerWidth, scroll: document.documentElement.scrollWidth }));
      expect(sizes.scroll, `${route} at ${width}`).toBeLessThanOrEqual(sizes.width);
      if (!route && width === 390) await page.screenshot({ path: info.outputPath('original-mobile.png'), animations: 'disabled' });
    }
  }
  expect(requests).toEqual([]);
  await page.goto('./#/posts?q=goroutines'); await page.reload(); await expect(page.locator('.post-card')).toHaveCount(1);
  await login(page);
  for (const width of [390, 320]) {
    await page.setViewportSize({ width, height: 900 });
    for (const route of ['create-post', 'profile', 'edit-profile', 'help']) {
      await page.goto(`./#/${route}`);
      const sizes = await page.evaluate(() => ({ width: innerWidth, scroll: document.documentElement.scrollWidth }));
      expect(sizes.scroll, `${route} at ${width}`).toBeLessThanOrEqual(sizes.width);
    }
  }
});

test('original forms use local data and never request a backend API', async ({ page }) => {
  const backendRequests = [];
  await page.route('**/*', route => {
    if (['fetch', 'xhr'].includes(route.request().resourceType())) { backendRequests.push(route.request().url()); return route.abort(); }
    return route.continue();
  });
  await login(page);
  await navigation(page, '/posts');
  await page.getByLabel('Search', { exact: true }).fill('goroutines');
  await page.getByRole('button', { name: 'Apply', exact: true }).click();
  await expect(page.locator('.post-card')).toHaveCount(1);
  await page.locator('.post-card').getByRole('button', { name: 'Open post' }).click();
  await page.getByLabel('Your comment').fill('This original form works with backend API requests disabled.');
  await page.getByRole('button', { name: 'Comment', exact: true }).click();
  await expect(page.locator('.vp-comment').last()).toContainText('backend API requests disabled');
  await page.reload();
  await expect(page.locator('.vp-comment').last()).toContainText('backend API requests disabled');
  expect(backendRequests).toEqual([]);
});
