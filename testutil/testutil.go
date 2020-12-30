package testutil

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func setCmdOut(cmd *cobra.Command, newOut io.Writer) {
	cmd.SetOut(newOut)
	cmd.SetErr(newOut)
	for _, c := range cmd.Commands() {
		setCmdOut(c, newOut)
	}
}

func ExecuteCommand(t *testing.T, cmd *cobra.Command, command string) string {
	buf := new(bytes.Buffer)
	setCmdOut(cmd, buf)
	cmdArgs := strings.Split(command, " ")
	cmd.SetArgs(cmdArgs)
	if err := cmd.Execute(); err != nil {
		t.Errorf("failed to execute box add command: %v", err)
	}

	return buf.String()
}

func Diff(t *testing.T, got, want interface{}) {
	t.Helper()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("(-got +want)\n%s", diff)
	}
}

func DeepEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
}
