package repoimpl_test

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
)

func TestBBoltTag_Update(t *testing.T) {
	fileName := "TestBBoltTag_Update.db"
	var wsName model.WSName = "workspace-for-test"
	oldTag := &model.TagWithIndex{
		Tag: &model.Tag{
			ID:   0,
			Name: "old",
		},
	}
	newTag := &model.TagWithIndex{
		Tag: &model.Tag{
			ID:   0,
			Name: "replaced",
		},
	}
	type args struct {
		tag *model.TagWithIndex
	}
	tests := []struct {
		name    string
		args    args
		oldTags []*model.TagWithIndex

		wantErr bool
	}{
		{
			oldTags: []*model.TagWithIndex{oldTag},
			args:    args{tag: newTag},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecases, teardown := testutil.SetUpUseCases(t, fileName, wsName)
			defer teardown()

			if err := usecases.Client.Tag.Update(wsName, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _, err := usecases.Client.Tag.Get(wsName, tt.args.tag.ID)
			if err != nil {
				t.Errorf("failed to get tag: %v: %v", newTag.ID, err)
			}
			if !reflect.DeepEqual(got, newTag) {
				t.Errorf("want: %#v, got: %#v", newTag, got)
			}
		})
	}
}
