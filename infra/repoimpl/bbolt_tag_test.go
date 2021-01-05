package repoimpl_test

import (
	"reflect"
	"testing"

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCasesWithTempDB(t, wsName)
			defer closer()
			defer remover()

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
