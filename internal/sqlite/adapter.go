package sqlite

import (
	"context"

	"github.com/inner-daydream/kvz/internal/kv"
)

type KvRepositoryAdapter struct {
	q Querier
}

func (r *KvRepositoryAdapter) AddHook(ctx context.Context, name string, script string) error {
	params := addHookParams{
		Name:   name,
		Script: script,
	}
	return r.q.addHook(ctx, params)
}

func (r *KvRepositoryAdapter) AttachHook(ctx context.Context, key string, hook string) error {
	params := attachHookParams{
		Key:  key,
		Hook: hook,
	}
	return r.q.attachHook(ctx, params)
}

func (r *KvRepositoryAdapter) GetVal(ctx context.Context, key string) (val string, err error) {
	return r.q.getVal(ctx, key)
}

func (r *KvRepositoryAdapter) ListHooks(ctx context.Context) ([]string, error) {
	return r.q.listHooks(ctx)
}

func (r *KvRepositoryAdapter) ListKeys(ctx context.Context) ([]string, error) {
	return r.q.listKeys(ctx)
}

func (r *KvRepositoryAdapter) SetVal(ctx context.Context, key string, val string) error {
	params := setValParams{
		Key: key,
		Val: val,
	}
	return r.q.setVal(ctx, params)
}

func NewRepository(querier Querier) *KvRepositoryAdapter {
	return &KvRepositoryAdapter{
		q: querier,
	}
}

func (r *KvRepositoryAdapter) GetAttachedHooks(ctx context.Context, key string) ([]kv.Hook, error) {
	sqliteHooks, err := r.q.getAttachedHooks(ctx, key)
	if err != nil {
		return nil, err
	}
	kvHooks := make([]kv.Hook, len(sqliteHooks))
	for i, sqliteHook := range sqliteHooks {
		kvHooks[i] = kv.Hook{
			Script: sqliteHook.Script,
			Name:   sqliteHook.Name,
		}
	}
	return kvHooks, nil
}

func (r *KvRepositoryAdapter) KeyExists(ctx context.Context, key string) (bool, error) {
	status, err := r.q.keyExists(ctx, key)
	if err != nil {
		return false, err
	}
	return status == 1, nil
}

func (r *KvRepositoryAdapter) HookExists(ctx context.Context, name string) (bool, error) {
	status, err := r.q.hookExists(ctx, name)
	if err != nil {
		return false, err
	}
	return status == 1, nil
}
