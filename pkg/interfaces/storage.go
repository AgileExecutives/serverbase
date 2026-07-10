package interfaces

import "context"

// Storage is a small interface for object storage used by modules.
// Implementations should be safe for concurrent use unless documented otherwise.
type Storage interface {
    // Put stores the data under the given key (overwriting if exists).
    Put(ctx context.Context, key string, data []byte) error

    // Get retrieves the data for key. Returns an error if not found.
    Get(ctx context.Context, key string) ([]byte, error)

    // Delete removes the object for key. Missing keys may return nil or an error.
    Delete(ctx context.Context, key string) error

    // Exists returns true if key exists.
    Exists(ctx context.Context, key string) (bool, error)
}
