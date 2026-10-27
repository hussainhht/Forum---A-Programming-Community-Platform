import { test, expect } from '@playwright/test';
import { readFile } from 'node:fs/promises';

const key = 'null-pointer:demo:v1';
const browserErrors = new WeakMap();
test.beforeEach(async ({ page }) => {
  const errors = [];
  browserErrors.set(page, errors);
  page.on('pageerror', error => errors.push(error.message));
  page.on('console', message => { if (['error', 'warning'].includes(message.type())) errors.push(message.text()); });
});
test.afterEach(async ({ page }) => { expect(browserErrors.get(page)).toEqual([]); });
async function demoLogin(page) {
  await page.goto('./#/login');
  await page.getByRole('main').getByRole('button', { name: 'Continue as Demo User', exact: true }).click();
  await expect(page.getByRole('link', { name: '@demo_dev', exact: true })).toBeVisible();
}

test('populated home, sidebar controls, all categories, help and invalid links', async ({ page }, testInfo) => {
  const errors = []; const badResponses = []; const requests = [];
  page.on('pageerror', error => errors.push(error.message));
  page.on('console', message => { if (message.type() === 'error' || message.type() === 'warning') errors.push(message.text()); });
  page.on('response', response => { if (response.status() >= 400) badResponses.push(`${response.status()} ${response.url()}`); });
  page.on('request', request => requests.push({ type: request.resourceType(), url: request.url() }));
  await page.goto('./#/');
  await expect(page.getByRole('heading', { name: 'Welcome to Null Pointer' })).toBeVisible();
  await page.getByRole('link', { name: 'Skip to content' }).focus();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('main')).toBeFocused();
  await expect(page).toHaveURL(/#\/$/);
  await page.evaluate(() => window.scrollTo(0, 0));
  await expect(page.locator('.forum-card')).toHaveCount(9);
  await page.screenshot({ path: testInfo.outputPath('home-desktop.png'), fullPage: true, animations: 'disabled' });
  const first = page.locator('.forum-card').first();
  await first.getByRole('button', { name: /^Minimize/ }).click();
  await expect(first.locator('.card-content')).toBeHidden();
  await first.getByRole('button', { name: /^Minimize/ }).click();
  await first.getByRole('button', { name: /^Hide/ }).click();
  await expect(first).toBeHidden();
  await page.getByRole('button', { name: 'Restore cards' }).click();
  await expect(first).toBeVisible();
  await page.locator('#sidebar').getByRole('button', { name: 'Show More' }).click();
  for (const [slug, label] of [['golang', 'Go'], ['javascript', 'JavaScript'], ['rust', 'Rust'], ['ruby', 'Ruby'], ['python', 'Python'], ['matlab', 'MATLAB'], ['brainfuck', 'Brainfuck'], ['java', 'Java']]) {
    await page.locator('#sidebar').getByRole('link', { name: label, exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`#/category/${slug}$`));
    await expect(page.getByRole('heading', { name: `The ${label} Zone` })).toBeVisible();
    expect(await page.locator('.forum-card').count()).toBeGreaterThanOrEqual(5);
    await expect(page.locator('.forum-card .badge').filter({ hasText: label }).first()).toBeVisible();
    if (slug === 'golang') await page.screenshot({ path: testInfo.outputPath('category-go.png'), fullPage: true, animations: 'disabled' });
  }
  for (const [route, title] of [['help', 'Help & Support'], ['privacy-terms', 'Privacy & Terms'], ['missing', '404 · Pointer not found'], ['post/invalid', '404 · Pointer not found'], ['profile/invalid', '404 · Pointer not found'], ['category/invalid', '404 · Pointer not found']]) {
    await page.goto(`./#/${route}`);
    await expect(page.getByRole('heading', { name: title, exact: true })).toBeVisible();
  }
  await page.goto('./#/help'); await page.reload();
  await expect(page.getByRole('heading', { name: 'Help & Support', exact: true })).toBeVisible();
  expect(errors).toEqual([]); expect(badResponses).toEqual([]);
  expect(requests.filter(request => ['fetch', 'xhr'].includes(request.type))).toEqual([]);
  expect(requests.every(request => request.url.startsWith('http://127.0.0.1:4173/Forum---A-Programming-Community-Platform/'))).toBeTruthy();
});

test('search, category union, sorting, filter reset, pagination, and browser history', async ({ page }) => {
  await page.goto('./#/posts');
  await expect(page.locator('.forum-card')).toHaveCount(9);
  const initialFirst = await page.locator('.forum-card h3').first().textContent();
  await page.getByLabel('Search', { exact: true }).fill('goroutines');
  await page.getByRole('button', { name: 'Apply filters' }).click();
  await expect(page.locator('.forum-card')).toHaveCount(1);
  await expect(page.locator('.forum-card h3')).toContainText('goroutines');
  await page.getByRole('link', { name: 'Reset filters' }).click();
  await page.getByLabel('Java', { exact: true }).check();
  await page.getByLabel('Ruby', { exact: true }).check();
  await page.getByLabel('Sort', { exact: true }).selectOption('comments');
  await page.getByRole('button', { name: 'Apply filters' }).click();
  await expect(page.locator('.result-count')).toContainText('10 discussions');
  await expect(page.locator('.forum-card .comment-link').first()).toContainText('3 comments');
  await page.getByRole('link', { name: 'Next →', exact: true }).click();
  await expect(page.locator('.pagination')).toContainText('Page 2 of 2');
  await expect(page.locator('.forum-card')).toHaveCount(1);
  await page.reload();
  await expect(page.getByLabel('Ruby', { exact: true })).toBeChecked();
  await page.goBack(); await expect(page.locator('.pagination')).toContainText('Page 1 of 2');
  await page.goForward(); await expect(page.locator('.pagination')).toContainText('Page 2 of 2');
  await page.getByRole('link', { name: 'Reset filters' }).click();
  await page.getByLabel('Sort', { exact: true }).selectOption('popular');
  await page.getByRole('button', { name: 'Apply filters' }).click();
  const counts = await page.locator('.forum-card [data-reaction="like"] .reaction-count').allTextContents();
  expect(counts.map(Number)).toEqual(counts.map(Number).sort((a, b) => b - a));
  await page.getByRole('link', { name: 'Reset filters' }).click();
  await expect(page.locator('.forum-card h3').first()).toHaveText(initialFirst);
  await page.getByLabel('Search', { exact: true }).fill('nothingmatches12345');
  await page.getByRole('button', { name: 'Apply filters' }).click();
  await expect(page.getByRole('heading', { name: 'No discussions yet' })).toBeVisible();
});

test('guest gate, authentication, post and comment reactions, comments and refresh', async ({ page }) => {
  await page.goto('./#/create-post');
  await expect(page.getByRole('heading', { name: 'Your next idea belongs here.' })).toBeVisible();
  await demoLogin(page);
  await page.goto('./#/post/1');
  await expect(page.locator('pre code').first()).toContainText('make(chan int)');
  const like = page.locator('.post-detail [data-reaction="like"]');
  const dislike = page.locator('.post-detail [data-reaction="dislike"]');
  const initial = Number(await like.locator('.reaction-count').textContent());
  await like.click(); await expect(like).toHaveAttribute('aria-pressed', 'true');
  await expect(like.locator('.reaction-count')).toHaveText(String(initial + 1));
  await dislike.click(); await expect(like).toHaveAttribute('aria-pressed', 'false');
  await expect(dislike).toHaveAttribute('aria-pressed', 'true');
  await dislike.click(); await expect(dislike).toHaveAttribute('aria-pressed', 'false');
  await like.click();
  const commentLike = page.locator('.comment [data-reaction="like"]').first();
  await commentLike.click(); await expect(commentLike).toHaveAttribute('aria-pressed', 'true');
  await page.getByLabel('Your comment').fill('I traced the channel close after the sender exits.\n\n```go\nclose(values)\n```');
  await page.getByRole('button', { name: 'Add comment', exact: true }).click();
  await expect(page.locator('.comment')).toHaveCount(4);
  await expect(page.locator('.comment').last()).toContainText('I traced the channel');
  await page.reload();
  await expect(page.locator('.comment')).toHaveCount(4);
  await expect(like).toHaveAttribute('aria-pressed', 'true');
  await expect(commentLike).toHaveAttribute('aria-pressed', 'true');
  await page.locator('#sidebar').getByRole('button', { name: 'Logout', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Welcome to Null Pointer' })).toBeVisible();
  await page.reload();
  await expect(page.locator('#sidebar').getByRole('link', { name: 'Login', exact: true })).toBeVisible();
});

test('create post validation, code preview, safe rendering, and immediately visible feed', async ({ page }, testInfo) => {
  await demoLogin(page);
  await page.locator('#sidebar').getByRole('link', { name: 'Create Post', exact: true }).click();
  await page.getByLabel('Post title', { exact: true }).fill('A safe demo post about channel ownership');
  const malicious = '<img src=x onerror="window.injected=true">';
  await page.getByLabel('Content', { exact: true }).fill(`Can someone review this example? ${malicious}`);
  await page.getByRole('button', { name: 'Publish post', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('Select at least one');
  await page.getByLabel('Go', { exact: true }).check();
  await page.getByRole('button', { name: '+ Code block', exact: true }).click();
  await expect(page.locator('#post-preview code')).toBeVisible();
  await page.screenshot({ path: testInfo.outputPath('create-post.png'), fullPage: true, animations: 'disabled' });
  await page.getByRole('button', { name: 'Publish post', exact: true }).click();
  await expect(page.locator('.post-detail h1')).toHaveText('A safe demo post about channel ownership');
  await expect(page.locator('.post-detail .formatted-content')).toContainText(malicious);
  await expect(page.locator('.post-detail .formatted-content img')).toHaveCount(0);
  expect(await page.evaluate(() => window.injected)).toBeUndefined();
  const postURL = page.url();
  await page.locator('#sidebar').getByRole('link', { name: 'Discover Posts', exact: true }).click();
  await expect(page.locator('.forum-card h3').first()).toHaveText('A safe demo post about channel ownership');
  await page.goto(postURL); await page.reload();
  await expect(page.locator('.post-detail h1')).toHaveText('A safe demo post about channel ownership');
});

test('registration, duplicate validation, choosing profiles and profile persistence', async ({ page }, testInfo) => {
  await page.goto('./#/register');
  await page.getByLabel('New fictional username').fill('async_owl');
  await page.getByLabel('Fictional display name').fill('Fictional Test');
  await page.getByRole('button', { name: 'Create demo profile', exact: true }).click();
  await expect(page.getByRole('alert')).toContainText('already in use');
  await page.getByLabel('New fictional username').fill('browser_tester');
  await page.getByRole('button', { name: 'Create demo profile', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Your Profile' })).toBeVisible();
  await page.getByRole('link', { name: 'Edit Profile', exact: true }).click();
  await page.getByLabel('Display name', { exact: true }).fill('Edited Fictional Dev');
  await page.getByLabel('Biography', { exact: true }).fill('I like deterministic tests and small programs.');
  await page.getByLabel('Rust', { exact: true }).check();
  await page.getByLabel('Color 3', { exact: true }).check();
  await page.getByRole('button', { name: 'Save Profile' }).click();
  await page.reload();
  await expect(page.getByLabel('Display name', { exact: true })).toHaveValue('Edited Fictional Dev');
  await expect(page.getByLabel('Rust', { exact: true })).toBeChecked();
  await page.getByRole('link', { name: 'Back to Profile' }).click();
  await expect(page.getByRole('heading', { name: 'Edited Fictional Dev' })).toBeVisible();
  await page.screenshot({ path: testInfo.outputPath('profile.png'), fullPage: true, animations: 'disabled' });
  await page.locator('#sidebar').getByRole('button', { name: 'Logout', exact: true }).click();
  await page.goto('./#/login');
  await page.getByLabel('Try an existing profile').selectOption('channel_chaser');
  await expect(page.getByLabel('Fictional username', { exact: true })).toHaveValue('channel_chaser');
  await page.getByRole('button', { name: 'Use this profile' }).click();
  await expect(page.getByRole('link', { name: '@channel_chaser', exact: true })).toBeVisible();
});

test('playground highlights safely and downloads exact UTF-8 CRLF BOM bytes', async ({ page }, testInfo) => {
  await page.goto('./#/playground');
  await expect(page.getByText('Preview only.', { exact: true })).toBeVisible();
  await page.getByLabel('Language', { exact: true }).selectOption('python');
  await page.getByRole('button', { name: 'Reset example' }).click();
  await expect(page.getByLabel('Your code')).toHaveValue(/def greet/);
  await page.getByLabel('Your code').fill('print("λ")\nprint("雪")\n');
  await page.getByRole('button', { name: 'Preview highlighted' }).click();
  await expect(page.locator('#code-preview .hljs-built_in').first()).toHaveText('print');
  await page.getByLabel('Filename (optional, no extension)').fill('unicode-example');
  await page.getByLabel('Line endings', { exact: true }).selectOption('crlf');
  await page.getByLabel('Encoding', { exact: true }).selectOption('bom');
  const downloadEvent = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Download file' }).click();
  const download = await downloadEvent;
  expect(download.suggestedFilename()).toBe('unicode-example.py');
  const bytes = await readFile(await download.path());
  expect([...bytes.subarray(0, 3)]).toEqual([239, 187, 191]);
  expect(bytes.subarray(3).toString('utf8')).toBe('print("λ")\r\nprint("雪")\r\n');
  await page.evaluate(() => window.scrollTo(0, 0));
  await page.screenshot({ path: testInfo.outputPath('playground.png'), fullPage: true, animations: 'disabled' });
  await page.getByRole('button', { name: 'Clear editor' }).click();
  await expect(page.getByLabel('Your code')).toHaveValue('');
  await page.getByLabel('Language', { exact: true }).selectOption('html');
  await page.getByLabel('Your code').fill('<script>window.playgroundInjected = true</script>');
  await expect(page.locator('#code-preview script')).toHaveCount(0);
  expect(await page.evaluate(() => window.playgroundInjected)).toBeUndefined();
});

test('confirmed reset restores seed, logs out, and leaves unrelated storage alone', async ({ page }) => {
  await demoLogin(page);
  await page.goto('./#/post/1');
  await page.getByLabel('Your comment').fill('This comment will be reset.');
  await page.getByRole('button', { name: 'Add comment', exact: true }).click();
  await page.evaluate(() => localStorage.setItem('unrelated-key', 'preserve-this'));
  await page.locator('#sidebar').getByRole('button', { name: 'Reset Demo Data', exact: true }).click();
  await expect(page.getByRole('dialog')).toBeVisible();
  await page.getByRole('button', { name: 'Cancel', exact: true }).click();
  await expect(page.locator('.comment')).toHaveCount(4);
  await page.locator('#sidebar').getByRole('button', { name: 'Reset Demo Data', exact: true }).click();
  await page.getByRole('dialog').getByRole('button', { name: 'Reset Demo Data', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Welcome to Null Pointer' })).toBeVisible();
  const result = await page.evaluate(key => ({ state: JSON.parse(localStorage.getItem(key)), unrelated: localStorage.getItem('unrelated-key') }), key);
  expect(result.state.users).toHaveLength(18); expect(result.state.posts).toHaveLength(40); expect(result.state.comments).toHaveLength(100);
  expect(result.state.currentUserId).toBeNull(); expect(result.unrelated).toBe('preserve-this');
  await page.reload();
  await expect(page.locator('#sidebar').getByRole('link', { name: 'Login', exact: true })).toBeVisible();
});

test('corrupt data and unavailable storage produce usable recovery notices', async ({ page }) => {
  await page.addInitScript(key => localStorage.setItem(key, '{bad-json'), key);
  await page.goto('./#/posts');
  await expect(page.locator('.forum-card')).toHaveCount(9);
  await expect(page.locator('#storage-warning')).toContainText('restored');
  await page.addInitScript(() => { Storage.prototype.setItem = function () { throw new DOMException('Full', 'QuotaExceededError'); }; });
  await page.reload();
  await page.locator('.topbar').getByRole('button', { name: 'Continue as Demo User' }).click();
  await expect(page.locator('#storage-warning')).toContainText('may be lost on refresh');
  await expect(page.locator('.forum-card')).toHaveCount(9);
});

test('desktop, tablet and mobile layouts have no overflow; mobile navigation works with keyboard', async ({ page }, testInfo) => {
  for (const viewport of [{ width: 1440, height: 1000 }, { width: 820, height: 1180 }, { width: 390, height: 844 }, { width: 320, height: 720 }]) {
    await page.setViewportSize(viewport);
    for (const route of ['', 'posts', 'post/1', 'playground', 'category/javascript', 'login']) {
      await page.goto(`./#/${route}`);
      await expect(page.locator('#main-content')).toBeVisible();
      const layout = await page.evaluate(() => ({ scroll: document.documentElement.scrollWidth, width: innerWidth }));
      expect(layout.scroll, `${route} at ${viewport.width}px`).toBeLessThanOrEqual(layout.width);
      if (viewport.width > 700) {
        const sidebar = await page.locator('#sidebar').boundingBox();
        const heading = await page.locator('#main-content h1').boundingBox();
        expect(heading.x, `Sidebar overlap at ${viewport.width}px on ${route}`).toBeGreaterThanOrEqual(sidebar.x + sidebar.width + 15);
      }
      if (!route && [820, 390].includes(viewport.width)) await page.screenshot({ path: testInfo.outputPath(`home-${viewport.width}.png`), fullPage: true, animations: 'disabled' });
    }
  }
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto('./#/');
  await page.getByRole('button', { name: 'Open navigation', exact: true }).click();
  await expect(page.locator('#sidebar')).toBeVisible();
  await page.locator('#sidebar').getByRole('link', { name: 'Discover Posts', exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Discover Posts' })).toBeVisible();
  await expect(page.locator('#sidebar')).toHaveAttribute('aria-hidden', 'true');
  await page.getByRole('button', { name: 'Open navigation', exact: true }).click();
  await page.keyboard.press('Escape');
  await expect(page.getByRole('button', { name: 'Open navigation', exact: true })).toBeFocused();
});

test('application data and interactions remain available offline after initial load', async ({ page, context }) => {
  await demoLogin(page);
  await context.setOffline(true);
  await page.locator('#sidebar').getByRole('link', { name: 'Discover Posts', exact: true }).click();
  await page.getByLabel('Search', { exact: true }).fill('ownership');
  await page.getByRole('button', { name: 'Apply filters' }).click();
  await expect(page.locator('.forum-card').first()).toContainText('Ownership');
  await page.locator('.forum-card h3 a').first().click();
  await page.getByLabel('Your comment').fill('This reply was written with network access disabled.');
  await page.getByRole('button', { name: 'Add comment', exact: true }).click();
  await expect(page.locator('.comment').last()).toContainText('network access disabled');
  await page.locator('.post-detail [data-reaction="like"]').click();
  await expect(page.locator('.post-detail [data-reaction="like"]')).toHaveAttribute('aria-pressed', 'true');
  await context.setOffline(false);
  await page.reload();
  await expect(page.locator('.comment').last()).toContainText('network access disabled');
});
