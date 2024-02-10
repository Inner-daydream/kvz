package kv

import (
	"context"
	"fmt"
)

type KvService interface {
	Set(key string, val string) (err error)
	Get(key string) (val string, err error)
}
type KvRepository interface {
	GetVal(ctx context.Context, Key string) (val string, err error)
	SetVal(ctx context.Context, Key string, val string) error
}

type kvService struct {
	r KvRepository
}

func (s *kvService) Set(key string, val string) error {
	if key == "" {
		return fmt.Errorf("key should not be empty")
	}
	if val == "" {
		return fmt.Errorf("value should not be empty for key: %s", key)
	}
	ctx := context.Background()
	err := s.r.SetVal(ctx, key, val)
	if err != nil {
		return fmt.Errorf("failed to set a value to the %s key: %w", key, err)
	}
	return nil
}

func (s *kvService) Get(key string) (val string, err error) {
	ctx := context.Background()
	val, err = s.r.GetVal(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get the value from the %s key: %w", key, err)
	}
	return val, nil
}

func NewServcice(r KvRepository) KvService {
	return &kvService{r: r}
}
