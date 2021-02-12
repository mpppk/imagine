package interactor_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestTag_SetTags(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name      string
		existTags []*model.Tag
		want      []*model.Tag
		args      args
		wantErr   bool
	}{
		{
			name:      "must be empty",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args:      args{ws: "default-workspace", tagNames: []string{}},
			want:      nil,
		},
		{
			name:      "set one tag",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args:      args{ws: "default-workspace", tagNames: []string{"new-tag1"}},
			want:      []*model.Tag{{ID: 1, Name: "new-tag1"}},
		},
		{
			name:      "set two tags",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args:      args{ws: "default-workspace", tagNames: []string{"new-tag1", "new-tag2"}},
			want:      []*model.Tag{{ID: 1, Name: "new-tag1"}, {ID: 2, Name: "new-tag2"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			idList, err := ut.Usecases.Tag.SetTags(tt.args.ws, tt.args.tagNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			var wantIdList []model.TagID
			for _, tag := range tt.want {
				wantIdList = append(wantIdList, tag.ID)
			}

			testutil.Diff(t, wantIdList, idList)

			tags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.want, tags)
		})
	}
}
