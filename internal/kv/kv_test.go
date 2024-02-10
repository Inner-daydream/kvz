package kv

import (
	"testing"

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
	db, err := sqlite.CreateDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal(err)
	}
	repo := sqlite.New(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &kvService{
				r: repo,
			}
			if err := s.Set(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("kvService.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kvService_Get(t *testing.T) {
	db, err := sqlite.CreateDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = sqlite.Migrate(db)
	repo := sqlite.New(db)
	if err != nil {
		t.Fatal(err)
	}
	service := NewServcice(repo)
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
