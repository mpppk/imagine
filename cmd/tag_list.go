package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/registry"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newTagListCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewTagListCmdConfigFromViper(args)
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

			tags, err := usecases.Tag.List(conf.WorkSpace)
			if err != nil {
				return err
			}
			for _, tag := range tags {
				contents, err := json.Marshal(tag)
				if err != nil {
					return fmt.Errorf("failed to marshal tag to json: %w", err)
				}
				cmd.Println(string(contents))
			}
			return nil
		},
	}

	return cmd, nil
}

func init() {
	tagSubCmdGenerator = append(tagSubCmdGenerator, newTagListCmd)
}
