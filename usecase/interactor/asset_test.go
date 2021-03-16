package interactor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/client"
	"github.com/mpppk/imagine/infra/queryimpl"

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

func newMockBoltTag(t *testing.T) *client.Tag {
	ctrl := gomock.NewController(t)
	repo := repoimpl.NewMockTag(ctrl)
	query := queryimpl.NewMockTag(ctrl)
	return &client.Tag{TagRepository: repo, TagQuery: query}
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
	mockTag := newMockBoltTag(t)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo.EXPECT().Update(gomock.Eq(testWSName), gomock.Any()).Return(nil)
			repo.EXPECT().Get(gomock.Eq(testWSName), gomock.Eq(tt.existAsset.ID)).Return(tt.existAsset, true, nil)

			a := interactor.NewAsset(repo, mockTag)
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

func TestAsset_BatchUpdateByID(t *testing.T) {
	type args struct {
		ws      model.WSName
		assets  []*model.Asset
		queries []*model.Query
	}
	tests := []struct {
		name               string
		args               args
		existAssets        []*model.Asset
		existTagNames      []string
		wantUpdatedAssets  []*model.Asset
		wantFilteredAssets []*model.Asset
		wantSkippedAssets  []*model.Asset
		wantAssets         []*model.Asset
		wantErr            bool
	}{
		{
			name: "update assets which match queries",
			args: args{
				ws: testWSName,
				assets: []*model.Asset{
					{
						ID: 1, Path: "updated-path1", Name: "updated-path1",
						BoundingBoxes: []*model.BoundingBox{{TagID: 2}},
					},
					{ID: 2, Path: "path2", Name: "updated-path2",
						// Even if arg asset have tag which matched query,
						// it will be filtered if asset on DB does not have query matched tag.
						BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
					},
					{ID: 3, Path: "path3", Name: "path3"},
				},
				queries: []*model.Query{
					{Op: model.EqualsQueryOP, Value: "tag1"},
				},
			},
			existTagNames: []string{"tag1", "tag2"},
			existAssets: []*model.Asset{
				{Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1},
				}},
				{Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2},
				}},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Path: "updated-path1", Name: "updated-path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}}},
			},
			wantFilteredAssets: []*model.Asset{
				{ID: 2, Path: "path2", Name: "updated-path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
			},
			wantSkippedAssets: []*model.Asset{
				{ID: 3, Path: "path3", Name: "path3"},
			},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "updated-path1", Path: "updated-path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				},
				{ID: 2, Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2},
				}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(tt.args.ws, tt.existAssets)
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			updatedAssets, filteredAssets, skippedAssets, err := ut.Usecases.Asset.BatchUpdateByID(tt.args.ws, tt.args.assets, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImportAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			testutil.Diff(t, tt.wantUpdatedAssets, updatedAssets)
			testutil.Diff(t, tt.wantFilteredAssets, filteredAssets)
			testutil.Diff(t, tt.wantSkippedAssets, skippedAssets)

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
		})
	}
}

func TestAsset_BatchUpdateByPath(t *testing.T) {
	type args struct {
		ws      model.WSName
		assets  []*model.Asset
		queries []*model.Query
	}
	tests := []struct {
		name               string
		args               args
		existAssets        []*model.Asset
		existTagNames      []string
		wantUpdatedAssets  []*model.Asset
		wantFilteredAssets []*model.Asset
		wantSkippedAssets  []*model.Asset
		wantAssets         []*model.Asset
		wantErr            bool
	}{
		{
			name: "update assets which match queries",
			args: args{
				ws: testWSName,
				assets: []*model.Asset{
					{
						Path: "path1", Name: "path1",
						BoundingBoxes: []*model.BoundingBox{{TagID: 2}},
					},
					{ID: 2, Path: "path2", Name: "updated-path2",
						// Even if arg asset have tag which matched query,
						// it will be filtered if asset on DB does not have query matched tag.
						BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
					},
					{ID: 3, Path: "path3", Name: "path3"},
				},
				queries: []*model.Query{
					{Op: model.EqualsQueryOP, Value: "tag1"},
				},
			},
			existTagNames: []string{"tag1", "tag2"},
			existAssets: []*model.Asset{
				{Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1},
				}},
				{Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2},
				}},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}}},
			},
			wantFilteredAssets: []*model.Asset{
				{ID: 2, Path: "path2", Name: "updated-path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
			},
			wantSkippedAssets: []*model.Asset{
				{ID: 3, Path: "path3", Name: "path3"},
			},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				},
				{ID: 2, Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2},
				}},
			},
			wantErr: false,
		},
		{
			name: "fail if arg asset have invalid ID",
			args: args{
				ws: testWSName,
				assets: []*model.Asset{
					{
						ID:   99,
						Path: "path1", Name: "path1",
						BoundingBoxes: []*model.BoundingBox{{TagID: 2}},
					},
				},
				queries: []*model.Query{
					{Op: model.EqualsQueryOP, Value: "tag1"},
				},
			},
			existTagNames: []string{"tag1", "tag2"},
			existAssets: []*model.Asset{
				{Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1},
				}},
				{Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2},
				}},
			},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
				},
				{ID: 2, Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2},
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(tt.args.ws, tt.existAssets)
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			updatedAssets, filteredAssets, skippedAssets, err := ut.Usecases.Asset.BatchUpdateByPath(tt.args.ws, tt.args.assets, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImportAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				assets := ut.Client.Asset.List(tt.args.ws)
				testutil.Diff(t, tt.wantAssets, assets)
				return
			}

			testutil.Diff(t, tt.wantUpdatedAssets, updatedAssets)
			testutil.Diff(t, tt.wantFilteredAssets, filteredAssets)
			testutil.Diff(t, tt.wantSkippedAssets, skippedAssets)

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
		})
	}
}

func TestAsset_SaveImportAssets(t *testing.T) {
	type args struct {
		ws      model.WSName
		assets  []*model.ImportAsset
		queries []*model.Query
	}
	tests := []struct {
		name              string
		args              args
		existAssets       []*model.Asset
		existTagNames     []string
		wantAddedAssets   []*model.Asset
		wantUpdatedAssets []*model.Asset
		wantSkippedAssets []*model.Asset
		wantAssets        []*model.Asset
		wantErr           bool
	}{
		{
			name: "add assets",
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					model.NewImportAsset(0, "path1", []*model.ImportBoundingBox{
						model.NewImportBoundingBoxFromTagID(1),
						model.NewImportBoundingBoxFromTagID(2),
					}),
					model.NewImportAsset(0, "path2", []*model.ImportBoundingBox{
						model.NewImportBoundingBoxFromTagID(1),
						model.NewImportBoundingBoxFromTagID(3),
					}),
					model.NewImportAsset(0, "path3", []*model.ImportBoundingBox{
						model.NewImportBoundingBoxFromTagID(2),
						model.NewImportBoundingBoxFromTagID(3),
					}),
				},
			},
			existTagNames: []string{"tag1", "tag2xxx", "xxxtag2", "tag4"},
			existAssets:   []*model.Asset{},
			wantAddedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				},
				{ID: 2, Name: "path2", Path: "path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 3}},
				},
				{ID: 3, Name: "path3", Path: "path3",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}, {TagID: 3}},
				},
			},
			wantUpdatedAssets: nil,
			wantSkippedAssets: nil,
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}},
				},
				{ID: 2, Name: "path2", Path: "path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 3}},
				},
				{ID: 3, Name: "path3", Path: "path3",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}, {TagID: 3}},
				},
			},
			wantErr: false,
		},
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
			existTagNames: []string{"tag1"},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			wantAddedAssets: []*model.Asset{
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
				},
			},
			wantSkippedAssets: nil,
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
			name: "add and update assets which equals and start-with queries",
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: &model.Asset{ID: 1, Path: "path1", Name: "path1"},
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag4"},
						},
					},
					{Asset: &model.Asset{Path: "path2", Name: "updated-path2"}},
					{Asset: &model.Asset{Path: "path3", Name: "path3"}},
				},
				queries: []*model.Query{
					{Op: model.EqualsQueryOP, Value: "tag1"},
					{Op: model.StartWithQueryOP, Value: "tag2"},
				},
			},
			existTagNames: []string{"tag1", "tag2xxx", "xxxtag2", "tag4"},
			existAssets: []*model.Asset{
				{Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}}},
				{Name: "path2", Path: "path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 3}}},
				{Name: "path3", Path: "path3",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}, {TagID: 3}}},
			},
			wantAddedAssets: nil,
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}, {TagID: 4}}},
			},
			wantSkippedAssets: []*model.Asset{
				{Path: "path2", Name: "updated-path2"},
				{Name: "path3", Path: "path3"},
			},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}, {TagID: 4}}},
				{ID: 2, Name: "path2", Path: "path2",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 3}}},
				{ID: 3, Name: "path3", Path: "path3",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}, {TagID: 3}}},
			},
			wantErr: false,
		},
		{
			name: "add and update assets which match queries",
			args: args{
				ws: testWSName,
				assets: []*model.ImportAsset{
					{
						Asset: &model.Asset{ID: 1, Path: "path1", Name: "path1"},
						BoundingBoxes: []*model.ImportBoundingBox{
							{TagName: "tag1"},
						},
					},
					{Asset: &model.Asset{Path: "path2", Name: "updated-path2"}},
					{Asset: &model.Asset{Path: "path3", Name: "path3"}},
				},
				queries: []*model.Query{
					{Op: model.PathEqualsQueryOP, Value: "path1"},
				},
			},
			existTagNames: []string{"tag1"},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			wantAddedAssets: []*model.Asset{
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
				},
			},
			wantSkippedAssets: []*model.Asset{
				{Path: "path2", Name: "updated-path2"},
			},
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
			existTagNames: []string{"tag1"},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			wantAddedAssets: []*model.Asset{
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 2}},
				},
			},
			wantSkippedAssets: nil,
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
			existTagNames: []string{"tag1", "tag2"},
			existAssets: []*model.Asset{
				{Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2},
				}},
			},
			wantAddedAssets: nil,
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}, {TagID: 2}, {TagID: 1, X: 1}},
				},
			},
			wantSkippedAssets: nil,
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
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			addedAssets, updatedAssets, skippedAssets, err := ut.Usecases.Asset.SaveImportAssets(tt.args.ws, tt.args.assets, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImportAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				assets := ut.Client.Asset.List(tt.args.ws)
				testutil.Diff(t, tt.wantAssets, assets)
				return
			}

			testutil.Diff(t, tt.wantAddedAssets, addedAssets)
			testutil.Diff(t, tt.wantUpdatedAssets, updatedAssets)
			testutil.Diff(t, tt.wantSkippedAssets, skippedAssets)

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
		})
	}
}

func TestAsset_ListAsyncByQueries(t *testing.T) {
	type args struct {
		ctx     context.Context
		ws      model.WSName
		queries []*model.Query
	}
	tests := []struct {
		name          string
		args          args
		existAssets   []*model.Asset
		existTagNames []string
		wantAssets    []*model.Asset
		wantErr       bool
	}{
		{
			name: "list",
			args: args{
				ctx:     context.Background(),
				ws:      "default-workspace",
				queries: []*model.Query{},
			},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
			},
			existTagNames: []string{"tag1", "tag2"},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1"},
				{ID: 2, Name: "path2", Path: "path2"},
			},
			wantErr: false,
		},
		{
			name: "list with equals query",
			args: args{
				ctx:     context.Background(),
				ws:      "default-workspace",
				queries: []*model.Query{{Op: "equals", Value: "tag1"}},
			},
			existAssets: []*model.Asset{
				{Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2}, {TagID: 3},
				}},
				{Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2}, {TagID: 3},
				}},
			},
			existTagNames: []string{"tag1", "tag2", "tag3"},
			wantAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2}, {TagID: 3},
				}},
			},
			wantErr: false,
		},
		{
			name: "list with not-start-with query",
			args: args{
				ctx:     context.Background(),
				ws:      "default-workspace",
				queries: []*model.Query{{Op: "not-start-with", Value: "tag1"}},
			},
			existAssets: []*model.Asset{
				{Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2}, {TagID: 3},
				}},
				{Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2}, {TagID: 3},
				}},
				{Path: "path3", Name: "path3", BoundingBoxes: []*model.BoundingBox{
					{TagID: 3},
				}},
				{Path: "path4", Name: "path4", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2},
				}},
			},
			existTagNames: []string{"tag1", "tag11", "tag2"},
			wantAssets: []*model.Asset{
				// path1 and path2 should match because they have tag which not start with "tag1"(in other words, "tag2")
				{ID: 1, Path: "path1", Name: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1}, {TagID: 2}, {TagID: 3},
				}},
				{ID: 2, Path: "path2", Name: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2}, {TagID: 3},
				}},
				{ID: 3, Path: "path3", Name: "path3", BoundingBoxes: []*model.BoundingBox{
					{TagID: 3},
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)
			ut.Client.Asset.BatchAdd(tt.args.ws, tt.existAssets)

			assetCh, err := ut.Usecases.Asset.ListAsyncByQueries(tt.args.ctx, tt.args.ws, tt.args.queries)
			gotAssets := testutil.ReadAllAssetsFromCh(assetCh)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAsyncByQueries() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			testutil.Diff(t, tt.wantAssets, gotAssets)
		})
	}
}
