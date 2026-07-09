# testutils

This package provides helpers and small fakes used across unit and integration tests.

Main helpers

- Database
  - `SetupTestDB(t *testing.T) *gorm.DB` — creates an in-memory SQLite DB and runs core migrations.
  - `SetupTestDBWithLogging(t *testing.T) *gorm.DB` — same, but with query logging enabled.
  - `CleanupTestDB(db *gorm.DB)` — closes DB connection.
  - `MigrateTestDB(t *testing.T, db *gorm.DB, entities ...interface{})` — run AutoMigrate for extra entities.
  - `BeginTestTransaction(t *testing.T, db *gorm.DB) *gorm.DB` — start a test transaction (defer rollback).

- HTTP / Gin
  - `SetupTestRouter() *gin.Engine` — Gin router in test mode with recovery.
  - `MakeJSONRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder` — send JSON request and return recorder.
  - `MakeAuthenticatedRequest(router *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder` — like above, with `Authorization` header.
  - `ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{})` — unmarshal response body.
  - `AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, target interface{})` — assert status and parse JSON when applicable.
  - `AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedMessageContains string)` — assert error response shape and message.
  - Context helpers: `CreateTestContext()`, `CreateAuthenticatedTestContext(tenantID, userID uint)`, `SetupAuthContext`, `SetJSONBody`, `SetURLParam`, `SetQueryParam`.

- Fixtures & factories
  - `CreateTestTenant`, `CreateTestUser`, invoice helpers and small utilities: `GenerateTestEmail`, `Ptr`, `NowPtr`, etc.

- Mocks & fakes
  - `NewMockLogger()` — a `core.Logger` mock that accepts common calls.
  - `MockEmailSender` — in-memory email sender that records sent emails (`NewMockEmailSender()`).
  - Mock PDF, storage and other test doubles live in `mocks.go`.
  - Simple in-memory repo helpers: `NewMemoryUserRepo()`.
   - Simple in-memory repo helpers: `NewMemoryUserRepo()`.
   - Email repo helper: `NewInMemoryEmailRepo()` — returns a test `EmailRepo` suitable for handler/service unit tests.

- Assertions
  - Time & numeric assertions: `AssertTimeEqual`, `AssertTimeNotZero`, `AssertFloatEqual`, `RequireNoDBError`, plus pointer/slice/map helpers.

Examples

- See the examples in [examples_test.go](examples_test.go) for usage patterns:
  - unit tests with in-memory repos,
  - service tests using `SetupTestDB`,
  - a small handler httptest using `SetupTestRouter`.

Quick run

Run the package tests locally:

```bash
cd serverbase
go test ./pkg/testutils -v
```

Notes

- The helpers are intentionally small and dependency-light so they can be used from other modules in the monorepo.
- If you need an additional fake or helper, add it here alongside tests so it stays discoverable.
# Test Utilities Package

This package provides common testing utilities for the AE Backend project.

## Components

### Database Helpers (`database.go`)

Create and manage test databases for unit and integration tests:

```go
func TestMyService(t *testing.T) {
    db := testutils.SetupTestDB(t)
    defer testutils.CleanupTestDB(db)
    
    testutils.MigrateTestDB(t, db, &MyEntity{})
    
    // Your tests here
}
```

**Functions:**
- `SetupTestDB(t)` - Creates SQLite in-memory database
- `SetupTestDBWithLogging(t)` - Same but with query logging
- `CleanupTestDB(db)` - Closes database connection
- `MigrateTestDB(t, db, entities...)` - Runs auto-migration
- `TruncateTable(db, tableName)` - Clears a table
- `GetRowCount(db, tableName)` - Counts rows
- `BeginTestTransaction(t, db)` - Starts transaction for test

### Test Fixtures (`fixtures.go`)

Generate test data for your tests:

```go
tenant, _ := testutils.CreateTestTenant(db, "Test Company")
user, _ := testutils.CreateTestUser(db, "test@example.com", tenant.ID)

email := testutils.GenerateTestEmail("user", 1) // user.1@test.example.com
invoiceNum := testutils.GenerateTestInvoiceNumber(2026, 1) // TEST-2026-00001
```

**Helper Functions:**
- `Ptr[T](value)` - Returns pointer to value
- `TimePtr(time)` - Returns pointer to time
- `NowPtr()` - Returns pointer to current time
- `PastTimePtr(duration)` - Returns pointer to past time
- `FutureTimePtr(duration)` - Returns pointer to future time

### HTTP Test Helpers (`http_helpers.go`)

Test HTTP handlers with ease:

```go
func TestMyHandler(t *testing.T) {
    router := testutils.SetupTestRouter()
    
    handler := NewMyHandler(service)
    router.POST("/endpoint", handler.HandleRequest)
    
    body := map[string]interface{}{"name": "test"}
    w := testutils.MakeJSONRequest(router, "POST", "/endpoint", body)
    
    var response MyResponse
    testutils.AssertJSONResponse(t, w, http.StatusOK, &response)
}
```

**Functions:**
- `SetupTestRouter()` - Creates Gin test router
- `SetupAuthContext(c, tenantID, userID)` - Adds auth to context
- `MakeJSONRequest(router, method, path, body)` - Makes JSON request
- `MakeAuthenticatedRequest(router, method, path, token, body)` - Authenticated request
- `ParseJSONResponse(t, w, target)` - Parses JSON response
- `AssertJSONResponse(t, w, status, target)` - Asserts and parses response
- `AssertErrorResponse(t, w, status, message)` - Asserts error response
- `CreateTestContext()` - Creates test Gin context
- `SetJSONBody(c, body)` - Sets JSON body in context
- `SetURLParam(c, key, value)` - Sets URL parameter
- `SetQueryParam(c, key, value)` - Sets query parameter

### Mocks (`mocks.go`)

Pre-built mocks for common services:

```go
func TestWithMockEmail(t *testing.T) {
    mockEmail := new(testutils.MockEmailService)
    mockEmail.On("SendEmail", mock.Anything, "test@example.com", 
        "Subject", "Body").Return(nil)
    
    service := NewService(mockEmail)
    err := service.DoSomething()
    
    assert.NoError(t, err)
    mockEmail.AssertExpectations(t)
}
```

**Available Mocks:**
- `MockEmailService` - Email sending
- `MockPDFGenerator` - PDF generation
- `MockStorageService` - Document storage (MinIO)
- `MockAuditService` - Audit logging
- `MockEventBus` - Event publishing
- `MockCache` - Redis cache
- `MockInvoiceNumberService` - Invoice number generation

### Custom Assertions (`assertions.go`)

Domain-specific assertions:

```go
testutils.AssertTimeEqual(t, expected, actual, time.Second)
testutils.AssertDecimalEqual(t, 100.50, invoice.Total)
testutils.AssertUintNotZero(t, invoice.ID)
testutils.AssertSliceLength(t, items, 5)
```

**Functions:**
- `AssertTimeEqual(t, expected, actual, tolerance)` - Time equality with tolerance
- `AssertTimeNotZero(t, time)` - Time is not zero
- `AssertTimeBefore/After(t, time1, time2)` - Time ordering
- `AssertFloatEqual(t, expected, actual, tolerance)` - Float equality
- `AssertDecimalEqual(t, expected, actual)` - Currency equality (0.01 tolerance)
- `AssertUintNotZero(t, val)` - Uint is not zero
- `AssertSliceLength(t, slice, len)` - Slice has expected length
- `AssertMapHasKey(t, map, key)` - Map contains key

## Examples

### Testing a Service

```go
package services_test

import (
    "testing"
    "github.com/ae-base-server/pkg/testutils"
    "github.com/stretchr/testify/assert"
)

func TestInvoiceService_Create(t *testing.T) {
    // Arrange
    db := testutils.SetupTestDB(t)
    defer testutils.CleanupTestDB(db)
    
    testutils.MigrateTestDB(t, db, &Invoice{}, &InvoiceItem{})
    
    service := NewInvoiceService(db)
    
    // Act
    invoice, err := service.CreateInvoice(testData)
    
    // Assert
    assert.NoError(t, err)
    testutils.AssertUintNotZero(t, invoice.ID)
    testutils.AssertDecimalEqual(t, 100.00, invoice.Total)
}
```

### Testing a Handler

```go
package handlers_test

import (
    "net/http"
    "testing"
    "github.com/ae-base-server/pkg/testutils"
    "github.com/stretchr/testify/mock"
)

func TestInvoiceHandler_Create(t *testing.T) {
    // Arrange
    router := testutils.SetupTestRouter()
    mockService := new(testutils.MockInvoiceService)
    handler := NewInvoiceHandler(mockService)
    
    mockService.On("CreateInvoice", mock.Anything).Return(&Invoice{ID: 1}, nil)
    
    router.POST("/invoices", handler.CreateInvoice)
    
    // Act
    body := map[string]interface{}{"customer_id": 1}
    w := testutils.MakeJSONRequest(router, "POST", "/invoices", body)
    
    // Assert
    var response InvoiceResponse
    testutils.AssertJSONResponse(t, w, http.StatusCreated, &response)
    assert.Equal(t, uint(1), response.ID)
    mockService.AssertExpectations(t)
}
```

### Integration Test

```go
package integration_test

import (
    "testing"
    "github.com/ae-base-server/pkg/testutils"
)

func TestInvoiceWorkflow_EndToEnd(t *testing.T) {
    // Setup
    db := testutils.SetupTestDB(t)
    defer testutils.CleanupTestDB(db)
    
    testutils.MigrateTestDB(t, db, AllEntities()...)
    testutils.SeedMinimalTestData(db)
    
    // Initialize services
    invoiceService := NewInvoiceService(db)
    pdfService := NewPDFService()
    
    // Execute workflow
    draft := invoiceService.CreateDraft(...)
    finalized := invoiceService.Finalize(draft.ID)
    pdf := pdfService.Generate(finalized.ID)
    
    // Verify
    testutils.AssertUintNotZero(t, finalized.ID)
    assert.Equal(t, StatusFinalized, finalized.Status)
    testutils.AssertSliceNotEmpty(t, pdf)
}
```

## Best Practices

1. **Always cleanup**: Use `defer testutils.CleanupTestDB(db)`
2. **Use mocks for external services**: Email, PDF, Storage
3. **Test with realistic data**: Use fixtures to generate valid test data
4. **Assert meaningful things**: Use custom assertions for domain logic
5. **Keep tests isolated**: Each test should be independent
6. **Test error paths**: Don't just test happy paths

## Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
./scripts/run-tests-with-coverage.sh

# Run specific package
go test -v ./modules/invoice/services/...

# Run with race detector
go test -race ./...
```

## Dependencies

```bash
go get github.com/stretchr/testify
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
```
