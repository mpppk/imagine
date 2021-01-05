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
		name string
		args args
		want []*model.Asset
	}{
		{
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
			assetsvc.Merge(tt.args.baseAssets, tt.args.otherAssets)
			testutil.Diff(t, tt.want, tt.args.baseAssets)
		})
	}
}
