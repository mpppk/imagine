//+build !wireinject

package registry

import (
	"github.com/mpppk/imagine/action"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
	"go.etcd.io/bbolt"
)

func NewHandlers(db *bbolt.DB) *fsa.Handlers {
	actionHandlers := InitializeHandlerCreator(db)
	handlers := fsa.NewHandlers()
	handlers.Handle(action.IndexClickAddDirectoryButtonType, actionHandlers.NewDirectoryScanHandler())
	handlers.Handle(action.GlobalRequestWorkSpacesType, actionHandlers.NewRequestWorkSpacesHandler())
	handlers.Handle(action.AssetRequestAssetsType, actionHandlers.NewRequestAssetsHandler())
	handlers.Handle(action.TagRequestType, actionHandlers.NewRequestAssetsHandler())
	return handlers
}
