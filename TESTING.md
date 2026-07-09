# Testing the serverbase packages

This document shows common testing patterns and commands for the `serverbase` module, and demonstrates how to use the helpers in `pkg/testutils`.

## Quick commands

Run all `serverbase` tests:

```bash
cd serverbase
go test ./... -v
```

Run a single package (example):

```bash
cd serverbase
go test ./pkg/testutils -v
```

Run a specific test in a package:
Run fast unit-only tests (prefers in-memory mocks):

```bash
make unit
```

This target sets `MOCK_EMAIL=true` and runs `go test -short` across modules to make unit runs faster.


```bash
cd serverbase
go test ./pkg/testutils -run TestHandler_RegisterEndpointSkeleton -v
```

## Using `pkg/testutils` helpers

- Create an in-memory sqlite DB and run migrations:

```go
db := testutils.SetupTestDB(t)
defer testutils.CleanupTestDB(db)
```

- Start a Gin router for handler tests:

```go
r := testutils.SetupTestRouter()
// register handler routes
w := testutils.MakeJSONRequest(t, r, "POST", "/auth/register", payload)
```

- Use fixtures for DB entities:

```go
tenant := testutils.CreateTestTenant(t, db, "My Tenant")
user := testutils.CreateTestUser(t, db, "u@example.com", "hash", tenant.ID)
```

- Use the in-memory mock email sender for assertions:

```go
mockEmail := testutils.NewMockEmailSender()
// inject into service or handler
last := mockEmail.Last()
require.NotNil(t, last)
require.Contains(t, last.Subject, "verify")
```

## Best practices

- Keep tests isolated: use `SetupTestDB` for DB-backed tests and `NewMemoryUserRepo` for fast pure-unit tests.
- Use `BeginTestTransaction` when a test needs to run inside a transaction and roll back at the end.
- Prefer `MakeJSONRequest` and `AssertJSONResponse` for concise handler tests.

## Troubleshooting

- If handlers expect configuration (e.g. `JWT_SECRET`), set it in the test process environment prior to generating tokens.
- If tests write mock files (e.g. `tmp/mock_emails.json`), add that path to `.gitignore` if necessary.

## Adding helpers

If you add a new common test helper, place it in `serverbase/pkg/testutils` and add an example unit test in that package so future contributors can discover it.
