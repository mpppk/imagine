package usecase

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/testutil"
	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/google/go-cmp/cmp"

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
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().Update(gomock.Eq(testWSName), gomock.Any()).Return(nil)
			repo.EXPECT().Get(gomock.Eq(testWSName), gomock.Eq(tt.existAsset.ID)).Return(tt.existAsset, true, nil)

			a := &Asset{assetRepository: repo}
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

func TestAsset_AppendBoundingBoxes(t *testing.T) {
	type args struct {
		ws     model.WSName
		assets []*model.ImportAsset
		cap    int
	}
	tests := []struct {
		name        string
		args        args
		tagSet      *model.TagSet
		existAssets []*model.Asset
		want        []model.AssetID
		wantAssets  []*model.Asset
		wantErr     bool
	}{
		{
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: model.NewAssetFromFilePath("path1"),
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag1"},
						},
					},
				},
				cap: 100,
			},
			tagSet: model.NewTagSet([]*model.Tag{{ID: 1, Name: "tag1"}}),
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			want:    []model.AssetID{1},
			wantErr: false,
		},
	}

	ctrl := gomock.NewController(t)
	assetRepo := repoimpl.NewMockAsset(ctrl)
	tagRepo := repoimpl.NewMockTag(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tagRepo.EXPECT().ListAsSet(gomock.Eq(testWSName)).Return(tt.tagSet, nil)
			assetRepo.EXPECT().BatchAppendBoundingBoxes(gomock.Eq(testWSName), gomock.Any()).Return(tt.want, nil)
			a := NewAsset(assetRepo, tagRepo)

			got, err := a.AppendBoundingBoxes(tt.args.ws, tt.args.assets, tt.args.cap)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppendBoundingBoxes() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}

func TestAsset_AddOrUpdateImportAssets(t *testing.T) {
	type args struct {
		ws     model.WSName
		assets []*model.ImportAsset
		cap    int
	}
	tests := []struct {
		name        string
		dbName      string
		args        args
		tagSet      *model.TagSet
		existAssets []*model.ImportAsset
		existTags   []*model.Tag
		want        []model.AssetID
		wantAssets  []*model.Asset
		wantErr     bool
	}{
		{
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: model.NewAssetFromFilePath("path1"),
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag1"},
						},
					},
				},
				cap: 100,
			},
			tagSet: model.NewTagSet([]*model.Tag{{ID: 1, Name: "tag1"}}),
			existAssets: []*model.ImportAsset{
				model.NewImportAssetFromFilePath("path1"),
				model.NewImportAssetFromFilePath("path2"),
			},
			want:    []model.AssetID{1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecases, err := registry.NewBoltUseCasesWithDBPath(tt.dbName)
			if err != nil {
				t.Fatalf("failed to create usecases instance: %v", err)
			}
			u := usecasetest.NewTestUseCaseUser(t, tt.dbName, tt.args.ws)
			defer u.RemoveDB()
			u.Use(func(tu *usecasetest.UseCases) {
				tu.Asset.AddOrUpdateImportAssets(tt.args.ws, tt.existAssets, 100)
				tu.Tag.SetTags(tt.args.ws, tt.existTags)
			})

			idList, err := usecases.Asset.AddOrUpdateImportAssets(tt.args.ws, tt.args.assets, tt.args.cap)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrUpdateImportAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			testutil.Diff(t, idList, tt.want)

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.List(tt.args.ws)
				testutil.Diff(t, assets, tt.wantAssets)
			})
		})
	}
}
