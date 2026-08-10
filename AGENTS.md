# AGENTS.md — GeocachingBrughia VZW

Read this **entire file** before making any change to this repository.
It exists because the codebase was hand-written (Vue frontend), then later
vibe-coded (admin panel, Golden Key, and the **entire Go backend** that
replaced the original Laravel API). That mix produced real bugs and several
silent-failure traps that an LLM can easily make worse. This file is the
defensive line against that.

> Short version: **production data is sacred, the migration system is fragile,
> several admin features silently no-op, and there are no tests.** Read on.

---

## 1. What this is

A single-deploy website for a Belgian non-profit (GeocachingBrughia VZW).
One `docker compose` runs a Go API + Vue SPA behind Nginx, with SQLite in a
named Docker volume. The site has a public side (home, events, geocaches,
contact form, Golden Key game) and an admin panel (auth-gated CRUD for all
content + contact-form inbox).

## 2. Tech stack & commands

**Backend** — Go 1.22, Chi router v5, SQLite (mattn/go-sqlite3, WAL mode),
JWT (golang-jwt/v5), bcrypt, gomail for SMTP. Distroless runtime image.

```bash
# from repo root
cd backend
cp ../.env.example ../.env        # then fill in real values
go mod download
go run .                           # serves on :8080

# verify (there is NO test suite — use these instead)
go build ./...                     # must compile
go vet ./...                       # must be clean
gofmt -l .                         # must print nothing
```

**Frontend** — Vue 3 (Composition API, `<script setup>`), Vite 5,
vue-router 4, TipTap 2 rich text, jszip + qrcode for event QR codes.

```bash
cd frontend
npm install
npm run dev      # Vite dev server, proxies /api/ to localhost:8080
npm run build    # production build → frontend/dist/ (must pass)
npm run preview  # preview the built bundle
```

**Full stack (production-like):**

```bash
docker compose up -d --build       # from repo root
docker compose logs backend | grep -A6 "DEFAULT ADMIN"   # first-boot password
docker compose logs -f
```

There is **no test framework, no linter config, and no CI** in the repo. Your
verification is: `go build ./...`, `go vet ./...`, `gofmt -l .`, and
`npm run build`. Run all of them before considering work done.

## 3. Repository layout

```
backend/
  main.go                          # bootstrap: config, db, migrations, seed, email, server
  internal/
    config/config.go               # env loading + prod validation (JWT secret ≥32 chars)
    database/
      database.go                  # *sql.DB wrapper, opens SQLite with WAL + FK on
      migrations.go                # ALL migrations, runs every boot (see §6)
      seed.go                      # default languages, static content, first admin
    handlers/                      # HTTP handlers, one file per resource
      handlers.go                  # shared Handler struct, image upload/serve, auth helper
      auth.go                      # login, change-password, profile, refresh, respondJSON
      contact.go                   # contact form + admin inbox (BUG — see §5.1)
      messages.go                  # site announcements
      events.go                    # events CRUD + QR-code ZIP generation
      geocaches.go                 # geocaches CRUD (see §7 quirk)
      golden_key.go                # golden key settings (BUG — see §5.2)
      golden_key_months.go         # monthly entries + hints (BUG — see §5.2)
      static.go                    # languages, static content, socials
      users.go                     # admin user management + invitations
    middleware/
      auth.go                      # JWTAuth, UserClaims in context
      cache.go                     # ETag (no-cache revalidate) + NoCache
      ratelimit.go                 # per-IP login rate limiter
    router/router.go               # route table — single source of truth for endpoints
    services/email/email.go        # SMTP + hourly reminder scheduler
frontend/
  src/
    App.vue                        # public chrome (header, side menu, language popup)
    router/index.js                # route table — public + admin + golden key
    data/config.js                 # apiUrl, defaultLanguage
    services/                      # API wrappers
      fetcher.js                   # fetchFromServer/fetchToServer/deleteFromServer
      AdminService.js              # admin API helpers (partially used — see §7)
      ContactService.js, EventService.js, GeocacheService.js,
      GoldenKeyService.js, GoldenKeyMonthService.js, MessageService.js,
      SocialService.js, StaticContentService.js, LanguageService.js,
      ConsoleService.js
    views/                         # one .vue per route
    components/                    # shared + admin layout components
      admin/AdminLayout.vue        # auth gate + sidebar + content slot
      admin/AdminSidebar.vue       # nav (NOT mobile-friendly — see §5.3)
      admin/ToastNotification.vue, TipTapEditor.vue, etc.
    css/admin.css                  # admin design system + CSS vars (sidebar width here)
data/                              # SQLite DB lives here (or Docker volume /data)
docker-compose.yml, Dockerfile.backend, Dockerfile.frontend, nginx.conf, certs/
```

## 4. Architecture & data flow (things that bite)

- **Auth:** JWT in `localStorage` as `admin_token`. Every admin fetch sets
  `Authorization: Bearer <token>`. `middleware.JWTAuth` validates and puts
  `UserClaims` in context. Token lifetime = `JWT_EXPIRY_HOURS` (24h default).
  There is **no token blocklist** — logout is client-side only.
- **Public API caching:** `middleware.CacheControl` wraps public GETs with
  ETag (MD5 of body) + `Cache-Control: no-cache` (revalidate). Nginx
  `proxy_cache` in front with the same ETag revalidation. Admin routes use
  `NoCache` (`no-store`).
- **First-boot admin:** `seed.go` creates one user with email literally
  `admin` (NOT an email address) and a random password printed once to stdout.
  `needs_password_update = 1` forces a password change on first login.
- **Email service:** optional. If `SMTP_HOST` is empty, emails are skipped
  but contact submissions are still saved. An hourly background job sends
  reminders for `status='new'` submissions older than `REMINDER_DAYS`.
- **Two parallel API helpers in the frontend:** (1) `services/fetcher.js`
  used by the `*Service.js` modules, and (2) **per-view `apiRequest()`
  functions** copy-pasted inside `AdminContactsView`, `AdminMessagesView`,
  `AdminGoldenKeyView`, `AdminGoldenKeyMonthView`, `AdminLayout`. These
  helpers behave differently — notably the contacts one **does not check
  `response.ok`**, which is the root cause of §5.1.

## 5. Known bugs — fix these, don't introduce new ones

These are confirmed, reproducible, and have been traced to root cause.
Do not "clean up" unrelated code while fixing them; make minimal, targeted
edits and verify with the commands in §2.

### 5.1 Contact form status does not persist (HIGH priority)

**Symptom:** In the admin "Contactberichten" inbox, changing
Nieuw / In Behandeling / Opgelost / Gesloten appears to work, then reverts
on refresh. Notes also never appear.

**Root cause — multiple mismatches in `AdminContactsView.vue`:**

1. **Wrong endpoint.** `updateStatus()` calls
   `PUT /api/admin/contacts/{id}/status` (`AdminContactsView.vue:347`), but
   the backend registers only `PUT /api/admin/contacts/{id}`
   (`backend/internal/router/router.go:132`). There is no `/status` route.
   Chi returns 404/405.
2. **Silent failure.** The view's local `apiRequest()` only checks for 401
   (`AdminContactsView.vue:238`) and never checks `response.ok`. The 404 is
   swallowed, the UI optimistically flips `contactDetails.value.status`,
   and the change is never persisted. On reload it reverts.
3. **Notes body mismatch.** `addNote()` sends `{ content: ... }`
   (`AdminContactsView.vue:376`), but `AddContactNote` decodes into
   `{ Note string `json:"note"` }` (`contact.go:265-267`) and rejects empty
   `note` with 400. Notes are never saved.
4. **Notes response/display mismatch.** Backend returns
   `{ id, message }` (`contact.go:296`); the view pushes the raw response
   as a note object and reads `note.admin_email` / `note.content`
   (`AdminContactsView.vue:170-173`). The backend `ContactNote` JSON fields
   are `user_name` and `note`. All three disagree.
5. **`contactDetails.name` does not exist.** There is no `name` column on
   `contact_submissions` and no `name` field in the `ContactSubmission`
   struct, yet the view shows `contactDetails.name` (`:112`) and uses it in
   the mailto body (`:425`). It is always `undefined`.

**Fix direction (pick the less invasive option unless told otherwise):**
- Add `r.Put("/contacts/{id}/status", h.UpdateContactStatus)` to the router
  **and** keep `UpdateContactSubmission` for the general case, OR change the
  view to call `PUT /api/admin/contacts/{id}` with `{ status }` (matching the
  existing handler). Either way, also: make `apiRequest` check `response.ok`
  and surface errors; align the note body field (`note` vs `content`) and the
  note response shape; remove or populate `name`. Whatever you choose, the
  **frontend and backend must agree on field names and routes** — verify by
  reading both sides, not just one.

### 5.2 Time handling: local time shifted on re-edit (HIGH priority)

**Symptom:** Admin saves a datetime; the next time they open it, the value is
shifted by the timezone offset.

**Root cause — naive strings sent to the browser:**

- The backend stores datetimes as **naive UTC TEXT** (`2006-01-02 15:04:05`,
  no timezone) for `golden_key_settings.activation_time`,
  `golden_key_months.live_date`, and `golden_key_months.found_date`
  (`golden_key.go:74`, `golden_key_months.go:275`).
- `GetGoldenKeySettings` **parses** the string via `parseFlexibleTime` and
  returns a `time.Time`, which JSON-marshals as RFC3339 with `Z`. That path
  is correct.
- `GetGoldenKeyMonths` / `GetGoldenKeyMonthByID` /
  `GetAdminGoldenKeyMonthByID` do **NOT** parse — they return
  `GoldenKeyMonth.LiveDate` and `FoundDate` as **raw naive strings**
  (`golden_key_months.go:18-19, 95, 143, 185, 226`). The browser then does
  `new Date("2026-04-12 10:12:00")`, which is **implementation-defined** for
  non-ISO space-separated strings: some engines treat it as local, some as
  UTC, Node treats it as UTC. Either way the round-trip is inconsistent with
  the save path (which sends proper RFC3339 with offset), so the value shifts
  by the TZ offset on the next edit.

**Fix direction:** Make the months endpoints return the same shape the
settings endpoint does — parse `liveDateStr`/`foundDate` with
`parseFlexibleTime` and expose them as `time.Time` (or as RFC3339 strings
with a `Z`). The frontend `toLocalDatetimeInput` already expects a real UTC
instant. Do **not** "fix" this by changing how the value is stored; the
storage format is shared with existing rows. Fix the **output** so it carries
an explicit timezone. Then re-verify the full round-trip: save → reload →
input shows the same wall-clock time the admin typed.

Also check `events.start_date`/`end_date` — they are passed through as raw
strings end-to-end and may have the same class of problem if any UI ever
round-trips them through a `datetime-local` input.

### 5.3 Golden Key scrolling issues on mobile (MEDIUM)

**Symptom:** The Golden Key "soon" countdown page and the active months grid
don't scroll / get clipped on mobile.

**Cause:** `GoldenKeyView.vue` uses `min-height: calc(100vh - 4.5rem)` and
`overflow-y: auto` on `.gk-soon`/`.gk-active`. `100vh` on mobile browsers is
the *large* viewport (includes browser chrome), so the element is taller than
the visible area, the inner `overflow-y` never engages, and the fixed header
overlaps content. `App.vue` also doesn't constrain the public layout height.
Use `100dvh` (dynamic viewport) where supported, and make sure the scrolling
container is the right element. Same applies to `GoldenKeyMonthView.vue`
(`.gkm { min-height: calc(100vh - 4.5rem); overflow-y: auto }`).

### 5.4 Admin panel not consistently mobile-friendly (MEDIUM)

**Symptom:** Some admin pages are usable on mobile, others aren't.

**Cause:** `AdminLayout.vue` + `AdminSidebar.vue` + `css/admin.css` have
**no responsive rules at all**. The sidebar is `position: fixed;
width: var(--sidebar-width)` = 260px always, and `.admin-main` has
`margin-left: var(--sidebar-width)`. On a 375px phone that leaves ~115px of
content. The only collapse mechanism is a 72px icon rail — there is no
mobile drawer/overlay. Meanwhile `AdminContactsView.vue` *does* have a
`@media (max-width: 1024px)` block. So responsiveness is sprinkled per-view,
inconsistently.

**Fix direction:** Introduce one mobile pattern in `AdminLayout`/`AdminSidebar`
(drawer + overlay + hamburger) and apply it everywhere; remove the per-view
workarounds that conflict. Keep it consistent.

### 5.5 Migration system is fragile (RISK — read §6 fully)

Not a "bug" to fix in code so much as a **trap to never trigger**. See §6.

### 5.6 Minor / low priority

- `AdminLayout.vue:67-78` calls `GET /api/setup-status` and `POST /api/register`
  which **do not exist** in the Go backend (Laravel-era leftovers). The errors
  are caught silently, so `needsSetup` stays false and the "Welcome! Create
  your administrator account" branch is dead code. Safe to remove.
- `AdminMessagesView.vue` reads `message.updated_at` (`:353`) but the
  `GetAdminMessages` `Message` struct (`messages.go:13-18`) does not include
  `updated_at`, so the column is always blank. The `getStateBadge` map still
  has uppercase `ONLINE`/`DRAFT`/`ARCHIVED` keys — also Laravel leftovers;
  the backend uses lowercase.
- `UpdateMessage` (`messages.go:140`) does **not validate** `state` — any
  string is accepted. The contact handler *does* validate. Make these
  consistent (validate against an allow-list) if you touch messages.
- `GoldenKeyService.updateGoldenKeySettings` (`GoldenKeyService.js:20`) does
  not send `rules` in the body; `AdminGoldenKeyView.vue` uses its own
  `apiRequest` and does send `rules`. The service function is effectively
  half-broken/unused — prefer the view's path or fix the service.
- `fetcher.js` `fetchFromServer` swallows network errors and returns `[]`,
  but returns `{access_denied: true}` on 401/403. Code that checks
  `if (!data)` will miss a network failure (empty array is truthy-but-empty).
  Always handle both.

## 6. The migration system — READ BEFORE TOUCHING THE DB

`backend/internal/database/migrations.go` is **the most dangerous file in
this repo for an LLM to edit.** Read this twice.

**How it actually works (it is not a normal migration system):**

- There is **no schema version table**. There is no `migrations` tracking.
- `db.Migrate()` runs **every migration, every boot**, in order.
- `CREATE TABLE IF NOT EXISTS` and `CREATE INDEX IF NOT EXISTS` make the
  CREATE-style migrations idempotent. Fine.
- `ALTER TABLE … ADD COLUMN` migrations are **not** idempotent. SQLite
  errors on a duplicate column. The code handles this by **string-matching
  the migration name** in a hardcoded list (`migrations.go:295`) and
  swallowing the error with a `(column already exists)` log. If you add a
  new ALTER migration and forget to add its name to that list, **the server
  will fail to boot on every restart after the first** because the second run
  returns a non-swallowed error.
- `seed_golden_key_months` uses `INSERT OR IGNORE` against the
  `month_number UNIQUE` constraint, so it only seeds missing rows. Editing
  an existing month's seed values will **not** update already-inserted rows.

**Rules for an LLM editing `migrations.go`:**

1. **Never modify an existing migration's SQL** to "fix" it. Existing rows in
   production depend on the exact shape that ran the first time. Add a **new**
   migration instead.
2. **Never reorder** migrations. They run in slice order.
3. If you add an `ALTER TABLE … ADD COLUMN`, you **must** also add the
   migration's `name` to the swallow-list at `migrations.go:295`, or the
   server will crash-loop on restart in production.
4. New `CREATE TABLE` migrations must use `IF NOT EXISTS`.
5. New seed-data migrations must use `INSERT OR IGNORE` (or equivalent) so
   they are re-runnable without clobbering admin edits.
6. **Never write a migration that drops or recreates a table that already
   has production data.** There is no data-migration tooling and no backup
   path in the deploy. Forward-only, additive, non-destructive.
7. SQLite has limited ALTER support (no DROP COLUMN before 3.35, no
   rename-with-rebuild). If you need a destructive schema change, you must
   plan a multi-deploy migration (add new → backfill → switch reads →
   switch writes → much later, remove old). Do not attempt it in one PR.

**The database is a Docker named volume (`geocaching-data`, mounted at
`/data`).** `git pull && docker compose up -d --build` is the update path and
it **must not touch that volume.** Never add a step to compose, a Dockerfile,
or a migration that removes or resets `/data`. See §8.

## 7. Codebase quirks you must not "fix" naively

These look wrong but are intentional (or load-bearing). Changing them breaks
things.

- **`geocaches.gc_code` stores a full geolink URL, not a GC code.**
  `geocaches.go` sets `gc.Geolink = gc.GCCode` on read and stores
  `gc.Geolink` into the `gc_code` column on write (`geocaches.go:89, 137,
  166`). The column name is historical. Don't rename it and don't "correct"
  it to store a short code — frontend and DB both depend on this overload.
- **Two API-helper paths in the frontend.** `services/fetcher.js` is the
  original; many admin views define their own inline `apiRequest()`. They
  differ in error handling (see §5.1). When editing an admin view, check
  which path it uses and preserve that behavior unless you are deliberately
  unifying them.
- **Response envelope is inconsistent across endpoints.** Some return a bare
  array (`GetAdminMessages`), some `{data, current_page, last_page, total}`
  (`GetContactSubmissions`), some `{status, data, total}` (`GetUsers`).
  Frontend code uses fallbacks like `data.data || data || []`. When adding an
  endpoint, pick the existing convention for that resource and update both
  sides; don't invent a third envelope.
- **`Event.UUID` backfill in `UpdateEvent`** (`events.go:358-369`): older rows
  may have empty UUID; update generates and writes one if missing. Keep this
  logic — public event links depend on UUID existing.
- **`getPathParam` + `r.PathValue` fallback** (`static.go:206, 425`): chi v5
  stores URL params in context. The fallback exists for safety. Prefer
  `chi.URLParam(r, "name")` for new code (that's what `contact.go`,
  `events.go`, `geocaches.go`, `golden_key_months.go` use).
- **`AdminService.js` is only partially used.** Several admin views bypass it
  with their own fetch logic. If you add a service function, also wire the
  view to use it, or you've added dead code.

## 8. Production data & deployment safety (CRITICAL)

- The production DB volume is the only copy of real data (contacts, events,
  golden-key state, admin users). **There is no automated backup.** Treat
  every migration and every DB-touching change as if a mistake is
  irreversible, because it nearly is.
- The deploy flow is `git pull && docker compose up -d --build`. This must
  remain safe. Do not add `volumes:` that bind-mount over `/data`, do not
  add init scripts that wipe the DB, do not change `DATABASE_PATH` semantics.
- `JWT_SECRET` must be ≥32 chars in production; `config.Validate()` refuses
  to start otherwise (`config.go:89-98`). Never lower this bar. Never commit
  a real secret — `.env` is gitignored.
- The seed admin email is the literal string `admin`. Real admin users
  created via the panel have proper emails. Don't "validate email format" on
  login in a way that blocks the seed account.
- Default admin password is printed **once** to stdout at first boot and
  never again. Don't add logging that repeats it, and don't remove the
  one-time print.
- SMTP is optional on purpose. Never make email-sending a hard dependency
  for saving a contact submission.

## 9. Security guardrails

- **SQL:** every query uses `?` placeholders. Keep it that way. Never
  `fmt.Sprintf` user input into SQL. The only string interpolation in SQL is
  for static query assembly (e.g. appending ` AND state = 'published'` based
  on a boolean, never based on user input).
- **Auth:** all admin routes are under `r.Route("/admin", …)` with
  `middleware.JWTAuth` + `NoCache` (`router.go:77-79`). Any new admin
  endpoint **must** live inside that group. Public endpoints must not.
- **JWT parsing** checks the signing method is HMAC (`auth.go:37`) to prevent
  the `none`/RS256 confusion attack. Preserve that check.
- **Image uploads** (`handlers.go:45`) validate extension **and** magic
  bytes, cap at 10MB, generate a UUID filename, and `ServeImage` runs
  `filepath.Base` to prevent traversal. Keep all of these.
- **HTML escaping in emails:** `email.go` escapes user-supplied fields with
  `html.EscapeString` before embedding in the HTML body. The contact form
  submission email is the main XSS-to-admin-inbox vector. Preserve this.
- **CORS:** `cfg.CORSOrigins` from `CORS_ORIGINS` env, comma-separated.
  `AllowCredentials: true`. Never set `AllowedOrigins: ["*"]` with credentials
  — that's an insecure combo browsers reject anyway.
- **Rate limiting:** login = 5 attempts / 15 min per IP
  (`router.go:40, 74`). Contact form = throttled. Nginx also rate-limits.
  Don't remove these to "make testing easier."
- **Secrets:** `.env` is gitignored. `.env.example` contains only placeholders.
  Never write a real password, SMTP pass, or JWT secret into a committed file.
- **CSP** is set in `nginx.conf` with `unsafe-inline` for scripts/styles
  (required by Vite's module preload + inline styles). Don't add inline event
  handlers or `eval`-style code that would force loosening CSP further.

## 10. Code conventions

**Go backend**
- Package layout: `internal/{config,database,handlers,middleware,router,services}`.
  One responsibility per package.
- Handlers hang off a shared `Handler` struct (`handlers.go:29`) holding `db`,
  `cfg`, `emailService`. New handlers follow the same receiver pattern.
- JSON responses go through `respondJSON(w, status, data)` (`auth.go:285`).
  Don't hand-roll `w.WriteHeader` + `json.Encode` in new code.
- Errors from `db.QueryRow(...).Scan(...)` must handle `sql.ErrNoRows`
  explicitly (404) and other errors (500). Several handlers do this; mirror it.
- Always `defer rows.Close()`. Always check `rows.Err()` after iteration in
  new code (existing code mostly doesn't — don't bulk-rewrite, but don't
  repeat the omission).
- `gofmt` formatting is enforced by convention. Run `gofmt -w .` before
  committing. No comments unless explaining *why* (the repo style is
  comment-light).

**Vue frontend**
- `<script setup>` Composition API only. No Options API in new code.
- Admin views wrap themselves in `<AdminLayout pageTitle="…">` and use a
  default `<slot>` plus an optional `#actions` slot. Follow that.
- API base is `config.apiUrl` (= `import.meta.env.VITE_API_URL`). Admin
  endpoints are prefixed `admin/...`, public endpoints are not.
- Tailwind is **not** used. Styling is plain CSS with CSS variables defined
  in `css/admin.css` (admin) and `css/main.css`/`base.css` (public). Reuse
  existing variables (`--admin-primary`, `--admin-surface`, etc.) — don't
  hardcode hex values that already exist as vars.
- Routing is flat (each admin view imports `AdminLayout` itself; there is no
  nested `<router-view>` in the layout). Preserve this unless doing a
  deliberate refactor.
- Don't add emojis to code or UI unless explicitly asked.

## 11. Where an LLM/agent tends to mess up here — specific guardrails

1. **Editing only one side of a frontend/backend contract.** This is how
   §5.1 happened. When you change a route, field name, or response shape,
   grep the other side and update both. Files to cross-check:
   `router.go` ↔ `services/*.js` and each view's `apiRequest`.
2. **Trusting `response.ok` is checked when it isn't.** Several inline
   `apiRequest` helpers only check 401. When adding/fixing a call, verify the
   response status is actually handled, or you'll reintroduce a silent
   no-op.
3. **"Cleaning up" the migration file.** Do not. Read §6. If you must add a
   migration, follow every rule there, especially the ALTER-name swallow-list.
4. **Treating `fetchFromServer` `[]` as "no data" vs "error".** It returns
   `[]` on network failure and `{access_denied: true}` on 401. Handle both.
5. **Renaming `gc_code` / dropping the geolink overload.** See §7. It's
   intentional.
6. **Adding a `DROP TABLE` or `DELETE FROM` migration** "to reset for
   testing." This wipes production. Never.
7. **Lowering security bars** (JWT length, CORS, rate limits, magic-byte
   validation, HTML escaping in emails) to make a feature "work faster."
   Don't.
8. **Inventing a new response envelope.** Reuse the existing one for that
   resource (see §7).
9. **Committing secrets / writing a real JWT or SMTP password into .env.**
   `.env` is gitignored for a reason. Leave it.
10. **Bulk-rewriting working code** to "improve" style while fixing a bug.
    Make the minimal change for the bug. Large refactors need explicit
    approval.
11. **Adding comments everywhere.** The repo is comment-light by choice.
    Only comment *why*, when non-obvious.
12. **Assuming tests exist.** They don't. Your safety net is
    `go build ./... && go vet ./... && gofmt -l .` and `npm run build`. Run
    them.
13. **Using `100vh` for mobile scroll containers.** Use `100dvh` or a
    flex-based min-height-0 layout. See §5.3.
14. **Forgetting that `AdminLayout` runs its own auth check on every admin
    route.** Don't add a second auth gate inside a child view; don't remove
    the one in `AdminLayout`.

## 12. Before you say "done"

- [ ] `go build ./...` passes from `backend/`
- [ ] `go vet ./...` is clean
- [ ] `gofmt -l .` prints nothing
- [ ] `npm run build` passes from `frontend/`
- [ ] If you touched a route or response shape, you updated **both** the
      Go handler/router and every frontend caller
- [ ] If you added an `ALTER TABLE` migration, you added its name to the
      swallow-list in `migrations.go:295`
- [ ] No new `fmt.Sprintf` into SQL with user input
- [ ] No secrets committed
- [ ] No production data was wiped or reset by your change
- [ ] You did not lower any security bar in §9
- [ ] If you fixed a bug, you reproduced the original symptom first and
      confirmed the fix removes it (don't claim a fix from reading alone)
