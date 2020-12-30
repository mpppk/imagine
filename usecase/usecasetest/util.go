package usecasetest

import (
	"os"
	"testing"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/usecase"
)

type UseCaseUser struct {
	t      *testing.T
	dbPath string
	wsName model.WSName
}

func (u *UseCaseUser) Use(f func(u *UseCases)) {
	usecases, closer, _ := SetUpTestUseCases(u.t, u.dbPath, u.wsName)
	defer closer()
	f(usecases)
}

func (u *UseCaseUser) RemoveDB() {
	_, _, remover := SetUpTestUseCases(u.t, u.dbPath, u.wsName)
	remover()
}

func NewTestUseCaseUser(t *testing.T, dbPath string, wsName model.WSName) *UseCaseUser {
	return &UseCaseUser{
		t:      t,
		dbPath: dbPath,
		wsName: wsName,
	}
}

func SetUpTestUseCases(t *testing.T, dbPath string, wsName model.WSName) (u *UseCases, closer func(), remover func()) {
	usecases, closer, remover := SetUpUseCases(t, dbPath, wsName)
	return NewUseCases(t, usecases), closer, remover
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
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("failed to remove test file: %v", err)
		}
	}

	closer = func() {
		if err := usecases.Close(); err != nil {
			t.Fatalf("failed to close usecases: %v", err)
		}
	}

	return usecases, closer, remover
}
