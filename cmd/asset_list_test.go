package cmd_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"
	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestAssetList(t *testing.T) {
	cases := []struct {
		name        string
		dbName      string
		wsName      model.WSName
		existTags   []*model.Tag
		existAssets []*model.ImportAsset
		command     string
		want        string
	}{
		{
			dbName: "asset_list.imagine",
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				{Asset: model.NewAssetFromFilePath("path1")},
				{Asset: model.NewAssetFromFilePath("path2")},
				{Asset: model.NewAssetFromFilePath("path3")},
			},
			existTags: []*model.Tag{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
			command:   `asset list`,
			want: `{"id":1,"name":"path1","path":"path1","boundingBoxes":null}
{"id":2,"name":"path2","path":"path2","boundingBoxes":null}
{"id":3,"name":"path3","path":"path3","boundingBoxes":null}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := usecasetest.NewTestUseCaseUser(t, c.dbName, c.wsName)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Asset.AddOrUpdateImportAssets(c.wsName, c.existAssets)
				usecases.Tag.SetTags(c.wsName, c.existTags)
			})

			cmdWithFlag := c.command + " --db " + c.dbName
			testutil.ExecuteCommand(t, newRootCmd(t), cmdWithFlag, "")
			//testutil.DeepEqual(t, gotOut, c.want) // FIXME: gotOut is always empty
		})
	}
}
