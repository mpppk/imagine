package assetsvc_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/service/assetsvc"

	"github.com/mpppk/imagine/domain/model"
)

func TestMerge(t *testing.T) {
	type args struct {
		baseAssets  []*model.Asset
		otherAssets []*model.Asset
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.Asset
		wantErr bool
	}{
		{
			args: args{
				baseAssets: []*model.Asset{
					{ID: 1, Name: "path1", Path: "path1"},
					nil,
				},
				otherAssets: []*model.Asset{
					{Path: "path1-updated"},
					{Name: "path2-updated", Path: "path2-updated"},
				},
			},
			want: []*model.Asset{
				{ID: 1, Name: "path1-updated", Path: "path1-updated"},
				nil,
			},
		},
		{
			name: "fail if Ids are different",
			args: args{
				baseAssets: []*model.Asset{
					{ID: 1, Name: "path1", Path: "path1"},
					nil,
				},
				otherAssets: []*model.Asset{
					{Path: "path1-updated"},
					{ID: 2, Name: "path2-updated", Path: "path2-updated"},
				},
			},
			want: []*model.Asset{
				{ID: 1, Name: "path1-updated", Path: "path1-updated"},
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := assetsvc.Update(tt.args.baseAssets, tt.args.otherAssets); (err != nil) != tt.wantErr {
				t.Errorf("assetsvc.Merge: want error: %v, got: %v", tt.wantErr, err)
			}
			if tt.wantErr {
				return
			}
			testutil.Diff(t, tt.want, tt.args.baseAssets)
		})
	}
}

func TestQuery(t *testing.T) {
	type args struct {
		assets  []*model.Asset
		queries []*model.Query
		tagSet  *model.TagSet
	}
	tests := []struct {
		name               string
		args               args
		wantNewAssets      []*model.Asset
		wantFilteredAssets []*model.Asset
	}{
		{
			name: "path-equals query",
			args: args{
				assets: []*model.Asset{
					{ID: 1, Path: "path1", Name: "path1"},
					{ID: 2, Path: "path2", Name: "path2"},
					{ID: 3, Path: "path3", Name: "path3"},
				},
				queries: []*model.Query{{Op: "path-equals", Value: "path1"}},
				tagSet:  &model.TagSet{},
			},
			wantNewAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1"},
			},
			wantFilteredAssets: []*model.Asset{
				{ID: 2, Path: "path2", Name: "path2"},
				{ID: 3, Path: "path3", Name: "path3"},
			},
		},
		{
			name: "tag based queries",
			args: args{
				assets: []*model.Asset{
					{ID: 1, Path: "path1", Name: "path1",
						BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}}},
					{ID: 2, Path: "path2", Name: "path2",
						BoundingBoxes: []*model.BoundingBox{{TagID: 2}}},
					{ID: 3, Path: "path3", Name: "path3"},
				},
				queries: []*model.Query{
					{Op: "equals", Value: "tag1"},
					{Op: "equals", Value: "tag2"},
				},
				tagSet: model.NewTagSet([]*model.Tag{
					testutil.MustNewTag(1, "tag1", 0),
					testutil.MustNewTag(2, "tag2", 1),
				}),
			},
			wantNewAssets: []*model.Asset{
				{
					ID: 1, Path: "path1", Name: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				},
			},
			wantFilteredAssets: []*model.Asset{
				{ID: 2, Path: "path2", Name: "path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}}},
				{ID: 3, Path: "path3", Name: "path3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewAssets, gotFilteredAssets := assetsvc.Query(tt.args.assets, tt.args.queries, tt.args.tagSet, false)
			testutil.Diff(t, tt.wantNewAssets, gotNewAssets)
			testutil.Diff(t, tt.wantFilteredAssets, gotFilteredAssets)
		})
	}
}
