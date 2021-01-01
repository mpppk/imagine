package usecase_test

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/usecase"

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

			a := usecase.NewAsset(repo, nil)
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
			a := usecase.NewAsset(assetRepo, tagRepo)

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

// TODO: テストが落ちてるので直すところから
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
		existAssets []*model.Asset
		existTags   []*model.Tag
		want        []model.AssetID
		wantAssets  []*model.Asset
		wantErr     bool
	}{
		{
			dbName: "TestAsset_AddOrUpdateImportAssets_add_box.db",
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
			tagSet: model.NewTagSet([]*model.Tag{{ID: 1, Name: "tag1"}}),
			existAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1"},
				{ID: 2, Path: "path2", Name: "path2"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := usecasetest.NewTestUseCaseUser(t, tt.dbName, tt.args.ws)
			defer u.RemoveDB()
			u.Use(func(tu *usecasetest.UseCases) {
				tu.Client.Asset.BatchAdd(tt.args.ws, tt.existAssets)
				tu.Tag.SetTags(tt.args.ws, tt.existTags)
			})

			usecases, err := registry.NewBoltUseCasesWithDBPath(tt.dbName)
			if err != nil {
				t.Fatalf("failed to create usecases instance: %v", err)
			}

			err = usecases.Asset.AddOrUpdateImportAssets(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrUpdateImportAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			//testutil.Diff(t, updatedIDList, tt.want) // FIXME

			if err := usecases.Close(); err != nil {
				t.Fatalf("failed to close db: %v", err)
			}

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.List(tt.args.ws)
				testutil.Diff(t, tt.wantAssets, assets)
			})
		})
	}
}
