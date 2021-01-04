package testutil

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func ExecuteCommand(t *testing.T, cmd *cobra.Command, command, in string) string {
	if in != "" {
		cmd.SetIn(strings.NewReader(in))
	}

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmdArgs := strings.Split(command, " ")
	cmd.SetArgs(cmdArgs)
	if err := cmd.Execute(); err != nil {
		t.Errorf("failed to execute box add command: %v", err)
	}

	return buf.String()
}

// Diff tests want and got are equals.
// If there are difference, fail test and print diff.
// See https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff
func Diff(t *testing.T, want, got interface{}) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got)\n%s", diff)
	}
}

func DeepEqual(t *testing.T, want, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
}
