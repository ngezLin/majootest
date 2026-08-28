# Test Case 2 — PHP Multi-Tenant Revenue Reporting Engine

A Laravel 11 REST API implementing JWT authentication, strict multi-tenancy enforcement, and daily revenue reporting with zero-fill for dates with no transactions.

---

## Stack

| Component | Technology |
|---|---|
| Framework | Laravel 11 |
| PHP | 8.3 |
| Authentication | `tymon/jwt-auth` v2 — JWT Bearer tokens |
| Database | MySQL 8.0 (provided schema) |
| Container | Docker (PHP-FPM + Nginx via Supervisord) |

---

## API Endpoints

| Method | Endpoint | Auth | Description |
|---|---|:---:|---|
| `POST` | `/api/auth/login` | ❌ | Login and receive JWT token |
| `POST` | `/api/auth/logout` | ✅ | Invalidate JWT token |
| `GET` | `/api/merchant/report` | ✅ | Monthly revenue report for authenticated merchant |
| `GET` | `/api/outlet/{id}/report` | ✅ | Monthly revenue report for a specific outlet |

### Query Parameters (Report Endpoints)

| Param | Type | Default | Description |
|---|---|---|---|
| `month` | int 1–12 | 11 (merchant) / 8 (outlet) | Report month |
| `year` | int | current year | Report year |
| `page` | int | 1 | Page number |
| `per_page` | int (max 31) | 10 | Days per page |

---

## Architecture

```
HTTP Request
     │
     ▼
[ JWT Middleware (auth:api) ]         — Verifies Bearer token
     │
     ▼
[ MultiTenancyMiddleware ]             — Resolves merchant_id from JWT user_id
     │                                   Injects it into request attributes
     ▼
[ Controller ]                         — Validates query params
     │
     ▼
[ RevenueReportService ]               — Queries DB, zero-fills missing dates, paginates
     │
     ▼
[ MySQL (Transactions table) ]         — Indexed on (merchant_id, created_at) and (outlet_id, created_at)
```

### Multi-Tenancy Implementation

The `MultiTenancyMiddleware` enforces strict data isolation:

1. Extracts `user_id` from the JWT payload.
2. Looks up `Merchants WHERE user_id = ?` to resolve `merchant_id`.
3. Injects `merchant_id` into `$request->attributes` — controllers **never trust URL/body input** for scoping.
4. For outlet reports: verifies `Outlets WHERE id = ? AND merchant_id = <from token>` before querying.

If the outlet doesn't belong to the authenticated merchant → `403 Forbidden`.

### Zero-Fill Revenue Logic

Revenue (Omzet) = `SUM(bill_total)` from `Transactions`.

For dates with no transactions, the report shows `0` instead of omitting the date.

**Implementation:**
1. Query DB: `SELECT DATE(created_at), SUM(bill_total) GROUP BY DATE(created_at)` for the given month.
2. Generate full month date range in PHP using `CarbonPeriod`.
3. Merge: for each date, use DB value or `0` (zero-fill).
4. Slice result with `array_slice` for pagination.

---

## Indexing Strategy (Optimization)

| Index | Table | Columns | Purpose |
|---|---|---|---|
| `idx_transactions_merchant_date` | Transactions | `(merchant_id, created_at)` | Merchant revenue queries — avoids full table scan |
| `idx_transactions_outlet_date` | Transactions | `(outlet_id, created_at)` | Outlet revenue queries |
| `idx_outlets_merchant` | Outlets | `(merchant_id)` | Multi-tenancy outlet ownership check |
| `idx_merchants_user` | Merchants | `(user_id)` | JWT → merchant_id resolution on every request |

---

## Getting Started

### Prerequisites
- Docker + Docker Compose installed

### 1. Setup

```bash
cd D:\projects\majootest

# Copy env
cp case2-php/.env.example case2-php/.env
```

### 2. Start containers

```bash
docker compose up -d
```

This will:
- Start MySQL on port `3307` and load the schema + seed data automatically.
- Build and start the PHP app on `http://localhost:8000`.

### 3. Run migrations (Users table + indexes)

```bash
docker exec majoo_php_app php artisan migrate
```

### 4. Generate JWT secret

```bash
docker exec majoo_php_app php artisan jwt:secret
```

---

## Testing the API

### Login (Merchant 1)

```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"merchant1@mail.com","password":"password123"}'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJ...",
    "token_type": "Bearer",
    "user": { "merchant_id": 1 }
  }
}
```

### Merchant Revenue Report (August)

```bash
curl "http://localhost:8000/api/merchant/report?month=8&year=2026&page=1&per_page=10" \
  -H "Authorization: Bearer eyJ..."
```

### Outlet Report (Outlet 1, August)

```bash
curl "http://localhost:8000/api/outlet/1/report?month=8&year=2026&page=1&per_page=10" \
  -H "Authorization: Bearer eyJ..."
```

### Cross-Tenant Access (Merchant 1 trying Outlet 2 of Merchant 2)

```bash
curl "http://localhost:8000/api/outlet/2/report?month=8&year=2026" \
  -H "Authorization: Bearer eyJ..."
# → 403 Forbidden
```

---

## Known Limitations & Future Improvements

- **No token refresh endpoint** — for simplicity; production would add `POST /auth/refresh`.
- **Zero-fill is done in PHP** — for a very large month range or many merchants, a MySQL recursive CTE date generator would be more efficient.
- **No rate limiting** — production would add `throttle:api` middleware.
- **SQLite fallback removed** — app targets MySQL only, matching the assessment schema.
