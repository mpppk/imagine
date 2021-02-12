package repoimpl_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/usecase/usecasetest"
)

func TestBBoltTag_Put(t *testing.T) {
	newTagWithIndex := func(id model.TagID, name string, index int) *model.TagWithIndex {
		return &model.TagWithIndex{
			Tag: &model.Tag{
				ID:   id,
				Name: name,
			},
			Index: index,
		}
	}
	type args struct {
		ws  model.WSName
		tag *model.TagWithIndex
	}
	tests := []struct {
		name      string
		args      args
		existTags []*model.Tag
		wantTags  []*model.Tag
		wantErr   bool
	}{
		{
			name:      "update tag",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args: args{
				ws:  "default-workspace",
				tag: newTagWithIndex(1, "updated-tag1", 0),
			},
			wantTags: []*model.Tag{{ID: 1, Name: "updated-tag1"}},
			wantErr:  false,
		},
		{
			name:      "add new tag",
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			args: args{
				ws:  "default-workspace",
				tag: newTagWithIndex(2, "updated-tag2", 1),
			},
			wantTags: []*model.Tag{{ID: 1, Name: "updated-tag1"}, {ID: 2, Name: "updated-tag2"}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.PutTags(tt.args.ws, tt.existTags)

			if err := ut.Usecases.Client.Tag.Put(tt.args.ws, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _, err := ut.Usecases.Client.Tag.Get(tt.args.ws, tt.args.tag.ID)
			if err != nil {
				t.Errorf("failed to get tag: %v: %v", tt.args.tag.ID, err)
			}

			testutil.Diff(t, tt.args.tag, got)
		})
	}
}
