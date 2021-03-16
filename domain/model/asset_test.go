package model_test

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/testutil"
)

func TestReplaceBoundingBoxByID(t *testing.T) {
	type args struct {
		boxes      []*model.BoundingBox
		replaceBox *model.BoundingBox
	}
	tests := []struct {
		name         string
		args         args
		wantNewBoxes []*model.BoundingBox
	}{
		{
			args: args{
				boxes: []*model.BoundingBox{
					{ID: 0},
					{ID: 1},
					{ID: 2},
				},
				replaceBox: &model.BoundingBox{ID: 1, X: 1},
			},
			wantNewBoxes: []*model.BoundingBox{
				{ID: 0}, {ID: 1, X: 1}, {ID: 2},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if gotNewBoxes := model.ReplaceBoundingBoxByID(tt.args.boxes, tt.args.replaceBox); !reflect.DeepEqual(gotNewBoxes, tt.wantNewBoxes) {
				t.Errorf("ReplaceBoundingBoxByID() = %v, want %v", gotNewBoxes, tt.wantNewBoxes)
			}
		})
	}
}

func TestAssetIsUpdatableByID(t *testing.T) {
	tests := []struct {
		name  string
		asset *model.Asset
		want  bool
	}{
		{
			name:  "return true if asset has ID",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
			want:  true,
		},
		{
			name: "return true if asset has ID and boxes have tag ID",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
				{TagID: 1}, {TagID: 2},
			}},
			want: true,
		},
		{
			name:  "return false if asset does not have ID",
			asset: &model.Asset{Name: "path1", Path: "path1"},
			want:  false,
		},
		{
			name: "return false if box does not have tag ID",
			asset: &model.Asset{Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
				{},
			}},
			want: false,
		},
		{
			name:  "return false if asset is nil",
			asset: nil,
			want:  false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.asset.IsUpdatableByID(); got != tt.want {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestAssetMerge(t *testing.T) {
	type args struct {
		asset *model.Asset
	}
	tests := []struct {
		name    string
		args    args
		asset   *model.Asset
		want    *model.Asset
		wantErr bool
	}{
		{
			name:  "do nothing if arg asset is nil",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: nil,
			},
			want: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
		},
		{
			name:  "update path",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &model.Asset{ID: 1, Path: "path2"},
			},
			want: &model.Asset{ID: 1, Name: "path2", Path: "path2"},
		},
		{
			name:  "update path even if arg asset does not have ID",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &model.Asset{Path: "path2"},
			},
			want: &model.Asset{ID: 1, Name: "path2", Path: "path2"},
		},
		{
			name:  "reserve path because the property is omitted",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &model.Asset{ID: 1},
			},
			want: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
		},
		{
			name:  "update bounding boxes",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &model.Asset{ID: 1, BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
			},
			want: &model.Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
		},
		{
			name:  "update path and reserve boxes",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
			args: args{
				asset: &model.Asset{ID: 1, Path: "path2"},
			},
			want: &model.Asset{ID: 1, Name: "path2", Path: "path2", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
		},
		{
			name:  "fail if IDs are different",
			asset: &model.Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
			args: args{
				asset: &model.Asset{ID: 2, Path: "path2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.asset.UpdateBy(tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Asset.Merge: want error: %v, got: %v", tt.wantErr, err)
			}
			if tt.wantErr {
				return
			}
			testutil.Diff(t, tt.want, tt.asset)
		})
	}
}

func TestNewImportAssetFromJson(t *testing.T) {
	type args struct {
		json string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.ImportAsset
		wantErr bool
	}{
		{
			name: "asset and box have all properties",
			args: args{json: `{"id": 1, "name": "path1", "path": "path1", "boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			want: &model.ImportAsset{
				Asset: &model.Asset{ID: 1, Name: "path1", Path: "path1"},
				BoundingBoxes: []*model.ImportBoundingBox{
					{TagName: "tag1", BoundingBox: &model.BoundingBox{ID: 2, TagID: 3}},
				},
			},
		},
		{
			name: "asset has only path",
			args: args{json: `{"path": "path1", "boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			want: &model.ImportAsset{
				Asset: &model.Asset{Path: "path1"},
				BoundingBoxes: []*model.ImportBoundingBox{
					{TagName: "tag1", BoundingBox: &model.BoundingBox{ID: 2, TagID: 3}},
				},
			},
		},
		{
			name: "asset has only id",
			args: args{json: `{"id": 1, "boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			want: &model.ImportAsset{
				Asset: &model.Asset{ID: 1},
				BoundingBoxes: []*model.ImportBoundingBox{
					{TagName: "tag1", BoundingBox: &model.BoundingBox{ID: 2, TagID: 3}},
				},
			},
		},
		{
			name: "box has only tag id",
			args: args{json: `{"id": 1, "boundingBoxes": [{"tagID": 3}]}`},
			want: &model.ImportAsset{
				Asset:         &model.Asset{ID: 1},
				BoundingBoxes: []*model.ImportBoundingBox{{BoundingBox: &model.BoundingBox{TagID: 3}}},
			},
		},
		{
			name: "box has only tag name",
			args: args{json: `{"id": 1, "boundingBoxes": [{"tagName": "tag1"}]}`},
			want: &model.ImportAsset{
				Asset:         &model.Asset{ID: 1},
				BoundingBoxes: []*model.ImportBoundingBox{{TagName: "tag1"}},
			},
		},
		{
			name:    "invalid json",
			args:    args{json: `{`},
			wantErr: true,
		},
		{
			name:    "asset has no id and path",
			args:    args{json: `{"boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			wantErr: true,
		},
		{
			name:    "box has no tag id and name",
			args:    args{json: `{"id": 1, "name": "path1", "path": "path1", "boundingBoxes": [{"id": 2}]}`},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := model.NewImportAssetFromJson([]byte(tt.args.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewImportAssetFromJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			testutil.Diff(t, got, tt.want)
		})
	}
}

func TestAsset_IsAddable(t *testing.T) {
	tests := []struct {
		name  string
		asset *model.Asset
		want  bool
	}{
		{
			name: "addable",
			asset: &model.Asset{
				ID:   0,
				Name: "path1",
				Path: "path1",
			},
			want: true,
		},
		{
			name:  "not addable because asset is nil",
			asset: nil,
			want:  false,
		},
		{
			name: "not addable because ID is not zero",
			asset: &model.Asset{
				ID:   1,
				Name: "path1",
				Path: "path1",
			},
			want: false,
		},
		{
			name: "not addable because Path is empty",
			asset: &model.Asset{
				ID:   0,
				Name: "path1",
				Path: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.asset.IsAddable()
			testutil.Diff(t, tt.want, got)
		})
	}
}
