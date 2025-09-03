## 1) Project Snapshot

* **Repo / Commit**: (local workspace, commit hash not provided)
* **Go version**: 1.23.1 (from `go.mod`)
* **Run command**: `go run .`
* **Binary port**: 7777 (see `main.go` / `api.NewServer(7777, ...)`)
* **SQLite file**: `forum.db` (path: project root, opened as `./forum.db` in `database/db.go`)
* **.env / config**: missing ☐ (no `.env` or config loader present)

---

## 2) Scoring (Weighted → 100%)

Heuristic provisional scores (0–5).

|  # | Area | Weight | Score | Weighted |
| -: | ----- | -----: | ----: | -------: |
|  1 | Project Setup & Run | 10% | 3 | 6.0% |
|  2 | DB Schema & Migrations | 8% | 3 | 4.8% |
|  3 | Authentication | 10% | 3 | 6.0% |
|  4 | Authorization | 5% | 1 | 1.0% |
|  5 | Validation & Inline Errors | 8% | 2 | 3.2% |
|  6 | Posts CRUD | 8% | 2 | 3.2% |
|  7 | Categories & Filtering | 6% | 3 | 3.6% |
|  8 | Comments | 5% | 3 | 3.0% |
|  9 | Likes/Dislikes | 4% | 2 | 1.6% |
| 10 | Templates & UX | 6% | 2 | 2.4% |
| 11 | Security | 8% | 2 | 3.2% |
| 12 | Static Assets & Build | 3% | 3 | 1.8% |
| 13 | Logging & Error Handling | 4% | 2 | 1.6% |
| 14 | Testing | 3% | 0 | 0.0% |
| 15 | Performance & Concurrency | 3% | 2 | 1.2% |
| 16 | Docker/Deployment | 6% | 0 | 0.0% |
| 17 | Code Quality | 3% | 3 | 1.8% |
| 18 | Documentation | 3% | 1 | 0.6% |
|    | **Total** | **100%** |  | **41.0%** |

Notes:
* Missing edit/delete, pagination, filters for "my posts" / "my liked posts".
* Sessions stored with expiry but not validated on access.
* No Dockerfile, no tests, no CSRF.

---

## 3) Mandatory Feature Checklist

### Authentication

* [x] Registration asks for email, username, password (+ confirm)
* [x] Duplicate email detected
* [x] Wrong email/password rejected (`bcrypt.CompareHashAndPassword`)
* [x] Login session via cookie (`session_id`)
* [x] Only one active session per user (removes existing before creating new)
* [ ] Sessions have effective expiration check (field stored, but not enforced in code)
* [ ] Cookie sets `Secure` flag (only `HttpOnly`, `SameSiteLax` present)
* [x] Posts/comments visible to guests; creation requires login (checked)

Manual Tests (status / needed):
* [ ] Multi‑browser invalidation (not yet manually verified)
* [ ] Attempt login with empty form returns inline error (basic error path exists but needs manual confirm)

### SQLite

* [x] CREATE (schema in `database/db.sql`)
* [x] INSERT (users/posts/comments/sessions/reactions)
* [x] SELECT queries for entities
* [x] New records visible after actions (queries present)

### Likes/Dislikes

* [x] Restricted to logged‑in users (session checks before handlers proceed)
* [x] Counts visible (post like count; comment counts not shown yet)
* [x] Cannot like & dislike same post simultaneously (update logic in `PostReaction`)
* [ ] Comment reactions toggle/update (currently always insert; duplicate toggle not handled / may error on change)
* [ ] Dislike counts surfaced (not displayed)

### Filter

* [x] Filter by category (`?category=...`)
* [ ] Filter by my created posts
* [ ] Filter by my liked posts

### Docker

* [ ] Dockerfile present
* [ ] Image build succeeds
* [ ] Container run documented
* [ ] Compose setup

### HTTP & Errors

* [x] Appropriate methods (GET/POST used; no PUT/DELETE yet)
* [ ] 404 custom for all unknown routes (root handler only; others may 404 plain)
* [ ] 400 handling consistent (some direct `http.Error` without template)
* [ ] 500 friendly page always (mixed usage)

### Allowed Packages Only

* [x] Stdlib + `sqlite3` + `bcrypt` + `uuid` (crypto/x packages). No disallowed frameworks.

---

## 4) Evidence

Routes (from `server.go`):
```
GET  /                 (home, category filtering)
GET  /register
POST /register
GET  /login
POST /login
GET  /create-post
POST /create-post
POST /create-comment
GET  /post/{id}
POST /post-reaction
POST /comment-reaction
GET  /logout
```

DB Schema (simplified):
```
users(id, username UNIQUE, email UNIQUE, password)
posts(id, user_id FK users, title, content, created_at)
categories(id, name UNIQUE)
post_categories(post_id FK posts, category_id FK categories, PK(post_id, category_id))
comments(id, user_id FK users, post_id FK posts, content)
post_reactions(user_id, post_id, reaction_type CHECK like|dislike, PK(user_id, post_id))
comments_reactions(user_id, comment_id, reaction_type CHECK like|dislike, PK(user_id, comment_id))
sessions(id, user_id FK users, session_id UNIQUE, expires_at)
```

Security Checks:
* Prepared statements / parameter binding used (`?` placeholders across queries)
* Template auto-escaping (Go `html/template`)
* Missing CSRF token pattern (no token generation or hidden field)
* Cookie lacks `Secure` flag; expiry not enforced (no check vs `expires_at` on read)
* No output encoding for user content sanitization beyond template escape (risk if HTML allowed in content)

---

## 5) Findings & TODO (Prioritized)

### High
* [ ] Add CSRF protection (generate token per session, store server-side, embed hidden input, verify on POST)
* [ ] Enforce session expiration (check `expires_at` each request; purge expired sessions)
* [ ] Add `Secure` & `SameSite=Strict` (or Lax justification) to session cookie in production
* [ ] Implement comment reaction update / toggle logic with UPSERT or update on conflict

### Medium
* [ ] Add filters: my posts, my liked posts
* [ ] Implement pagination on posts list (LIMIT/OFFSET with `?page=`)
* [ ] Unify error rendering (use `error.html` for validation errors instead of raw `http.Error`)
* [ ] Add edit/delete for posts (owner-only) and delete for comments
* [ ] Enhance validation feedback on create-post (inline vs plain 400)

### Low
* [ ] Dockerfile + optional compose for reproducible run
* [ ] Unit tests for auth, post creation, reactions
* [ ] Session cleanup job for expired sessions
* [ ] Logging wrapper with contextual info (user ID, route, latency)
* [ ] Markdown or plain text sanitization / allowlist for post & comment content
* [ ] Display dislike & comment reaction counts
* [ ] Add schema & routes diagrams to docs

---

## 6) Verdict

* **Pass / Fail**: Fail (core features missing: pagination, filters (beyond category), Docker, CSRF, tests)
* **Why fail**: Security gaps (no CSRF, cookie flags incomplete), missing mandatory extended filters, lack of deployment artifacts, partial reactions implementation.
* **Final Completion %**: ~41%

---

## 7) Quick Commands

```bash
go run .
go test ./...
sqlite3 forum.db ".schema"
sqlite3 forum.db "SELECT COUNT(*) FROM users;"
```

When Dockerfile added:
```bash
docker build -t forum:latest .
docker run --rm -p 7777:7777 --name forum forum:latest
```

---

## 8) Notes

* Add `.env.example` to configure port, DB path, and secure cookie toggle.
* Consider middleware layer for auth/session + CSRF to DRY handlers.
* Introduce `context.Context` request scoping for user ID after session validation.
* Use `PRAGMA foreign_keys = ON;` (SQLite) if not already implicitly active.
