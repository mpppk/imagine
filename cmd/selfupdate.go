package cmd

import (
	"github.com/mpppk/imagine/util"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newSelfUpdateCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "selfupdate",
		Short: "Update imagine",
		//Long: `Update imagine`,
		Run: func(cmd *cobra.Command, args []string) {
			updated, err := util.Do()
			if err != nil {
				cmd.Println("Binary update failed:", err)
				return
			}
			if updated {
				cmd.Println("Current binary is the latest version", util.Version)
			} else {
				cmd.Println("Successfully updated to version", util.Version)
			}
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newSelfUpdateCmd)
}
