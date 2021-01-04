package cmd

import (
	"fmt"

	"github.com/mpppk/imagine/cmd/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

var boxSubCmdGenerator []cmdGenerator

func newBoxCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "box",
		Short: "Manage bounding boxes",
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

	if err := registerSubCommands(fs, cmd, boxSubCmdGenerator); err != nil {
		return nil, fmt.Errorf("failed to rergister asset sub commands: %w", err)
	}

	return cmd, nil
}

func init() {
	rootSubCmdGenerator = append(rootSubCmdGenerator, newBoxCmd)
}
