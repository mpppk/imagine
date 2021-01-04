package cmd_test

import (
	"testing"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
)

func TestBoXAdd(t *testing.T) {
	cases := []struct {
		name         string
		dbName       string
		wsName       model.WSName
		importTags   []*model.Tag
		importAssets []*model.ImportAsset
		command      string
		stdInText    string
		want         string
		wantAssets   []*model.Asset
	}{
		{
			dbName: "box_add_test.imagine",
			wsName: "default-workspace",
			importAssets: []*model.ImportAsset{
				{Asset: model.NewAssetFromFilePath("path1")},
				{Asset: model.NewAssetFromFilePath("path2")},
				{Asset: model.NewAssetFromFilePath("path3")},
			},
			importTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1},
				}},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			command: `box add`, want: "",
			stdInText: `{"path": "path1", "boundingBoxes": [{"tagName": "tag1"}]}`,
		},
		{
			name:   "add same id box",
			dbName: "box_add_test.imagine",
			wsName: "default-workspace",
			importAssets: []*model.ImportAsset{
				{Asset: &model.Asset{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1, X: 5}},
				}},
				{Asset: model.NewAssetFromFilePath("path1")},
				{Asset: model.NewAssetFromFilePath("path2")},
				{Asset: model.NewAssetFromFilePath("path3")},
			},
			importTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1},
					{TagID: 1, X: 5},
				}},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			command: `box add`, want: "",
			stdInText: `{"path": "path1", "boundingBoxes": [{"tagName": "tag1"}]}`,
		},
		{
			name:   "skip if same box is added",
			dbName: "box_add_test.imagine",
			wsName: "default-workspace",
			importAssets: []*model.ImportAsset{
				{Asset: &model.Asset{ID: 1, Name: "path1", Path: "path1",
					BoundingBoxes: []*model.BoundingBox{{TagID: 1}},
				}},
				{Asset: model.NewAssetFromFilePath("path1")},
				{Asset: model.NewAssetFromFilePath("path2")},
				{Asset: model.NewAssetFromFilePath("path3")},
			},
			importTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 1},
				}},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
			command: `box add`, want: "",
			stdInText: `{"path": "path1", "boundingBoxes": [{"tagName": "tag1"}]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := usecasetest.NewTestUseCaseUser(t, c.dbName, c.wsName)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Asset.AddOrUpdateImportAssets(c.wsName, c.importAssets)
				usecases.Tag.SetTags(c.wsName, c.importTags)
			})
			rootCmd := newRootCmd(t)
			cmdWithFlag := c.command + " --db " + c.dbName
			testutil.ExecuteCommand(t, rootCmd, cmdWithFlag, c.stdInText)

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.ListBy(c.wsName, func(a *model.Asset) bool { return true })
				testutil.Diff(t, assets, c.wantAssets)
			})
		})
	}
}
