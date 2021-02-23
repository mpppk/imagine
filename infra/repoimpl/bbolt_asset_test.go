package repoimpl_test

import (
	"context"
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/infra/repoimpl"

	"github.com/mpppk/imagine/domain/model"
)

func TestBBoltAsset_add(t *testing.T) {
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
		usecasetest.RunParallelWithUseCases(t, tt.name, wsName, func(t *testing.T, ut *usecasetest.UseCases) {
			if _, err := ut.Usecases.Client.Asset.Add(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("AddWithIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, exist, err := ut.Usecases.Client.Asset.Get(wsName, tt.args.asset.ID)
			if err != nil || !exist {
				t.Errorf("failed to get asset: %v: %v", newAsset.ID, err)
			}

			testutil.Diff(t, newAsset, got)
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
		usecasetest.RunParallelWithUseCases(t, tt.name, wsName, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(wsName, tt.existAssets)

			ch, errCh, err := ut.Usecases.Client.Asset.ListByIDListAsync(context.Background(), wsName, tt.args.idList, tt.args.cap)
			if (err != nil) && !tt.wantErr {
				t.Fatalf("unexpected ListByIDListAsync error = %v, wantErr %v", err, tt.wantErr)
			}

			assets, err := consumeAllAssetsFromChan(ch, errCh)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantAssets err: %v, got: %v", tt.wantErr, err)
			}

			testutil.Diff(t, tt.want, assets)
		})
	}
}

func TestBBoltAsset_Update(t *testing.T) {
	var wsName model.WSName = "workspace-for-test"
	type args struct {
		asset *model.Asset
	}
	tests := []struct {
		name        string
		args        args
		existAssets []*model.Asset
		want        *model.Asset
		wantErr     bool
	}{
		{
			existAssets: []*model.Asset{
				{
					Name:          "old",
					Path:          "path/to/1",
					BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(1)},
				},
			},
			args: args{asset: &model.Asset{
				ID:            1,
				Name:          "replaced",
				Path:          "path/to/0",
				BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0)},
			}},
			want: &model.Asset{
				ID:            1,
				Name:          "replaced",
				Path:          "path/to/0",
				BoundingBoxes: []*model.BoundingBox{repoimpl.CreateBoundingBox(0)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, wsName, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(wsName, tt.existAssets)

			if err := ut.Usecases.Client.Asset.Update(wsName, tt.args.asset); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _, err := ut.Usecases.Client.Asset.Get(wsName, tt.args.asset.ID)
			if err != nil {
				t.Errorf("failed to get asset: %v: %v", tt.want.ID, err)
			}

			testutil.Diff(t, tt.want, got)
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
		usecasetest.RunParallelWithUseCases(t, tt.name, wsName, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(wsName, tt.existAssets)

			assets, err := ut.Usecases.Client.Asset.ListByIDList(wsName, tt.args.idList)
			if (err != nil) != tt.wantErr {
				t.Errorf("failed to list assets. error: %v", err)
			}

			testutil.Diff(t, tt.want, assets)
		})
	}
}

func TestBBoltAsset_GetByPath(t *testing.T) {
	var wsName model.WSName = "workspace-for-test"
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		existAssets []*model.Asset
		want        *model.Asset
		wantErr     bool
	}{
		{
			args: args{"path/to/0.png"},
			existAssets: []*model.Asset{
				model.NewAssetFromFilePath("path/to/0.png"),
			},
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
		usecasetest.RunParallelWithUseCases(t, tt.name, wsName, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Client.Asset.BatchAdd(wsName, tt.existAssets)

			got, exist, err := ut.Usecases.Client.Asset.GetByPath(wsName, tt.args.path)
			if err != nil || !exist {
				t.Errorf("failed to get asset by GetByPath. exist: %v error: %v", exist, err)
			}

			testutil.Diff(t, tt.want, got)
		})
	}
}

func TestBBoltAsset_AddByFilePathListIfDoesNotExist(t *testing.T) {
	type args struct {
		ws           model.WSName
		filePathList []string
	}
	tests := []struct {
		name       string
		args       args
		want       []model.AssetID
		wantAssets []*model.Asset
		wantErr    bool
	}{
		{
			args: args{
				ws:           "workspace-for-test",
				filePathList: []string{"0.png", "1.png"},
			},
			want: []model.AssetID{1, 2},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "0", Path: "0.png"},
				{ID: 2, Name: "1", Path: "1.png"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			idList, err := ut.Usecases.Client.Asset.AddByFilePathListIfDoesNotExist(tt.args.ws, tt.args.filePathList)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}

			testutil.Diff(t, tt.want, idList)

			gotAssets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, gotAssets)
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
		existTagNames     []string
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
			existTagNames: []string{"tag1", "tag2", "tag3"},
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
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Asset.AddOrMergeImportAssets(tt.args.ws, tt.existAssets)
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			updatedAssets, skippedAssets, err := ut.Usecases.Client.Asset.BatchUpdateByID(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
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
		name          string
		args          args
		existTagNames []string
		want          []model.AssetID
		wantAssets    []*model.Asset
		wantErr       bool
	}{
		{
			name:          "add assets",
			existTagNames: []string{"tag1", "tag2", "tag3"},
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
			name:          "error if arg asset have ID",
			existTagNames: []string{"tag1", "tag2", "tag3"},
			args: args{
				ws: "workspace-for-test",
				assets: []*model.Asset{
					{ID: 1, Name: "path1", Path: "path1"},
				},
			},
			wantErr: true,
		},
		{
			name:          "error if arg asset does not have ID",
			existTagNames: []string{"tag1", "tag2", "tag3"},
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
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			idList, err := ut.Usecases.Client.Asset.BatchAdd(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
			testutil.Diff(t, tt.want, idList)

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
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
		existTagNames     []string
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
			existTagNames: []string{"tag1", "tag2", "tag3"},
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
		usecasetest.RunParallelWithUseCases(t, tt.name, tt.args.ws, func(t *testing.T, ut *usecasetest.UseCases) {
			ut.Asset.AddOrMergeImportAssets(tt.args.ws, tt.existAssets)
			ut.Tag.SetTagByNames(tt.args.ws, tt.existTagNames)

			updatedAssets, skippedAssets, err := ut.Usecases.Client.Asset.BatchUpdateByPath(tt.args.ws, tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddByFilePathListifDoesNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}

			assets := ut.Client.Asset.List(tt.args.ws)
			testutil.Diff(t, tt.wantAssets, assets)
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
