# Geocaching Brughia Backend

Go API server for Geocaching Brughia website.

## Development

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run with hot reload (install air first: go install github.com/air-verse/air@latest)
air

# Or run directly
go run .
```

## Project Structure

```
├── main.go                 # Entry point
├── internal/
│   ├── config/            # Configuration loading
│   ├── database/          # SQLite connection & migrations
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Auth, caching middleware
│   ├── router/            # Route definitions
│   └── services/
│       └── email/         # Email notifications
```

## Building

```bash
# Build for Linux
CGO_ENABLED=1 GOOS=linux go build -o server .

# Build with optimizations
CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server .
```

## Database

SQLite database with WAL mode for better concurrency.
Migrations run automatically on startup.

### Tables
- `users` - Admin users
- `languages` - Supported languages
- `events` - Events with translations
- `event_translations` - Event descriptions per language
- `geocaches` - Geocache listings
- `messages` - Site announcements with translations
- `message_translations` - Message content per language
- `static_content` - UI translations
- `socials` - Social media links
- `contact_submissions` - Contact form submissions
- `contact_notes` - Internal notes on submissions
