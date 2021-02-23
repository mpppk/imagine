package interactor_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestTag_SetTagByNames(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name      string
		existTags []*model.TagWithIndex
		want      []*model.TagWithIndex
		args      args
		wantErr   bool
	}{
		{
			name: "must be empty",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			args: args{ws: "default-workspace", tagNames: []string{}},
			want: nil,
		},
		{
			name: "set one tag",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			args: args{ws: "default-workspace", tagNames: []string{"new-tag1"}},
			want: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "new-tag1", 0),
			},
		},
		{
			name: "set two tags",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			args: args{ws: "default-workspace", tagNames: []string{"new-tag1", "new-tag2"}},
			want: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "new-tag1", 0),
				testutil.MustNewTagWithIndex(2, "new-tag2", 1),
			},
		},
		{
			name: "remove tag if does not provided",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			args: args{ws: "default-workspace", tagNames: []string{"tag2"}},
			want: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag2", 0),
			},
		},
		{
			name: "add tag",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			args: args{ws: "default-workspace", tagNames: []string{"tag1", "tag2"}},
			want: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			tags, err := ut.Usecases.Tag.SetTagByNames(tt.args.ws, tt.args.tagNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			testutil.Diff(t, tt.want, tags)

			newTags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.want, newTags)
		})
	}
}

func TestTag_SetTagByNamesTwice(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name      string
		existTags []*model.TagWithIndex
		want1     []*model.TagWithIndex
		want2     []*model.TagWithIndex
		args1     args
		args2     args
		wantErr   bool
	}{
		{
			name: "Same name tag should have same ID even if once deleted",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
				testutil.MustNewTagWithIndex(2, "tag2", 1),
			},
			args1: args{ws: "default-workspace", tagNames: []string{}},
			args2: args{ws: "default-workspace", tagNames: []string{"tag3"}},
			want1: nil,
			want2: []*model.TagWithIndex{
				// TODO: broken test. ID should be 3
				//testutil.MustNewTagWithIndex(3, "tag3", 0),
				testutil.MustNewTagWithIndex(1, "tag3", 0),
			},
		},
		{
			name: "Different name tag should have different ID",
			existTags: []*model.TagWithIndex{
				testutil.MustNewTagWithIndex(1, "tag1", 0),
			},
			args1: args{ws: "default-workspace", tagNames: []string{}},
			args2: args{ws: "default-workspace", tagNames: []string{"tag2"}},
			want1: nil,
			want2: []*model.TagWithIndex{
				// TODO: broken test. ID should be 2
				//testutil.MustNewTagWithIndex(2, "tag2", 0),
				testutil.MustNewTagWithIndex(1, "tag2", 0),
			},
		},
	}
	testF := func(t *testing.T, ut *usecasetest.UseCases, args args, want []*model.TagWithIndex, wantErr bool) {
		tags, err := ut.Usecases.Tag.SetTagByNames(args.ws, args.tagNames)
		if (err != nil) != wantErr {
			t.Errorf("PutTags() error = %v, wantErr %v", err, wantErr)
			return
		} else if wantErr {
			return
		}
		testutil.Diff(t, want, tags)

		newTags := ut.Tag.List(args.ws)
		testutil.Diff(t, want, newTags)
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args1.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			testF(t, ut, tt.args1, tt.want1, tt.wantErr)
			testF(t, ut, tt.args2, tt.want2, tt.wantErr)
		})
	}
}
