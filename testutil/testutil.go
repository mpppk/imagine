package testutil

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/spf13/cobra"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/usecase"
)

// SetUpUseCases setup usecases instance and cleanup function
func SetUpUseCases(t *testing.T, dbPath string, wsName model.WSName) (*usecase.UseCases, func()) {
	t.Helper()
	usecases, err := registry.NewBoltUseCasesWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("failed to create usecases instance: %v", err)
	}

	if err := usecases.Client.Init(); err != nil {
		t.Fatalf("failed to initialize client: %v", err)
	}

	if wsName != "" {
		if err := usecases.InitializeWorkSpace(wsName); err != nil {
			t.Fatalf("failed to initialize workspace(%s): %v", wsName, err)
		}
	}

	f := func() {
		if err := usecases.Close(); err != nil {
			t.Fatalf("failed to close usecases: %v", err)
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("failed to remove test file: %v", err)
		}
	}
	return usecases, f
}

func ExecuteCommand(cmd *cobra.Command, command string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmdArgs := strings.Split(command, " ")
	cmd.SetArgs(cmdArgs)
	if err := cmd.Execute(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Diff(t *testing.T, got, want interface{}) {
	t.Helper()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("(-got +want)\n%s", diff)
	}
}
