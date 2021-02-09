package usecasetest

import (
	"os"
	"testing"

	"github.com/mpppk/imagine/usecase/interactor"

	"github.com/mpppk/imagine/testutil"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/registry"
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

// SetUpUseCases setup Usecases instance and cleanup function
func SetUpUseCases(t *testing.T, dbPath string, wsName model.WSName) (u *interactor.UseCases, closer func(), remover func()) {
	t.Helper()
	usecases, err := registry.NewBoltUseCasesWithDBPath(dbPath)
	if err != nil {
		t.Fatalf("failed to create Usecases instance: %v", err)
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
		if err := usecases.Close(); err != nil {
			t.Fatalf("failed to close Usecases: %v", err)
		}
		if err := os.Remove(dbPath); err != nil {
			t.Errorf("failed to remove test file: %v", err)
		}
	}

	closer = func() {
		t.Helper()
		if err := usecases.Close(); err != nil {
			t.Fatalf("failed to close Usecases: %v", err)
		}
	}

	return usecases, closer, remover
}

func RunParallelWithUseCases(t *testing.T, name string, wsName model.WSName, f func(t *testing.T, ut *UseCases)) {
	f2 := func(t *testing.T) {
		t.Parallel()
		u := NewTestUseCaseUser(t, wsName)
		defer u.RemoveDB()
		u.Use(func(ut *UseCases) {
			f(t, ut)
		})
	}
	t.Run(name, f2)
}
