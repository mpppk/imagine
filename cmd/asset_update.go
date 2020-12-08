package cmd

import (
	"fmt"
	"os"

	"github.com/mpppk/imagine/registry"

	bolt "go.etcd.io/bbolt"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newAssetUpdateCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewAssetAddCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}

			assetUseCase := registry.InitializeAssetUseCase(db)
			if err := assetUseCase.Init(conf.WorkSpace); err != nil {
				return fmt.Errorf("failed to initialize asset usecase: %w", err)
			}

			if err := assetUseCase.ImportFromReader(conf.WorkSpace, os.Stdin, conf.New); err != nil {
				return fmt.Errorf("failed to import asset from reader: %w", err)
			}

			return nil
		},
	}

	registerFlags := func(cmd *cobra.Command) error {
		flags := []option.Flag{
			&option.BoolFlag{
				BaseFlag: &option.BaseFlag{
					Name:  "new",
					Usage: "If the asset with the specified ID does not exist, create a new one",
				},
				Value: false,
			},
		}
		return option.RegisterFlags(cmd, flags)
	}

	if err := registerFlags(cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func init() {
	// FIXME: fs
	assetUpdateCmd, err := newAssetUpdateCmd(nil)
	if err != nil {
		panic(err)
	}
	assetCmd.AddCommand(assetUpdateCmd)
}
