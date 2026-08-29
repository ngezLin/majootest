# Test Case 2 — PHP Multi-Tenant Revenue Reporting Engine

A Laravel 11 RESTful API implementing JWT authentication, strict multi-tenancy enforcement, and daily revenue reporting with zero-fill for dates with no transactions.

---

## 1. Stack & Technologies

| Component | Technology | Description |
|---|---|---|
| **Framework** | Laravel 11 | RESTful routing, controllers, Eloquent ORM |
| **Authentication** | `tymon/jwt-auth` | Stateless JWT Bearer tokens |
| **Database** | MySQL 8.0 | Relational database with composite indexing |
| **Infrastructure** | Docker | Containerized MySQL 8.0 service |

---

## 2. API Endpoints

| Method | Endpoint | Auth | Description |
|---|---|:---:|---|
| `POST` | `/api/auth/login` | ❌ | Authenticate via email & password, returns JWT token |
| `POST` | `/api/auth/logout` | ✅ | Invalidate current JWT token |
| `GET` | `/api/merchant/report` | ✅ | Monthly daily revenue report for authenticated merchant |
| `GET` | `/api/outlet/{id}/report` | ✅ | Monthly daily revenue report for a specific outlet |

### Query Parameters (Report Endpoints)

| Parameter | Type | Default | Description |
|---|---|---|---|
| `month` | int (1–12) | `11` (merchant) / `8` (outlet) | Target report month |
| `year` | int | Current Year | Target report year |
| `page` | int | `1` | Pagination page number |
| `per_page` | int (1–31) | `10` | Number of days per page |

---

## 3. Architecture & Business Logic

### A. Strict Multi-Tenancy Enforcement
Multi-tenancy is enforced in `App\Http\Middleware\MultiTenancyMiddleware`:
1. The middleware parses and verifies the Bearer JWT token on incoming requests.
2. It extracts `user_id` and looks up `Merchants WHERE user_id = ?` to resolve the tenant's `merchant_id`.
3. The resolved `merchant_id` is injected into request attributes.
4. Controllers **never trust user-supplied merchant IDs from the request URL or body**, ensuring strict isolation.
5. For outlet endpoints (`/api/outlet/{id}/report`), it checks that `Outlets.merchant_id == authenticated merchant_id`. If a user attempts to view another merchant's outlet, it returns **`403 Forbidden`**.

### B. Daily Revenue Calculation with Zero-Fill
Assessment Constraint: *Reports must include pagination and display "0" revenue for dates with no transactions.*

**Implementation (`App\Services\RevenueReportService`)**:
1. Queries the database using `SUM(bill_total)` grouped by `DATE(created_at)` for the requested month.
2. Generates the continuous calendar day sequence for the entire month using `CarbonPeriod`.
3. Merges the DB results with the full date range — dates without transactions are filled with `0.0`.
4. Applies pagination using `array_slice` and returns pagination metadata (`total_days`, `total_pages`, `current_page`, `per_page`).

---

## 4. Indexing Strategy & DML Justification

| Index | Target Table | Target Columns | Purpose & DML Justification |
|---|---|---|---|
| `idx_transactions_merchant_date` | `Transactions` | `(merchant_id, created_at)` | **Optimizes Merchant Reports**: Enables index seek on `merchant_id` combined with range scan on `created_at BETWEEN start AND end`, avoiding full table scans as transaction volume scales. |
| `idx_transactions_outlet_date` | `Transactions` | `(outlet_id, created_at)` | **Optimizes Outlet Reports**: Speeds up outlet-specific date filtering and aggregation. |
| `idx_outlets_merchant` | `Outlets` | `(merchant_id)` | **Optimizes Multi-Tenancy Checks**: Accelerates lookup when verifying outlet ownership in `OutletReportController`. |
| `idx_merchants_user` | `Merchants` | `(user_id)` | **Optimizes JWT Tenant Resolution**: Allows $O(\log N)$ lookup of `merchant_id` from authenticated `user_id` on every request. |

---

## 5. Local Setup & Running Guide

### Prerequisites
- Docker (for MySQL)
- PHP 8.2+ with `pdo_mysql` and `mbstring` extensions enabled
- Composer

### Step 1: Start MySQL via Docker
From the project root:
```powershell
docker compose up -d
```
*MySQL 8.0 will be running on `127.0.0.1:3306` with database `reporting_db`.*

### Step 2: Configure Environment
In `case2-php/`:
```powershell
cp .env.example .env
```
Ensure database settings match:
```ini
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=reporting_db
DB_USERNAME=root
DB_PASSWORD=root
```

### Step 3: Install Dependencies & Run Migrations
From `case2-php/`:
```powershell
composer install
php artisan key:generate
php artisan jwt:secret
php artisan migrate
```

### Step 4: Start Development Server
```powershell
php artisan serve
```
*Server runs on **`http://127.0.0.1:8000`**.*

---

## 6. API Verification Examples

### 1. Login (Merchant 1)
```powershell
$login = Invoke-RestMethod -Uri "http://127.0.0.1:8000/api/auth/login" -Method Post -ContentType "application/json" -Body '{"email":"merchant1@mail.com","password":"password123"}'
$token = $login.data.token
```

### 2. Merchant Report (August, showing zero-fills on Aug 6–31)
```powershell
Invoke-RestMethod -Uri "http://127.0.0.1:8000/api/merchant/report?month=8&year=2026&page=1&per_page=10" -Headers @{ Authorization = "Bearer $token" }
```

### 3. Outlet Report (Outlet 1)
```powershell
Invoke-RestMethod -Uri "http://127.0.0.1:8000/api/outlet/1/report?month=8&year=2026" -Headers @{ Authorization = "Bearer $token" }
```

### 4. Cross-Tenant Security Check
Merchant 1 trying to access Outlet 2 (owned by Merchant 2):
```powershell
Invoke-RestMethod -Uri "http://127.0.0.1:8000/api/outlet/2/report?month=8&year=2026" -Headers @{ Authorization = "Bearer $token" }
# → 403 Forbidden: "You do not have access to this outlet."
```
