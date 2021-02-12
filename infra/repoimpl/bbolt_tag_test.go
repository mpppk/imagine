package repoimpl_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/usecase/usecasetest"
)

func TestBBoltTag_Save(t *testing.T) {
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
		name          string
		args          args
		existTagNames []string
		want          model.TagID
		wantTags      []*model.Tag
		wantErr       bool
	}{
		{
			name:          "update tag",
			existTagNames: []string{"tag1"},
			args: args{
				ws:  "default-workspace",
				tag: newTagWithIndex(1, "updated-tag1", 0),
			},
			want:     1,
			wantTags: []*model.Tag{{ID: 1, Name: "updated-tag1"}},
			wantErr:  false,
		},
		{
			name:          "add new tag",
			existTagNames: []string{"tag1"},
			args: args{
				ws:  "default-workspace",
				tag: newTagWithIndex(0, "updated-tag2", 1),
			},
			want:     2,
			wantTags: []*model.Tag{{ID: 1, Name: "updated-tag1"}, {ID: 2, Name: "updated-tag2"}},
			wantErr:  false,
		},
		{
			name:          "fail if tag ID which does not exist is provided",
			existTagNames: []string{"tag1"},
			args: args{
				ws:  "default-workspace",
				tag: newTagWithIndex(2, "updated-tag2", 1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTags(tt.args.ws, tt.existTagNames)

			id, err := ut.Usecases.Client.Tag.Save(tt.args.ws, tt.args.tag)
			if testutil.FatalIfErrIsUnexpected(t, tt.wantErr, err) {
				return
			}
			testutil.Diff(t, tt.want, id)

			got, _, err := ut.Usecases.Client.Tag.Get(tt.args.ws, tt.args.tag.ID)
			if err != nil {
				t.Errorf("failed to get tag: %v: %v", tt.args.tag.ID, err)
			}

			testutil.Diff(t, tt.args.tag, got)
		})
	}
}

func TestBBoltTag_AddByNames(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name          string
		args          args
		existTagNames []string
		want          []model.TagID
		wantTags      []*model.Tag
		wantErr       bool
	}{
		{
			name:          "add tag",
			existTagNames: []string{"tag1"},
			args: args{
				ws:       "default-workspace",
				tagNames: []string{"tag2"},
			},
			want:     []model.TagID{2},
			wantTags: []*model.Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}},
			wantErr:  false,
		},
		{
			name:          "add tag to empty DB",
			existTagNames: []string{"tag1"},
			args: args{
				ws:       "default-workspace",
				tagNames: []string{"tag1"},
			},
			want:     []model.TagID{1},
			wantTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTags(tt.args.ws, tt.existTagNames)

			idList, err := ut.Usecases.Client.Tag.AddByNames(tt.args.ws, tt.args.tagNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.Diff(t, tt.want, idList)

			gotTags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.wantTags, gotTags)
		})
	}
}
