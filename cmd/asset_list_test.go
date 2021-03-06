package cmd_test

import (
	"testing"

	"github.com/mpppk/imagine/testutil"
	"github.com/mpppk/imagine/usecase/usecasetest"

	"github.com/mpppk/imagine/domain/model"
)

func TestAssetList(t *testing.T) {
	cases := []struct {
		name          string
		wsName        model.WSName
		existTagNames []string
		existAssets   []*model.ImportAsset
		command       string
		want          string
	}{
		{
			wsName: "default-workspace",
			existAssets: []*model.ImportAsset{
				{Asset: model.NewAssetFromFilePath("path1")},
				{Asset: model.NewAssetFromFilePath("path2")},
				{Asset: model.NewAssetFromFilePath("path3")},
			},
			existTagNames: []string{"tag1", "tag2", "tag3"},
			command:       `asset list`,
			want: `{"id":1,"name":"path1","path":"path1","boundingBoxes":null}
{"id":2,"name":"path2","path":"path2","boundingBoxes":null}
{"id":3,"name":"path3","path":"path3","boundingBoxes":null}
`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u := usecasetest.NewTestUseCaseUser(t, c.wsName)
			defer u.RemoveDB()
			u.Use(func(usecases *usecasetest.UseCases) {
				usecases.Asset.SaveImportAssets(c.wsName, c.existAssets, nil)
				usecases.Tag.SetTagByNames(c.wsName, c.existTagNames)
			})

			cmdWithFlag := c.command + " --db " + u.DBPath
			gotOut := testutil.ExecuteCommand(t, newRootCmd(t), cmdWithFlag, "")
			testutil.Diff(t, gotOut, c.want)
		})
	}
}
