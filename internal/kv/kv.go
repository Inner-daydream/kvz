package kv

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

type KvService interface {
	Set(key string, val string) (err error)
	Get(key string) (val string, err error)
	AddHook(name string, script string) (err error)
	AttachHook(key string, hook string) error
	ListKeys() ([]string, error)
	ListHooks() ([]string, error)
	GetAttachedHooks(key string) ([]Hook, error)
	ExecHooks(hooks []Hook, newVal string) []CmdOutput
}
type KvRepository interface {
	GetVal(ctx context.Context, key string) (val string, err error)
	SetVal(ctx context.Context, key string, val string) error
	AddHook(ctx context.Context, name string, script string) error
	AttachHook(ctx context.Context, key string, hook string) error
	ListKeys(ctx context.Context) ([]string, error)
	ListHooks(ctx context.Context) ([]string, error)
	GetAttachedHooks(ctx context.Context, key string) ([]Hook, error)
	KeyExists(ctx context.Context, key string) (bool, error)
	HookExists(ctx context.Context, name string) (bool, error)
}

type Hook struct {
	Name   string
	Script string
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

func (s *kvService) AddHook(name string, script string) error {
	ctx := context.Background()
	if name == "" || script == "" {
		return fmt.Errorf("key or hook name may not be empty")
	}
	err := s.r.AddHook(ctx, name, script)
	if err != nil {
		return fmt.Errorf("failed to create the hook: %w", err)
	}
	return nil
}

func (s *kvService) AttachHook(key string, hook string) error {
	if key == "" || hook == "" {
		return fmt.Errorf("key or hook name may not be empty")
	}
	ctx := context.Background()
	keyExists, err := s.r.KeyExists(ctx, key)
	if err != nil {
		return fmt.Errorf("could not check if key exists: %w", err)
	}
	if !keyExists {
		return fmt.Errorf("specified key does not exist")
	}
	ctx = context.Background()
	hookExists, err := s.r.HookExists(ctx, hook)
	if err != nil {
		return fmt.Errorf("could not check if hook exists: %w", err)
	}
	if !hookExists {
		return fmt.Errorf("specified hook does not exist")
	}
	err = s.r.AttachHook(ctx, key, hook)
	if err != nil {
		return fmt.Errorf("failed to attach the %s hook to the %s key: %w", hook, key, err)
	}
	return nil
}

func (s *kvService) ListHooks() (hookNames []string, err error) {
	ctx := context.Background()
	hookNames, err = s.r.ListHooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the list of hooks: %w", err)
	}
	return hookNames, nil
}

func (s *kvService) ListKeys() (keys []string, err error) {
	ctx := context.Background()
	keys, err = s.r.ListKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the list of keys: %w", err)
	}
	return keys, nil
}

func (s *kvService) GetAttachedHooks(key string) ([]Hook, error) {
	if key == "" {
		return nil, fmt.Errorf("key may not be empty")
	}
	ctx := context.Background()
	keyExists, err := s.r.KeyExists(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("could not determine if the %s key is stored: %w", key, err)
	}
	if !keyExists {
		return nil, fmt.Errorf("the key %s is not stored: %w", key, err)
	}
	ctx = context.Background()
	hooks, err := s.r.GetAttachedHooks(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get the hooks attached to the %s key", err)
	}
	return hooks, nil
}

type CmdOutput struct {
	Stdout string
	Stderr string
	Error  error
	Caller string
}

func (s *kvService) ExecHooks(hooks []Hook, newVal string) []CmdOutput {
	if len(hooks) == 0 {
		return nil
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	cmdOutputs := make([]CmdOutput, len(hooks))
	var stdout, stderr bytes.Buffer
	for i, hook := range hooks {
		cmd := exec.Command(shell, "-c", hook.Script)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Env = append(cmd.Env, fmt.Sprintf("NEW_VAL=%s", newVal))
		err := cmd.Run()
		cmdOutputs[i] = CmdOutput{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			Error:  err,
			Caller: hook.Name,
		}
	}
	return cmdOutputs
}

func NewServcice(r KvRepository) KvService {
	return &kvService{r: r}
}
