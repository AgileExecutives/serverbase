package storage

// Minimal storage adapter used by other packages.

type Storage struct{}

func NewStorage() *Storage { return &Storage{} }

func (s *Storage) Save(data []byte) (string, error) { return "ok", nil }
