# HTTP Integration Test Documentation

Complete documentation for testing the Books HTTP API.

## Overview

This test suite includes:
- **Unit Tests**: Tests for individual business logic
- **Integration Tests**: End-to-end tests for HTTP endpoints
- **Benchmark Tests**: Performance tests to measure API speed
- **Error Handling Tests**: Tests to ensure proper error handling

## Test Structure

```
server/
├── http_integration_test.go  # Complete HTTP endpoint tests
├── benchmark_test.go         # Performance tests
├── server_test.go           # Test setup and infrastructure
└── server.go                # Server implementation
```

## Prerequisites

1. **Docker**: Required for PostgreSQL test container
2. **Go 1.23+**: To run tests
3. **Dependencies**: Install with `go mod download`

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test-all

# Run integration tests only
make test-integration

# Run benchmark tests
make test-benchmark

# Run with coverage report
make test-coverage
```

### Manual Test Commands

```bash
# Integration tests
go test -v ./server/http_integration_test.go ./server/server.go ./server/router.go ./server/db.go

# Benchmark tests
go test -v ./server/benchmark_test.go ./server/server.go ./server/router.go ./server/db.go -bench=.

# With race detection
go test -race -v ./server/...

# Parallel execution
go test -parallel 4 -v ./server/...
```

## Test Coverage

This test covers the following endpoints:

### ✅ Book API Endpoints

| Method | Endpoint | Test Coverage |
|--------|----------|---------------|
| POST | `/api/v1/books` | ✅ Create book (valid/invalid data) |
| GET | `/api/v1/books` | ✅ Get all books |
| GET | `/api/v1/books/{id}` | ✅ Get book by ID (valid/invalid ID) |
| PUT | `/api/v1/books/{id}` | ✅ Update book (valid/invalid data) |
| DELETE | `/api/v1/books/{id}` | ✅ Delete book |

### ✅ Utility Endpoints

| Method | Endpoint | Test Coverage |
|--------|----------|---------------|
| GET | `/api/v1/healthcheck` | ✅ Health check |
| POST | `/api/v1/process-url` | ✅ URL processing |

### ✅ Infrastructure Tests

- **CORS Headers**: Test CORS configuration
- **Content-Type Validation**: Test JSON content type validation
- **Concurrency**: Test concurrent requests handling
- **Error Handling**: Test various error scenarios
- **Rate Limiting**: Test server under load

## Test Scenarios

### 1. Book Creation Tests
```go
// Valid book creation
{
    "title": "Test Book",
    "author": "Test Author", 
    "published_year": 2023
}

// Invalid scenarios tested:
- Missing title
- Missing author  
- Invalid published year (0 or negative)
- Invalid JSON format
- Wrong content type
```

### 2. Book Retrieval Tests
```go
// Get all books - tests pagination and data consistency
// Get by ID - tests:
- Valid existing ID
- Non-existent ID (404)
- Invalid ID format (400)
```

### 3. Book Update Tests
```go
// Update scenarios:
- Valid complete update
- Partial updates
- Non-existent book update
- Invalid data validation
```

### 4. Book Deletion Tests
```go
// Delete scenarios:
- Valid deletion
- Non-existent book deletion
- Verify deletion (subsequent GET returns 404)
```

### 5. Concurrency Tests
```go
// Tests concurrent operations:
- Multiple simultaneous book creations
- Read/write race conditions
- Database connection pooling
```

## Benchmark Results

Run benchmarks to see performance:

```bash
make test-benchmark
```

Example output:
```
BenchmarkCreateBook-8         1000   1543 ns/op   312 B/op   4 allocs/op
BenchmarkGetAllBooks-8        2000   2134 ns/op   456 B/op   8 allocs/op
BenchmarkGetBookByID-8        3000   1012 ns/op   248 B/op   3 allocs/op
BenchmarkHealthCheck-8       10000    156 ns/op    48 B/op   1 allocs/op
```

## CI/CD Integration

For continuous integration:

```bash
# Complete CI pipeline
make ci

# This will run:
- Code formatting check
- Linting
- Race condition detection
- All tests with coverage
```

## Docker Testing

Tests can also be run in Docker:

```bash
make docker-test
```

## Test Database

Tests use PostgreSQL test containers that are:
- Isolated for each test suite
- Automatically cleaned up after tests
- Migration applied automatically
- Random port assignment for parallel runs

## Error Scenarios Tested

1. **HTTP Errors**:
   - 400 Bad Request (invalid JSON, validation errors)
   - 404 Not Found (non-existent resources)
   - 405 Method Not Allowed
   - 500 Internal Server Error

2. **Database Errors**:
   - Connection failures
   - Transaction rollbacks
   - Constraint violations

3. **Validation Errors**:
   - Required field validation
   - Data type validation
   - Business rule validation

## Performance Monitoring

Benchmark tests include:
- **Throughput**: Requests per second
- **Latency**: Response time distribution
- **Memory**: Allocation patterns
- **Concurrency**: Parallel request handling

## Debugging Tests

To debug failing tests:

```bash
# Verbose output
go test -v ./server/... -run TestSpecificTest

# Debug with delve
dlv test ./server/... -- -test.run TestSpecificTest

# Test with timeout
go test -timeout 30s -v ./server/...
```

## Test Maintenance

1. **Adding New Tests**: 
   - Follow existing patterns
   - Use table-driven tests for multiple scenarios
   - Include both positive and negative test cases

2. **Updating Tests**:
   - Update tests when API changes
   - Maintain backward compatibility test cases
   - Update benchmarks for performance regressions

3. **Test Data**:
   - Use realistic test data
   - Clean up test data after tests
   - Avoid hard-coded values

## Tips & Best Practices

1. **Isolation**: Each test is independent
2. **Cleanup**: Always cleanup resources
3. **Realistic Data**: Use realistic test data
4. **Error Cases**: Test error scenarios extensively
5. **Performance**: Monitor performance regressions
6. **Documentation**: Keep test documentation updated

## Troubleshooting

### Common Issues

1. **Docker not running**: Ensure Docker daemon is running
2. **Port conflicts**: Tests use random ports, shouldn't conflict
3. **Database connection**: Check PostgreSQL container logs
4. **Timeout issues**: Increase test timeout for slow systems

### Log Debugging

Enable debug logging in tests:
```go
// In test setup
log.Level(zerolog.DebugLevel)
```
