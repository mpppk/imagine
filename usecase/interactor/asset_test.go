package interactor_test

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/usecase/interactor"

	"github.com/mpppk/imagine/testutil"
	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/golang/mock/gomock"

	"github.com/mpppk/imagine/infra/repoimpl"

	"github.com/mpppk/imagine/domain/model"
)

var testWSName model.WSName = "test-ws"

func newBox(boxId model.BoundingBoxID, tagId model.TagID, tagName string) *model.BoundingBox {
	return &model.BoundingBox{
		ID:    boxId,
		TagID: tagId,
	}
}

func TestAsset_AssignBoundingBox(t *testing.T) {
	type args struct {
		ws      model.WSName
		assetId model.AssetID
		box     *model.BoundingBox
	}
	tests := []struct {
		name       string
		args       args
		existAsset *model.Asset
		want       *model.Asset
		wantErr    bool
	}{
		{
			args: args{
				ws:      testWSName,
				assetId: 0,
				box:     newBox(0, 0, "test-tag"),
			},
			existAsset: &model.Asset{
				ID:            0,
				Name:          "test",
				Path:          "test/path",
				BoundingBoxes: nil,
			},
			want: &model.Asset{
				ID:   0,
				Name: "test",
				Path: "test/path",
				BoundingBoxes: []*model.BoundingBox{
					newBox(1, 0, "test-tag"),
				},
			},
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	repo := repoimpl.NewMockAsset(ctrl)
	repo.EXPECT().Init(gomock.Eq(testWSName)).Return(nil)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo.EXPECT().Update(gomock.Eq(testWSName), gomock.Any()).Return(nil)
			repo.EXPECT().Get(gomock.Eq(testWSName), gomock.Eq(tt.existAsset.ID)).Return(tt.existAsset, true, nil)

			a := interactor.NewAsset(repo, nil)
			got, err := a.AssignBoundingBox(tt.args.ws, tt.args.assetId, tt.args.box)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssignBoundingBox() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AssignBoundingBox() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAsset_AddOrMergeImportAssets(t *testing.T) {
	type args struct {
		ws     model.WSName
		assets []*model.ImportAsset
	}
	tests := []struct {
		name        string
		args        args
		existAssets []*model.Asset
		existTags   []*model.Tag
		want        []model.AssetID
		wantAssets  []*model.Asset
		wantErr     bool
	}{
		{
			name: "add and update assets",
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: &model.Asset{ID: 1, Path: "path1", Name: "path1"},
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag1"},
						},
					},
					{Asset: &model.Asset{Path: "path3", Name: "path3"}},
				},
			},
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			want: []model.AssetID{1},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
				},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantErr: false,
		},
		{
			name: "refer a tag that doesn't exist",
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: &model.Asset{ID: 1, Path: "path1", Name: "path1"},
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag2"},
						},
					},
					{Asset: &model.Asset{Path: "path3", Name: "path3"}},
				},
			},
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			want: []model.AssetID{1},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}},
				},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantErr: false,
		},
		{
			name: "append boxes",
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: &model.Asset{ID: 1, Path: "path1", Name: "path1"},
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag1", BoundingBox: &model.BoundingBox{X: 1}},
							{TagName: "tag2"},
						},
					},
				},
			},
			existTags: []*model.Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}},
			existAssets: []*model.Asset{
				{Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2},
				}},
			},
			want: []model.AssetID{1},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}, {TagID: 1, X: 1}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(tt.args.ws, tt.existAssets)
			ut.Tag.SetTags(tt.args.ws, tt.existTags)

			err := ut.Usecases.Asset.AddOrMergeImportAssets(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrMergeImportAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
		})
	}
}
