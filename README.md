# GameVault Store: Gaming Purchase Processing System

## Overview

Build a simple gaming store purchase processing system in Go that demonstrates:
- **NDJSON file processing** without memory bloat
- **HTTP API** with proper timeouts  
- **Idempotent operations** and transaction handling

**Time allocation:** ~45 minutes

---

## Directory Structure

```
GOTEST/
├── data/                    # Sample data files
│   └── purchases.ndjson    # Newline-delimited JSON data
├── docker-compose.yml      # Docker setup for PostgreSQL
├── Makefile               # Build and development commands
├── README.md              # This file
├── solution/              # Completed solution (if available)
├── sql/                   # Database schema and seed data
│   ├── schema.sql         # Table definitions
│   └── seed.sql           # Sample data insertion
├── starter/               # Go source code to complete
│   ├── go.mod             # Go module file
│   ├── go.sum             # Go module checksums
│   ├── ingest.go          # NDJSON parsing
│   ├── main.go            # Application entry point
│   ├── server.go          # HTTP server and handlers
│   ├── store.go           # Database interface and implementation
│   ├── enrich.go          # Purchase enrichment worker pool
│   └── retry.go           # Retry utilities
├── tests/                 # Test files
│   └── ingest_test.go     # File ingestion tests
├── STRUCTURE.md           # Additional project structure details
└── todo.md                # Development task list
```

---

## Problem Statement

You operate GameVault, a digital game distribution platform. You receive purchase events from various gaming platforms (Steam, Epic, Xbox) and need to:

1. **Ingest purchases** from gaming platforms

---

## Database Setup

**Quick Start with Docker (Recommended):**

```bash
# 1. Start PostgreSQL with schema and seed data
docker compose up -d

# 2. Verify database is ready
docker compose logs postgres

# 3. Database will be available at: postgres://developer:devpass123@localhost:5432/gamevault_test
```

**Manual Setup (Alternative):**

```bash
# 1. Create database and load schema  
export DATABASE_URL="postgres://username:password@localhost:5432/gamevault_test?sslmode=disable"

# 2. Run schema creation
psql $DATABASE_URL -f sql/schema.sql

# 3. Load seed data (optional)
psql $DATABASE_URL -f sql/seed.sql
```

---

## Core Requirements

### 1. Purchase Ingestion (`POST /ingest`)

- Accept **multipart file upload** (NDJSON format only)
- **Stream parse** files without loading entirely into memory
- **Handle purchases** by `transaction_id`
- Return JSON with `created`, `updated`, and `total` counts
- Use `context.Context` with timeouts on all DB operations

**Example request:**
```bash
curl -F "file=@data/purchases.ndjson" http://localhost:8080/ingest
# Returns: {"created": 5, "updated": 2, "total": 7}
```



---

## Implementation Guidelines

### File Structure
Complete the TODO sections in these starter files:

- `main.go` - Entry point, flag parsing, graceful shutdown
- `store.go` - Database interface and PostgreSQL implementation  
- `server.go` - HTTP handlers with proper error handling
- `ingest.go` - Streaming NDJSON parser
- `tests/*_test.go` - Table-driven tests and benchmarks

### Key Technical Requirements

**Go Language Mastery:**
- Use small, focused interfaces (avoid `interface{}`)
- Proper error wrapping with `fmt.Errorf("%w", err)`
- Context propagation from HTTP → service → database

**Database Engineering:**
- Prepared statements for hot paths
- Proper transaction boundaries and isolation
- Parameterized queries (never string concatenation)

**Concurrency:**
- No data races (test with `go run -race`)
- Context cancellation support throughout

**HTTP & I/O:**
- Server timeouts (read/write/idle)
- Streaming file processing
- Graceful shutdown with `http.Server.Shutdown()`
- Proper JSON error responses

---

## Testing Requirements

### Unit Tests
Write table-driven tests for:
- NDJSON parsing edge cases
- Timestamp parsing and validation

### Integration Tests  
Test database operations:
- Upsert behavior (created vs updated)

**Testing Strategy:**
- Use transactions that rollback to avoid test data pollution
- Or create temporary schema per test run

### Benchmarks
Include at least one benchmark:
- `BenchmarkStreamNDJSON` - test parsing throughput

---

## Acceptance Criteria

### Core Functionality (Must Have)
- [ ] File ingestion works with NDJSON
- [ ] Duplicate transactions are handled correctly
- [ ] All database queries use prepared statements
- [ ] Context timeouts enforced throughout

### Code Quality (Must Have)  
- [ ] Proper error wrapping and handling
- [ ] Small, testable functions and interfaces
- [ ] No data races when running with `-race`
- [ ] Clean separation of concerns (handlers/business logic/data)

### Testing (Must Have)
- [ ] Table-driven unit tests with edge cases
- [ ] At least one meaningful benchmark
- [ ] Database tests use transactions or cleanup properly

---

## Stretch Goals (If Time Permits)

Pick 1-2 if you finish early:

1. **Observability** - Add request logging and metrics counters
2. **Database Optimization** - Discuss connection pooling and prepared statement caching  
3. **Retry Policies** - Use your `Retry[T]()` function for transient DB errors
4. **Isolation Levels** - When would `SERIALIZABLE` be needed vs `READ COMMITTED`?

---

## Sample Data

Test files are provided in `data/`:
- `purchases.ndjson` - Sample newline-delimited JSON

### Generate Large Test Data

Use the included data generator for performance testing:

```bash
# Generate 1,000 records (default)
go run generate-data.go

# Generate 10,000 records for stress testing
go run generate-data.go -count 10000

# Custom output file
go run generate-data.go -count 50000 -output data/large-test.ndjson
```

**Features:**
- **Player ID reuse**: 70% of records reuse existing players for realistic loyalty calculations
- **Duplicate transactions**: 5% chance of duplicate `transaction_id` for testing
- **Timestamped output**: Defaults to `data/test-YYYYMMDD-HHMMSS.ndjson`
- **Realistic variety**: Multiple platforms, item types, and price ranges

---

## Scoring Rubric

| Area | Strong (2 pts) | Adequate (1 pt) | Weak (0 pts) |
|------|----------------|-----------------|--------------|
| **Go Mastery** | Clean APIs, proper error handling, context usage | Mostly idiomatic, minor issues | Leaky abstractions, no error wrapping |
| **Concurrency** | Bounded pools, no races, proper cancellation | Works but rough edges | Ad-hoc goroutines, possible races |
| **Database** | Correct TX boundaries, prepared statements, proper isolation | Mostly correct, minor bugs | String-concat queries, no transactions |
| **HTTP/IO** | Timeouts, graceful shutdown, streaming | Works, some blocking | No timeouts, memory inefficient |
| **Testing** | Table-driven, benchmarks, proper cleanup | Some tests, no bench | No meaningful tests |
| **Performance** | Smart allocations, connection pooling | Minor optimizations | No performance considerations |

**Scoring:**
- **9-12 points**: Strong hire
- **6-8 points**: On the fence, followup questions
- **≤5 points**: No hire

---

## Getting Started

1. **Setup database:**
   ```bash
   # Quick start
   make db-up
   
   # Or manually
   docker compose up -d
   ```

2. **Install dependencies:**
   ```bash
   cd starter/
   go mod tidy
   ```

3. **Complete the TODOs** in starter files

4. **Test as you go:**
   ```bash
   make test          # Run tests
   make test-race     # Test for race conditions
   make bench         # Run benchmarks
   
   # Or manually
   cd tests && go test -v
   cd tests && go test -bench=.
   cd starter && go run -race .
   ```

5. **Development workflow:**
   ```bash
   make dev           # Setup DB + run tests
   make run           # Start HTTP server
      ```

## Running the System

```bash
# Start HTTP server
make run
# Or manually: go run . -addr=:8080 -db="postgres://developer:devpass123@localhost:5432/gamevault_test?sslmode=disable"

# Test ingestion
curl -F "file=@data/purchases.ndjson" http://localhost:8080/ingest

```

---

**Good luck!** Focus on clean, idiomatic Go code that demonstrates production-ready patterns.