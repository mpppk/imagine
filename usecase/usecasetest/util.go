package usecasetest

import (
	"github.com/mpppk/imagine/testutil"
	"os"
	"testing"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/usecase"
)

type UseCaseUser struct {
	t        *testing.T
	DBPath   string
	wsName   model.WSName
	RemoveDB func()
}

func (u *UseCaseUser) Use(f func(u *UseCases)) {
	usecases, closer, _ := SetUpTestUseCases(u.t, u.DBPath, u.wsName)
	defer closer()
	f(usecases)
}

func NewTestUseCaseUser(t *testing.T, wsName model.WSName) *UseCaseUser {
	file, closeF, removeF := testutil.NewTempDBFile(t)
	defer closeF()

	return &UseCaseUser{
		t:        t,
		DBPath:   file.Name(),
		wsName:   wsName,
		RemoveDB: removeF,
	}
}

func SetUpTestUseCases(t *testing.T, dbPath string, wsName model.WSName) (u *UseCases, closer func(), remover func()) {
	usecases, closer, remover := SetUpUseCases(t, dbPath, wsName)
	return NewUseCases(t, usecases), closer, remover
}

func SetUpUseCasesWithTempDB(t *testing.T, wsName model.WSName) (u *usecase.UseCases, closer func(), remover func()) {
	file, closeF, removeF := testutil.NewTempDBFile(t)
	closeF()
	usecases, closeF, _ := SetUpUseCases(t, file.Name(), wsName)
	return usecases, closeF, removeF
}

// SetUpUseCases setup usecases instance and cleanup function
func SetUpUseCases(t *testing.T, dbPath string, wsName model.WSName) (u *usecase.UseCases, closer func(), remover func()) {
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

	remover = func() {
		t.Helper()
		if err := usecases.Close();err != nil {
			t.Fatalf("failed to close usecases: %v", err)
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("failed to remove test file: %v", err)
		}
	}

	closer = func() {
		t.Helper()
		if err := usecases.Close(); err != nil {
			t.Fatalf("failed to close usecases: %v", err)
		}
	}

	return usecases, closer, remover
}
