package cmd_test

import (
	"testing"

	"github.com/mpppk/imagine/cmd"

	"github.com/spf13/cobra"
)

func newRootCmd(t *testing.T) *cobra.Command {
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		t.Fatalf("failed to create root cmd: %v", err)
	}
	return rootCmd
}
