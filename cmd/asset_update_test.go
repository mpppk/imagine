package cmd_test

import (
	"testing"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
)

func TestAssetUpdate(t *testing.T) {
	cases := []struct {
		name        string
		wsName      model.WSName
		existTags   []*model.Tag
		existAssets []*model.ImportAsset
		stdInText   string
		command     string
		want        string
		wantAssets  []*model.Asset
	}{
		{
			name:   "Do nothing",
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				model.NewImportAssetFromFilePath("path1"),
				model.NewImportAssetFromFilePath("path2"),
				model.NewImportAssetFromFilePath("path3"),
			},
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			command:   `asset update`,
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1"},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
		},
		{
			name:   "find by tag and update bounding box",
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				model.NewImportAssetFromFilePath("path1"),
				model.NewImportAssetFromFilePath("path2"),
				model.NewImportAssetFromFilePath("path3"),
			},
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			command:   `asset update`,
			stdInText: `{"id": 1, "boundingBoxes": [{"tagID": 1}]}`,
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
		},
		{
			name:   "find by path and update bounding box",
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				{Asset: model.NewAssetFromFilePath("path1")},
				{Asset: model.NewAssetFromFilePath("path2")},
				{Asset: model.NewAssetFromFilePath("path3")},
			},
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			command:   `asset update`,
			stdInText: `{"path": "path1", "boundingBoxes": [{"tagID": 1}]}`,
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{{TagID: 1}}},
				{ID: 2, Name: "path2", Path: "path2"},
				{ID: 3, Name: "path3", Path: "path3"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := usecasetest.NewTestUseCaseUser(t, c.wsName)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Asset.AddOrMergeImportAssets(c.wsName, c.existAssets)
				usecases.Tag.SetTags(c.wsName, c.existTags)
			})

			cmdWithFlag := c.command + " --db " + u.DBPath
			testutil.ExecuteCommand(t, newRootCmd(t), cmdWithFlag, c.stdInText)

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.List(c.wsName)
				testutil.Diff(t, c.wantAssets, assets)
			})
		})
	}
}
