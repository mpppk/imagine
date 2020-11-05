package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/mpppk/imagine/infra/repoimpl"
	"github.com/mpppk/imagine/usecase"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

	"github.com/spf13/cobra"
)

func newTagListCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: tag一覧表示を実装
			conf, err := option.NewTagListCmdConfigFromViper(args)
			if err != nil {
				return err
			}
			db, err := bolt.Open(conf.DB, 0600, &bolt.Options{ReadOnly: true})
			if err != nil {
				return err
			}
			tagRepository := repoimpl.NewBBoltTag(db)
			tagUseCase := usecase.NewTag(tagRepository)
			tags, err := tagUseCase.List(conf.WorkSpace)
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
	// FIXME: fs
	tagListCmd, err := newTagListCmd(nil)
	if err != nil {
		panic(err)
	}
	tagCmd.AddCommand(tagListCmd)
}
