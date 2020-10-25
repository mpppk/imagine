package cmd

import (
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// FIXME: fs
var assetCmd, _ = newAssetCmd(nil)

func newAssetCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "asset",
	}

	return cmd, nil
}

func init() {
	rootCmd.AddCommand(assetCmd)
}
