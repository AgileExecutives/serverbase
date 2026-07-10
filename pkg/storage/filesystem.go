package storage

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileSystemStorage stores objects under a base directory on the local filesystem.
type FileSystemStorage struct {
	BaseDir string
}

// NewFileSystemStorage returns a new FileSystemStorage rooted at baseDir.
// The directory is created if it does not exist.
func NewFileSystemStorage(baseDir string) (*FileSystemStorage, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	return &FileSystemStorage{BaseDir: baseDir}, nil
}

func (f *FileSystemStorage) pathFor(key string) string {
	// keep it simple: treat key as relative path under BaseDir
	return filepath.Join(f.BaseDir, key)
}

func (f *FileSystemStorage) Put(ctx context.Context, key string, data []byte) error {
	p := f.pathFor(key)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return ioutil.WriteFile(p, data, 0o644)
}

func (f *FileSystemStorage) Get(ctx context.Context, key string) ([]byte, error) {
	p := f.pathFor(key)
	return ioutil.ReadFile(p)
}

func (f *FileSystemStorage) Delete(ctx context.Context, key string) error {
	p := f.pathFor(key)
	if err := os.Remove(p); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (f *FileSystemStorage) Exists(ctx context.Context, key string) (bool, error) {
	p := f.pathFor(key)
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
