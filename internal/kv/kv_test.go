package kv_test

import (
	"reflect"
	"testing"

	"github.com/inner-daydream/kvz/internal/kv"
	"github.com/inner-daydream/kvz/internal/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func Test_kvService_Set(t *testing.T) {

	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid kv pair",
			args: args{
				key: "testKey1",
				val: "testValue1",
			},
			wantErr: false,
		},
		{
			name: "Empty key",
			args: args{
				key: "",
				val: "testValue2",
			},
			wantErr: true,
		},
		{
			name: "Empty value",
			args: args{
				key: "testKey3",
				val: "",
			},
			wantErr: true,
		},
		{
			name: "Assign new value",
			args: args{
				key: "testKey1",
				val: "newValue",
			},
			wantErr: false,
		},
	}
	db, err := sqlite.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	queries := sqlite.New(db)
	repo := sqlite.NewRepository(queries)
	service := kv.NewServcice(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.Set(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("kvService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kvService_Get(t *testing.T) {
	db, err := sqlite.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := sqlite.New(db)
	if err != nil {
		t.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.NewRepository(queries)
	service := kv.NewServcice(repo)
	key := "testKey1"
	val := "testValue1"
	service.Set("testKey1", "testValue1")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantVal string
		wantErr bool
	}{
		{
			name: "Non-existent key",
			args: args{
				key: "none",
			},
			wantVal: "",
			wantErr: true,
		},
		{
			name: "Empty key",
			args: args{
				key: "",
			},
			wantVal: "",
			wantErr: true,
		},
		{
			name: "Existing key",
			args: args{
				key: key,
			},
			wantVal: val,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotVal, err := service.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("kvService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("kvService.Get() = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_kvService_ListKeys(t *testing.T) {
	db, err := sqlite.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := sqlite.New(db)
	if err != nil {
		t.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.NewRepository(queries)
	service := kv.NewServcice(repo)
	keys := []string{"k1", "k2", "k3"}
	for _, key := range keys {
		err := service.Set(key, "testValue")
		if err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		name     string
		wantKeys []string
		wantErr  bool
	}{
		{
			name:     "List keys in the order they were added in",
			wantKeys: keys,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotKeys, err := service.ListKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("kvService.ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("kvService.ListKeys() = %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func Test_kvService_ListHooks(t *testing.T) {
	db, err := sqlite.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := sqlite.New(db)
	if err != nil {
		t.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.NewRepository(queries)
	service := kv.NewServcice(repo)
	hookNames := []string{"h1", "h2", "h3"}
	for _, hookName := range hookNames {
		err := service.AddScriptHook(hookName, "echo hello")
		if err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		name          string
		wantHookNames []string
		wantErr       bool
	}{
		{
			name:          "List hooks in the order they were added in",
			wantHookNames: hookNames,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotHookNames, err := service.ListHooks()
			if (err != nil) != tt.wantErr {
				t.Errorf("kvService.ListHooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHookNames, tt.wantHookNames) {
				t.Errorf("kvService.ListHooks() = %v, want %v", gotHookNames, tt.wantHookNames)
			}
		})
	}
}

func Test_kvService_AttachHook(t *testing.T) {
	db, err := sqlite.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := sqlite.New(db)
	if err != nil {
		t.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.NewRepository(queries)
	service := kv.NewServcice(repo)
	testKey := "test1"
	err = service.Set(testKey, "val1")
	if err != nil {
		t.Fatal(err)
	}
	testHook := "h1"
	err = service.AddScriptHook(testHook, "echo test")
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		key  string
		hook string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "attach hook to existing key",
			args: args{
				key:  testKey,
				hook: testHook,
			},
			wantErr: false,
		},
		{
			name: "attach hook to missing key",
			args: args{
				key:  "missingKey",
				hook: testHook,
			},
			wantErr: true,
		},
		{
			name: "attach missing hook to key",
			args: args{
				key:  testKey,
				hook: "missingHook",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := service.AttachHook(tt.args.key, tt.args.hook); (err != nil) != tt.wantErr {
				t.Errorf("kvService.AttachHook() - %s -  error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func Test_kvService_GetAttachedHooks(t *testing.T) {
	db, err := sqlite.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	queries := sqlite.New(db)
	if err != nil {
		t.Fatal(err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.NewRepository(queries)
	service := kv.NewServcice(repo)
	keys := []string{"k1", "k2", "k3"}
	for _, key := range keys {
		err := service.Set(key, "testValue")
		if err != nil {
			t.Fatal(err)
		}
	}
	hookNames := []string{"h1", "h2", "h3"}
	hookContent := "echo hello"
	testkey := keys[0]
	wantedHooks := make([]kv.Hook, 3)
	for i, hookName := range hookNames {
		err := service.AddScriptHook(hookName, hookContent)
		if err != nil {
			t.Fatal(err)
		}
		err = service.AttachHook(testkey, hookName)
		if err != nil {
			t.Fatal(err)
		}
		wantedHooks[i] = kv.Hook{
			Name:   hookName,
			Script: hookContent,
		}

	}

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []kv.Hook
		wantErr bool
	}{
		{
			name: "Get back only hooks attached to the key in the order they were attached",
			args: args{
				key: testkey,
			},
			want:    wantedHooks,
			wantErr: false,
		},
		{
			name: "Try to get hooks from a key where none are attached",
			args: args{
				key: keys[1],
			},
			want:    []kv.Hook{},
			wantErr: false,
		},
		{
			name: "Try to get hooks from a non-existing key",
			args: args{
				key: "NoKey",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty key provided",
			args: args{
				key: "",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetAttachedHooks(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("kvService.GetAttachedHooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("kvService.GetAttachedHooks() = %v, want %v", got, tt.want)
			}
		})
	}
}
