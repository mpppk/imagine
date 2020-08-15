package repoimpl_test

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/repository"

	bolt "go.etcd.io/bbolt"

	"github.com/mpppk/imagine/infra/repoimpl"

	"github.com/mpppk/imagine/domain/model"
)

func TestBBoltTag_Update(t *testing.T) {
	fileName := "TestBBoltTag_Update.db"
	var wsName model.WSName = "workspace-for-test"
	oldTag := &model.Tag{
		ID:   0,
		Name: "old",
	}
	newTag := &model.Tag{
		ID:   0,
		Name: "replaced",
	}
	type args struct {
		tag *model.Tag
	}
	tests := []struct {
		name    string
		args    args
		oldTags []*model.Tag

		wantErr bool
	}{
		{
			oldTags: []*model.Tag{oldTag},
			args:    args{tag: newTag},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, db := newTagRepository(t, wsName, fileName)
			defer teardown(t, fileName, db)

			if err := repo.Update(wsName, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := repo.Get(wsName, tt.args.tag.ID)
			if err != nil {
				t.Errorf("failed to get tag: %v: %v", newTag.ID, err)
			}
			if !reflect.DeepEqual(got, newTag) {
				t.Errorf("want: %#v, got: %#v", newTag, got)
			}
		})
	}
}

func newTagRepository(t *testing.T, wsName model.WSName, fileName string) (repository.Tag, *bolt.DB) {
	t.Helper()
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		t.Errorf("failed to create bolt db: %v", err)
	}

	repo := repoimpl.NewBBoltTag(db)
	if err != nil {
		t.Errorf("failed to create BBoltAsset: %v", err)
	}

	if err := repo.Init(wsName); err != nil {
		t.Errorf("failed to init tag repository: %v", err)
	}
	return repo, db
}
