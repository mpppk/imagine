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

func TestBBoltAsset_SearchByTags(t *testing.T) {
	fileName := "testfile.db"
	var wsName model.WSName = "workspace-for-test"
	tagA, tagB, tagC := model.Tag("a"), model.Tag("b"), model.Tag("c")
	assets := []*model.Asset{
		{
			ID:   0,
			Name: "0",
			Path: "path/to/0",
			Tags: nil,
		},
		{
			ID:   1,
			Name: "1",
			Path: "path/to/1",
			Tags: []model.Tag{tagA},
		},
		{
			ID:   2,
			Name: "2",
			Path: "path/to/2",
			Tags: []model.Tag{tagA, tagB},
		},
		{
			ID:   3,
			Name: "3",
			Path: "path/to/3",
			Tags: []model.Tag{tagB},
		},
		{
			ID:   4,
			Name: "4",
			Path: "path/to/4",
			Tags: []model.Tag{tagC},
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
			args:    args{tags: []model.Tag{"a"}},
			want:    []*model.Asset{assets[1], assets[2]},
			wantErr: false,
		},
		{
			args:    args{tags: []model.Tag{"b"}},
			want:    []*model.Asset{assets[2], assets[3]},
			wantErr: false,
		},
		{
			args:    args{tags: []model.Tag{"c"}},
			want:    []*model.Asset{assets[4]},
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
				if err := repo.Add(wsName, asset); err != nil {
					t.Errorf("failed to add asset:%v, error:%v", asset, err)
				}
			}

			got, err := repo.ListByTags(wsName, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchByTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchByTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBBoltAsset_Update(t *testing.T) {
	fileName := "testfile.db"
	var wsName model.WSName = "workspace-for-test"
	tagA := model.Tag("a")
	newAsset := &model.Asset{
		ID:   0,
		Name: "replaced",
		Path: "path/to/0",
		Tags: []model.Tag{tagA},
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

			if err := repo.Update(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := repo.Get(wsName, tt.args.asset.ID)
			if err != nil {
				t.Errorf("failed to get asset: %v: %v", newAsset.ID, err)
			}
			if !reflect.DeepEqual(got, newAsset) {
				t.Errorf("want: %#v, got: %#v", newAsset, got)
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
	if err := os.Remove(fileName); err != nil {
		t.Errorf("failed to remove test file: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Errorf("failed to close db: %v", err)
	}
}
