# Chirpy

A RESTful HTTP server for a Twitter-like social media platform built with Go and PostgreSQL.

## What It Does

Chirpy is a lightweight social media API that allows users to:

- **Create and manage user accounts** with secure password hashing (Argon2id)
- **Post short messages** (chirps) up to 140 characters with automatic profanity filtering
- **Authenticate securely** using JWT tokens with refresh token support
- **Upgrade to premium** (Chirpy Red) through webhook integration with Polka payment processor
- **Browse and filter chirps** by author and sort order (ascending/descending)
- **Manage personal chirps** with full CRUD operations

## Why You Should Care

This project demonstrates:

- **RESTful API Design**: Clean, organized HTTP endpoints following REST principles
- **Database Integration**: Using PostgreSQL with SQLC for type-safe queries
- **Authentication & Authorization**: JWT-based auth with access and refresh tokens
- **Security Best Practices**: Argon2id password hashing, API key validation, token-based auth
- **Webhook Integration**: Handling third-party payment webhooks with API key verification
- **Modern Go Practices**: Standard library HTTP server, context usage, proper error handling

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **Auth**: JWT (golang-jwt/jwt), Argon2id password hashing
- **Database Migrations**: SQL schema files
- **Query Generation**: SQLC for type-safe SQL
- **Environment Config**: godotenv

## Prerequisites

- Go (1.18+ recommended)
- PostgreSQL database
- Environment variables configured (see below)

## Installation & Setup

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- SQLC (for regenerating database queries)

### Steps

1. **Clone the repository**

   ```bash
   git clone <your-repo-url>
   cd httpServer
   ```

2. **Set up PostgreSQL database**

   ```bash
   createdb chirpy
   ```

3. **Run database migrations**

   ```bash
   psql -d chirpy -f sql/schema/001_users.sql
   psql -d chirpy -f sql/schema/002_chirps.sql
   psql -d chirpy -f sql/schema/003_users_hashed_password.sql
   psql -d chirpy -f sql/schema/004_refresh_tokens.sql
   psql -d chirpy -f sql/schema/005_chirpy_red.sql
   ```

4. **Configure environment variables**

   Create a `.env` file in the root directory:

   ```env
   DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
   PLATFORM="dev"
   JWT_SECRET="your-secret-key-here"
   POLKA_KEY="your-polka-api-key"
   ```

5. **Install dependencies**

   ```bash
   go mod download
   ```

6. **Run the server**

   ```bash
   go build -o httpserver && ./httpserver
   ```

   The server will start on `http://localhost:8080`

## API Endpoints

### Health & Metrics

- `GET /api/healthz` - Health check
- `GET /admin/metrics` - Server metrics
- `POST /admin/reset` - Reset database (dev only)

### Users

- `POST /api/users` - Create a new user
- `PUT /api/users` - Update user information (authenticated)
- `POST /api/login` - Login and receive JWT tokens
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Chirps

- `POST /api/chirps` - Create a new chirp (authenticated)
- `GET /api/chirps` - Get all chirps (supports `?author_id=<uuid>` and `?sort=asc|desc`)
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete your own chirp (authenticated)

### Webhooks

- `POST /api/polka/webhooks` - Polka payment webhook (requires API key)

### Static Files

- `GET /app/<path>` - Serves static files from the project root

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
- The profanity filter replaces words like "kerfuffle", "sharbert", and "fornax" with "\*\*\*\*"
- All database operations use SQLC-generated type-safe queries
- Users and chirps are stored in PostgreSQL with proper foreign key relationships
- If you want a different port or file root, update `main.go` constants.

## Database Schema

Project Structure

```
httpServer/
├── internal/
│   ├── auth/          # Authentication utilities (JWT, password hashing, API keys)
│   └── database/      # SQLC generated database code
├── sql/
│   ├── queries/       # SQL queries for SQLC
│   └── schema/        # Database migration files
├── assets/            # Static assets
├── main.go            # Server entry point
├── *.go               # HTTP handlers
├── .env               # Environment configuration
└── go.mod             # Go module dependencies
```

## Development

To regenerate database queries after modifying SQL files:

```bash
sqlc generate
```

## Database Schema

The application uses the following main tables:

- **users**: User accounts with email, hashed passwords, and Chirpy Red status
- **chirps**: User-generated chirps with references to user accounts
- **refresh_tokens**: JWT refresh tokens for authentication

## Notes

- The server maps `/app/` to the current working directory (`.`)
- Passwords are hashed using Argon2id for secure storage
- Chirps have a 140-character limit and are automatically filtered for profanity
- All database operations use SQLC-generated type-safe queries
- JWT tokens expire after 1 hour; refresh tokens are valid for 60 days
- Polka webhooks require API key authentication

## License

This project was built as part of a learning exercise
