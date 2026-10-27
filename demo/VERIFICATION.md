# Original-interface verification

The custom demo design has been replaced. The browser now renders the original application templates from `web/templates/` and original CSS from `web/static/css/`, instead of a separate visual implementation.

## Checks

- Production build succeeds using Node/Vite only. No Go or SQLite process is started.
- 16 Node tests cover data relationships, storage, reactions, profile edits, bounded photos, escaping, downloads, routing, the Go-template interpreter, all real templates, and original home structure.
- 10 Chromium scenarios cover the original screens and controls: home/sidebar; all eight category themes; Discover filters/search/sort/pagination/history; post/comments/reactions; block editing/publishing; animated login/registration/profile tabs/photo upload; playground bytes; storage/reset recovery; responsive widths; and local interactions with backend API requests blocked.
- The original home is explicitly checked for centered 48px heading, collapsed sidebar, four featured/four latest cards, and 360x300 card dimensions at desktop width. The removed custom header, statistics strip, tile grid, and callout are explicitly checked to be absent.
- Category feature grids and floating keyword art are retained and checked.
- Browser errors, console warnings, and failed HTTP asset responses fail the walkthrough.
- Responsive checks cover widths 1440, 820, 390, and 320 px; additional phone checks cover Create Post, Profile, Settings, and Help.
- Development server returns HTTP 200 for the entry point, application module, imported original template module, and original assets.

## Visual review

Screenshots of the original home, Go category, Discover page, animated login, block editor, profile, post/comments, Code Playground, and mobile home were reviewed. Representative current captures replace the former screenshots in `verification/`. All captures can be regenerated with `npm run test:e2e` under ignored `test-results/`.

## Preservation and deployment

All changes are confined to `demo/`. The original `web/`, `Internal/`, `database/`, `main.go`, Go modules, Dockerfile, and root README are unchanged. Work remains on `live-demo`; `main` is unchanged. Obsolete demo pages/components, custom stylesheet, duplicated asset copies, and previous screenshots have been removed.

The original Pages workflow remains in place and its environment already permits `live-demo`. This replacement has not been committed or pushed automatically. The public website keeps the previously deployed version until the owner pushes these changes.

Browser verification uses Chromium. Firefox/Safari and the Go backend were not exercised. Authentication remains simulated and the playground previews/downloads code without executing it. The app has no offline service worker: hosting must deliver the initial site/assets before browser-only interactions can work.
