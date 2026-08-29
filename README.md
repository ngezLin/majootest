# Majoo Technical Assessment — Complete Solutions

This repository contains the complete solutions for the **Majoo Technical Take-Home Assessment**, covering **Go Concurrency**, **REST API Development (Go & PHP)**, and **Social Media Database Schema Design & Optimization**.

---

## 📁 Repository Structure & Overview

```
majootest/
├── case1-concurrency/      # 🐹 Test Case 1: Go Concurrent CSV File Processor
│   ├── cmd/                # CLI entrypoint & synthetic data generator
│   ├── internal/           # Streaming reader, worker pool, progress tracker, aggregator
│   ├── tests/              # Unit tests and performance benchmarks (> 1M rec/sec)
│   └── README.md           # Documentation, concurrency model, and benchmarks
│
├── case2-go/               # 🐹 Test Case 2 (Go): Blog & Comments REST API
│   ├── cmd/                # Server entrypoint with graceful shutdown
│   ├── internal/           # Gin handlers, JWT auth, database TxManager transactions
│   ├── docs/               # OpenAPI 3.0 specification & interactive Swagger UI
│   └── README.md           # API documentation and endpoints guide
│
├── case2-php/              # 🐘 Test Case 2 (PHP): Multi-Tenant Revenue Reporting Engine
│   ├── app/                # Laravel 11 Controllers, MultiTenancyMiddleware, RevenueReportService
│   ├── database/           # Migrations, composite indexing, and schema seed
│   ├── public/             # OpenAPI 3.0 specification & interactive Swagger UI
│   └── README.md           # Multi-tenancy architecture, zero-fill omzet logic, and DML justifications
│
├── case3-database/         # 🗄️ Test Case 3: Social Media Database Design & Optimization
│   ├── schema.sql          # Complete 3NF normalized DDL schema
│   ├── seed.sql            # Realistic sample test data
│   ├── queries.sql         # Newsfeed generation, unread counts, and engagement queries
│   └── README.md           # Mermaid ERD diagram, indexing strategy, and Redis caching architecture
│
└── docker-compose.yml      # Unified 1-command Docker Compose stack for all services
```

---

## 🚀 Quick Start Guide

### 1. Run Case 2 (PHP & Go APIs) with Docker in 1 Command
From the project root directory:
```powershell
docker compose up -d
```

#### Interactive Documentation & Testing:
* **PHP Multi-Tenant Reporting API Swagger UI**: 👉 **[http://localhost:8000/docs](http://localhost:8000/docs)**
* **Go Blog REST API Swagger UI**: 👉 **[http://localhost:8080/docs](http://localhost:8080/docs)**

---

### 2. Run Case 1 (Go Concurrent CSV Processor)
```powershell
cd case1-concurrency

# 1. Generate sample CSV test data (50,000 rows across 5 files):
go run ./cmd/generator --files=5 --rows=10000

# 2. Run concurrent processor with 8 workers:
go run ./cmd --workers=8 --input=./data/sample

# 3. Run unit tests & benchmarks:
cd tests
go test -v .
go test -bench=BenchmarkWorkerPool -benchmem .
```

---

### 3. Run & Verify Case 3 (Social Media Database)
```powershell
# Execute schema, seed data, and complex queries directly on MySQL:
Get-Content case3-database/schema.sql | docker exec -i majoo_mysql mysql -uroot -proot
Get-Content case3-database/seed.sql | docker exec -i majoo_mysql mysql -uroot -proot
Get-Content case3-database/queries.sql | docker exec -i majoo_mysql mysql -uroot -proot
```

---

## 📊 Summary of Test Cases & Assessment Deliverables

| Case | Category | Language / Framework | Deliverables & Key Highlights |
|---|---|---|---|
| **Case 1** | Concurrent Processing | Go 1.24 | Worker Pool Pattern, streaming $O(1)$ memory reading, `sync/atomic` live progress bar, row-level error isolation, **> 1,000,000 records/sec** throughput. |
| **Case 2 (Go)** | REST API | Go / Gin / MySQL | JWT Authentication, CRUD for Posts & Comments, atomic database transactions via `TxManager`, interactive Swagger UI at `/docs`. |
| **Case 2 (PHP)** | REST API | PHP 8.4 / Laravel 11 | JWT Authentication, strict `MultiTenancyMiddleware` data isolation, Zero-Fill daily revenue engine (`SUM(bill_total)`), composite indexing, interactive Swagger UI at `/docs`. |
| **Case 3** | Database Design | MySQL / Redis | 3NF Normalized schema, Mermaid ERD, Keyset cursor newsfeed queries, composite indexing, Hybrid Fan-out & Redis caching strategy. |