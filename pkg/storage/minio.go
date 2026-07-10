package storage

import (
    "context"
    "io"
    "bytes"
    "errors"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "fmt"
)

// MinIOStorage is an S3-compatible adapter using minio client.
type MinIOStorage struct {
    client *minio.Client
    bucket string
}

// NewMinIOStorage constructs a MinIOStorage.
func NewMinIOStorage(endpoint, accessKey, secretKey, bucket string, secure bool) (*MinIOStorage, error) {
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: secure,
    })
    if err != nil {
        return nil, err
    }
    // ensure bucket exists (create if missing)
    ctx := context.Background()
    exists, err := client.BucketExists(ctx, bucket)
    if err != nil {
        return nil, err
    }
    if !exists {
        if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
            return nil, fmt.Errorf("failed create bucket: %w", err)
        }
    }
    return &MinIOStorage{client: client, bucket: bucket}, nil
}

func (m *MinIOStorage) Put(ctx context.Context, key string, data []byte) error {
    r := bytes.NewReader(data)
    _, err := m.client.PutObject(ctx, m.bucket, key, r, int64(len(data)), minio.PutObjectOptions{})
    return err
}

func (m *MinIOStorage) Get(ctx context.Context, key string) ([]byte, error) {
    obj, err := m.client.GetObject(ctx, m.bucket, key, minio.GetObjectOptions{})
    if err != nil {
        return nil, err
    }
    defer obj.Close()
    buf := new(bytes.Buffer)
    if _, err := io.Copy(buf, obj); err != nil {
        return nil, err
    }
    if buf.Len() == 0 {
        // double-check existence
        _, err := m.client.StatObject(ctx, m.bucket, key, minio.StatObjectOptions{})
        if err != nil {
            return nil, err
        }
    }
    return buf.Bytes(), nil
}

func (m *MinIOStorage) Delete(ctx context.Context, key string) error {
    return m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
}

func (m *MinIOStorage) Exists(ctx context.Context, key string) (bool, error) {
    _, err := m.client.StatObject(ctx, m.bucket, key, minio.StatObjectOptions{})
    if err != nil {
        var re minio.ErrorResponse
        if errors.As(err, &re) {
            if re.Code == "NoSuchKey" {
                return false, nil
            }
        }
        return false, err
    }
    return true, nil
}
