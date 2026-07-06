# serverbase

Lightweight extracted runtime utilities for AE base server used by `base-server` and other apps.

This module provides:

- Module registry and bootstrapping helpers
- Minimal HTTP server wiring used by `cmd/minimal-server`
- Shared config helpers

Local development

```bash
# Run module tests
cd serverbase
go test ./...

# Use with base-server locally (base-server/go.mod contains a replace to ../serverbase)
```
