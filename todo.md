# Project Simplification Plan

## Phase 1: Remove CSV Support

### Documentation Changes
- [x] Update README.md overview to mention only NDJSON
- [x] Remove CSV references from directory structure
- [x] Update ingestion endpoint documentation to only show NDJSON examples
- [x] Remove CSV from file structure requirements
- [x] Update sample data references to only mention purchases.ndjson

### Code Changes
- [x] Remove CSV parsing logic from ingest.go
- [x] Update multipart file handler to only accept NDJSON
- [x] Remove CSV-related test cases from ingest_test.go
- [x] Update error messages to only reference NDJSON format
- [x] Remove CSV sample data file from data/ directory

### Testing Updates
- [x] Remove CSV parsing tests
- [x] Update integration tests to only use NDJSON (no integration tests found)
- [x] Remove CSV benchmarks if they exist (none found)
- [x] Update acceptance criteria to remove CSV requirements (already done)

### Implementation Notes
- Keep streaming NDJSON parser
- Maintain idempotent upsert behavior
- Preserve error handling for unsupported formats
- Keep file upload validation but restrict to NDJSON only