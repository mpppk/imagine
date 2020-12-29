package cmd

import (
	"fmt"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

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
			usecases, err := registry.NewBoltUseCasesWithDBPath(conf.DB)
			if err != nil {
				return fmt.Errorf("failed to create usecases instance: %w", err)
			}
			defer func() {
				if err := usecases.Close(); err != nil {
					panic(err)
				}
			}()

			fCh, errCh, err := usecases.Asset.ListAsyncWithFormat(conf.WorkSpace, conf.Format, 100)
			if err != nil {
				return err
			}
			for {
				select {
				case formattedAsset, ok := <-fCh:
					if !ok {
						return nil
					}
					cmd.Println(formattedAsset)
				case err := <-errCh:
					return err
				}
			}
		},
	}

	registerFlags := func(cmd *cobra.Command) error {
		flags := []option.Flag{
			&option.StringFlag{
				BaseFlag: &option.BaseFlag{
					Name:  "format",
					Usage: "output format",
				},
				Value: "json",
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
	assetListCmd, err := newAssetListCmd(nil)
	if err != nil {
		panic(err)
	}
	assetCmd.AddCommand(assetListCmd)
}
