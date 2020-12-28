package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	bolt "go.etcd.io/bbolt"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/cmd"

	"github.com/spf13/cobra"
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
			defer teardown(t, c.dbName)

			if err := addAssets(c.dbName, c.wsName, c.importAssets); err != nil {
				t.Fatal(err)
			}

			if err := addTags(c.dbName, c.wsName, c.importTags); err != nil {
				t.Fatal(err)
			}

			reader := strings.NewReader(c.stdInText)
			cmd.RootCmd.SetIn(reader)
			cmdWithFlag := c.command + " --db " + c.dbName
			if _, err := executeCommand(cmd.RootCmd, cmdWithFlag); err != nil {
				t.Errorf("failed to execute box add command: %v", err)
			}

			assets, err := listAssets(c.dbName, c.wsName)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(assets, c.wantAssets); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}

func addTags(dbPath string, ws model.WSName, tags []*model.Tag) error {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open db file from %v: %v", dbPath, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	usecases := registry.NewBoltUseCases(db)
	if err := usecases.Init(ws); err != nil {
		return fmt.Errorf("failed to initialize workspace(%s): %w", ws, err)
	}
	return usecases.Tag.SetTags(ws, tags)
}

func addAssets(dbPath string, ws model.WSName, importAssets []*model.ImportAsset) error {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open db file from %v: %v", dbPath, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	usecases := registry.NewBoltUseCases(db)
	if err := usecases.Init(ws); err != nil {
		return fmt.Errorf("failed to initialize workspace(%s): %w", ws, err)
	}
	_, err = usecases.Asset.AddImportAssets(ws, importAssets, 100)
	return err
}

func listAssets(dbPath string, ws model.WSName) ([]*model.Asset, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db file from %v: %v", dbPath, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	usecases := registry.NewBoltUseCases(db)
	if err := usecases.Init(ws); err != nil {
		return nil, fmt.Errorf("failed to initialize workspace(%s): %w", ws, err)
	}
	assets, err := usecases.Client.Asset.ListBy(ws, func(a *model.Asset) bool { return true })
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	return assets, nil
}

func executeCommand(cmd *cobra.Command, command string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmdArgs := strings.Split(command, " ")
	cmd.SetArgs(cmdArgs)
	if err := cmd.Execute(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func teardown(t *testing.T, fileName string) {
	t.Helper()
	if err := os.Remove(fileName); err != nil {
		t.Errorf("failed to remove test file: %v", err)
	}
}
