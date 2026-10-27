# Current demo files

The previous custom visual implementation has been removed. The original `web/templates/` and `web/static/css/` files are now the interface source. Those original files are unchanged.

Removed: the former `src/pages/`, `src/components/`, `src/styles/`, duplicated `public/assets/`, and previous screenshots. The current screenshots show the original interface.

Kept: local seed data, persistence, hash routing, safe content utilities, and the existing Pages workflow.

Generated original assets under `public/web/` and `public/profile-img/`, dependencies, builds, and test output are ignored by Git.

## Current files

- [demo/.gitignore](.gitignore)
- [demo/FILES.md](FILES.md)
- [demo/README.md](README.md)
- [demo/VERIFICATION.md](VERIFICATION.md)
- [demo/index.html](index.html)
- [demo/package-lock.json](package-lock.json)
- [demo/package.json](package.json)
- [demo/playwright.config.js](playwright.config.js)
- [demo/scripts/sync-original.mjs](scripts/sync-original.mjs)
- [demo/src/data/categories.js](src/data/categories.js)
- [demo/src/data/discussions.json](src/data/discussions.json)
- [demo/src/data/seed.js](src/data/seed.js)
- [demo/src/main.js](src/main.js)
- [demo/src/original/adapt.js](src/original/adapt.js)
- [demo/src/original/adapter.css](src/original/adapter.css)
- [demo/src/original/model.js](src/original/model.js)
- [demo/src/original/sources.js](src/original/sources.js)
- [demo/src/original/template.js](src/original/template.js)
- [demo/src/router.js](src/router.js)
- [demo/src/services/feed.js](src/services/feed.js)
- [demo/src/services/storage.js](src/services/storage.js)
- [demo/src/utils/content.js](src/utils/content.js)
- [demo/tests/browser/demo.spec.js](tests/browser/demo.spec.js)
- [demo/tests/state.test.js](tests/state.test.js)
- [demo/tests/templates.test.js](tests/templates.test.js)
- [demo/verification/category-go.png](verification/category-go.png)
- [demo/verification/create-post.png](verification/create-post.png)
- [demo/verification/discover.png](verification/discover.png)
- [demo/verification/home-desktop.png](verification/home-desktop.png)
- [demo/verification/home-mobile.png](verification/home-mobile.png)
- [demo/verification/login.png](verification/login.png)
- [demo/verification/playground.png](verification/playground.png)
- [demo/verification/profile.png](verification/profile.png)
- [demo/vite.config.js](vite.config.js)
