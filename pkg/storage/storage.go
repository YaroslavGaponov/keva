package storage

import (
	"context"
)

type Storage interface {
	Open() error
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetHash(ctx context.Context, key string) ([]byte, error)
	Remove(ctx context.Context, key string) error
	Close() error
}
