# Forum live demo - original interface

This static demo now renders the **actual templates in `web/templates/`** and loads the **actual styles in `web/static/css/`**. The previous custom demo design has been removed.

The home page, collapsed sidebar, window cards, login artwork and animations, block editor, profile tabs, post/comments cards, category feature sections, Code Playground, Help, Privacy, and error page come from the original application. The browser adapter replaces only Go-dependent data and actions. No original backend, template, CSS, database, or Docker file was modified.

## Run

Requires Node 22.12+ and npm. From the repository root:

```powershell
cd demo
npm ci
npm run dev
```

Open http://127.0.0.1:5173/Forum---A-Programming-Community-Platform/.

```powershell
npm run build
npm run preview
```

Production preview: http://127.0.0.1:4173/Forum---A-Programming-Community-Platform/.

The predev/prebuild script copies original static assets into ignored generated public folders. Vite bundles the original HTML/CSS as text. Neither development nor deployment needs Go or SQLite. Do not open the entry point as a file URL.

## How original rendering works

- `src/original/sources.js` imports real templates and styles directly from the repository, avoiding duplicate visual implementations.
- `src/original/template.js` interprets only the Go template directives used by these files: conditions, ranges, root/loop variables, simple functions and date formatting. It does not evaluate JavaScript. Dynamic text is escaped; only known rendered templates or syntax-highlighted code can be inserted as trusted HTML.
- `src/original/model.js` maps local state to the original Go view-model field names.
- `src/original/adapt.js` rewires server navigation/forms and adds demo-only behavior to the original markup.
- `src/original/adapter.css` contains limited accessibility, responsive overflow, and browser compatibility fixes. Original CSS supplies layout, colors, typography, cards, forms, and animation.
- `src/main.js` handles hash routes, browser actions, and page-specific stylesheet loading in the original order.
- `scripts/sync-original.mjs` copies original images/icons and supplies missing Terms/MATLAB/Brainfuck icon files. Original files stay untouched.

There is no custom top bar, home category tile grid, statistics strip, replacement card design, or additional marketing copy. The home uses four featured and four latest posts, matching the original Go query limits. The sidebar starts collapsed, as it does in the original root template.

The original Create Post editor supports Text/Code/Link blocks, Add Block, Remove Last Block, Clear All, and Submit Post. Original profile tabs switch between created and liked posts. Card minimize/open/hide controls work; the small Interactive Demo menu can restore hidden cards.

## Static-only differences

- Login retains the original animated page and form layout. Identity selection is simulated. Password inputs are disabled and are never collected or stored. Continue as Demo User signs in as `demo_dev`; entering an existing fictional username selects that profile.
- Registration creates a local fictional profile. The original email slot is used for a fictional display name. Read-only profile email text uses the reserved `example.invalid` domain.
- The original photo-upload UI accepts PNG/JPEG/WebP up to 5 MB, resizes locally to 128?128 JPEG, and persists only a bounded small image. No photo is uploaded to a server.
- The password-change card in Settings uses the original card/form styling for local display-name, biography, and interests editing.
- Playground previews and downloads run locally. Code is never executed or compiled. UTF-8, optional BOM, and LF/CRLF downloads are supported.
- Original Help/Privacy sections remain, with server/account claims corrected for the static demo.
- A small collapsed Interactive Demo menu provides the notice, reset, demo identity shortcut, source link, and hidden-card restoration. It does not change the original page layout.
- The missing `bg-tech-purple.png` reference falls back to the existing original `bg-tech.png`. The redundant remote Inter import is omitted; the original root stylesheet supplies the forum's system font.

## Data and persistence

The demonstration still has 18 fictional profiles, 40 discussions, and 100 comments. Content and seed relationships are in `src/data/`; state mutation/storage and feed selection are in `src/services/`. Search, category union filters, sorting, pagination, publishing, comments, reactions, profiles, and downloads remain functional behind the original UI.

Data uses `null-pointer:demo:v1`. Existing browser edits survive this UI replacement. Loading validates IDs, relationships, category names, dates, reaction maps, and optional bounded photo data. Corrupt data falls back to the seed with a notice. Storage failures leave an in-memory app and show a warning. Reset Demo Data requires confirmation, restores the seed, and logs out without clearing unrelated storage. No real authentication or shared backend is provided. Use fictional information.

## Tests

```powershell
npm test
npm run build
npx playwright install chromium
npm run test:e2e
```

CI installs Chromium using `npx playwright install --with-deps chromium`. Tests cover the original templates, visual identity and dimensions, original forms, state/persistence, malicious input, downloads, direct hash links, and desktop/tablet/phone overflow. Captures are generated under ignored `test-results/`. See [VERIFICATION.md](VERIFICATION.md).

## GitHub Pages

Expected public URL: https://hussainhht.github.io/Forum---A-Programming-Community-Platform/.

The existing `.github/workflows/deploy-demo.yml` builds/tests `live-demo` and deploys only `demo/dist/`. Repository Settings > Pages > Source must be GitHub Actions. The `github-pages` environment must allow `live-demo`; that branch permission was fixed earlier. No merge into `main` is required.

Changes to original templates/CSS should also trigger the workflow when pushed on `live-demo`, because they are now the demo's source. No Go or SQLite process is installed or run by the workflow.

To publish these local changes, from the repository root:

```powershell
git switch live-demo
git add demo
git diff --cached --stat
git commit -m "Restore original forum interface in static demo"
git push origin live-demo
```

No push is performed automatically. The public site keeps its existing version until these changes are pushed.
