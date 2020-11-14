package repoimpl

import (
	"os"
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/model"
	"go.etcd.io/bbolt"
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
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

			id, exist, err := repo.Get(tt.args.ws, tt.args.path)
			if err != nil || !exist {
				t.Errorf("Add() exist = %v, error = %v, wantErr %v", exist, err, tt.wantErr)
			}
			if id != tt.args.assetID {
				t.Errorf("Add() error. added asset id is %v but return value of Get() is %v", tt.args.assetID, id)
			}
		})
	}
}

func Test_bboltPathRepository_AddList(t *testing.T) {
	type fields struct {
		base *boltRepository
	}
	type args struct {
		ws          model.WSName
		paths       []string
		assetIDList []model.AssetID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &bboltPathRepository{
				base: tt.fields.base,
			}
			if err := p.AddList(tt.args.ws, tt.args.paths, tt.args.assetIDList); (err != nil) != tt.wantErr {
				t.Errorf("AddList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bboltPathRepository_FilterExistPath(t *testing.T) {
	type fields struct {
		base *boltRepository
	}
	type args struct {
		ws    model.WSName
		paths []string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantNotExistPaths []string
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &bboltPathRepository{
				base: tt.fields.base,
			}
			gotNotExistPaths, err := p.FilterExistPath(tt.args.ws, tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterExistPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNotExistPaths, tt.wantNotExistPaths) {
				t.Errorf("FilterExistPath() gotNotExistPaths = %v, want %v", gotNotExistPaths, tt.wantNotExistPaths)
			}
		})
	}
}

func Test_bboltPathRepository_Get(t *testing.T) {
	type fields struct {
		base *boltRepository
	}
	type args struct {
		ws   model.WSName
		path string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantId    model.AssetID
		wantExist bool
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &bboltPathRepository{
				base: tt.fields.base,
			}
			gotId, gotExist, err := p.Get(tt.args.ws, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("Get() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotExist != tt.wantExist {
				t.Errorf("Get() gotExist = %v, want %v", gotExist, tt.wantExist)
			}
		})
	}
}

func Test_bboltPathRepository_GetList(t *testing.T) {
	type fields struct {
		base *boltRepository
	}
	type args struct {
		ws    model.WSName
		paths []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantIdList []uint64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &bboltPathRepository{
				base: tt.fields.base,
			}
			gotIdList, err := p.GetList(tt.args.ws, tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIdList, tt.wantIdList) {
				t.Errorf("GetList() gotIdList = %v, want %v", gotIdList, tt.wantIdList)
			}
		})
	}
}

func Test_newBBoltPathRepository(t *testing.T) {
	type args struct {
		b *bbolt.DB
	}
	tests := []struct {
		name string
		args args
		want *bboltPathRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newBBoltPathRepository(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newBBoltPathRepository() = %v, want %v", got, tt.want)
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
