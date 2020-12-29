package cmd

import (
	"fmt"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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
			usecases, err := registry.NewBoltUseCasesWithDBPath(conf.DB)
			if err != nil {
				return fmt.Errorf("failed to create usecases instance: %w", err)
			}
			defer func() {
				if err := usecases.Close(); err != nil {
					panic(err)
				}
			}()
			if err := usecases.Client.Asset.Revalidate(conf.WorkSpace); err != nil {
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
