# Chirpy

A Twitter-like social media API built with Go, featuring user authentication, chirp (tweet) management, and webhooks.

## Features

- 🔐 **User Authentication**: JWT-based authentication with refresh tokens
- 📝 **Chirp Management**: Create, read, and delete chirps (short messages)
- 👥 **User Management**: User registration, login, and profile updates
- 🔍 **Advanced Queries**: Filter chirps by author and sort by date
- 🎣 **Webhooks**: Webhook integration for external services
- 📊 **Metrics**: Built-in metrics and monitoring endpoints
- 🗄️ **Database**: PostgreSQL integration with SQLC for type-safe queries

## API Endpoints

### Authentication
- `POST /api/users` - Register a new user
- `POST /api/login` - User login
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token
- `PUT /api/users` - Update user profile

### Chirps
- `GET /api/chirps` - Get all chirps (supports filtering and sorting)
  - Query parameters:
    - `author_id` - Filter by author UUID
    - `sort` - Sort order (`asc` or `desc`, default: `asc`)
- `POST /api/chirps` - Create a new chirp
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp

### Webhooks
- `POST /api/polka/webhooks` - Polka webhook endpoint

### Admin & Monitoring
- `GET /api/healthz` - Health check endpoint
- `GET /admin/metrics` - View metrics dashboard
- `POST /admin/reset` - Reset application state

### Static Files
- `/app/*` - Serve static files

## Technology Stack

- **Language**: Go 1.25
- **Database**: PostgreSQL
- **Query Builder**: [SQLC](https://sqlc.dev/)
- **Authentication**: JWT tokens with Argon2id password hashing
- **HTTP Router**: Go standard library `net/http`
- **Database Driver**: `lib/pq`

## Dependencies

```go
github.com/alexedwards/argon2id v1.0.0    // Password hashing
github.com/golang-jwt/jwt/v5 v5.3.0       // JWT tokens
github.com/google/uuid v1.6.0             // UUID generation
github.com/joho/godotenv v1.5.1           // Environment variables
github.com/lib/pq v1.10.9                 // PostgreSQL driver
```

## Setup

### Prerequisites
- Go 1.25+
- PostgreSQL
- [SQLC](https://sqlc.dev/) (for database code generation)

### Environment Variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-secret-key-here
POLKA_KEY=your-polka-api-key
```

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/swokamoto/chirpy.git
   cd chirpy
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   ```bash
   # Apply database migrations from sql/schema/
   psql $DB_URL -f sql/schema/001_users.sql
   psql $DB_URL -f sql/schema/002_chirps.sql
   psql $DB_URL -f sql/schema/003_password.sql
   psql $DB_URL -f sql/schema/004_refresh_tokens.sql
   psql $DB_URL -f sql/schema/005_is_red.sql
   ```

4. Generate database code (if modified):
   ```bash
   sqlc generate
   ```

5. Build and run:
   ```bash
   go build -o chirpy
   ./chirpy
   ```

The server will start on `http://localhost:8080`

## Usage Examples

### Register a user
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

### Create a chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"body": "Hello, world! This is my first chirp!"}'
```

### Get chirps with filtering and sorting
```bash
# Get all chirps, sorted by newest first
curl "http://localhost:8080/api/chirps?sort=desc"

# Get chirps from specific author
curl "http://localhost:8080/api/chirps?author_id=USER_UUID"

# Combine filters
curl "http://localhost:8080/api/chirps?author_id=USER_UUID&sort=desc"
```

## Project Structure

```
chirpy/
├── sql/
│   ├── queries/          # SQLC query definitions
│   └── schema/           # Database migrations
├── internal/
│   ├── auth/             # Authentication utilities
│   └── database/         # Generated SQLC code
├── handler_*.go          # HTTP handlers
├── main.go               # Application entry point
├── json.go               # JSON utilities
├── metrics.go            # Metrics middleware
├── readiness.go          # Health check handler
├── reset.go              # Reset functionality
└── assets/               # Static assets
```

## Development

### Database Changes

1. Add migration to `sql/schema/`
2. Update queries in `sql/queries/`
3. Regenerate Go code: `sqlc generate`

### Adding New Endpoints

1. Create handler function in appropriate `handler_*.go` file
2. Register route in `main.go`
3. Add tests and documentation

## License

This project was created as part of the Boot.dev backend development course.