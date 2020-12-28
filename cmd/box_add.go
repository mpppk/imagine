package cmd

import (
	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

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
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			defer func() {
				if err := db.Close(); err != nil {
					panic(err)
				}
			}()

			usecases := registry.NewBoltUseCases(db)
			if err := usecases.Asset.ImportBoundingBoxesFromReader(conf.WorkSpace, cmd.InOrStdin()); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd, nil
}

func init() {
	// FIXME: fs
	boxAddCmd, err := newBoxAddCmd(nil)
	if err != nil {
		panic(err)
	}
	boxCmd.AddCommand(boxAddCmd)
}
