package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

	"github.com/spf13/cobra"
)

func newExportCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export tags",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("db", cmd.Flags().Lookup("db")); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewExportCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			assetRepository := repoimpl.NewBBoltAsset(db)
			assetUseCase := usecase.NewAsset(assetRepository)
			assetChan, err := assetUseCase.ListAsync(conf.WorkSpace)
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

	registerFlags := func(cmd *cobra.Command) error {
		flags := []option.Flag{
			&option.StringFlag{
				BaseFlag: &option.BaseFlag{
					Name:       "db",
					Usage:      "db file path",
					IsRequired: true,
				},
			},
			&option.StringFlag{
				BaseFlag: &option.BaseFlag{
					Name:  "workspace",
					Usage: "workspace name",
				},
				Value: "default-workspace",
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
	cmdGenerators = append(cmdGenerators, newExportCmd)
}
