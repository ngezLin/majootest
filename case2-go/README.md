# Test Case 2: Blog REST API (Go)

A production-ready RESTful API for a blog platform built in Go using the **Controller-Service-Repository (3-Tier)** architecture.

---

## Features

- **User Authentication & Authorization**:
  - Registration & Login with `bcrypt` (salted password hashing).
  - JWT Bearer token generation and middleware verification.
  - Resource ownership checks (only authors can modify/delete their own posts and comments).
- **CRUD for Posts & Comments**:
  - Paginated list queries with search filters.
  - Granular CRUD endpoints for blog posts and nested/associated comments.
- **Database Transactions (ACID)**:
  - Atomic comment creation coupled with real-time `comment_count` increments on posts.
  - Atomic post deletion with cascade cleanup of comments inside a single transactional boundary (`WithTransaction`).
- **Robust Error Handling & Input Validation**:
  - Request body validation with descriptive error details (`go-playground/validator/v10`).
  - Standardized JSON responses (success & error envelopes).
  - Panic recovery middleware preventing server crashes.
- **Observability & Resilience**:
  - Structured request logging with status codes, latency, and client IP.
  - Graceful shutdown handling `SIGINT` / `SIGTERM` with context timeout.

---

## Architecture Overview

```text
HTTP Request
     │
     ▼
[ Middlewares: Logger ──► Recovery ──► CORS ──► JWT Auth ]
     │
     ▼
[ Controller Layer ]      <-- Request JSON decoding, validation, response envelope
     │
     ▼
[ Service Layer ]         <-- Business logic, ownership rules, transaction orchestration
     │
     ▼
[ Repository Layer ]      <-- Raw SQL queries, database connection pool, TxManager
     │
     ▼
[ MySQL Database ]
```

### Folder Layout

```text
case2-go/
├── cmd/
│   └── main.go                     # Application bootstrap, routing & graceful shutdown
├── internal/
│   ├── config/                     # Environment configuration loader
│   ├── models/                     # Data models, DTOs & response envelopes
│   ├── controllers/                # HTTP handlers (AuthController, PostController, CommentController)
│   ├── services/                   # Business logic (AuthService, PostService, CommentService)
│   ├── repositories/               # Data access & ACID Transaction Manager
│   ├── middlewares/                # Auth, Logger, Recovery, CORS middlewares
│   └── utils/                      # Password hash, JWT token, Validator, Response helpers
├── migrations/                     # Database schema migration scripts (up/down)
├── docs/                           # OpenAPI / Swagger 3.0 specification
├── Dockerfile                      # Multi-stage production container build
├── .env.example                    # Environment template
└── README.md
```

---

## Technology Choices Justification

| Technology | Justification |
| :--- | :--- |
| **Gin Framework (`gin-gonic/gin`)** | High-performance, widely-adopted HTTP framework in the Go ecosystem providing idiomatic JSON binding, route grouping, and robust middleware support. |
| **Raw SQL (`database/sql`)** | Demonstrates clear command over SQL, transactions (`BeginTx`, `Commit`, `Rollback`), prepared statements, indexing, and connection pool optimization without ORM overhead. |
| **`golang-jwt/jwt/v5`** | Industry-standard JWT library for secure token generation and cryptographic HMAC-SHA256 signature verification. |
| **`golang.org/x/crypto/bcrypt`** | Secure adaptive password hashing with automatic per-user salt generation to resist rainbow table attacks. |
| **`validator/v10`** | Declarative tag-based struct validation with detailed error field reporting. |

---

## Getting Started

### Prerequisites
- Go 1.22+ installed
- MySQL 8.0+ (or Docker)

### 1. Environment Setup
Copy `.env.example` to `.env` and adjust database credentials:
```bash
cp .env.example .env
```

### 2. Database Migrations
Execute the migration scripts located in `migrations/` on your MySQL instance:
```bash
mysql -u root -p blog_db < migrations/000001_create_users_table.up.sql
mysql -u root -p blog_db < migrations/000002_create_posts_table.up.sql
mysql -u root -p blog_db < migrations/000003_create_comments_table.up.sql
```

### 3. Run Locally
```bash
go run ./cmd
# or: go run ./cmd/main.go
```
The server will start listening on `http://localhost:8080`.

---

## Running with Docker

Run the entire service inside Docker:
```bash
# Build Docker image
docker build -t blog-api:latest .

# Run container
docker run -p 8080:8080 --env-file .env blog-api:latest
```

---

## Running Automated Tests & Coverage

Execute the complete test suite (unit tests, integration tests, mock repositories):
```bash
# Run all tests with verbose output
go test ./... -v

# Run with test coverage
go test ./... -cover
```

---

## API Endpoints Summary

Refer to [`docs/openapi.yaml`](./docs/openapi.yaml) for the full OpenAPI specification.

### Authentication
- `POST /api/v1/auth/register` — Register a new user
- `POST /api/v1/auth/login` — Login and receive JWT access token
- `GET  /api/v1/auth/me` — Get current user profile *(Protected)*

### Posts
- `GET    /api/v1/posts` — List posts with pagination (`?page=1&limit=10&search=keyword`)
- `GET    /api/v1/posts/:id` — Get single post by ID
- `POST   /api/v1/posts` — Create a new post *(Protected)*
- `PUT    /api/v1/posts/:id` — Update post *(Protected, Owner only)*
- `DELETE /api/v1/posts/:id` — Delete post and comments atomically *(Protected, Owner only)*

### Comments
- `GET    /api/v1/posts/:postId/comments` — List comments for a post
- `GET    /api/v1/posts/:postId/comments/:commentId` — Get single comment
- `POST   /api/v1/posts/:postId/comments` — Add comment *(Protected, atomic counter increment)*
- `PUT    /api/v1/posts/:postId/comments/:commentId` — Update comment *(Protected, Owner only)*
- `DELETE /api/v1/posts/:postId/comments/:commentId` — Delete comment *(Protected, Owner only, atomic counter decrement)*

---

## Standard JSON Envelopes

### Success Envelope
```json
{
  "success": true,
  "message": "Post created successfully",
  "data": {
    "id": 1,
    "user_id": 42,
    "title": "Understanding Concurrency in Go",
    "content": "Goroutines make concurrency simple...",
    "comment_count": 0,
    "created_at": "2026-08-27T15:30:00Z"
  }
}
```

### Error Envelope
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "Password",
        "message": "Field 'Password' must be at least 8 characters long"
      }
    ]
  }
}
```

---

## Known Limitations & Future Roadmap
1. **Refresh Token Rotation**: Currently relies on short-to-medium lifespan access tokens; adding Redis-backed refresh token rotation would allow seamless token refresh.
2. **Caching Layer**: Integrating Redis caching for high-traffic post lookups (`GET /posts/:id`).
3. **Full-Text Search**: Implementing Elasticsearch or MySQL Full-Text indexes for larger blog datasets.
