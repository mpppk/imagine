package repoimpl_test

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"

	"github.com/mpppk/imagine/domain/repository"

	bolt "go.etcd.io/bbolt"

	"github.com/mpppk/imagine/infra/repoimpl"
)

func TestBoltMeta_SetAndGetVersion(t *testing.T) {
	fileName := "TestBoltMeta_SetAndGetVersion.db"
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			version: "0.0.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newMetaRepository(t, fileName)
			defer teardown(t, fileName, db)

			v, err := semver.New(tt.version)
			if err != nil {
				t.Errorf("failed to create semver struct: %v", err)
			}
			if err := repo.SetDBVersion(v); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotVersion, err := repo.GetDBVersion()
			if err != nil {
				t.Errorf("failed to get version: %v", err)
			}
			if !reflect.DeepEqual(v, gotVersion) {
				t.Errorf("want: %#v, got: %#v", v, gotVersion)
			}
		})
	}
}

func newMetaRepository(t *testing.T, fileName string) (repository.Meta, *bolt.DB) {
	t.Helper()
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		t.Errorf("failed to create bolt db: %v", err)
	}

	repo := repoimpl.NewBoltMeta(db)
	if err != nil {
		t.Errorf("failed to create BBoltAsset: %v", err)
	}

	return repo, db
}
