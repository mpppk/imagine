package cmd

import (
	"fmt"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func newValidateCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "revalidate",
		Short: "Revalidate",
		//Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRevalidateCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			assetRepository := repoimpl.NewBBoltAsset(db)
			if err := assetRepository.Revalidate(conf.WorkSpace); err != nil {
				return fmt.Errorf("failed to revalidate: %w", err)
			}
			return nil
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newValidateCmd)
}
