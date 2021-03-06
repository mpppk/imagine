package model_test

import (
	"testing"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/testutil"
)

func TestQuery_Match(t *testing.T) {
	type args struct {
		asset  *model.Asset
		tagSet *model.TagSet
	}
	tests := []struct {
		name  string
		query *model.Query
		args  args
		want  bool
	}{
		{
			name: "equals op match if asset have same name tag",
			query: &model.Query{
				Op:    model.EqualsQueryOP,
				Value: "tag1",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: true,
		},
		{
			name: "equals op does not match if specified tag is not assigned to asset",
			query: &model.Query{
				Op:    model.EqualsQueryOP,
				Value: "tag2",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: false,
		},
		{
			name: "not-equals op match if asset does not have same name tag",
			query: &model.Query{
				Op:    model.NotEqualsQueryOP,
				Value: "tag2",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: true,
		},
		{
			name: "not-equals op does not match if asset have same name tag",
			query: &model.Query{
				Op:    model.NotEqualsQueryOP,
				Value: "tag1",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: false,
		},
		{
			name: "start-with op match if asset have tag which have name started with specified value",
			query: &model.Query{
				Op:    model.StartWithQueryOP,
				Value: "tag1",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1xxx", 0),
					testutil.MustNewTag(2, "tag2xxx", 1),
				}),
			},
			want: true,
		},
		{
			name: "start-with op match if asset have same name tag",
			query: &model.Query{
				Op:    model.StartWithQueryOP,
				Value: "tag1",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: true,
		},
		{
			name: "start-with op does not match if asset have tag which does not have matched prefix",
			query: &model.Query{
				Op:    model.StartWithQueryOP,
				Value: "tag1",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "xtag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: false,
		},
		{
			name: "start-with op does not match if asset have tag which does not have matched prefix2",
			query: &model.Query{
				Op:    model.StartWithQueryOP,
				Value: "tag2",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: false,
		},
		{
			name: "no-tags op match if asset does not have any tags",
			query: &model.Query{
				Op: model.NoTagsQueryOP,
			},
			args: args{
				asset: &model.Asset{
					ID:            1,
					Name:          "path1",
					Path:          "path1",
					BoundingBoxes: []*model.BoundingBox{},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: true,
		},
		{
			name: "no-tags op does not match if asset have some tags",
			query: &model.Query{
				Op: model.NoTagsQueryOP,
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
					BoundingBoxes: []*model.BoundingBox{
						{TagID: 1},
					},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: false,
		},
		{
			name: "path-equals op match if asset have same name path",
			query: &model.Query{
				Op:    model.PathEqualsQueryOP,
				Value: "path1",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: true,
		},
		{
			name: "path-equals op does not match if asset does not have same name path",
			query: &model.Query{
				Op:    model.PathEqualsQueryOP,
				Value: "path2",
			},
			args: args{
				asset: &model.Asset{
					ID:   1,
					Name: "path1",
					Path: "path1",
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.query.Match(tt.args.asset, tt.args.tagSet); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
