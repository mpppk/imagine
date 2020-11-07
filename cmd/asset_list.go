package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

	"github.com/spf13/cobra"
)

func newAssetListCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewAssetListCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			assetRepository := repoimpl.NewBBoltAsset(db)
			assetUseCase := usecase.NewAsset(assetRepository)
			assetChan, err := assetUseCase.ListAsync(context.Background(), conf.WorkSpace)
			if err != nil {
				return err
			}
			for asset := range assetChan {
				contents, err := json.Marshal(asset)
				if err != nil {
					return fmt.Errorf("failed to marshal asset to json: %w", err)
				}
				cmd.Println(string(contents))
			}
			return nil
		},
	}

	return cmd, nil
}

func init() {
	// FIXME: fs
	assetListCmd, err := newAssetListCmd(nil)
	if err != nil {
		panic(err)
	}
	assetCmd.AddCommand(assetListCmd)
}
