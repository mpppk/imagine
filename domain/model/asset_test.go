package model

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/testutil"
)

func TestReplaceBoundingBoxByID(t *testing.T) {
	type args struct {
		boxes      []*BoundingBox
		replaceBox *BoundingBox
	}
	tests := []struct {
		name         string
		args         args
		wantNewBoxes []*BoundingBox
	}{
		{
			args: args{
				boxes: []*BoundingBox{
					{ID: 0},
					{ID: 1},
					{ID: 2},
				},
				replaceBox: &BoundingBox{ID: 1, X: 1},
			},
			wantNewBoxes: []*BoundingBox{
				{ID: 0}, {ID: 1, X: 1}, {ID: 2},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if gotNewBoxes := ReplaceBoundingBoxByID(tt.args.boxes, tt.args.replaceBox); !reflect.DeepEqual(gotNewBoxes, tt.wantNewBoxes) {
				t.Errorf("ReplaceBoundingBoxByID() = %v, want %v", gotNewBoxes, tt.wantNewBoxes)
			}
		})
	}
}

func TestAssetIsUpdatableByID(t *testing.T) {
	tests := []struct {
		name  string
		asset *Asset
		want  bool
	}{
		{
			name:  "return true if asset has ID",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1"},
			want:  true,
		},
		{
			name: "return true if asset has ID and boxes have tag ID",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*BoundingBox{
				{TagID: 1}, {TagID: 2},
			}},
			want: true,
		},
		{
			name:  "return false if asset does not have ID",
			asset: &Asset{Name: "path1", Path: "path1"},
			want:  false,
		},
		{
			name: "return false if box does not have tag ID",
			asset: &Asset{Name: "path1", Path: "path1", BoundingBoxes: []*BoundingBox{
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
		asset *Asset
	}
	tests := []struct {
		name  string
		args  args
		asset *Asset
		want  *Asset
	}{
		{
			name:  "do nothing if arg asset is nil",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: nil,
			},
			want: &Asset{ID: 1, Name: "path1", Path: "path1"},
		},
		{
			name:  "update path",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &Asset{ID: 1, Path: "path2"},
			},
			want: &Asset{ID: 1, Name: "path2", Path: "path2"},
		},
		{
			name:  "reserve path because the property is omitted",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &Asset{ID: 1},
			},
			want: &Asset{ID: 1, Name: "path1", Path: "path1"},
		},
		{
			name:  "update bounding boxes",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1"},
			args: args{
				asset: &Asset{ID: 1, BoundingBoxes: []*BoundingBox{{TagID: 1}}},
			},
			want: &Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*BoundingBox{{TagID: 1}}},
		},
		{
			name:  "update path and reserve boxes",
			asset: &Asset{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*BoundingBox{{TagID: 1}}},
			args: args{
				asset: &Asset{ID: 1, Path: "path2"},
			},
			want: &Asset{ID: 1, Name: "path2", Path: "path2", BoundingBoxes: []*BoundingBox{{TagID: 1}}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.asset.Merge(tt.args.asset)
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
		want    *ImportAsset
		wantErr bool
	}{
		{
			name: "asset and box have all properties",
			args: args{json: `{"id": 1, "name": "path1", "path": "path1", "boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			want: &ImportAsset{
				Asset: &Asset{ID: 1, Name: "path1", Path: "path1"},
				BoundingBoxes: []*ImportBoundingBox{
					{TagName: "tag1", BoundingBox: &BoundingBox{ID: 2, TagID: 3}},
				},
			},
		},
		{
			name: "asset has only path",
			args: args{json: `{"path": "path1", "boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			want: &ImportAsset{
				Asset: &Asset{Path: "path1"},
				BoundingBoxes: []*ImportBoundingBox{
					{TagName: "tag1", BoundingBox: &BoundingBox{ID: 2, TagID: 3}},
				},
			},
		},
		{
			name: "asset has only id",
			args: args{json: `{"id": 1, "boundingBoxes": [{"id": 2, "tagID": 3,"tagName": "tag1"}]}`},
			want: &ImportAsset{
				Asset: &Asset{ID: 1},
				BoundingBoxes: []*ImportBoundingBox{
					{TagName: "tag1", BoundingBox: &BoundingBox{ID: 2, TagID: 3}},
				},
			},
		},
		{
			name: "box has only tag id",
			args: args{json: `{"id": 1, "boundingBoxes": [{"tagID": 3}]}`},
			want: &ImportAsset{
				Asset:         &Asset{ID: 1},
				BoundingBoxes: []*ImportBoundingBox{{BoundingBox: &BoundingBox{TagID: 3}}},
			},
		},
		{
			name: "box has only tag name",
			args: args{json: `{"id": 1, "boundingBoxes": [{"tagName": "tag1"}]}`},
			want: &ImportAsset{
				Asset:         &Asset{ID: 1},
				BoundingBoxes: []*ImportBoundingBox{{TagName: "tag1"}},
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
			got, err := NewImportAssetFromJson([]byte(tt.args.json))
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
