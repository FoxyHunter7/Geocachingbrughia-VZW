# 🌲 Geocaching Brughia VZW

Website for Geocaching Brughia, a non-profit organization promoting geocaching in the Bruges region.

## 📁 Project Structure

```
├── frontend/          # Vue 3 SPA (public site + admin panel)
├── backend/           # Go API server with SQLite
├── backend-laravel/   # Legacy Laravel backend (reference only)
├── docker/            # Docker configuration
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   ├── docker-compose.yml
│   └── nginx.conf
└── data/              # SQLite database (mounted volume)
```

## 🚀 Quick Start

### Development

**Backend (Go):**
```bash
cd backend
cp .env.example .env  # Edit with your settings
go mod download
go run .
```

**Frontend (Vue):**
```bash
cd frontend
npm install
npm run dev
```

### Production (Docker)

```bash
cd docker
cp .env.example .env  # Edit with your settings

# Set up data directory permissions
mkdir -p ../data
sudo chown -R 65532:65532 ../data

# Build and run
docker compose up -d --build
```

## 🏗️ Architecture

### Backend (Go)
- **Framework:** Chi router (lightweight, fast)
- **Database:** SQLite with WAL mode
- **Auth:** JWT tokens
- **Email:** gomail for SMTP notifications

### Frontend (Vue 3)
- **Framework:** Vue 3 with Composition API
- **Build:** Vite
- **Features:**
  - Multi-language support (NL, EN, FR, DE)
  - Public site with events, geocaches, contact form
  - Admin panel for content management

### Docker Stack
- **Backend:** Distroless container (rootless, secure)
- **Frontend:** Static files served by Nginx
- **Nginx:** Reverse proxy with caching

## 🔒 Security Features

- Rootless containers (UID 65532)
- Read-only filesystems where possible
- Rate limiting on API endpoints
- JWT authentication for admin
- CORS configuration
- Security headers (X-Frame-Options, CSP, etc.)

## 📬 Contact Form System

When someone submits the contact form:
1. Submission saved to database with status `new`
2. Notification email sent to configured admin email
3. Background job checks hourly for unanswered submissions
4. Reminder sent if submission is `new` for more than 3 days

## 🗄️ Caching Strategy

- **ETag-based caching** for static content
- **stale-while-revalidate** for best UX
- **Nginx proxy cache** for API responses
- Cache headers set by backend, respected by Nginx

```
Cache-Control: public, max-age=300, stale-while-revalidate=3600
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_PATH` | SQLite database path | `./data/geocaching.db` |
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_EXPIRY_HOURS` | Token expiry | `24` |
| `SMTP_HOST` | SMTP server | - |
| `SMTP_PORT` | SMTP port | `587` |
| `SMTP_USER` | SMTP username | - |
| `SMTP_PASS` | SMTP password | - |
| `SMTP_FROM` | From address | - |
| `NOTIFICATION_EMAIL` | Admin notification email | - |
| `REMINDER_DAYS` | Days before reminder | `3` |
| `CORS_ORIGINS` | Allowed origins (comma-separated) | - |

## 📊 API Endpoints

### Public
- `GET /api/languages` - Supported languages
- `GET /api/static` - Static content/translations
- `GET /api/events` - Published events
- `GET /api/home_events` - Homepage events
- `GET /api/messages` - Site messages
- `GET /api/geocaches` - Active geocaches
- `GET /api/socials` - Social media links
- `POST /api/contact` - Submit contact form

### Admin (requires JWT)
- `GET/POST/PUT/DELETE /api/admin/events`
- `GET/POST/PUT/DELETE /api/admin/geocaches`
- `GET/POST/PUT/DELETE /api/admin/messages`
- `GET/POST/PUT/DELETE /api/admin/languages`
- `GET/POST/PUT/DELETE /api/admin/static`
- `GET/POST/PUT/DELETE /api/admin/socials`
- `GET/PUT/DELETE /api/admin/contacts` - Contact submissions

## 🔄 Migration from Laravel

The `backend-laravel/` folder contains the original Laravel implementation for reference.
Data can be migrated using a simple script (TODO).

## 📝 License

All rights reserved - Geocaching Brughia VZW

## 🤝 Contributing

This is a private project for Geocaching Brughia VZW.
Developed out of goodwill to support the organisation.
