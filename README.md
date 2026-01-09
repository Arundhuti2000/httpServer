# httpServer

A Go-based HTTP server for a social media platform featuring user authentication, chirp (short message) management, and PostgreSQL database integration. The server provides RESTful API endpoints for creating users, managing chirps (tweets-like messages), user login, and administrative functions.

## Overview

This server implements a Twitter-like application backend with the following key features:

- **User Management**: Create user accounts with secure password hashing using Argon2id
- **Authentication**: User login system with password verification
- **Chirp System**: Create and retrieve short messages (chirps) with a 140-character limit
- **Profanity Filter**: Automatic content moderation for chirps
- **Database Integration**: PostgreSQL database with SQLC for type-safe queries
- **Static File Serving**: Serves frontend files from the project root under `/app/`
- **Admin Dashboard**: Metrics tracking and system reset capabilities
- Default port: `8080`

## Prerequisites

- Go (1.18+ recommended)
- PostgreSQL database
- Environment variables configured (see below)

## Configuration

Create a `.env` file in the project root with:

```env
DB_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
PLATFORM=dev
```

## Build & Run

From the project root:

```bash
# Run without building
go run .

# Build and run
go build -o httpserver .
./httpserver    # Unix
.\httpserver.exe # Windows
```

The server listens on port `8080` by default.

## Endpoints

### Static Files

- `GET /app/<path>` — Serves static files from the repository root. Example: `/app/index.html`

### API Endpoints

#### Health & Monitoring

- `GET /api/healthz` — Health check endpoint (returns 200 OK)

#### User Management

- `POST /api/users` — Create a new user account
  - Request body: `{"email": "user@example.com", "password": "securepassword"}`
  - Returns user details (ID, email, created_at, updated_at)
- `POST /api/login` — Authenticate user and login
  - Request body: `{"email": "user@example.com", "password": "password"}`
  - Returns user details upon successful authentication

#### Chirp Management

- `POST /api/chirps` — Create a new chirp (max 140 characters)
  - Request body: `{"body": "Your chirp message here", "user_id": "uuid"}`
  - Automatically filters profanity
  - Returns created chirp details
- `GET /api/chirps` — Retrieve all chirps
  - Returns array of all chirps
- `GET /api/chirps/{chirpID}` — Get a specific chirp by ID
  - Returns single chirp details

### Admin Endpoints

- `GET /admin/metrics` — View server metrics (file server hit count)
- `POST /admin/reset` — Reset all metrics and potentially clear database (dev/platform dependent)

## Examples

Create a user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "mypassword"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "mypassword"}'
```

Create a chirp:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -d '{"body": "This is my first chirp!", "user_id": "your-user-uuid"}'
```

Get all chirps:

```bash
curl http://localhost:8080/api/chirps
```

Get specific chirp:

```bash
curl http://localhost:8080/api/chirps/{chirp-uuid}
```

Fetch index:

```bash
curl http://localhost:8080/app/index.html
```

Check health:

```bash
curl http://localhost:8080/api/healthz
```

View metrics:

```bash
curl http://localhost:8080/admin/metrics
```

Reset metrics:

```bash
curl -X POST http://localhost:8080/admin/reset
```

## Notes

- The server maps `/app/` to the current working directory (`.`) so any files (for eg `index.html` or the `assets/` folder) will be served there.
- Passwords are hashed using Argon2id for secure storage
- Chirps have a 140-character limit and are automatically filtered for profanity
- The profanity filter replaces words like "kerfuffle", "sharbert", and "fornax" with "****"
- All database operations use SQLC-generated type-safe queries
- Users and chirps are stored in PostgreSQL with proper foreign key relationships
- If you want a different port or file root, update `main.go` constants.

## Database Schema

The application uses two main tables:

- **users**: Stores user accounts with email and hashed passwords
- **chirps**: Stores user-generated chirps with references to user accounts

## License

This project has no license; add one if you plan to publish.
