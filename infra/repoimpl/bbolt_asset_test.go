package repoimpl_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/google/go-cmp/cmp"

	"github.com/mpppk/imagine/infra/repoimpl"

	"github.com/mpppk/imagine/domain/model"
)

func TestBBoltAsset_add(t *testing.T) {
	fileName := "TestBBoltAsset_Add.db"
	var wsName model.WSName = "workspace-for-test"
	newAsset := &model.Asset{
		ID:            0,
		Name:          "replaced",
		Path:          "path/to/0",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0)},
	}
	type args struct {
		asset *model.Asset
	}
	tests := []struct {
		name string
		args args

		wantErr bool
	}{
		{
			args:    args{asset: newAsset},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCases(t, fileName, wsName)
			defer closer()
			defer remover()

			if _, err := usecases.Client.Asset.Add(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, exist, err := usecases.Client.Asset.Get(wsName, tt.args.asset.ID)
			if err != nil || !exist {
				t.Errorf("failed to get asset: %v: %v", newAsset.ID, err)
			}
			if !reflect.DeepEqual(got, newAsset) {
				t.Errorf("wantAssets: %#v, got: %#v", newAsset, got)
			}
		})
	}
}

func TestBBoltAsset_ListByIDListAsync(t *testing.T) {
	var wsName model.WSName = "workspace-for-test"
	type args struct {
		idList []model.AssetID
		cap    int
	}
	tests := []struct {
		name        string
		existAssets []*model.Asset
		args        args

		want    []*model.Asset
		wantErr bool
	}{
		{
			existAssets: []*model.Asset{},
			args:        args{idList: []model.AssetID{1}, cap: 100},
			wantErr:     true,
		},
		{
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
				model.NewAssetFromFilePath("path3"),
			},
			args: args{idList: []model.AssetID{1, 3}, cap: 100},
			want: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantErr: false,
		},
		{
			name: "cap 1",
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
				model.NewAssetFromFilePath("path3"),
			},
			args: args{idList: []model.AssetID{1, 3}, cap: 1},
			want: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCasesWithTempDB(t, wsName)
			defer closer()
			defer remover()

			if _, err := usecases.Client.Asset.BatchAdd(wsName, tt.existAssets); err != nil {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			ch, errCh, err := usecases.Client.Asset.ListByIDListAsync(context.Background(), wsName, tt.args.idList, tt.args.cap)
			if (err != nil) && !tt.wantErr {
				t.Fatalf("unexpected ListByIDListAsync error = %v, wantErr %v", err, tt.wantErr)
			}

			assets, err := consumeAllAssetsFromChan(ch, errCh)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantAssets err: %v, got: %v", tt.wantErr, err)
			}
			if diff := cmp.Diff(assets, tt.want); diff != "" {
				t.Errorf("(-got +wantAssets)\n%s", diff)
			}
		})
	}
}

func TestBBoltAsset_Update(t *testing.T) {
	fileName := "TestBBoltAsset_Update.db"
	var wsName model.WSName = "workspace-for-test"
	oldAsset := &model.Asset{
		ID:            0,
		Name:          "old",
		Path:          "path/to/1",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(1)},
	}
	newAsset := &model.Asset{
		ID:            0,
		Name:          "replaced",
		Path:          "path/to/0",
		BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0)},
	}
	type args struct {
		asset *model.Asset
	}
	tests := []struct {
		name      string
		args      args
		oldAssets []*model.Asset
		wantErr   bool
	}{
		{
			oldAssets: []*model.Asset{oldAsset},
			args:      args{asset: newAsset},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCases(t, fileName, wsName)
			defer closer()
			defer remover()

			for _, asset := range tt.oldAssets {
				if _, _, err := usecases.Client.Asset.AddByFilePathIfDoesNotExist(wsName, asset.Path); err != nil {
					t.Fatalf("failed to addWithID assets: %v: %#v", err, asset)
				}
			}

			if err := usecases.Client.Asset.Update(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _, err := usecases.Client.Asset.Get(wsName, tt.args.asset.ID)
			if err != nil {
				t.Errorf("failed to get asset: %v: %v", newAsset.ID, err)
			}
			if !reflect.DeepEqual(got, newAsset) {
				t.Errorf("wantAssets: %#v, got: %#v", newAsset, got)
			}
		})
	}
}

func TestBBoltAsset_ListByIDList(t *testing.T) {
	var wsName model.WSName = "workspace-for-test"
	type args struct {
		idList []model.AssetID
	}
	tests := []struct {
		name        string
		args        args
		existAssets []*model.Asset
		want        []*model.Asset
		wantErr     bool
	}{
		{
			name:        "return nil if ID does not exist",
			args:        args{[]model.AssetID{1}},
			existAssets: []*model.Asset{},
			want:        []*model.Asset{nil},
		},
		{
			args: args{[]model.AssetID{1, 3}},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
			},
			want: []*model.Asset{{ID: 1, Name: "path1", Path: "path1"}, nil},
		},
		{
			args: args{[]model.AssetID{1, 3}},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path1"),
				model.NewAssetFromFilePath("path2"),
				model.NewAssetFromFilePath("path3"),
			},
			want: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCasesWithTempDB(t, wsName)
			defer closer()
			defer remover()

			if _, err := usecases.Client.Asset.BatchAdd(wsName, tt.existAssets); err != nil {
				t.Fatalf("BatchAdd() error = %v, wantErr %v", err, tt.wantErr)
			}

			assets, err := usecases.Client.Asset.ListByIDList(wsName, tt.args.idList)
			if (err != nil) != tt.wantErr {
				t.Errorf("failed to list assets. error: %v", err)
			}
			if diff := cmp.Diff(assets, tt.want); diff != "" {
				t.Errorf("(-got +wantAssets)\n%s", diff)
			}
		})
	}
}

func TestBBoltAsset_GetByPath(t *testing.T) {
	fileName := "TestBBoltAsset_GetByPath.db"
	var wsName model.WSName = "workspace-for-test"
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		existPaths []string
		want       *model.Asset
		wantErr    bool
	}{
		{
			args:       args{"path/to/0.png"},
			existPaths: []string{"path/to/0.png"},
			want: &model.Asset{
				ID:   1,
				Name: "0",
				Path: "path/to/0.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCases(t, fileName, wsName)
			defer closer()
			defer remover()

			for _, path := range tt.existPaths {
				if _, ok, err := usecases.Client.Asset.AddByFilePathIfDoesNotExist(wsName, path); (err != nil) != tt.wantErr {
					t.Errorf("GetByPath() error whiile assets adding. error = %v, wantErr %v", err, tt.wantErr)
				} else if !ok {
					t.Errorf("failed to add asset in AddByFilePathIfDoesNotExist test exist: %v", ok)
				}
			}
			_, err := usecases.Client.Asset.ListBy(wsName, func(a *model.Asset) bool { return true })
			if err != nil {
				t.Errorf("failed to list assets. error: %v", err)
			}
			got, exist, err := usecases.Client.Asset.GetByPath(wsName, tt.args.path)
			if err != nil || !exist {
				t.Errorf("failed to get asset by GetByPath. exist: %v error: %v", exist, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wantAssets: %#v, got: %#v", tt.want, got)
			}
		})
	}
}

func TestBBoltAsset_AddByFilePathListIfDoesNotExist(t *testing.T) {
	fileName := "TestBBoltAsset_AddByFlePathListIfDoesNotExist.db"
	type args struct {
		ws           model.WSName
		filePathList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.AssetID
		wantErr bool
	}{
		{
			args: args{
				ws:           "workspace-for-test",
				filePathList: []string{"0.png", "1.png"},
			},
			want: []model.AssetID{1, 2},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			usecases, closer, remover := usecasetest.SetUpUseCases(t, fileName, tt.args.ws)
			defer closer()
			defer remover()

			idList, err := usecases.Client.Asset.AddByFilePathListIfDoesNotExist(tt.args.ws, tt.args.filePathList)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, id := range idList {
				asset, exist, err := usecases.Client.Asset.Get(tt.args.ws, id)
				if err != nil || !exist {
					t.Errorf("failed to get asset which added by AddByFilePathListifDoesNotExist() exist = %v error = %v", exist, err)
				}
				asset2, exist, err := usecases.Client.Asset.GetByPath(tt.args.ws, asset.Path)
				if err != nil || !exist {
					t.Errorf("failed to get asset which added by AddByFilePathListifDoesNotExist() exist = %v error = %v", exist, err)
				}

				if !reflect.DeepEqual(asset, asset2) {
					t.Errorf("inconsistecy detected on AddByFilePathListifDoesNotExisti() asset1: %#v, asset2: %#v", asset, asset2)
				}
			}
		})
	}
}

func TestBBoltAsset_BatchUpdateByID(t *testing.T) {
	type args struct {
		ws     model.WSName
		assets []*model.Asset
	}
	tests := []struct {
		name              string
		args              args
		existAssets       []*model.ImportAsset
		existTags         []*model.Tag
		wantUpdatedAssets []*model.Asset
		wantSkippedAssets []*model.Asset
		wantAssets        []*model.Asset
		wantErr           bool
	}{
		{
			existAssets: []*model.ImportAsset{
				model.NewImportAssetFromFilePath("path1"),
				model.NewImportAssetFromFilePath("path2"),
				model.NewImportAssetFromFilePath("path3"),
			},
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			args: args{
				ws: "workspace-for-test",
				assets: []*model.Asset{
					{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{ID: 1, TagID: 1}}},
					{Name: "path2", Path: "path2"}, // does not have ID
				},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{ID: 1, TagID: 1}}},
			},
			wantSkippedAssets: []*model.Asset{
				{Name: "path2", Path: "path2"},
			},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{ID: 1, TagID: 1}}},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := usecasetest.NewTestUseCaseUser(t, tt.args.ws)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Asset.AddOrMergeImportAssets(tt.args.ws, tt.existAssets)
				usecases.Tag.SetTags(tt.args.ws, tt.existTags)
			})

			usecases, closeDB, _ := usecasetest.SetUpUseCases(t, u.DBPath, tt.args.ws)
			defer closeDB()

			updatedAssets, skippedAssets, err := usecases.Client.Asset.BatchUpdateByID(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}

			closeDB()

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.ListBy(tt.args.ws, func(a *model.Asset) bool { return true })
				testutil.Diff(t, tt.wantAssets, assets)
			})

			testutil.Diff(t, tt.wantUpdatedAssets, updatedAssets)
			testutil.Diff(t, tt.wantSkippedAssets, skippedAssets)
		})
	}
}

func TestBBoltAsset_BatchAdd(t *testing.T) {
	type args struct {
		ws     model.WSName
		assets []*model.Asset
	}
	tests := []struct {
		name       string
		args       args
		existTags  []*model.Tag
		want       []model.AssetID
		wantAssets []*model.Asset
		wantErr    bool
	}{
		{
			name:      "add assets",
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			args: args{
				ws: "workspace-for-test",
				assets: []*model.Asset{
					{Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
					{Name: "path2", Path: "path2"}, // does not have ID but have Path
				},
			},
			want: []model.AssetID{1, 2},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
				{ID: 2, Name: "path2", Path: "path2"},
			},
		},
		{
			name:      "error if arg asset have ID",
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			args: args{
				ws: "workspace-for-test",
				assets: []*model.Asset{
					{ID: 1, Name: "path1", Path: "path1"},
				},
			},
			wantErr: true,
		},
		{
			name:      "error if arg asset does not have ID",
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			args: args{
				ws: "workspace-for-test",
				assets: []*model.Asset{
					{Name: "path1"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := usecasetest.NewTestUseCaseUser(t, tt.args.ws)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Tag.SetTags(tt.args.ws, tt.existTags)
			})

			usecases, closeDB, _ := usecasetest.SetUpUseCases(t, u.DBPath, tt.args.ws)
			defer closeDB()

			idList, err := usecases.Client.Asset.BatchAdd(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}

			closeDB()

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.ListBy(tt.args.ws, func(a *model.Asset) bool { return true })
				testutil.Diff(t, tt.wantAssets, assets)
			})

			testutil.Diff(t, tt.want, idList)
		})
	}
}

func TestBBoltAsset_BatchUpdateByPath(t *testing.T) {
	type args struct {
		ws     model.WSName
		assets []*model.Asset
	}
	tests := []struct {
		name              string
		args              args
		existAssets       []*model.ImportAsset
		existTags         []*model.Tag
		wantUpdatedAssets []*model.Asset
		wantSkippedAssets []*model.Asset
		wantAssets        []*model.Asset
		wantErr           bool
	}{
		{
			existAssets: []*model.ImportAsset{
				model.NewImportAssetFromFilePath("path1"),
				model.NewImportAssetFromFilePath("path2"),
				model.NewImportAssetFromFilePath("path3"),
			},
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			args: args{
				ws: "workspace-for-test",
				assets: []*model.Asset{
					{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{ID: 1, TagID: 1}}},
					{Name: "path2-update", Path: "path2"}, // does not have ID but have Path
					{Name: "path3-update"},                // does not have ID and Path, so updating will be skipped
				},
			},
			wantUpdatedAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{ID: 1, TagID: 1}}},
				{ID: 2, Name: "path2-update", Path: "path2"}, // does not have ID but have Path
			},
			wantSkippedAssets: []*model.Asset{
				{Name: "path3-update"},
			},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{ID: 1, TagID: 1}}},
				{ID: 2, Name: "path2-update", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u := usecasetest.NewTestUseCaseUser(t, tt.args.ws)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Asset.AddOrMergeImportAssets(tt.args.ws, tt.existAssets)
				usecases.Tag.SetTags(tt.args.ws, tt.existTags)
			})

			usecases, closeDB, _ := usecasetest.SetUpUseCases(t, u.DBPath, tt.args.ws)
			defer closeDB()

			updatedAssets, skippedAssets, err := usecases.Client.Asset.BatchUpdateByPath(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}

			closeDB()

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.ListBy(tt.args.ws, func(a *model.Asset) bool { return true })
				testutil.Diff(t, tt.wantAssets, assets)
			})

			testutil.Diff(t, tt.wantUpdatedAssets, updatedAssets)
			testutil.Diff(t, tt.wantSkippedAssets, skippedAssets)
		})
	}
}

func consumeAllAssetsFromChan(ch <-chan *model.Asset, errCh <-chan error) (assets []*model.Asset, err error) {
	for {
		select {
		case asset, ok := <-ch:
			if !ok {
				return assets, nil
			}
			assets = append(assets, asset)
		case err := <-errCh:
			return nil, err
		}
	}
}
