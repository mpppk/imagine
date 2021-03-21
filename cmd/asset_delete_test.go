package cmd_test

import (
	"testing"

	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
)

func TestAssetDelete(t *testing.T) {
	cases := []struct {
		name          string
		wsName        model.WSName
		existTagNames []string
		existAssets   []*model.ImportAsset
		stdInText     string
		command       string
		want          string
		wantAssets    []*model.Asset
	}{
		{
			name:   "delete tag",
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				model.NewImportAsset(0, "path1", []*model.ImportBoundingBox{
					{TagName: "tag1"}, {TagName: "tag2"}, {TagName: "tag3"},
				}),
				model.NewImportAsset(0, "path2", []*model.ImportBoundingBox{
					{TagName: "tag2"}, {TagName: "tag3"},
				}),
				model.NewImportAsset(0, "path3", []*model.ImportBoundingBox{
					{TagName: "tag1"}, {TagName: "tag3"},
				}),
			},
			existTagNames: []string{"tag1", "tag2", "tag3"},
			command:       `tag delete --query equals=tag1`,
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2}, {TagID: 3},
				}},
				{ID: 2, Name: "path2", Path: "path2", BoundingBoxes: []*model.BoundingBox{
					{TagID: 2}, {TagID: 3},
				}},
				{ID: 3, Name: "path3", Path: "path3", BoundingBoxes: []*model.BoundingBox{
					{TagID: 3},
				}},
			},
		},
		{
			name:   "delete tags",
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				model.NewImportAsset(0, "path1", []*model.ImportBoundingBox{
					{TagName: "tag1"}, {TagName: "tag12"}, {TagName: "tag3"},
				}),
				model.NewImportAsset(0, "path2", []*model.ImportBoundingBox{
					{TagName: "tag12"},
				}),
				model.NewImportAsset(0, "path3", []*model.ImportBoundingBox{
					{TagName: "tag1"}, {TagName: "tag3"},
				}),
			},
			existTagNames: []string{"tag1", "tag2", "tag3"},
			command:       `tag delete --query start-with=tag1`,
			wantAssets: []*model.Asset{
				{ID: 1, Name: "path1", Path: "path1", BoundingBoxes: []*model.BoundingBox{
					{TagID: 3},
				}},
				{ID: 2, Name: "path2", Path: "path2", BoundingBoxes: []*model.BoundingBox{}},
				{ID: 3, Name: "path3", Path: "path3", BoundingBoxes: []*model.BoundingBox{
					{TagID: 3},
				}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := usecasetest.NewTestUseCaseUser(t, c.wsName)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Tag.SetTagByNames(c.wsName, c.existTagNames)
				usecases.Asset.SaveImportAssets(c.wsName, c.existAssets, nil)
			})

			cmdWithFlag := c.command + " --db " + u.DBPath
			testutil.ExecuteCommand(t, newRootCmd(t), cmdWithFlag, c.stdInText)

			u.Use(func(usecases *usecasetest.UseCases) {
				assets := usecases.Client.Asset.List(c.wsName)
				for _, asset := range assets {
					testutil.SortBoundingBoxesByTagID(asset)
				}
				testutil.Diff(t, c.wantAssets, assets)
			})
		})
	}
}
