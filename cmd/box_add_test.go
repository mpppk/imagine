package cmd_test

import (
	"strings"
	"testing"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/cmd"
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			usecases, teardown := testutil.SetUpUseCases(t, c.dbName, c.wsName)
			defer teardown()

			if _, err := usecases.Asset.AddImportAssets(c.wsName, c.importAssets, 100); err != nil {
				t.Fatalf("failed to get assets: %v", err)
			}

			if err := usecases.Tag.SetTags(c.wsName, c.importTags); err != nil {
				t.Fatalf("failed to set tags: %v", err)
			}

			reader := strings.NewReader(c.stdInText)
			cmd.RootCmd.SetIn(reader)
			cmdWithFlag := c.command + " --db " + c.dbName
			if _, err := testutil.ExecuteCommand(cmd.RootCmd, cmdWithFlag); err != nil {
				t.Errorf("failed to execute box add command: %v", err)
			}

			assets, err := usecases.Client.Asset.ListBy(c.wsName, func(a *model.Asset) bool { return true })
			if err != nil {
				t.Fatalf("faled to list assets: %v", err)
			}

			testutil.Diff(t, assets, c.wantAssets)
		})
	}
}
