package testutil

import (
	"os"
	"testing"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/usecase"
)

// SetUpUseCases setup usecases instance and cleanup function
func SetUpUseCases(t *testing.T, fileName string, wsName model.WSName) (*usecase.UseCases, func()) {
	t.Helper()
	usecases, err := registry.NewBoltUseCasesWithDBPath(fileName)
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
			panic(err)
		}
		if err := os.Remove(fileName); err != nil {
			t.Errorf("failed to remove test file: %v", err)
		}
	}
	return usecases, f
}
