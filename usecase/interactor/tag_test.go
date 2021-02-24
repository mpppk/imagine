package interactor_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestTag_SaveTags(t *testing.T) {
	type args struct {
		ws   model.WSName
		tags []*model.Tag
	}
	tests := []struct {
		name          string
		existTagNames []string
		want          []*model.Tag
		wantTags      []*model.Tag
		args          args
		wantErr       bool
	}{
		{
			name:          "update tag",
			existTagNames: []string{"tag1", "tag2"},
			args: args{ws: "default-workspace", tags: []*model.Tag{
				testutil.MustNewTag(1, "updated-tag1", 0),
			}},
			want: []*model.Tag{
				testutil.MustNewTag(1, "updated-tag1", 0),
			},
			wantTags: []*model.Tag{
				testutil.MustNewTag(1, "updated-tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			},
		},
		{
			name:          "add tag",
			existTagNames: []string{"tag1", "tag2"},
			args: args{ws: "default-workspace", tags: []*model.Tag{
				testutil.MustNewTag(0, "tag3", 2),
			}},
			want: []*model.Tag{
				testutil.MustNewTag(3, "tag3", 2),
			},
			wantTags: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
				testutil.MustNewTag(3, "tag3", 2),
			},
		},
		{
			name:          "fail if non exist tag is provided",
			existTagNames: []string{"tag1", "tag2"},
			args: args{ws: "default-workspace", tags: []*model.Tag{
				testutil.MustNewTag(99, "tag99", 2),
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)
			tags, err := ut.Usecases.Tag.SaveTags(tt.args.ws, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			testutil.Diff(t, tt.want, tags)

			newTags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.wantTags, newTags)
		})
	}
}

func TestTag_SetTags(t *testing.T) {
	type args struct {
		ws   model.WSName
		tags []*model.UnindexedTag
	}
	tests := []struct {
		name          string
		existTagNames []string
		want          []*model.Tag
		wantTags      []*model.Tag
		args          args
		wantErr       bool
	}{
		{
			name:          "update tag",
			existTagNames: []string{"tag1", "tag2"},
			args: args{ws: "default-workspace", tags: []*model.UnindexedTag{
				testutil.MustNewUnindexedTag(1, "updated-tag1"),
			}},
			want: []*model.Tag{
				testutil.MustNewTag(1, "updated-tag1", 0),
			},
			wantTags: []*model.Tag{
				testutil.MustNewTag(1, "updated-tag1", 0),
			},
		},
		{
			name:          "replace tag",
			existTagNames: []string{"tag1", "tag2"},
			args: args{ws: "default-workspace", tags: []*model.UnindexedTag{
				testutil.MustNewUnindexedTag(0, "tag3"),
			}},
			want: []*model.Tag{
				testutil.MustNewTag(3, "tag3", 0),
			},
			wantTags: []*model.Tag{
				testutil.MustNewTag(3, "tag3", 0),
			},
		},
		{
			name:          "fail if non exist tag is provided",
			existTagNames: []string{"tag1", "tag2"},
			args: args{ws: "default-workspace", tags: []*model.UnindexedTag{
				testutil.MustNewUnindexedTag(99, "tag99"),
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)
			tags, err := ut.Usecases.Tag.SetTags(tt.args.ws, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			testutil.Diff(t, tt.want, tags)

			newTags := ut.Tag.List(tt.args.ws)
			testutil.Diff(t, tt.wantTags, newTags)
		})
	}
}

func TestTag_SetTagByNames(t *testing.T) {
	type args struct {
		ws       model.WSName
		tagNames []string
	}
	tests := []struct {
		name    string
		want    []*model.Tag
		args    args
		wantErr bool
	}{
		{
			name: "must be empty",
			args: args{ws: "default-workspace", tagNames: []string{}},
			want: nil,
		},
		{
			name: "set one tag",
			args: args{ws: "default-workspace", tagNames: []string{"new-tag1"}},
			want: []*model.Tag{
				testutil.MustNewTag(1, "new-tag1", 0),
			},
		},
		{
			name: "set two tags",
			args: args{ws: "default-workspace", tagNames: []string{"new-tag1", "new-tag2"}},
			want: []*model.Tag{
				testutil.MustNewTag(1, "new-tag1", 0),
				testutil.MustNewTag(2, "new-tag2", 1),
			},
		},
		{
			name: "remove tag if does not provided",
			args: args{ws: "default-workspace", tagNames: []string{"tag2"}},
			want: []*model.Tag{
				testutil.MustNewTag(1, "tag2", 0),
			},
		},
		{
			name: "add tag",
			args: args{ws: "default-workspace", tagNames: []string{"tag1", "tag2"}},
			want: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
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
		name    string
		want1   []*model.Tag
		want2   []*model.Tag
		args1   args
		args2   args
		wantErr bool
	}{
		{
			name:  "Replace by new tag",
			args1: args{ws: "default-workspace", tagNames: []string{"tag1", "tag2"}},
			args2: args{ws: "default-workspace", tagNames: []string{"tag3"}},
			want1: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			},
			want2: []*model.Tag{
				testutil.MustNewTag(3, "tag3", 0),
			},
		},
		{
			name:  "Add only new tag",
			args1: args{ws: "default-workspace", tagNames: []string{"tag1", "tag2"}},
			args2: args{ws: "default-workspace", tagNames: []string{"tag2", "tag3"}},
			want1: []*model.Tag{
				testutil.MustNewTag(1, "tag1", 0),
				testutil.MustNewTag(2, "tag2", 1),
			},
			want2: []*model.Tag{
				testutil.MustNewTag(2, "tag2", 0),
				testutil.MustNewTag(3, "tag3", 1),
			},
		},
	}
	testF := func(t *testing.T, ut *usecasetest.UseCases, args args, want []*model.Tag, wantErr bool) {
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
