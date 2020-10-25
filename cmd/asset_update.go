package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	bolt "go.etcd.io/bbolt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newAssetUpdateCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update assets",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("db", cmd.Flags().Lookup("db")); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewAssetAddCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			assetRepository := repoimpl.NewBBoltAsset(db)
			assetUseCase := usecase.NewAsset(assetRepository)

			scanner := bufio.NewScanner(os.Stdin)
			var asset model.Asset
			for scanner.Scan() {
				if err := json.Unmarshal(scanner.Bytes(), &asset); err != nil {
					return fmt.Errorf("failed to unmarshal json to asset")
				}
				if asset.ID == 0 {
					if asset.Path == "" {
						log.Printf("warning: image path is empty: %s", scanner.Text())
						continue
					}
					if err := assetUseCase.AddAssetFromImagePath(conf.WorkSpace, asset.Path); err != nil {
						return fmt.Errorf("failed to add asset. image path: %s: %w", asset.Path, err)
					}
					log.Printf("debug: asset added: %#v", asset)
				} else {
					ok, err := assetRepository.Has(conf.WorkSpace, asset.ID)
					if err != nil {
						return fmt.Errorf("failed to check asset. image path: %s: %w", asset.Path, err)
					} else if !ok {
						if conf.New {
							if err := assetRepository.Add(conf.WorkSpace, &asset); err != nil {
								return fmt.Errorf("failed to add asset: %w", err)
							}
							log.Printf("debug: asset added: %#v", asset)
						} else {
							log.Printf("debug: asset skipped because it does not exist: id:%d", asset.ID)
						}
						continue
					}
					if err := assetRepository.Update(conf.WorkSpace, &asset); err != nil {
						return fmt.Errorf("failed to update asset: %w", err)
					}
					log.Printf("debug: asset updated: %#v", asset)
				}
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
