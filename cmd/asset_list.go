package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mpppk/imagine/infra/repoimpl"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

	"github.com/spf13/cobra"
)

func boxesToTagIDList(boxes []*model.BoundingBox) (idList []model.TagID) {
	tagM := map[model.TagID]struct{}{}
	for _, box := range boxes {
		tagM[box.TagID] = struct{}{}
	}

	for id := range tagM {
		idList = append(idList, id)
	}
	return
}

func newAssetListCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewAssetListCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, nil)
			if err != nil {
				return err
			}
			assetUseCase := registry.InitializeAssetUseCase(db)
			assetChan, err := assetUseCase.ListAsync(context.Background(), conf.WorkSpace)
			if err != nil {
				return err
			}

			tagRepository := repoimpl.NewBBoltTag(db)
			tagSet, err := tagRepository.ListAsSet(conf.WorkSpace)
			if err != nil {
				return err
			}

			format := func(format string, asset *model.Asset) (string, error) {
				switch format {
				case "json":
					contents, err := json.Marshal(asset)
					if err != nil {
						return "", fmt.Errorf("failed to marshal asset to json: %w", err)
					}
					return string(contents), nil
				case "csv":
					var tagNames []string
					for _, tagID := range boxesToTagIDList(asset.BoundingBoxes) {
						tag, ok := tagSet.Get(tagID)
						if !ok {
							log.Printf("warning: tag not found. id:%v", tagID)
							continue
						}
						tagNames = append(tagNames, tag.Name)
					}

					line := []string{
						strconv.Quote(strconv.Itoa(int(asset.ID))),
						strconv.Quote(asset.Path),
						strconv.Quote(strings.Join(tagNames, ",")),
					}

					return strings.Join(line, ","), nil
				default:
					return "", fmt.Errorf("unknown output format: %s", format)
				}
			}

			if conf.Format == "csv" {
				header := []string{strconv.Quote("id"), strconv.Quote("path"), strconv.Quote("tags")}
				cmd.Println(strings.Join(header, ","))
			}

			for asset := range assetChan {
				t, err := format(conf.Format, asset)
				if err != nil {
					return err
				}
				cmd.Println(t)
			}
			return nil
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
