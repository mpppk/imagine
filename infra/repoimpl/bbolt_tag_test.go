package repoimpl_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/usecase/usecasetest"
)

func TestBBoltTag_Update(t *testing.T) {
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
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, wsName, func(t *testing.T, ut *usecasetest.UseCases) {
			if err := ut.Usecases.Client.Tag.Update(wsName, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _, err := ut.Usecases.Client.Tag.Get(wsName, tt.args.tag.ID)
			if err != nil {
				t.Errorf("failed to get tag: %v: %v", newTag.ID, err)
			}

			testutil.Diff(t, newTag, got)
		})
	}
}
