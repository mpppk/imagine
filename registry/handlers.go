//+build !wireinject

package registry

import (
	"github.com/mpppk/imagine/action"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
	"go.etcd.io/bbolt"
)

func NewHandlers(db *bbolt.DB) *fsa.Handlers {
	handlers := fsa.NewHandlers()
	handlers.Handle(action.IndexClickAddDirectoryButtonType, InitializeDirectoryScanHandler(db))
	handlers.Handle(action.GlobalRequestWorkSpacesType, InitializeRequestWorkSpacesHandler(db))
	handlers.Handle(action.AssetRequestAssetsType, InitializeRequestAssetsHandler(db))
	return handlers
}
