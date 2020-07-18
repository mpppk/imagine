package registry

import (
	"github.com/mpppk/imagine/action"
	fsa "github.com/mpppk/lorca-fsa"
	"go.etcd.io/bbolt"
)

func NewHandlers(db *bbolt.DB) *fsa.Handlers {
	handlers := fsa.NewHandlers()
	handlers.Handle(action.IndexClickAddDirectoryButtonType, InitializeDirectoryScanHandler(db))
	return handlers
}
