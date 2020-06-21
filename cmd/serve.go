package cmd

import (
	"fmt"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/mpppk/imagine/registry"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

func newServeCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// FIXME
			db, err := bolt.Open("test.db", 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open DB: %w", err)
			}
			conf, err := option.NewServeCmdConfigFromViper()
			if err != nil {
				return err
			}
			e := registry.InitializeServer(nil, db)
			e.Logger.Fatal(e.Start(":" + conf.Port))
			return nil
		},
	}
	if err := registerServeCommandFlags(cmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func registerServeCommandFlags(cmd *cobra.Command) error {
	flags := []option.Flag{
		&option.Uint16Flag{
			BaseFlag: &option.BaseFlag{
				Name:  "port",
				Usage: "server port",
			},
			Value: 1323,
		},
	}

	if err := viper.BindEnv("port"); err != nil {
		return err
	}
	return option.RegisterFlags(cmd, flags)
}

func init() {
	cmdGenerators = append(cmdGenerators, newServeCmd)
}
