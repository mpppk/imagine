package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

	"github.com/spf13/cobra"
)

func newAssetDeleteCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete assets",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("db", cmd.Flags().Lookup("db")); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewAssetDeleteCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			assetRepository := repoimpl.NewBBoltAsset(db)
			scanner := bufio.NewScanner(os.Stdin)
			var asset model.Asset
			for scanner.Scan() {
				if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
					return fmt.Errorf("failed to unmarshal json to asset")
				}
				if asset.ID == 0 {
					log.Printf("warning: missing ID: %#v", asset)
					continue
				}

				if ok, err := assetRepository.Has(conf.WorkSpace, asset.ID); err != nil {
					return fmt.Errorf("failed to check asset. image path: %s: %w", asset.Path, err)
				} else if !ok {
					log.Printf("debug: asset skipped because it does not exist: id:%d", asset.ID)
					continue
				}

				if err := assetRepository.Delete(conf.WorkSpace, asset.ID); err != nil {
					return fmt.Errorf("failed to add asset. image path: %s: %w", asset.Path, err)
				}
				log.Printf("debug: asset deleted: %#v", asset)
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("faield to scan asset op: %w", err)
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
	// FIXME: fs
	assetListCmd, err := newAssetDeleteCmd(nil)
	if err != nil {
		panic(err)
	}
	assetCmd.AddCommand(assetListCmd)
}
