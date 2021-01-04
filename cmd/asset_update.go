package cmd

import (
	"fmt"
	"os"

	"github.com/mpppk/imagine/registry"

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
			usecases, err := registry.NewBoltUseCasesWithDBPath(conf.DB)
			if err != nil {
				return fmt.Errorf("failed to create usecases instance: %w", err)
			}
			defer func() {
				if err := usecases.Close(); err != nil {
					panic(err)
				}
			}()
			if err := usecases.InitializeWorkSpace(conf.WorkSpace); err != nil {
				return fmt.Errorf("failed to initialize asset usecase: %w", err)
			}

			// FIXME: capacity
			if err := usecases.Asset.AddOrUpdateImportAssetsFromReader(conf.WorkSpace, os.Stdin, 10000); err != nil {
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
	assetSubCmdGenerator = append(assetSubCmdGenerator, newAssetUpdateCmd)
}
