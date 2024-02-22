package sqlite

import (
	"context"
	"database/sql"

	"github.com/inner-daydream/kvz/internal/kv"
)

type KvRepositoryAdapter struct {
	q Querier
}

func (r *KvRepositoryAdapter) DeleteHook(ctx context.Context, name string) error {
	return r.q.deleteHook(ctx, name)
}

func (r *KvRepositoryAdapter) DeleteKey(ctx context.Context, key string) error {
	return r.q.deleteKey(ctx, key)
}

func (r *KvRepositoryAdapter) AddFileHook(ctx context.Context, name string, content string) error {
	params := addFileHookParams{
		Name: name,
		Script: sql.NullString{
			Valid:  true,
			String: content,
		},
	}
	return r.q.addFileHook(ctx, params)
}

// AddFilePathHook implements kv.KvRepository.
func (r *KvRepositoryAdapter) AddFilePathHook(ctx context.Context, name string, filepath string) error {
	params := addFilePathHookParams{
		Name: name,
		Filepath: sql.NullString{
			Valid:  true,
			String: filepath,
		},
	}
	return r.q.addFilePathHook(ctx, params)
}

func (r *KvRepositoryAdapter) AddScriptHook(ctx context.Context, name string, script string) error {
	params := addScriptHookParams{
		Name: name,
		Script: sql.NullString{
			Valid:  true,
			String: script,
		},
	}
	return r.q.addScriptHook(ctx, params)
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
		script := ""
		if sqliteHook.Script.Valid {
			script = sqliteHook.Script.String
		}
		kvHooks[i] = kv.Hook{
			Script:      script,
			Name:        sqliteHook.Name,
			IsFile:      sqliteHook.IsFile,
			IsLocalFile: sqliteHook.Filepath.Valid,
			Filepath:    sqliteHook.Filepath.String,
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
