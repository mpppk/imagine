package cmd

import (
	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// FIXME: fs
var assetCmd, _ = newAssetCmd(nil)

func newAssetCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "asset",
		Short: "Manage assets",
	}

	registerFlags := func(cmd *cobra.Command) error {
		flags := []option.Flag{
			&option.StringFlag{
				BaseFlag: &option.BaseFlag{
					Name:         "workspace",
					Usage:        "workspace name",
					IsPersistent: true,
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
	rootCmd.AddCommand(assetCmd)
}
