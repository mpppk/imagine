package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newAssetDeleteCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewAssetDeleteCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			usecases, err := registry.NewBoltUseCasesWithDBPath(conf.DB)
			if err != nil {
				return fmt.Errorf("failed to create usecases instance: %w", err)
			}
			defer func() {
				if err := usecases.Close(); err != nil {
					panic(err)
				}
			}()
			scanner := bufio.NewScanner(os.Stdin)
			var asset model.Asset
			for scanner.Scan() {
				if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
					return fmt.Errorf("failed to unmarshal json to asset")
				}
				if asset.ID == 0 {
					log.Printf("warning: missing ID: %#v", asset)
					continue
				}

				if ok, err := usecases.Client.Asset.Has(conf.WorkSpace, asset.ID); err != nil {
					return fmt.Errorf("failed to check asset. image path: %s: %w", asset.Path, err)
				} else if !ok {
					log.Printf("debug: asset skipped because it does not exist: id:%d", asset.ID)
					continue
				}

				if err := usecases.Client.Asset.Delete(conf.WorkSpace, asset.ID); err != nil {
					return fmt.Errorf("failed to add asset. image path: %s: %w", asset.Path, err)
				}
				log.Printf("debug: asset deleted: %#v", asset)
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("faield to scan asset op: %w", err)
			}
			return nil
		},
	}

	return cmd, nil
}

func init() {
	// FIXME: fs
	assetListCmd, err := newAssetDeleteCmd(nil)
	if err != nil {
		panic(err)
	}
	assetCmd.AddCommand(assetListCmd)
}
