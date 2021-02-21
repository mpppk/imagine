package repoimpl

import (
	"os"
	"testing"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

func Test_bboltPathRepository_Add(t *testing.T) {
	fileName := "Test_bboltPathRepository_Add.db"
	type args struct {
		ws      model.WSName
		path    string
		assetID model.AssetID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				ws:      "bbolt-repository-test",
				path:    "path/to/0.png",
				assetID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, tt.args.ws, fileName)
			defer teardown(t, fileName, db)

			if err := repo.Add(tt.args.ws, tt.args.path, tt.args.assetID); (err != nil) != tt.wantErr {
				t.Errorf("AddWithIndex() error = %v, wantErr %v", err, tt.wantErr)
			}

			id, exist, err := repo.Get(tt.args.ws, tt.args.path)
			if err != nil || !exist {
				t.Errorf("AddWithIndex() exist = %v, error = %v, wantErr %v", exist, err, tt.wantErr)
			}
			if id != tt.args.assetID {
				t.Errorf("AddWithIndex() error. added asset ID is %v but return value of Get() is %v", tt.args.assetID, id)
			}
		})
	}
}

func newRepository(t *testing.T, wsName model.WSName, fileName string) (*bboltPathRepository, *bolt.DB) {
	t.Helper()
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		t.Errorf("failed to create bolt db: %v", err)
	}

	repo := newBBoltPathRepository(db)
	if err != nil {
		t.Errorf("failed to create BBoltAsset: %v", err)
	}

	return repo, db
}

func teardown(t *testing.T, fileName string, db *bolt.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Errorf("failed to close db: %v", err)
	}
	if err := os.Remove(fileName); err != nil {
		t.Errorf("failed to remove test file: %v", err)
	}
}
