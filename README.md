# GeocachingBrughia-VZW

This repository contains the code for the GeocachingBrughia VZW website,
a non-profit organisation.
The project has been developed out of goodwill to support the organisation.

## Contact Form Submissions

1. Submission is saved to the database with status `new`.
2. A notification email is sent to `NOTIFICATION_EMAIL`.
3. A background job runs every hour and sends a reminder if any submission
   has been `new` for more than `REMINDER_DAYS` days.
4. Admin users can update status, add notes,
   and assign submissions via the admin panel.

## Configuration reference

| Variable             | Description                                 | Default                 |
| -------------------- | ------------------------------------------- | ----------------------- |
| `PORT`               | Backend listen port                         | `8080`                  |
| `DATABASE_PATH`      | Path to SQLite file                         | `./data/geocaching.db`  |
| `JWT_SECRET`         | JWT signing secret (required in production) | —                       |
| `JWT_EXPIRY_HOURS`   | Token lifetime in hours                     | `24`                    |
| `SMTP_HOST`          | SMTP server hostname                        | —                       |
| `SMTP_PORT`          | SMTP port                                   | `587`                   |
| `SMTP_USER`          | SMTP username                               | —                       |
| `SMTP_PASS`          | SMTP password                               | —                       |
| `SMTP_FROM`          | Sender address for outgoing mail            | —                       |
| `NOTIFICATION_EMAIL` | Where contact form notifications go         | —                       |
| `REMINDER_DAYS`      | Days before sending a follow-up reminder    | `3`                     |
| `CORS_ORIGINS`       | Comma-separated allowed origins             | `http://localhost:5173` |

SMTP is optional, if `SMTP_HOST` is not set the email service is disabled
and contact form submissions are still saved to the database,
just no emails are sent.

## Development

### Tech stack

- **Backend:** Go, Chi router, SQLite (WAL mode), JWT auth, gomail
- **Frontend:** Vue 3, Vite, Composition API
- **Infrastructure:** Docker, Nginx (reverse proxy + static file serving)
- **Containers:** Distroless backend image (UID 65532, no shell)

The production setup is a single docker compose:
A Go backend and Vue frontend behind Nginx.
The database is SQLite stored in a named Docker volume. It is setup this way
for simplicity, given the minimal amount of data stored.

### Project Structure

```txt
├── frontend/          - Vue 3 SPA (public site + admin panel)
├── backend/           - Go API server with SQLite
├── data/              - SQLite database (Docker named volume)
├── docker-compose.yml
├── Dockerfile.backend
├── Dockerfile.frontend
└── nginx.conf
```

### Local Deployment

**Backend:**

```bash
cd backend
cp .env.example .env   # fill in your values
go mod download
go run .
```

**Frontend** _(separate terminal)_:

```bash
cd frontend
npm install
npm run dev
```

The frontend dev server proxies `/api/` to `localhost:8080` automatically.
You don't need Nginx locally.

## Poduction

### First deploy

1. **Copy the repo to the server** (or `git clone`).

2. **Create a `.env` file** in the project root (next to `docker-compose.yml`):

   ```env
   JWT_SECRET=<generate with: openssl rand -hex 32>
   JWT_EXPIRY_HOURS=24

   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USER=user@example.com
   SMTP_PASS=yourpassword
   SMTP_FROM=noreply@geocachingbrughia.be
   NOTIFICATION_EMAIL=bestuur@geocachingbrughia.be

   REMINDER_DAYS=3
   CORS_ORIGINS=https://geocachingbrughia.be
   ```

   `JWT_SECRET` is required and must be at least 32 characters,
   the server refuses to start without it in production.

3. **Place SSL certificates** in `./certs/`
   _(the volume is already mounted in `docker-compose.yml`)_:

   ```txt
   certs/
   ├── fullchain.pem
   └── privkey.pem
   ```

   If you're using Let's Encrypt / Certbot, point it at this directory.
   Then uncomment the HTTPS server block in `nginx.conf`
   and the HTTP→HTTPS redirect.

4. **Build and start**:

   ```bash
   docker compose up -d --build
   ```

5. **Get the first-boot admin password** from the container logs:

   ```bash
   docker compose logs backend | grep -A 6 "DEFAULT ADMIN"
   ```

   The backend seeds an `admin` account on first run with a random password
   printed once to stdout.
   Log in at `/admin` and change the password immediately.
   The login identifier is literally `admin` (not an email address)
   for this seed account.

### Health Check

```bash
docker compose ps
curl http://localhost/health
```

### Update

To deploy a new version:

```bash
git pull

# Rebuild images and restart — the database volume is untouched
docker compose up -d --build
```

If you only changed frontend or backend independently
you can rebuild just that service:

```bash
docker compose up -d --build backend
docker compose up -d --build frontend
```

To see logs during/after an update:

```bash
docker compose logs -f
```

### Rolling back

Docker keeps the previous image layers around until pruned.
If something goes wrong after an update,
you can check out the previous commit and rebuild:

```bash
git checkout HEAD~1
docker compose up -d --build
```
