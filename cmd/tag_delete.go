package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newTagDeleteCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewTagDeleteCmdConfigFromViper(args)
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

			deletedTags, err := usecases.Tag.Delete(conf.WorkSpace, conf.Queries)
			if err != nil {
				return err
			}
			for _, tag := range deletedTags {
				contents, err := json.Marshal(tag)
				if err != nil {
					return fmt.Errorf("failed to marshal tag to json: %w", err)
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
					Name:  "query",
					Usage: "Only tags that match query will be deleted",
				},
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
	tagSubCmdGenerator = append(tagSubCmdGenerator, newTagDeleteCmd)
}
