package repoimpl_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
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
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0, "a")},
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
func TestBBoltAsset_SearchByTags(t *testing.T) {
	fileName := "TestBBoltAsset_SearchByTags.db"
	var wsName model.WSName = "workspace-for-test"
	boxA := repoimpl.CreateBoundingBox(0, "a")
	boxB := repoimpl.CreateBoundingBox(1, "b")
	boxC := repoimpl.CreateBoundingBox(2, "c")
	assets := []*model.Asset{
		{
			ID:            0,
			Name:          "0",
			Path:          "path/to/0",
			BoundingBoxes: nil,
		},
		{
			ID:            1,
			Name:          "1",
			Path:          "path/to/1",
			BoundingBoxes: []*model.BoundingBox{boxA},
		},
		{
			ID:            2,
			Name:          "2",
			Path:          "path/to/2",
			BoundingBoxes: []*model.BoundingBox{boxA, boxB},
		},
		{
			ID:            3,
			Name:          "3",
			Path:          "path/to/3",
			BoundingBoxes: []*model.BoundingBox{boxB},
		},
		{
			ID:            4,
			Name:          "4",
			Path:          "path/to/4",
			BoundingBoxes: []*model.BoundingBox{boxC},
		},
	}
	type args struct {
		tags []model.Tag
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.Asset
		wantErr bool
	}{
		{
			args:    args{tags: []model.Tag{*boxA.Tag}},
			want:    []*model.Asset{assets[1], assets[2]},
			wantErr: false,
		},
		{
			args:    args{tags: []model.Tag{*boxB.Tag}},
			want:    []*model.Asset{assets[2], assets[3]},
			wantErr: false,
		},
		{
			args:    args{tags: []model.Tag{*boxC.Tag}},
			want:    []*model.Asset{assets[4]},
			wantErr: false,
		},
		{
			args:    args{tags: []model.Tag{*boxA.Tag, *boxB.Tag}},
			want:    []*model.Asset{assets[2]},
			wantErr: false,
		},
		{
			args:    args{tags: []model.Tag{}},
			wantErr: true,
		},
		{
			args:    args{tags: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, wsName, fileName)
			defer teardown(t, fileName, db)
			for _, asset := range assets {
				if _, err := repo.Add(wsName, asset); err != nil {
					t.Errorf("failed to addByID asset:%v, error:%v", asset, err)
				}
			}

			got, err := repo.ListByTags(wsName, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				gotMsg := "\n"
				for i, asset := range got {
					gotMsg += fmt.Sprint(i) + pp.Sprintln(asset)
					// gotMsg += fmt.Sprintf("%d: %#v\n", i, asset)
				}
				wantMsg := "\n"
				for i, asset := range tt.want {
					wantMsg += fmt.Sprint(i) + pp.Sprintln(asset)
					// wantMsg += fmt.Sprintf("%d: %#v\n", i, asset)
				}
				t.Errorf("SearchByTags() got = %s\nwant = %s", gotMsg, wantMsg)
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
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(1, "b")},
	}
	newAsset := &model.Asset{
		ID:            0,
		Name:          "replaced",
		Path:          "path/to/0",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0, "a")},
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
	var wsName model.WSName = "workspace-for-test"
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
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newRepository(t, wsName, fileName)
			defer teardown(t, fileName, db)

			idList, err := repo.AddByFilePathListIfDoesNotExist(wsName, tt.args.filePathList)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, id := range idList {
				asset, exist, err := repo.Get(wsName, id)
				if err != nil || !exist {
					t.Errorf("failed to get asset which added by AddByFilePathListifDoesNotExist() exist = %v error = %v", exist, err)
				}
				asset2, exist, err := repo.GetByPath(wsName, asset.Path)
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
