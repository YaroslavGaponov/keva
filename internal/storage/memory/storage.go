package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/YaroslavGaponov/keva/pkg/utils"
)

type ValueType struct {
	Hash  []byte
	Value []byte
}

type MemStorage struct {
	mu      sync.Mutex
	options Options
	data    map[string]ValueType
}

func New(options Options) *MemStorage {
	s := MemStorage{
		options: options,
	}
	return &s
}

func (s *MemStorage) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]ValueType)
	return nil
}

func (s *MemStorage) Set(ctx context.Context, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = ValueType{
		Hash:  utils.GetHash(value),
		Value: value,
	}
	return nil
}

func (s *MemStorage) Get(ctx context.Context, key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.data[key]
	if !found {
		return nil, fmt.Errorf("key %s is not found", key)
	}
	return value.Value, nil
}

func (s *MemStorage) GetHash(ctx context.Context, key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.data[key]
	if !found {
		return nil, fmt.Errorf("key %s is not found", key)
	}
	return value.Hash, nil
}

func (s *MemStorage) Remove(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *MemStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = nil
	return nil
}
