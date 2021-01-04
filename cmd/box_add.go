package cmd

import (
	"fmt"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newBoxAddCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add bounding boxes",
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

			if err := usecases.Asset.ImportBoundingBoxesFromReader(conf.WorkSpace, cmd.InOrStdin()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd, nil
}

func init() {
	boxSubCmdGenerator = append(boxSubCmdGenerator, newBoxAddCmd)
}
