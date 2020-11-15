package repoimpl_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/repository"

	bolt "go.etcd.io/bbolt"

	"github.com/mpppk/imagine/infra/repoimpl"

	"github.com/mpppk/imagine/domain/model"
)

func TestBBoltAsset_add(t *testing.T) {
	fileName := "TestBBoltAsset_Add.db"
	var wsName model.WSName = "workspace-for-test"
	newAsset := &model.Asset{
		ID:            0,
		Name:          "replaced",
		Path:          "path/to/0",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0)},
	}
	type args struct {
		asset *model.Asset
	}
	tests := []struct {
		name string
		args args

		wantErr bool
	}{
		{
			args:    args{asset: newAsset},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, wsName, fileName)
			defer teardown(t, fileName, db)

			if _, err := repo.Add(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, exist, err := repo.Get(wsName, tt.args.asset.ID)
			if err != nil || !exist {
				t.Errorf("failed to get asset: %v: %v", newAsset.ID, err)
			}
			if !reflect.DeepEqual(got, newAsset) {
				t.Errorf("want: %#v, got: %#v", newAsset, got)
			}
		})
	}
}

func TestBBoltAsset_Update(t *testing.T) {
	fileName := "TestBBoltAsset_Update.db"
	var wsName model.WSName = "workspace-for-test"
	oldAsset := &model.Asset{
		ID:            0,
		Name:          "old",
		Path:          "path/to/1",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(1)},
	}
	newAsset := &model.Asset{
		ID:            0,
		Name:          "replaced",
		Path:          "path/to/0",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0)},
	}
	type args struct {
		asset *model.Asset
	}
	tests := []struct {
		name      string
		args      args
		oldAssets []*model.Asset
		wantErr   bool
	}{
		{
			oldAssets: []*model.Asset{oldAsset},
			args:      args{asset: newAsset},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, wsName, fileName)
			defer teardown(t, fileName, db)

			for _, asset := range tt.oldAssets {
				if _, _, err := repo.AddByFilePathIfDoesNotExist(wsName, asset.Path); (err != nil) != tt.wantErr {
					t.Errorf("failed to addByID assets: %#v", asset)
				}
			}

			if err := repo.Update(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _, err := repo.Get(wsName, tt.args.asset.ID)
			if err != nil {
				t.Errorf("failed to get asset: %v: %v", newAsset.ID, err)
			}
			if !reflect.DeepEqual(got, newAsset) {
				t.Errorf("want: %#v, got: %#v", newAsset, got)
			}
		})
	}
}

func TestBBoltAsset_GetByPath(t *testing.T) {
	fileName := "TestBBoltAsset_GetByPath.db"
	var wsName model.WSName = "workspace-for-test"
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		existPaths []string
		want       *model.Asset
		wantErr    bool
	}{
		{
			args:       args{"path/to/0.png"},
			existPaths: []string{"path/to/0.png"},
			want: &model.Asset{
				ID:   1,
				Name: "0",
				Path: "path/to/0.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, wsName, fileName)
			defer teardown(t, fileName, db)

			for _, path := range tt.existPaths {
				if _, ok, err := repo.AddByFilePathIfDoesNotExist(wsName, path); (err != nil) != tt.wantErr {
					t.Errorf("GetByPath() error whiile assets adding. error = %v, wantErr %v", err, tt.wantErr)
				} else if !ok {
					t.Errorf("failed to add asset in AddByFilePathIfDoesNotExist test exist: %v", ok)
				}
			}
			_, err := repo.ListBy(wsName, func(a *model.Asset) bool { return true })
			if err != nil {
				t.Errorf("failed to list assets. error: %v", err)
			}
			got, exist, err := repo.GetByPath(wsName, tt.args.path)
			if err != nil || !exist {
				t.Errorf("failed to get asset by GetByPath. exist: %v error: %v", exist, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %#v, got: %#v", tt.want, got)
			}
		})
	}
}

func TestBBoltAsset_AddByFilePathListIfDoesNotExist(t *testing.T) {
	fileName := "TestBBoltAsset_AddByFlePathListIfDoesNotExist.db"
	type args struct {
		ws           model.WSName
		filePathList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.AssetID
		wantErr bool
	}{
		{
			args: args{
				ws:           "workspace-for-test",
				filePathList: []string{"0.png", "1.png"},
			},
			want: []model.AssetID{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, tt.args.ws, fileName)
			defer teardown(t, fileName, db)

			idList, err := repo.AddByFilePathListIfDoesNotExist(tt.args.ws, tt.args.filePathList)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, id := range idList {
				asset, exist, err := repo.Get(tt.args.ws, id)
				if err != nil || !exist {
					t.Errorf("failed to get asset which added by AddByFilePathListifDoesNotExist() exist = %v error = %v", exist, err)
				}
				asset2, exist, err := repo.GetByPath(tt.args.ws, asset.Path)
				if err != nil || !exist {
					t.Errorf("failed to get asset which added by AddByFilePathListifDoesNotExist() exist = %v error = %v", exist, err)
				}

				if !reflect.DeepEqual(asset, asset2) {
					t.Errorf("inconsistecy detected on AddByFilePathListifDoesNotExisti() asset1: %#v, asset2: %#v", asset, asset2)
				}
			}
		})
	}
}

func newRepository(t *testing.T, wsName model.WSName, fileName string) (repository.Asset, *bolt.DB) {
	t.Helper()
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		t.Errorf("failed to create bolt db: %v", err)
	}

	repo := repoimpl.NewBBoltAsset(db)
	if err != nil {
		t.Errorf("failed to create BBoltAsset: %v", err)
	}

	if err := repo.Init(wsName); err != nil {
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
