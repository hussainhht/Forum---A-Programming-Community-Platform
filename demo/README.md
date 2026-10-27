# Forum · Interactive Live Demo

A frontend-only demonstration of the original **Forum – A Programming Community Platform**, branded **Null Pointer** in its home template. Visitors can explore a populated programming community without Go, SQLite, accounts, or a backend server.

Expected public URL after the owner enables Pages and pushes `live-demo`:

**https://hussainhht.github.io/Forum---A-Programming-Community-Platform/**

The demo is implemented separately in `demo/`. The original Go application, database schema, Dockerfile, templates, assets, and root README are unchanged. The original application still uses its existing instructions and port 7777.

## Run locally

Requires Node.js **22.12+** and npm. From the repository root:

```powershell
git switch live-demo
cd demo
npm ci
npm run dev
```

Open **http://127.0.0.1:5173/Forum---A-Programming-Community-Platform/**.

The first dependency installation needs internet access. Once installed, application data and interactions do not require a network service. Serve the app over HTTP; opening `index.html` as a `file://` URL is not supported.

To build and preview the exact deployment:

```powershell
npm run build
npm run preview
```

Open **http://127.0.0.1:4173/Forum---A-Programming-Community-Platform/**. Only `demo/dist/` is deployed. Vite preview is a local verification server; GitHub Pages serves the production files.

## Included features

- Home with featured posts, recent posts, programming categories, and working window-card controls.
- Discover Posts: title/content search, multiple category selections (match any), latest/popular/comment sorting, reset, and pagination with filters preserved in the URL.
- All eight category pages with the original theme colors and logos: Go, JavaScript, Rust, Ruby, Python, MATLAB, Brainfuck, and Java.
- Post details, safe paragraph and fenced-code rendering, author links, post/comment reactions, and local comments.
- Create Post with validation, multiple categories, code insertion, live preview, and immediate feed visibility.
- Simulated login, fictional registration, logout, profile switching, and the **Continue as Demo User** shortcut. No password fields or external authentication.
- Public profiles, created/liked posts, local display-name/biography/interests edits, and lightweight avatar color choices.
- Code Playground: 15 language/text choices, browser syntax highlighting, filename, LF/CRLF, optional UTF-8 BOM, file download, reset example, and clear editor. Code is **preview-only**, never executed.
- Help, demo-specific Privacy & Terms, missing-route/entity recovery, confirmation-based reset, storage warnings, mobile navigation, reduced-motion support, semantic forms, and keyboard focus states.

## Data and storage

The seed contains **18 fictional users, 40 distinct programming discussions, and 100 relevant comments**. Content includes goroutines, async/await, ownership, collection semantics, SQLite indexes, matrix operations, numerical methods, and esolang tape tracing. Dates and IDs are deterministic; posts and comments always reference existing users. Initial engagement uses actual per-user reaction maps, so counts and active buttons derive from the same data.

`src/data/discussions.json` contains authored discussions, code examples, and technical replies. `src/data/seed.js` contains fictional profiles, subject-specific follow-ups, and deterministic relationships/timestamps. `src/data/categories.js` preserves the Go handler's category palette.

On first visit, the app saves a complete version-1 snapshot under **`null-pointer:demo:v1`**. Reloads restore this snapshot without appending the seed again. New posts/comments, reactions, profile edits, fictional registrations, and the selected identity persist on the same browser origin.

Malformed JSON, invalid shapes, duplicate identities, and broken references fall back to the seed with a visible warning. Unreadable data is not immediately overwritten; the next successful edit or explicit reset writes a valid snapshot. Storage read/write failures are caught. When a write fails, the app continues in memory and warns that changes may be lost on refresh.

**Reset Demo Data** asks for confirmation, replaces only the demo's key with the original seed, and logs out. It leaves unrelated localStorage keys alone. Other browser profiles/devices do not share edits. No passwords, real emails, or uploaded images are collected. Use fictional content.

Unpublished post and playground drafts survive navigation within the running app but are not persisted after refresh. Published content is persisted.

## Architecture

Vanilla JavaScript ES modules, HTML, CSS, JSON, localStorage, Vite, and locally bundled highlight.js. Playwright and Node's test runner are development-only dependencies.

```text
demo/
  index.html                       # Static entry point
  package.json / package-lock.json  # Frontend dependencies and commands
  vite.config.js                   # Repository base path
  playwright.config.js             # Production-build browser walkthrough
  public/assets/                   # Reused original icons/art and two missing logos
  src/
    main.js                        # App lifecycle and DOM event wiring
    router.js                      # Hash parsing and query links
    components/                    # Shared shell, cards, author/avatar, forms
    data/                          # Category definitions, JSON discussions, seed
    services/                      # State persistence, mutations, feed selectors
    pages/                         # Feed, post editor/detail, account, tools, info
    styles/original/               # Copies of original root/home/category CSS
    styles/demo.css                # Responsive adaptations and app styles
    utils/content.js               # Escaping, safe formatting, downloads
  tests/state.test.js              # State, storage, formatting, byte-level checks
  tests/browser/demo.spec.js        # End-to-end walkthrough
  VERIFICATION.md                  # Actual verification results and limitations
  FILES.md                         # Complete created-file inventory
```

Hash routing means Pages does not need SPA rewrite support. Examples:

```text
#/
#/posts?q=goroutines&category=golang&sort=popular&page=1
#/post/1
#/create-post
#/category/rust
#/profile
#/profile/user-3
#/edit-profile
#/playground
#/login
#/register
#/help
#/privacy-terms
```

All local asset URLs use Vite's configured base; external documentation and GitHub links are explicit. No application requests are sent to the Go routes.

User text is escaped before HTML rendering. The formatter supports paragraphs and fenced code, with highlight.js escaping the code input. Raw HTML, arbitrary embeds, and executable user links are never inserted into the DOM. Download output is generated with `TextEncoder`, `Blob`, and revocable object URLs.

## Test

From `demo/`:

```powershell
npm test
npm run build
npx playwright install chromium
npm run test:e2e
```

Linux CI uses `npx playwright install --with-deps chromium`. Browser tests start and stop the production preview automatically on port 4173. They cover actual UI interactions, direct hash links, refresh/history, browser storage, downloads, malicious input, broken requests, and desktop/tablet/mobile sizes. Screenshots and failure traces go into ignored `test-results/`.

See [VERIFICATION.md](VERIFICATION.md) for the tested behavior and deployment limits. The initial dependency audit reported no vulnerabilities.

## Publish using GitHub Pages

The workflow **`.github/workflows/deploy-demo.yml`** runs only on pushes to **`live-demo`** (or manual dispatch of that branch). It installs Node dependencies, runs the unit tests, builds, installs Chromium, runs the browser tests, uploads **only `demo/dist`**, and deploys through the official Pages actions. Go and SQLite are not installed or started. The deploy job has `pages: write` and `id-token: write`; repository content access remains read-only.

One-time owner settings:

1. In the repository, open **Settings → Pages → Build and deployment → Source**, then select **GitHub Actions**.
2. Ensure **Settings → Actions → General** allows GitHub Actions and the official `actions/*` actions used in the workflow.
3. If the **github-pages** environment restricts deployment branches, allow **`live-demo`**. Any protection requiring owner approval still applies.
4. Push `live-demo`, then inspect **Actions → Deploy Forum interactive demo**. Use the deployment URL from the successful deploy job.

These steps follow [GitHub's custom Pages workflow documentation](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages) and [Vite's repository-subdirectory deployment guidance](https://vite.dev/guide/static-deploy.html#github-pages).

No repository settings, commits, pushes, or public deployments were performed automatically.

From the **repository root**, review and publish when ready:

```powershell
git switch live-demo
git status --short
git add demo .github/workflows/deploy-demo.yml
git diff --cached --stat
git commit -m "Add interactive frontend demo and GitHub Pages deployment"
git push -u origin live-demo
```

There is no need to merge into or modify `main`. If already inside `demo/`, run `cd ..` before the Git commands.

## Differences from the Go application

| Original Go/SQLite feature | Static demo behavior |
| --- | --- |
| Password verification, sessions, bcrypt | Fictional profile selection persisted locally; no security boundary |
| Database-backed posts/comments/reactions | Versioned local browser snapshot; no shared community writes |
| Server-rendered HTML | Browser components with hash routes |
| Chroma preview and Go download endpoints | Bundled highlight.js, UTF-8 bytes, Blob download |
| Profile photo upload/password change | Initials, avatar color, display name, bio, interests; no uploads/passwords |
| Text/code/link block editor | Plain paragraphs and fenced code with live preview; links remain literal text |
| Server Help/Privacy descriptions | Demo-specific help and accurate local-data explanation |

The original source mentions code execution in some descriptions, but its playground service produces highlighted HTML and file downloads. This demo explicitly offers preview and download only. It has no interpreter, compiler, email, external auth, moderation service, uploads, or analytics.

The sidebar, home/category card styling, palette, category colors, and language/UI icons derive from the original project. Original category assets had no MATLAB or Brainfuck SVGs, so local compact SVGs fill those gaps. The copied Python logo has a corrected viewBox so the original artwork is readable. Dark monochrome logos get a light CSS treatment. Existing files under `web/` remain unchanged.
